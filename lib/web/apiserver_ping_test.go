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

package web

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/gravitational/roundtrip"
	"github.com/gravitational/trace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gravitational/teleport/api"
	"github.com/gravitational/teleport/api/client/webclient"
	"github.com/gravitational/teleport/api/constants"
	autoupdatev1pb "github.com/gravitational/teleport/api/gen/proto/go/teleport/autoupdate/v1"
	headerv1 "github.com/gravitational/teleport/api/gen/proto/go/teleport/header/v1"
	"github.com/gravitational/teleport/api/types"
	"github.com/gravitational/teleport/api/types/autoupdate"
	"github.com/gravitational/teleport/lib/client"
	"github.com/gravitational/teleport/lib/modules"
	"github.com/gravitational/teleport/lib/modules/modulestest"
)

func TestPing(t *testing.T) {
	ctx := context.Background()
	env := newWebPack(t, 1)
	authServer := env.server.Auth()

	clt, err := client.NewWebClient(env.proxies[0].webURL.String(), roundtrip.HTTPClient(client.NewInsecureWebClient()))
	require.NoError(t, err)

	tests := []struct {
		name       string
		buildType  string // defaults to modules.BuildOSS
		spec       *types.AuthPreferenceSpecV2
		assertResp func(cap types.AuthPreference, resp *webclient.PingResponse)
	}{
		{
			name: "OK local auth",
			spec: &types.AuthPreferenceSpecV2{
				Type:         constants.Local,
				SecondFactor: constants.SecondFactorOptional,
				U2F: &types.U2F{
					AppID: "https://example.com",
				},
				Webauthn: &types.Webauthn{
					RPID: "example.com",
				},
			},
			assertResp: func(cap types.AuthPreference, resp *webclient.PingResponse) {
				assert.Equal(t, cap.GetType(), resp.Auth.Type)
				assert.Equal(t, constants.SecondFactorOn, resp.Auth.SecondFactor)
				assert.NotEmpty(t, cap.GetPreferredLocalMFA(), "preferred local MFA empty")
				assert.NotNil(t, resp.Auth.Local, "Auth.Local expected")

				u2f, _ := cap.GetU2F()
				require.NotNil(t, resp.Auth.U2F)
				assert.Equal(t, u2f.AppID, resp.Auth.U2F.AppID)

				webCfg, _ := cap.GetWebauthn()
				require.NotNil(t, resp.Auth.Webauthn)
				assert.Equal(t, webCfg.RPID, resp.Auth.Webauthn.RPID)

				assert.Equal(t, types.SignatureAlgorithmSuite_SIGNATURE_ALGORITHM_SUITE_UNSPECIFIED, resp.Auth.SignatureAlgorithmSuite)
			},
		},
		{
			name: "OK local auth SecondFactors",
			spec: &types.AuthPreferenceSpecV2{
				Type: constants.Local,
				SecondFactors: []types.SecondFactorType{
					types.SecondFactorType_SECOND_FACTOR_TYPE_WEBAUTHN,
				},
				U2F: &types.U2F{
					AppID: "https://example.com",
				},
				Webauthn: &types.Webauthn{
					RPID: "example.com",
				},
			},
			assertResp: func(cap types.AuthPreference, resp *webclient.PingResponse) {
				assert.Equal(t, cap.GetType(), resp.Auth.Type)
				assert.Equal(t, constants.SecondFactorWebauthn, resp.Auth.SecondFactor)
				assert.NotEmpty(t, cap.GetPreferredLocalMFA(), "preferred local MFA empty")
				assert.NotNil(t, resp.Auth.Local, "Auth.Local expected")

				u2f, _ := cap.GetU2F()
				require.NotNil(t, resp.Auth.U2F)
				assert.Equal(t, u2f.AppID, resp.Auth.U2F.AppID)

				webCfg, _ := cap.GetWebauthn()
				require.NotNil(t, resp.Auth.Webauthn)
				assert.Equal(t, webCfg.RPID, resp.Auth.Webauthn.RPID)

				assert.Equal(t, types.SignatureAlgorithmSuite_SIGNATURE_ALGORITHM_SUITE_UNSPECIFIED, resp.Auth.SignatureAlgorithmSuite)
			},
		},
		{
			name: "OK signature algorithm suite",
			spec: &types.AuthPreferenceSpecV2{
				SignatureAlgorithmSuite: types.SignatureAlgorithmSuite_SIGNATURE_ALGORITHM_SUITE_BALANCED_V1,
			},
			assertResp: func(cap types.AuthPreference, resp *webclient.PingResponse) {
				assert.Equal(t, types.SignatureAlgorithmSuite_SIGNATURE_ALGORITHM_SUITE_BALANCED_V1, resp.Auth.SignatureAlgorithmSuite)
			},
		},
		{
			name: "OK passwordless connector",
			spec: &types.AuthPreferenceSpecV2{
				Type:         constants.Local,
				SecondFactor: constants.SecondFactorOptional,
				Webauthn: &types.Webauthn{
					RPID: "example.com",
				},
				ConnectorName: constants.PasswordlessConnector,
			},
			assertResp: func(_ types.AuthPreference, resp *webclient.PingResponse) {
				assert.True(t, resp.Auth.AllowPasswordless, "Auth.AllowPasswordless")
				require.NotNil(t, resp.Auth.Local, "Auth.Local")
				assert.Equal(t, constants.PasswordlessConnector, resp.Auth.Local.Name, "Auth.Local.Name")
			},
		},
		{
			name: "OK headless connector",
			spec: &types.AuthPreferenceSpecV2{
				Type:         constants.Local,
				SecondFactor: constants.SecondFactorOptional,
				Webauthn: &types.Webauthn{
					RPID: "example.com",
				},
				ConnectorName: constants.HeadlessConnector,
			},
			assertResp: func(_ types.AuthPreference, resp *webclient.PingResponse) {
				assert.True(t, resp.Auth.AllowHeadless, "Auth.AllowHeadless")
				require.NotNil(t, resp.Auth.Local, "Auth.Local")
				assert.Equal(t, constants.HeadlessConnector, resp.Auth.Local.Name, "Auth.Local.Name")
			},
		},
		{
			name:      "OK device trust mode=off",
			buildType: modules.BuildOSS,
			spec: &types.AuthPreferenceSpecV2{
				// Configuration is unimportant, what counts here is that the build
				// is OSS.
				Type:         constants.Local,
				SecondFactor: constants.SecondFactorOptional,
				Webauthn: &types.Webauthn{
					RPID: "example.com",
				},
			},
			assertResp: func(_ types.AuthPreference, resp *webclient.PingResponse) {
				assert.True(t, resp.Auth.DeviceTrust.Disabled, "Auth.DeviceTrust.Disabled")
			},
		},
		{
			name:      "OK device trust mode=optional",
			buildType: modules.BuildEnterprise,
			spec: &types.AuthPreferenceSpecV2{
				Type:         constants.Local,
				SecondFactor: constants.SecondFactorOptional,
				Webauthn: &types.Webauthn{
					RPID: "example.com",
				},
				DeviceTrust: &types.DeviceTrust{
					Mode: constants.DeviceTrustModeOptional,
				},
			},
			assertResp: func(_ types.AuthPreference, resp *webclient.PingResponse) {
				assert.False(t, resp.Auth.DeviceTrust.Disabled, "Auth.DeviceTrust.Disabled")
			},
		},
		{
			name:      "OK device trust auto-enroll",
			buildType: modules.BuildEnterprise,
			spec: &types.AuthPreferenceSpecV2{
				Type:         constants.Local,
				SecondFactor: constants.SecondFactorOptional,
				Webauthn: &types.Webauthn{
					RPID: "example.com",
				},
				DeviceTrust: &types.DeviceTrust{
					Mode:       constants.DeviceTrustModeOptional,
					AutoEnroll: true,
				},
			},
			assertResp: func(_ types.AuthPreference, resp *webclient.PingResponse) {
				assert.False(t, resp.Auth.DeviceTrust.Disabled, "Auth.DeviceTrust.Disabled")
				assert.True(t, resp.Auth.DeviceTrust.AutoEnroll, "Auth.DeviceTrust.AutoEnroll")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buildType := test.buildType
			if buildType == "" {
				buildType = modules.BuildOSS
			}
			modulestest.SetTestModules(t, modulestest.Modules{
				TestBuildType: buildType,
			})

			cap, err := types.NewAuthPreference(*test.spec)
			require.NoError(t, err)
			cap, err = authServer.UpsertAuthPreference(ctx, cap)
			require.NoError(t, err)

			resp, err := clt.Get(ctx, clt.Endpoint("webapi", "ping"), url.Values{})
			require.NoError(t, err)
			var pingResp webclient.PingResponse
			require.NoError(t, json.Unmarshal(resp.Bytes(), &pingResp))

			test.assertResp(cap, &pingResp)
		})
	}
}

