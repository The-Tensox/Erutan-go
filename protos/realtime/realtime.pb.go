// Code generated by protoc-gen-go. DO NOT EDIT.
// source: realtime.proto

package erutan

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type NetObject_Type int32

const (
	NetObject_ANIMAL NetObject_Type = 0
	NetObject_FOOD   NetObject_Type = 1
	NetObject_GROUND NetObject_Type = 2
)

var NetObject_Type_name = map[int32]string{
	0: "ANIMAL",
	1: "FOOD",
	2: "GROUND",
}

var NetObject_Type_value = map[string]int32{
	"ANIMAL": 0,
	"FOOD":   1,
	"GROUND": 2,
}

func (x NetObject_Type) String() string {
	return proto.EnumName(NetObject_Type_name, int32(x))
}

func (NetObject_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_dcbdca058206953b, []int{3, 0}
}

type Metadata struct {
	Timestamp            *timestamp.Timestamp `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Metadata) Reset()         { *m = Metadata{} }
func (m *Metadata) String() string { return proto.CompactTextString(m) }
func (*Metadata) ProtoMessage()    {}
func (*Metadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_dcbdca058206953b, []int{0}
}

func (m *Metadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Metadata.Unmarshal(m, b)
}
func (m *Metadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Metadata.Marshal(b, m, deterministic)
}
func (m *Metadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Metadata.Merge(m, src)
}
func (m *Metadata) XXX_Size() int {
	return xxx_messageInfo_Metadata.Size(m)
}
func (m *Metadata) XXX_DiscardUnknown() {
	xxx_messageInfo_Metadata.DiscardUnknown(m)
}

var xxx_messageInfo_Metadata proto.InternalMessageInfo

func (m *Metadata) GetTimestamp() *timestamp.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

type NetVector3 struct {
	X                    float32  `protobuf:"fixed32,1,opt,name=x,proto3" json:"x,omitempty"`
	Y                    float32  `protobuf:"fixed32,2,opt,name=y,proto3" json:"y,omitempty"`
	Z                    float32  `protobuf:"fixed32,3,opt,name=z,proto3" json:"z,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NetVector3) Reset()         { *m = NetVector3{} }
func (m *NetVector3) String() string { return proto.CompactTextString(m) }
func (*NetVector3) ProtoMessage()    {}
func (*NetVector3) Descriptor() ([]byte, []int) {
	return fileDescriptor_dcbdca058206953b, []int{1}
}

func (m *NetVector3) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetVector3.Unmarshal(m, b)
}
func (m *NetVector3) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetVector3.Marshal(b, m, deterministic)
}
func (m *NetVector3) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetVector3.Merge(m, src)
}
func (m *NetVector3) XXX_Size() int {
	return xxx_messageInfo_NetVector3.Size(m)
}
func (m *NetVector3) XXX_DiscardUnknown() {
	xxx_messageInfo_NetVector3.DiscardUnknown(m)
}

var xxx_messageInfo_NetVector3 proto.InternalMessageInfo

func (m *NetVector3) GetX() float32 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *NetVector3) GetY() float32 {
	if m != nil {
		return m.Y
	}
	return 0
}

func (m *NetVector3) GetZ() float32 {
	if m != nil {
		return m.Z
	}
	return 0
}

