// Copyright 2023 Gravitational, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: teleport/discoveryconfig/v1/discoveryconfig.proto

package discoveryconfigv1

import (
	v1 "github.com/gravitational/teleport/api/gen/proto/go/teleport/header/v1"
	types "github.com/gravitational/teleport/api/types"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// DiscoveryConfigState is the state of the discovery config resource.
type DiscoveryConfigState int32

const (
	DiscoveryConfigState_DISCOVERY_CONFIG_STATE_UNSPECIFIED DiscoveryConfigState = 0
	// DISCOVERY_CONFIG_STATE_RUNNING is used when the operation doesn't report
	// incidents.
	DiscoveryConfigState_DISCOVERY_CONFIG_STATE_RUNNING DiscoveryConfigState = 1
	// DISCOVERY_CONFIG_STATE_ERROR is used when the operation reports
	// incidents.
	DiscoveryConfigState_DISCOVERY_CONFIG_STATE_ERROR DiscoveryConfigState = 2
	// DISCOVERY_CONFIG_STATE_SYNCING is used when the discovery process has started but didn't finished yet.
	DiscoveryConfigState_DISCOVERY_CONFIG_STATE_SYNCING DiscoveryConfigState = 3
)

// Enum value maps for DiscoveryConfigState.
var (
	DiscoveryConfigState_name = map[int32]string{
		0: "DISCOVERY_CONFIG_STATE_UNSPECIFIED",
		1: "DISCOVERY_CONFIG_STATE_RUNNING",
		2: "DISCOVERY_CONFIG_STATE_ERROR",
		3: "DISCOVERY_CONFIG_STATE_SYNCING",
	}
	DiscoveryConfigState_value = map[string]int32{
		"DISCOVERY_CONFIG_STATE_UNSPECIFIED": 0,
		"DISCOVERY_CONFIG_STATE_RUNNING":     1,
		"DISCOVERY_CONFIG_STATE_ERROR":       2,
		"DISCOVERY_CONFIG_STATE_SYNCING":     3,
	}
)

func (x DiscoveryConfigState) Enum() *DiscoveryConfigState {
	p := new(DiscoveryConfigState)
	*p = x
	return p
}

func (x DiscoveryConfigState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DiscoveryConfigState) Descriptor() protoreflect.EnumDescriptor {
	return file_teleport_discoveryconfig_v1_discoveryconfig_proto_enumTypes[0].Descriptor()
}

func (DiscoveryConfigState) Type() protoreflect.EnumType {
	return &file_teleport_discoveryconfig_v1_discoveryconfig_proto_enumTypes[0]
}

func (x DiscoveryConfigState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DiscoveryConfigState.Descriptor instead.
func (DiscoveryConfigState) EnumDescriptor() ([]byte, []int) {
	return file_teleport_discoveryconfig_v1_discoveryconfig_proto_rawDescGZIP(), []int{0}
}

// DiscoveryConfig is a resource that has Discovery Resource Matchers and a Discovery Group.
//
// Teleport Discovery Services will load the dynamic DiscoveryConfigs whose Discovery Group matches the discovery_group defined in their configuration.
type DiscoveryConfig struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Header is the resource header.
	Header *v1.ResourceHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	// Spec is an DiscoveryConfig specification.
	Spec *DiscoveryConfigSpec `protobuf:"bytes,2,opt,name=spec,proto3" json:"spec,omitempty"`
	// Status is the resource Status
	Status        *DiscoveryConfigStatus `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DiscoveryConfig) Reset() {
	*x = DiscoveryConfig{}
	mi := &file_teleport_discoveryconfig_v1_discoveryconfig_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DiscoveryConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscoveryConfig) ProtoMessage() {}

func (x *DiscoveryConfig) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_discoveryconfig_v1_discoveryconfig_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscoveryConfig.ProtoReflect.Descriptor instead.
func (*DiscoveryConfig) Descriptor() ([]byte, []int) {
	return file_teleport_discoveryconfig_v1_discoveryconfig_proto_rawDescGZIP(), []int{0}
}

func (x *DiscoveryConfig) GetHeader() *v1.ResourceHeader {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *DiscoveryConfig) GetSpec() *DiscoveryConfigSpec {
	if x != nil {
		return x.Spec
	}
	return nil
}

func (x *DiscoveryConfig) GetStatus() *DiscoveryConfigStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

// DiscoveryConfigSpec contains properties required to create matchers to be used by discovery_service.
// Those matchers are used by discovery_service to watch for cloud resources and create them in Teleport.
type DiscoveryConfigSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// DiscoveryGroup is used by discovery_service to add extra matchers.
	// All the discovery_services that have the same discovery_group, will load the matchers of this resource.
	DiscoveryGroup string `protobuf:"bytes,1,opt,name=discovery_group,json=discoveryGroup,proto3" json:"discovery_group,omitempty"`
	// AWS is a list of AWS Matchers.
	Aws []*types.AWSMatcher `protobuf:"bytes,2,rep,name=aws,proto3" json:"aws,omitempty"`
	// Azure is a list of Azure Matchers.
	Azure []*types.AzureMatcher `protobuf:"bytes,3,rep,name=azure,proto3" json:"azure,omitempty"`
	// GCP is a list of GCP Matchers.
	Gcp []*types.GCPMatcher `protobuf:"bytes,4,rep,name=gcp,proto3" json:"gcp,omitempty"`
	// Kube is a list of Kubernetes Matchers.
	Kube []*types.KubernetesMatcher `protobuf:"bytes,5,rep,name=kube,proto3" json:"kube,omitempty"`
	// AccessGraph is the configurations for syncing Cloud accounts into Access Graph.
	AccessGraph   *types.AccessGraphSync `protobuf:"bytes,6,opt,name=access_graph,json=accessGraph,proto3" json:"access_graph,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DiscoveryConfigSpec) Reset() {
	*x = DiscoveryConfigSpec{}
	mi := &file_teleport_discoveryconfig_v1_discoveryconfig_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DiscoveryConfigSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscoveryConfigSpec) ProtoMessage() {}

func (x *DiscoveryConfigSpec) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_discoveryconfig_v1_discoveryconfig_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscoveryConfigSpec.ProtoReflect.Descriptor instead.
func (*DiscoveryConfigSpec) Descriptor() ([]byte, []int) {
	return file_teleport_discoveryconfig_v1_discoveryconfig_proto_rawDescGZIP(), []int{1}
}

func (x *DiscoveryConfigSpec) GetDiscoveryGroup() string {
	if x != nil {
		return x.DiscoveryGroup
	}
	return ""
}

func (x *DiscoveryConfigSpec) GetAws() []*types.AWSMatcher {
	if x != nil {
		return x.Aws
	}
	return nil
}

func (x *DiscoveryConfigSpec) GetAzure() []*types.AzureMatcher {
	if x != nil {
		return x.Azure
	}
	return nil
}

func (x *DiscoveryConfigSpec) GetGcp() []*types.GCPMatcher {
	if x != nil {
		return x.Gcp
	}
	return nil
}

func (x *DiscoveryConfigSpec) GetKube() []*types.KubernetesMatcher {
	if x != nil {
		return x.Kube
	}
	return nil
}

func (x *DiscoveryConfigSpec) GetAccessGraph() *types.AccessGraphSync {
	if x != nil {
		return x.AccessGraph
	}
	return nil
}

// DiscoveryConfigStatus holds dynamic information about the discovery configuration
// running status such as errors, state and count of the resources.
type DiscoveryConfigStatus struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// State reports the Discovery config state.
	State DiscoveryConfigState `protobuf:"varint,1,opt,name=state,proto3,enum=teleport.discoveryconfig.v1.DiscoveryConfigState" json:"state,omitempty"`
	// error_message holds the error message when state is DISCOVERY_CONFIG_STATE_ERROR.
	ErrorMessage *string `protobuf:"bytes,2,opt,name=error_message,json=errorMessage,proto3,oneof" json:"error_message,omitempty"`
	// discovered_resources holds the count of the discovered resources in the previous iteration.
	DiscoveredResources uint64 `protobuf:"varint,3,opt,name=discovered_resources,json=discoveredResources,proto3" json:"discovered_resources,omitempty"`
	// last_sync_time is the timestamp when the Discovery Config was last sync.
	LastSyncTime *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=last_sync_time,json=lastSyncTime,proto3" json:"last_sync_time,omitempty"`
	// IntegrationDiscoveredResources maps an integration to discovered resources summary.
	IntegrationDiscoveredResources map[string]*IntegrationDiscoveredSummary `protobuf:"bytes,6,rep,name=integration_discovered_resources,json=integrationDiscoveredResources,proto3" json:"integration_discovered_resources,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields                  protoimpl.UnknownFields
	sizeCache                      protoimpl.SizeCache
}

func (x *DiscoveryConfigStatus) Reset() {
	*x = DiscoveryConfigStatus{}
	mi := &file_teleport_discoveryconfig_v1_discoveryconfig_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DiscoveryConfigStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscoveryConfigStatus) ProtoMessage() {}

func (x *DiscoveryConfigStatus) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_discoveryconfig_v1_discoveryconfig_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscoveryConfigStatus.ProtoReflect.Descriptor instead.
func (*DiscoveryConfigStatus) Descriptor() ([]byte, []int) {
	return file_teleport_discoveryconfig_v1_discoveryconfig_proto_rawDescGZIP(), []int{2}
}

func (x *DiscoveryConfigStatus) GetState() DiscoveryConfigState {
	if x != nil {
		return x.State
	}
	return DiscoveryConfigState_DISCOVERY_CONFIG_STATE_UNSPECIFIED
}

func (x *DiscoveryConfigStatus) GetErrorMessage() string {
	if x != nil && x.ErrorMessage != nil {
		return *x.ErrorMessage
	}
	return ""
}

func (x *DiscoveryConfigStatus) GetDiscoveredResources() uint64 {
	if x != nil {
		return x.DiscoveredResources
	}
	return 0
}

func (x *DiscoveryConfigStatus) GetLastSyncTime() *timestamppb.Timestamp {
	if x != nil {
		return x.LastSyncTime
	}
	return nil
}

func (x *DiscoveryConfigStatus) GetIntegrationDiscoveredResources() map[string]*IntegrationDiscoveredSummary {
	if x != nil {
		return x.IntegrationDiscoveredResources
	}
	return nil
}

// IntegrationDiscoveredSummary contains the a summary for each resource type that was discovered.
type IntegrationDiscoveredSummary struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// AWSEC2 contains the summary for the AWS EC2 discovered instances.
	AwsEc2 *ResourcesDiscoveredSummary `protobuf:"bytes,1,opt,name=aws_ec2,json=awsEc2,proto3" json:"aws_ec2,omitempty"`
	// AWSRDS contains the summary for the AWS RDS discovered databases.
	AwsRds *ResourcesDiscoveredSummary `protobuf:"bytes,2,opt,name=aws_rds,json=awsRds,proto3" json:"aws_rds,omitempty"`
	// AWSEKS contains the summary for the AWS EKS discovered clusters.
	AwsEks        *ResourcesDiscoveredSummary `protobuf:"bytes,3,opt,name=aws_eks,json=awsEks,proto3" json:"aws_eks,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *IntegrationDiscoveredSummary) Reset() {
	*x = IntegrationDiscoveredSummary{}
	mi := &file_teleport_discoveryconfig_v1_discoveryconfig_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IntegrationDiscoveredSummary) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IntegrationDiscoveredSummary) ProtoMessage() {}

