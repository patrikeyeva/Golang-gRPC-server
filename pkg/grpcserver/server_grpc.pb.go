// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.0
// source: api/server.proto

package grpcserver

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	DbHandler_CreateArticle_FullMethodName = "/grpcserver.DbHandler/CreateArticle"
	DbHandler_GetArticle_FullMethodName    = "/grpcserver.DbHandler/GetArticle"
	DbHandler_DeleteArticle_FullMethodName = "/grpcserver.DbHandler/DeleteArticle"
	DbHandler_UpdateArticle_FullMethodName = "/grpcserver.DbHandler/UpdateArticle"
	DbHandler_CreateComment_FullMethodName = "/grpcserver.DbHandler/CreateComment"
)

// DbHandlerClient is the client API for DbHandler service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DbHandlerClient interface {
	CreateArticle(ctx context.Context, in *ArticleRequest, opts ...grpc.CallOption) (*ArticleResponse, error)
	GetArticle(ctx context.Context, in *ArticleID, opts ...grpc.CallOption) (*GetArticleResponse, error)
	DeleteArticle(ctx context.Context, in *ArticleID, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UpdateArticle(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	CreateComment(ctx context.Context, in *CommentRequest, opts ...grpc.CallOption) (*CommentResponse, error)
}

type dbHandlerClient struct {
	cc grpc.ClientConnInterface
}

func NewDbHandlerClient(cc grpc.ClientConnInterface) DbHandlerClient {
	return &dbHandlerClient{cc}
}

func (c *dbHandlerClient) CreateArticle(ctx context.Context, in *ArticleRequest, opts ...grpc.CallOption) (*ArticleResponse, error) {
	out := new(ArticleResponse)
	err := c.cc.Invoke(ctx, DbHandler_CreateArticle_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dbHandlerClient) GetArticle(ctx context.Context, in *ArticleID, opts ...grpc.CallOption) (*GetArticleResponse, error) {
	out := new(GetArticleResponse)
	err := c.cc.Invoke(ctx, DbHandler_GetArticle_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dbHandlerClient) DeleteArticle(ctx context.Context, in *ArticleID, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, DbHandler_DeleteArticle_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dbHandlerClient) UpdateArticle(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, DbHandler_UpdateArticle_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dbHandlerClient) CreateComment(ctx context.Context, in *CommentRequest, opts ...grpc.CallOption) (*CommentResponse, error) {
	out := new(CommentResponse)
	err := c.cc.Invoke(ctx, DbHandler_CreateComment_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DbHandlerServer is the server API for DbHandler service.
// All implementations must embed UnimplementedDbHandlerServer
// for forward compatibility
type DbHandlerServer interface {
	CreateArticle(context.Context, *ArticleRequest) (*ArticleResponse, error)
	GetArticle(context.Context, *ArticleID) (*GetArticleResponse, error)
	DeleteArticle(context.Context, *ArticleID) (*emptypb.Empty, error)
	UpdateArticle(context.Context, *UpdateRequest) (*emptypb.Empty, error)
	CreateComment(context.Context, *CommentRequest) (*CommentResponse, error)
	mustEmbedUnimplementedDbHandlerServer()
}

// UnimplementedDbHandlerServer must be embedded to have forward compatible implementations.
type UnimplementedDbHandlerServer struct {
}

func (UnimplementedDbHandlerServer) CreateArticle(context.Context, *ArticleRequest) (*ArticleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateArticle not implemented")
}
func (UnimplementedDbHandlerServer) GetArticle(context.Context, *ArticleID) (*GetArticleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArticle not implemented")
}
func (UnimplementedDbHandlerServer) DeleteArticle(context.Context, *ArticleID) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteArticle not implemented")
}
func (UnimplementedDbHandlerServer) UpdateArticle(context.Context, *UpdateRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateArticle not implemented")
}
func (UnimplementedDbHandlerServer) CreateComment(context.Context, *CommentRequest) (*CommentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateComment not implemented")
}
func (UnimplementedDbHandlerServer) mustEmbedUnimplementedDbHandlerServer() {}

// UnsafeDbHandlerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DbHandlerServer will
// result in compilation errors.
type UnsafeDbHandlerServer interface {
	mustEmbedUnimplementedDbHandlerServer()
}

func RegisterDbHandlerServer(s grpc.ServiceRegistrar, srv DbHandlerServer) {
	s.RegisterService(&DbHandler_ServiceDesc, srv)
}

func _DbHandler_CreateArticle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ArticleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DbHandlerServer).CreateArticle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DbHandler_CreateArticle_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DbHandlerServer).CreateArticle(ctx, req.(*ArticleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DbHandler_GetArticle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ArticleID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DbHandlerServer).GetArticle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DbHandler_GetArticle_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DbHandlerServer).GetArticle(ctx, req.(*ArticleID))
	}
	return interceptor(ctx, in, info, handler)
}

func _DbHandler_DeleteArticle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ArticleID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DbHandlerServer).DeleteArticle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DbHandler_DeleteArticle_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DbHandlerServer).DeleteArticle(ctx, req.(*ArticleID))
	}
	return interceptor(ctx, in, info, handler)
}

func _DbHandler_UpdateArticle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DbHandlerServer).UpdateArticle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DbHandler_UpdateArticle_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DbHandlerServer).UpdateArticle(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DbHandler_CreateComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DbHandlerServer).CreateComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DbHandler_CreateComment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DbHandlerServer).CreateComment(ctx, req.(*CommentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DbHandler_ServiceDesc is the grpc.ServiceDesc for DbHandler service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DbHandler_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpcserver.DbHandler",
	HandlerType: (*DbHandlerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateArticle",
			Handler:    _DbHandler_CreateArticle_Handler,
		},
		{
			MethodName: "GetArticle",
			Handler:    _DbHandler_GetArticle_Handler,
		},
		{
			MethodName: "DeleteArticle",
			Handler:    _DbHandler_DeleteArticle_Handler,
		},
		{
			MethodName: "UpdateArticle",
			Handler:    _DbHandler_UpdateArticle_Handler,
		},
		{
			MethodName: "CreateComment",
			Handler:    _DbHandler_CreateComment_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/server.proto",
}
