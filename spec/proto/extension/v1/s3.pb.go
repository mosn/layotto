// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.19.1
// source: s3.proto

package s3

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/anypb"
	_ "google.golang.org/protobuf/types/known/emptypb"
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

type GetObjectInput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Bucket       string `protobuf:"bytes,1,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Key          string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	ChecksumMode string `protobuf:"bytes,3,opt,name=checksum_mode,json=checksumMode,proto3" json:"checksum_mode,omitempty"`
	PartNumber   string `protobuf:"bytes,4,opt,name=part_number,json=partNumber,proto3" json:"part_number,omitempty"`
}

func (x *GetObjectInput) Reset() {
	*x = GetObjectInput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s3_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetObjectInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetObjectInput) ProtoMessage() {}

func (x *GetObjectInput) ProtoReflect() protoreflect.Message {
	mi := &file_s3_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetObjectInput.ProtoReflect.Descriptor instead.
func (*GetObjectInput) Descriptor() ([]byte, []int) {
	return file_s3_proto_rawDescGZIP(), []int{0}
}

func (x *GetObjectInput) GetBucket() string {
	if x != nil {
		return x.Bucket
	}
	return ""
}

func (x *GetObjectInput) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *GetObjectInput) GetChecksumMode() string {
	if x != nil {
		return x.ChecksumMode
	}
	return ""
}

func (x *GetObjectInput) GetPartNumber() string {
	if x != nil {
		return x.PartNumber
	}
	return ""
}

type GetObjectOutput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AcceptRanges string `protobuf:"bytes,1,opt,name=accept_ranges,json=acceptRanges,proto3" json:"accept_ranges,omitempty"`
	Body         []byte `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *GetObjectOutput) Reset() {
	*x = GetObjectOutput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s3_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetObjectOutput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetObjectOutput) ProtoMessage() {}

func (x *GetObjectOutput) ProtoReflect() protoreflect.Message {
	mi := &file_s3_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetObjectOutput.ProtoReflect.Descriptor instead.
func (*GetObjectOutput) Descriptor() ([]byte, []int) {
	return file_s3_proto_rawDescGZIP(), []int{1}
}

func (x *GetObjectOutput) GetAcceptRanges() string {
	if x != nil {
		return x.AcceptRanges
	}
	return ""
}

func (x *GetObjectOutput) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

var File_s3_proto protoreflect.FileDescriptor

var file_s3_proto_rawDesc = []byte{
	0x0a, 0x08, 0x73, 0x33, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x6d, 0x6f, 0x73, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x76,
	0x31, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x80, 0x01, 0x0a, 0x0e, 0x47, 0x65,
	0x74, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x16, 0x0a, 0x06,
	0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x62, 0x75,
	0x63, 0x6b, 0x65, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x73,
	0x75, 0x6d, 0x5f, 0x6d, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x63,
	0x68, 0x65, 0x63, 0x6b, 0x73, 0x75, 0x6d, 0x4d, 0x6f, 0x64, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x70,
	0x61, 0x72, 0x74, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x70, 0x61, 0x72, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x4a, 0x0a, 0x0f,
	0x47, 0x65, 0x74, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12,
	0x23, 0x0a, 0x0d, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x5f, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x52, 0x61,
	0x6e, 0x67, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x32, 0x6e, 0x0a, 0x0c, 0x53, 0x33, 0x4f, 0x75,
	0x74, 0x42, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x12, 0x5e, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x4f,
	0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x25, 0x2e, 0x6d, 0x6f, 0x73, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65,
	0x74, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x1a, 0x26, 0x2e, 0x6d,
	0x6f, 0x73, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x4f, 0x75,
	0x74, 0x70, 0x75, 0x74, 0x22, 0x00, 0x30, 0x01, 0x42, 0x3e, 0x5a, 0x3c, 0x67, 0x69, 0x74, 0x6c,
	0x61, 0x62, 0x2e, 0x61, 0x6c, 0x69, 0x70, 0x61, 0x79, 0x2d, 0x69, 0x6e, 0x63, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x61, 0x6e, 0x74, 0x2d, 0x6d, 0x65, 0x73, 0x68, 0x2f, 0x6c, 0x61, 0x79, 0x6f, 0x74,
	0x74, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69,
	0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x3b, 0x73, 0x33, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_s3_proto_rawDescOnce sync.Once
	file_s3_proto_rawDescData = file_s3_proto_rawDesc
)

func file_s3_proto_rawDescGZIP() []byte {
	file_s3_proto_rawDescOnce.Do(func() {
		file_s3_proto_rawDescData = protoimpl.X.CompressGZIP(file_s3_proto_rawDescData)
	})
	return file_s3_proto_rawDescData
}

var file_s3_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_s3_proto_goTypes = []interface{}{
	(*GetObjectInput)(nil),  // 0: mosn.proto.runtime.v1.GetObjectInput
	(*GetObjectOutput)(nil), // 1: mosn.proto.runtime.v1.GetObjectOutput
}
var file_s3_proto_depIdxs = []int32{
	0, // 0: mosn.proto.runtime.v1.S3OutBinding.GetObject:input_type -> mosn.proto.runtime.v1.GetObjectInput
	1, // 1: mosn.proto.runtime.v1.S3OutBinding.GetObject:output_type -> mosn.proto.runtime.v1.GetObjectOutput
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_s3_proto_init() }
func file_s3_proto_init() {
	if File_s3_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_s3_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetObjectInput); i {
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
		file_s3_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetObjectOutput); i {
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
			RawDescriptor: file_s3_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_s3_proto_goTypes,
		DependencyIndexes: file_s3_proto_depIdxs,
		MessageInfos:      file_s3_proto_msgTypes,
	}.Build()
	File_s3_proto = out.File
	file_s3_proto_rawDesc = nil
	file_s3_proto_goTypes = nil
	file_s3_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// S3OutBindingClient is the client API for S3OutBinding service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type S3OutBindingClient interface {
	GetObject(ctx context.Context, in *GetObjectInput, opts ...grpc.CallOption) (S3OutBinding_GetObjectClient, error)
}

