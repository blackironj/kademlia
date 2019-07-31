// Code generated by protoc-gen-go. DO NOT EDIT.
// source: kademlia_service.proto

package kademlia

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

type Target struct {
	TargetId             string   `protobuf:"bytes,1,opt,name=target_id,json=targetId,proto3" json:"target_id,omitempty"`
	SenderId             string   `protobuf:"bytes,2,opt,name=sender_id,json=senderId,proto3" json:"sender_id,omitempty"`
	SenderIp             string   `protobuf:"bytes,3,opt,name=sender_ip,json=senderIp,proto3" json:"sender_ip,omitempty"`
	SenderKadPort        string   `protobuf:"bytes,4,opt,name=sender_kad_port,json=senderKadPort,proto3" json:"sender_kad_port,omitempty"`
	SenderServPort       string   `protobuf:"bytes,5,opt,name=sender_serv_port,json=senderServPort,proto3" json:"sender_serv_port,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Target) Reset()         { *m = Target{} }
func (m *Target) String() string { return proto.CompactTextString(m) }
func (*Target) ProtoMessage()    {}
func (*Target) Descriptor() ([]byte, []int) {
	return fileDescriptor_9e023dba0ca8fba6, []int{0}
}

func (m *Target) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Target.Unmarshal(m, b)
}
func (m *Target) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Target.Marshal(b, m, deterministic)
}
func (m *Target) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Target.Merge(m, src)
}
func (m *Target) XXX_Size() int {
	return xxx_messageInfo_Target.Size(m)
}
func (m *Target) XXX_DiscardUnknown() {
	xxx_messageInfo_Target.DiscardUnknown(m)
}

var xxx_messageInfo_Target proto.InternalMessageInfo

func (m *Target) GetTargetId() string {
	if m != nil {
		return m.TargetId
	}
	return ""
}

func (m *Target) GetSenderId() string {
	if m != nil {
		return m.SenderId
	}
	return ""
}

func (m *Target) GetSenderIp() string {
	if m != nil {
		return m.SenderIp
	}
	return ""
}

func (m *Target) GetSenderKadPort() string {
	if m != nil {
		return m.SenderKadPort
	}
	return ""
}

func (m *Target) GetSenderServPort() string {
	if m != nil {
		return m.SenderServPort
	}
	return ""
}

type Neighbors struct {
	Neighbors            []*NeighborInfo `protobuf:"bytes,1,rep,name=neighbors,proto3" json:"neighbors,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Neighbors) Reset()         { *m = Neighbors{} }
func (m *Neighbors) String() string { return proto.CompactTextString(m) }
func (*Neighbors) ProtoMessage()    {}
func (*Neighbors) Descriptor() ([]byte, []int) {
	return fileDescriptor_9e023dba0ca8fba6, []int{1}
}

func (m *Neighbors) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Neighbors.Unmarshal(m, b)
}
func (m *Neighbors) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Neighbors.Marshal(b, m, deterministic)
}
func (m *Neighbors) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Neighbors.Merge(m, src)
}
func (m *Neighbors) XXX_Size() int {
	return xxx_messageInfo_Neighbors.Size(m)
}
func (m *Neighbors) XXX_DiscardUnknown() {
	xxx_messageInfo_Neighbors.DiscardUnknown(m)
}

var xxx_messageInfo_Neighbors proto.InternalMessageInfo

func (m *Neighbors) GetNeighbors() []*NeighborInfo {
	if m != nil {
		return m.Neighbors
	}
	return nil
}

