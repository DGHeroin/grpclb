// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: msg.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Payload   []byte `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
	ErrorCode int32  `protobuf:"varint,3,opt,name=errorCode,proto3" json:"errorCode,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Message) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *Message) GetErrorCode() int32 {
	if x != nil {
		return x.ErrorCode
	}
	return 0
}

var File_msg_proto protoreflect.FileDescriptor

var file_msg_proto_rawDesc = []byte{
	0x0a, 0x09, 0x6d, 0x73, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x22,
	0x55, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x32, 0x67, 0x0a, 0x0e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x12, 0x25, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x0b, 0x2e, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x1a, 0x0b, 0x2e, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x00, 0x12,
	0x2e, 0x0a, 0x0c, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x50, 0x75, 0x73, 0x68, 0x12,
	0x0b, 0x2e, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x0b, 0x2e, 0x70,
	0x62, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42,
	0x06, 0x5a, 0x04, 0x2e, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_msg_proto_rawDescOnce sync.Once
	file_msg_proto_rawDescData = file_msg_proto_rawDesc
)

func file_msg_proto_rawDescGZIP() []byte {
	file_msg_proto_rawDescOnce.Do(func() {
		file_msg_proto_rawDescData = protoimpl.X.CompressGZIP(file_msg_proto_rawDescData)
	})
	return file_msg_proto_rawDescData
}

var file_msg_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_msg_proto_goTypes = []interface{}{
	(*Message)(nil), // 0: pb.Message
}
var file_msg_proto_depIdxs = []int32{
	0, // 0: pb.MessageHandler.Request:input_type -> pb.Message
	0, // 1: pb.MessageHandler.RegisterPush:input_type -> pb.Message
	0, // 2: pb.MessageHandler.Request:output_type -> pb.Message
	0, // 3: pb.MessageHandler.RegisterPush:output_type -> pb.Message
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_msg_proto_init() }
func file_msg_proto_init() {
	if File_msg_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_msg_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_msg_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_msg_proto_goTypes,
		DependencyIndexes: file_msg_proto_depIdxs,
		MessageInfos:      file_msg_proto_msgTypes,
	}.Build()
	File_msg_proto = out.File
	file_msg_proto_rawDesc = nil
	file_msg_proto_goTypes = nil
	file_msg_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// MessageHandlerClient is the client API for MessageHandler service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MessageHandlerClient interface {
	Request(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error)
	RegisterPush(ctx context.Context, opts ...grpc.CallOption) (MessageHandler_RegisterPushClient, error)
}

type messageHandlerClient struct {
	cc grpc.ClientConnInterface
}

func NewMessageHandlerClient(cc grpc.ClientConnInterface) MessageHandlerClient {
	return &messageHandlerClient{cc}
}

func (c *messageHandlerClient) Request(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, "/pb.MessageHandler/Request", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageHandlerClient) RegisterPush(ctx context.Context, opts ...grpc.CallOption) (MessageHandler_RegisterPushClient, error) {
	stream, err := c.cc.NewStream(ctx, &_MessageHandler_serviceDesc.Streams[0], "/pb.MessageHandler/RegisterPush", opts...)
	if err != nil {
		return nil, err
	}
	x := &messageHandlerRegisterPushClient{stream}
	return x, nil
}

type MessageHandler_RegisterPushClient interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ClientStream
}

type messageHandlerRegisterPushClient struct {
	grpc.ClientStream
}

func (x *messageHandlerRegisterPushClient) Send(m *Message) error {
	return x.ClientStream.SendMsg(m)
}

func (x *messageHandlerRegisterPushClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MessageHandlerServer is the server API for MessageHandler service.
type MessageHandlerServer interface {
	Request(context.Context, *Message) (*Message, error)
	RegisterPush(MessageHandler_RegisterPushServer) error
}

// UnimplementedMessageHandlerServer can be embedded to have forward compatible implementations.
type UnimplementedMessageHandlerServer struct {
}

func (*UnimplementedMessageHandlerServer) Request(context.Context, *Message) (*Message, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Request not implemented")
}
func (*UnimplementedMessageHandlerServer) RegisterPush(MessageHandler_RegisterPushServer) error {
	return status.Errorf(codes.Unimplemented, "method RegisterPush not implemented")
}

func RegisterMessageHandlerServer(s *grpc.Server, srv MessageHandlerServer) {
	s.RegisterService(&_MessageHandler_serviceDesc, srv)
}

func _MessageHandler_Request_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageHandlerServer).Request(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.MessageHandler/Request",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageHandlerServer).Request(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessageHandler_RegisterPush_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MessageHandlerServer).RegisterPush(&messageHandlerRegisterPushServer{stream})
}

type MessageHandler_RegisterPushServer interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ServerStream
}

type messageHandlerRegisterPushServer struct {
	grpc.ServerStream
}

func (x *messageHandlerRegisterPushServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func (x *messageHandlerRegisterPushServer) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _MessageHandler_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.MessageHandler",
	HandlerType: (*MessageHandlerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Request",
			Handler:    _MessageHandler_Request_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RegisterPush",
			Handler:       _MessageHandler_RegisterPush_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "msg.proto",
}