func (x *IntegrationDiscoveredSummary) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_discoveryconfig_v1_discoveryconfig_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IntegrationDiscoveredSummary.ProtoReflect.Descriptor instead.
func (*IntegrationDiscoveredSummary) Descriptor() ([]byte, []int) {
	return file_teleport_discoveryconfig_v1_discoveryconfig_proto_rawDescGZIP(), []int{3}
}

func (x *IntegrationDiscoveredSummary) GetAwsEc2() *ResourcesDiscoveredSummary {
	if x != nil {
		return x.AwsEc2
	}
	return nil
}

func (x *IntegrationDiscoveredSummary) GetAwsRds() *ResourcesDiscoveredSummary {
	if x != nil {
		return x.AwsRds
	}
	return nil
}

func (x *IntegrationDiscoveredSummary) GetAwsEks() *ResourcesDiscoveredSummary {
	if x != nil {
		return x.AwsEks
	}
	return nil
}

// ResourcesDiscoveredSummary represents the AWS resources that were discovered.
type ResourcesDiscoveredSummary struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Found holds the count of resources found.
	// After a resource is found, it starts the sync process and ends in either an enrolled or a failed resource.
	Found uint64 `protobuf:"varint,1,opt,name=found,proto3" json:"found,omitempty"`
	// Enrolled holds the count of the resources that were successfully enrolled.
	Enrolled uint64 `protobuf:"varint,2,opt,name=enrolled,proto3" json:"enrolled,omitempty"`
	// Failed holds the count of the resources that failed to enroll.
	Failed        uint64 `protobuf:"varint,3,opt,name=failed,proto3" json:"failed,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ResourcesDiscoveredSummary) Reset() {
	*x = ResourcesDiscoveredSummary{}
	mi := &file_teleport_discoveryconfig_v1_discoveryconfig_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ResourcesDiscoveredSummary) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResourcesDiscoveredSummary) ProtoMessage() {}

func (x *ResourcesDiscoveredSummary) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_discoveryconfig_v1_discoveryconfig_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResourcesDiscoveredSummary.ProtoReflect.Descriptor instead.
func (*ResourcesDiscoveredSummary) Descriptor() ([]byte, []int) {
	return file_teleport_discoveryconfig_v1_discoveryconfig_proto_rawDescGZIP(), []int{4}
}

func (x *ResourcesDiscoveredSummary) GetFound() uint64 {
	if x != nil {
		return x.Found
	}
	return 0
}

func (x *ResourcesDiscoveredSummary) GetEnrolled() uint64 {
	if x != nil {
		return x.Enrolled
	}
	return 0
}

func (x *ResourcesDiscoveredSummary) GetFailed() uint64 {
	if x != nil {
		return x.Failed
	}
	return 0
}

var File_teleport_discoveryconfig_v1_discoveryconfig_proto protoreflect.FileDescriptor

const file_teleport_discoveryconfig_v1_discoveryconfig_proto_rawDesc = "" +
	"\n" +
	"1teleport/discoveryconfig/v1/discoveryconfig.proto\x12\x1bteleport.discoveryconfig.v1\x1a\x1fgoogle/protobuf/timestamp.proto\x1a'teleport/header/v1/resourceheader.proto\x1a!teleport/legacy/types/types.proto\"\xdf\x01\n" +
	"\x0fDiscoveryConfig\x12:\n" +
	"\x06header\x18\x01 \x01(\v2\".teleport.header.v1.ResourceHeaderR\x06header\x12D\n" +
	"\x04spec\x18\x02 \x01(\v20.teleport.discoveryconfig.v1.DiscoveryConfigSpecR\x04spec\x12J\n" +
	"\x06status\x18\x03 \x01(\v22.teleport.discoveryconfig.v1.DiscoveryConfigStatusR\x06status\"\x9c\x02\n" +
	"\x13DiscoveryConfigSpec\x12'\n" +
	"\x0fdiscovery_group\x18\x01 \x01(\tR\x0ediscoveryGroup\x12#\n" +
	"\x03aws\x18\x02 \x03(\v2\x11.types.AWSMatcherR\x03aws\x12)\n" +
	"\x05azure\x18\x03 \x03(\v2\x13.types.AzureMatcherR\x05azure\x12#\n" +
	"\x03gcp\x18\x04 \x03(\v2\x11.types.GCPMatcherR\x03gcp\x12,\n" +
	"\x04kube\x18\x05 \x03(\v2\x18.types.KubernetesMatcherR\x04kube\x129\n" +
	"\faccess_graph\x18\x06 \x01(\v2\x16.types.AccessGraphSyncR\vaccessGraph\"\xe7\x04\n" +
	"\x15DiscoveryConfigStatus\x12G\n" +
	"\x05state\x18\x01 \x01(\x0e21.teleport.discoveryconfig.v1.DiscoveryConfigStateR\x05state\x12(\n" +
	"\rerror_message\x18\x02 \x01(\tH\x00R\ferrorMessage\x88\x01\x01\x121\n" +
	"\x14discovered_resources\x18\x03 \x01(\x04R\x13discoveredResources\x12@\n" +
	"\x0elast_sync_time\x18\x04 \x01(\v2\x1a.google.protobuf.TimestampR\flastSyncTime\x12\xa0\x01\n" +
	" integration_discovered_resources\x18\x06 \x03(\v2V.teleport.discoveryconfig.v1.DiscoveryConfigStatus.IntegrationDiscoveredResourcesEntryR\x1eintegrationDiscoveredResources\x1a\x8c\x01\n" +
	"#IntegrationDiscoveredResourcesEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12O\n" +
	"\x05value\x18\x02 \x01(\v29.teleport.discoveryconfig.v1.IntegrationDiscoveredSummaryR\x05value:\x028\x01B\x10\n" +
	"\x0e_error_messageJ\x04\b\x05\x10\x06R\x1caws_ec2_instances_discovered\"\x94\x02\n" +
	"\x1cIntegrationDiscoveredSummary\x12P\n" +
	"\aaws_ec2\x18\x01 \x01(\v27.teleport.discoveryconfig.v1.ResourcesDiscoveredSummaryR\x06awsEc2\x12P\n" +
	"\aaws_rds\x18\x02 \x01(\v27.teleport.discoveryconfig.v1.ResourcesDiscoveredSummaryR\x06awsRds\x12P\n" +
	"\aaws_eks\x18\x03 \x01(\v27.teleport.discoveryconfig.v1.ResourcesDiscoveredSummaryR\x06awsEks\"f\n" +
	"\x1aResourcesDiscoveredSummary\x12\x14\n" +
	"\x05found\x18\x01 \x01(\x04R\x05found\x12\x1a\n" +
	"\benrolled\x18\x02 \x01(\x04R\benrolled\x12\x16\n" +
	"\x06failed\x18\x03 \x01(\x04R\x06failed*\xa8\x01\n" +
	"\x14DiscoveryConfigState\x12&\n" +
	"\"DISCOVERY_CONFIG_STATE_UNSPECIFIED\x10\x00\x12\"\n" +
	"\x1eDISCOVERY_CONFIG_STATE_RUNNING\x10\x01\x12 \n" +
	"\x1cDISCOVERY_CONFIG_STATE_ERROR\x10\x02\x12\"\n" +
	"\x1eDISCOVERY_CONFIG_STATE_SYNCING\x10\x03BbZ`github.com/gravitational/teleport/api/gen/proto/go/teleport/discoveryconfig/v1;discoveryconfigv1b\x06proto3"

