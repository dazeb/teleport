// Copyright 2025 Gravitational, Inc
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
// source: teleport/scopes/joining/v1/service.proto

package joiningv1

import (
	v1 "github.com/gravitational/teleport/api/gen/proto/go/teleport/scopes/v1"
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

// GetScopedTokenRequest is the request to get a scoped token.
type GetScopedTokenRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Name is the name of the scoped token.
	Name          string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetScopedTokenRequest) Reset() {
	*x = GetScopedTokenRequest{}
	mi := &file_teleport_scopes_joining_v1_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetScopedTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetScopedTokenRequest) ProtoMessage() {}

func (x *GetScopedTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_scopes_joining_v1_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetScopedTokenRequest.ProtoReflect.Descriptor instead.
func (*GetScopedTokenRequest) Descriptor() ([]byte, []int) {
	return file_teleport_scopes_joining_v1_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetScopedTokenRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// GetScopedTokenResponse is the response to get a scoped token.
type GetScopedTokenResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Token is the scoped token.
	Token         *ScopedToken `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetScopedTokenResponse) Reset() {
	*x = GetScopedTokenResponse{}
	mi := &file_teleport_scopes_joining_v1_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetScopedTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetScopedTokenResponse) ProtoMessage() {}

func (x *GetScopedTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_scopes_joining_v1_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetScopedTokenResponse.ProtoReflect.Descriptor instead.
func (*GetScopedTokenResponse) Descriptor() ([]byte, []int) {
	return file_teleport_scopes_joining_v1_service_proto_rawDescGZIP(), []int{1}
}

func (x *GetScopedTokenResponse) GetToken() *ScopedToken {
	if x != nil {
		return x.Token
	}
	return nil
}

// ListScopedTokensRequest is the request to list scoped tokens.
type ListScopedTokensRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// ResourceScope filters tokens by their resource scope if specified.
	ResourceScope *v1.Filter `protobuf:"bytes,1,opt,name=resource_scope,json=resourceScope,proto3" json:"resource_scope,omitempty"`
	// AssignedScope filters tokens by their assigned scope if specified.
	AssignedScope *v1.Filter `protobuf:"bytes,2,opt,name=assigned_scope,json=assignedScope,proto3" json:"assigned_scope,omitempty"`
	// Cursor is the pagination cursor.
	Cursor string `protobuf:"bytes,3,opt,name=cursor,proto3" json:"cursor,omitempty"`
	// Limit is the maximum number of results to return.
	Limit         uint32 `protobuf:"varint,4,opt,name=limit,proto3" json:"limit,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListScopedTokensRequest) Reset() {
	*x = ListScopedTokensRequest{}
	mi := &file_teleport_scopes_joining_v1_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListScopedTokensRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListScopedTokensRequest) ProtoMessage() {}

func (x *ListScopedTokensRequest) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_scopes_joining_v1_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListScopedTokensRequest.ProtoReflect.Descriptor instead.
func (*ListScopedTokensRequest) Descriptor() ([]byte, []int) {
	return file_teleport_scopes_joining_v1_service_proto_rawDescGZIP(), []int{2}
}

func (x *ListScopedTokensRequest) GetResourceScope() *v1.Filter {
	if x != nil {
		return x.ResourceScope
	}
	return nil
}

func (x *ListScopedTokensRequest) GetAssignedScope() *v1.Filter {
	if x != nil {
		return x.AssignedScope
	}
	return nil
}

func (x *ListScopedTokensRequest) GetCursor() string {
	if x != nil {
		return x.Cursor
	}
	return ""
}

func (x *ListScopedTokensRequest) GetLimit() uint32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

// ListScopedTokensResponse is the response to list scoped tokens.
type ListScopedTokensResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Tokens is the list of scoped tokens.
	Tokens []*ScopedToken `protobuf:"bytes,1,rep,name=tokens,proto3" json:"tokens,omitempty"`
	// Cursor is the pagination cursor.
	Cursor        string `protobuf:"bytes,2,opt,name=cursor,proto3" json:"cursor,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListScopedTokensResponse) Reset() {
	*x = ListScopedTokensResponse{}
	mi := &file_teleport_scopes_joining_v1_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListScopedTokensResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListScopedTokensResponse) ProtoMessage() {}

func (x *ListScopedTokensResponse) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_scopes_joining_v1_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListScopedTokensResponse.ProtoReflect.Descriptor instead.
func (*ListScopedTokensResponse) Descriptor() ([]byte, []int) {
	return file_teleport_scopes_joining_v1_service_proto_rawDescGZIP(), []int{3}
}

func (x *ListScopedTokensResponse) GetTokens() []*ScopedToken {
	if x != nil {
		return x.Tokens
	}
	return nil
}

func (x *ListScopedTokensResponse) GetCursor() string {
	if x != nil {
		return x.Cursor
	}
	return ""
}

// CreateScopedTokenRequest is the request to create a scoped token.
type CreateScopedTokenRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Token is the scoped token to create.
	Token         *ScopedToken `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateScopedTokenRequest) Reset() {
	*x = CreateScopedTokenRequest{}
	mi := &file_teleport_scopes_joining_v1_service_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateScopedTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateScopedTokenRequest) ProtoMessage() {}

