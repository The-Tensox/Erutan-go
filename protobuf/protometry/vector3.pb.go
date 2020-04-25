// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protobuf/protometry/vector3.proto

package protometry

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type Vector3 struct {
	X                    float64  `protobuf:"fixed64,1,opt,name=x,proto3" json:"x,omitempty"`
	Y                    float64  `protobuf:"fixed64,2,opt,name=y,proto3" json:"y,omitempty"`
	Z                    float64  `protobuf:"fixed64,3,opt,name=z,proto3" json:"z,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Vector3) Reset()         { *m = Vector3{} }
func (m *Vector3) String() string { return proto.CompactTextString(m) }
func (*Vector3) ProtoMessage()    {}
func (*Vector3) Descriptor() ([]byte, []int) {
	return fileDescriptor_b35e8f0b9325e45c, []int{0}
}

func (m *Vector3) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Vector3.Unmarshal(m, b)
}
func (m *Vector3) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Vector3.Marshal(b, m, deterministic)
}
func (m *Vector3) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Vector3.Merge(m, src)
}
func (m *Vector3) XXX_Size() int {
	return xxx_messageInfo_Vector3.Size(m)
}
func (m *Vector3) XXX_DiscardUnknown() {
	xxx_messageInfo_Vector3.DiscardUnknown(m)
}

var xxx_messageInfo_Vector3 proto.InternalMessageInfo

func (m *Vector3) GetX() float64 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *Vector3) GetY() float64 {
	if m != nil {
		return m.Y
	}
	return 0
}

func (m *Vector3) GetZ() float64 {
	if m != nil {
		return m.Z
	}
	return 0
}

func init() {
	proto.RegisterType((*Vector3)(nil), "protometry.Vector3")
}

func init() {
	proto.RegisterFile("protobuf/protometry/vector3.proto", fileDescriptor_b35e8f0b9325e45c)
}

var fileDescriptor_b35e8f0b9325e45c = []byte{
	// 130 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x2c, 0x28, 0xca, 0x2f,
	0xc9, 0x4f, 0x2a, 0x4d, 0xd3, 0x07, 0x33, 0x72, 0x53, 0x4b, 0x8a, 0x2a, 0xf5, 0xcb, 0x52, 0x93,
	0x4b, 0xf2, 0x8b, 0x8c, 0xf5, 0xc0, 0x42, 0x42, 0x5c, 0x08, 0x19, 0x25, 0x63, 0x2e, 0xf6, 0x30,
	0x88, 0xa4, 0x10, 0x0f, 0x17, 0x63, 0x85, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x63, 0x10, 0x63, 0x05,
	0x88, 0x57, 0x29, 0xc1, 0x04, 0xe1, 0x55, 0x82, 0x78, 0x55, 0x12, 0xcc, 0x10, 0x5e, 0x95, 0x93,
	0x7e, 0x94, 0x42, 0x7a, 0x66, 0x49, 0x46, 0x69, 0x92, 0x5e, 0x72, 0x7e, 0xae, 0x7e, 0x48, 0x46,
	0xaa, 0x6e, 0x48, 0x6a, 0x5e, 0x71, 0x7e, 0x05, 0x92, 0x95, 0xab, 0x98, 0xb8, 0x02, 0xe0, 0x9c,
	0x24, 0x36, 0xb0, 0x84, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x43, 0x56, 0x17, 0x91, 0x9d, 0x00,
	0x00, 0x00,
}