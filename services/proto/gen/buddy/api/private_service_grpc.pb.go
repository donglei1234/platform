// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: private_service.proto

package buddy

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

const (
	PrivateService_GetProfileBuddies_FullMethodName     = "/buddy.pb.PrivateService/GetProfileBuddies"
	PrivateService_WatchProfileBuddies_FullMethodName   = "/buddy.pb.PrivateService/WatchProfileBuddies"
	PrivateService_GetBuddies_FullMethodName            = "/buddy.pb.PrivateService/GetBuddies"
	PrivateService_WatchBuddies_FullMethodName          = "/buddy.pb.PrivateService/WatchBuddies"
	PrivateService_GetProfileBlockedList_FullMethodName = "/buddy.pb.PrivateService/GetProfileBlockedList"
)

// PrivateServiceClient is the client API for PrivateService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PrivateServiceClient interface {
	// GetProfileBuddies returns the provided user's buddies.
	GetProfileBuddies(ctx context.Context, in *GetProfileBuddiesRequest, opts ...grpc.CallOption) (*Buddies, error)
	// WatchProfileBuddies returns a stream on which changes to the provided user's buddies will be sent.
	WatchProfileBuddies(ctx context.Context, in *WatchProfileBuddiesRequest, opts ...grpc.CallOption) (PrivateService_WatchProfileBuddiesClient, error)
	// DEPRECATED
	// GetBuddies returns the provided user's buddies.
	GetBuddies(ctx context.Context, in *Name, opts ...grpc.CallOption) (*Buddies, error)
	// DEPRECATED
	// WatchBuddies returns a stream on which changes to the provided user's buddies will be sent.
	WatchBuddies(ctx context.Context, in *Name, opts ...grpc.CallOption) (PrivateService_WatchBuddiesClient, error)
	GetProfileBlockedList(ctx context.Context, in *GetProfileBlockedListRequest, opts ...grpc.CallOption) (*GetProfileBlockedListResponse, error)
}

type privateServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPrivateServiceClient(cc grpc.ClientConnInterface) PrivateServiceClient {
	return &privateServiceClient{cc}
}