func (x *CreateScopedTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_scopes_joining_v1_service_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateScopedTokenRequest.ProtoReflect.Descriptor instead.
func (*CreateScopedTokenRequest) Descriptor() ([]byte, []int) {
	return file_teleport_scopes_joining_v1_service_proto_rawDescGZIP(), []int{4}
}

func (x *CreateScopedTokenRequest) GetToken() *ScopedToken {
	if x != nil {
		return x.Token
	}
	return nil
}

// CreateScopedTokenResponse is the response to create a scoped token.
type CreateScopedTokenResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Token is the scoped token that was created.
	Token         *ScopedToken `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateScopedTokenResponse) Reset() {
	*x = CreateScopedTokenResponse{}
	mi := &file_teleport_scopes_joining_v1_service_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateScopedTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateScopedTokenResponse) ProtoMessage() {}

func (x *CreateScopedTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_scopes_joining_v1_service_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateScopedTokenResponse.ProtoReflect.Descriptor instead.
func (*CreateScopedTokenResponse) Descriptor() ([]byte, []int) {
	return file_teleport_scopes_joining_v1_service_proto_rawDescGZIP(), []int{5}
}

func (x *CreateScopedTokenResponse) GetToken() *ScopedToken {
	if x != nil {
		return x.Token
	}
	return nil
}

// UpdateScopedTokenRequest is the request to update a scoped token.
type UpdateScopedTokenRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Token is the scoped token to update.
	Token         *ScopedToken `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateScopedTokenRequest) Reset() {
	*x = UpdateScopedTokenRequest{}
	mi := &file_teleport_scopes_joining_v1_service_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateScopedTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateScopedTokenRequest) ProtoMessage() {}

func (x *UpdateScopedTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_scopes_joining_v1_service_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateScopedTokenRequest.ProtoReflect.Descriptor instead.
func (*UpdateScopedTokenRequest) Descriptor() ([]byte, []int) {
	return file_teleport_scopes_joining_v1_service_proto_rawDescGZIP(), []int{6}
}

func (x *UpdateScopedTokenRequest) GetToken() *ScopedToken {
	if x != nil {
		return x.Token
	}
	return nil
}

// UpdateScopedTokenResponse is the response to update a scoped token.
type UpdateScopedTokenResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Token is the post-update scoped token.
	Token         *ScopedToken `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateScopedTokenResponse) Reset() {
	*x = UpdateScopedTokenResponse{}
	mi := &file_teleport_scopes_joining_v1_service_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateScopedTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateScopedTokenResponse) ProtoMessage() {}

func (x *UpdateScopedTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_scopes_joining_v1_service_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateScopedTokenResponse.ProtoReflect.Descriptor instead.
func (*UpdateScopedTokenResponse) Descriptor() ([]byte, []int) {
	return file_teleport_scopes_joining_v1_service_proto_rawDescGZIP(), []int{7}
}

func (x *UpdateScopedTokenResponse) GetToken() *ScopedToken {
	if x != nil {
		return x.Token
	}
	return nil
}

// DeleteScopedTokenRequest is the request to delete a scoped token.
type DeleteScopedTokenRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Name is the name of the scoped token to delete.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Revision asserts the revision of the scoped token to delete (optional).
	Revision      string `protobuf:"bytes,2,opt,name=revision,proto3" json:"revision,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteScopedTokenRequest) Reset() {
	*x = DeleteScopedTokenRequest{}
	mi := &file_teleport_scopes_joining_v1_service_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteScopedTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteScopedTokenRequest) ProtoMessage() {}

