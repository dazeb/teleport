//
// Teleport
// Copyright (C) 2023  Gravitational, Inc.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: teleport/lib/teleterm/v1/gateway.proto

package teletermv1

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

// Gateway is Teleterm's name for a connection to a resource like a database or a web app
// established through our ALPN proxy.
//
// The term "gateway" is used to avoid using the term "proxy" itself which could be confusing as
// "proxy" means a couple of different things depending on the context. But for Teleterm, a gateway
// is always an ALPN proxy connection.
//
// See RFD 39 for more info on ALPN.
type Gateway struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// uri is the gateway uri
	Uri string `protobuf:"bytes,1,opt,name=uri,proto3" json:"uri,omitempty"`
	// target_name is the target resource name
	TargetName string `protobuf:"bytes,2,opt,name=target_name,json=targetName,proto3" json:"target_name,omitempty"`
	// target_uri is the target uri
	TargetUri string `protobuf:"bytes,3,opt,name=target_uri,json=targetUri,proto3" json:"target_uri,omitempty"`
	// target_user is the target user
	TargetUser string `protobuf:"bytes,4,opt,name=target_user,json=targetUser,proto3" json:"target_user,omitempty"`
	// local_address is the gateway address on localhost
	LocalAddress string `protobuf:"bytes,5,opt,name=local_address,json=localAddress,proto3" json:"local_address,omitempty"`
	// local_port is the gateway address on localhost
	LocalPort string `protobuf:"bytes,6,opt,name=local_port,json=localPort,proto3" json:"local_port,omitempty"`
	// protocol is the protocol used by the gateway. For databases, it matches the type of the
	// database that the gateway targets. For apps, it's either "HTTP" or "TCP".
	Protocol string `protobuf:"bytes,7,opt,name=protocol,proto3" json:"protocol,omitempty"`
	// target_subresource_name points at a subresource of the remote resource, for example a
	// database name on a database server or a target port of a multi-port TCP app.
	TargetSubresourceName string `protobuf:"bytes,9,opt,name=target_subresource_name,json=targetSubresourceName,proto3" json:"target_subresource_name,omitempty"`
	// gateway_cli_client represents a command that the user can execute to connect to the resource
	// through the gateway.
	//
	// Instead of generating those commands in in the frontend code, they are returned from the tsh
	// daemon. This means that the Database Access team can add support for a new protocol and
	// Connect will support it right away with no extra changes.
	GatewayCliCommand *GatewayCLICommand `protobuf:"bytes,10,opt,name=gateway_cli_command,json=gatewayCliCommand,proto3" json:"gateway_cli_command,omitempty"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *Gateway) Reset() {
	*x = Gateway{}
	mi := &file_teleport_lib_teleterm_v1_gateway_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Gateway) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Gateway) ProtoMessage() {}

func (x *Gateway) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_lib_teleterm_v1_gateway_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Gateway.ProtoReflect.Descriptor instead.
func (*Gateway) Descriptor() ([]byte, []int) {
	return file_teleport_lib_teleterm_v1_gateway_proto_rawDescGZIP(), []int{0}
}

func (x *Gateway) GetUri() string {
	if x != nil {
		return x.Uri
	}
	return ""
}

func (x *Gateway) GetTargetName() string {
	if x != nil {
		return x.TargetName
	}
	return ""
}

func (x *Gateway) GetTargetUri() string {
	if x != nil {
		return x.TargetUri
	}
	return ""
}

func (x *Gateway) GetTargetUser() string {
	if x != nil {
		return x.TargetUser
	}
	return ""
}

func (x *Gateway) GetLocalAddress() string {
	if x != nil {
		return x.LocalAddress
	}
	return ""
}

func (x *Gateway) GetLocalPort() string {
	if x != nil {
		return x.LocalPort
	}
	return ""
}

func (x *Gateway) GetProtocol() string {
	if x != nil {
		return x.Protocol
	}
	return ""
}

func (x *Gateway) GetTargetSubresourceName() string {
	if x != nil {
		return x.TargetSubresourceName
	}
	return ""
}

func (x *Gateway) GetGatewayCliCommand() *GatewayCLICommand {
	if x != nil {
		return x.GatewayCliCommand
	}
	return nil
}

// GatewayCLICommand represents a command that the user can execute to connect to a gateway
// resource. It is a combination of two different os/exec.Cmd structs, where path, args and env are
// directly taken from one Cmd and the preview field is constructed from another Cmd.
type GatewayCLICommand struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// path is the absolute path to the CLI client of a gateway if the client is
	// in PATH. Otherwise, the name of the program we were trying to find.
	Path string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	// args is a list containing the name of the program as the first element
	// and the actual args as the other elements
	Args []string `protobuf:"bytes,2,rep,name=args,proto3" json:"args,omitempty"`
	// env is a list of env vars that need to be set for the command
	// invocation. The elements of the list are in the format of NAME=value.
	Env []string `protobuf:"bytes,3,rep,name=env,proto3" json:"env,omitempty"`
	// preview is used to show the user what command will be executed before they decide to run it.
	// It can also be copied and then pasted into a terminal.
	// It's like os/exec.Cmd.String with two exceptions:
	//
	// 1) It is prepended with Cmd.Env.
	// 2) The command name is relative and not absolute.
	// 3) It is taken from a different Cmd than the other fields in this message. This Cmd uses a
	// special print format which makes the args suitable to be entered into a terminal, but not to
	// directly spawn a process.
	//
	// Should not be used to execute the command in the shell. Instead, use path, args, and env.
	Preview       string `protobuf:"bytes,4,opt,name=preview,proto3" json:"preview,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GatewayCLICommand) Reset() {
	*x = GatewayCLICommand{}
	mi := &file_teleport_lib_teleterm_v1_gateway_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GatewayCLICommand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GatewayCLICommand) ProtoMessage() {}