type NetQuaternion struct {
	X                    float32  `protobuf:"fixed32,1,opt,name=x,proto3" json:"x,omitempty"`
	Y                    float32  `protobuf:"fixed32,2,opt,name=y,proto3" json:"y,omitempty"`
	Z                    float32  `protobuf:"fixed32,3,opt,name=z,proto3" json:"z,omitempty"`
	W                    float32  `protobuf:"fixed32,4,opt,name=w,proto3" json:"w,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NetQuaternion) Reset()         { *m = NetQuaternion{} }
func (m *NetQuaternion) String() string { return proto.CompactTextString(m) }
func (*NetQuaternion) ProtoMessage()    {}
func (*NetQuaternion) Descriptor() ([]byte, []int) {
	return fileDescriptor_dcbdca058206953b, []int{2}
}

func (m *NetQuaternion) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetQuaternion.Unmarshal(m, b)
}
func (m *NetQuaternion) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetQuaternion.Marshal(b, m, deterministic)
}
func (m *NetQuaternion) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetQuaternion.Merge(m, src)
}
func (m *NetQuaternion) XXX_Size() int {
	return xxx_messageInfo_NetQuaternion.Size(m)
}
func (m *NetQuaternion) XXX_DiscardUnknown() {
	xxx_messageInfo_NetQuaternion.DiscardUnknown(m)
}

var xxx_messageInfo_NetQuaternion proto.InternalMessageInfo

func (m *NetQuaternion) GetX() float32 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *NetQuaternion) GetY() float32 {
	if m != nil {
		return m.Y
	}
	return 0
}

func (m *NetQuaternion) GetZ() float32 {
	if m != nil {
		return m.Z
	}
	return 0
}

func (m *NetQuaternion) GetW() float32 {
	if m != nil {
		return m.W
	}
	return 0
}

// TODO: parent ...
type NetObject struct {
	ObjectId             string         `protobuf:"bytes,1,opt,name=object_id,json=objectId,proto3" json:"object_id,omitempty"`
	OwnerId              string         `protobuf:"bytes,2,opt,name=owner_id,json=ownerId,proto3" json:"owner_id,omitempty"`
	Position             *NetVector3    `protobuf:"bytes,3,opt,name=position,proto3" json:"position,omitempty"`
	Rotation             *NetQuaternion `protobuf:"bytes,4,opt,name=rotation,proto3" json:"rotation,omitempty"`
	Scale                *NetVector3    `protobuf:"bytes,5,opt,name=scale,proto3" json:"scale,omitempty"`
	Type                 NetObject_Type `protobuf:"varint,6,opt,name=type,proto3,enum=erutan.NetObject_Type" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *NetObject) Reset()         { *m = NetObject{} }
func (m *NetObject) String() string { return proto.CompactTextString(m) }
func (*NetObject) ProtoMessage()    {}
func (*NetObject) Descriptor() ([]byte, []int) {
	return fileDescriptor_dcbdca058206953b, []int{3}
}

func (m *NetObject) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetObject.Unmarshal(m, b)
}
func (m *NetObject) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetObject.Marshal(b, m, deterministic)
}
func (m *NetObject) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetObject.Merge(m, src)
}
func (m *NetObject) XXX_Size() int {
	return xxx_messageInfo_NetObject.Size(m)
}
func (m *NetObject) XXX_DiscardUnknown() {
	xxx_messageInfo_NetObject.DiscardUnknown(m)
}

var xxx_messageInfo_NetObject proto.InternalMessageInfo

func (m *NetObject) GetObjectId() string {
	if m != nil {
		return m.ObjectId
	}
	return ""
}

func (m *NetObject) GetOwnerId() string {
	if m != nil {
		return m.OwnerId
	}
	return ""
}

func (m *NetObject) GetPosition() *NetVector3 {
	if m != nil {
		return m.Position
	}
	return nil
}

func (m *NetObject) GetRotation() *NetQuaternion {
	if m != nil {
		return m.Rotation
	}
	return nil
}

func (m *NetObject) GetScale() *NetVector3 {
	if m != nil {
		return m.Scale
	}
	return nil
}

func (m *NetObject) GetType() NetObject_Type {
	if m != nil {
		return m.Type
	}
	return NetObject_ANIMAL
}