func (x *DeleteScopedTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_scopes_joining_v1_service_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteScopedTokenRequest.ProtoReflect.Descriptor instead.
func (*DeleteScopedTokenRequest) Descriptor() ([]byte, []int) {
	return file_teleport_scopes_joining_v1_service_proto_rawDescGZIP(), []int{8}
}

func (x *DeleteScopedTokenRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *DeleteScopedTokenRequest) GetRevision() string {
	if x != nil {
		return x.Revision
	}
	return ""
}

// DeleteScopedTokenResponse is the response to delete a scoped token.
type DeleteScopedTokenResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteScopedTokenResponse) Reset() {
	*x = DeleteScopedTokenResponse{}
	mi := &file_teleport_scopes_joining_v1_service_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteScopedTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteScopedTokenResponse) ProtoMessage() {}

func (x *DeleteScopedTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_scopes_joining_v1_service_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteScopedTokenResponse.ProtoReflect.Descriptor instead.
func (*DeleteScopedTokenResponse) Descriptor() ([]byte, []int) {
	return file_teleport_scopes_joining_v1_service_proto_rawDescGZIP(), []int{9}
}

var File_teleport_scopes_joining_v1_service_proto protoreflect.FileDescriptor

const file_teleport_scopes_joining_v1_service_proto_rawDesc = "" +
	"\n" +
	"(teleport/scopes/joining/v1/service.proto\x12\x1ateleport.scopes.joining.v1\x1a&teleport/scopes/joining/v1/token.proto\x1a\x1fteleport/scopes/v1/scopes.proto\"+\n" +
	"\x15GetScopedTokenRequest\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\"W\n" +
	"\x16GetScopedTokenResponse\x12=\n" +
	"\x05token\x18\x01 \x01(\v2'.teleport.scopes.joining.v1.ScopedTokenR\x05token\"\xcd\x01\n" +
	"\x17ListScopedTokensRequest\x12A\n" +
	"\x0eresource_scope\x18\x01 \x01(\v2\x1a.teleport.scopes.v1.FilterR\rresourceScope\x12A\n" +
	"\x0eassigned_scope\x18\x02 \x01(\v2\x1a.teleport.scopes.v1.FilterR\rassignedScope\x12\x16\n" +
	"\x06cursor\x18\x03 \x01(\tR\x06cursor\x12\x14\n" +
	"\x05limit\x18\x04 \x01(\rR\x05limit\"s\n" +
	"\x18ListScopedTokensResponse\x12?\n" +
	"\x06tokens\x18\x01 \x03(\v2'.teleport.scopes.joining.v1.ScopedTokenR\x06tokens\x12\x16\n" +
	"\x06cursor\x18\x02 \x01(\tR\x06cursor\"Y\n" +
	"\x18CreateScopedTokenRequest\x12=\n" +
	"\x05token\x18\x01 \x01(\v2'.teleport.scopes.joining.v1.ScopedTokenR\x05token\"Z\n" +
	"\x19CreateScopedTokenResponse\x12=\n" +
	"\x05token\x18\x01 \x01(\v2'.teleport.scopes.joining.v1.ScopedTokenR\x05token\"Y\n" +
	"\x18UpdateScopedTokenRequest\x12=\n" +
	"\x05token\x18\x01 \x01(\v2'.teleport.scopes.joining.v1.ScopedTokenR\x05token\"Z\n" +
	"\x19UpdateScopedTokenResponse\x12=\n" +
	"\x05token\x18\x01 \x01(\v2'.teleport.scopes.joining.v1.ScopedTokenR\x05token\"J\n" +
	"\x18DeleteScopedTokenRequest\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12\x1a\n" +
	"\brevision\x18\x02 \x01(\tR\brevision\"\x1b\n" +
	"\x19DeleteScopedTokenResponse2\x97\x05\n" +
	"\x14ScopedJoiningService\x12w\n" +
	"\x0eGetScopedToken\x121.teleport.scopes.joining.v1.GetScopedTokenRequest\x1a2.teleport.scopes.joining.v1.GetScopedTokenResponse\x12}\n" +
	"\x10ListScopedTokens\x123.teleport.scopes.joining.v1.ListScopedTokensRequest\x1a4.teleport.scopes.joining.v1.ListScopedTokensResponse\x12\x80\x01\n" +
	"\x11CreateScopedToken\x124.teleport.scopes.joining.v1.CreateScopedTokenRequest\x1a5.teleport.scopes.joining.v1.CreateScopedTokenResponse\x12\x80\x01\n" +
	"\x11UpdateScopedToken\x124.teleport.scopes.joining.v1.UpdateScopedTokenRequest\x1a5.teleport.scopes.joining.v1.UpdateScopedTokenResponse\x12\x80\x01\n" +
	"\x11DeleteScopedToken\x124.teleport.scopes.joining.v1.DeleteScopedTokenRequest\x1a5.teleport.scopes.joining.v1.DeleteScopedTokenResponseBYZWgithub.com/gravitational/teleport/api/gen/proto/go/teleport/scopes/joining/v1;joiningv1b\x06proto3"

var (
	file_teleport_scopes_joining_v1_service_proto_rawDescOnce sync.Once
	file_teleport_scopes_joining_v1_service_proto_rawDescData []byte
)

func file_teleport_scopes_joining_v1_service_proto_rawDescGZIP() []byte {
	file_teleport_scopes_joining_v1_service_proto_rawDescOnce.Do(func() {
		file_teleport_scopes_joining_v1_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_teleport_scopes_joining_v1_service_proto_rawDesc), len(file_teleport_scopes_joining_v1_service_proto_rawDesc)))
	})
	return file_teleport_scopes_joining_v1_service_proto_rawDescData
}

