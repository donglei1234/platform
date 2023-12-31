// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: iap/api/message.proto

package pb

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

type SYS int32

const (
	SYS_IOS     SYS = 0
	SYS_ANDROID SYS = 1
)

// Enum value maps for SYS.
var (
	SYS_name = map[int32]string{
		0: "IOS",
		1: "ANDROID",
	}
	SYS_value = map[string]int32{
		"IOS":     0,
		"ANDROID": 1,
	}
)

func (x SYS) Enum() *SYS {
	p := new(SYS)
	*p = x
	return p
}

func (x SYS) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SYS) Descriptor() protoreflect.EnumDescriptor {
	return file_iap_api_common_proto_enumTypes[0].Descriptor()
}

func (SYS) Type() protoreflect.EnumType {
	return &file_iap_api_common_proto_enumTypes[0]
}

func (x SYS) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SYS.Descriptor instead.
func (SYS) EnumDescriptor() ([]byte, []int) {
	return file_iap_api_common_proto_rawDescGZIP(), []int{0}
}

type IAPRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AppStoreId   string `protobuf:"bytes,1,opt,name=appStoreId,proto3" json:"appStoreId,omitempty"`     // app的名字
	Sys          SYS    `protobuf:"varint,2,opt,name=sys,proto3,enum=iap.pb.SYS" json:"sys,omitempty"`  //客户端
	ProductType  int32  `protobuf:"varint,3,opt,name=productType,proto3" json:"productType,omitempty"`  //产品类型
	ProductId    string `protobuf:"bytes,4,opt,name=productId,proto3" json:"productId,omitempty"`       //产品ID
	ProductToken string `protobuf:"bytes,5,opt,name=productToken,proto3" json:"productToken,omitempty"` //Token
}

func (x *IAPRequest) Reset() {
	*x = IAPRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iap_api_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IAPRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IAPRequest) ProtoMessage() {}

func (x *IAPRequest) ProtoReflect() protoreflect.Message {
	mi := &file_iap_api_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IAPRequest.ProtoReflect.Descriptor instead.
func (*IAPRequest) Descriptor() ([]byte, []int) {
	return file_iap_api_common_proto_rawDescGZIP(), []int{0}
}

func (x *IAPRequest) GetAppStoreId() string {
	if x != nil {
		return x.AppStoreId
	}
	return ""
}

func (x *IAPRequest) GetSys() SYS {
	if x != nil {
		return x.Sys
	}
	return SYS_IOS
}

func (x *IAPRequest) GetProductType() int32 {
	if x != nil {
		return x.ProductType
	}
	return 0
}

func (x *IAPRequest) GetProductId() string {
	if x != nil {
		return x.ProductId
	}
	return ""
}

func (x *IAPRequest) GetProductToken() string {
	if x != nil {
		return x.ProductToken
	}
	return ""
}

type IAPResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg  string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Sys  SYS    `protobuf:"varint,3,opt,name=sys,proto3,enum=iap.pb.SYS" json:"sys,omitempty"`
	Data string `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *IAPResponse) Reset() {
	*x = IAPResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iap_api_common_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IAPResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IAPResponse) ProtoMessage() {}

func (x *IAPResponse) ProtoReflect() protoreflect.Message {
	mi := &file_iap_api_common_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IAPResponse.ProtoReflect.Descriptor instead.
func (*IAPResponse) Descriptor() ([]byte, []int) {
	return file_iap_api_common_proto_rawDescGZIP(), []int{1}
}

func (x *IAPResponse) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *IAPResponse) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *IAPResponse) GetSys() SYS {
	if x != nil {
		return x.Sys
	}
	return SYS_IOS
}

func (x *IAPResponse) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

var File_iap_api_common_proto protoreflect.FileDescriptor

var file_iap_api_common_proto_rawDesc = []byte{
	0x0a, 0x14, 0x69, 0x61, 0x70, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x69, 0x61, 0x70, 0x2e, 0x70, 0x62, 0x22, 0xaf,
	0x01, 0x0a, 0x0a, 0x49, 0x41, 0x50, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a,
	0x0a, 0x61, 0x70, 0x70, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x61, 0x70, 0x70, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x49, 0x64, 0x12, 0x1d, 0x0a,
	0x03, 0x73, 0x79, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0b, 0x2e, 0x69, 0x61, 0x70,
	0x2e, 0x70, 0x62, 0x2e, 0x53, 0x59, 0x53, 0x52, 0x03, 0x73, 0x79, 0x73, 0x12, 0x20, 0x0a, 0x0b,
	0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1c,
	0x0a, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0c,
	0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x22, 0x66, 0x0a, 0x0b, 0x49, 0x41, 0x50, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63,
	0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x1d, 0x0a, 0x03, 0x73, 0x79, 0x73, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x0b, 0x2e, 0x69, 0x61, 0x70, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x59, 0x53, 0x52,
	0x03, 0x73, 0x79, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x2a, 0x1b, 0x0a, 0x03, 0x53, 0x59, 0x53, 0x12,
	0x07, 0x0a, 0x03, 0x49, 0x4f, 0x53, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x41, 0x4e, 0x44, 0x52,
	0x4f, 0x49, 0x44, 0x10, 0x01, 0x42, 0x0c, 0x5a, 0x0a, 0x69, 0x61, 0x70, 0x2f, 0x61, 0x70, 0x69,
	0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_iap_api_common_proto_rawDescOnce sync.Once
	file_iap_api_common_proto_rawDescData = file_iap_api_common_proto_rawDesc
)

func file_iap_api_common_proto_rawDescGZIP() []byte {
	file_iap_api_common_proto_rawDescOnce.Do(func() {
		file_iap_api_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_iap_api_common_proto_rawDescData)
	})
	return file_iap_api_common_proto_rawDescData
}

var file_iap_api_common_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_iap_api_common_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_iap_api_common_proto_goTypes = []interface{}{
	(SYS)(0),            // 0: iap.pb.SYS
	(*IAPRequest)(nil),  // 1: iap.pb.IAPRequest
	(*IAPResponse)(nil), // 2: iap.pb.IAPResponse
}
var file_iap_api_common_proto_depIdxs = []int32{
	0, // 0: iap.pb.IAPRequest.sys:type_name -> iap.pb.SYS
	0, // 1: iap.pb.IAPResponse.sys:type_name -> iap.pb.SYS
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_iap_api_common_proto_init() }
func file_iap_api_common_proto_init() {
	if File_iap_api_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_iap_api_common_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IAPRequest); i {
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
		file_iap_api_common_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IAPResponse); i {
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
			RawDescriptor: file_iap_api_common_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_iap_api_common_proto_goTypes,
		DependencyIndexes: file_iap_api_common_proto_depIdxs,
		EnumInfos:         file_iap_api_common_proto_enumTypes,
		MessageInfos:      file_iap_api_common_proto_msgTypes,
	}.Build()
	File_iap_api_common_proto = out.File
	file_iap_api_common_proto_rawDesc = nil
	file_iap_api_common_proto_goTypes = nil
	file_iap_api_common_proto_depIdxs = nil
}
