// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

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
	GetScoreListByUser(ctx context.Context, in *GetScoreListByUserIDReq, opts ...grpc.CallOption) (*GetScoreListByUserIDResp, error)
}

type scoreServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewScoreServiceClient(cc grpc.ClientConnInterface) ScoreServiceClient {
	return &scoreServiceClient{cc}
}

func (c *scoreServiceClient) AddScoreByUserID(ctx context.Context, in *AddScoreByUserIDReq, opts ...grpc.CallOption) (*AddScoreByUserIDResp, error) {
	out := new(AddScoreByUserIDResp)
	err := c.cc.Invoke(ctx, "/ScoreService/AddScoreByUserID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scoreServiceClient) GetScoreListByUser(ctx context.Context, in *GetScoreListByUserIDReq, opts ...grpc.CallOption) (*GetScoreListByUserIDResp, error) {
	out := new(GetScoreListByUserIDResp)
	err := c.cc.Invoke(ctx, "/ScoreService/GetScoreListByUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ScoreServiceServer is the server API for ScoreService service.
// All implementations must embed UnimplementedScoreServiceServer
// for forward compatibility
type ScoreServiceServer interface {
	// @alias =/score/add/byUser
	AddScoreByUserID(context.Context, *AddScoreByUserIDReq) (*AddScoreByUserIDResp, error)
	// @alias =/score/list
	GetScoreListByUser(context.Context, *GetScoreListByUserIDReq) (*GetScoreListByUserIDResp, error)
	mustEmbedUnimplementedScoreServiceServer()
}

// UnimplementedScoreServiceServer must be embedded to have forward compatible implementations.
type UnimplementedScoreServiceServer struct {
}

func (UnimplementedScoreServiceServer) AddScoreByUserID(context.Context, *AddScoreByUserIDReq) (*AddScoreByUserIDResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddScoreByUserID not implemented")
}
func (UnimplementedScoreServiceServer) GetScoreListByUser(context.Context, *GetScoreListByUserIDReq) (*GetScoreListByUserIDResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetScoreListByUser not implemented")
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
		FullMethod: "/ScoreService/AddScoreByUserID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScoreServiceServer).AddScoreByUserID(ctx, req.(*AddScoreByUserIDReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScoreService_GetScoreListByUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetScoreListByUserIDReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScoreServiceServer).GetScoreListByUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ScoreService/GetScoreListByUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScoreServiceServer).GetScoreListByUser(ctx, req.(*GetScoreListByUserIDReq))
	}
	return interceptor(ctx, in, info, handler)
}

// ScoreService_ServiceDesc is the grpc.ServiceDesc for ScoreService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ScoreService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ScoreService",
	HandlerType: (*ScoreServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddScoreByUserID",
			Handler:    _ScoreService_AddScoreByUserID_Handler,
		},
		{
			MethodName: "GetScoreListByUser",
			Handler:    _ScoreService_GetScoreListByUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "scoreService/proto/score.proto",
}