var file_teleport_scopes_joining_v1_service_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_teleport_scopes_joining_v1_service_proto_goTypes = []any{
	(*GetScopedTokenRequest)(nil),     // 0: teleport.scopes.joining.v1.GetScopedTokenRequest
	(*GetScopedTokenResponse)(nil),    // 1: teleport.scopes.joining.v1.GetScopedTokenResponse
	(*ListScopedTokensRequest)(nil),   // 2: teleport.scopes.joining.v1.ListScopedTokensRequest
	(*ListScopedTokensResponse)(nil),  // 3: teleport.scopes.joining.v1.ListScopedTokensResponse
	(*CreateScopedTokenRequest)(nil),  // 4: teleport.scopes.joining.v1.CreateScopedTokenRequest
	(*CreateScopedTokenResponse)(nil), // 5: teleport.scopes.joining.v1.CreateScopedTokenResponse
	(*UpdateScopedTokenRequest)(nil),  // 6: teleport.scopes.joining.v1.UpdateScopedTokenRequest
	(*UpdateScopedTokenResponse)(nil), // 7: teleport.scopes.joining.v1.UpdateScopedTokenResponse
	(*DeleteScopedTokenRequest)(nil),  // 8: teleport.scopes.joining.v1.DeleteScopedTokenRequest
	(*DeleteScopedTokenResponse)(nil), // 9: teleport.scopes.joining.v1.DeleteScopedTokenResponse
	(*ScopedToken)(nil),               // 10: teleport.scopes.joining.v1.ScopedToken
	(*v1.Filter)(nil),                 // 11: teleport.scopes.v1.Filter
}
var file_teleport_scopes_joining_v1_service_proto_depIdxs = []int32{
	10, // 0: teleport.scopes.joining.v1.GetScopedTokenResponse.token:type_name -> teleport.scopes.joining.v1.ScopedToken
	11, // 1: teleport.scopes.joining.v1.ListScopedTokensRequest.resource_scope:type_name -> teleport.scopes.v1.Filter
	11, // 2: teleport.scopes.joining.v1.ListScopedTokensRequest.assigned_scope:type_name -> teleport.scopes.v1.Filter
	10, // 3: teleport.scopes.joining.v1.ListScopedTokensResponse.tokens:type_name -> teleport.scopes.joining.v1.ScopedToken
	10, // 4: teleport.scopes.joining.v1.CreateScopedTokenRequest.token:type_name -> teleport.scopes.joining.v1.ScopedToken
	10, // 5: teleport.scopes.joining.v1.CreateScopedTokenResponse.token:type_name -> teleport.scopes.joining.v1.ScopedToken
	10, // 6: teleport.scopes.joining.v1.UpdateScopedTokenRequest.token:type_name -> teleport.scopes.joining.v1.ScopedToken
	10, // 7: teleport.scopes.joining.v1.UpdateScopedTokenResponse.token:type_name -> teleport.scopes.joining.v1.ScopedToken
	0,  // 8: teleport.scopes.joining.v1.ScopedJoiningService.GetScopedToken:input_type -> teleport.scopes.joining.v1.GetScopedTokenRequest
	2,  // 9: teleport.scopes.joining.v1.ScopedJoiningService.ListScopedTokens:input_type -> teleport.scopes.joining.v1.ListScopedTokensRequest
	4,  // 10: teleport.scopes.joining.v1.ScopedJoiningService.CreateScopedToken:input_type -> teleport.scopes.joining.v1.CreateScopedTokenRequest
	6,  // 11: teleport.scopes.joining.v1.ScopedJoiningService.UpdateScopedToken:input_type -> teleport.scopes.joining.v1.UpdateScopedTokenRequest
	8,  // 12: teleport.scopes.joining.v1.ScopedJoiningService.DeleteScopedToken:input_type -> teleport.scopes.joining.v1.DeleteScopedTokenRequest
	1,  // 13: teleport.scopes.joining.v1.ScopedJoiningService.GetScopedToken:output_type -> teleport.scopes.joining.v1.GetScopedTokenResponse
	3,  // 14: teleport.scopes.joining.v1.ScopedJoiningService.ListScopedTokens:output_type -> teleport.scopes.joining.v1.ListScopedTokensResponse
	5,  // 15: teleport.scopes.joining.v1.ScopedJoiningService.CreateScopedToken:output_type -> teleport.scopes.joining.v1.CreateScopedTokenResponse
	7,  // 16: teleport.scopes.joining.v1.ScopedJoiningService.UpdateScopedToken:output_type -> teleport.scopes.joining.v1.UpdateScopedTokenResponse
	9,  // 17: teleport.scopes.joining.v1.ScopedJoiningService.DeleteScopedToken:output_type -> teleport.scopes.joining.v1.DeleteScopedTokenResponse
	13, // [13:18] is the sub-list for method output_type
	8,  // [8:13] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_teleport_scopes_joining_v1_service_proto_init() }
func file_teleport_scopes_joining_v1_service_proto_init() {
	if File_teleport_scopes_joining_v1_service_proto != nil {
		return
	}
	file_teleport_scopes_joining_v1_token_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_teleport_scopes_joining_v1_service_proto_rawDesc), len(file_teleport_scopes_joining_v1_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_teleport_scopes_joining_v1_service_proto_goTypes,
		DependencyIndexes: file_teleport_scopes_joining_v1_service_proto_depIdxs,
		MessageInfos:      file_teleport_scopes_joining_v1_service_proto_msgTypes,
	}.Build()
	File_teleport_scopes_joining_v1_service_proto = out.File
	file_teleport_scopes_joining_v1_service_proto_goTypes = nil
	file_teleport_scopes_joining_v1_service_proto_depIdxs = nil
}