type NeighborInfo struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Ip                   string   `protobuf:"bytes,2,opt,name=ip,proto3" json:"ip,omitempty"`
	KadPort              string   `protobuf:"bytes,3,opt,name=kad_port,json=kadPort,proto3" json:"kad_port,omitempty"`
	ServPort             string   `protobuf:"bytes,4,opt,name=serv_port,json=servPort,proto3" json:"serv_port,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NeighborInfo) Reset()         { *m = NeighborInfo{} }
func (m *NeighborInfo) String() string { return proto.CompactTextString(m) }
func (*NeighborInfo) ProtoMessage()    {}
func (*NeighborInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_9e023dba0ca8fba6, []int{2}
}

func (m *NeighborInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NeighborInfo.Unmarshal(m, b)
}
func (m *NeighborInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NeighborInfo.Marshal(b, m, deterministic)
}
func (m *NeighborInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NeighborInfo.Merge(m, src)
}
func (m *NeighborInfo) XXX_Size() int {
	return xxx_messageInfo_NeighborInfo.Size(m)
}
func (m *NeighborInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_NeighborInfo.DiscardUnknown(m)
}

var xxx_messageInfo_NeighborInfo proto.InternalMessageInfo

func (m *NeighborInfo) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *NeighborInfo) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

func (m *NeighborInfo) GetKadPort() string {
	if m != nil {
		return m.KadPort
	}
	return ""
}

func (m *NeighborInfo) GetServPort() string {
	if m != nil {
		return m.ServPort
	}
	return ""
}

func init() {
	proto.RegisterType((*Target)(nil), "kademlia.Target")
	proto.RegisterType((*Neighbors)(nil), "kademlia.Neighbors")
	proto.RegisterType((*NeighborInfo)(nil), "kademlia.NeighborInfo")
}

func init() { proto.RegisterFile("kademlia_service.proto", fileDescriptor_9e023dba0ca8fba6) }

var fileDescriptor_9e023dba0ca8fba6 = []byte{
	// 277 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x91, 0xc1, 0x4e, 0xf2, 0x40,
	0x14, 0x85, 0xff, 0x16, 0x7e, 0x6c, 0xaf, 0x0a, 0x64, 0x4c, 0x48, 0xd5, 0x0d, 0xe9, 0xc2, 0x74,
	0xd5, 0x05, 0xf8, 0x02, 0x6e, 0x48, 0x08, 0x09, 0x31, 0xe0, 0xbe, 0x19, 0x9c, 0x0b, 0x4e, 0xaa,
	0x9d, 0xc9, 0x74, 0xc2, 0x7b, 0xf9, 0x86, 0xa6, 0x73, 0x67, 0x40, 0xe2, 0x6e, 0xe6, 0x9c, 0x6f,
	0xf1, 0x9d, 0x5c, 0x98, 0xd4, 0x5c, 0xe0, 0xd7, 0xa7, 0xe4, 0x55, 0x8b, 0xe6, 0x28, 0xdf, 0xb1,
	0xd4, 0x46, 0x59, 0xc5, 0x92, 0x90, 0xe7, 0xdf, 0x11, 0x0c, 0xde, 0xb8, 0x39, 0xa0, 0x65, 0x8f,
	0x90, 0x5a, 0xf7, 0xaa, 0xa4, 0xc8, 0xa2, 0x69, 0x54, 0xa4, 0x9b, 0x84, 0x82, 0xa5, 0xe8, 0xca,
	0x16, 0x1b, 0x81, 0xa6, 0x2b, 0x63, 0x2a, 0x29, 0xb8, 0x2c, 0x75, 0xd6, 0xbb, 0x28, 0x35, 0x7b,
	0x82, 0x91, 0x2f, 0x6b, 0x2e, 0x2a, 0xad, 0x8c, 0xcd, 0xfa, 0x0e, 0xb9, 0xa5, 0x78, 0xc5, 0xc5,
	0xab, 0x32, 0x96, 0x15, 0x30, 0xf6, 0x5c, 0xe7, 0x4a, 0xe0, 0x7f, 0x07, 0x0e, 0x29, 0xdf, 0xa2,
	0x39, 0x76, 0x64, 0xfe, 0x02, 0xe9, 0x1a, 0xe5, 0xe1, 0x63, 0xa7, 0x4c, 0xcb, 0x9e, 0x21, 0x6d,
	0xc2, 0x27, 0x8b, 0xa6, 0xbd, 0xe2, 0x7a, 0x36, 0x29, 0xc3, 0xbc, 0x32, 0x70, 0xcb, 0x66, 0xaf,
	0x36, 0x67, 0x30, 0xdf, 0xc3, 0xcd, 0xef, 0x8a, 0x0d, 0x21, 0x3e, 0x8d, 0x8e, 0xa5, 0x70, 0x7f,
	0xed, 0x77, 0xc6, 0x52, 0xb3, 0x7b, 0x48, 0x4e, 0xf6, 0x34, 0xf0, 0xaa, 0xf6, 0xde, 0x6e, 0x7c,
	0x10, 0xee, 0x87, 0xf1, 0xa4, 0x3a, 0x5b, 0xc0, 0x68, 0xe5, 0x5d, 0xb6, 0x74, 0x01, 0x36, 0x87,
	0x64, 0x21, 0x1b, 0xb1, 0x56, 0x02, 0xd9, 0xf8, 0x6c, 0x4a, 0x47, 0x78, 0xb8, 0xfb, 0xeb, 0xde,
	0xe6, 0xff, 0x76, 0x03, 0x77, 0xb7, 0xf9, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x35, 0x06, 0x4c,
	0xe7, 0xd1, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// KademliaServiceClient is the client API for KademliaService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type KademliaServiceClient interface {
	FindNode(ctx context.Context, in *Target, opts ...grpc.CallOption) (*Neighbors, error)
}

type kademliaServiceClient struct {
	cc *grpc.ClientConn
}

func NewKademliaServiceClient(cc *grpc.ClientConn) KademliaServiceClient {
	return &kademliaServiceClient{cc}
}

func (c *kademliaServiceClient) FindNode(ctx context.Context, in *Target, opts ...grpc.CallOption) (*Neighbors, error) {
	out := new(Neighbors)
	err := c.cc.Invoke(ctx, "/kademlia.KademliaService/FindNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KademliaServiceServer is the server API for KademliaService service.
type KademliaServiceServer interface {
	FindNode(context.Context, *Target) (*Neighbors, error)
}

// UnimplementedKademliaServiceServer can be embedded to have forward compatible implementations.
type UnimplementedKademliaServiceServer struct {
}

func (*UnimplementedKademliaServiceServer) FindNode(ctx context.Context, req *Target) (*Neighbors, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindNode not implemented")
}

func RegisterKademliaServiceServer(s *grpc.Server, srv KademliaServiceServer) {
	s.RegisterService(&_KademliaService_serviceDesc, srv)
}

func _KademliaService_FindNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Target)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KademliaServiceServer).FindNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kademlia.KademliaService/FindNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KademliaServiceServer).FindNode(ctx, req.(*Target))
	}
	return interceptor(ctx, in, info, handler)
}

var _KademliaService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "kademlia.KademliaService",
	HandlerType: (*KademliaServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindNode",
			Handler:    _KademliaService_FindNode_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kademlia_service.proto",
}