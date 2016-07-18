// Code generated by protoc-gen-go.
// source: anakin.proto
// DO NOT EDIT!

/*
Package remote is a generated protocol buffer package.

It is generated from these files:
	anakin.proto

It has these top-level messages:
	InstanceStats
	Instance
	RpcRequest
*/
package remote

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/timestamp"

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

type State int32

const (
	State_Active    State = 0
	State_Passive   State = 1
	State_Suspended State = 2
	State_Failing   State = 3
	State_Trying    State = 4
)

var State_name = map[int32]string{
	0: "Active",
	1: "Passive",
	2: "Suspended",
	3: "Failing",
	4: "Trying",
}
var State_value = map[string]int32{
	"Active":    0,
	"Passive":   1,
	"Suspended": 2,
	"Failing":   3,
	"Trying":    4,
}

func (x State) String() string {
	return proto.EnumName(State_name, int32(x))
}
func (State) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type InstanceStats struct {
	Os       string  `protobuf:"bytes,1,opt,name=os" json:"os,omitempty"`
	CpuCores int32   `protobuf:"varint,2,opt,name=cpuCores" json:"cpuCores,omitempty"`
	Mem      string  `protobuf:"bytes,3,opt,name=mem" json:"mem,omitempty"`
	Rps      float64 `protobuf:"fixed64,4,opt,name=rps" json:"rps,omitempty"`
}