type s3OutBindingClient struct {
	cc grpc.ClientConnInterface
}

func NewS3OutBindingClient(cc grpc.ClientConnInterface) S3OutBindingClient {
	return &s3OutBindingClient{cc}
}

func (c *s3OutBindingClient) GetObject(ctx context.Context, in *GetObjectInput, opts ...grpc.CallOption) (S3OutBinding_GetObjectClient, error) {
	stream, err := c.cc.NewStream(ctx, &_S3OutBinding_serviceDesc.Streams[0], "/mosn.proto.runtime.v1.S3OutBinding/GetObject", opts...)
	if err != nil {
		return nil, err
	}
	x := &s3OutBindingGetObjectClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type S3OutBinding_GetObjectClient interface {
	Recv() (*GetObjectOutput, error)
	grpc.ClientStream
}

type s3OutBindingGetObjectClient struct {
	grpc.ClientStream
}

func (x *s3OutBindingGetObjectClient) Recv() (*GetObjectOutput, error) {
	m := new(GetObjectOutput)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// S3OutBindingServer is the server API for S3OutBinding service.
type S3OutBindingServer interface {
	GetObject(*GetObjectInput, S3OutBinding_GetObjectServer) error
}

// UnimplementedS3OutBindingServer can be embedded to have forward compatible implementations.
type UnimplementedS3OutBindingServer struct {
}

func (*UnimplementedS3OutBindingServer) GetObject(*GetObjectInput, S3OutBinding_GetObjectServer) error {
	return status.Errorf(codes.Unimplemented, "method GetObject not implemented")
}

func RegisterS3OutBindingServer(s *grpc.Server, srv S3OutBindingServer) {
	s.RegisterService(&_S3OutBinding_serviceDesc, srv)
}

func _S3OutBinding_GetObject_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetObjectInput)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(S3OutBindingServer).GetObject(m, &s3OutBindingGetObjectServer{stream})
}

type S3OutBinding_GetObjectServer interface {
	Send(*GetObjectOutput) error
	grpc.ServerStream
}

type s3OutBindingGetObjectServer struct {
	grpc.ServerStream
}

func (x *s3OutBindingGetObjectServer) Send(m *GetObjectOutput) error {
	return x.ServerStream.SendMsg(m)
}

var _S3OutBinding_serviceDesc = grpc.ServiceDesc{
	ServiceName: "mosn.proto.runtime.v1.S3OutBinding",
	HandlerType: (*S3OutBindingServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetObject",
			Handler:       _S3OutBinding_GetObject_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "s3.proto",
}
