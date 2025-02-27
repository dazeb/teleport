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

// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: prehog/v1alpha/connect.proto

package prehogv1alphaconnect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1alpha "github.com/gravitational/teleport/gen/proto/go/prehog/v1alpha"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// ConnectReportingServiceName is the fully-qualified name of the ConnectReportingService service.
	ConnectReportingServiceName = "prehog.v1alpha.ConnectReportingService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ConnectReportingServiceSubmitConnectEventProcedure is the fully-qualified name of the
	// ConnectReportingService's SubmitConnectEvent RPC.
	ConnectReportingServiceSubmitConnectEventProcedure = "/prehog.v1alpha.ConnectReportingService/SubmitConnectEvent"
)

// ConnectReportingServiceClient is a client for the prehog.v1alpha.ConnectReportingService service.
type ConnectReportingServiceClient interface {
	SubmitConnectEvent(context.Context, *connect.Request[v1alpha.SubmitConnectEventRequest]) (*connect.Response[v1alpha.SubmitConnectEventResponse], error)
}

// NewConnectReportingServiceClient constructs a client for the
// prehog.v1alpha.ConnectReportingService service. By default, it uses the Connect protocol with the
// binary Protobuf Codec, asks for gzipped responses, and sends uncompressed requests. To use the
// gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewConnectReportingServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) ConnectReportingServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	connectReportingServiceMethods := v1alpha.File_prehog_v1alpha_connect_proto.Services().ByName("ConnectReportingService").Methods()
	return &connectReportingServiceClient{
		submitConnectEvent: connect.NewClient[v1alpha.SubmitConnectEventRequest, v1alpha.SubmitConnectEventResponse](
			httpClient,
			baseURL+ConnectReportingServiceSubmitConnectEventProcedure,
			connect.WithSchema(connectReportingServiceMethods.ByName("SubmitConnectEvent")),
			connect.WithClientOptions(opts...),
		),
	}
}

// connectReportingServiceClient implements ConnectReportingServiceClient.
type connectReportingServiceClient struct {
	submitConnectEvent *connect.Client[v1alpha.SubmitConnectEventRequest, v1alpha.SubmitConnectEventResponse]
}

// SubmitConnectEvent calls prehog.v1alpha.ConnectReportingService.SubmitConnectEvent.
func (c *connectReportingServiceClient) SubmitConnectEvent(ctx context.Context, req *connect.Request[v1alpha.SubmitConnectEventRequest]) (*connect.Response[v1alpha.SubmitConnectEventResponse], error) {
	return c.submitConnectEvent.CallUnary(ctx, req)
}

// ConnectReportingServiceHandler is an implementation of the prehog.v1alpha.ConnectReportingService
// service.
type ConnectReportingServiceHandler interface {
	SubmitConnectEvent(context.Context, *connect.Request[v1alpha.SubmitConnectEventRequest]) (*connect.Response[v1alpha.SubmitConnectEventResponse], error)
}

// NewConnectReportingServiceHandler builds an HTTP handler from the service implementation. It
// returns the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewConnectReportingServiceHandler(svc ConnectReportingServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	connectReportingServiceMethods := v1alpha.File_prehog_v1alpha_connect_proto.Services().ByName("ConnectReportingService").Methods()
	connectReportingServiceSubmitConnectEventHandler := connect.NewUnaryHandler(
		ConnectReportingServiceSubmitConnectEventProcedure,
		svc.SubmitConnectEvent,
		connect.WithSchema(connectReportingServiceMethods.ByName("SubmitConnectEvent")),
		connect.WithHandlerOptions(opts...),
	)
	return "/prehog.v1alpha.ConnectReportingService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ConnectReportingServiceSubmitConnectEventProcedure:
			connectReportingServiceSubmitConnectEventHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedConnectReportingServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedConnectReportingServiceHandler struct{}

func (UnimplementedConnectReportingServiceHandler) SubmitConnectEvent(context.Context, *connect.Request[v1alpha.SubmitConnectEventRequest]) (*connect.Response[v1alpha.SubmitConnectEventResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("prehog.v1alpha.ConnectReportingService.SubmitConnectEvent is not implemented"))
}
