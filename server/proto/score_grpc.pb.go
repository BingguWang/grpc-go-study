// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: server/proto/score.proto

package proto

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

// ScoreServiceClient is the client API for ScoreService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ScoreServiceClient interface {
	// @alias =/score/add/byUser
	AddScoreByUserID(ctx context.Context, in *AddScoreByUserIDReq, opts ...grpc.CallOption) (*AddScoreByUserIDResp, error)
	// @alias =/score/list
	GetStreamScoreListByUser(ctx context.Context, in *GetScoreListByUserIDReq, opts ...grpc.CallOption) (ScoreService_GetStreamScoreListByUserClient, error)
	AddStreamScoreByUserID(ctx context.Context, opts ...grpc.CallOption) (ScoreService_AddStreamScoreByUserIDClient, error)
	AddAndGetScore(ctx context.Context, opts ...grpc.CallOption) (ScoreService_AddAndGetScoreClient, error)
}

type scoreServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewScoreServiceClient(cc grpc.ClientConnInterface) ScoreServiceClient {
	return &scoreServiceClient{cc}
}

func (c *scoreServiceClient) AddScoreByUserID(ctx context.Context, in *AddScoreByUserIDReq, opts ...grpc.CallOption) (*AddScoreByUserIDResp, error) {
	out := new(AddScoreByUserIDResp)
	err := c.cc.Invoke(ctx, "/proto.ScoreService/AddScoreByUserID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scoreServiceClient) GetStreamScoreListByUser(ctx context.Context, in *GetScoreListByUserIDReq, opts ...grpc.CallOption) (ScoreService_GetStreamScoreListByUserClient, error) {
	stream, err := c.cc.NewStream(ctx, &ScoreService_ServiceDesc.Streams[0], "/proto.ScoreService/GetStreamScoreListByUser", opts...)
	if err != nil {
		return nil, err
	}
	x := &scoreServiceGetStreamScoreListByUserClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ScoreService_GetStreamScoreListByUserClient interface {
	Recv() (*GetScoreListByUserIDResp, error)
	grpc.ClientStream
}

type scoreServiceGetStreamScoreListByUserClient struct {
	grpc.ClientStream
}

func (x *scoreServiceGetStreamScoreListByUserClient) Recv() (*GetScoreListByUserIDResp, error) {
	m := new(GetScoreListByUserIDResp)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *scoreServiceClient) AddStreamScoreByUserID(ctx context.Context, opts ...grpc.CallOption) (ScoreService_AddStreamScoreByUserIDClient, error) {
	stream, err := c.cc.NewStream(ctx, &ScoreService_ServiceDesc.Streams[1], "/proto.ScoreService/AddStreamScoreByUserID", opts...)
	if err != nil {
		return nil, err
	}
	x := &scoreServiceAddStreamScoreByUserIDClient{stream}
	return x, nil
}

type ScoreService_AddStreamScoreByUserIDClient interface {
	Send(*AddScoreByUserIDReq) error
	CloseAndRecv() (*AddScoreByUserIDResp, error)
	grpc.ClientStream
}

type scoreServiceAddStreamScoreByUserIDClient struct {
	grpc.ClientStream
}

func (x *scoreServiceAddStreamScoreByUserIDClient) Send(m *AddScoreByUserIDReq) error {
	return x.ClientStream.SendMsg(m)
}

func (x *scoreServiceAddStreamScoreByUserIDClient) CloseAndRecv() (*AddScoreByUserIDResp, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(AddScoreByUserIDResp)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *scoreServiceClient) AddAndGetScore(ctx context.Context, opts ...grpc.CallOption) (ScoreService_AddAndGetScoreClient, error) {
	stream, err := c.cc.NewStream(ctx, &ScoreService_ServiceDesc.Streams[2], "/proto.ScoreService/AddAndGetScore", opts...)
	if err != nil {
		return nil, err
	}
	x := &scoreServiceAddAndGetScoreClient{stream}
	return x, nil
}

type ScoreService_AddAndGetScoreClient interface {
	Send(*AddScoreByUserIDReq) error
	Recv() (*GetScoreListByUserIDResp, error)
	grpc.ClientStream
}

type scoreServiceAddAndGetScoreClient struct {
	grpc.ClientStream
}

func (x *scoreServiceAddAndGetScoreClient) Send(m *AddScoreByUserIDReq) error {
	return x.ClientStream.SendMsg(m)
}

func (x *scoreServiceAddAndGetScoreClient) Recv() (*GetScoreListByUserIDResp, error) {
	m := new(GetScoreListByUserIDResp)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ScoreServiceServer is the server API for ScoreService service.
// All implementations must embed UnimplementedScoreServiceServer
// for forward compatibility
type ScoreServiceServer interface {
	// @alias =/score/add/byUser
	AddScoreByUserID(context.Context, *AddScoreByUserIDReq) (*AddScoreByUserIDResp, error)
	// @alias =/score/list
	GetStreamScoreListByUser(*GetScoreListByUserIDReq, ScoreService_GetStreamScoreListByUserServer) error
	AddStreamScoreByUserID(ScoreService_AddStreamScoreByUserIDServer) error
	AddAndGetScore(ScoreService_AddAndGetScoreServer) error
	mustEmbedUnimplementedScoreServiceServer()
}

// UnimplementedScoreServiceServer must be embedded to have forward compatible implementations.
type UnimplementedScoreServiceServer struct {
}

func (UnimplementedScoreServiceServer) AddScoreByUserID(context.Context, *AddScoreByUserIDReq) (*AddScoreByUserIDResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddScoreByUserID not implemented")
}
func (UnimplementedScoreServiceServer) GetStreamScoreListByUser(*GetScoreListByUserIDReq, ScoreService_GetStreamScoreListByUserServer) error {
	return status.Errorf(codes.Unimplemented, "method GetStreamScoreListByUser not implemented")
}
func (UnimplementedScoreServiceServer) AddStreamScoreByUserID(ScoreService_AddStreamScoreByUserIDServer) error {
	return status.Errorf(codes.Unimplemented, "method AddStreamScoreByUserID not implemented")
}
func (UnimplementedScoreServiceServer) AddAndGetScore(ScoreService_AddAndGetScoreServer) error {
	return status.Errorf(codes.Unimplemented, "method AddAndGetScore not implemented")
}
func (UnimplementedScoreServiceServer) mustEmbedUnimplementedScoreServiceServer() {}

// UnsafeScoreServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ScoreServiceServer will
// result in compilation errors.
type UnsafeScoreServiceServer interface {
	mustEmbedUnimplementedScoreServiceServer()
}

func RegisterScoreServiceServer(s grpc.ServiceRegistrar, srv ScoreServiceServer) {
	s.RegisterService(&ScoreService_ServiceDesc, srv)
}

func _ScoreService_AddScoreByUserID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddScoreByUserIDReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScoreServiceServer).AddScoreByUserID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ScoreService/AddScoreByUserID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScoreServiceServer).AddScoreByUserID(ctx, req.(*AddScoreByUserIDReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScoreService_GetStreamScoreListByUser_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetScoreListByUserIDReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ScoreServiceServer).GetStreamScoreListByUser(m, &scoreServiceGetStreamScoreListByUserServer{stream})
}

type ScoreService_GetStreamScoreListByUserServer interface {
	Send(*GetScoreListByUserIDResp) error
	grpc.ServerStream
}

type scoreServiceGetStreamScoreListByUserServer struct {
	grpc.ServerStream
}

func (x *scoreServiceGetStreamScoreListByUserServer) Send(m *GetScoreListByUserIDResp) error {
	return x.ServerStream.SendMsg(m)
}

func _ScoreService_AddStreamScoreByUserID_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ScoreServiceServer).AddStreamScoreByUserID(&scoreServiceAddStreamScoreByUserIDServer{stream})
}

