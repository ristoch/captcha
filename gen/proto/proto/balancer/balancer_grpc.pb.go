package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const _ = grpc.SupportPackageIsVersion9

const (
	BalancerService_RegisterInstance_FullMethodName = "/balancer.v1.BalancerService/RegisterInstance"
	BalancerService_CheckUserBlocked_FullMethodName = "/balancer.v1.BalancerService/CheckUserBlocked"
	BalancerService_BlockUser_FullMethodName        = "/balancer.v1.BalancerService/BlockUser"
	BalancerService_GetInstances_FullMethodName     = "/balancer.v1.BalancerService/GetInstances"
)

type BalancerServiceClient interface {
	RegisterInstance(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[RegisterInstanceRequest, RegisterInstanceResponse], error)
	CheckUserBlocked(ctx context.Context, in *CheckUserBlockedRequest, opts ...grpc.CallOption) (*CheckUserBlockedResponse, error)
	BlockUser(ctx context.Context, in *BlockUserRequest, opts ...grpc.CallOption) (*BlockUserResponse, error)
	GetInstances(ctx context.Context, in *GetInstancesRequest, opts ...grpc.CallOption) (*GetInstancesResponse, error)
}

type balancerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBalancerServiceClient(cc grpc.ClientConnInterface) BalancerServiceClient {
	return &balancerServiceClient{cc}
}

func (c *balancerServiceClient) RegisterInstance(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[RegisterInstanceRequest, RegisterInstanceResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &BalancerService_ServiceDesc.Streams[0], BalancerService_RegisterInstance_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[RegisterInstanceRequest, RegisterInstanceResponse]{ClientStream: stream}
	return x, nil
}

type BalancerService_RegisterInstanceClient = grpc.BidiStreamingClient[RegisterInstanceRequest, RegisterInstanceResponse]

func (c *balancerServiceClient) CheckUserBlocked(ctx context.Context, in *CheckUserBlockedRequest, opts ...grpc.CallOption) (*CheckUserBlockedResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CheckUserBlockedResponse)
	err := c.cc.Invoke(ctx, BalancerService_CheckUserBlocked_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *balancerServiceClient) BlockUser(ctx context.Context, in *BlockUserRequest, opts ...grpc.CallOption) (*BlockUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BlockUserResponse)
	err := c.cc.Invoke(ctx, BalancerService_BlockUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *balancerServiceClient) GetInstances(ctx context.Context, in *GetInstancesRequest, opts ...grpc.CallOption) (*GetInstancesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetInstancesResponse)
	err := c.cc.Invoke(ctx, BalancerService_GetInstances_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type BalancerServiceServer interface {
	RegisterInstance(grpc.BidiStreamingServer[RegisterInstanceRequest, RegisterInstanceResponse]) error
	CheckUserBlocked(context.Context, *CheckUserBlockedRequest) (*CheckUserBlockedResponse, error)
	BlockUser(context.Context, *BlockUserRequest) (*BlockUserResponse, error)
	GetInstances(context.Context, *GetInstancesRequest) (*GetInstancesResponse, error)
	mustEmbedUnimplementedBalancerServiceServer()
}

type UnimplementedBalancerServiceServer struct{}

func (UnimplementedBalancerServiceServer) RegisterInstance(grpc.BidiStreamingServer[RegisterInstanceRequest, RegisterInstanceResponse]) error {
	return status.Errorf(codes.Unimplemented, "method RegisterInstance not implemented")
}
func (UnimplementedBalancerServiceServer) CheckUserBlocked(context.Context, *CheckUserBlockedRequest) (*CheckUserBlockedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckUserBlocked not implemented")
}
func (UnimplementedBalancerServiceServer) BlockUser(context.Context, *BlockUserRequest) (*BlockUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BlockUser not implemented")
}
func (UnimplementedBalancerServiceServer) GetInstances(context.Context, *GetInstancesRequest) (*GetInstancesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInstances not implemented")
}
func (UnimplementedBalancerServiceServer) mustEmbedUnimplementedBalancerServiceServer() {}
func (UnimplementedBalancerServiceServer) testEmbeddedByValue()                         {}

type UnsafeBalancerServiceServer interface {
	mustEmbedUnimplementedBalancerServiceServer()
}

func RegisterBalancerServiceServer(s grpc.ServiceRegistrar, srv BalancerServiceServer) {
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&BalancerService_ServiceDesc, srv)
}

func _BalancerService_RegisterInstance_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BalancerServiceServer).RegisterInstance(&grpc.GenericServerStream[RegisterInstanceRequest, RegisterInstanceResponse]{ServerStream: stream})
}

type BalancerService_RegisterInstanceServer = grpc.BidiStreamingServer[RegisterInstanceRequest, RegisterInstanceResponse]

func _BalancerService_CheckUserBlocked_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckUserBlockedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BalancerServiceServer).CheckUserBlocked(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BalancerService_CheckUserBlocked_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BalancerServiceServer).CheckUserBlocked(ctx, req.(*CheckUserBlockedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BalancerService_BlockUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BlockUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BalancerServiceServer).BlockUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BalancerService_BlockUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BalancerServiceServer).BlockUser(ctx, req.(*BlockUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BalancerService_GetInstances_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetInstancesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BalancerServiceServer).GetInstances(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BalancerService_GetInstances_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BalancerServiceServer).GetInstances(ctx, req.(*GetInstancesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var BalancerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "balancer.v1.BalancerService",
	HandlerType: (*BalancerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckUserBlocked",
			Handler:    _BalancerService_CheckUserBlocked_Handler,
		},
		{
			MethodName: "BlockUser",
			Handler:    _BalancerService_BlockUser_Handler,
		},
		{
			MethodName: "GetInstances",
			Handler:    _BalancerService_GetInstances_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RegisterInstance",
			Handler:       _BalancerService_RegisterInstance_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/balancer/balancer.proto",
}