func (c *privateServiceClient) GetProfileBuddies(ctx context.Context, in *GetProfileBuddiesRequest, opts ...grpc.CallOption) (*Buddies, error) {
	out := new(Buddies)
	err := c.cc.Invoke(ctx, PrivateService_GetProfileBuddies_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *privateServiceClient) WatchProfileBuddies(ctx context.Context, in *WatchProfileBuddiesRequest, opts ...grpc.CallOption) (PrivateService_WatchProfileBuddiesClient, error) {
	stream, err := c.cc.NewStream(ctx, &PrivateService_ServiceDesc.Streams[0], PrivateService_WatchProfileBuddies_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &privateServiceWatchProfileBuddiesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PrivateService_WatchProfileBuddiesClient interface {
	Recv() (*ROSUpdate, error)
	grpc.ClientStream
}

type privateServiceWatchProfileBuddiesClient struct {
	grpc.ClientStream
}

func (x *privateServiceWatchProfileBuddiesClient) Recv() (*ROSUpdate, error) {
	m := new(ROSUpdate)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *privateServiceClient) GetBuddies(ctx context.Context, in *Name, opts ...grpc.CallOption) (*Buddies, error) {
	out := new(Buddies)
	err := c.cc.Invoke(ctx, PrivateService_GetBuddies_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *privateServiceClient) WatchBuddies(ctx context.Context, in *Name, opts ...grpc.CallOption) (PrivateService_WatchBuddiesClient, error) {
	stream, err := c.cc.NewStream(ctx, &PrivateService_ServiceDesc.Streams[1], PrivateService_WatchBuddies_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &privateServiceWatchBuddiesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PrivateService_WatchBuddiesClient interface {
	Recv() (*ROSUpdate, error)
	grpc.ClientStream
}

type privateServiceWatchBuddiesClient struct {
	grpc.ClientStream
}

func (x *privateServiceWatchBuddiesClient) Recv() (*ROSUpdate, error) {
	m := new(ROSUpdate)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *privateServiceClient) GetProfileBlockedList(ctx context.Context, in *GetProfileBlockedListRequest, opts ...grpc.CallOption) (*GetProfileBlockedListResponse, error) {
	out := new(GetProfileBlockedListResponse)
	err := c.cc.Invoke(ctx, PrivateService_GetProfileBlockedList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PrivateServiceServer is the server API for PrivateService service.
// All implementations should embed UnimplementedPrivateServiceServer
// for forward compatibility
type PrivateServiceServer interface {
	// GetProfileBuddies returns the provided user's buddies.
	GetProfileBuddies(context.Context, *GetProfileBuddiesRequest) (*Buddies, error)
	// WatchProfileBuddies returns a stream on which changes to the provided user's buddies will be sent.
	WatchProfileBuddies(*WatchProfileBuddiesRequest, PrivateService_WatchProfileBuddiesServer) error
	// DEPRECATED
	// GetBuddies returns the provided user's buddies.
	GetBuddies(context.Context, *Name) (*Buddies, error)
	// DEPRECATED
	// WatchBuddies returns a stream on which changes to the provided user's buddies will be sent.
	WatchBuddies(*Name, PrivateService_WatchBuddiesServer) error
	GetProfileBlockedList(context.Context, *GetProfileBlockedListRequest) (*GetProfileBlockedListResponse, error)
}

// UnimplementedPrivateServiceServer should be embedded to have forward compatible implementations.
type UnimplementedPrivateServiceServer struct {
}

func (UnimplementedPrivateServiceServer) GetProfileBuddies(context.Context, *GetProfileBuddiesRequest) (*Buddies, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProfileBuddies not implemented")
}
func (UnimplementedPrivateServiceServer) WatchProfileBuddies(*WatchProfileBuddiesRequest, PrivateService_WatchProfileBuddiesServer) error {
	return status.Errorf(codes.Unimplemented, "method WatchProfileBuddies not implemented")
}
func (UnimplementedPrivateServiceServer) GetBuddies(context.Context, *Name) (*Buddies, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBuddies not implemented")
}
func (UnimplementedPrivateServiceServer) WatchBuddies(*Name, PrivateService_WatchBuddiesServer) error {
	return status.Errorf(codes.Unimplemented, "method WatchBuddies not implemented")
}
func (UnimplementedPrivateServiceServer) GetProfileBlockedList(context.Context, *GetProfileBlockedListRequest) (*GetProfileBlockedListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProfileBlockedList not implemented")
}

// UnsafePrivateServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PrivateServiceServer will
// result in compilation errors.
type UnsafePrivateServiceServer interface {
	mustEmbedUnimplementedPrivateServiceServer()
}

func RegisterPrivateServiceServer(s grpc.ServiceRegistrar, srv PrivateServiceServer) {
	s.RegisterService(&PrivateService_ServiceDesc, srv)
}

func _PrivateService_GetProfileBuddies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProfileBuddiesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PrivateServiceServer).GetProfileBuddies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PrivateService_GetProfileBuddies_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PrivateServiceServer).GetProfileBuddies(ctx, req.(*GetProfileBuddiesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PrivateService_WatchProfileBuddies_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WatchProfileBuddiesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PrivateServiceServer).WatchProfileBuddies(m, &privateServiceWatchProfileBuddiesServer{stream})
}

type PrivateService_WatchProfileBuddiesServer interface {
	Send(*ROSUpdate) error
	grpc.ServerStream
}

type privateServiceWatchProfileBuddiesServer struct {
	grpc.ServerStream
}

func (x *privateServiceWatchProfileBuddiesServer) Send(m *ROSUpdate) error {
	return x.ServerStream.SendMsg(m)
}

func _PrivateService_GetBuddies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Name)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PrivateServiceServer).GetBuddies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PrivateService_GetBuddies_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PrivateServiceServer).GetBuddies(ctx, req.(*Name))
	}
	return interceptor(ctx, in, info, handler)
}

func _PrivateService_WatchBuddies_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Name)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PrivateServiceServer).WatchBuddies(m, &privateServiceWatchBuddiesServer{stream})
}

type PrivateService_WatchBuddiesServer interface {
	Send(*ROSUpdate) error
	grpc.ServerStream
}

type privateServiceWatchBuddiesServer struct {
	grpc.ServerStream
}

func (x *privateServiceWatchBuddiesServer) Send(m *ROSUpdate) error {
	return x.ServerStream.SendMsg(m)
}

func _PrivateService_GetProfileBlockedList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProfileBlockedListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PrivateServiceServer).GetProfileBlockedList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PrivateService_GetProfileBlockedList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PrivateServiceServer).GetProfileBlockedList(ctx, req.(*GetProfileBlockedListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PrivateService_ServiceDesc is the grpc.ServiceDesc for PrivateService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PrivateService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "buddy.pb.PrivateService",
	HandlerType: (*PrivateServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetProfileBuddies",
			Handler:    _PrivateService_GetProfileBuddies_Handler,
		},
		{
			MethodName: "GetBuddies",
			Handler:    _PrivateService_GetBuddies_Handler,
		},
		{
			MethodName: "GetProfileBlockedList",
			Handler:    _PrivateService_GetProfileBlockedList_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "WatchProfileBuddies",
			Handler:       _PrivateService_WatchProfileBuddies_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "WatchBuddies",
			Handler:       _PrivateService_WatchBuddies_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "private_service.proto",
}
