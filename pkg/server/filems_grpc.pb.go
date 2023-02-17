// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.11.2
// source: filems.proto

package server

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

// FileMsServiceClient is the client API for FileMsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FileMsServiceClient interface {
	// Create file by content
	CreateFile(ctx context.Context, in *CreateFileRequest, opts ...grpc.CallOption) (*FileRequest, error)
	// update file content
	UpdateFileContent(ctx context.Context, in *UpdateFileRequest, opts ...grpc.CallOption) (*FileRequest, error)
	// get file, return url
	GetFile(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileResponse, error)
	// get file content
	GetFileContent(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileContentResponse, error)
	// delete file
	DelFile(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileResponse, error)
	// list file versions
	ListFileVersions(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileVersionResponse, error)
}

type fileMsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFileMsServiceClient(cc grpc.ClientConnInterface) FileMsServiceClient {
	return &fileMsServiceClient{cc}
}

func (c *fileMsServiceClient) CreateFile(ctx context.Context, in *CreateFileRequest, opts ...grpc.CallOption) (*FileRequest, error) {
	out := new(FileRequest)
	err := c.cc.Invoke(ctx, "/proto.FileMsService/CreateFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileMsServiceClient) UpdateFileContent(ctx context.Context, in *UpdateFileRequest, opts ...grpc.CallOption) (*FileRequest, error) {
	out := new(FileRequest)
	err := c.cc.Invoke(ctx, "/proto.FileMsService/UpdateFileContent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileMsServiceClient) GetFile(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileResponse, error) {
	out := new(FileResponse)
	err := c.cc.Invoke(ctx, "/proto.FileMsService/GetFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileMsServiceClient) GetFileContent(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileContentResponse, error) {
	out := new(FileContentResponse)
	err := c.cc.Invoke(ctx, "/proto.FileMsService/GetFileContent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileMsServiceClient) DelFile(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileResponse, error) {
	out := new(FileResponse)
	err := c.cc.Invoke(ctx, "/proto.FileMsService/DelFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileMsServiceClient) ListFileVersions(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileVersionResponse, error) {
	out := new(FileVersionResponse)
	err := c.cc.Invoke(ctx, "/proto.FileMsService/ListFileVersions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FileMsServiceServer is the server API for FileMsService service.
// All implementations should embed UnimplementedFileMsServiceServer
// for forward compatibility
type FileMsServiceServer interface {
	// Create file by content
	CreateFile(context.Context, *CreateFileRequest) (*FileRequest, error)
	// update file content
	UpdateFileContent(context.Context, *UpdateFileRequest) (*FileRequest, error)
	// get file, return url
	GetFile(context.Context, *FileRequest) (*FileResponse, error)
	// get file content
	GetFileContent(context.Context, *FileRequest) (*FileContentResponse, error)
	// delete file
	DelFile(context.Context, *FileRequest) (*FileResponse, error)
	// list file versions
	ListFileVersions(context.Context, *FileRequest) (*FileVersionResponse, error)
}

// UnimplementedFileMsServiceServer should be embedded to have forward compatible implementations.
type UnimplementedFileMsServiceServer struct {
}

func (UnimplementedFileMsServiceServer) CreateFile(context.Context, *CreateFileRequest) (*FileRequest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFile not implemented")
}
func (UnimplementedFileMsServiceServer) UpdateFileContent(context.Context, *UpdateFileRequest) (*FileRequest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateFileContent not implemented")
}
func (UnimplementedFileMsServiceServer) GetFile(context.Context, *FileRequest) (*FileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFile not implemented")
}
func (UnimplementedFileMsServiceServer) GetFileContent(context.Context, *FileRequest) (*FileContentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFileContent not implemented")
}
func (UnimplementedFileMsServiceServer) DelFile(context.Context, *FileRequest) (*FileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelFile not implemented")
}
func (UnimplementedFileMsServiceServer) ListFileVersions(context.Context, *FileRequest) (*FileVersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListFileVersions not implemented")
}

// UnsafeFileMsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FileMsServiceServer will
// result in compilation errors.
type UnsafeFileMsServiceServer interface {
	mustEmbedUnimplementedFileMsServiceServer()
}

func RegisterFileMsServiceServer(s grpc.ServiceRegistrar, srv FileMsServiceServer) {
	s.RegisterService(&FileMsService_ServiceDesc, srv)
}

func _FileMsService_CreateFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileMsServiceServer).CreateFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.FileMsService/CreateFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileMsServiceServer).CreateFile(ctx, req.(*CreateFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileMsService_UpdateFileContent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileMsServiceServer).UpdateFileContent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.FileMsService/UpdateFileContent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileMsServiceServer).UpdateFileContent(ctx, req.(*UpdateFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileMsService_GetFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileMsServiceServer).GetFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.FileMsService/GetFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileMsServiceServer).GetFile(ctx, req.(*FileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileMsService_GetFileContent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileMsServiceServer).GetFileContent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.FileMsService/GetFileContent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileMsServiceServer).GetFileContent(ctx, req.(*FileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileMsService_DelFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileMsServiceServer).DelFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.FileMsService/DelFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileMsServiceServer).DelFile(ctx, req.(*FileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileMsService_ListFileVersions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileMsServiceServer).ListFileVersions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.FileMsService/ListFileVersions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileMsServiceServer).ListFileVersions(ctx, req.(*FileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FileMsService_ServiceDesc is the grpc.ServiceDesc for FileMsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FileMsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.FileMsService",
	HandlerType: (*FileMsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateFile",
			Handler:    _FileMsService_CreateFile_Handler,
		},
		{
			MethodName: "UpdateFileContent",
			Handler:    _FileMsService_UpdateFileContent_Handler,
		},
		{
			MethodName: "GetFile",
			Handler:    _FileMsService_GetFile_Handler,
		},
		{
			MethodName: "GetFileContent",
			Handler:    _FileMsService_GetFileContent_Handler,
		},
		{
			MethodName: "DelFile",
			Handler:    _FileMsService_DelFile_Handler,
		},
		{
			MethodName: "ListFileVersions",
			Handler:    _FileMsService_ListFileVersions_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "filems.proto",
}
