// Code generated by protoc-gen-go. DO NOT EDIT.
// source: push.proto

package push

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
	return fileDescriptor_d1e4bfd2e9d102bb, []int{0}
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
	Data []string `protobuf:"bytes,3,rep,name=data,proto3" json:"data"`
	// 是否是重试发送的包
	IsRetry              bool     `protobuf:"varint,4,opt,name=isRetry,proto3" json:"isRetry"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Packet) Reset()         { *m = Packet{} }
func (m *Packet) String() string { return proto.CompactTextString(m) }
func (*Packet) ProtoMessage()    {}
func (*Packet) Descriptor() ([]byte, []int) {
	return fileDescriptor_d1e4bfd2e9d102bb, []int{1}
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

func (m *Packet) GetData() []string {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Packet) GetIsRetry() bool {
	if m != nil {
		return m.IsRetry
	}
	return false
}

func init() {
	proto.RegisterType((*BaseResp)(nil), "push.BaseResp")
	proto.RegisterType((*Packet)(nil), "push.Packet")
}

func init() { proto.RegisterFile("push.proto", fileDescriptor_d1e4bfd2e9d102bb) }

var fileDescriptor_d1e4bfd2e9d102bb = []byte{
	// 186 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x34, 0x8f, 0xb1, 0xaa, 0x83, 0x30,
	0x14, 0x86, 0x89, 0xe6, 0x7a, 0xf5, 0x70, 0x71, 0xc8, 0x70, 0x09, 0x9d, 0x82, 0x53, 0xe8, 0xe0,
	0x60, 0x97, 0xce, 0x7d, 0x02, 0xc9, 0x03, 0x14, 0x52, 0x13, 0xaa, 0xb4, 0x12, 0xf1, 0xc4, 0xa1,
	0x6f, 0x5f, 0x12, 0xcd, 0x74, 0xfe, 0x6f, 0xf8, 0xff, 0xc3, 0x07, 0xb0, 0x6c, 0x38, 0xb6, 0xcb,
	0xea, 0xbc, 0x63, 0x34, 0xe4, 0xe6, 0x0a, 0xe5, 0x4d, 0xa3, 0x55, 0x16, 0x17, 0xc6, 0x80, 0x0e,
	0xce, 0x58, 0x4e, 0x04, 0x91, 0x3f, 0x2a, 0x66, 0xc6, 0xe1, 0x77, 0xb6, 0x88, 0xfa, 0x69, 0x79,
	0x26, 0x88, 0xac, 0x54, 0xc2, 0xe6, 0x0e, 0x45, 0xaf, 0x87, 0x97, 0xf5, 0xac, 0x86, 0x6c, 0x32,
	0xb1, 0x55, 0xa9, 0x6c, 0x32, 0xec, 0x1f, 0x8a, 0xd9, 0x99, 0xed, 0x9d, 0x2a, 0x07, 0x85, 0x7d,
	0xa3, 0xbd, 0xe6, 0xb9, 0xc8, 0x65, 0xa5, 0x62, 0x0e, 0xfb, 0x13, 0x2a, 0xeb, 0xd7, 0x0f, 0xa7,
	0x82, 0xc8, 0x52, 0x25, 0xec, 0x3a, 0xa0, 0xfd, 0x86, 0x23, 0x3b, 0x03, 0x84, 0x7b, 0xfc, 0xfa,
	0x6b, 0xa3, 0xc2, 0x4e, 0xa7, 0x7a, 0xa7, 0x64, 0xf0, 0x28, 0xa2, 0xda, 0xe5, 0x1b, 0x00, 0x00,
	0xff, 0xff, 0xee, 0xb1, 0xa9, 0x74, 0xe8, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// PushClient is the client API for Push service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PushClient interface {
	// 推送数据包
	PushPacket(ctx context.Context, in *Packet, opts ...grpc.CallOption) (*BaseResp, error)
}

type pushClient struct {
	cc *grpc.ClientConn
}

func NewPushClient(cc *grpc.ClientConn) PushClient {
	return &pushClient{cc}
}

func (c *pushClient) PushPacket(ctx context.Context, in *Packet, opts ...grpc.CallOption) (*BaseResp, error) {
	out := new(BaseResp)
	err := c.cc.Invoke(ctx, "/push.Push/PushPacket", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PushServer is the server API for Push service.
type PushServer interface {
	// 推送数据包
	PushPacket(context.Context, *Packet) (*BaseResp, error)
}

// UnimplementedPushServer can be embedded to have forward compatible implementations.
type UnimplementedPushServer struct {
}

func (*UnimplementedPushServer) PushPacket(ctx context.Context, req *Packet) (*BaseResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushPacket not implemented")
}

func RegisterPushServer(s *grpc.Server, srv PushServer) {
	s.RegisterService(&_Push_serviceDesc, srv)
}

func _Push_PushPacket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Packet)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PushServer).PushPacket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/push.Push/PushPacket",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PushServer).PushPacket(ctx, req.(*Packet))
	}
	return interceptor(ctx, in, info, handler)
}

var _Push_serviceDesc = grpc.ServiceDesc{
	ServiceName: "push.Push",
	HandlerType: (*PushServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PushPacket",
			Handler:    _Push_PushPacket_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "push.proto",
}