type Packet struct {
	Metadata *Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	// Types that are valid to be assigned to Type:
	//	*Packet_CreateObject
	//	*Packet_UpdatePosition
	//	*Packet_DestroyObject
	Type                 isPacket_Type `protobuf_oneof:"type"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Packet) Reset()         { *m = Packet{} }
func (m *Packet) String() string { return proto.CompactTextString(m) }
func (*Packet) ProtoMessage()    {}
func (*Packet) Descriptor() ([]byte, []int) {
	return fileDescriptor_dcbdca058206953b, []int{4}
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

func (m *Packet) GetMetadata() *Metadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type isPacket_Type interface {
	isPacket_Type()
}

type Packet_CreateObject struct {
	CreateObject *Packet_CreateObjectPacket `protobuf:"bytes,2,opt,name=create_object,json=createObject,proto3,oneof"`
}

type Packet_UpdatePosition struct {
	UpdatePosition *Packet_UpdatePositionPacket `protobuf:"bytes,3,opt,name=update_position,json=updatePosition,proto3,oneof"`
}

type Packet_DestroyObject struct {
	DestroyObject *Packet_DestroyObjectPacket `protobuf:"bytes,4,opt,name=destroy_object,json=destroyObject,proto3,oneof"`
}

func (*Packet_CreateObject) isPacket_Type() {}

func (*Packet_UpdatePosition) isPacket_Type() {}

func (*Packet_DestroyObject) isPacket_Type() {}

func (m *Packet) GetType() isPacket_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (m *Packet) GetCreateObject() *Packet_CreateObjectPacket {
	if x, ok := m.GetType().(*Packet_CreateObject); ok {
		return x.CreateObject
	}
	return nil
}

func (m *Packet) GetUpdatePosition() *Packet_UpdatePositionPacket {
	if x, ok := m.GetType().(*Packet_UpdatePosition); ok {
		return x.UpdatePosition
	}
	return nil
}

func (m *Packet) GetDestroyObject() *Packet_DestroyObjectPacket {
	if x, ok := m.GetType().(*Packet_DestroyObject); ok {
		return x.DestroyObject
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Packet) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Packet_CreateObject)(nil),
		(*Packet_UpdatePosition)(nil),
		(*Packet_DestroyObject)(nil),
	}
}

type Packet_CreateObjectPacket struct {
	Object               *NetObject `protobuf:"bytes,1,opt,name=object,proto3" json:"object,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Packet_CreateObjectPacket) Reset()         { *m = Packet_CreateObjectPacket{} }
func (m *Packet_CreateObjectPacket) String() string { return proto.CompactTextString(m) }
func (*Packet_CreateObjectPacket) ProtoMessage()    {}
func (*Packet_CreateObjectPacket) Descriptor() ([]byte, []int) {
	return fileDescriptor_dcbdca058206953b, []int{4, 0}
}

func (m *Packet_CreateObjectPacket) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Packet_CreateObjectPacket.Unmarshal(m, b)
}
func (m *Packet_CreateObjectPacket) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Packet_CreateObjectPacket.Marshal(b, m, deterministic)
}
func (m *Packet_CreateObjectPacket) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Packet_CreateObjectPacket.Merge(m, src)
}
func (m *Packet_CreateObjectPacket) XXX_Size() int {
	return xxx_messageInfo_Packet_CreateObjectPacket.Size(m)
}
func (m *Packet_CreateObjectPacket) XXX_DiscardUnknown() {
	xxx_messageInfo_Packet_CreateObjectPacket.DiscardUnknown(m)
}

var xxx_messageInfo_Packet_CreateObjectPacket proto.InternalMessageInfo

func (m *Packet_CreateObjectPacket) GetObject() *NetObject {
	if m != nil {
		return m.Object
	}
	return nil
}

type Packet_UpdatePositionPacket struct {
	ObjectId             string      `protobuf:"bytes,1,opt,name=object_id,json=objectId,proto3" json:"object_id,omitempty"`
	Position             *NetVector3 `protobuf:"bytes,2,opt,name=position,proto3" json:"position,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Packet_UpdatePositionPacket) Reset()         { *m = Packet_UpdatePositionPacket{} }
func (m *Packet_UpdatePositionPacket) String() string { return proto.CompactTextString(m) }
func (*Packet_UpdatePositionPacket) ProtoMessage()    {}
func (*Packet_UpdatePositionPacket) Descriptor() ([]byte, []int) {
	return fileDescriptor_dcbdca058206953b, []int{4, 1}
}

func (m *Packet_UpdatePositionPacket) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Packet_UpdatePositionPacket.Unmarshal(m, b)
}
func (m *Packet_UpdatePositionPacket) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Packet_UpdatePositionPacket.Marshal(b, m, deterministic)
}
func (m *Packet_UpdatePositionPacket) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Packet_UpdatePositionPacket.Merge(m, src)
}
func (m *Packet_UpdatePositionPacket) XXX_Size() int {
	return xxx_messageInfo_Packet_UpdatePositionPacket.Size(m)
}
func (m *Packet_UpdatePositionPacket) XXX_DiscardUnknown() {
	xxx_messageInfo_Packet_UpdatePositionPacket.DiscardUnknown(m)
}

var xxx_messageInfo_Packet_UpdatePositionPacket proto.InternalMessageInfo

func (m *Packet_UpdatePositionPacket) GetObjectId() string {
	if m != nil {
		return m.ObjectId
	}
	return ""
}

func (m *Packet_UpdatePositionPacket) GetPosition() *NetVector3 {
	if m != nil {
		return m.Position
	}
	return nil
}

type Packet_DestroyObjectPacket struct {
	ObjectId             string   `protobuf:"bytes,1,opt,name=object_id,json=objectId,proto3" json:"object_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Packet_DestroyObjectPacket) Reset()         { *m = Packet_DestroyObjectPacket{} }
func (m *Packet_DestroyObjectPacket) String() string { return proto.CompactTextString(m) }
func (*Packet_DestroyObjectPacket) ProtoMessage()    {}
func (*Packet_DestroyObjectPacket) Descriptor() ([]byte, []int) {
	return fileDescriptor_dcbdca058206953b, []int{4, 2}
}

func (m *Packet_DestroyObjectPacket) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Packet_DestroyObjectPacket.Unmarshal(m, b)
}
func (m *Packet_DestroyObjectPacket) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Packet_DestroyObjectPacket.Marshal(b, m, deterministic)
}
func (m *Packet_DestroyObjectPacket) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Packet_DestroyObjectPacket.Merge(m, src)
}
func (m *Packet_DestroyObjectPacket) XXX_Size() int {
	return xxx_messageInfo_Packet_DestroyObjectPacket.Size(m)
}
func (m *Packet_DestroyObjectPacket) XXX_DiscardUnknown() {
	xxx_messageInfo_Packet_DestroyObjectPacket.DiscardUnknown(m)
}