// TestPing_multiProxyAddr makes sure ping endpoint can be called over any of
// the proxy's configured public addresses.
func TestPing_multiProxyAddr(t *testing.T) {
	env := newWebPack(t, 1)
	proxy := env.proxies[0]
	req, err := http.NewRequest(http.MethodGet, proxy.newClient(t).Endpoint("webapi", "ping"), nil)
	require.NoError(t, err)
	// Make sure ping endpoint can be reached over all proxy public addrs.
	for _, proxyAddr := range proxy.handler.handler.cfg.ProxyPublicAddrs {
		req.Host = proxyAddr.Host()
		resp, err := client.NewInsecureWebClient().Do(req)
		require.NoError(t, err)
		require.NoError(t, resp.Body.Close())
	}
}

// TestPing_minimalAPI tests that pinging the minimal web API works correctly.
func TestPing_minimalAPI(t *testing.T) {
	env := newWebPack(t, 1, func(cfg *proxyConfig) {
		cfg.minimalHandler = true
	})
	proxy := env.proxies[0]
	tests := []struct {
		name string
		host string
	}{
		{
			name: "Default ping",
			host: proxy.handler.handler.cfg.ProxyPublicAddrs[0].Host(),
		},
		{
			// This test ensures that the API doesn't try to launch an application.
			name: "Ping with alternate host",
			host: "example.com",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, proxy.newClient(t).Endpoint("webapi", "ping"), nil)
			require.NoError(t, err)
			req.Host = tc.host
			resp, err := client.NewInsecureWebClient().Do(req)
			require.NoError(t, err)
			require.NoError(t, resp.Body.Close())
		})
	}
}

