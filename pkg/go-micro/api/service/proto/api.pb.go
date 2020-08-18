// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/service/proto/api.proto

package go_micro_api

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type Endpoint struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Host                 []string `protobuf:"bytes,2,rep,name=host,proto3" json:"host,omitempty"`
	Path                 []string `protobuf:"bytes,3,rep,name=path,proto3" json:"path,omitempty"`
	Method               []string `protobuf:"bytes,4,rep,name=method,proto3" json:"method,omitempty"`
	Stream               bool     `protobuf:"varint,5,opt,name=stream,proto3" json:"stream,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Endpoint) Reset()         { *m = Endpoint{} }
func (m *Endpoint) String() string { return proto.CompactTextString(m) }
func (*Endpoint) ProtoMessage()    {}
func (*Endpoint) Descriptor() ([]byte, []int) {
	return fileDescriptor_c4a48b6b680b5c31, []int{0}
}

func (m *Endpoint) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Endpoint.Unmarshal(m, b)
}
func (m *Endpoint) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Endpoint.Marshal(b, m, deterministic)
}
func (m *Endpoint) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Endpoint.Merge(m, src)
}
func (m *Endpoint) XXX_Size() int {
	return xxx_messageInfo_Endpoint.Size(m)
}
func (m *Endpoint) XXX_DiscardUnknown() {
	xxx_messageInfo_Endpoint.DiscardUnknown(m)
}

var xxx_messageInfo_Endpoint proto.InternalMessageInfo

func (m *Endpoint) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Endpoint) GetHost() []string {
	if m != nil {
		return m.Host
	}
	return nil
}

func (m *Endpoint) GetPath() []string {
	if m != nil {
		return m.Path
	}
	return nil
}

func (m *Endpoint) GetMethod() []string {
	if m != nil {
		return m.Method
	}
	return nil
}

func (m *Endpoint) GetStream() bool {
	if m != nil {
		return m.Stream
	}
	return false
}

type EmptyResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmptyResponse) Reset()         { *m = EmptyResponse{} }
func (m *EmptyResponse) String() string { return proto.CompactTextString(m) }
func (*EmptyResponse) ProtoMessage()    {}
func (*EmptyResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c4a48b6b680b5c31, []int{1}
}

func (m *EmptyResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmptyResponse.Unmarshal(m, b)
}
func (m *EmptyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmptyResponse.Marshal(b, m, deterministic)
}
func (m *EmptyResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmptyResponse.Merge(m, src)
}
func (m *EmptyResponse) XXX_Size() int {
	return xxx_messageInfo_EmptyResponse.Size(m)
}
func (m *EmptyResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_EmptyResponse.DiscardUnknown(m)
}

var xxx_messageInfo_EmptyResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Endpoint)(nil), "go.micro.api.Endpoint")
	proto.RegisterType((*EmptyResponse)(nil), "go.micro.api.EmptyResponse")
}

func init() { proto.RegisterFile("api/service/proto/api.proto", fileDescriptor_c4a48b6b680b5c31) }

var fileDescriptor_c4a48b6b680b5c31 = []byte{
	// 212 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0xd0, 0xc1, 0x4a, 0x03, 0x31,
	0x10, 0x80, 0x61, 0xd7, 0xad, 0x65, 0x1d, 0x14, 0x21, 0x87, 0x12, 0xec, 0x65, 0xd9, 0x53, 0x4f,
	0x59, 0xd0, 0x27, 0x28, 0xda, 0x17, 0xd8, 0x37, 0x88, 0xed, 0xd0, 0x9d, 0x43, 0x32, 0x43, 0x32,
	0x14, 0x7c, 0x08, 0xdf, 0x59, 0x12, 0x2b, 0x2c, 0x5e, 0xbd, 0xfd, 0xf3, 0x1d, 0x86, 0x61, 0x60,
	0xeb, 0x85, 0xc6, 0x8c, 0xe9, 0x42, 0x47, 0x1c, 0x25, 0xb1, 0xf2, 0xe8, 0x85, 0x5c, 0x2d, 0xf3,
	0x70, 0x66, 0x17, 0xe8, 0x98, 0xd8, 0x79, 0xa1, 0xe1, 0x02, 0xdd, 0x21, 0x9e, 0x84, 0x29, 0xaa,
	0x31, 0xb0, 0x8a, 0x3e, 0xa0, 0x6d, 0xfa, 0x66, 0x77, 0x3f, 0xd5, 0x2e, 0x36, 0x73, 0x56, 0x7b,
	0xdb, 0xb7, 0xc5, 0x4a, 0x17, 0x13, 0xaf, 0xb3, 0x6d, 0x7f, 0xac, 0xb4, 0xd9, 0xc0, 0x3a, 0xa0,
	0xce, 0x7c, 0xb2, 0xab, 0xaa, 0xd7, 0xa9, 0x78, 0xd6, 0x84, 0x3e, 0xd8, 0xbb, 0xbe, 0xd9, 0x75,
	0xd3, 0x75, 0x1a, 0x9e, 0xe0, 0xf1, 0x10, 0x44, 0x3f, 0x27, 0xcc, 0xc2, 0x31, 0xe3, 0xcb, 0x57,
	0x03, 0xed, 0x5e, 0xc8, 0xec, 0xa1, 0x9b, 0xf0, 0x4c, 0x59, 0x31, 0x99, 0x8d, 0x5b, 0xde, 0xea,
	0x7e, 0x0f, 0x7d, 0xde, 0xfe, 0xf1, 0xe5, 0xa2, 0xe1, 0xc6, 0xbc, 0x01, 0xbc, 0x63, 0xfa, 0xdf,
	0x92, 0x8f, 0x75, 0xfd, 0xd6, 0xeb, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0x46, 0x62, 0x67, 0x30,
	0x4c, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ApiClient is the client API for Api service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ApiClient interface {
	Register(ctx context.Context, in *Endpoint, opts ...grpc.CallOption) (*EmptyResponse, error)
	Deregister(ctx context.Context, in *Endpoint, opts ...grpc.CallOption) (*EmptyResponse, error)
}

type apiClient struct {
	cc *grpc.ClientConn
}

func NewApiClient(cc *grpc.ClientConn) ApiClient {
	return &apiClient{cc}
}

func (c *apiClient) Register(ctx context.Context, in *Endpoint, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/go.micro.api.Api/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiClient) Deregister(ctx context.Context, in *Endpoint, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/go.micro.api.Api/Deregister", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ApiServer is the server API for Api service.
type ApiServer interface {
	Register(context.Context, *Endpoint) (*EmptyResponse, error)
	Deregister(context.Context, *Endpoint) (*EmptyResponse, error)
}

// UnimplementedApiServer can be embedded to have forward compatible implementations.
type UnimplementedApiServer struct {
}

func (*UnimplementedApiServer) Register(ctx context.Context, req *Endpoint) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (*UnimplementedApiServer) Deregister(ctx context.Context, req *Endpoint) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deregister not implemented")
}

func RegisterApiServer(s *grpc.Server, srv ApiServer) {
	s.RegisterService(&_Api_serviceDesc, srv)
}

func _Api_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Endpoint)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/go.micro.api.Api/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).Register(ctx, req.(*Endpoint))
	}
	return interceptor(ctx, in, info, handler)
}

func _Api_Deregister_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Endpoint)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).Deregister(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/go.micro.api.Api/Deregister",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).Deregister(ctx, req.(*Endpoint))
	}
	return interceptor(ctx, in, info, handler)
}

var _Api_serviceDesc = grpc.ServiceDesc{
	ServiceName: "go.micro.api.Api",
	HandlerType: (*ApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _Api_Register_Handler,
		},
		{
			MethodName: "Deregister",
			Handler:    _Api_Deregister_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/service/proto/api.proto",
}
