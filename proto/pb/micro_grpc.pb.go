// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MicroClient is the client API for Micro service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MicroClient interface {
	// Sends a greeting
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
}

type microClient struct {
	cc grpc.ClientConnInterface
}

func NewMicroClient(cc grpc.ClientConnInterface) MicroClient {
	return &microClient{cc}
}

func (c *microClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/pb.Micro/SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MicroServer is the server API for Micro service.
// All implementations must embed UnimplementedMicroServer
// for forward compatibility
type MicroServer interface {
	// Sends a greeting
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
	mustEmbedUnimplementedMicroServer()
}

// UnimplementedMicroServer must be embedded to have forward compatible implementations.
type UnimplementedMicroServer struct {
}

func (UnimplementedMicroServer) SayHello(context.Context, *HelloRequest) (*HelloReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (UnimplementedMicroServer) mustEmbedUnimplementedMicroServer() {}

// UnsafeMicroServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MicroServer will
// result in compilation errors.
type UnsafeMicroServer interface {
	mustEmbedUnimplementedMicroServer()
}

func RegisterMicroServer(s grpc.ServiceRegistrar, srv MicroServer) {
	s.RegisterService(&Micro_ServiceDesc, srv)
}

func _Micro_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MicroServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Micro/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MicroServer).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Micro_ServiceDesc is the grpc.ServiceDesc for Micro service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Micro_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Micro",
	HandlerType: (*MicroServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Micro_SayHello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/micro.proto",
}