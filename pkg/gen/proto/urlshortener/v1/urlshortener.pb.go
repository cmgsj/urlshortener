// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3-devel
// 	protoc        (unknown)
// source: urlshortener/v1/urlshortener.proto

package urlshortenerv1

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type URL struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UrlId         string                 `protobuf:"bytes,1,opt,name=url_id,json=urlId,proto3" json:"url_id,omitempty"`
	RedirectUrl   string                 `protobuf:"bytes,2,opt,name=redirect_url,json=redirectUrl,proto3" json:"redirect_url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *URL) Reset() {
	*x = URL{}
	mi := &file_urlshortener_v1_urlshortener_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *URL) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*URL) ProtoMessage() {}

func (x *URL) ProtoReflect() protoreflect.Message {
	mi := &file_urlshortener_v1_urlshortener_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use URL.ProtoReflect.Descriptor instead.
func (*URL) Descriptor() ([]byte, []int) {
	return file_urlshortener_v1_urlshortener_proto_rawDescGZIP(), []int{0}
}

func (x *URL) GetUrlId() string {
	if x != nil {
		return x.UrlId
	}
	return ""
}

func (x *URL) GetRedirectUrl() string {
	if x != nil {
		return x.RedirectUrl
	}
	return ""
}

type ListURLsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListURLsRequest) Reset() {
	*x = ListURLsRequest{}
	mi := &file_urlshortener_v1_urlshortener_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListURLsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListURLsRequest) ProtoMessage() {}

func (x *ListURLsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_urlshortener_v1_urlshortener_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListURLsRequest.ProtoReflect.Descriptor instead.
func (*ListURLsRequest) Descriptor() ([]byte, []int) {
	return file_urlshortener_v1_urlshortener_proto_rawDescGZIP(), []int{1}
}

type ListURLsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Urls          []*URL                 `protobuf:"bytes,1,rep,name=urls,proto3" json:"urls,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListURLsResponse) Reset() {
	*x = ListURLsResponse{}
	mi := &file_urlshortener_v1_urlshortener_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListURLsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListURLsResponse) ProtoMessage() {}

func (x *ListURLsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_urlshortener_v1_urlshortener_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListURLsResponse.ProtoReflect.Descriptor instead.
func (*ListURLsResponse) Descriptor() ([]byte, []int) {
	return file_urlshortener_v1_urlshortener_proto_rawDescGZIP(), []int{2}
}

func (x *ListURLsResponse) GetUrls() []*URL {
	if x != nil {
		return x.Urls
	}
	return nil
}

type GetURLRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UrlId         string                 `protobuf:"bytes,1,opt,name=url_id,json=urlId,proto3" json:"url_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetURLRequest) Reset() {
	*x = GetURLRequest{}
	mi := &file_urlshortener_v1_urlshortener_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetURLRequest) ProtoMessage() {}

func (x *GetURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_urlshortener_v1_urlshortener_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetURLRequest.ProtoReflect.Descriptor instead.
func (*GetURLRequest) Descriptor() ([]byte, []int) {
	return file_urlshortener_v1_urlshortener_proto_rawDescGZIP(), []int{3}
}

func (x *GetURLRequest) GetUrlId() string {
	if x != nil {
		return x.UrlId
	}
	return ""
}

type GetURLResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Url           *URL                   `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetURLResponse) Reset() {
	*x = GetURLResponse{}
	mi := &file_urlshortener_v1_urlshortener_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetURLResponse) ProtoMessage() {}

func (x *GetURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_urlshortener_v1_urlshortener_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetURLResponse.ProtoReflect.Descriptor instead.
func (*GetURLResponse) Descriptor() ([]byte, []int) {
	return file_urlshortener_v1_urlshortener_proto_rawDescGZIP(), []int{4}
}

func (x *GetURLResponse) GetUrl() *URL {
	if x != nil {
		return x.Url
	}
	return nil
}

type CreateURLRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RedirectUrl   string                 `protobuf:"bytes,1,opt,name=redirect_url,json=redirectUrl,proto3" json:"redirect_url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateURLRequest) Reset() {
	*x = CreateURLRequest{}
	mi := &file_urlshortener_v1_urlshortener_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateURLRequest) ProtoMessage() {}

func (x *CreateURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_urlshortener_v1_urlshortener_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateURLRequest.ProtoReflect.Descriptor instead.
func (*CreateURLRequest) Descriptor() ([]byte, []int) {
	return file_urlshortener_v1_urlshortener_proto_rawDescGZIP(), []int{5}
}

func (x *CreateURLRequest) GetRedirectUrl() string {
	if x != nil {
		return x.RedirectUrl
	}
	return ""
}

type CreateURLResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UrlId         string                 `protobuf:"bytes,1,opt,name=url_id,json=urlId,proto3" json:"url_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateURLResponse) Reset() {
	*x = CreateURLResponse{}
	mi := &file_urlshortener_v1_urlshortener_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateURLResponse) ProtoMessage() {}

