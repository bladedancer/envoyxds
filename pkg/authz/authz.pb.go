// Code generated by protoc-gen-go. DO NOT EDIT.
// source: bladedancer/envoyxds/pkg/authz/authz.proto

package authz

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	protobuf "google/protobuf"
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

type ChangeType int32

const (
	ChangeType_UNKNOWN_CTX_TYPE ChangeType = 0
	ChangeType_BASIC            ChangeType = 1
	ChangeType_API              ChangeType = 2
	ChangeType_OAUTH            ChangeType = 3
)

var ChangeType_name = map[int32]string{
	0: "UNKNOWN_CTX_TYPE",
	1: "BASIC",
	2: "API",
	3: "OAUTH",
}

var ChangeType_value = map[string]int32{
	"UNKNOWN_CTX_TYPE": 0,
	"BASIC":            1,
	"API":              2,
	"OAUTH":            3,
}

func (x ChangeType) String() string {
	return proto.EnumName(ChangeType_name, int32(x))
}

func (ChangeType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_799bd7a37f25fa36, []int{0}
}

type AuthEnvelope struct {
	CtxType              ChangeType    `protobuf:"varint,1,opt,name=ctxType,proto3,enum=authz.ChangeType" json:"ctxType,omitempty"`
	Context              *protobuf.Any `protobuf:"bytes,2,opt,name=context,proto3" json:"context,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *AuthEnvelope) Reset()         { *m = AuthEnvelope{} }
func (m *AuthEnvelope) String() string { return proto.CompactTextString(m) }
func (*AuthEnvelope) ProtoMessage()    {}
func (*AuthEnvelope) Descriptor() ([]byte, []int) {
	return fileDescriptor_799bd7a37f25fa36, []int{0}
}

func (m *AuthEnvelope) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AuthEnvelope.Unmarshal(m, b)
}
func (m *AuthEnvelope) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AuthEnvelope.Marshal(b, m, deterministic)
}
func (m *AuthEnvelope) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthEnvelope.Merge(m, src)
}
func (m *AuthEnvelope) XXX_Size() int {
	return xxx_messageInfo_AuthEnvelope.Size(m)
}
func (m *AuthEnvelope) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthEnvelope.DiscardUnknown(m)
}

var xxx_messageInfo_AuthEnvelope proto.InternalMessageInfo

func (m *AuthEnvelope) GetCtxType() ChangeType {
	if m != nil {
		return m.CtxType
	}
	return ChangeType_UNKNOWN_CTX_TYPE
}

func (m *AuthEnvelope) GetContext() *protobuf.Any {
	if m != nil {
		return m.Context
	}
	return nil
}

type BasicAuthCtx struct {
	User                 string   `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	Pass                 string   `protobuf:"bytes,2,opt,name=pass,proto3" json:"pass,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BasicAuthCtx) Reset()         { *m = BasicAuthCtx{} }
func (m *BasicAuthCtx) String() string { return proto.CompactTextString(m) }
func (*BasicAuthCtx) ProtoMessage()    {}
func (*BasicAuthCtx) Descriptor() ([]byte, []int) {
	return fileDescriptor_799bd7a37f25fa36, []int{1}
}

func (m *BasicAuthCtx) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BasicAuthCtx.Unmarshal(m, b)
}
func (m *BasicAuthCtx) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BasicAuthCtx.Marshal(b, m, deterministic)
}
func (m *BasicAuthCtx) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BasicAuthCtx.Merge(m, src)
}
func (m *BasicAuthCtx) XXX_Size() int {
	return xxx_messageInfo_BasicAuthCtx.Size(m)
}
func (m *BasicAuthCtx) XXX_DiscardUnknown() {
	xxx_messageInfo_BasicAuthCtx.DiscardUnknown(m)
}

var xxx_messageInfo_BasicAuthCtx proto.InternalMessageInfo

func (m *BasicAuthCtx) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

func (m *BasicAuthCtx) GetPass() string {
	if m != nil {
		return m.Pass
	}
	return ""
}

type ApiKeyCtx struct {
	ApiKey               string   `protobuf:"bytes,1,opt,name=apiKey,proto3" json:"apiKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ApiKeyCtx) Reset()         { *m = ApiKeyCtx{} }
func (m *ApiKeyCtx) String() string { return proto.CompactTextString(m) }
func (*ApiKeyCtx) ProtoMessage()    {}
func (*ApiKeyCtx) Descriptor() ([]byte, []int) {
	return fileDescriptor_799bd7a37f25fa36, []int{2}
}

func (m *ApiKeyCtx) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ApiKeyCtx.Unmarshal(m, b)
}
func (m *ApiKeyCtx) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ApiKeyCtx.Marshal(b, m, deterministic)
}
func (m *ApiKeyCtx) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ApiKeyCtx.Merge(m, src)
}
func (m *ApiKeyCtx) XXX_Size() int {
	return xxx_messageInfo_ApiKeyCtx.Size(m)
}
func (m *ApiKeyCtx) XXX_DiscardUnknown() {
	xxx_messageInfo_ApiKeyCtx.DiscardUnknown(m)
}

