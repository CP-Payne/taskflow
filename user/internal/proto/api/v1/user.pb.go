// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.27.0
// source: api/v1/user.proto

package user_v1

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

type AuthenticateUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email    string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *AuthenticateUserRequest) Reset() {
	*x = AuthenticateUserRequest{}
	mi := &file_api_v1_user_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AuthenticateUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthenticateUserRequest) ProtoMessage() {}

func (x *AuthenticateUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_user_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthenticateUserRequest.ProtoReflect.Descriptor instead.
func (*AuthenticateUserRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_user_proto_rawDescGZIP(), []int{0}
}

func (x *AuthenticateUserRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *AuthenticateUserRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type AuthenticateUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Jwt string `protobuf:"bytes,1,opt,name=jwt,proto3" json:"jwt,omitempty"`
}

func (x *AuthenticateUserResponse) Reset() {
	*x = AuthenticateUserResponse{}
	mi := &file_api_v1_user_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AuthenticateUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthenticateUserResponse) ProtoMessage() {}

func (x *AuthenticateUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_user_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthenticateUserResponse.ProtoReflect.Descriptor instead.
func (*AuthenticateUserResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_user_proto_rawDescGZIP(), []int{1}
}

func (x *AuthenticateUserResponse) GetJwt() string {
	if x != nil {
		return x.Jwt
	}
	return ""
}

type RegisterUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email    string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Username string `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Password string `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *RegisterUserRequest) Reset() {
	*x = RegisterUserRequest{}
	mi := &file_api_v1_user_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterUserRequest) ProtoMessage() {}

func (x *RegisterUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_user_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterUserRequest.ProtoReflect.Descriptor instead.
func (*RegisterUserRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_user_proto_rawDescGZIP(), []int{2}
}

func (x *RegisterUserRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *RegisterUserRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *RegisterUserRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type RegisterUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Jwt string `protobuf:"bytes,1,opt,name=jwt,proto3" json:"jwt,omitempty"`
}

func (x *RegisterUserResponse) Reset() {
	*x = RegisterUserResponse{}
	mi := &file_api_v1_user_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterUserResponse) ProtoMessage() {}

func (x *RegisterUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_user_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterUserResponse.ProtoReflect.Descriptor instead.
func (*RegisterUserResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_user_proto_rawDescGZIP(), []int{3}
}

func (x *RegisterUserResponse) GetJwt() string {
	if x != nil {
		return x.Jwt
	}
	return ""
}

type GetByIDRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *GetByIDRequest) Reset() {
	*x = GetByIDRequest{}
	mi := &file_api_v1_user_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetByIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetByIDRequest) ProtoMessage() {}

func (x *GetByIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_user_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetByIDRequest.ProtoReflect.Descriptor instead.
func (*GetByIDRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_user_proto_rawDescGZIP(), []int{4}
}

func (x *GetByIDRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type GetByIDResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId   string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Email    string `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Username string `protobuf:"bytes,3,opt,name=username,proto3" json:"username,omitempty"`
}

func (x *GetByIDResponse) Reset() {
	*x = GetByIDResponse{}
	mi := &file_api_v1_user_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetByIDResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetByIDResponse) ProtoMessage() {}

func (x *GetByIDResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_user_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetByIDResponse.ProtoReflect.Descriptor instead.
func (*GetByIDResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_user_proto_rawDescGZIP(), []int{5}
}

func (x *GetByIDResponse) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *GetByIDResponse) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *GetByIDResponse) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

var File_api_v1_user_proto protoreflect.FileDescriptor

var file_api_v1_user_proto_rawDesc = []byte{
	0x0a, 0x11, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x07, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x22, 0x4b, 0x0a, 0x17,
	0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1a, 0x0a,
	0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x2c, 0x0a, 0x18, 0x41, 0x75, 0x74,
	0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6a, 0x77, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6a, 0x77, 0x74, 0x22, 0x63, 0x0a, 0x13, 0x52, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14,
	0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x28, 0x0a, 0x14,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6a, 0x77, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6a, 0x77, 0x74, 0x22, 0x29, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x42, 0x79, 0x49,
	0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x22, 0x5c, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x14, 0x0a,
	0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d,
	0x61, 0x69, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x32,
	0xf0, 0x01, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x59, 0x0a, 0x10, 0x41, 0x75, 0x74, 0x68,
	0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x20, 0x2e, 0x75,
	0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63,
	0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21,
	0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74,
	0x69, 0x63, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x4d, 0x0a, 0x0c, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x55,
	0x73, 0x65, 0x72, 0x12, 0x1c, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1d, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x3e, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x42, 0x79, 0x49, 0x44, 0x12, 0x17, 0x2e,
	0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x79, 0x49, 0x44, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31,
	0x2e, 0x47, 0x65, 0x74, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x43, 0x50, 0x2d, 0x50, 0x61, 0x79, 0x6e, 0x65, 0x2f, 0x74, 0x61, 0x73, 0x6b, 0x66, 0x6c,
	0x6f, 0x77, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x76, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_v1_user_proto_rawDescOnce sync.Once
	file_api_v1_user_proto_rawDescData = file_api_v1_user_proto_rawDesc
)

func file_api_v1_user_proto_rawDescGZIP() []byte {
	file_api_v1_user_proto_rawDescOnce.Do(func() {
		file_api_v1_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_user_proto_rawDescData)
	})
	return file_api_v1_user_proto_rawDescData
}

var file_api_v1_user_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_api_v1_user_proto_goTypes = []any{
	(*AuthenticateUserRequest)(nil),  // 0: user.v1.AuthenticateUserRequest
	(*AuthenticateUserResponse)(nil), // 1: user.v1.AuthenticateUserResponse
	(*RegisterUserRequest)(nil),      // 2: user.v1.RegisterUserRequest
	(*RegisterUserResponse)(nil),     // 3: user.v1.RegisterUserResponse
	(*GetByIDRequest)(nil),           // 4: user.v1.GetByIDRequest
	(*GetByIDResponse)(nil),          // 5: user.v1.GetByIDResponse
}
var file_api_v1_user_proto_depIdxs = []int32{
	0, // 0: user.v1.User.AuthenticateUser:input_type -> user.v1.AuthenticateUserRequest
	2, // 1: user.v1.User.RegisterUser:input_type -> user.v1.RegisterUserRequest
	4, // 2: user.v1.User.GetByID:input_type -> user.v1.GetByIDRequest
	1, // 3: user.v1.User.AuthenticateUser:output_type -> user.v1.AuthenticateUserResponse
	3, // 4: user.v1.User.RegisterUser:output_type -> user.v1.RegisterUserResponse
	5, // 5: user.v1.User.GetByID:output_type -> user.v1.GetByIDResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_v1_user_proto_init() }
func file_api_v1_user_proto_init() {
	if File_api_v1_user_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_v1_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_v1_user_proto_goTypes,
		DependencyIndexes: file_api_v1_user_proto_depIdxs,
		MessageInfos:      file_api_v1_user_proto_msgTypes,
	}.Build()
	File_api_v1_user_proto = out.File
	file_api_v1_user_proto_rawDesc = nil
	file_api_v1_user_proto_goTypes = nil
	file_api_v1_user_proto_depIdxs = nil
}