var xxx_messageInfo_Packet_DestroyObjectPacket proto.InternalMessageInfo

func (m *Packet_DestroyObjectPacket) GetObjectId() string {
	if m != nil {
		return m.ObjectId
	}
	return ""
}

func init() {
	proto.RegisterEnum("erutan.NetObject_Type", NetObject_Type_name, NetObject_Type_value)
	proto.RegisterType((*Metadata)(nil), "erutan.Metadata")
	proto.RegisterType((*NetVector3)(nil), "erutan.NetVector3")
	proto.RegisterType((*NetQuaternion)(nil), "erutan.NetQuaternion")
	proto.RegisterType((*NetObject)(nil), "erutan.NetObject")
	proto.RegisterType((*Packet)(nil), "erutan.Packet")
	proto.RegisterType((*Packet_CreateObjectPacket)(nil), "erutan.Packet.CreateObjectPacket")
	proto.RegisterType((*Packet_UpdatePositionPacket)(nil), "erutan.Packet.UpdatePositionPacket")
	proto.RegisterType((*Packet_DestroyObjectPacket)(nil), "erutan.Packet.DestroyObjectPacket")
}

func init() { proto.RegisterFile("realtime.proto", fileDescriptor_dcbdca058206953b) }

var fileDescriptor_dcbdca058206953b = []byte{
	// 551 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x94, 0x5f, 0x6f, 0xd3, 0x3c,
	0x14, 0xc6, 0x9b, 0xac, 0xcb, 0x9b, 0x9e, 0xae, 0x7d, 0x8b, 0x61, 0x28, 0x84, 0x0b, 0x46, 0xb8,
	0x29, 0x08, 0x79, 0x90, 0x49, 0x68, 0x77, 0x68, 0xa5, 0xb0, 0x55, 0xb0, 0xb4, 0x64, 0x1b, 0xb7,
	0x93, 0x9b, 0x98, 0x29, 0xd0, 0xd6, 0x91, 0x73, 0xaa, 0xae, 0xfb, 0x40, 0x5c, 0xf0, 0x25, 0x41,
	0xb1, 0x93, 0xb6, 0xeb, 0xc6, 0xe0, 0x2e, 0xcf, 0xf1, 0xcf, 0x8f, 0xcf, 0x3f, 0x05, 0x9a, 0x92,
	0xb3, 0x11, 0x26, 0x63, 0x4e, 0x53, 0x29, 0x50, 0x10, 0x8b, 0xcb, 0x29, 0xb2, 0x89, 0xfb, 0xe4,
	0x42, 0x88, 0x8b, 0x11, 0xdf, 0x55, 0xd1, 0xe1, 0xf4, 0xeb, 0x6e, 0xce, 0x64, 0xc8, 0xc6, 0xa9,
	0x06, 0xbd, 0x2e, 0xd8, 0xc7, 0x1c, 0x59, 0xcc, 0x90, 0x91, 0x7d, 0xa8, 0x2d, 0x8e, 0x1d, 0x63,
	0xc7, 0x68, 0xd7, 0x7d, 0x97, 0x6a, 0x03, 0x5a, 0x1a, 0xd0, 0xd3, 0x92, 0x08, 0x97, 0xb0, 0xf7,
	0x06, 0x20, 0xe0, 0xf8, 0x85, 0x47, 0x28, 0xe4, 0x1e, 0xd9, 0x02, 0xe3, 0x52, 0xdd, 0x37, 0x43,
	0xe3, 0x32, 0x57, 0x73, 0xc7, 0xd4, 0x6a, 0x9e, 0xab, 0x2b, 0x67, 0x43, 0xab, 0x2b, 0xef, 0x10,
	0x1a, 0x01, 0xc7, 0xcf, 0x53, 0x86, 0x5c, 0x4e, 0x12, 0x31, 0xf9, 0xf7, 0xab, 0xb9, 0x9a, 0x39,
	0x55, 0xad, 0x66, 0xde, 0x0f, 0x13, 0x6a, 0x01, 0xc7, 0xfe, 0xf0, 0x1b, 0x8f, 0x90, 0x3c, 0x86,
	0x9a, 0x50, 0x5f, 0xe7, 0x49, 0xac, 0xdc, 0x6a, 0xa1, 0xad, 0x03, 0xbd, 0x98, 0x3c, 0x02, 0x5b,
	0xcc, 0x26, 0x5c, 0xe6, 0x67, 0xa6, 0x3a, 0xfb, 0x4f, 0xe9, 0x5e, 0x4c, 0x28, 0xd8, 0xa9, 0xc8,
	0x12, 0x4c, 0xc4, 0x44, 0x3d, 0x54, 0xf7, 0x09, 0xd5, 0x8d, 0xa4, 0xcb, 0xf2, 0xc2, 0x05, 0x43,
	0x5e, 0x83, 0x2d, 0x05, 0x32, 0xc5, 0x57, 0x15, 0xbf, 0xbd, 0xc2, 0x2f, 0xcb, 0x0a, 0x17, 0x18,
	0x69, 0xc3, 0x66, 0x16, 0xb1, 0x11, 0x77, 0x36, 0xff, 0xe8, 0xaf, 0x01, 0xf2, 0x02, 0xaa, 0x38,
	0x4f, 0xb9, 0x63, 0xed, 0x18, 0xed, 0xa6, 0xff, 0x70, 0x05, 0xd4, 0x55, 0xd2, 0xd3, 0x79, 0xca,
	0x43, 0xc5, 0x78, 0x6d, 0xa8, 0xe6, 0x8a, 0x00, 0x58, 0x07, 0x41, 0xef, 0xf8, 0xe0, 0x53, 0xab,
	0x42, 0x6c, 0xa8, 0x7e, 0xe8, 0xf7, 0xbb, 0x2d, 0x23, 0x8f, 0x1e, 0x86, 0xfd, 0xb3, 0xa0, 0xdb,
	0x32, 0xbd, 0x5f, 0x1b, 0x60, 0x0d, 0x58, 0xf4, 0x9d, 0x23, 0x79, 0x09, 0xf6, 0xb8, 0x18, 0x7d,
	0x31, 0xed, 0x56, 0xf9, 0x48, 0xb9, 0x12, 0xe1, 0x82, 0x20, 0x47, 0xd0, 0x88, 0x24, 0x67, 0xc8,
	0xcf, 0x75, 0x27, 0x55, 0xef, 0xea, 0xfe, 0xd3, 0xf2, 0x8a, 0x36, 0xa5, 0xef, 0x14, 0xa3, 0x33,
	0xd4, 0xa1, 0xa3, 0x4a, 0xb8, 0x15, 0xad, 0x44, 0x49, 0x00, 0xff, 0x4f, 0xd3, 0x38, 0x77, 0x5a,
	0x6b, 0xf6, 0xb3, 0x35, 0xaf, 0x33, 0x45, 0x0d, 0x0a, 0x68, 0xe1, 0xd6, 0x9c, 0x5e, 0x8b, 0x93,
	0x8f, 0xd0, 0x8c, 0x79, 0x86, 0x52, 0xcc, 0xcb, 0xd4, 0xf4, 0x2c, 0xbc, 0x35, 0xbb, 0xae, 0x86,
	0xd6, 0x72, 0x6b, 0xc4, 0xab, 0x61, 0xf7, 0x2d, 0x90, 0x9b, 0x25, 0x90, 0xe7, 0x60, 0x15, 0xd6,
	0xba, 0x51, 0xf7, 0x6e, 0x4c, 0x23, 0x2c, 0x00, 0x37, 0x82, 0x07, 0xb7, 0xe5, 0x7d, 0xf7, 0x4e,
	0xae, 0x2e, 0x9e, 0xf9, 0xf7, 0xc5, 0x73, 0x7d, 0xb8, 0x7f, 0x4b, 0x35, 0x77, 0xbe, 0xd1, 0xb1,
	0xf4, 0x3e, 0xf9, 0xfb, 0x60, 0xbd, 0x57, 0xd6, 0x84, 0x82, 0x75, 0x82, 0x92, 0xb3, 0x31, 0x69,
	0x5e, 0x6f, 0x95, 0xbb, 0xa6, 0xbd, 0x4a, 0xdb, 0x78, 0x65, 0x74, 0x28, 0x40, 0x24, 0xc6, 0xc5,
	0x51, 0xa7, 0x70, 0x19, 0x18, 0x3f, 0xcd, 0x6d, 0xfd, 0x49, 0x4f, 0x22, 0x99, 0xa4, 0x98, 0xd1,
	0x41, 0xfe, 0xb3, 0xc8, 0x86, 0x96, 0xfa, 0x69, 0xec, 0xfd, 0x0e, 0x00, 0x00, 0xff, 0xff, 0x72,
	0xdf, 0xe4, 0x25, 0x9d, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ErutanClient is the client API for Erutan service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ErutanClient interface {
	//
	//Blabla
	Stream(ctx context.Context, opts ...grpc.CallOption) (Erutan_StreamClient, error)
}

type erutanClient struct {
	cc *grpc.ClientConn
}

func NewErutanClient(cc *grpc.ClientConn) ErutanClient {
	return &erutanClient{cc}
}

func (c *erutanClient) Stream(ctx context.Context, opts ...grpc.CallOption) (Erutan_StreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Erutan_serviceDesc.Streams[0], "/erutan.Erutan/Stream", opts...)
	if err != nil {
		return nil, err
	}
	x := &erutanStreamClient{stream}
	return x, nil
}