var xxx_messageInfo_ApiKeyCtx proto.InternalMessageInfo

func (m *ApiKeyCtx) GetApiKey() string {
	if m != nil {
		return m.ApiKey
	}
	return ""
}

type OAuthCtx struct {
	Oath                 string   `protobuf:"bytes,1,opt,name=oath,proto3" json:"oath,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OAuthCtx) Reset()         { *m = OAuthCtx{} }
func (m *OAuthCtx) String() string { return proto.CompactTextString(m) }
func (*OAuthCtx) ProtoMessage()    {}
func (*OAuthCtx) Descriptor() ([]byte, []int) {
	return fileDescriptor_799bd7a37f25fa36, []int{3}
}

func (m *OAuthCtx) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OAuthCtx.Unmarshal(m, b)
}
func (m *OAuthCtx) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OAuthCtx.Marshal(b, m, deterministic)
}
func (m *OAuthCtx) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OAuthCtx.Merge(m, src)
}
func (m *OAuthCtx) XXX_Size() int {
	return xxx_messageInfo_OAuthCtx.Size(m)
}
func (m *OAuthCtx) XXX_DiscardUnknown() {
	xxx_messageInfo_OAuthCtx.DiscardUnknown(m)
}

var xxx_messageInfo_OAuthCtx proto.InternalMessageInfo

func (m *OAuthCtx) GetOath() string {
	if m != nil {
		return m.Oath
	}
	return ""
}

type ApiKeyMessage struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ApiKeyMessage) Reset()         { *m = ApiKeyMessage{} }
func (m *ApiKeyMessage) String() string { return proto.CompactTextString(m) }
func (*ApiKeyMessage) ProtoMessage()    {}
func (*ApiKeyMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_799bd7a37f25fa36, []int{4}
}

func (m *ApiKeyMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ApiKeyMessage.Unmarshal(m, b)
}
func (m *ApiKeyMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ApiKeyMessage.Marshal(b, m, deterministic)
}
func (m *ApiKeyMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ApiKeyMessage.Merge(m, src)
}
func (m *ApiKeyMessage) XXX_Size() int {
	return xxx_messageInfo_ApiKeyMessage.Size(m)
}
func (m *ApiKeyMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_ApiKeyMessage.DiscardUnknown(m)
}

var xxx_messageInfo_ApiKeyMessage proto.InternalMessageInfo

func (m *ApiKeyMessage) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func init() {
	proto.RegisterEnum("authz.ChangeType", ChangeType_name, ChangeType_value)
	proto.RegisterType((*AuthEnvelope)(nil), "authz.AuthEnvelope")
	proto.RegisterType((*BasicAuthCtx)(nil), "authz.BasicAuthCtx")
	proto.RegisterType((*ApiKeyCtx)(nil), "authz.ApiKeyCtx")
	proto.RegisterType((*OAuthCtx)(nil), "authz.OAuthCtx")
	proto.RegisterType((*ApiKeyMessage)(nil), "authz.ApiKeyMessage")
}

func init() {
	proto.RegisterFile("bladedancer/envoyxds/pkg/authz/authz.proto", fileDescriptor_799bd7a37f25fa36)
}

var fileDescriptor_799bd7a37f25fa36 = []byte{
	// 313 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x4c, 0x8f, 0xe1, 0x4b, 0xc2, 0x40,
	0x18, 0xc6, 0x9b, 0xa6, 0xb6, 0x37, 0x8b, 0x75, 0x48, 0x58, 0x1f, 0xc2, 0xd6, 0x17, 0x31, 0xd8,
	0xc0, 0xa0, 0xef, 0xe7, 0x10, 0x12, 0x49, 0x65, 0x4d, 0xaa, 0x4f, 0x72, 0xce, 0x6b, 0x13, 0xe5,
	0x6e, 0xec, 0x6e, 0xb2, 0xf5, 0xd7, 0xc7, 0xdd, 0x5c, 0xf5, 0xe5, 0x78, 0xde, 0xe7, 0x9e, 0xf7,
	0xf7, 0xf0, 0xc2, 0x60, 0xbd, 0x27, 0x1b, 0xba, 0x21, 0x2c, 0xa4, 0xa9, 0x4b, 0xd9, 0x81, 0x17,
	0xf9, 0x46, 0xb8, 0xc9, 0x2e, 0x72, 0x49, 0x26, 0xe3, 0xef, 0xf2, 0x75, 0x92, 0x94, 0x4b, 0x8e,
	0x1a, 0x7a, 0xb8, 0xbd, 0x89, 0x38, 0x8f, 0xf6, 0xd4, 0xd5, 0xe6, 0x3a, 0xfb, 0x72, 0x09, 0x2b,
	0xca, 0x84, 0xbd, 0x83, 0x36, 0xce, 0x64, 0x3c, 0x66, 0x07, 0xba, 0xe7, 0x09, 0x45, 0x8f, 0xd0,
	0x0a, 0x65, 0x1e, 0x14, 0x09, 0xed, 0x1a, 0x3d, 0xa3, 0x7f, 0x39, 0xbc, 0x72, 0x4a, 0xa0, 0x17,
	0x13, 0x16, 0x51, 0xf5, 0xe1, 0x57, 0x09, 0xe4, 0x40, 0x2b, 0xe4, 0x4c, 0xd2, 0x5c, 0x76, 0x6b,
	0x3d, 0xa3, 0x7f, 0x3e, 0xec, 0x38, 0x65, 0x93, 0x53, 0x35, 0x39, 0x98, 0x15, 0x7e, 0x15, 0xb2,
	0x9f, 0xa1, 0x3d, 0x22, 0x62, 0x1b, 0xaa, 0x46, 0x4f, 0xe6, 0x08, 0xc1, 0x69, 0x26, 0x68, 0xaa,
	0x9b, 0x4c, 0x5f, 0x6b, 0xe5, 0x25, 0x44, 0x08, 0x0d, 0x34, 0x7d, 0xad, 0xed, 0x07, 0x30, 0x71,
	0xb2, 0x9d, 0xd2, 0x42, 0x2d, 0x5d, 0x43, 0x93, 0xe8, 0xe1, 0xb8, 0x76, 0x9c, 0xec, 0x3b, 0x38,
	0x9b, 0xff, 0x03, 0x73, 0x22, 0xe3, 0x0a, 0xac, 0xb4, 0x7d, 0x0f, 0x17, 0x25, 0xe4, 0x95, 0x0a,
	0x41, 0x22, 0x8a, 0x2c, 0xa8, 0xef, 0x7e, 0x29, 0x4a, 0x0e, 0x30, 0xc0, 0xdf, 0x99, 0xa8, 0x03,
	0xd6, 0x72, 0x36, 0x9d, 0xcd, 0xdf, 0x67, 0x2b, 0x2f, 0xf8, 0x58, 0x05, 0x9f, 0x8b, 0xb1, 0x75,
	0x82, 0x4c, 0x68, 0x8c, 0xf0, 0xdb, 0xc4, 0xb3, 0x0c, 0xd4, 0x82, 0x3a, 0x5e, 0x4c, 0xac, 0x9a,
	0xf2, 0xe6, 0x78, 0x19, 0xbc, 0x58, 0xf5, 0x75, 0x53, 0x5f, 0xfe, 0xf4, 0x13, 0x00, 0x00, 0xff,
	0xff, 0x1a, 0xa0, 0x1d, 0x86, 0xa6, 0x01, 0x00, 0x00,
}