var (
	file_teleport_discoveryconfig_v1_discoveryconfig_proto_rawDescOnce sync.Once
	file_teleport_discoveryconfig_v1_discoveryconfig_proto_rawDescData []byte
)

func file_teleport_discoveryconfig_v1_discoveryconfig_proto_rawDescGZIP() []byte {
	file_teleport_discoveryconfig_v1_discoveryconfig_proto_rawDescOnce.Do(func() {
		file_teleport_discoveryconfig_v1_discoveryconfig_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_teleport_discoveryconfig_v1_discoveryconfig_proto_rawDesc), len(file_teleport_discoveryconfig_v1_discoveryconfig_proto_rawDesc)))
	})
	return file_teleport_discoveryconfig_v1_discoveryconfig_proto_rawDescData
}

var file_teleport_discoveryconfig_v1_discoveryconfig_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_teleport_discoveryconfig_v1_discoveryconfig_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_teleport_discoveryconfig_v1_discoveryconfig_proto_goTypes = []any{
	(DiscoveryConfigState)(0),            // 0: teleport.discoveryconfig.v1.DiscoveryConfigState
	(*DiscoveryConfig)(nil),              // 1: teleport.discoveryconfig.v1.DiscoveryConfig
	(*DiscoveryConfigSpec)(nil),          // 2: teleport.discoveryconfig.v1.DiscoveryConfigSpec
	(*DiscoveryConfigStatus)(nil),        // 3: teleport.discoveryconfig.v1.DiscoveryConfigStatus
	(*IntegrationDiscoveredSummary)(nil), // 4: teleport.discoveryconfig.v1.IntegrationDiscoveredSummary
	(*ResourcesDiscoveredSummary)(nil),   // 5: teleport.discoveryconfig.v1.ResourcesDiscoveredSummary
	nil,                                  // 6: teleport.discoveryconfig.v1.DiscoveryConfigStatus.IntegrationDiscoveredResourcesEntry
	(*v1.ResourceHeader)(nil),            // 7: teleport.header.v1.ResourceHeader
	(*types.AWSMatcher)(nil),             // 8: types.AWSMatcher
	(*types.AzureMatcher)(nil),           // 9: types.AzureMatcher
	(*types.GCPMatcher)(nil),             // 10: types.GCPMatcher
	(*types.KubernetesMatcher)(nil),      // 11: types.KubernetesMatcher
	(*types.AccessGraphSync)(nil),        // 12: types.AccessGraphSync
	(*timestamppb.Timestamp)(nil),        // 13: google.protobuf.Timestamp
}
var file_teleport_discoveryconfig_v1_discoveryconfig_proto_depIdxs = []int32{
	7,  // 0: teleport.discoveryconfig.v1.DiscoveryConfig.header:type_name -> teleport.header.v1.ResourceHeader
	2,  // 1: teleport.discoveryconfig.v1.DiscoveryConfig.spec:type_name -> teleport.discoveryconfig.v1.DiscoveryConfigSpec
	3,  // 2: teleport.discoveryconfig.v1.DiscoveryConfig.status:type_name -> teleport.discoveryconfig.v1.DiscoveryConfigStatus
	8,  // 3: teleport.discoveryconfig.v1.DiscoveryConfigSpec.aws:type_name -> types.AWSMatcher
	9,  // 4: teleport.discoveryconfig.v1.DiscoveryConfigSpec.azure:type_name -> types.AzureMatcher
	10, // 5: teleport.discoveryconfig.v1.DiscoveryConfigSpec.gcp:type_name -> types.GCPMatcher
	11, // 6: teleport.discoveryconfig.v1.DiscoveryConfigSpec.kube:type_name -> types.KubernetesMatcher
	12, // 7: teleport.discoveryconfig.v1.DiscoveryConfigSpec.access_graph:type_name -> types.AccessGraphSync
	0,  // 8: teleport.discoveryconfig.v1.DiscoveryConfigStatus.state:type_name -> teleport.discoveryconfig.v1.DiscoveryConfigState
	13, // 9: teleport.discoveryconfig.v1.DiscoveryConfigStatus.last_sync_time:type_name -> google.protobuf.Timestamp
	6,  // 10: teleport.discoveryconfig.v1.DiscoveryConfigStatus.integration_discovered_resources:type_name -> teleport.discoveryconfig.v1.DiscoveryConfigStatus.IntegrationDiscoveredResourcesEntry
	5,  // 11: teleport.discoveryconfig.v1.IntegrationDiscoveredSummary.aws_ec2:type_name -> teleport.discoveryconfig.v1.ResourcesDiscoveredSummary
	5,  // 12: teleport.discoveryconfig.v1.IntegrationDiscoveredSummary.aws_rds:type_name -> teleport.discoveryconfig.v1.ResourcesDiscoveredSummary
	5,  // 13: teleport.discoveryconfig.v1.IntegrationDiscoveredSummary.aws_eks:type_name -> teleport.discoveryconfig.v1.ResourcesDiscoveredSummary
	4,  // 14: teleport.discoveryconfig.v1.DiscoveryConfigStatus.IntegrationDiscoveredResourcesEntry.value:type_name -> teleport.discoveryconfig.v1.IntegrationDiscoveredSummary
	15, // [15:15] is the sub-list for method output_type
	15, // [15:15] is the sub-list for method input_type
	15, // [15:15] is the sub-list for extension type_name
	15, // [15:15] is the sub-list for extension extendee
	0,  // [0:15] is the sub-list for field type_name
}

func init() { file_teleport_discoveryconfig_v1_discoveryconfig_proto_init() }
func file_teleport_discoveryconfig_v1_discoveryconfig_proto_init() {
	if File_teleport_discoveryconfig_v1_discoveryconfig_proto != nil {
		return
	}
	file_teleport_discoveryconfig_v1_discoveryconfig_proto_msgTypes[2].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_teleport_discoveryconfig_v1_discoveryconfig_proto_rawDesc), len(file_teleport_discoveryconfig_v1_discoveryconfig_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_teleport_discoveryconfig_v1_discoveryconfig_proto_goTypes,
		DependencyIndexes: file_teleport_discoveryconfig_v1_discoveryconfig_proto_depIdxs,
		EnumInfos:         file_teleport_discoveryconfig_v1_discoveryconfig_proto_enumTypes,
		MessageInfos:      file_teleport_discoveryconfig_v1_discoveryconfig_proto_msgTypes,
	}.Build()
	File_teleport_discoveryconfig_v1_discoveryconfig_proto = out.File
	file_teleport_discoveryconfig_v1_discoveryconfig_proto_goTypes = nil
	file_teleport_discoveryconfig_v1_discoveryconfig_proto_depIdxs = nil
}