func (m *InstanceStats) Reset()                    { *m = InstanceStats{} }
func (m *InstanceStats) String() string            { return proto.CompactTextString(m) }
func (*InstanceStats) ProtoMessage()               {}
func (*InstanceStats) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Instance struct {
	Id        string                     `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Version   string                     `protobuf:"bytes,2,opt,name=version" json:"version,omitempty"`
	AdminPort string                     `protobuf:"bytes,3,opt,name=adminPort" json:"adminPort,omitempty"`
	AdminIp   string                     `protobuf:"bytes,4,opt,name=adminIp" json:"adminIp,omitempty"`
	ProxyIp   string                     `protobuf:"bytes,5,opt,name=proxyIp" json:"proxyIp,omitempty"`
	ProxyPort string                     `protobuf:"bytes,6,opt,name=proxyPort" json:"proxyPort,omitempty"`
	Started   *google_protobuf.Timestamp `protobuf:"bytes,7,opt,name=started" json:"started,omitempty"`
	State     State                      `protobuf:"varint,8,opt,name=state,enum=remote.State" json:"state,omitempty"`
	Stats     *InstanceStats             `protobuf:"bytes,9,opt,name=stats" json:"stats,omitempty"`
}

func (m *Instance) Reset()                    { *m = Instance{} }
func (m *Instance) String() string            { return proto.CompactTextString(m) }
func (*Instance) ProtoMessage()               {}
func (*Instance) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Instance) GetStarted() *google_protobuf.Timestamp {
	if m != nil {
		return m.Started
	}
	return nil
}

func (m *Instance) GetStats() *InstanceStats {
	if m != nil {
		return m.Stats
	}
	return nil
}

type RpcRequest struct {
	Time *google_protobuf.Timestamp `protobuf:"bytes,1,opt,name=time" json:"time,omitempty"`
}

func (m *RpcRequest) Reset()                    { *m = RpcRequest{} }
func (m *RpcRequest) String() string            { return proto.CompactTextString(m) }
func (*RpcRequest) ProtoMessage()               {}
func (*RpcRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *RpcRequest) GetTime() *google_protobuf.Timestamp {
	if m != nil {
		return m.Time
	}
	return nil
}

func init() {
	proto.RegisterType((*InstanceStats)(nil), "remote.InstanceStats")
	proto.RegisterType((*Instance)(nil), "remote.Instance")
	proto.RegisterType((*RpcRequest)(nil), "remote.RpcRequest")
	proto.RegisterEnum("remote.State", State_name, State_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion2

// Client API for Anakin service

type AnakinClient interface {
	GetInstance(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (*Instance, error)
}

type anakinClient struct {
	cc *grpc.ClientConn
}

func NewAnakinClient(cc *grpc.ClientConn) AnakinClient {
	return &anakinClient{cc}
}

func (c *anakinClient) GetInstance(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (*Instance, error) {
	out := new(Instance)
	err := grpc.Invoke(ctx, "/remote.Anakin/GetInstance", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Anakin service

type AnakinServer interface {
	GetInstance(context.Context, *RpcRequest) (*Instance, error)
}

func RegisterAnakinServer(s *grpc.Server, srv AnakinServer) {
	s.RegisterService(&_Anakin_serviceDesc, srv)
}

func _Anakin_GetInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RpcRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnakinServer).GetInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/remote.Anakin/GetInstance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnakinServer).GetInstance(ctx, req.(*RpcRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Anakin_serviceDesc = grpc.ServiceDesc{
	ServiceName: "remote.Anakin",
	HandlerType: (*AnakinServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetInstance",
			Handler:    _Anakin_GetInstance_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

var fileDescriptor0 = []byte{
	// 357 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x91, 0xc1, 0x4a, 0xf3, 0x40,
	0x10, 0xc7, 0xbf, 0xb4, 0x4d, 0xd2, 0x4c, 0xbf, 0x6a, 0x5c, 0x10, 0x42, 0x11, 0x94, 0xe2, 0xa1,
	0x28, 0xa4, 0xd0, 0x82, 0x37, 0x0f, 0x22, 0x58, 0x7b, 0x2b, 0x6d, 0x5f, 0x20, 0x4d, 0xc6, 0xb2,
	0xd8, 0xec, 0xc6, 0xdd, 0x4d, 0xb1, 0x6f, 0xeb, 0xa3, 0x38, 0xd9, 0x34, 0x8a, 0x1e, 0xbc, 0xed,
	0xfc, 0x67, 0xe6, 0x37, 0xf3, 0x9f, 0x85, 0xff, 0x89, 0x48, 0x5e, 0xb9, 0x88, 0x0b, 0x25, 0x8d,
	0x64, 0x9e, 0xc2, 0x5c, 0x1a, 0x1c, 0x5c, 0x6e, 0xa5, 0xdc, 0xee, 0x70, 0x6c, 0xd5, 0x4d, 0xf9,
	0x32, 0x36, 0x3c, 0x47, 0x6d, 0x92, 0xbc, 0xa8, 0x0b, 0x87, 0x33, 0xe8, 0xcf, 0x05, 0x09, 0x22,
	0xc5, 0x95, 0x49, 0x8c, 0x66, 0x00, 0x2d, 0xa9, 0x23, 0xe7, 0xca, 0x19, 0x05, 0x2c, 0x84, 0x6e,
	0x5a, 0x94, 0x8f, 0x52, 0xa1, 0x8e, 0x5a, 0xa4, 0xb8, 0xac, 0x07, 0xed, 0x1c, 0xf3, 0xa8, 0x6d,
	0xd3, 0x14, 0xa8, 0x42, 0x47, 0x1d, 0x0a, 0x9c, 0xe1, 0x87, 0x03, 0xdd, 0x86, 0x54, 0x41, 0x78,
	0x76, 0x84, 0x9c, 0x82, 0xbf, 0x47, 0xa5, 0xb9, 0x14, 0x96, 0x11, 0xb0, 0x33, 0x08, 0x92, 0x2c,
	0xe7, 0x62, 0x21, 0x95, 0x39, 0x92, 0xa8, 0xc6, 0x4a, 0xf3, 0xc2, 0xd2, 0xac, 0x40, 0xfb, 0xbd,
	0x1f, 0x48, 0x70, 0x9b, 0x26, 0x2b, 0xd8, 0x26, 0xcf, 0x4a, 0xb7, 0xe0, 0xd3, 0x38, 0x65, 0x30,
	0x8b, 0x7c, 0x12, 0x7a, 0x93, 0x41, 0x5c, 0xbb, 0x8d, 0x1b, 0xb7, 0xf1, 0xba, 0x71, 0xcb, 0x2e,
	0xc0, 0xa5, 0x87, 0xc1, 0xa8, 0x4b, 0xa5, 0x27, 0x93, 0x7e, 0x5c, 0x1f, 0x28, 0xae, 0x4c, 0x23,
	0xbb, 0xae, 0xb3, 0x3a, 0x0a, 0x2c, 0xe8, 0xbc, 0xc9, 0xfe, 0x38, 0xcd, 0xf0, 0x0e, 0x60, 0x59,
	0xa4, 0x4b, 0x7c, 0x2b, 0x09, 0xca, 0x46, 0xd0, 0xa9, 0x8e, 0x69, 0x5d, 0xfe, 0x39, 0xfb, 0xe6,
	0x19, 0xdc, 0x7a, 0x0c, 0x80, 0xf7, 0x90, 0x1a, 0xbe, 0xc7, 0xf0, 0x1f, 0x1d, 0xcf, 0x5f, 0x24,
	0x5a, 0x57, 0x81, 0xc3, 0xfa, 0x10, 0xac, 0x4a, 0x5d, 0xa0, 0xc8, 0x30, 0x0b, 0x5b, 0x55, 0xee,
	0x29, 0xe1, 0x3b, 0x2e, 0xb6, 0x61, 0xbb, 0x6a, 0x5a, 0xab, 0x43, 0xf5, 0xee, 0x4c, 0xee, 0x09,
	0x60, 0xbf, 0x99, 0x4d, 0xa1, 0x37, 0x43, 0xf3, 0x75, 0x70, 0xd6, 0x6c, 0xfc, 0xbd, 0xe0, 0x20,
	0xfc, 0xed, 0x62, 0xe3, 0xd9, 0xe5, 0xa6, 0x9f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x3e, 0x63, 0x9a,
	0x3f, 0x2c, 0x02, 0x00, 0x00,
}