func (x *CreateURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_urlshortener_v1_urlshortener_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateURLResponse.ProtoReflect.Descriptor instead.
func (*CreateURLResponse) Descriptor() ([]byte, []int) {
	return file_urlshortener_v1_urlshortener_proto_rawDescGZIP(), []int{6}
}

func (x *CreateURLResponse) GetUrlId() string {
	if x != nil {
		return x.UrlId
	}
	return ""
}

type UpdateURLRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Url           *URL                   `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateURLRequest) Reset() {
	*x = UpdateURLRequest{}
	mi := &file_urlshortener_v1_urlshortener_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateURLRequest) ProtoMessage() {}

func (x *UpdateURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_urlshortener_v1_urlshortener_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateURLRequest.ProtoReflect.Descriptor instead.
func (*UpdateURLRequest) Descriptor() ([]byte, []int) {
	return file_urlshortener_v1_urlshortener_proto_rawDescGZIP(), []int{7}
}

func (x *UpdateURLRequest) GetUrl() *URL {
	if x != nil {
		return x.Url
	}
	return nil
}

type UpdateURLResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateURLResponse) Reset() {
	*x = UpdateURLResponse{}
	mi := &file_urlshortener_v1_urlshortener_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateURLResponse) ProtoMessage() {}

func (x *UpdateURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_urlshortener_v1_urlshortener_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateURLResponse.ProtoReflect.Descriptor instead.
func (*UpdateURLResponse) Descriptor() ([]byte, []int) {
	return file_urlshortener_v1_urlshortener_proto_rawDescGZIP(), []int{8}
}

type DeleteURLRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UrlId         string                 `protobuf:"bytes,1,opt,name=url_id,json=urlId,proto3" json:"url_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteURLRequest) Reset() {
	*x = DeleteURLRequest{}
	mi := &file_urlshortener_v1_urlshortener_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteURLRequest) ProtoMessage() {}

func (x *DeleteURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_urlshortener_v1_urlshortener_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteURLRequest.ProtoReflect.Descriptor instead.
func (*DeleteURLRequest) Descriptor() ([]byte, []int) {
	return file_urlshortener_v1_urlshortener_proto_rawDescGZIP(), []int{9}
}

func (x *DeleteURLRequest) GetUrlId() string {
	if x != nil {
		return x.UrlId
	}
	return ""
}

type DeleteURLResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteURLResponse) Reset() {
	*x = DeleteURLResponse{}
	mi := &file_urlshortener_v1_urlshortener_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteURLResponse) ProtoMessage() {}