type Erutan_StreamClient interface {
	Send(*Packet) error
	Recv() (*Packet, error)
	grpc.ClientStream
}

type erutanStreamClient struct {
	grpc.ClientStream
}

func (x *erutanStreamClient) Send(m *Packet) error {
	return x.ClientStream.SendMsg(m)
}

func (x *erutanStreamClient) Recv() (*Packet, error) {
	m := new(Packet)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ErutanServer is the server API for Erutan service.
type ErutanServer interface {
	//
	//Blabla
	Stream(Erutan_StreamServer) error
}

// UnimplementedErutanServer can be embedded to have forward compatible implementations.
type UnimplementedErutanServer struct {
}

func (*UnimplementedErutanServer) Stream(srv Erutan_StreamServer) error {
	return status.Errorf(codes.Unimplemented, "method Stream not implemented")
}

func RegisterErutanServer(s *grpc.Server, srv ErutanServer) {
	s.RegisterService(&_Erutan_serviceDesc, srv)
}

func _Erutan_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ErutanServer).Stream(&erutanStreamServer{stream})
}

type Erutan_StreamServer interface {
	Send(*Packet) error
	Recv() (*Packet, error)
	grpc.ServerStream
}

type erutanStreamServer struct {
	grpc.ServerStream
}

func (x *erutanStreamServer) Send(m *Packet) error {
	return x.ServerStream.SendMsg(m)
}

func (x *erutanStreamServer) Recv() (*Packet, error) {
	m := new(Packet)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Erutan_serviceDesc = grpc.ServiceDesc{
	ServiceName: "erutan.Erutan",
	HandlerType: (*ErutanServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _Erutan_Stream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "realtime.proto",
}