func (x *GatewayCLICommand) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_lib_teleterm_v1_gateway_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GatewayCLICommand.ProtoReflect.Descriptor instead.
func (*GatewayCLICommand) Descriptor() ([]byte, []int) {
	return file_teleport_lib_teleterm_v1_gateway_proto_rawDescGZIP(), []int{1}
}

func (x *GatewayCLICommand) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *GatewayCLICommand) GetArgs() []string {
	if x != nil {
		return x.Args
	}
	return nil
}

func (x *GatewayCLICommand) GetEnv() []string {
	if x != nil {
		return x.Env
	}
	return nil
}

func (x *GatewayCLICommand) GetPreview() string {
	if x != nil {
		return x.Preview
	}
	return ""
}

var File_teleport_lib_teleterm_v1_gateway_proto protoreflect.FileDescriptor

const file_teleport_lib_teleterm_v1_gateway_proto_rawDesc = "" +
	"\n" +
	"&teleport/lib/teleterm/v1/gateway.proto\x12\x18teleport.lib.teleterm.v1\"\x84\x03\n" +
	"\aGateway\x12\x10\n" +
	"\x03uri\x18\x01 \x01(\tR\x03uri\x12\x1f\n" +
	"\vtarget_name\x18\x02 \x01(\tR\n" +
	"targetName\x12\x1d\n" +
	"\n" +
	"target_uri\x18\x03 \x01(\tR\ttargetUri\x12\x1f\n" +
	"\vtarget_user\x18\x04 \x01(\tR\n" +
	"targetUser\x12#\n" +
	"\rlocal_address\x18\x05 \x01(\tR\flocalAddress\x12\x1d\n" +
	"\n" +
	"local_port\x18\x06 \x01(\tR\tlocalPort\x12\x1a\n" +
	"\bprotocol\x18\a \x01(\tR\bprotocol\x126\n" +
	"\x17target_subresource_name\x18\t \x01(\tR\x15targetSubresourceName\x12[\n" +
	"\x13gateway_cli_command\x18\n" +
	" \x01(\v2+.teleport.lib.teleterm.v1.GatewayCLICommandR\x11gatewayCliCommandJ\x04\b\b\x10\tR\vcli_command\"g\n" +
	"\x11GatewayCLICommand\x12\x12\n" +
	"\x04path\x18\x01 \x01(\tR\x04path\x12\x12\n" +
	"\x04args\x18\x02 \x03(\tR\x04args\x12\x10\n" +
	"\x03env\x18\x03 \x03(\tR\x03env\x12\x18\n" +
	"\apreview\x18\x04 \x01(\tR\apreviewBTZRgithub.com/gravitational/teleport/gen/proto/go/teleport/lib/teleterm/v1;teletermv1b\x06proto3"

var (
	file_teleport_lib_teleterm_v1_gateway_proto_rawDescOnce sync.Once
	file_teleport_lib_teleterm_v1_gateway_proto_rawDescData []byte
)

func file_teleport_lib_teleterm_v1_gateway_proto_rawDescGZIP() []byte {
	file_teleport_lib_teleterm_v1_gateway_proto_rawDescOnce.Do(func() {
		file_teleport_lib_teleterm_v1_gateway_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_teleport_lib_teleterm_v1_gateway_proto_rawDesc), len(file_teleport_lib_teleterm_v1_gateway_proto_rawDesc)))
	})
	return file_teleport_lib_teleterm_v1_gateway_proto_rawDescData
}

var file_teleport_lib_teleterm_v1_gateway_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_teleport_lib_teleterm_v1_gateway_proto_goTypes = []any{
	(*Gateway)(nil),           // 0: teleport.lib.teleterm.v1.Gateway
	(*GatewayCLICommand)(nil), // 1: teleport.lib.teleterm.v1.GatewayCLICommand
}
var file_teleport_lib_teleterm_v1_gateway_proto_depIdxs = []int32{
	1, // 0: teleport.lib.teleterm.v1.Gateway.gateway_cli_command:type_name -> teleport.lib.teleterm.v1.GatewayCLICommand
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_teleport_lib_teleterm_v1_gateway_proto_init() }
func file_teleport_lib_teleterm_v1_gateway_proto_init() {
	if File_teleport_lib_teleterm_v1_gateway_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_teleport_lib_teleterm_v1_gateway_proto_rawDesc), len(file_teleport_lib_teleterm_v1_gateway_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_teleport_lib_teleterm_v1_gateway_proto_goTypes,
		DependencyIndexes: file_teleport_lib_teleterm_v1_gateway_proto_depIdxs,
		MessageInfos:      file_teleport_lib_teleterm_v1_gateway_proto_msgTypes,
	}.Build()
	File_teleport_lib_teleterm_v1_gateway_proto = out.File
	file_teleport_lib_teleterm_v1_gateway_proto_goTypes = nil
	file_teleport_lib_teleterm_v1_gateway_proto_depIdxs = nil
}
