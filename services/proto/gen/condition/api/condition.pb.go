// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        (unknown)
// source: condition.proto

package pb

import (
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

type Condition_Status int32

const (
	Condition_NONE     Condition_Status = 0
	Condition_ACTIVE   Condition_Status = 1
	Condition_FINISHED Condition_Status = 2
)

// Enum value maps for Condition_Status.
var (
	Condition_Status_name = map[int32]string{
		0: "NONE",
		1: "ACTIVE",
		2: "FINISHED",
	}
	Condition_Status_value = map[string]int32{
		"NONE":     0,
		"ACTIVE":   1,
		"FINISHED": 2,
	}
)

func (x Condition_Status) Enum() *Condition_Status {
	p := new(Condition_Status)
	*p = x
	return p
}

func (x Condition_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Condition_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_condition_proto_enumTypes[0].Descriptor()
}

func (Condition_Status) Type() protoreflect.EnumType {
	return &file_condition_proto_enumTypes[0]
}

func (x Condition_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Condition_Status.Descriptor instead.
func (Condition_Status) EnumDescriptor() ([]byte, []int) {
	return file_condition_proto_rawDescGZIP(), []int{4, 0}
}

type Condition_UpdateStrategy int32

const (
	Condition_STRATEGY_NONE    Condition_UpdateStrategy = 0
	Condition_STRATEGY_REPLACE Condition_UpdateStrategy = 1
	Condition_STRATEGY_ADD     Condition_UpdateStrategy = 2
)

// Enum value maps for Condition_UpdateStrategy.
var (
	Condition_UpdateStrategy_name = map[int32]string{
		0: "STRATEGY_NONE",
		1: "STRATEGY_REPLACE",
		2: "STRATEGY_ADD",
	}
	Condition_UpdateStrategy_value = map[string]int32{
		"STRATEGY_NONE":    0,
		"STRATEGY_REPLACE": 1,
		"STRATEGY_ADD":     2,
	}
)

func (x Condition_UpdateStrategy) Enum() *Condition_UpdateStrategy {
	p := new(Condition_UpdateStrategy)
	*p = x
	return p
}

func (x Condition_UpdateStrategy) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Condition_UpdateStrategy) Descriptor() protoreflect.EnumDescriptor {
	return file_condition_proto_enumTypes[1].Descriptor()
}

func (Condition_UpdateStrategy) Type() protoreflect.EnumType {
	return &file_condition_proto_enumTypes[1]
}

func (x Condition_UpdateStrategy) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Condition_UpdateStrategy.Descriptor instead.
func (Condition_UpdateStrategy) EnumDescriptor() ([]byte, []int) {
	return file_condition_proto_rawDescGZIP(), []int{4, 1}
}

type Nothing struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Nothing) Reset() {
	*x = Nothing{}
	if protoimpl.UnsafeEnabled {
		mi := &file_condition_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Nothing) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Nothing) ProtoMessage() {}

func (x *Nothing) ProtoReflect() protoreflect.Message {
	mi := &file_condition_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Nothing.ProtoReflect.Descriptor instead.
func (*Nothing) Descriptor() ([]byte, []int) {
	return file_condition_proto_rawDescGZIP(), []int{0}
}

type UpdateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Update []*Condition `protobuf:"bytes,1,rep,name=update,proto3" json:"update,omitempty"`
}

func (x *UpdateRequest) Reset() {
	*x = UpdateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_condition_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateRequest) ProtoMessage() {}

func (x *UpdateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_condition_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateRequest.ProtoReflect.Descriptor instead.
func (*UpdateRequest) Descriptor() ([]byte, []int) {
	return file_condition_proto_rawDescGZIP(), []int{1}
}

func (x *UpdateRequest) GetUpdate() []*Condition {
	if x != nil {
		return x.Update
	}
	return nil
}

type RegisterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Conditions []*Condition `protobuf:"bytes,1,rep,name=conditions,proto3" json:"conditions,omitempty"`
}

