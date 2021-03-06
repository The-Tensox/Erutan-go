// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0-devel
// 	protoc        v3.11.4
// source: protobuf/protometry/volume.proto

package protometry

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Sphere struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Center *Vector3 `protobuf:"bytes,1,opt,name=center,proto3" json:"center,omitempty"`
	Radius float64  `protobuf:"fixed64,2,opt,name=radius,proto3" json:"radius,omitempty"`
}

func (x *Sphere) Reset() {
	*x = Sphere{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_protometry_volume_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Sphere) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Sphere) ProtoMessage() {}

func (x *Sphere) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_protometry_volume_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Sphere.ProtoReflect.Descriptor instead.
func (*Sphere) Descriptor() ([]byte, []int) {
	return file_protobuf_protometry_volume_proto_rawDescGZIP(), []int{0}
}

func (x *Sphere) GetCenter() *Vector3 {
	if x != nil {
		return x.Center
	}
	return nil
}

func (x *Sphere) GetRadius() float64 {
	if x != nil {
		return x.Radius
	}
	return 0
}

type Capsule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Center *Vector3 `protobuf:"bytes,1,opt,name=center,proto3" json:"center,omitempty"`
	Width  float64  `protobuf:"fixed64,2,opt,name=width,proto3" json:"width,omitempty"`
}

func (x *Capsule) Reset() {
	*x = Capsule{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_protometry_volume_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Capsule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Capsule) ProtoMessage() {}

func (x *Capsule) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_protometry_volume_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Capsule.ProtoReflect.Descriptor instead.
func (*Capsule) Descriptor() ([]byte, []int) {
	return file_protobuf_protometry_volume_proto_rawDescGZIP(), []int{1}
}

func (x *Capsule) GetCenter() *Vector3 {
	if x != nil {
		return x.Center
	}
	return nil
}

func (x *Capsule) GetWidth() float64 {
	if x != nil {
		return x.Width
	}
	return 0
}

// Box is an AABB volume
type Box struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Min *Vector3 `protobuf:"bytes,1,opt,name=min,proto3" json:"min,omitempty"`
	Max *Vector3 `protobuf:"bytes,2,opt,name=max,proto3" json:"max,omitempty"`
}

func (x *Box) Reset() {
	*x = Box{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_protometry_volume_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Box) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Box) ProtoMessage() {}

func (x *Box) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_protometry_volume_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Box.ProtoReflect.Descriptor instead.
func (*Box) Descriptor() ([]byte, []int) {
	return file_protobuf_protometry_volume_proto_rawDescGZIP(), []int{2}
}

func (x *Box) GetMin() *Vector3 {
	if x != nil {
		return x.Min
	}
	return nil
}

func (x *Box) GetMax() *Vector3 {
	if x != nil {
		return x.Max
	}
	return nil
}

type Mesh struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Center   *Vector3   `protobuf:"bytes,1,opt,name=center,proto3" json:"center,omitempty"` // I.e "pivot"
	Vertices []*Vector3 `protobuf:"bytes,2,rep,name=vertices,proto3" json:"vertices,omitempty"`
	Tris     []int32    `protobuf:"varint,3,rep,packed,name=tris,proto3" json:"tris,omitempty"`
	Normals  []*Vector3 `protobuf:"bytes,4,rep,name=normals,proto3" json:"normals,omitempty"`
	Uvs      []*Vector3 `protobuf:"bytes,5,rep,name=uvs,proto3" json:"uvs,omitempty"`
}

func (x *Mesh) Reset() {
	*x = Mesh{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_protometry_volume_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Mesh) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Mesh) ProtoMessage() {}

func (x *Mesh) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_protometry_volume_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Mesh.ProtoReflect.Descriptor instead.
func (*Mesh) Descriptor() ([]byte, []int) {
	return file_protobuf_protometry_volume_proto_rawDescGZIP(), []int{3}
}

func (x *Mesh) GetCenter() *Vector3 {
	if x != nil {
		return x.Center
	}
	return nil
}

func (x *Mesh) GetVertices() []*Vector3 {
	if x != nil {
		return x.Vertices
	}
	return nil
}

func (x *Mesh) GetTris() []int32 {
	if x != nil {
		return x.Tris
	}
	return nil
}

func (x *Mesh) GetNormals() []*Vector3 {
	if x != nil {
		return x.Normals
	}
	return nil
}

func (x *Mesh) GetUvs() []*Vector3 {
	if x != nil {
		return x.Uvs
	}
	return nil
}

var File_protobuf_protometry_volume_proto protoreflect.FileDescriptor

