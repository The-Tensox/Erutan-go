// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protobuf/protometry/volume.proto

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

type Sphere struct {
	Center               *Vector3 `protobuf:"bytes,1,opt,name=center,proto3" json:"center,omitempty"`
	Radius               float64  `protobuf:"fixed64,2,opt,name=radius,proto3" json:"radius,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Sphere) Reset()         { *m = Sphere{} }
func (m *Sphere) String() string { return proto.CompactTextString(m) }
func (*Sphere) ProtoMessage()    {}
func (*Sphere) Descriptor() ([]byte, []int) {
	return fileDescriptor_d31fb1d8f7886926, []int{0}
}

func (m *Sphere) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Sphere.Unmarshal(m, b)
}
func (m *Sphere) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Sphere.Marshal(b, m, deterministic)
}
func (m *Sphere) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Sphere.Merge(m, src)
}
func (m *Sphere) XXX_Size() int {
	return xxx_messageInfo_Sphere.Size(m)
}
func (m *Sphere) XXX_DiscardUnknown() {
	xxx_messageInfo_Sphere.DiscardUnknown(m)
}

var xxx_messageInfo_Sphere proto.InternalMessageInfo

func (m *Sphere) GetCenter() *Vector3 {
	if m != nil {
		return m.Center
	}
	return nil
}

func (m *Sphere) GetRadius() float64 {
	if m != nil {
		return m.Radius
	}
	return 0
}

type Capsule struct {
	Center               *Vector3 `protobuf:"bytes,1,opt,name=center,proto3" json:"center,omitempty"`
	Width                float64  `protobuf:"fixed64,2,opt,name=width,proto3" json:"width,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Capsule) Reset()         { *m = Capsule{} }
func (m *Capsule) String() string { return proto.CompactTextString(m) }
func (*Capsule) ProtoMessage()    {}
func (*Capsule) Descriptor() ([]byte, []int) {
	return fileDescriptor_d31fb1d8f7886926, []int{1}
}

func (m *Capsule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Capsule.Unmarshal(m, b)
}
func (m *Capsule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Capsule.Marshal(b, m, deterministic)
}
func (m *Capsule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Capsule.Merge(m, src)
}
func (m *Capsule) XXX_Size() int {
	return xxx_messageInfo_Capsule.Size(m)
}
func (m *Capsule) XXX_DiscardUnknown() {
	xxx_messageInfo_Capsule.DiscardUnknown(m)
}

var xxx_messageInfo_Capsule proto.InternalMessageInfo

func (m *Capsule) GetCenter() *Vector3 {
	if m != nil {
		return m.Center
	}
	return nil
}

func (m *Capsule) GetWidth() float64 {
	if m != nil {
		return m.Width
	}
	return 0
}

// Box is an AABB volume
type Box struct {
	Min                  *Vector3 `protobuf:"bytes,1,opt,name=min,proto3" json:"min,omitempty"`
	Max                  *Vector3 `protobuf:"bytes,2,opt,name=max,proto3" json:"max,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Box) Reset()         { *m = Box{} }
func (m *Box) String() string { return proto.CompactTextString(m) }
func (*Box) ProtoMessage()    {}
func (*Box) Descriptor() ([]byte, []int) {
	return fileDescriptor_d31fb1d8f7886926, []int{2}
}

func (m *Box) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Box.Unmarshal(m, b)
}
func (m *Box) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Box.Marshal(b, m, deterministic)
}
func (m *Box) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Box.Merge(m, src)
}
func (m *Box) XXX_Size() int {
	return xxx_messageInfo_Box.Size(m)
}
func (m *Box) XXX_DiscardUnknown() {
	xxx_messageInfo_Box.DiscardUnknown(m)
}

var xxx_messageInfo_Box proto.InternalMessageInfo

func (m *Box) GetMin() *Vector3 {
	if m != nil {
		return m.Min
	}
	return nil
}

func (m *Box) GetMax() *Vector3 {
	if m != nil {
		return m.Max
	}
	return nil
}

type Mesh struct {
	Center               *Vector3   `protobuf:"bytes,1,opt,name=center,proto3" json:"center,omitempty"`
	Vertices             []*Vector3 `protobuf:"bytes,2,rep,name=vertices,proto3" json:"vertices,omitempty"`
	Tris                 []int32    `protobuf:"varint,3,rep,packed,name=tris,proto3" json:"tris,omitempty"`
	Normals              []*Vector3 `protobuf:"bytes,4,rep,name=normals,proto3" json:"normals,omitempty"`
	Uvs                  []*Vector3 `protobuf:"bytes,5,rep,name=uvs,proto3" json:"uvs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Mesh) Reset()         { *m = Mesh{} }
func (m *Mesh) String() string { return proto.CompactTextString(m) }
func (*Mesh) ProtoMessage()    {}
func (*Mesh) Descriptor() ([]byte, []int) {
	return fileDescriptor_d31fb1d8f7886926, []int{3}
}

func (m *Mesh) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Mesh.Unmarshal(m, b)
}
func (m *Mesh) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Mesh.Marshal(b, m, deterministic)
}
func (m *Mesh) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Mesh.Merge(m, src)
}
func (m *Mesh) XXX_Size() int {
	return xxx_messageInfo_Mesh.Size(m)
}
func (m *Mesh) XXX_DiscardUnknown() {
	xxx_messageInfo_Mesh.DiscardUnknown(m)
}