func (x *DeleteURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_urlshortener_v1_urlshortener_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteURLResponse.ProtoReflect.Descriptor instead.
func (*DeleteURLResponse) Descriptor() ([]byte, []int) {
	return file_urlshortener_v1_urlshortener_proto_rawDescGZIP(), []int{10}
}

var File_urlshortener_v1_urlshortener_proto protoreflect.FileDescriptor

var file_urlshortener_v1_urlshortener_proto_rawDesc = string([]byte{
	0x0a, 0x22, 0x75, 0x72, 0x6c, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2f, 0x76,
	0x31, 0x2f, 0x75, 0x72, 0x6c, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x75, 0x72, 0x6c, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e,
	0x65, 0x72, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d,
	0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x3f, 0x0a, 0x03, 0x55, 0x52, 0x4c, 0x12, 0x15, 0x0a, 0x06, 0x75, 0x72,
	0x6c, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x75, 0x72, 0x6c, 0x49,
	0x64, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x5f, 0x75, 0x72,
	0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63,
	0x74, 0x55, 0x72, 0x6c, 0x22, 0x11, 0x0a, 0x0f, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x52, 0x4c, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x3c, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x55,
	0x52, 0x4c, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a, 0x04, 0x75,
	0x72, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x75, 0x72, 0x6c, 0x73,
	0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x52, 0x4c, 0x52,
	0x04, 0x75, 0x72, 0x6c, 0x73, 0x22, 0x26, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x15, 0x0a, 0x06, 0x75, 0x72, 0x6c, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x75, 0x72, 0x6c, 0x49, 0x64, 0x22, 0x38, 0x0a,
	0x0e, 0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x26, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x75,
	0x72, 0x6c, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55,
	0x52, 0x4c, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x35, 0x0a, 0x10, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x72,
	0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x55, 0x72, 0x6c, 0x22, 0x2a,
	0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x15, 0x0a, 0x06, 0x75, 0x72, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x75, 0x72, 0x6c, 0x49, 0x64, 0x22, 0x3a, 0x0a, 0x10, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x26,
	0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x75, 0x72,
	0x6c, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x52,
	0x4c, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x13, 0x0a, 0x11, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x29, 0x0a, 0x10, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x15, 0x0a, 0x06, 0x75, 0x72, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x75, 0x72, 0x6c, 0x49, 0x64, 0x22, 0x13, 0x0a, 0x11, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x82, 0x07, 0x0a, 0x13,
	0x55, 0x52, 0x4c, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0xa4, 0x01, 0x0a, 0x08, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x52, 0x4c, 0x73,
	0x12, 0x20, 0x2e, 0x75, 0x72, 0x6c, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e,
	0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x52, 0x4c, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x21, 0x2e, 0x75, 0x72, 0x6c, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65,
	0x72, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x52, 0x4c, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x53, 0x92, 0x41, 0x3d, 0x0a, 0x13, 0x55, 0x52, 0x4c, 0x53,
	0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x09, 0x4c, 0x69, 0x73, 0x74, 0x20, 0x55, 0x52, 0x4c, 0x73, 0x2a, 0x09, 0x6c, 0x69, 0x73, 0x74,
	0x5f, 0x75, 0x72, 0x6c, 0x73, 0x62, 0x10, 0x0a, 0x0e, 0x0a, 0x0a, 0x41, 0x50, 0x49, 0x4b, 0x65,
	0x79, 0x41, 0x75, 0x74, 0x68, 0x12, 0x00, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0d, 0x12, 0x0b, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x72, 0x6c, 0x12, 0xa3, 0x01, 0x0a, 0x06, 0x47,
	0x65, 0x74, 0x55, 0x52, 0x4c, 0x12, 0x1e, 0x2e, 0x75, 0x72, 0x6c, 0x73, 0x68, 0x6f, 0x72, 0x74,
	0x65, 0x6e, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x75, 0x72, 0x6c, 0x73, 0x68, 0x6f, 0x72, 0x74,
	0x65, 0x6e, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x58, 0x92, 0x41, 0x39, 0x0a, 0x13, 0x55, 0x52, 0x4c,
	0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x07, 0x47, 0x65, 0x74, 0x20, 0x55, 0x52, 0x4c, 0x2a, 0x07, 0x67, 0x65, 0x74, 0x5f, 0x75,
	0x72, 0x6c, 0x62, 0x10, 0x0a, 0x0e, 0x0a, 0x0a, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x41, 0x75,
	0x74, 0x68, 0x12, 0x00, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x16, 0x12, 0x14, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x76, 0x31, 0x2f, 0x75, 0x72, 0x6c, 0x2f, 0x7b, 0x75, 0x72, 0x6c, 0x5f, 0x69, 0x64, 0x7d,
	0x12, 0xac, 0x01, 0x0a, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x52, 0x4c, 0x12, 0x21,
	0x2e, 0x75, 0x72, 0x6c, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e, 0x76, 0x31,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x22, 0x2e, 0x75, 0x72, 0x6c, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72,
	0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x58, 0x92, 0x41, 0x3f, 0x0a, 0x13, 0x55, 0x52, 0x4c, 0x53,
	0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x20, 0x55, 0x52, 0x4c, 0x2a, 0x0a, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x5f, 0x75, 0x72, 0x6c, 0x62, 0x10, 0x0a, 0x0e, 0x0a, 0x0a, 0x41, 0x50, 0x49,
	0x4b, 0x65, 0x79, 0x41, 0x75, 0x74, 0x68, 0x12, 0x00, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10, 0x3a,
	0x01, 0x2a, 0x22, 0x0b, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x72, 0x6c, 0x12,
	0xb9, 0x01, 0x0a, 0x09, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x52, 0x4c, 0x12, 0x21, 0x2e,
	0x75, 0x72, 0x6c, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x22, 0x2e, 0x75, 0x72, 0x6c, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e,
	0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x65, 0x92, 0x41, 0x3f, 0x0a, 0x13, 0x55, 0x52, 0x4c, 0x53, 0x68,
	0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x0a,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x20, 0x55, 0x52, 0x4c, 0x2a, 0x0a, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x5f, 0x75, 0x72, 0x6c, 0x62, 0x10, 0x0a, 0x0e, 0x0a, 0x0a, 0x41, 0x50, 0x49, 0x4b,
	0x65, 0x79, 0x41, 0x75, 0x74, 0x68, 0x12, 0x00, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1d, 0x3a, 0x01,
	0x2a, 0x1a, 0x18, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x72, 0x6c, 0x2f, 0x7b,
	0x75, 0x72, 0x6c, 0x2e, 0x75, 0x72, 0x6c, 0x5f, 0x69, 0x64, 0x7d, 0x12, 0xb2, 0x01, 0x0a, 0x09,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x52, 0x4c, 0x12, 0x21, 0x2e, 0x75, 0x72, 0x6c, 0x73,
	0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x75,
	0x72, 0x6c, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x5e, 0x92, 0x41, 0x3f, 0x0a, 0x13, 0x55, 0x52, 0x4c, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65,
	0x6e, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x0a, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x20, 0x55, 0x52, 0x4c, 0x2a, 0x0a, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x5f, 0x75,
	0x72, 0x6c, 0x62, 0x10, 0x0a, 0x0e, 0x0a, 0x0a, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x41, 0x75,
	0x74, 0x68, 0x12, 0x00, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x16, 0x2a, 0x14, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x76, 0x31, 0x2f, 0x75, 0x72, 0x6c, 0x2f, 0x7b, 0x75, 0x72, 0x6c, 0x5f, 0x69, 0x64, 0x7d,
	0x42, 0xf6, 0x02, 0x92, 0x41, 0xa9, 0x01, 0x12, 0x1f, 0x0a, 0x18, 0x55, 0x52, 0x4c, 0x53, 0x68,
	0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x20, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x20,
	0x41, 0x50, 0x49, 0x32, 0x03, 0x31, 0x2e, 0x30, 0x2a, 0x02, 0x01, 0x02, 0x32, 0x10, 0x61, 0x70,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x3a, 0x10,
	0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6a, 0x73, 0x6f, 0x6e,
	0x5a, 0x1f, 0x0a, 0x1d, 0x0a, 0x0a, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x41, 0x75, 0x74, 0x68,
	0x12, 0x0f, 0x08, 0x02, 0x1a, 0x09, 0x58, 0x2d, 0x41, 0x50, 0x49, 0x2d, 0x4b, 0x65, 0x79, 0x20,
	0x02, 0x62, 0x10, 0x0a, 0x0e, 0x0a, 0x0a, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x41, 0x75, 0x74,
	0x68, 0x12, 0x00, 0x6a, 0x2b, 0x0a, 0x13, 0x55, 0x52, 0x4c, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65,
	0x6e, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x14, 0x55, 0x52, 0x4c, 0x53,
	0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x20, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x0a, 0x13, 0x63, 0x6f, 0x6d, 0x2e, 0x75, 0x72, 0x6c, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e,
	0x65, 0x72, 0x2e, 0x76, 0x31, 0x42, 0x11, 0x55, 0x72, 0x6c, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65,
	0x6e, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x42, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6d, 0x67, 0x73, 0x6a, 0x2f, 0x75, 0x72, 0x6c,
	0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x75, 0x72, 0x6c, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x3b,
	0x75, 0x72, 0x6c, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x76, 0x31, 0xa2, 0x02,
	0x03, 0x55, 0x58, 0x58, 0xaa, 0x02, 0x0f, 0x55, 0x72, 0x6c, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65,
	0x6e, 0x65, 0x72, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0f, 0x55, 0x72, 0x6c, 0x73, 0x68, 0x6f, 0x72,
	0x74, 0x65, 0x6e, 0x65, 0x72, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1b, 0x55, 0x72, 0x6c, 0x73, 0x68,
	0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x10, 0x55, 0x72, 0x6c, 0x73, 0x68, 0x6f, 0x72,
	0x74, 0x65, 0x6e, 0x65, 0x72, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
})

var (
	file_urlshortener_v1_urlshortener_proto_rawDescOnce sync.Once
	file_urlshortener_v1_urlshortener_proto_rawDescData = file_urlshortener_v1_urlshortener_proto_rawDesc
)

func file_urlshortener_v1_urlshortener_proto_rawDescGZIP() []byte {
	file_urlshortener_v1_urlshortener_proto_rawDescOnce.Do(func() {
		file_urlshortener_v1_urlshortener_proto_rawDescData = string(protoimpl.X.CompressGZIP([]byte(file_urlshortener_v1_urlshortener_proto_rawDescData)))
	})
	return []byte(file_urlshortener_v1_urlshortener_proto_rawDescData)
}

var file_urlshortener_v1_urlshortener_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_urlshortener_v1_urlshortener_proto_goTypes = []any{
	(*URL)(nil),               // 0: urlshortener.v1.URL
	(*ListURLsRequest)(nil),   // 1: urlshortener.v1.ListURLsRequest
	(*ListURLsResponse)(nil),  // 2: urlshortener.v1.ListURLsResponse
	(*GetURLRequest)(nil),     // 3: urlshortener.v1.GetURLRequest
	(*GetURLResponse)(nil),    // 4: urlshortener.v1.GetURLResponse
	(*CreateURLRequest)(nil),  // 5: urlshortener.v1.CreateURLRequest
	(*CreateURLResponse)(nil), // 6: urlshortener.v1.CreateURLResponse
	(*UpdateURLRequest)(nil),  // 7: urlshortener.v1.UpdateURLRequest
	(*UpdateURLResponse)(nil), // 8: urlshortener.v1.UpdateURLResponse
	(*DeleteURLRequest)(nil),  // 9: urlshortener.v1.DeleteURLRequest
	(*DeleteURLResponse)(nil), // 10: urlshortener.v1.DeleteURLResponse
}
var file_urlshortener_v1_urlshortener_proto_depIdxs = []int32{
	0,  // 0: urlshortener.v1.ListURLsResponse.urls:type_name -> urlshortener.v1.URL
	0,  // 1: urlshortener.v1.GetURLResponse.url:type_name -> urlshortener.v1.URL
	0,  // 2: urlshortener.v1.UpdateURLRequest.url:type_name -> urlshortener.v1.URL
	1,  // 3: urlshortener.v1.URLShortenerService.ListURLs:input_type -> urlshortener.v1.ListURLsRequest
	3,  // 4: urlshortener.v1.URLShortenerService.GetURL:input_type -> urlshortener.v1.GetURLRequest
	5,  // 5: urlshortener.v1.URLShortenerService.CreateURL:input_type -> urlshortener.v1.CreateURLRequest
	7,  // 6: urlshortener.v1.URLShortenerService.UpdateURL:input_type -> urlshortener.v1.UpdateURLRequest
	9,  // 7: urlshortener.v1.URLShortenerService.DeleteURL:input_type -> urlshortener.v1.DeleteURLRequest
	2,  // 8: urlshortener.v1.URLShortenerService.ListURLs:output_type -> urlshortener.v1.ListURLsResponse
	4,  // 9: urlshortener.v1.URLShortenerService.GetURL:output_type -> urlshortener.v1.GetURLResponse
	6,  // 10: urlshortener.v1.URLShortenerService.CreateURL:output_type -> urlshortener.v1.CreateURLResponse
	8,  // 11: urlshortener.v1.URLShortenerService.UpdateURL:output_type -> urlshortener.v1.UpdateURLResponse
	10, // 12: urlshortener.v1.URLShortenerService.DeleteURL:output_type -> urlshortener.v1.DeleteURLResponse
	8,  // [8:13] is the sub-list for method output_type
	3,  // [3:8] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_urlshortener_v1_urlshortener_proto_init() }
func file_urlshortener_v1_urlshortener_proto_init() {
	if File_urlshortener_v1_urlshortener_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_urlshortener_v1_urlshortener_proto_rawDesc), len(file_urlshortener_v1_urlshortener_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_urlshortener_v1_urlshortener_proto_goTypes,
		DependencyIndexes: file_urlshortener_v1_urlshortener_proto_depIdxs,
		MessageInfos:      file_urlshortener_v1_urlshortener_proto_msgTypes,
	}.Build()
	File_urlshortener_v1_urlshortener_proto = out.File
	file_urlshortener_v1_urlshortener_proto_goTypes = nil
	file_urlshortener_v1_urlshortener_proto_depIdxs = nil
}