func (x *RegisterRequest) Reset() {
	*x = RegisterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_condition_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterRequest) ProtoMessage() {}

func (x *RegisterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_condition_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterRequest.ProtoReflect.Descriptor instead.
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return file_condition_proto_rawDescGZIP(), []int{2}
}

func (x *RegisterRequest) GetConditions() []*Condition {
	if x != nil {
		return x.Conditions
	}
	return nil
}

type UnRegisterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Conditions []*Condition `protobuf:"bytes,1,rep,name=conditions,proto3" json:"conditions,omitempty"`
}

func (x *UnRegisterRequest) Reset() {
	*x = UnRegisterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_condition_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnRegisterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnRegisterRequest) ProtoMessage() {}

func (x *UnRegisterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_condition_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnRegisterRequest.ProtoReflect.Descriptor instead.
func (*UnRegisterRequest) Descriptor() ([]byte, []int) {
	return file_condition_proto_rawDescGZIP(), []int{3}
}

func (x *UnRegisterRequest) GetConditions() []*Condition {
	if x != nil {
		return x.Conditions
	}
	return nil
}

type Condition struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OwnerId        int32                    `protobuf:"varint,1,opt,name=ownerId,proto3" json:"ownerId,omitempty"`
	Id             int32                    `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	Type           int32                    `protobuf:"varint,3,opt,name=type,proto3" json:"type,omitempty"`
	Params         []int32                  `protobuf:"varint,4,rep,packed,name=params,proto3" json:"params,omitempty"`
	Progress       int32                    `protobuf:"varint,5,opt,name=progress,proto3" json:"progress,omitempty"`
	Theme          string                   `protobuf:"bytes,6,opt,name=theme,proto3" json:"theme,omitempty"`
	Status         Condition_Status         `protobuf:"varint,7,opt,name=status,proto3,enum=condition.pb.Condition_Status" json:"status,omitempty"`
	UpdateStrategy Condition_UpdateStrategy `protobuf:"varint,8,opt,name=updateStrategy,proto3,enum=condition.pb.Condition_UpdateStrategy" json:"updateStrategy,omitempty"`
}

func (x *Condition) Reset() {
	*x = Condition{}
	if protoimpl.UnsafeEnabled {
		mi := &file_condition_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Condition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Condition) ProtoMessage() {}

func (x *Condition) ProtoReflect() protoreflect.Message {
	mi := &file_condition_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Condition.ProtoReflect.Descriptor instead.
func (*Condition) Descriptor() ([]byte, []int) {
	return file_condition_proto_rawDescGZIP(), []int{4}
}

func (x *Condition) GetOwnerId() int32 {
	if x != nil {
		return x.OwnerId
	}
	return 0
}

func (x *Condition) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Condition) GetType() int32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *Condition) GetParams() []int32 {
	if x != nil {
		return x.Params
	}
	return nil
}

func (x *Condition) GetProgress() int32 {
	if x != nil {
		return x.Progress
	}
	return 0
}

func (x *Condition) GetTheme() string {
	if x != nil {
		return x.Theme
	}
	return ""
}

func (x *Condition) GetStatus() Condition_Status {
	if x != nil {
		return x.Status
	}
	return Condition_NONE
}

func (x *Condition) GetUpdateStrategy() Condition_UpdateStrategy {
	if x != nil {
		return x.UpdateStrategy
	}
	return Condition_STRATEGY_NONE
}

type Changes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Conditions []*Condition `protobuf:"bytes,1,rep,name=conditions,proto3" json:"conditions,omitempty"`
}

func (x *Changes) Reset() {
	*x = Changes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_condition_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Changes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Changes) ProtoMessage() {}

func (x *Changes) ProtoReflect() protoreflect.Message {
	mi := &file_condition_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Changes.ProtoReflect.Descriptor instead.
func (*Changes) Descriptor() ([]byte, []int) {
	return file_condition_proto_rawDescGZIP(), []int{5}
}

func (x *Changes) GetConditions() []*Condition {
	if x != nil {
		return x.Conditions
	}
	return nil
}

var File_condition_proto protoreflect.FileDescriptor

var file_condition_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0c, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x62, 0x22,
	0x09, 0x0a, 0x07, 0x4e, 0x6f, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x22, 0x40, 0x0a, 0x0d, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2f, 0x0a, 0x06, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x63, 0x6f,
	0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6e, 0x64, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x22, 0x4a, 0x0a, 0x0f,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x37, 0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x63, 0x6f,
	0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x4c, 0x0a, 0x11, 0x55, 0x6e, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x37, 0x0a,
	0x0a, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x17, 0x2e, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x62,
	0x2e, 0x43, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x63, 0x6f, 0x6e, 0x64,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x96, 0x03, 0x0a, 0x09, 0x43, 0x6f, 0x6e, 0x64, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18, 0x04, 0x20, 0x03,
	0x28, 0x05, 0x52, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72,
	0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x72,
	0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x68, 0x65, 0x6d, 0x65, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x68, 0x65, 0x6d, 0x65, 0x12, 0x36, 0x0a, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x63,
	0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6e, 0x64,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x4e, 0x0a, 0x0e, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x74,
	0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x26, 0x2e, 0x63,
	0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6e, 0x64,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x74, 0x72, 0x61,
	0x74, 0x65, 0x67, 0x79, 0x52, 0x0e, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x74, 0x72, 0x61,
	0x74, 0x65, 0x67, 0x79, 0x22, 0x2c, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x08,
	0x0a, 0x04, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x41, 0x43, 0x54, 0x49,
	0x56, 0x45, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x46, 0x49, 0x4e, 0x49, 0x53, 0x48, 0x45, 0x44,
	0x10, 0x02, 0x22, 0x4b, 0x0a, 0x0e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x74, 0x72, 0x61,
	0x74, 0x65, 0x67, 0x79, 0x12, 0x11, 0x0a, 0x0d, 0x53, 0x54, 0x52, 0x41, 0x54, 0x45, 0x47, 0x59,
	0x5f, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12, 0x14, 0x0a, 0x10, 0x53, 0x54, 0x52, 0x41, 0x54,
	0x45, 0x47, 0x59, 0x5f, 0x52, 0x45, 0x50, 0x4c, 0x41, 0x43, 0x45, 0x10, 0x01, 0x12, 0x10, 0x0a,
	0x0c, 0x53, 0x54, 0x52, 0x41, 0x54, 0x45, 0x47, 0x59, 0x5f, 0x41, 0x44, 0x44, 0x10, 0x02, 0x22,
	0x42, 0x0a, 0x07, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x12, 0x37, 0x0a, 0x0a, 0x63, 0x6f,
	0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17,
	0x2e, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x6f,
	0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x32, 0x99, 0x02, 0x0a, 0x10, 0x43, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x39, 0x0a, 0x05, 0x57, 0x61, 0x74, 0x63,
	0x68, 0x12, 0x15, 0x2e, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x62,
	0x2e, 0x4e, 0x6f, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x1a, 0x15, 0x2e, 0x63, 0x6f, 0x6e, 0x64, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x22,
	0x00, 0x30, 0x01, 0x12, 0x42, 0x0a, 0x08, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x12,
	0x1d, 0x2e, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x62, 0x2e, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15,
	0x2e, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x62, 0x2e, 0x4e, 0x6f,
	0x74, 0x68, 0x69, 0x6e, 0x67, 0x22, 0x00, 0x12, 0x46, 0x0a, 0x0a, 0x55, 0x6e, 0x72, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x12, 0x1f, 0x2e, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x70, 0x62, 0x2e, 0x55, 0x6e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x70, 0x62, 0x2e, 0x4e, 0x6f, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x22, 0x00, 0x12,
	0x3e, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1b, 0x2e, 0x63, 0x6f, 0x6e, 0x64,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x62, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x70, 0x62, 0x2e, 0x4e, 0x6f, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x22, 0x00, 0x42,
	0x12, 0x5a, 0x10, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x61, 0x70, 0x69,
	0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_condition_proto_rawDescOnce sync.Once
	file_condition_proto_rawDescData = file_condition_proto_rawDesc
)

func file_condition_proto_rawDescGZIP() []byte {
	file_condition_proto_rawDescOnce.Do(func() {
		file_condition_proto_rawDescData = protoimpl.X.CompressGZIP(file_condition_proto_rawDescData)
	})
	return file_condition_proto_rawDescData
}

var file_condition_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_condition_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_condition_proto_goTypes = []interface{}{
	(Condition_Status)(0),         // 0: condition.pb.Condition.Status
	(Condition_UpdateStrategy)(0), // 1: condition.pb.Condition.UpdateStrategy
	(*Nothing)(nil),               // 2: condition.pb.Nothing
	(*UpdateRequest)(nil),         // 3: condition.pb.UpdateRequest
	(*RegisterRequest)(nil),       // 4: condition.pb.RegisterRequest
	(*UnRegisterRequest)(nil),     // 5: condition.pb.UnRegisterRequest
	(*Condition)(nil),             // 6: condition.pb.Condition
	(*Changes)(nil),               // 7: condition.pb.Changes
}
var file_condition_proto_depIdxs = []int32{
	6,  // 0: condition.pb.UpdateRequest.update:type_name -> condition.pb.Condition
	6,  // 1: condition.pb.RegisterRequest.conditions:type_name -> condition.pb.Condition
	6,  // 2: condition.pb.UnRegisterRequest.conditions:type_name -> condition.pb.Condition
	0,  // 3: condition.pb.Condition.status:type_name -> condition.pb.Condition.Status
	1,  // 4: condition.pb.Condition.updateStrategy:type_name -> condition.pb.Condition.UpdateStrategy
	6,  // 5: condition.pb.Changes.conditions:type_name -> condition.pb.Condition
	2,  // 6: condition.pb.ConditionService.Watch:input_type -> condition.pb.Nothing
	4,  // 7: condition.pb.ConditionService.Register:input_type -> condition.pb.RegisterRequest
	5,  // 8: condition.pb.ConditionService.Unregister:input_type -> condition.pb.UnRegisterRequest
	3,  // 9: condition.pb.ConditionService.Update:input_type -> condition.pb.UpdateRequest
	7,  // 10: condition.pb.ConditionService.Watch:output_type -> condition.pb.Changes
	2,  // 11: condition.pb.ConditionService.Register:output_type -> condition.pb.Nothing
	2,  // 12: condition.pb.ConditionService.Unregister:output_type -> condition.pb.Nothing
	2,  // 13: condition.pb.ConditionService.Update:output_type -> condition.pb.Nothing
	10, // [10:14] is the sub-list for method output_type
	6,  // [6:10] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_condition_proto_init() }
func file_condition_proto_init() {
	if File_condition_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_condition_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Nothing); i {
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
		file_condition_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateRequest); i {
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
		file_condition_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterRequest); i {
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
		file_condition_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnRegisterRequest); i {
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
		file_condition_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Condition); i {
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
		file_condition_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Changes); i {
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
			RawDescriptor: file_condition_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_condition_proto_goTypes,
		DependencyIndexes: file_condition_proto_depIdxs,
		EnumInfos:         file_condition_proto_enumTypes,
		MessageInfos:      file_condition_proto_msgTypes,
	}.Build()
	File_condition_proto = out.File
	file_condition_proto_rawDesc = nil
	file_condition_proto_goTypes = nil
	file_condition_proto_depIdxs = nil
}