var file_protobuf_protometry_volume_proto_rawDesc = []byte{
	0x0a, 0x20, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x6d, 0x65, 0x74, 0x72, 0x79, 0x2f, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x1a, 0x21,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x6d, 0x65,
	0x74, 0x72, 0x79, 0x2f, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x33, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x4d, 0x0a, 0x06, 0x53, 0x70, 0x68, 0x65, 0x72, 0x65, 0x12, 0x2b, 0x0a, 0x06, 0x63,
	0x65, 0x6e, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x33,
	0x52, 0x06, 0x63, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x61, 0x64, 0x69,
	0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06, 0x72, 0x61, 0x64, 0x69, 0x75, 0x73,
	0x22, 0x4c, 0x0a, 0x07, 0x43, 0x61, 0x70, 0x73, 0x75, 0x6c, 0x65, 0x12, 0x2b, 0x0a, 0x06, 0x63,
	0x65, 0x6e, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x33,
	0x52, 0x06, 0x63, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x77, 0x69, 0x64, 0x74,
	0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x77, 0x69, 0x64, 0x74, 0x68, 0x22, 0x53,
	0x0a, 0x03, 0x42, 0x6f, 0x78, 0x12, 0x25, 0x0a, 0x03, 0x6d, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e,
	0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x33, 0x52, 0x03, 0x6d, 0x69, 0x6e, 0x12, 0x25, 0x0a, 0x03,
	0x6d, 0x61, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x33, 0x52, 0x03,
	0x6d, 0x61, 0x78, 0x22, 0xce, 0x01, 0x0a, 0x04, 0x4d, 0x65, 0x73, 0x68, 0x12, 0x2b, 0x0a, 0x06,
	0x63, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x33, 0x52, 0x06, 0x63, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x12, 0x2f, 0x0a, 0x08, 0x76, 0x65, 0x72,
	0x74, 0x69, 0x63, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x33,
	0x52, 0x08, 0x76, 0x65, 0x72, 0x74, 0x69, 0x63, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x72,
	0x69, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x05, 0x52, 0x04, 0x74, 0x72, 0x69, 0x73, 0x12, 0x2d,
	0x0a, 0x07, 0x6e, 0x6f, 0x72, 0x6d, 0x61, 0x6c, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e, 0x56, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x33, 0x52, 0x07, 0x6e, 0x6f, 0x72, 0x6d, 0x61, 0x6c, 0x73, 0x12, 0x25, 0x0a,
	0x03, 0x75, 0x76, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x33, 0x52,
	0x03, 0x75, 0x76, 0x73, 0x42, 0x2f, 0x5a, 0x20, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x54, 0x68, 0x65, 0x2d, 0x54, 0x65, 0x6e, 0x73, 0x6f, 0x78, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x6d, 0x65, 0x74, 0x72, 0x79, 0xaa, 0x02, 0x0a, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x6d, 0x65, 0x74, 0x72, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protobuf_protometry_volume_proto_rawDescOnce sync.Once
	file_protobuf_protometry_volume_proto_rawDescData = file_protobuf_protometry_volume_proto_rawDesc
)

func file_protobuf_protometry_volume_proto_rawDescGZIP() []byte {
	file_protobuf_protometry_volume_proto_rawDescOnce.Do(func() {
		file_protobuf_protometry_volume_proto_rawDescData = protoimpl.X.CompressGZIP(file_protobuf_protometry_volume_proto_rawDescData)
	})
	return file_protobuf_protometry_volume_proto_rawDescData
}

var file_protobuf_protometry_volume_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_protobuf_protometry_volume_proto_goTypes = []interface{}{
	(*Sphere)(nil),  // 0: protometry.Sphere
	(*Capsule)(nil), // 1: protometry.Capsule
	(*Box)(nil),     // 2: protometry.Box
	(*Mesh)(nil),    // 3: protometry.Mesh
	(*Vector3)(nil), // 4: protometry.Vector3
}
var file_protobuf_protometry_volume_proto_depIdxs = []int32{
	4, // 0: protometry.Sphere.center:type_name -> protometry.Vector3
	4, // 1: protometry.Capsule.center:type_name -> protometry.Vector3
	4, // 2: protometry.Box.min:type_name -> protometry.Vector3
	4, // 3: protometry.Box.max:type_name -> protometry.Vector3
	4, // 4: protometry.Mesh.center:type_name -> protometry.Vector3
	4, // 5: protometry.Mesh.vertices:type_name -> protometry.Vector3
	4, // 6: protometry.Mesh.normals:type_name -> protometry.Vector3
	4, // 7: protometry.Mesh.uvs:type_name -> protometry.Vector3
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_protobuf_protometry_volume_proto_init() }
func file_protobuf_protometry_volume_proto_init() {
	if File_protobuf_protometry_volume_proto != nil {
		return
	}
	file_protobuf_protometry_vector3_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_protobuf_protometry_volume_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Sphere); i {
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
		file_protobuf_protometry_volume_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Capsule); i {
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
		file_protobuf_protometry_volume_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Box); i {
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
		file_protobuf_protometry_volume_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Mesh); i {
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
			RawDescriptor: file_protobuf_protometry_volume_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_protobuf_protometry_volume_proto_goTypes,
		DependencyIndexes: file_protobuf_protometry_volume_proto_depIdxs,
		MessageInfos:      file_protobuf_protometry_volume_proto_msgTypes,
	}.Build()
	File_protobuf_protometry_volume_proto = out.File
	file_protobuf_protometry_volume_proto_rawDesc = nil
	file_protobuf_protometry_volume_proto_goTypes = nil
	file_protobuf_protometry_volume_proto_depIdxs = nil
}
