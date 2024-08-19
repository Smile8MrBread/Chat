// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.27.1
// source: chat.proto

package chatGrpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	Chat_AddContact_FullMethodName   = "/chat.Chat/AddContact"
	Chat_AllContacts_FullMethodName  = "/chat.Chat/AllContacts"
	Chat_IsMessaged_FullMethodName   = "/chat.Chat/IsMessaged"
	Chat_AllMessaged_FullMethodName  = "/chat.Chat/AllMessaged"
	Chat_AllMessages_FullMethodName  = "/chat.Chat/AllMessages"
	Chat_IdentMessage_FullMethodName = "/chat.Chat/IdentMessage"
)

// ChatClient is the client API for Chat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatClient interface {
	AddContact(ctx context.Context, in *AddContactRequest, opts ...grpc.CallOption) (*Nothing, error)
	// rpc DeleteContact(DeleteContactRequest) returns (Nothing);
	AllContacts(ctx context.Context, in *AllContactsRequest, opts ...grpc.CallOption) (*AllContactsResponse, error)
	IsMessaged(ctx context.Context, in *IsMessagedRequest, opts ...grpc.CallOption) (*Nothing, error)
	AllMessaged(ctx context.Context, in *AllMessagedRequest, opts ...grpc.CallOption) (*AllMessagedResponse, error)
	// rpc CreateMessage(CreateMessageRequest) returns (CreateMessageResponse);
	// rpc UpdateMessage(UpdateMessageRequest) returns (UpdateMessageResponse);
	AllMessages(ctx context.Context, in *AllMessagesRequest, opts ...grpc.CallOption) (*AllMessagesResponse, error)
	IdentMessage(ctx context.Context, in *IdentMessageRequest, opts ...grpc.CallOption) (*IdentMessageResponse, error)
}

type chatClient struct {
	cc grpc.ClientConnInterface
}

func NewChatClient(cc grpc.ClientConnInterface) ChatClient {
	return &chatClient{cc}
}

func (c *chatClient) AddContact(ctx context.Context, in *AddContactRequest, opts ...grpc.CallOption) (*Nothing, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Nothing)
	err := c.cc.Invoke(ctx, Chat_AddContact_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) AllContacts(ctx context.Context, in *AllContactsRequest, opts ...grpc.CallOption) (*AllContactsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AllContactsResponse)
	err := c.cc.Invoke(ctx, Chat_AllContacts_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) IsMessaged(ctx context.Context, in *IsMessagedRequest, opts ...grpc.CallOption) (*Nothing, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Nothing)
	err := c.cc.Invoke(ctx, Chat_IsMessaged_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) AllMessaged(ctx context.Context, in *AllMessagedRequest, opts ...grpc.CallOption) (*AllMessagedResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AllMessagedResponse)
	err := c.cc.Invoke(ctx, Chat_AllMessaged_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) AllMessages(ctx context.Context, in *AllMessagesRequest, opts ...grpc.CallOption) (*AllMessagesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AllMessagesResponse)
	err := c.cc.Invoke(ctx, Chat_AllMessages_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) IdentMessage(ctx context.Context, in *IdentMessageRequest, opts ...grpc.CallOption) (*IdentMessageResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(IdentMessageResponse)
	err := c.cc.Invoke(ctx, Chat_IdentMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatServer is the server API for Chat service.
// All implementations must embed UnimplementedChatServer
// for forward compatibility
type ChatServer interface {
	AddContact(context.Context, *AddContactRequest) (*Nothing, error)
	// rpc DeleteContact(DeleteContactRequest) returns (Nothing);
	AllContacts(context.Context, *AllContactsRequest) (*AllContactsResponse, error)
	IsMessaged(context.Context, *IsMessagedRequest) (*Nothing, error)
	AllMessaged(context.Context, *AllMessagedRequest) (*AllMessagedResponse, error)
	// rpc CreateMessage(CreateMessageRequest) returns (CreateMessageResponse);
	// rpc UpdateMessage(UpdateMessageRequest) returns (UpdateMessageResponse);
	AllMessages(context.Context, *AllMessagesRequest) (*AllMessagesResponse, error)
	IdentMessage(context.Context, *IdentMessageRequest) (*IdentMessageResponse, error)
	mustEmbedUnimplementedChatServer()
}

// UnimplementedChatServer must be embedded to have forward compatible implementations.
type UnimplementedChatServer struct {
}

func (UnimplementedChatServer) AddContact(context.Context, *AddContactRequest) (*Nothing, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddContact not implemented")
}
func (UnimplementedChatServer) AllContacts(context.Context, *AllContactsRequest) (*AllContactsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AllContacts not implemented")
}
func (UnimplementedChatServer) IsMessaged(context.Context, *IsMessagedRequest) (*Nothing, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsMessaged not implemented")
}
func (UnimplementedChatServer) AllMessaged(context.Context, *AllMessagedRequest) (*AllMessagedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AllMessaged not implemented")
}
func (UnimplementedChatServer) AllMessages(context.Context, *AllMessagesRequest) (*AllMessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AllMessages not implemented")
}
func (UnimplementedChatServer) IdentMessage(context.Context, *IdentMessageRequest) (*IdentMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IdentMessage not implemented")
}
func (UnimplementedChatServer) mustEmbedUnimplementedChatServer() {}

// UnsafeChatServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatServer will
// result in compilation errors.
type UnsafeChatServer interface {
	mustEmbedUnimplementedChatServer()
}

func RegisterChatServer(s grpc.ServiceRegistrar, srv ChatServer) {
	s.RegisterService(&Chat_ServiceDesc, srv)
}

func _Chat_AddContact_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddContactRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).AddContact(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chat_AddContact_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).AddContact(ctx, req.(*AddContactRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_AllContacts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AllContactsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).AllContacts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chat_AllContacts_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).AllContacts(ctx, req.(*AllContactsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_IsMessaged_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsMessagedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).IsMessaged(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chat_IsMessaged_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).IsMessaged(ctx, req.(*IsMessagedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_AllMessaged_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AllMessagedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).AllMessaged(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chat_AllMessaged_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).AllMessaged(ctx, req.(*AllMessagedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_AllMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AllMessagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).AllMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chat_AllMessages_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).AllMessages(ctx, req.(*AllMessagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_IdentMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdentMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).IdentMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chat_IdentMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).IdentMessage(ctx, req.(*IdentMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Chat_ServiceDesc is the grpc.ServiceDesc for Chat service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Chat_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chat.Chat",
	HandlerType: (*ChatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddContact",
			Handler:    _Chat_AddContact_Handler,
		},
		{
			MethodName: "AllContacts",
			Handler:    _Chat_AllContacts_Handler,
		},
		{
			MethodName: "IsMessaged",
			Handler:    _Chat_IsMessaged_Handler,
		},
		{
			MethodName: "AllMessaged",
			Handler:    _Chat_AllMessaged_Handler,
		},
		{
			MethodName: "AllMessages",
			Handler:    _Chat_AllMessages_Handler,
		},
		{
			MethodName: "IdentMessage",
			Handler:    _Chat_IdentMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chat.proto",
}