// TestPing_autoUpdateResources tests that find endpoint return correct data related to auto updates.
func TestPing_autoUpdateResources(t *testing.T) {
	env := newWebPack(t, 1)
	proxy := env.proxies[0]

	ctx := t.Context()

	testGroup := "test-group"
	testUpdaterID := uuid.NewString()

	// Note: we are not using webclient.Ping() here because we don't have a valid certificate for 127.0.0.1 and
	// The webclient doesn't support being passed custom transport.
	clt := proxy.newClient(t)
	tests := []struct {
		name     string
		config   *autoupdatev1pb.AutoUpdateConfigSpec
		version  *autoupdatev1pb.AutoUpdateVersionSpec
		rollout  *autoupdatev1pb.AutoUpdateAgentRollout
		cleanup  bool
		expected webclient.AutoUpdateSettings
	}{
		{
			name: "resources not defined",
			expected: webclient.AutoUpdateSettings{
				ToolsVersion:             api.Version,
				ToolsAutoUpdate:          false,
				AgentUpdateJitterSeconds: DefaultAgentUpdateJitterSeconds,
				AgentAutoUpdate:          false,
				AgentVersion:             api.Version,
			},
		},
		{
			name: "enable tools auto update",
			config: &autoupdatev1pb.AutoUpdateConfigSpec{
				Tools: &autoupdatev1pb.AutoUpdateConfigSpecTools{
					Mode: autoupdate.ToolsUpdateModeEnabled,
				},
			},
			expected: webclient.AutoUpdateSettings{
				ToolsAutoUpdate:          true,
				ToolsVersion:             api.Version,
				AgentUpdateJitterSeconds: DefaultAgentUpdateJitterSeconds,
				AgentAutoUpdate:          false,
				AgentVersion:             api.Version,
			},
			cleanup: true,
		},
		{
			name: "enable agent auto update, immediate schedule",
			rollout: &autoupdatev1pb.AutoUpdateAgentRollout{
				Metadata: &headerv1.Metadata{
					Name: types.MetaNameAutoUpdateAgentRollout,
				},
				Spec: &autoupdatev1pb.AutoUpdateAgentRolloutSpec{
					AutoupdateMode: autoupdate.AgentsUpdateModeEnabled,
					Strategy:       autoupdate.AgentsStrategyHaltOnError,
					Schedule:       autoupdate.AgentsScheduleImmediate,
					StartVersion:   "1.2.3",
					TargetVersion:  "1.2.4",
				},
			},
			expected: webclient.AutoUpdateSettings{
				ToolsVersion:             api.Version,
				ToolsAutoUpdate:          false,
				AgentUpdateJitterSeconds: DefaultAgentUpdateJitterSeconds,
				AgentAutoUpdate:          true,
				AgentVersion:             "1.2.4",
			},
			cleanup: true,
		},
		{
			name: "agent rollout present but AU mode is disabled",
			rollout: &autoupdatev1pb.AutoUpdateAgentRollout{
				Metadata: &headerv1.Metadata{
					Name: types.MetaNameAutoUpdateAgentRollout,
				},
				Spec: &autoupdatev1pb.AutoUpdateAgentRolloutSpec{
					AutoupdateMode: autoupdate.AgentsUpdateModeDisabled,
					Strategy:       autoupdate.AgentsStrategyHaltOnError,
					Schedule:       autoupdate.AgentsScheduleImmediate,
					StartVersion:   "1.2.3",
					TargetVersion:  "1.2.4",
				},
			},
			expected: webclient.AutoUpdateSettings{
				ToolsVersion:             api.Version,
				ToolsAutoUpdate:          false,
				AgentUpdateJitterSeconds: DefaultAgentUpdateJitterSeconds,
				AgentAutoUpdate:          false,
				AgentVersion:             "1.2.4",
			},
			cleanup: true,
		},
		{
			name:    "empty config and version",
			config:  &autoupdatev1pb.AutoUpdateConfigSpec{},
			version: &autoupdatev1pb.AutoUpdateVersionSpec{},
			expected: webclient.AutoUpdateSettings{
				ToolsVersion:             api.Version,
				ToolsAutoUpdate:          false,
				AgentUpdateJitterSeconds: DefaultAgentUpdateJitterSeconds,
				AgentAutoUpdate:          false,
				AgentVersion:             api.Version,
			},
			cleanup: true,
		},
		{
			name: "set tools auto update version",
			version: &autoupdatev1pb.AutoUpdateVersionSpec{
				Tools: &autoupdatev1pb.AutoUpdateVersionSpecTools{
					TargetVersion: "1.2.3",
				},
			},
			expected: webclient.AutoUpdateSettings{
				ToolsVersion:             "1.2.3",
				ToolsAutoUpdate:          false,
				AgentUpdateJitterSeconds: DefaultAgentUpdateJitterSeconds,
				AgentAutoUpdate:          false,
				AgentVersion:             api.Version,
			},
			cleanup: true,
		},
		{
			name: "enable tools auto update and set version",
			config: &autoupdatev1pb.AutoUpdateConfigSpec{
				Tools: &autoupdatev1pb.AutoUpdateConfigSpecTools{
					Mode: autoupdate.ToolsUpdateModeEnabled,
				},
			},
			version: &autoupdatev1pb.AutoUpdateVersionSpec{
				Tools: &autoupdatev1pb.AutoUpdateVersionSpecTools{
					TargetVersion: "1.2.3",
				},
			},
			expected: webclient.AutoUpdateSettings{
				ToolsAutoUpdate:          true,
				ToolsVersion:             "1.2.3",
				AgentUpdateJitterSeconds: DefaultAgentUpdateJitterSeconds,
				AgentAutoUpdate:          false,
				AgentVersion:             api.Version,
			},
		},
		{
			name: "modify auto update config and version",
			config: &autoupdatev1pb.AutoUpdateConfigSpec{
				Tools: &autoupdatev1pb.AutoUpdateConfigSpecTools{
					Mode: autoupdate.ToolsUpdateModeDisabled,
				},
			},
			version: &autoupdatev1pb.AutoUpdateVersionSpec{
				Tools: &autoupdatev1pb.AutoUpdateVersionSpecTools{
					TargetVersion: "3.2.1",
				},
			},
			expected: webclient.AutoUpdateSettings{
				ToolsAutoUpdate:          false,
				ToolsVersion:             "3.2.1",
				AgentUpdateJitterSeconds: DefaultAgentUpdateJitterSeconds,
				AgentAutoUpdate:          false,
				AgentVersion:             api.Version,
			},
			cleanup: true,
		},
		{
			name: "group must be updated",
			rollout: &autoupdatev1pb.AutoUpdateAgentRollout{
				Metadata: &headerv1.Metadata{
					Name: types.MetaNameAutoUpdateAgentRollout,
				},
				Spec: &autoupdatev1pb.AutoUpdateAgentRolloutSpec{
					AutoupdateMode: autoupdate.AgentsUpdateModeEnabled,
					Strategy:       autoupdate.AgentsStrategyHaltOnError,
					Schedule:       autoupdate.AgentsScheduleRegular,
					StartVersion:   "1.2.3",
					TargetVersion:  "1.2.4",
				},
				Status: &autoupdatev1pb.AutoUpdateAgentRolloutStatus{
					Groups: []*autoupdatev1pb.AutoUpdateAgentRolloutStatusGroup{
						{
							Name:       testGroup,
							State:      autoupdatev1pb.AutoUpdateAgentGroupState_AUTO_UPDATE_AGENT_GROUP_STATE_ACTIVE,
							ConfigDays: []string{"*"},
						},
						{
							Name:       "unstarted",
							State:      autoupdatev1pb.AutoUpdateAgentGroupState_AUTO_UPDATE_AGENT_GROUP_STATE_UNSTARTED,
							ConfigDays: []string{"*"},
						},
					},
				},
			},
			expected: webclient.AutoUpdateSettings{
				ToolsVersion:             api.Version,
				ToolsAutoUpdate:          false,
				AgentVersion:             "1.2.4",
				AgentAutoUpdate:          true,
				AgentUpdateJitterSeconds: DefaultAgentUpdateJitterSeconds,
			},
			cleanup: true,
		},
		{
			name: "group must not be updated",
			rollout: &autoupdatev1pb.AutoUpdateAgentRollout{
				Metadata: &headerv1.Metadata{
					Name: types.MetaNameAutoUpdateAgentRollout,
				},
				Spec: &autoupdatev1pb.AutoUpdateAgentRolloutSpec{
					AutoupdateMode: autoupdate.AgentsUpdateModeEnabled,
					Strategy:       autoupdate.AgentsStrategyHaltOnError,
					Schedule:       autoupdate.AgentsScheduleRegular,
					StartVersion:   "1.2.3",
					TargetVersion:  "1.2.4",
				},
				Status: &autoupdatev1pb.AutoUpdateAgentRolloutStatus{
					Groups: []*autoupdatev1pb.AutoUpdateAgentRolloutStatusGroup{
						{
							Name:       "done",
							State:      autoupdatev1pb.AutoUpdateAgentGroupState_AUTO_UPDATE_AGENT_GROUP_STATE_ACTIVE,
							ConfigDays: []string{"*"},
						},
						{
							Name:       testGroup,
							State:      autoupdatev1pb.AutoUpdateAgentGroupState_AUTO_UPDATE_AGENT_GROUP_STATE_UNSTARTED,
							ConfigDays: []string{"*"},
						},
					},
				},
			},
			expected: webclient.AutoUpdateSettings{
				ToolsVersion:             api.Version,
				ToolsAutoUpdate:          false,
				AgentVersion:             "1.2.3",
				AgentAutoUpdate:          false,
				AgentUpdateJitterSeconds: DefaultAgentUpdateJitterSeconds,
			},
			cleanup: true,
		},
		{
			name: "canary must be updated",
			rollout: &autoupdatev1pb.AutoUpdateAgentRollout{
				Metadata: &headerv1.Metadata{
					Name: types.MetaNameAutoUpdateAgentRollout,
				},
				Spec: &autoupdatev1pb.AutoUpdateAgentRolloutSpec{
					AutoupdateMode: autoupdate.AgentsUpdateModeEnabled,
					Strategy:       autoupdate.AgentsStrategyHaltOnError,
					Schedule:       autoupdate.AgentsScheduleRegular,
					StartVersion:   "1.2.3",
					TargetVersion:  "1.2.4",
				},
				Status: &autoupdatev1pb.AutoUpdateAgentRolloutStatus{
					Groups: []*autoupdatev1pb.AutoUpdateAgentRolloutStatusGroup{
						{
							Name:       testGroup,
							State:      autoupdatev1pb.AutoUpdateAgentGroupState_AUTO_UPDATE_AGENT_GROUP_STATE_CANARY,
							ConfigDays: []string{"*"},
							Canaries: []*autoupdatev1pb.Canary{
								{
									UpdaterId: testUpdaterID,
									HostId:    uuid.NewString(),
									Hostname:  "test-host",
									Success:   false,
								},
							},
						},
						{
							Name:       "unstarted",
							State:      autoupdatev1pb.AutoUpdateAgentGroupState_AUTO_UPDATE_AGENT_GROUP_STATE_UNSTARTED,
							ConfigDays: []string{"*"},
						},
					},
				},
			},
			expected: webclient.AutoUpdateSettings{
				ToolsVersion:             api.Version,
				ToolsAutoUpdate:          false,
				AgentVersion:             "1.2.4",
				AgentAutoUpdate:          true,
				AgentUpdateJitterSeconds: DefaultAgentUpdateJitterSeconds,
			},
			cleanup: true,
		},
		{
			name: "canary must not be updated",
			rollout: &autoupdatev1pb.AutoUpdateAgentRollout{
				Metadata: &headerv1.Metadata{
					Name: types.MetaNameAutoUpdateAgentRollout,
				},
				Spec: &autoupdatev1pb.AutoUpdateAgentRolloutSpec{
					AutoupdateMode: autoupdate.AgentsUpdateModeEnabled,
					Strategy:       autoupdate.AgentsStrategyHaltOnError,
					Schedule:       autoupdate.AgentsScheduleRegular,
					StartVersion:   "1.2.3",
					TargetVersion:  "1.2.4",
				},
				Status: &autoupdatev1pb.AutoUpdateAgentRolloutStatus{
					Groups: []*autoupdatev1pb.AutoUpdateAgentRolloutStatusGroup{
						{
							Name:  testGroup,
							State: autoupdatev1pb.AutoUpdateAgentGroupState_AUTO_UPDATE_AGENT_GROUP_STATE_CANARY,
							Canaries: []*autoupdatev1pb.Canary{
								{
									UpdaterId: uuid.NewString(),
									HostId:    uuid.NewString(),
									Hostname:  "test-host",
									Success:   false,
								},
							},
							ConfigDays: []string{"*"},
						},
						{
							Name:       "unstarted",
							State:      autoupdatev1pb.AutoUpdateAgentGroupState_AUTO_UPDATE_AGENT_GROUP_STATE_UNSTARTED,
							ConfigDays: []string{"*"},
						},
					},
				},
			},
			expected: webclient.AutoUpdateSettings{
				ToolsVersion:             api.Version,
				ToolsAutoUpdate:          false,
				AgentVersion:             "1.2.3",
				AgentAutoUpdate:          false,
				AgentUpdateJitterSeconds: DefaultAgentUpdateJitterSeconds,
			},
			cleanup: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.config != nil {
				config, err := autoupdate.NewAutoUpdateConfig(tc.config)
				require.NoError(t, err)
				_, err = env.server.Auth().UpsertAutoUpdateConfig(ctx, config)
				require.NoError(t, err)
			}
			if tc.version != nil {
				version, err := autoupdate.NewAutoUpdateVersion(tc.version)
				require.NoError(t, err)
				_, err = env.server.Auth().UpsertAutoUpdateVersion(ctx, version)
				require.NoError(t, err)
			}
			if tc.rollout != nil {
				_, err := env.server.Auth().UpsertAutoUpdateAgentRollout(ctx, tc.rollout)
				require.NoError(t, err)
			}

			// expire the fn cache to force the next answer to be fresh
			for _, proxy := range env.proxies {
				proxy.clock.Advance(2 * findEndpointCacheTTL)
			}

			resp, err := clt.Get(ctx, clt.Endpoint("webapi", "ping"), url.Values{
				webclient.AgentUpdateIDParameter:    []string{testUpdaterID},
				webclient.AgentUpdateGroupParameter: []string{testGroup},
			})
			require.NoError(t, err)

			pr := &webclient.PingResponse{}
			require.NoError(t, json.NewDecoder(resp.Reader()).Decode(pr))

			assert.Equal(t, tc.expected, pr.AutoUpdate)

			if tc.cleanup {
				require.NotErrorIs(t, env.server.Auth().DeleteAutoUpdateConfig(ctx), &trace.NotFoundError{})
				require.NotErrorIs(t, env.server.Auth().DeleteAutoUpdateVersion(ctx), &trace.NotFoundError{})
				require.NotErrorIs(t, env.server.Auth().DeleteAutoUpdateAgentRollout(ctx), &trace.NotFoundError{})
			}
		})
	}
}
