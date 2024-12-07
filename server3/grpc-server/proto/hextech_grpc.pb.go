// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.3
// source: hextech.proto

package hextech

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Broker_ProcessCommand_FullMethodName     = "/hextech.Broker/ProcessCommand"
	Broker_ObtenerProducto_FullMethodName    = "/hextech.Broker/ObtenerProducto"
	Broker_GetClockByRegion_FullMethodName   = "/hextech.Broker/GetClockByRegion"
	Broker_GetClock_FullMethodName           = "/hextech.Broker/GetClock"
	Broker_DistributeLogs_FullMethodName     = "/hextech.Broker/DistributeLogs"
	Broker_DistributeRegions_FullMethodName  = "/hextech.Broker/DistributeRegions"
	Broker_GetFile_FullMethodName            = "/hextech.Broker/GetFile"
	Broker_ReceiveMergedFile_FullMethodName  = "/hextech.Broker/ReceiveMergedFile"
	Broker_ReceiveMergedFile2_FullMethodName = "/hextech.Broker/ReceiveMergedFile2"
)

// BrokerClient is the client API for Broker service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BrokerClient interface {
	ProcessCommand(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (*CommandResponse, error)
	ObtenerProducto(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (*CommandResponse, error)
	GetClockByRegion(ctx context.Context, in *RegionRequest, opts ...grpc.CallOption) (*ClockResponse, error)
	GetClock(ctx context.Context, in *ClockRequest, opts ...grpc.CallOption) (*ClockResponse, error)
	DistributeLogs(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (*MergeResponse, error)
	DistributeRegions(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*MergeResponse, error)
	GetFile(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileResponse, error)
	ReceiveMergedFile(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileResponse, error)
	ReceiveMergedFile2(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileResponse, error)
}

type brokerClient struct {
	cc grpc.ClientConnInterface
}

func NewBrokerClient(cc grpc.ClientConnInterface) BrokerClient {
	return &brokerClient{cc}
}

func (c *brokerClient) ProcessCommand(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (*CommandResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CommandResponse)
	err := c.cc.Invoke(ctx, Broker_ProcessCommand_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerClient) ObtenerProducto(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (*CommandResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CommandResponse)
	err := c.cc.Invoke(ctx, Broker_ObtenerProducto_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerClient) GetClockByRegion(ctx context.Context, in *RegionRequest, opts ...grpc.CallOption) (*ClockResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ClockResponse)
	err := c.cc.Invoke(ctx, Broker_GetClockByRegion_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerClient) GetClock(ctx context.Context, in *ClockRequest, opts ...grpc.CallOption) (*ClockResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ClockResponse)
	err := c.cc.Invoke(ctx, Broker_GetClock_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerClient) DistributeLogs(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (*MergeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MergeResponse)
	err := c.cc.Invoke(ctx, Broker_DistributeLogs_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerClient) DistributeRegions(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*MergeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MergeResponse)
	err := c.cc.Invoke(ctx, Broker_DistributeRegions_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerClient) GetFile(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FileResponse)
	err := c.cc.Invoke(ctx, Broker_GetFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerClient) ReceiveMergedFile(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FileResponse)
	err := c.cc.Invoke(ctx, Broker_ReceiveMergedFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerClient) ReceiveMergedFile2(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FileResponse)
	err := c.cc.Invoke(ctx, Broker_ReceiveMergedFile2_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BrokerServer is the server API for Broker service.
// All implementations must embed UnimplementedBrokerServer
// for forward compatibility.
type BrokerServer interface {
	ProcessCommand(context.Context, *CommandRequest) (*CommandResponse, error)
	ObtenerProducto(context.Context, *CommandRequest) (*CommandResponse, error)
	GetClockByRegion(context.Context, *RegionRequest) (*ClockResponse, error)
	GetClock(context.Context, *ClockRequest) (*ClockResponse, error)
	DistributeLogs(context.Context, *LogRequest) (*MergeResponse, error)
	DistributeRegions(context.Context, *FileRequest) (*MergeResponse, error)
	GetFile(context.Context, *FileRequest) (*FileResponse, error)
	ReceiveMergedFile(context.Context, *FileRequest) (*FileResponse, error)
	ReceiveMergedFile2(context.Context, *FileRequest) (*FileResponse, error)
	mustEmbedUnimplementedBrokerServer()
}

// UnimplementedBrokerServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedBrokerServer struct{}

func (UnimplementedBrokerServer) ProcessCommand(context.Context, *CommandRequest) (*CommandResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessCommand not implemented")
}
func (UnimplementedBrokerServer) ObtenerProducto(context.Context, *CommandRequest) (*CommandResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ObtenerProducto not implemented")
}
func (UnimplementedBrokerServer) GetClockByRegion(context.Context, *RegionRequest) (*ClockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClockByRegion not implemented")
}
func (UnimplementedBrokerServer) GetClock(context.Context, *ClockRequest) (*ClockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClock not implemented")
}
func (UnimplementedBrokerServer) DistributeLogs(context.Context, *LogRequest) (*MergeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DistributeLogs not implemented")
}
func (UnimplementedBrokerServer) DistributeRegions(context.Context, *FileRequest) (*MergeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DistributeRegions not implemented")
}
func (UnimplementedBrokerServer) GetFile(context.Context, *FileRequest) (*FileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFile not implemented")
}
func (UnimplementedBrokerServer) ReceiveMergedFile(context.Context, *FileRequest) (*FileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReceiveMergedFile not implemented")
}
func (UnimplementedBrokerServer) ReceiveMergedFile2(context.Context, *FileRequest) (*FileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReceiveMergedFile2 not implemented")
}
func (UnimplementedBrokerServer) mustEmbedUnimplementedBrokerServer() {}
func (UnimplementedBrokerServer) testEmbeddedByValue()                {}

// UnsafeBrokerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BrokerServer will
// result in compilation errors.
type UnsafeBrokerServer interface {
	mustEmbedUnimplementedBrokerServer()
}

func RegisterBrokerServer(s grpc.ServiceRegistrar, srv BrokerServer) {
	// If the following call pancis, it indicates UnimplementedBrokerServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Broker_ServiceDesc, srv)
}

func _Broker_ProcessCommand_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommandRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).ProcessCommand(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Broker_ProcessCommand_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).ProcessCommand(ctx, req.(*CommandRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broker_ObtenerProducto_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommandRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).ObtenerProducto(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Broker_ObtenerProducto_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).ObtenerProducto(ctx, req.(*CommandRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broker_GetClockByRegion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).GetClockByRegion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Broker_GetClockByRegion_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).GetClockByRegion(ctx, req.(*RegionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broker_GetClock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).GetClock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Broker_GetClock_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).GetClock(ctx, req.(*ClockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broker_DistributeLogs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).DistributeLogs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Broker_DistributeLogs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).DistributeLogs(ctx, req.(*LogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broker_DistributeRegions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).DistributeRegions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Broker_DistributeRegions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).DistributeRegions(ctx, req.(*FileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broker_GetFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).GetFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Broker_GetFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).GetFile(ctx, req.(*FileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broker_ReceiveMergedFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).ReceiveMergedFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Broker_ReceiveMergedFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).ReceiveMergedFile(ctx, req.(*FileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broker_ReceiveMergedFile2_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).ReceiveMergedFile2(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Broker_ReceiveMergedFile2_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).ReceiveMergedFile2(ctx, req.(*FileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Broker_ServiceDesc is the grpc.ServiceDesc for Broker service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Broker_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "hextech.Broker",
	HandlerType: (*BrokerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ProcessCommand",
			Handler:    _Broker_ProcessCommand_Handler,
		},
		{
			MethodName: "ObtenerProducto",
			Handler:    _Broker_ObtenerProducto_Handler,
		},
		{
			MethodName: "GetClockByRegion",
			Handler:    _Broker_GetClockByRegion_Handler,
		},
		{
			MethodName: "GetClock",
			Handler:    _Broker_GetClock_Handler,
		},
		{
			MethodName: "DistributeLogs",
			Handler:    _Broker_DistributeLogs_Handler,
		},
		{
			MethodName: "DistributeRegions",
			Handler:    _Broker_DistributeRegions_Handler,
		},
		{
			MethodName: "GetFile",
			Handler:    _Broker_GetFile_Handler,
		},
		{
			MethodName: "ReceiveMergedFile",
			Handler:    _Broker_ReceiveMergedFile_Handler,
		},
		{
			MethodName: "ReceiveMergedFile2",
			Handler:    _Broker_ReceiveMergedFile2_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "hextech.proto",
}
