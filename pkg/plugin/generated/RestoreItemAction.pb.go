// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.2
// source: RestoreItemAction.proto

package generated

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

type RestoreItemActionExecuteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Plugin         string `protobuf:"bytes,1,opt,name=plugin,proto3" json:"plugin,omitempty"`
	Item           []byte `protobuf:"bytes,2,opt,name=item,proto3" json:"item,omitempty"`
	Restore        []byte `protobuf:"bytes,3,opt,name=restore,proto3" json:"restore,omitempty"`
	ItemFromBackup []byte `protobuf:"bytes,4,opt,name=itemFromBackup,proto3" json:"itemFromBackup,omitempty"`
}

func (x *RestoreItemActionExecuteRequest) Reset() {
	*x = RestoreItemActionExecuteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_RestoreItemAction_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RestoreItemActionExecuteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RestoreItemActionExecuteRequest) ProtoMessage() {}

func (x *RestoreItemActionExecuteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_RestoreItemAction_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RestoreItemActionExecuteRequest.ProtoReflect.Descriptor instead.
func (*RestoreItemActionExecuteRequest) Descriptor() ([]byte, []int) {
	return file_RestoreItemAction_proto_rawDescGZIP(), []int{0}
}

func (x *RestoreItemActionExecuteRequest) GetPlugin() string {
	if x != nil {
		return x.Plugin
	}
	return ""
}

func (x *RestoreItemActionExecuteRequest) GetItem() []byte {
	if x != nil {
		return x.Item
	}
	return nil
}

func (x *RestoreItemActionExecuteRequest) GetRestore() []byte {
	if x != nil {
		return x.Restore
	}
	return nil
}

func (x *RestoreItemActionExecuteRequest) GetItemFromBackup() []byte {
	if x != nil {
		return x.ItemFromBackup
	}
	return nil
}

type RestoreItemActionExecuteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Item            []byte                `protobuf:"bytes,1,opt,name=item,proto3" json:"item,omitempty"`
	AdditionalItems []*ResourceIdentifier `protobuf:"bytes,2,rep,name=additionalItems,proto3" json:"additionalItems,omitempty"`
	SkipRestore     bool                  `protobuf:"varint,3,opt,name=skipRestore,proto3" json:"skipRestore,omitempty"`
}

func (x *RestoreItemActionExecuteResponse) Reset() {
	*x = RestoreItemActionExecuteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_RestoreItemAction_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RestoreItemActionExecuteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RestoreItemActionExecuteResponse) ProtoMessage() {}

func (x *RestoreItemActionExecuteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_RestoreItemAction_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RestoreItemActionExecuteResponse.ProtoReflect.Descriptor instead.
func (*RestoreItemActionExecuteResponse) Descriptor() ([]byte, []int) {
	return file_RestoreItemAction_proto_rawDescGZIP(), []int{1}
}

func (x *RestoreItemActionExecuteResponse) GetItem() []byte {
	if x != nil {
		return x.Item
	}
	return nil
}

func (x *RestoreItemActionExecuteResponse) GetAdditionalItems() []*ResourceIdentifier {
	if x != nil {
		return x.AdditionalItems
	}
	return nil
}

func (x *RestoreItemActionExecuteResponse) GetSkipRestore() bool {
	if x != nil {
		return x.SkipRestore
	}
	return false
}

type RestoreItemActionAppliesToRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Plugin string `protobuf:"bytes,1,opt,name=plugin,proto3" json:"plugin,omitempty"`
}

func (x *RestoreItemActionAppliesToRequest) Reset() {
	*x = RestoreItemActionAppliesToRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_RestoreItemAction_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RestoreItemActionAppliesToRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RestoreItemActionAppliesToRequest) ProtoMessage() {}

