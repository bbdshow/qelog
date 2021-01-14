// Code generated by protoc-gen-go. DO NOT EDIT.
// source: receiver.proto

// protoc --go_out=plugins=grpc:.  *.proto

package receiverpb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type BaseResp struct {
	Code                 int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BaseResp) Reset()         { *m = BaseResp{} }
func (m *BaseResp) String() string { return proto.CompactTextString(m) }
func (*BaseResp) ProtoMessage()    {}
func (*BaseResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_4b7296e1d2b388c5, []int{0}
}

func (m *BaseResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BaseResp.Unmarshal(m, b)
}
func (m *BaseResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BaseResp.Marshal(b, m, deterministic)
}
func (m *BaseResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BaseResp.Merge(m, src)
}
func (m *BaseResp) XXX_Size() int {
	return xxx_messageInfo_BaseResp.Size(m)
}
func (m *BaseResp) XXX_DiscardUnknown() {
	xxx_messageInfo_BaseResp.DiscardUnknown(m)
}

var xxx_messageInfo_BaseResp proto.InternalMessageInfo

func (m *BaseResp) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *BaseResp) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type Packet struct {
	// 包ID
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	// 包属于哪一个模块
	Module string `protobuf:"bytes,2,opt,name=module,proto3" json:"module"`
	// 包具体的数据
	Data                 []byte   `protobuf:"bytes,3,opt,name=data,proto3" json:"data"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Packet) Reset()         { *m = Packet{} }
func (m *Packet) String() string { return proto.CompactTextString(m) }
func (*Packet) ProtoMessage()    {}
func (*Packet) Descriptor() ([]byte, []int) {
	return fileDescriptor_4b7296e1d2b388c5, []int{1}
}

func (m *Packet) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Packet.Unmarshal(m, b)
}
func (m *Packet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Packet.Marshal(b, m, deterministic)
}
func (m *Packet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Packet.Merge(m, src)
}
func (m *Packet) XXX_Size() int {
	return xxx_messageInfo_Packet.Size(m)
}
func (m *Packet) XXX_DiscardUnknown() {
	xxx_messageInfo_Packet.DiscardUnknown(m)
}

var xxx_messageInfo_Packet proto.InternalMessageInfo

func (m *Packet) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Packet) GetModule() string {
	if m != nil {
		return m.Module
	}
	return ""
}

func (m *Packet) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*BaseResp)(nil), "receiverpb.BaseResp")
	proto.RegisterType((*Packet)(nil), "receiverpb.Packet")
}

func init() { proto.RegisterFile("receiver.proto", fileDescriptor_4b7296e1d2b388c5) }

var fileDescriptor_4b7296e1d2b388c5 = []byte{
	// 180 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x8f, 0x41, 0xcb, 0x82, 0x40,
	0x10, 0x86, 0x59, 0xbf, 0x2f, 0xd3, 0x21, 0x3c, 0x0c, 0x11, 0x4b, 0x27, 0xf1, 0xe4, 0xc9, 0x43,
	0x41, 0x74, 0x96, 0x7e, 0x80, 0xcc, 0x3f, 0x58, 0xdd, 0xa1, 0xa4, 0x64, 0xc5, 0xd5, 0x7e, 0x7f,
	0xb4, 0xed, 0x52, 0xb7, 0x79, 0x5e, 0x78, 0x5f, 0x9e, 0x81, 0x6c, 0xe2, 0x8e, 0xfb, 0x27, 0x4f,
	0xd5, 0x38, 0x99, 0xd9, 0x20, 0x04, 0x1e, 0xdb, 0xe2, 0x0c, 0x49, 0xad, 0x2c, 0x13, 0xdb, 0x11,
	0x11, 0xfe, 0x3b, 0xa3, 0x59, 0x8a, 0x5c, 0x94, 0x2b, 0x72, 0x37, 0x4a, 0x58, 0x0f, 0x6c, 0xad,
	0xba, 0xb2, 0x8c, 0x72, 0x51, 0xa6, 0x14, 0xb0, 0xb8, 0x40, 0xdc, 0xa8, 0xee, 0xce, 0x33, 0x66,
	0x10, 0xf5, 0xda, 0xb5, 0x52, 0x8a, 0x7a, 0x8d, 0x3b, 0x88, 0x07, 0xa3, 0x97, 0x47, 0xa8, 0x78,
	0x7a, 0xef, 0x6b, 0x35, 0x2b, 0xf9, 0x97, 0x8b, 0x72, 0x43, 0xee, 0x3e, 0xd4, 0x90, 0x90, 0xb7,
	0xc1, 0x13, 0x40, 0xb3, 0xd8, 0x9b, 0x5f, 0xc5, 0xea, 0xab, 0x59, 0x7d, 0xb2, 0xfd, 0xf6, 0x37,
	0x0b, 0xde, 0x6d, 0xec, 0xde, 0x3a, 0xbe, 0x02, 0x00, 0x00, 0xff, 0xff, 0x4c, 0xa4, 0xa5, 0x06,
	0xe8, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ReceiverClient is the client API for Receiver service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ReceiverClient interface {
	// 推送数据包
	PushPacket(ctx context.Context, in *Packet, opts ...grpc.CallOption) (*BaseResp, error)
}

type receiverClient struct {
	cc *grpc.ClientConn
}

func NewReceiverClient(cc *grpc.ClientConn) ReceiverClient {
	return &receiverClient{cc}
}

func (c *receiverClient) PushPacket(ctx context.Context, in *Packet, opts ...grpc.CallOption) (*BaseResp, error) {
	out := new(BaseResp)
	err := c.cc.Invoke(ctx, "/receiverpb.Receiver/PushPacket", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReceiverServer is the server API for Receiver service.
type ReceiverServer interface {
	// 推送数据包
	PushPacket(context.Context, *Packet) (*BaseResp, error)
}

// UnimplementedReceiverServer can be embedded to have forward compatible implementations.
type UnimplementedReceiverServer struct {
}

func (*UnimplementedReceiverServer) PushPacket(ctx context.Context, req *Packet) (*BaseResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushPacket not implemented")
}

func RegisterReceiverServer(s *grpc.Server, srv ReceiverServer) {
	s.RegisterService(&_Receiver_serviceDesc, srv)
}

func _Receiver_PushPacket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Packet)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReceiverServer).PushPacket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/receiverpb.Receiver/PushPacket",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReceiverServer).PushPacket(ctx, req.(*Packet))
	}
	return interceptor(ctx, in, info, handler)
}

var _Receiver_serviceDesc = grpc.ServiceDesc{
	ServiceName: "receiverpb.Receiver",
	HandlerType: (*ReceiverServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PushPacket",
			Handler:    _Receiver_PushPacket_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "receiver.proto",
}