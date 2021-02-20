// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package service

import (
	context "context"

	domain "github.com/luigizuccarelli/golang-eventbus-grpc/internal/gRPC/domain"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// DataSchemaServiceClient is the client API for DataSchemaService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DataSchemaServiceClient interface {
	Get(ctx context.Context, in *domain.DataSchema, opts ...grpc.CallOption) (*GetDataSchemaResponse, error)
}

type dataSchemaServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDataSchemaServiceClient(cc grpc.ClientConnInterface) DataSchemaServiceClient {
	return &dataSchemaServiceClient{cc}
}

func (c *dataSchemaServiceClient) Get(ctx context.Context, in *domain.DataSchema, opts ...grpc.CallOption) (*GetDataSchemaResponse, error) {
	out := new(GetDataSchemaResponse)
	err := c.cc.Invoke(ctx, "/service.DataSchemaService/get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DataSchemaServiceServer is the server API for DataSchemaService service.
// All implementations must embed UnimplementedDataSchemaServiceServer
// for forward compatibility
type DataSchemaServiceServer interface {
	Get(context.Context, *domain.DataSchema) (*GetDataSchemaResponse, error)
	//mustEmbedUnimplementedDataSchemaServiceServer()
}

// UnimplementedDataSchemaServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDataSchemaServiceServer struct {
}

func (UnimplementedDataSchemaServiceServer) Get(context.Context, *domain.DataSchema) (*GetDataSchemaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedDataSchemaServiceServer) mustEmbedUnimplementedDataSchemaServiceServer() {}

// UnsafeDataSchemaServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DataSchemaServiceServer will
// result in compilation errors.
type UnsafeDataSchemaServiceServer interface {
	mustEmbedUnimplementedDataSchemaServiceServer()
}

func RegisterDataSchemaServiceServer(s grpc.ServiceRegistrar, srv DataSchemaServiceServer) {
	s.RegisterService(&DataSchemaService_ServiceDesc, srv)
}

func _DataSchemaService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(domain.DataSchema)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataSchemaServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.DataSchemaService/get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataSchemaServiceServer).Get(ctx, req.(*domain.DataSchema))
	}
	return interceptor(ctx, in, info, handler)
}

// DataSchemaService_ServiceDesc is the grpc.ServiceDesc for DataSchemaService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DataSchemaService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.DataSchemaService",
	HandlerType: (*DataSchemaServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "get",
			Handler:    _DataSchemaService_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/proto-files/service/schema-service.proto",
}