func (x *RestoreItemActionAppliesToRequest) ProtoReflect() protoreflect.Message {
	mi := &file_RestoreItemAction_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RestoreItemActionAppliesToRequest.ProtoReflect.Descriptor instead.
func (*RestoreItemActionAppliesToRequest) Descriptor() ([]byte, []int) {
	return file_RestoreItemAction_proto_rawDescGZIP(), []int{2}
}

func (x *RestoreItemActionAppliesToRequest) GetPlugin() string {
	if x != nil {
		return x.Plugin
	}
	return ""
}

type RestoreItemActionAppliesToResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceSelector *ResourceSelector `protobuf:"bytes,1,opt,name=ResourceSelector,proto3" json:"ResourceSelector,omitempty"`
}

func (x *RestoreItemActionAppliesToResponse) Reset() {
	*x = RestoreItemActionAppliesToResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_RestoreItemAction_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RestoreItemActionAppliesToResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RestoreItemActionAppliesToResponse) ProtoMessage() {}

func (x *RestoreItemActionAppliesToResponse) ProtoReflect() protoreflect.Message {
	mi := &file_RestoreItemAction_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RestoreItemActionAppliesToResponse.ProtoReflect.Descriptor instead.
func (*RestoreItemActionAppliesToResponse) Descriptor() ([]byte, []int) {
	return file_RestoreItemAction_proto_rawDescGZIP(), []int{3}
}

func (x *RestoreItemActionAppliesToResponse) GetResourceSelector() *ResourceSelector {
	if x != nil {
		return x.ResourceSelector
	}
	return nil
}

var File_RestoreItemAction_proto protoreflect.FileDescriptor

var file_RestoreItemAction_proto_rawDesc = []byte{
	0x0a, 0x17, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x41, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x67, 0x65, 0x6e, 0x65, 0x72,
	0x61, 0x74, 0x65, 0x64, 0x1a, 0x0c, 0x53, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x8f, 0x01, 0x0a, 0x1f, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x49, 0x74,
	0x65, 0x6d, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x12, 0x12,
	0x0a, 0x04, 0x69, 0x74, 0x65, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x69, 0x74,
	0x65, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x07, 0x72, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x12, 0x26, 0x0a, 0x0e,
	0x69, 0x74, 0x65, 0x6d, 0x46, 0x72, 0x6f, 0x6d, 0x42, 0x61, 0x63, 0x6b, 0x75, 0x70, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x0e, 0x69, 0x74, 0x65, 0x6d, 0x46, 0x72, 0x6f, 0x6d, 0x42, 0x61,
	0x63, 0x6b, 0x75, 0x70, 0x22, 0xa1, 0x01, 0x0a, 0x20, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x49, 0x74, 0x65, 0x6d, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x69, 0x74, 0x65,
	0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x69, 0x74, 0x65, 0x6d, 0x12, 0x47, 0x0a,
	0x0f, 0x61, 0x64, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x49, 0x74, 0x65, 0x6d, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74,
	0x65, 0x64, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x66, 0x69, 0x65, 0x72, 0x52, 0x0f, 0x61, 0x64, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x61,
	0x6c, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x73, 0x6b, 0x69, 0x70, 0x52, 0x65,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x73, 0x6b, 0x69,
	0x70, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x22, 0x3b, 0x0a, 0x21, 0x52, 0x65, 0x73, 0x74,
	0x6f, 0x72, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x70, 0x70,
	0x6c, 0x69, 0x65, 0x73, 0x54, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a,
	0x06, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x22, 0x6d, 0x0a, 0x22, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x49, 0x74, 0x65, 0x6d, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x65,
	0x73, 0x54, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x47, 0x0a, 0x10, 0x52,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65,
	0x64, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x52, 0x10, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x65, 0x6c, 0x65,
	0x63, 0x74, 0x6f, 0x72, 0x32, 0xe1, 0x01, 0x0a, 0x11, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x49, 0x74, 0x65, 0x6d, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x68, 0x0a, 0x09, 0x41, 0x70,
	0x70, 0x6c, 0x69, 0x65, 0x73, 0x54, 0x6f, 0x12, 0x2c, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61,
	0x74, 0x65, 0x64, 0x2e, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x41,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x65, 0x73, 0x54, 0x6f, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2d, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65,
	0x64, 0x2e, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x41, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x65, 0x73, 0x54, 0x6f, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x62, 0x0a, 0x07, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x12,
	0x2a, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x52, 0x65, 0x73, 0x74,
	0x6f, 0x72, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x78, 0x65,
	0x63, 0x75, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2b, 0x2e, 0x67, 0x65,
	0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x49,
	0x74, 0x65, 0x6d, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x76, 0x6d, 0x77, 0x61, 0x72, 0x65, 0x2d, 0x74, 0x61,
	0x6e, 0x7a, 0x75, 0x2f, 0x76, 0x65, 0x6c, 0x65, 0x72, 0x6f, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_RestoreItemAction_proto_rawDescOnce sync.Once
	file_RestoreItemAction_proto_rawDescData = file_RestoreItemAction_proto_rawDesc
)

func file_RestoreItemAction_proto_rawDescGZIP() []byte {
	file_RestoreItemAction_proto_rawDescOnce.Do(func() {
		file_RestoreItemAction_proto_rawDescData = protoimpl.X.CompressGZIP(file_RestoreItemAction_proto_rawDescData)
	})
	return file_RestoreItemAction_proto_rawDescData
}

var file_RestoreItemAction_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_RestoreItemAction_proto_goTypes = []interface{}{
	(*RestoreItemActionExecuteRequest)(nil),    // 0: generated.RestoreItemActionExecuteRequest
	(*RestoreItemActionExecuteResponse)(nil),   // 1: generated.RestoreItemActionExecuteResponse
	(*RestoreItemActionAppliesToRequest)(nil),  // 2: generated.RestoreItemActionAppliesToRequest
	(*RestoreItemActionAppliesToResponse)(nil), // 3: generated.RestoreItemActionAppliesToResponse
	(*ResourceIdentifier)(nil),                 // 4: generated.ResourceIdentifier
	(*ResourceSelector)(nil),                   // 5: generated.ResourceSelector
}
var file_RestoreItemAction_proto_depIdxs = []int32{
	4, // 0: generated.RestoreItemActionExecuteResponse.additionalItems:type_name -> generated.ResourceIdentifier
	5, // 1: generated.RestoreItemActionAppliesToResponse.ResourceSelector:type_name -> generated.ResourceSelector
	2, // 2: generated.RestoreItemAction.AppliesTo:input_type -> generated.RestoreItemActionAppliesToRequest
	0, // 3: generated.RestoreItemAction.Execute:input_type -> generated.RestoreItemActionExecuteRequest
	3, // 4: generated.RestoreItemAction.AppliesTo:output_type -> generated.RestoreItemActionAppliesToResponse
	1, // 5: generated.RestoreItemAction.Execute:output_type -> generated.RestoreItemActionExecuteResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_RestoreItemAction_proto_init() }
func file_RestoreItemAction_proto_init() {
	if File_RestoreItemAction_proto != nil {
		return
	}
	file_Shared_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_RestoreItemAction_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RestoreItemActionExecuteRequest); i {
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
		file_RestoreItemAction_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RestoreItemActionExecuteResponse); i {
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
		file_RestoreItemAction_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RestoreItemActionAppliesToRequest); i {
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
		file_RestoreItemAction_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RestoreItemActionAppliesToResponse); i {
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
			RawDescriptor: file_RestoreItemAction_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_RestoreItemAction_proto_goTypes,
		DependencyIndexes: file_RestoreItemAction_proto_depIdxs,
		MessageInfos:      file_RestoreItemAction_proto_msgTypes,
	}.Build()
	File_RestoreItemAction_proto = out.File
	file_RestoreItemAction_proto_rawDesc = nil
	file_RestoreItemAction_proto_goTypes = nil
	file_RestoreItemAction_proto_depIdxs = nil
}
