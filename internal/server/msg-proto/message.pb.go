// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.1
// source: msg-proto/message.proto

package msg_proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MessageDelivery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message *Message `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	UserIds []string `protobuf:"bytes,2,rep,name=userIds,proto3" json:"userIds,omitempty"`
}

func (x *MessageDelivery) Reset() {
	*x = MessageDelivery{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageDelivery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageDelivery) ProtoMessage() {}

func (x *MessageDelivery) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageDelivery.ProtoReflect.Descriptor instead.
func (*MessageDelivery) Descriptor() ([]byte, []int) {
	return file_msg_proto_message_proto_rawDescGZIP(), []int{0}
}

func (x *MessageDelivery) GetMessage() *Message {
	if x != nil {
		return x.Message
	}
	return nil
}

func (x *MessageDelivery) GetUserIds() []string {
	if x != nil {
		return x.UserIds
	}
	return nil
}

type MessageStreamResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message *Message `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *MessageStreamResponse) Reset() {
	*x = MessageStreamResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageStreamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageStreamResponse) ProtoMessage() {}

func (x *MessageStreamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageStreamResponse.ProtoReflect.Descriptor instead.
func (*MessageStreamResponse) Descriptor() ([]byte, []int) {
	return file_msg_proto_message_proto_rawDescGZIP(), []int{1}
}

func (x *MessageStreamResponse) GetMessage() *Message {
	if x != nil {
		return x.Message
	}
	return nil
}

var File_msg_proto_message_proto protoreflect.FileDescriptor

var file_msg_proto_message_proto_rawDesc = []byte{
	0x0a, 0x17, 0x6d, 0x73, 0x67, 0x2d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x15, 0x6d, 0x73, 0x67, 0x2d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x55, 0x0a, 0x0f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x12, 0x28, 0x0a, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x6d, 0x6f, 0x64,
	0x65, 0x6c, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x73, 0x22, 0x41, 0x0a,
	0x15, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x32, 0x59, 0x0a, 0x0e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x47, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1e, 0x2e, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x42, 0x3a, 0x5a, 0x38, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x41, 0x72, 0x74, 0x79, 0x6f, 0x6d,
	0x41, 0x72, 0x74, 0x61, 0x6d, 0x6f, 0x6e, 0x6f, 0x76, 0x2f, 0x6d, 0x73, 0x67, 0x2f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x6d, 0x73,
	0x67, 0x2d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_msg_proto_message_proto_rawDescOnce sync.Once
	file_msg_proto_message_proto_rawDescData = file_msg_proto_message_proto_rawDesc
)

func file_msg_proto_message_proto_rawDescGZIP() []byte {
	file_msg_proto_message_proto_rawDescOnce.Do(func() {
		file_msg_proto_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_msg_proto_message_proto_rawDescData)
	})
	return file_msg_proto_message_proto_rawDescData
}

var file_msg_proto_message_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_msg_proto_message_proto_goTypes = []interface{}{
	(*MessageDelivery)(nil),       // 0: message.MessageDelivery
	(*MessageStreamResponse)(nil), // 1: message.MessageStreamResponse
	(*Message)(nil),               // 2: model.Message
	(*emptypb.Empty)(nil),         // 3: google.protobuf.Empty
}
var file_msg_proto_message_proto_depIdxs = []int32{
	2, // 0: message.MessageDelivery.message:type_name -> model.Message
	2, // 1: message.MessageStreamResponse.message:type_name -> model.Message
	3, // 2: message.MessageService.GetMessages:input_type -> google.protobuf.Empty
	1, // 3: message.MessageService.GetMessages:output_type -> message.MessageStreamResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_msg_proto_message_proto_init() }
func file_msg_proto_message_proto_init() {
	if File_msg_proto_message_proto != nil {
		return
	}
	file_msg_proto_model_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_msg_proto_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageDelivery); i {
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
		file_msg_proto_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageStreamResponse); i {
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
			RawDescriptor: file_msg_proto_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_msg_proto_message_proto_goTypes,
		DependencyIndexes: file_msg_proto_message_proto_depIdxs,
		MessageInfos:      file_msg_proto_message_proto_msgTypes,
	}.Build()
	File_msg_proto_message_proto = out.File
	file_msg_proto_message_proto_rawDesc = nil
	file_msg_proto_message_proto_goTypes = nil
	file_msg_proto_message_proto_depIdxs = nil
}
