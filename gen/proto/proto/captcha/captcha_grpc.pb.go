package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const _ = grpc.SupportPackageIsVersion9

const (
	CaptchaService_NewChallenge_FullMethodName      = "/captcha.v1.CaptchaService/NewChallenge"
	CaptchaService_ValidateChallenge_FullMethodName = "/captcha.v1.CaptchaService/ValidateChallenge"
	CaptchaService_MakeEventStream_FullMethodName   = "/captcha.v1.CaptchaService/MakeEventStream"
)

type CaptchaServiceClient interface {
	NewChallenge(ctx context.Context, in *ChallengeRequest, opts ...grpc.CallOption) (*ChallengeResponse, error)
	ValidateChallenge(ctx context.Context, in *ValidateRequest, opts ...grpc.CallOption) (*ValidateResponse, error)
	MakeEventStream(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[ClientEvent, ServerEvent], error)
}

type captchaServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCaptchaServiceClient(cc grpc.ClientConnInterface) CaptchaServiceClient {
	return &captchaServiceClient{cc}
}

func (c *captchaServiceClient) NewChallenge(ctx context.Context, in *ChallengeRequest, opts ...grpc.CallOption) (*ChallengeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ChallengeResponse)
	err := c.cc.Invoke(ctx, CaptchaService_NewChallenge_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *captchaServiceClient) ValidateChallenge(ctx context.Context, in *ValidateRequest, opts ...grpc.CallOption) (*ValidateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ValidateResponse)
	err := c.cc.Invoke(ctx, CaptchaService_ValidateChallenge_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *captchaServiceClient) MakeEventStream(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[ClientEvent, ServerEvent], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &CaptchaService_ServiceDesc.Streams[0], CaptchaService_MakeEventStream_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[ClientEvent, ServerEvent]{ClientStream: stream}
	return x, nil
}

type CaptchaService_MakeEventStreamClient = grpc.BidiStreamingClient[ClientEvent, ServerEvent]

type CaptchaServiceServer interface {
	NewChallenge(context.Context, *ChallengeRequest) (*ChallengeResponse, error)
	ValidateChallenge(context.Context, *ValidateRequest) (*ValidateResponse, error)
	MakeEventStream(grpc.BidiStreamingServer[ClientEvent, ServerEvent]) error
	mustEmbedUnimplementedCaptchaServiceServer()
}

type UnimplementedCaptchaServiceServer struct{}

func (UnimplementedCaptchaServiceServer) NewChallenge(context.Context, *ChallengeRequest) (*ChallengeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewChallenge not implemented")
}
func (UnimplementedCaptchaServiceServer) ValidateChallenge(context.Context, *ValidateRequest) (*ValidateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateChallenge not implemented")
}
func (UnimplementedCaptchaServiceServer) MakeEventStream(grpc.BidiStreamingServer[ClientEvent, ServerEvent]) error {
	return status.Errorf(codes.Unimplemented, "method MakeEventStream not implemented")
}
func (UnimplementedCaptchaServiceServer) mustEmbedUnimplementedCaptchaServiceServer() {}
func (UnimplementedCaptchaServiceServer) testEmbeddedByValue()                        {}

type UnsafeCaptchaServiceServer interface {
	mustEmbedUnimplementedCaptchaServiceServer()
}

func RegisterCaptchaServiceServer(s grpc.ServiceRegistrar, srv CaptchaServiceServer) {
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CaptchaService_ServiceDesc, srv)
}

func _CaptchaService_NewChallenge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChallengeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CaptchaServiceServer).NewChallenge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CaptchaService_NewChallenge_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CaptchaServiceServer).NewChallenge(ctx, req.(*ChallengeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CaptchaService_ValidateChallenge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CaptchaServiceServer).ValidateChallenge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CaptchaService_ValidateChallenge_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CaptchaServiceServer).ValidateChallenge(ctx, req.(*ValidateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CaptchaService_MakeEventStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CaptchaServiceServer).MakeEventStream(&grpc.GenericServerStream[ClientEvent, ServerEvent]{ServerStream: stream})
}

type CaptchaService_MakeEventStreamServer = grpc.BidiStreamingServer[ClientEvent, ServerEvent]

var CaptchaService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "captcha.v1.CaptchaService",
	HandlerType: (*CaptchaServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewChallenge",
			Handler:    _CaptchaService_NewChallenge_Handler,
		},
		{
			MethodName: "ValidateChallenge",
			Handler:    _CaptchaService_ValidateChallenge_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "MakeEventStream",
			Handler:       _CaptchaService_MakeEventStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/captcha/captcha.proto",
}