type ScoreService_AddStreamScoreByUserIDServer interface {
	SendAndClose(*AddScoreByUserIDResp) error
	Recv() (*AddScoreByUserIDReq, error)
	grpc.ServerStream
}

type scoreServiceAddStreamScoreByUserIDServer struct {
	grpc.ServerStream
}

func (x *scoreServiceAddStreamScoreByUserIDServer) SendAndClose(m *AddScoreByUserIDResp) error {
	return x.ServerStream.SendMsg(m)
}

func (x *scoreServiceAddStreamScoreByUserIDServer) Recv() (*AddScoreByUserIDReq, error) {
	m := new(AddScoreByUserIDReq)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ScoreService_AddAndGetScore_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ScoreServiceServer).AddAndGetScore(&scoreServiceAddAndGetScoreServer{stream})
}

type ScoreService_AddAndGetScoreServer interface {
	Send(*GetScoreListByUserIDResp) error
	Recv() (*AddScoreByUserIDReq, error)
	grpc.ServerStream
}

type scoreServiceAddAndGetScoreServer struct {
	grpc.ServerStream
}

func (x *scoreServiceAddAndGetScoreServer) Send(m *GetScoreListByUserIDResp) error {
	return x.ServerStream.SendMsg(m)
}

func (x *scoreServiceAddAndGetScoreServer) Recv() (*AddScoreByUserIDReq, error) {
	m := new(AddScoreByUserIDReq)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ScoreService_ServiceDesc is the grpc.ServiceDesc for ScoreService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ScoreService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ScoreService",
	HandlerType: (*ScoreServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddScoreByUserID",
			Handler:    _ScoreService_AddScoreByUserID_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetStreamScoreListByUser",
			Handler:       _ScoreService_GetStreamScoreListByUser_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "AddStreamScoreByUserID",
			Handler:       _ScoreService_AddStreamScoreByUserID_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "AddAndGetScore",
			Handler:       _ScoreService_AddAndGetScore_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "server/proto/score.proto",
}
