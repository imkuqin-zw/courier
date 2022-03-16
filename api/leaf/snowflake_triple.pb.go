// Code generated by protoc-gen-go-triple. DO NOT EDIT.
// versions:
// - protoc-gen-go-triple v1.0.5
// - protoc             v3.12.3
// source: snowflake.proto

package pbLeaf

import (
	context "context"
	protocol "dubbo.apache.org/dubbo-go/v3/protocol"
	dubbo3 "dubbo.apache.org/dubbo-go/v3/protocol/dubbo3"
	invocation "dubbo.apache.org/dubbo-go/v3/protocol/invocation"
	grpc_go "github.com/dubbogo/grpc-go"
	codes "github.com/dubbogo/grpc-go/codes"
	metadata "github.com/dubbogo/grpc-go/metadata"
	status "github.com/dubbogo/grpc-go/status"
	common "github.com/dubbogo/triple/pkg/common"
	constant "github.com/dubbogo/triple/pkg/common/constant"
	triple "github.com/dubbogo/triple/pkg/triple"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc_go.SupportPackageIsVersion7

// SnowflakeClient is the client API for Snowflake service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SnowflakeClient interface {
	FetchNext(ctx context.Context, in *SnowflakeFetchNextReq, opts ...grpc_go.CallOption) (*SnowflakeFetchNextRes, common.ErrorWithAttachment)
}

type snowflakeClient struct {
	cc *triple.TripleConn
}

type SnowflakeClientImpl struct {
	FetchNext func(ctx context.Context, in *SnowflakeFetchNextReq) (*SnowflakeFetchNextRes, error)
}

func (c *SnowflakeClientImpl) GetDubboStub(cc *triple.TripleConn) SnowflakeClient {
	return NewSnowflakeClient(cc)
}

func (c *SnowflakeClientImpl) XXX_InterfaceName() string {
	return "com.github.imkuqin_zw.courier.api.leaf.Snowflake"
}

func NewSnowflakeClient(cc *triple.TripleConn) SnowflakeClient {
	return &snowflakeClient{cc}
}

func (c *snowflakeClient) FetchNext(ctx context.Context, in *SnowflakeFetchNextReq, opts ...grpc_go.CallOption) (*SnowflakeFetchNextRes, common.ErrorWithAttachment) {
	out := new(SnowflakeFetchNextRes)
	interfaceKey := ctx.Value(constant.InterfaceKey).(string)
	return out, c.cc.Invoke(ctx, "/"+interfaceKey+"/FetchNext", in, out)
}

// SnowflakeServer is the server API for Snowflake service.
// All implementations must embed UnimplementedSnowflakeServer
// for forward compatibility
type SnowflakeServer interface {
	FetchNext(context.Context, *SnowflakeFetchNextReq) (*SnowflakeFetchNextRes, error)
	mustEmbedUnimplementedSnowflakeServer()
}

// UnimplementedSnowflakeServer must be embedded to have forward compatible implementations.
type UnimplementedSnowflakeServer struct {
	proxyImpl protocol.Invoker
}

func (UnimplementedSnowflakeServer) FetchNext(context.Context, *SnowflakeFetchNextReq) (*SnowflakeFetchNextRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FetchNext not implemented")
}
func (s *UnimplementedSnowflakeServer) XXX_SetProxyImpl(impl protocol.Invoker) {
	s.proxyImpl = impl
}

func (s *UnimplementedSnowflakeServer) XXX_GetProxyImpl() protocol.Invoker {
	return s.proxyImpl
}

func (s *UnimplementedSnowflakeServer) XXX_ServiceDesc() *grpc_go.ServiceDesc {
	return &Snowflake_ServiceDesc
}
func (s *UnimplementedSnowflakeServer) XXX_InterfaceName() string {
	return "com.github.imkuqin_zw.courier.api.leaf.Snowflake"
}

func (UnimplementedSnowflakeServer) mustEmbedUnimplementedSnowflakeServer() {}

// UnsafeSnowflakeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SnowflakeServer will
// result in compilation errors.
type UnsafeSnowflakeServer interface {
	mustEmbedUnimplementedSnowflakeServer()
}

func RegisterSnowflakeServer(s grpc_go.ServiceRegistrar, srv SnowflakeServer) {
	s.RegisterService(&Snowflake_ServiceDesc, srv)
}

func _Snowflake_FetchNext_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc_go.UnaryServerInterceptor) (interface{}, error) {
	in := new(SnowflakeFetchNextReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	base := srv.(dubbo3.Dubbo3GrpcService)
	args := []interface{}{}
	args = append(args, in)
	md, _ := metadata.FromIncomingContext(ctx)
	invAttachment := make(map[string]interface{}, len(md))
	for k, v := range md {
		invAttachment[k] = v
	}
	invo := invocation.NewRPCInvocation("FetchNext", args, invAttachment)
	if interceptor == nil {
		result := base.XXX_GetProxyImpl().Invoke(ctx, invo)
		return result, result.Error()
	}
	info := &grpc_go.UnaryServerInfo{
		Server:     srv,
		FullMethod: ctx.Value("XXX_TRIPLE_GO_INTERFACE_NAME").(string),
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		result := base.XXX_GetProxyImpl().Invoke(ctx, invo)
		return result, result.Error()
	}
	return interceptor(ctx, in, info, handler)
}

// Snowflake_ServiceDesc is the grpc_go.ServiceDesc for Snowflake service.
// It's only intended for direct use with grpc_go.RegisterService,
// and not to be introspected or modified (even as a copy)
var Snowflake_ServiceDesc = grpc_go.ServiceDesc{
	ServiceName: "com.github.imkuqin_zw.courier.api.leaf.Snowflake",
	HandlerType: (*SnowflakeServer)(nil),
	Methods: []grpc_go.MethodDesc{
		{
			MethodName: "FetchNext",
			Handler:    _Snowflake_FetchNext_Handler,
		},
	},
	Streams:  []grpc_go.StreamDesc{},
	Metadata: "snowflake.proto",
}