var xxx_messageInfo_Mesh proto.InternalMessageInfo

func (m *Mesh) GetCenter() *Vector3 {
	if m != nil {
		return m.Center
	}
	return nil
}

func (m *Mesh) GetVertices() []*Vector3 {
	if m != nil {
		return m.Vertices
	}
	return nil
}

func (m *Mesh) GetTris() []int32 {
	if m != nil {
		return m.Tris
	}
	return nil
}

func (m *Mesh) GetNormals() []*Vector3 {
	if m != nil {
		return m.Normals
	}
	return nil
}

func (m *Mesh) GetUvs() []*Vector3 {
	if m != nil {
		return m.Uvs
	}
	return nil
}

func init() {
	proto.RegisterType((*Sphere)(nil), "protometry.Sphere")
	proto.RegisterType((*Capsule)(nil), "protometry.Capsule")
	proto.RegisterType((*Box)(nil), "protometry.Box")
	proto.RegisterType((*Mesh)(nil), "protometry.Mesh")
}

func init() {
	proto.RegisterFile("protobuf/protometry/volume.proto", fileDescriptor_d31fb1d8f7886926)
}

var fileDescriptor_d31fb1d8f7886926 = []byte{
	// 281 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0xdf, 0x4a, 0x84, 0x40,
	0x14, 0xc6, 0x71, 0x5d, 0xdd, 0x38, 0xdd, 0x4d, 0x11, 0xd2, 0x95, 0x09, 0x81, 0x10, 0xab, 0xd0,
	0xbe, 0xc1, 0x76, 0xdb, 0x42, 0xb8, 0x4b, 0x17, 0xdd, 0xa9, 0x7b, 0x6a, 0x06, 0x1c, 0x47, 0xe6,
	0x8f, 0xd9, 0x2b, 0xf5, 0x30, 0x3d, 0x53, 0x38, 0xda, 0x6e, 0x17, 0x19, 0xec, 0xdd, 0xf9, 0xfc,
	0x7e, 0xe7, 0x87, 0x87, 0x81, 0xb0, 0x91, 0x42, 0x8b, 0xc2, 0xbc, 0xa6, 0x76, 0xe0, 0xa8, 0xe5,
	0x47, 0xda, 0x8a, 0xca, 0x70, 0x4c, 0xec, 0x17, 0x02, 0xc7, 0xe2, 0xfa, 0xe6, 0x4f, 0x1a, 0x4b,
	0x2d, 0xe4, 0x6a, 0xc0, 0xa3, 0x0d, 0xf8, 0xdb, 0x86, 0xa2, 0x44, 0x72, 0x07, 0x7e, 0x89, 0xb5,
	0x46, 0x19, 0x38, 0xa1, 0x13, 0x9f, 0xdf, 0x5f, 0x24, 0xc7, 0xa5, 0xe4, 0x79, 0x58, 0xca, 0x46,
	0x84, 0x5c, 0x81, 0x2f, 0xf3, 0x3d, 0x33, 0x2a, 0x98, 0x85, 0x4e, 0xec, 0x64, 0x63, 0x8a, 0x1e,
	0x61, 0xf1, 0x90, 0x37, 0xca, 0x54, 0x27, 0xfa, 0x2e, 0xc1, 0x7b, 0x67, 0x7b, 0x4d, 0x47, 0xdd,
	0x10, 0xa2, 0x2d, 0xb8, 0x6b, 0xd1, 0x91, 0x5b, 0x70, 0x39, 0xab, 0xff, 0xd3, 0xf4, 0xbd, 0xc5,
	0xf2, 0xce, 0x1a, 0x26, 0xb1, 0xbc, 0x8b, 0xbe, 0x1c, 0x98, 0x6f, 0x50, 0xd1, 0xd3, 0x7e, 0x30,
	0x85, 0xb3, 0x16, 0xa5, 0x66, 0x25, 0xf6, 0x27, 0xbb, 0x53, 0xf8, 0x01, 0x22, 0x04, 0xe6, 0x5a,
	0x32, 0x15, 0xb8, 0xa1, 0x1b, 0x7b, 0x99, 0x9d, 0xc9, 0x12, 0x16, 0xb5, 0x90, 0x3c, 0xaf, 0x54,
	0x30, 0x9f, 0x76, 0xfc, 0x30, 0xfd, 0x41, 0xa6, 0x55, 0x81, 0x37, 0x8d, 0xf6, 0xfd, 0x3a, 0x7d,
	0x09, 0xdf, 0x98, 0xa6, 0xa6, 0x48, 0x4a, 0xc1, 0xd3, 0x1d, 0xc5, 0xe5, 0x0e, 0x6b, 0x25, 0xba,
	0x5f, 0x8f, 0xfe, 0x39, 0x83, 0xa7, 0x43, 0x28, 0x7c, 0x5b, 0xac, 0xbe, 0x03, 0x00, 0x00, 0xff,
	0xff, 0x0c, 0x31, 0x05, 0xdb, 0x4d, 0x02, 0x00, 0x00,
}
