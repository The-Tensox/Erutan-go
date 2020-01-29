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
	X                    float64  `protobuf:"fixed64,1,opt,name=x,proto3" json:"x,omitempty"`
	Y                    float64  `protobuf:"fixed64,2,opt,name=y,proto3" json:"y,omitempty"`
	Z                    float64  `protobuf:"fixed64,3,opt,name=z,proto3" json:"z,omitempty"`
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

func (m *NetVector3) GetX() float64 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *NetVector3) GetY() float64 {
	if m != nil {
		return m.Y
	}
	return 0
}

func (m *NetVector3) GetZ() float64 {
	if m != nil {
		return m.Z
	}
	return 0
}

type NetQuaternion struct {
	X                    float64  `protobuf:"fixed64,1,opt,name=x,proto3" json:"x,omitempty"`
	Y                    float64  `protobuf:"fixed64,2,opt,name=y,proto3" json:"y,omitempty"`
	Z                    float64  `protobuf:"fixed64,3,opt,name=z,proto3" json:"z,omitempty"`
	W                    float64  `protobuf:"fixed64,4,opt,name=w,proto3" json:"w,omitempty"`
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

func (m *NetQuaternion) GetX() float64 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *NetQuaternion) GetY() float64 {
	if m != nil {
		return m.Y
	}
	return 0
}

func (m *NetQuaternion) GetZ() float64 {
	if m != nil {
		return m.Z
	}
	return 0
}

func (m *NetQuaternion) GetW() float64 {
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
	Components           []*Component   `protobuf:"bytes,7,rep,name=components,proto3" json:"components,omitempty"`
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

func (m *NetObject) GetComponents() []*Component {
	if m != nil {
		return m.Components
	}
	return nil
}

// Data-oriented design like
type Component struct {
	// Types that are valid to be assigned to Type:
	//	*Component_Animal
	//	*Component_Food
	//	*Component_Ground
	Type                 isComponent_Type `protobuf_oneof:"type"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Component) Reset()         { *m = Component{} }
func (m *Component) String() string { return proto.CompactTextString(m) }
func (*Component) ProtoMessage()    {}
func (*Component) Descriptor() ([]byte, []int) {
	return fileDescriptor_dcbdca058206953b, []int{4}
}

func (m *Component) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Component.Unmarshal(m, b)
}
func (m *Component) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Component.Marshal(b, m, deterministic)
}
func (m *Component) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Component.Merge(m, src)
}
func (m *Component) XXX_Size() int {
	return xxx_messageInfo_Component.Size(m)
}
func (m *Component) XXX_DiscardUnknown() {
	xxx_messageInfo_Component.DiscardUnknown(m)
}

var xxx_messageInfo_Component proto.InternalMessageInfo

type isComponent_Type interface {
	isComponent_Type()
}

type Component_Animal struct {
	Animal *Component_AnimalComponent `protobuf:"bytes,1,opt,name=animal,proto3,oneof"`
}

type Component_Food struct {
	Food *Component_FoodComponent `protobuf:"bytes,2,opt,name=food,proto3,oneof"`
}

type Component_Ground struct {
	Ground *Component_GroundComponent `protobuf:"bytes,3,opt,name=ground,proto3,oneof"`
}

func (*Component_Animal) isComponent_Type() {}

func (*Component_Food) isComponent_Type() {}

func (*Component_Ground) isComponent_Type() {}

func (m *Component) GetType() isComponent_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (m *Component) GetAnimal() *Component_AnimalComponent {
	if x, ok := m.GetType().(*Component_Animal); ok {
		return x.Animal
	}
	return nil
}

func (m *Component) GetFood() *Component_FoodComponent {
	if x, ok := m.GetType().(*Component_Food); ok {
		return x.Food
	}
	return nil
}

func (m *Component) GetGround() *Component_GroundComponent {
	if x, ok := m.GetType().(*Component_Ground); ok {
		return x.Ground
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Component) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Component_Animal)(nil),
		(*Component_Food)(nil),
		(*Component_Ground)(nil),
	}
}

type Component_AnimalComponent struct {
	Life                 float64    `protobuf:"fixed64,1,opt,name=life,proto3" json:"life,omitempty"`
	Food                 *NetObject `protobuf:"bytes,2,opt,name=food,proto3" json:"food,omitempty"`
	Target               *NetObject `protobuf:"bytes,3,opt,name=target,proto3" json:"target,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Component_AnimalComponent) Reset()         { *m = Component_AnimalComponent{} }
func (m *Component_AnimalComponent) String() string { return proto.CompactTextString(m) }
func (*Component_AnimalComponent) ProtoMessage()    {}
func (*Component_AnimalComponent) Descriptor() ([]byte, []int) {
	return fileDescriptor_dcbdca058206953b, []int{4, 0}
}

func (m *Component_AnimalComponent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Component_AnimalComponent.Unmarshal(m, b)
}
func (m *Component_AnimalComponent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Component_AnimalComponent.Marshal(b, m, deterministic)
}
func (m *Component_AnimalComponent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Component_AnimalComponent.Merge(m, src)
}
func (m *Component_AnimalComponent) XXX_Size() int {
	return xxx_messageInfo_Component_AnimalComponent.Size(m)
}
func (m *Component_AnimalComponent) XXX_DiscardUnknown() {
	xxx_messageInfo_Component_AnimalComponent.DiscardUnknown(m)
}

var xxx_messageInfo_Component_AnimalComponent proto.InternalMessageInfo

func (m *Component_AnimalComponent) GetLife() float64 {
	if m != nil {
		return m.Life
	}
	return 0
}

func (m *Component_AnimalComponent) GetFood() *NetObject {
	if m != nil {
		return m.Food
	}
	return nil
}

func (m *Component_AnimalComponent) GetTarget() *NetObject {
	if m != nil {
		return m.Target
	}
	return nil
}

type Component_FoodComponent struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Component_FoodComponent) Reset()         { *m = Component_FoodComponent{} }
func (m *Component_FoodComponent) String() string { return proto.CompactTextString(m) }
func (*Component_FoodComponent) ProtoMessage()    {}
func (*Component_FoodComponent) Descriptor() ([]byte, []int) {
	return fileDescriptor_dcbdca058206953b, []int{4, 1}
}

func (m *Component_FoodComponent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Component_FoodComponent.Unmarshal(m, b)
}
func (m *Component_FoodComponent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Component_FoodComponent.Marshal(b, m, deterministic)
}
func (m *Component_FoodComponent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Component_FoodComponent.Merge(m, src)
}
func (m *Component_FoodComponent) XXX_Size() int {
	return xxx_messageInfo_Component_FoodComponent.Size(m)
}
func (m *Component_FoodComponent) XXX_DiscardUnknown() {
	xxx_messageInfo_Component_FoodComponent.DiscardUnknown(m)
}

var xxx_messageInfo_Component_FoodComponent proto.InternalMessageInfo

type Component_GroundComponent struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Component_GroundComponent) Reset()         { *m = Component_GroundComponent{} }
func (m *Component_GroundComponent) String() string { return proto.CompactTextString(m) }
func (*Component_GroundComponent) ProtoMessage()    {}
func (*Component_GroundComponent) Descriptor() ([]byte, []int) {
	return fileDescriptor_dcbdca058206953b, []int{4, 2}
}

