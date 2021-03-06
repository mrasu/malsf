// Code generated by protoc-gen-go.
// source: action.proto
// DO NOT EDIT!

/*
Package structs is a generated protocol buffer package.

It is generated from these files:
	action.proto

It has these top-level messages:
	Action
	Reaction
*/
package structs

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Action struct {
	Id       int32  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Message  string `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
	Type     string `protobuf:"bytes,3,opt,name=type" json:"type,omitempty"`
	NodeName string `protobuf:"bytes,4,opt,name=nodeName" json:"nodeName,omitempty"`
	Service  string `protobuf:"bytes,5,opt,name=service" json:"service,omitempty"`
}

func (m *Action) Reset()                    { *m = Action{} }
func (m *Action) String() string            { return proto.CompactTextString(m) }
func (*Action) ProtoMessage()               {}
func (*Action) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Reaction struct {
	Id       int32  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	FromId   int32  `protobuf:"varint,2,opt,name=fromId" json:"fromId,omitempty"`
	Code     int32  `protobuf:"varint,3,opt,name=code" json:"code,omitempty"`
	Message  string `protobuf:"bytes,4,opt,name=message" json:"message,omitempty"`
	NodeName string `protobuf:"bytes,5,opt,name=nodeName" json:"nodeName,omitempty"`
	Service  string `protobuf:"bytes,6,opt,name=service" json:"service,omitempty"`
}

func (m *Reaction) Reset()                    { *m = Reaction{} }
func (m *Reaction) String() string            { return proto.CompactTextString(m) }
func (*Reaction) ProtoMessage()               {}
func (*Reaction) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*Action)(nil), "structs.Action")
	proto.RegisterType((*Reaction)(nil), "structs.Reaction")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for ActionService service

type ActionServiceClient interface {
	Notify(ctx context.Context, in *Action, opts ...grpc.CallOption) (*Reaction, error)
}

type actionServiceClient struct {
	cc *grpc.ClientConn
}

func NewActionServiceClient(cc *grpc.ClientConn) ActionServiceClient {
	return &actionServiceClient{cc}
}

func (c *actionServiceClient) Notify(ctx context.Context, in *Action, opts ...grpc.CallOption) (*Reaction, error) {
	out := new(Reaction)
	err := grpc.Invoke(ctx, "/structs.ActionService/Notify", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ActionService service

type ActionServiceServer interface {
	Notify(context.Context, *Action) (*Reaction, error)
}

func RegisterActionServiceServer(s *grpc.Server, srv ActionServiceServer) {
	s.RegisterService(&_ActionService_serviceDesc, srv)
}

func _ActionService_Notify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Action)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ActionServiceServer).Notify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/structs.ActionService/Notify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ActionServiceServer).Notify(ctx, req.(*Action))
	}
	return interceptor(ctx, in, info, handler)
}

var _ActionService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "structs.ActionService",
	HandlerType: (*ActionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Notify",
			Handler:    _ActionService_Notify_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("action.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 229 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0x4c, 0x2e, 0xc9,
	0xcc, 0xcf, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2f, 0x2e, 0x29, 0x2a, 0x4d, 0x2e,
	0x29, 0x56, 0xaa, 0xe1, 0x62, 0x73, 0x04, 0x4b, 0x08, 0xf1, 0x71, 0x31, 0x65, 0xa6, 0x48, 0x30,
	0x2a, 0x30, 0x6a, 0xb0, 0x06, 0x31, 0x65, 0xa6, 0x08, 0x49, 0x70, 0xb1, 0xe7, 0xa6, 0x16, 0x17,
	0x27, 0xa6, 0xa7, 0x4a, 0x30, 0x29, 0x30, 0x6a, 0x70, 0x06, 0xc1, 0xb8, 0x42, 0x42, 0x5c, 0x2c,
	0x25, 0x95, 0x05, 0xa9, 0x12, 0xcc, 0x60, 0x61, 0x30, 0x5b, 0x48, 0x8a, 0x8b, 0x23, 0x2f, 0x3f,
	0x25, 0xd5, 0x2f, 0x31, 0x37, 0x55, 0x82, 0x05, 0x2c, 0x0e, 0xe7, 0x83, 0x4c, 0x2a, 0x4e, 0x2d,
	0x2a, 0xcb, 0x4c, 0x4e, 0x95, 0x60, 0x85, 0x98, 0x04, 0xe5, 0x2a, 0x4d, 0x63, 0xe4, 0xe2, 0x08,
	0x4a, 0x4d, 0xc4, 0xee, 0x00, 0x31, 0x2e, 0xb6, 0xb4, 0xa2, 0xfc, 0x5c, 0xcf, 0x14, 0xb0, 0xfd,
	0xac, 0x41, 0x50, 0x1e, 0xc8, 0xfa, 0xe4, 0xfc, 0x14, 0x88, 0xf5, 0xac, 0x41, 0x60, 0x36, 0xb2,
	0x63, 0x59, 0x50, 0x1d, 0x8b, 0xec, 0x30, 0x56, 0xdc, 0x0e, 0x63, 0x43, 0x71, 0x98, 0x91, 0x3d,
	0x17, 0x2f, 0x24, 0x58, 0x82, 0x21, 0x02, 0x42, 0x7a, 0x5c, 0x6c, 0x7e, 0xf9, 0x25, 0x99, 0x69,
	0x95, 0x42, 0xfc, 0x7a, 0xd0, 0xb0, 0xd3, 0x83, 0xa8, 0x90, 0x12, 0x84, 0x0b, 0xc0, 0xbc, 0xa2,
	0xc4, 0x90, 0xc4, 0x06, 0x0e, 0x67, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x4b, 0xdd, 0xf6,
	0x7e, 0x77, 0x01, 0x00, 0x00,
}
