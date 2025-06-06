/*
 * Teleport
 * Copyright (C) 2023  Gravitational, Inc.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package proxy

import (
	"context"
	"encoding/base64"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/gravitational/trace"
	"github.com/jonboulle/clockwork"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/gravitational/teleport/api/types"
	"github.com/gravitational/teleport/api/utils/retryutils"
	"github.com/gravitational/teleport/lib/cloud"
	"github.com/gravitational/teleport/lib/cloud/awsconfig"
	"github.com/gravitational/teleport/lib/cloud/azure"
	"github.com/gravitational/teleport/lib/cloud/gcp"
	kubeutils "github.com/gravitational/teleport/lib/kube/utils"
	"github.com/gravitational/teleport/lib/labels"
	"github.com/gravitational/teleport/lib/service/servicecfg"
	"github.com/gravitational/teleport/lib/services"
)

// kubeDetails contain the cluster-related details including authentication.
type kubeDetails struct {
	kubeCreds

	// dynamicLabels is the dynamic labels executor for this cluster.
	dynamicLabels *labels.Dynamic
	// kubeCluster is the dynamic kube_cluster or a static generated from kubeconfig and that only has the name populated.
	kubeCluster types.KubeCluster
	// kubeClusterVersion is the version of the kube_cluster's related Kubernetes server.
	kubeClusterVersion *version.Info

	// rwMu is the mutex to protect the kubeCodecs, gvkSupportedResources, and rbacSupportedTypes.
	rwMu sync.RWMutex
	// kubeCodecs is the codec factory for the cluster resources.
	// The codec factory includes the default resources and the namespaced resources
	// that are supported by the cluster.
	// The codec factory is updated periodically to include the latest custom resources
	// that are added to the cluster.
	kubeCodecs *serializer.CodecFactory
	// rbacSupportedTypes is the list of supported types for RBAC for the cluster.
	// The list is updated periodically to include the latest custom resources
	// that are added to the cluster.
	rbacSupportedTypes rbacSupportedResources
	// gvkSupportedResources is the list of registered API path resources and their
	// GVK definition.
	gvkSupportedResources gvkSupportedResources
	// isClusterOffline is true if the cluster is offline.
	// An offline cluster will not be able to serve any requests until it comes back online.
	// The cluster is marked as offline if the cluster schema cannot be created
	// and the list of supported types for RBAC cannot be generated.
	isClusterOffline bool

	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
}

// clusterDetailsConfig contains the configuration for creating a proxied cluster.
type clusterDetailsConfig struct {
	// cloudClients is the cloud clients to use for dynamic clusters.
	cloudClients cloud.Clients
	// awsCloudClients provides AWS SDK clients.
	awsCloudClients AWSClientGetter
	// kubeCreds is the credentials to use for the cluster.
	kubeCreds kubeCreds
	// cluster is the cluster to create a proxied cluster for.
	cluster types.KubeCluster
	// log is the logger to use.
	log *slog.Logger
	// checker is the permissions checker to use.
	checker servicecfg.ImpersonationPermissionsChecker
	// resourceMatchers is the list of resource matchers to match the cluster against
	// to determine if we should assume the role or not for AWS.
	resourceMatchers []services.ResourceMatcher
	// clock is the clock to use.
	clock clockwork.Clock
	// component is the Kubernetes component that serves this cluster.
	component KubeServiceType
}

const (
	defaultRefreshPeriod = 5 * time.Minute
	backoffRefreshStep   = 10 * time.Second
)

// newClusterDetails creates a proxied kubeDetails structure given a dynamic cluster.
func newClusterDetails(ctx context.Context, cfg clusterDetailsConfig) (_ *kubeDetails, err error) {
	creds := cfg.kubeCreds
	if creds == nil {
		creds, err = getKubeClusterCredentials(ctx, cfg)
		if err != nil {
			return nil, trace.Wrap(err)
		}
	}

	var dynLabels *labels.Dynamic
	if len(cfg.cluster.GetDynamicLabels()) > 0 {
		dynLabels, err = labels.NewDynamic(
			ctx,
			&labels.DynamicConfig{
				Labels: cfg.cluster.GetDynamicLabels(),
				Log:    cfg.log,
			})
		if err != nil {
			return nil, trace.Wrap(err)
		}
		dynLabels.Sync()
		go dynLabels.Start()
	}

	var isClusterOffline bool
	// Create the codec factory and the list of supported types for RBAC.
	codecFactory, rbacSupportedTypes, gvkSupportedRes, err := newClusterSchemaBuilder(cfg.log, creds.getKubeClient())
	if err != nil {
		cfg.log.WarnContext(ctx, "Failed to create cluster schema, the cluster may be offline", "error", err)
		// If the cluster is offline, we will not be able to create the codec factory
		// and the list of supported types for RBAC.
		// We mark the cluster as offline and continue to create the kubeDetails but
		// the offline cluster will not be able to serve any requests until it comes back online.
		isClusterOffline = true
	}

	kubeVersion, err := creds.getKubeClient().Discovery().ServerVersion()
	if err != nil {
		cfg.log.WarnContext(ctx, "Failed to get Kubernetes cluster version, the cluster may be offline", "error", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	k := &kubeDetails{
		kubeCreds:             creds,
		dynamicLabels:         dynLabels,
		kubeCluster:           cfg.cluster,
		kubeClusterVersion:    kubeVersion,
		kubeCodecs:            codecFactory,
		rbacSupportedTypes:    rbacSupportedTypes,
		cancelFunc:            cancel,
		isClusterOffline:      isClusterOffline,
		gvkSupportedResources: gvkSupportedRes,
	}

	// If cluster is online and there's no errors, we refresh details seldom (every 5 minutes),
	// but if the cluster is offline, we try to refresh details more often to catch it getting back online earlier.
	firstPeriod := defaultRefreshPeriod
	if isClusterOffline {
		firstPeriod = backoffRefreshStep
	}
	refreshDelay, err := retryutils.NewLinear(retryutils.LinearConfig{
		First:  firstPeriod,
		Step:   backoffRefreshStep,
		Max:    defaultRefreshPeriod,
		Jitter: retryutils.SeventhJitter,
		Clock:  cfg.clock,
	})
	if err != nil {
		k.Close()
		return nil, trace.Wrap(err)
	}

	k.wg.Add(1)
	// Start the periodic update of the codec factory and the list of supported types for RBAC.
	go func() {
		defer k.wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case <-refreshDelay.After():
				codecFactory, rbacSupportedTypes, gvkSupportedResources, err := newClusterSchemaBuilder(cfg.log, creds.getKubeClient())
				if err != nil {
					// If this is first time we get an error, we reset retry mechanism so it will start trying to refresh details quicker, with linear backoff.
					if refreshDelay.First == defaultRefreshPeriod {
						refreshDelay.First = backoffRefreshStep
						refreshDelay.Reset()
					} else {
						refreshDelay.Inc()
					}
					cfg.log.ErrorContext(ctx, "Failed to update cluster schema", "error", err)
					continue
				}

				kubeVersion, err := creds.getKubeClient().Discovery().ServerVersion()
				if err != nil {
					cfg.log.WarnContext(ctx, "Failed to get Kubernetes cluster version, the cluster may be offline", "error", err)
				}

				// Restore details refresh delay to the default value, in case previously cluster was offline.
				refreshDelay.First = defaultRefreshPeriod

				k.rwMu.Lock()
				k.kubeCodecs = codecFactory
				k.rbacSupportedTypes = rbacSupportedTypes
				k.gvkSupportedResources = gvkSupportedResources
				k.isClusterOffline = false
				k.kubeClusterVersion = kubeVersion
				k.rwMu.Unlock()
			}
		}
	}()
	return k, nil
}

func (k *kubeDetails) Close() {
	// send a close signal and wait for the close to finish.
	k.cancelFunc()
	k.wg.Wait()
	if k.dynamicLabels != nil {
		k.dynamicLabels.Close()
	}
	// it is safe to call close even for static creds.
	k.kubeCreds.close()
}

// getClusterSupportedResources returns the codec factory and the list of supported types for RBAC.
func (k *kubeDetails) getClusterSupportedResources() (*serializer.CodecFactory, rbacSupportedResources, error) {
	k.rwMu.RLock()
	defer k.rwMu.RUnlock()
	// If the cluster is offline, return an error because we don't have the schema
	// for the cluster.
	if k.isClusterOffline {
		return nil, nil, trace.ConnectionProblem(nil, "kubernetes cluster %q is offline", k.kubeCluster.GetName())
	}
	return k.kubeCodecs, k.rbacSupportedTypes, nil
}

// getObjectGVK returns the default GVK (if any) registered for the specified request path.
func (k *kubeDetails) getObjectGVK(resource apiResource) *schema.GroupVersionKind {
	k.rwMu.RLock()
	defer k.rwMu.RUnlock()
	return k.gvkSupportedResources[gvkSupportedResourcesKey{
		name:     strings.Split(resource.resourceKind, "/")[0],
		apiGroup: resource.apiGroup,
		version:  resource.apiGroupVersion,
	}]
}

// getKubeClusterCredentials generates kube credentials for dynamic clusters.
func getKubeClusterCredentials(ctx context.Context, cfg clusterDetailsConfig) (kubeCreds, error) {
	switch dynCredsCfg := (dynamicCredsConfig{
		kubeCluster:      cfg.cluster,
		log:              cfg.log,
		checker:          cfg.checker,
		resourceMatchers: cfg.resourceMatchers,
		clock:            cfg.clock,
		component:        cfg.component,
	}); {
	case cfg.cluster.IsKubeconfig():
		return getStaticCredentialsFromKubeconfig(ctx, cfg.component, cfg.cluster, cfg.log, cfg.checker)
	case cfg.cluster.IsAzure():
		return getAzureCredentials(ctx, cfg.cloudClients, dynCredsCfg)
	case cfg.cluster.IsAWS():
		return getAWSCredentials(ctx, cfg.awsCloudClients, dynCredsCfg)
	case cfg.cluster.IsGCP():
		return getGCPCredentials(ctx, cfg.cloudClients, dynCredsCfg)
	default:
		return nil, trace.BadParameter("authentication method provided for cluster %q not supported", cfg.cluster.GetName())
	}
}

// getAzureCredentials creates a dynamicCreds that generates and updates the access credentials to a AKS Kubernetes cluster.
func getAzureCredentials(ctx context.Context, cloudClients cloud.Clients, cfg dynamicCredsConfig) (*dynamicKubeCreds, error) {
	// create a client that returns the credentials for kubeCluster
	cfg.client = azureRestConfigClient(cloudClients)

	creds, err := newDynamicKubeCreds(
		ctx,
		cfg,
	)
	return creds, trace.Wrap(err)
}

// azureRestConfigClient creates a dynamicCredsClient that returns credentials to a AKS cluster.
func azureRestConfigClient(cloudClients cloud.Clients) dynamicCredsClient {
	return func(ctx context.Context, cluster types.KubeCluster) (*rest.Config, time.Time, error) {
		aksClient, err := cloudClients.GetAzureKubernetesClient(cluster.GetAzureConfig().SubscriptionID)
		if err != nil {
			return nil, time.Time{}, trace.Wrap(err)
		}
		cfg, exp, err := aksClient.ClusterCredentials(ctx, azure.ClusterCredentialsConfig{
			ResourceGroup:                   cluster.GetAzureConfig().ResourceGroup,
			ResourceName:                    cluster.GetAzureConfig().ResourceName,
			TenantID:                        cluster.GetAzureConfig().TenantID,
			ImpersonationPermissionsChecker: checkImpersonationPermissions,
		})
		return cfg, exp, trace.Wrap(err)
	}
}

// getAWSCredentials creates a dynamicKubeCreds that generates and updates the access credentials to a EKS kubernetes cluster.
func getAWSCredentials(ctx context.Context, cloudClients AWSClientGetter, cfg dynamicCredsConfig) (*dynamicKubeCreds, error) {
	// create a client that returns the credentials for kubeCluster
	cfg.client = getAWSClientRestConfig(cloudClients, cfg.clock, cfg.resourceMatchers)
	creds, err := newDynamicKubeCreds(ctx, cfg)
	return creds, trace.Wrap(err)
}

// getAWSResourceMatcherToCluster returns the AWS assume role ARN and external ID for the cluster that matches the kubeCluster.
// If no match is found, nil is returned, which means that we should not attempt to assume a role.
func getAWSResourceMatcherToCluster(kubeCluster types.KubeCluster, resourceMatchers []services.ResourceMatcher) *services.ResourceMatcherAWS {
	if !kubeCluster.IsAWS() {
		return nil
	}
	for _, matcher := range resourceMatchers {
		if len(matcher.Labels) == 0 || matcher.AWS.AssumeRoleARN == "" {
			continue
		}
		if match, _, _ := services.MatchLabels(matcher.Labels, kubeCluster.GetAllLabels()); !match {
			continue
		}
		return &matcher.AWS
	}
	return nil
}

// STSPresignClient is the subset of the STS presign interface we use in fetchers.
type STSPresignClient = kubeutils.STSPresignClient

// EKSClient is the subset of the EKS Client interface we use.
type EKSClient interface {
	eks.DescribeClusterAPIClient
}

// AWSClientGetter is an interface for getting an EKS client and an STS client.
type AWSClientGetter interface {
	awsconfig.Provider
	// GetAWSEKSClient returns AWS EKS client for the specified config.
	GetAWSEKSClient(aws.Config) EKSClient
	// GetAWSSTSPresignClient returns AWS STS presign client for the specified config.
	GetAWSSTSPresignClient(aws.Config) STSPresignClient
}

// getAWSClientRestConfig creates a dynamicCredsClient that generates returns credentials to EKS clusters.
func getAWSClientRestConfig(cloudClients AWSClientGetter, clock clockwork.Clock, resourceMatchers []services.ResourceMatcher) dynamicCredsClient {
	return func(ctx context.Context, cluster types.KubeCluster) (*rest.Config, time.Time, error) {
		region := cluster.GetAWSConfig().Region
		opts := []awsconfig.OptionsFn{
			awsconfig.WithAmbientCredentials(),
		}
		if awsAssume := getAWSResourceMatcherToCluster(cluster, resourceMatchers); awsAssume != nil {
			opts = append(opts, awsconfig.WithAssumeRole(awsAssume.AssumeRoleARN, awsAssume.ExternalID))
		}

		cfg, err := cloudClients.GetConfig(ctx, region, opts...)
		if err != nil {
			return nil, time.Time{}, trace.Wrap(err)
		}

		regionalClient := cloudClients.GetAWSEKSClient(cfg)

		eksCfg, err := regionalClient.DescribeCluster(ctx, &eks.DescribeClusterInput{
			Name: aws.String(cluster.GetAWSConfig().Name),
		})
		if err != nil {
			return nil, time.Time{}, trace.Wrap(err)
		}

		ca, err := base64.StdEncoding.DecodeString(aws.ToString(eksCfg.Cluster.CertificateAuthority.Data))
		if err != nil {
			return nil, time.Time{}, trace.Wrap(err)
		}

		apiEndpoint := aws.ToString(eksCfg.Cluster.Endpoint)
		if len(apiEndpoint) == 0 {
			return nil, time.Time{}, trace.BadParameter("invalid api endpoint for cluster %q", cluster.GetAWSConfig().Name)
		}

		stsPresignClient := cloudClients.GetAWSSTSPresignClient(cfg)

		token, exp, err := kubeutils.GenAWSEKSToken(ctx, stsPresignClient, cluster.GetAWSConfig().Name, clock)
		if err != nil {
			return nil, time.Time{}, trace.Wrap(err)
		}

		return &rest.Config{
			Host:        apiEndpoint,
			BearerToken: token,
			TLSClientConfig: rest.TLSClientConfig{
				CAData: ca,
			},
		}, exp, nil
	}
}

// getStaticCredentialsFromKubeconfig loads a kubeconfig from the cluster and returns the access credentials for the cluster.
// If the config defines multiple contexts, it will pick one (the order is not guaranteed).
func getStaticCredentialsFromKubeconfig(ctx context.Context, component KubeServiceType, cluster types.KubeCluster, log *slog.Logger, checker servicecfg.ImpersonationPermissionsChecker) (*staticKubeCreds, error) {
	config, err := clientcmd.Load(cluster.GetKubeconfig())
	if err != nil {
		return nil, trace.WrapWithMessage(err, "unable to parse kubeconfig for cluster %q", cluster.GetName())
	}
	if len(config.CurrentContext) == 0 && len(config.Contexts) > 0 {
		// select the first context key as default context
		for k := range config.Contexts {
			config.CurrentContext = k
			break
		}
	}
	restConfig, err := clientcmd.NewDefaultClientConfig(*config, nil).ClientConfig()
	if err != nil {
		return nil, trace.WrapWithMessage(err, "unable to create client from kubeconfig for cluster %q", cluster.GetName())
	}

	creds, err := extractKubeCreds(ctx, component, cluster.GetName(), restConfig, log, checker)
	return creds, trace.Wrap(err)
}

// getGCPCredentials creates a dynamicKubeCreds that generates and updates the access credentials to a GKE kubernetes cluster.
func getGCPCredentials(ctx context.Context, cloudClients cloud.Clients, cfg dynamicCredsConfig) (*dynamicKubeCreds, error) {
	// create a client that returns the credentials for kubeCluster
	cfg.client = gcpRestConfigClient(cloudClients)
	creds, err := newDynamicKubeCreds(ctx, cfg)
	return creds, trace.Wrap(err)
}

// gcpRestConfigClient creates a dynamicCredsClient that returns credentials to a GKE cluster.
func gcpRestConfigClient(cloudClients cloud.Clients) dynamicCredsClient {
	return func(ctx context.Context, cluster types.KubeCluster) (*rest.Config, time.Time, error) {
		gkeClient, err := cloudClients.GetGCPGKEClient(ctx)
		if err != nil {
			return nil, time.Time{}, trace.Wrap(err)
		}
		cfg, exp, err := gkeClient.GetClusterRestConfig(ctx,
			gcp.ClusterDetails{
				ProjectID: cluster.GetGCPConfig().ProjectID,
				Location:  cluster.GetGCPConfig().Location,
				Name:      cluster.GetGCPConfig().Name,
			},
		)
		return cfg, exp, trace.Wrap(err)
	}
}