func (m *Component_GroundComponent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Component_GroundComponent.Unmarshal(m, b)
}
func (m *Component_GroundComponent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Component_GroundComponent.Marshal(b, m, deterministic)
}
func (m *Component_GroundComponent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Component_GroundComponent.Merge(m, src)
}
func (m *Component_GroundComponent) XXX_Size() int {
	return xxx_messageInfo_Component_GroundComponent.Size(m)
}
func (m *Component_GroundComponent) XXX_DiscardUnknown() {
	xxx_messageInfo_Component_GroundComponent.DiscardUnknown(m)
}

var xxx_messageInfo_Component_GroundComponent proto.InternalMessageInfo

type Packet struct {
	Metadata *Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	// Types that are valid to be assigned to Type:
	//	*Packet_CreateObject
	//	*Packet_UpdatePosition
	//	*Packet_UpdateRotation
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
	return fileDescriptor_dcbdca058206953b, []int{5}
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

type Packet_UpdateRotation struct {
	UpdateRotation *Packet_UpdateRotationPacket `protobuf:"bytes,4,opt,name=update_rotation,json=updateRotation,proto3,oneof"`
}

type Packet_DestroyObject struct {
	DestroyObject *Packet_DestroyObjectPacket `protobuf:"bytes,5,opt,name=destroy_object,json=destroyObject,proto3,oneof"`
}

func (*Packet_CreateObject) isPacket_Type() {}

func (*Packet_UpdatePosition) isPacket_Type() {}

func (*Packet_UpdateRotation) isPacket_Type() {}

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

func (m *Packet) GetUpdateRotation() *Packet_UpdateRotationPacket {
	if x, ok := m.GetType().(*Packet_UpdateRotation); ok {
		return x.UpdateRotation
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
		(*Packet_UpdateRotation)(nil),
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
	return fileDescriptor_dcbdca058206953b, []int{5, 0}
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
	return fileDescriptor_dcbdca058206953b, []int{5, 1}
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

type Packet_UpdateRotationPacket struct {
	ObjectId             string         `protobuf:"bytes,1,opt,name=object_id,json=objectId,proto3" json:"object_id,omitempty"`
	Rotation             *NetQuaternion `protobuf:"bytes,2,opt,name=rotation,proto3" json:"rotation,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Packet_UpdateRotationPacket) Reset()         { *m = Packet_UpdateRotationPacket{} }
func (m *Packet_UpdateRotationPacket) String() string { return proto.CompactTextString(m) }
func (*Packet_UpdateRotationPacket) ProtoMessage()    {}
func (*Packet_UpdateRotationPacket) Descriptor() ([]byte, []int) {
	return fileDescriptor_dcbdca058206953b, []int{5, 2}
}

func (m *Packet_UpdateRotationPacket) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Packet_UpdateRotationPacket.Unmarshal(m, b)
}
func (m *Packet_UpdateRotationPacket) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Packet_UpdateRotationPacket.Marshal(b, m, deterministic)
}
func (m *Packet_UpdateRotationPacket) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Packet_UpdateRotationPacket.Merge(m, src)
}
func (m *Packet_UpdateRotationPacket) XXX_Size() int {
	return xxx_messageInfo_Packet_UpdateRotationPacket.Size(m)
}
func (m *Packet_UpdateRotationPacket) XXX_DiscardUnknown() {
	xxx_messageInfo_Packet_UpdateRotationPacket.DiscardUnknown(m)
}

var xxx_messageInfo_Packet_UpdateRotationPacket proto.InternalMessageInfo

func (m *Packet_UpdateRotationPacket) GetObjectId() string {
	if m != nil {
		return m.ObjectId
	}
	return ""
}

func (m *Packet_UpdateRotationPacket) GetRotation() *NetQuaternion {
	if m != nil {
		return m.Rotation
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
	return fileDescriptor_dcbdca058206953b, []int{5, 3}
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
	proto.RegisterType((*Component)(nil), "erutan.Component")
	proto.RegisterType((*Component_AnimalComponent)(nil), "erutan.Component.AnimalComponent")
	proto.RegisterType((*Component_FoodComponent)(nil), "erutan.Component.FoodComponent")
	proto.RegisterType((*Component_GroundComponent)(nil), "erutan.Component.GroundComponent")
	proto.RegisterType((*Packet)(nil), "erutan.Packet")
	proto.RegisterType((*Packet_CreateObjectPacket)(nil), "erutan.Packet.CreateObjectPacket")
	proto.RegisterType((*Packet_UpdatePositionPacket)(nil), "erutan.Packet.UpdatePositionPacket")
	proto.RegisterType((*Packet_UpdateRotationPacket)(nil), "erutan.Packet.UpdateRotationPacket")
	proto.RegisterType((*Packet_DestroyObjectPacket)(nil), "erutan.Packet.DestroyObjectPacket")
}

func init() { proto.RegisterFile("realtime.proto", fileDescriptor_dcbdca058206953b) }

var fileDescriptor_dcbdca058206953b = []byte{
	// 721 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x95, 0xdd, 0x4e, 0xdb, 0x4a,
	0x10, 0xc7, 0x63, 0x13, 0x4c, 0x32, 0x21, 0x1f, 0xec, 0x39, 0x1c, 0xe5, 0xf8, 0x5c, 0xc0, 0x71,
	0x55, 0x29, 0xad, 0xaa, 0xa5, 0x04, 0xb5, 0x42, 0xea, 0x45, 0x05, 0xa4, 0x7c, 0xa8, 0x25, 0x49,
	0x0d, 0xf4, 0x16, 0x2d, 0xf6, 0x26, 0x72, 0x9b, 0x78, 0x2d, 0x7b, 0xa3, 0x10, 0x1e, 0xa9, 0x2f,
	0xd0, 0x67, 0xe9, 0x65, 0xdf, 0xa4, 0xf2, 0xee, 0xda, 0x38, 0x4e, 0xa0, 0xdc, 0x79, 0x66, 0x7e,
	0xf3, 0x9f, 0xd9, 0xd9, 0x0f, 0x43, 0x2d, 0xa4, 0x64, 0xc4, 0xbd, 0x31, 0xc5, 0x41, 0xc8, 0x38,
	0x43, 0x06, 0x0d, 0x27, 0x9c, 0xf8, 0xe6, 0xd6, 0x90, 0xb1, 0xe1, 0x88, 0xee, 0x08, 0xef, 0xcd,
	0x64, 0xb0, 0x13, 0x33, 0x11, 0x27, 0xe3, 0x40, 0x82, 0x56, 0x07, 0x4a, 0xe7, 0x94, 0x13, 0x97,
	0x70, 0x82, 0xf6, 0xa1, 0x9c, 0x86, 0x9b, 0xda, 0xb6, 0xd6, 0xaa, 0xb4, 0x4d, 0x2c, 0x05, 0x70,
	0x22, 0x80, 0x2f, 0x13, 0xc2, 0xbe, 0x87, 0xad, 0xb7, 0x00, 0x5d, 0xca, 0xbf, 0x50, 0x87, 0xb3,
	0x70, 0x0f, 0xad, 0x83, 0x76, 0x2b, 0xf2, 0x35, 0x5b, 0xbb, 0x8d, 0xad, 0x59, 0x53, 0x97, 0xd6,
	0x2c, 0xb6, 0xee, 0x9a, 0x2b, 0xd2, 0xba, 0xb3, 0x4e, 0xa0, 0xda, 0xa5, 0xfc, 0xf3, 0x84, 0x70,
	0x1a, 0xfa, 0x1e, 0xf3, 0x9f, 0x9e, 0x1a, 0x5b, 0xd3, 0x66, 0x51, 0x5a, 0x53, 0xeb, 0xa7, 0x0e,
	0xe5, 0x2e, 0xe5, 0xbd, 0x9b, 0xaf, 0xd4, 0xe1, 0xe8, 0x3f, 0x28, 0x33, 0xf1, 0x75, 0xed, 0xb9,
	0x42, 0xad, 0x6c, 0x97, 0xa4, 0xe3, 0xcc, 0x45, 0xff, 0x42, 0x89, 0x4d, 0x7d, 0x1a, 0xc6, 0x31,
	0x5d, 0xc4, 0xd6, 0x84, 0x7d, 0xe6, 0x22, 0x0c, 0xa5, 0x80, 0x45, 0x1e, 0xf7, 0x98, 0x2f, 0x0a,
	0x55, 0xda, 0x08, 0xcb, 0x41, 0xe2, 0xfb, 0xe5, 0xd9, 0x29, 0x83, 0x76, 0xa1, 0x14, 0x32, 0x4e,
	0x04, 0x5f, 0x14, 0xfc, 0x66, 0x86, 0xbf, 0x5f, 0x96, 0x9d, 0x62, 0xa8, 0x05, 0xab, 0x91, 0x43,
	0x46, 0xb4, 0xb9, 0xfa, 0xa0, 0xbe, 0x04, 0xd0, 0x4b, 0x28, 0xf2, 0x59, 0x40, 0x9b, 0xc6, 0xb6,
	0xd6, 0xaa, 0xb5, 0xff, 0xc9, 0x80, 0x72, 0x95, 0xf8, 0x72, 0x16, 0x50, 0x5b, 0x30, 0x68, 0x17,
	0xc0, 0x61, 0xe3, 0x80, 0xf9, 0xd4, 0xe7, 0x51, 0x73, 0x6d, 0x7b, 0xa5, 0x55, 0x69, 0x6f, 0x24,
	0x19, 0x47, 0x49, 0xc4, 0xce, 0x40, 0x56, 0x0b, 0x8a, 0xb1, 0x00, 0x02, 0x30, 0x0e, 0xba, 0x67,
	0xe7, 0x07, 0x9f, 0x1a, 0x05, 0x54, 0x82, 0xe2, 0x71, 0xaf, 0xd7, 0x69, 0x68, 0xb1, 0xf7, 0xc4,
	0xee, 0x5d, 0x75, 0x3b, 0x0d, 0xdd, 0xfa, 0xa5, 0x43, 0x39, 0xd5, 0x40, 0xef, 0xc0, 0x20, 0xbe,
	0x37, 0x26, 0x23, 0x75, 0x42, 0xfe, 0x5f, 0x28, 0x83, 0x0f, 0x44, 0x3c, 0xb5, 0x4f, 0x0b, 0xb6,
	0x4a, 0x41, 0x6f, 0xa0, 0x38, 0x60, 0x4c, 0xce, 0xbd, 0xd2, 0xde, 0x5a, 0x4c, 0x3d, 0x66, 0xcc,
	0xcd, 0x26, 0x0a, 0x3c, 0xae, 0x39, 0x0c, 0xd9, 0xc4, 0x77, 0xd5, 0xae, 0x2c, 0xa9, 0x79, 0x22,
	0xe2, 0x73, 0x35, 0x65, 0x8a, 0x39, 0x85, 0x7a, 0xae, 0x21, 0x84, 0xa0, 0x38, 0xf2, 0x06, 0x54,
	0x1d, 0x34, 0xf1, 0x8d, 0x9e, 0xcf, 0xb5, 0xb6, 0xb1, 0x30, 0x6e, 0xd5, 0xca, 0x0b, 0x30, 0x38,
	0x09, 0x87, 0x94, 0xab, 0x56, 0x96, 0x80, 0x0a, 0x30, 0xeb, 0x50, 0x9d, 0x5b, 0x8e, 0xb9, 0x01,
	0xf5, 0x5c, 0x9b, 0x87, 0x86, 0xdc, 0x64, 0xeb, 0xc7, 0x2a, 0x18, 0x7d, 0xe2, 0x7c, 0xa3, 0x1c,
	0xbd, 0x82, 0xd2, 0x58, 0xdd, 0x48, 0x35, 0xe2, 0x46, 0x52, 0x23, 0xb9, 0xa9, 0x76, 0x4a, 0xa0,
	0x53, 0xa8, 0x3a, 0x21, 0x25, 0x9c, 0x5e, 0xcb, 0x03, 0xae, 0xfa, 0x4f, 0x27, 0x24, 0x45, 0xf1,
	0x91, 0x60, 0x64, 0x83, 0xd2, 0x75, 0x5a, 0xb0, 0xd7, 0x9d, 0x8c, 0x17, 0x75, 0xa1, 0x3e, 0x09,
	0xdc, 0x58, 0x29, 0x77, 0x07, 0x9e, 0xe5, 0xb4, 0xae, 0x04, 0xd5, 0x57, 0x50, 0xaa, 0x56, 0x9b,
	0xcc, 0xf9, 0x33, 0x7a, 0xb9, 0x3b, 0xb2, 0x5c, 0xcf, 0x56, 0x50, 0x5e, 0x2f, 0xf1, 0xa3, 0x8f,
	0x50, 0x73, 0x69, 0xc4, 0x43, 0x36, 0x4b, 0x96, 0x2a, 0xaf, 0x90, 0x95, 0x93, 0xeb, 0x48, 0x28,
	0xb7, 0xd6, 0xaa, 0x9b, 0x75, 0x9b, 0xef, 0x01, 0x2d, 0x8e, 0x24, 0xde, 0x5c, 0x25, 0xad, 0x3d,
	0xb8, 0xb9, 0x12, 0x30, 0x1d, 0xf8, 0x7b, 0xd9, 0x1c, 0x1e, 0x7f, 0x7a, 0xb2, 0xef, 0x8b, 0xfe,
	0xe7, 0xf7, 0xc5, 0x1c, 0x24, 0x45, 0xe6, 0x87, 0xf3, 0x78, 0x91, 0xec, 0xa3, 0xa4, 0x3f, 0xe9,
	0x51, 0x32, 0xdb, 0xf0, 0xd7, 0x92, 0xa9, 0x3d, 0x5a, 0x26, 0x39, 0xb9, 0xed, 0x7d, 0x30, 0x3e,
	0x08, 0x75, 0x84, 0xc1, 0xb8, 0xe0, 0x21, 0x25, 0x63, 0x54, 0x9b, 0xdf, 0x12, 0x33, 0x67, 0x5b,
	0x85, 0x96, 0xf6, 0x5a, 0x3b, 0xc4, 0xe2, 0xd1, 0x52, 0xa1, 0x43, 0xa5, 0xd2, 0xd7, 0xbe, 0xeb,
	0x9b, 0xf2, 0x13, 0x5f, 0x38, 0xa1, 0x17, 0xf0, 0x08, 0xf7, 0xe3, 0x7f, 0x4f, 0x74, 0x63, 0x88,
	0x7f, 0xd0, 0xde, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0x39, 0xf1, 0x9c, 0xa2, 0xec, 0x06, 0x00,
	0x00,
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
