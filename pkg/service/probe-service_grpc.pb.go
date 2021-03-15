// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package service

import (
	context "context"

	domain "github.com/luigizuccarelli/golang-eventbus-grpc/pkg/domain"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// HealthClient is the client API for Health service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HealthClient interface {
	Check(ctx context.Context, in *domain.HealthCheckRequest, opts ...grpc.CallOption) (*domain.HealthCheckResponse, error)
	Watch(ctx context.Context, in *domain.HealthCheckRequest, opts ...grpc.CallOption) (Health_WatchClient, error)
}

type healthClient struct {
	cc grpc.ClientConnInterface
}

func NewHealthClient(cc grpc.ClientConnInterface) HealthClient {
	return &healthClient{cc}
}

func (c *healthClient) Check(ctx context.Context, in *domain.HealthCheckRequest, opts ...grpc.CallOption) (*domain.HealthCheckResponse, error) {
	out := new(domain.HealthCheckResponse)
	err := c.cc.Invoke(ctx, "/service.Health/Check", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *healthClient) Watch(ctx context.Context, in *domain.HealthCheckRequest, opts ...grpc.CallOption) (Health_WatchClient, error) {
	stream, err := c.cc.NewStream(ctx, &Health_ServiceDesc.Streams[0], "/service.Health/Watch", opts...)
	if err != nil {
		return nil, err
	}
	x := &healthWatchClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Health_WatchClient interface {
	Recv() (*domain.HealthCheckResponse, error)
	grpc.ClientStream
}

type healthWatchClient struct {
	grpc.ClientStream
}

func (x *healthWatchClient) Recv() (*domain.HealthCheckResponse, error) {
	m := new(domain.HealthCheckResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// HealthServer is the server API for Health service.
// All implementations must embed UnimplementedHealthServer
// for forward compatibility
type HealthServer interface {
	Check(context.Context, *domain.HealthCheckRequest) (*domain.HealthCheckResponse, error)
	Watch(*domain.HealthCheckRequest, Health_WatchServer) error
	//mustEmbedUnimplementedHealthServer()
}

// UnimplementedHealthServer must be embedded to have forward compatible implementations.
type UnimplementedHealthServer struct {
}

func (UnimplementedHealthServer) Check(context.Context, *domain.HealthCheckRequest) (*domain.HealthCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}
func (UnimplementedHealthServer) Watch(*domain.HealthCheckRequest, Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}
func (UnimplementedHealthServer) mustEmbedUnimplementedHealthServer() {}

// UnsafeHealthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HealthServer will
// result in compilation errors.
type UnsafeHealthServer interface {
	mustEmbedUnimplementedHealthServer()
}

func RegisterHealthServer(s grpc.ServiceRegistrar, srv HealthServer) {
	s.RegisterService(&Health_ServiceDesc, srv)
}

func _Health_Check_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(domain.HealthCheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HealthServer).Check(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Health/Check",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HealthServer).Check(ctx, req.(*domain.HealthCheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Health_Watch_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(domain.HealthCheckRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(HealthServer).Watch(m, &healthWatchServer{stream})
}

type Health_WatchServer interface {
	Send(*domain.HealthCheckResponse) error
	grpc.ServerStream
}

type healthWatchServer struct {
	grpc.ServerStream
}

func (x *healthWatchServer) Send(m *domain.HealthCheckResponse) error {
	return x.ServerStream.SendMsg(m)
}

// Health_ServiceDesc is the grpc.ServiceDesc for Health service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Health_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.Health",
	HandlerType: (*HealthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Check",
			Handler:    _Health_Check_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Watch",
			Handler:       _Health_Watch_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "internal/proto-files/service/probe-service.proto",
}
