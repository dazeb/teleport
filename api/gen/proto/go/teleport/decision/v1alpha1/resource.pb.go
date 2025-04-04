// Copyright 2024 Gravitational, Inc
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
// source: teleport/decision/v1alpha1/resource.proto

package decisionpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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

// Resource is the conventional reference type used to refer to the "object" of
// an action that is being considered for an authorization decision. For
// example, a call to EvaluateSSHAccess would use the Resource type to reference
// the ssh node being accessed.
type Resource struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Kind is the type of the resource. Required for requests that support
	// multiple types, otherwise safe to omit.
	Kind string `protobuf:"bytes,1,opt,name=kind,proto3" json:"kind,omitempty"`
	// SubKind is the subtype of the resource. Usually not required as most
	// resources don't have subkinds, or their subkinds do not have an effect on
	// authorization decisions.
	SubKind string `protobuf:"bytes,2,opt,name=sub_kind,json=subKind,proto3" json:"sub_kind,omitempty"`
	// Name is the unique name of the resource.
	Name          string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Resource) Reset() {
	*x = Resource{}
	mi := &file_teleport_decision_v1alpha1_resource_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Resource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Resource) ProtoMessage() {}

func (x *Resource) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_decision_v1alpha1_resource_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Resource.ProtoReflect.Descriptor instead.
func (*Resource) Descriptor() ([]byte, []int) {
	return file_teleport_decision_v1alpha1_resource_proto_rawDescGZIP(), []int{0}
}

func (x *Resource) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *Resource) GetSubKind() string {
	if x != nil {
		return x.SubKind
	}
	return ""
}

func (x *Resource) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_teleport_decision_v1alpha1_resource_proto protoreflect.FileDescriptor

const file_teleport_decision_v1alpha1_resource_proto_rawDesc = "" +
	"\n" +
	")teleport/decision/v1alpha1/resource.proto\x12\x1ateleport.decision.v1alpha1\"M\n" +
	"\bResource\x12\x12\n" +
	"\x04kind\x18\x01 \x01(\tR\x04kind\x12\x19\n" +
	"\bsub_kind\x18\x02 \x01(\tR\asubKind\x12\x12\n" +
	"\x04name\x18\x03 \x01(\tR\x04nameBZZXgithub.com/gravitational/teleport/api/gen/proto/go/teleport/decision/v1alpha1;decisionpbb\x06proto3"

var (
	file_teleport_decision_v1alpha1_resource_proto_rawDescOnce sync.Once
	file_teleport_decision_v1alpha1_resource_proto_rawDescData []byte
)

func file_teleport_decision_v1alpha1_resource_proto_rawDescGZIP() []byte {
	file_teleport_decision_v1alpha1_resource_proto_rawDescOnce.Do(func() {
		file_teleport_decision_v1alpha1_resource_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_teleport_decision_v1alpha1_resource_proto_rawDesc), len(file_teleport_decision_v1alpha1_resource_proto_rawDesc)))
	})
	return file_teleport_decision_v1alpha1_resource_proto_rawDescData
}

var file_teleport_decision_v1alpha1_resource_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_teleport_decision_v1alpha1_resource_proto_goTypes = []any{
	(*Resource)(nil), // 0: teleport.decision.v1alpha1.Resource
}
var file_teleport_decision_v1alpha1_resource_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_teleport_decision_v1alpha1_resource_proto_init() }
func file_teleport_decision_v1alpha1_resource_proto_init() {
	if File_teleport_decision_v1alpha1_resource_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_teleport_decision_v1alpha1_resource_proto_rawDesc), len(file_teleport_decision_v1alpha1_resource_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_teleport_decision_v1alpha1_resource_proto_goTypes,
		DependencyIndexes: file_teleport_decision_v1alpha1_resource_proto_depIdxs,
		MessageInfos:      file_teleport_decision_v1alpha1_resource_proto_msgTypes,
	}.Build()
	File_teleport_decision_v1alpha1_resource_proto = out.File
	file_teleport_decision_v1alpha1_resource_proto_goTypes = nil
	file_teleport_decision_v1alpha1_resource_proto_depIdxs = nil
}
