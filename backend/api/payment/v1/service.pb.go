// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        v5.29.2
// source: v1/service.proto

package payment

import (
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

type CreditCardInfo struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	Number          string                 `protobuf:"bytes,1,opt,name=number,proto3" json:"number,omitempty"`
	Cvv             int32                  `protobuf:"varint,2,opt,name=cvv,proto3" json:"cvv,omitempty"`
	ExpirationYear  int32                  `protobuf:"varint,3,opt,name=expiration_year,json=expirationYear,proto3" json:"expiration_year,omitempty"`
	ExpirationMonth int32                  `protobuf:"varint,4,opt,name=expiration_month,json=expirationMonth,proto3" json:"expiration_month,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *CreditCardInfo) Reset() {
	*x = CreditCardInfo{}
	mi := &file_v1_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreditCardInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreditCardInfo) ProtoMessage() {}

func (x *CreditCardInfo) ProtoReflect() protoreflect.Message {
	mi := &file_v1_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreditCardInfo.ProtoReflect.Descriptor instead.
func (*CreditCardInfo) Descriptor() ([]byte, []int) {
	return file_v1_service_proto_rawDescGZIP(), []int{0}
}

func (x *CreditCardInfo) GetNumber() string {
	if x != nil {
		return x.Number
	}
	return ""
}

func (x *CreditCardInfo) GetCvv() int32 {
	if x != nil {
		return x.Cvv
	}
	return 0
}

func (x *CreditCardInfo) GetExpirationYear() int32 {
	if x != nil {
		return x.ExpirationYear
	}
	return 0
}

func (x *CreditCardInfo) GetExpirationMonth() int32 {
	if x != nil {
		return x.ExpirationMonth
	}
	return 0
}

type ChargeReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Amount        float32                `protobuf:"fixed32,1,opt,name=amount,proto3" json:"amount,omitempty"`
	CreditCard    *CreditCardInfo        `protobuf:"bytes,2,opt,name=credit_card,json=creditCard,proto3" json:"credit_card,omitempty"`
	OrderId       string                 `protobuf:"bytes,3,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	UserId        string                 `protobuf:"bytes,4,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChargeReq) Reset() {
	*x = ChargeReq{}
	mi := &file_v1_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChargeReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChargeReq) ProtoMessage() {}

func (x *ChargeReq) ProtoReflect() protoreflect.Message {
	mi := &file_v1_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChargeReq.ProtoReflect.Descriptor instead.
func (*ChargeReq) Descriptor() ([]byte, []int) {
	return file_v1_service_proto_rawDescGZIP(), []int{1}
}

func (x *ChargeReq) GetAmount() float32 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *ChargeReq) GetCreditCard() *CreditCardInfo {
	if x != nil {
		return x.CreditCard
	}
	return nil
}

func (x *ChargeReq) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

func (x *ChargeReq) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type ChargeResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TransactionId string                 `protobuf:"bytes,1,opt,name=transaction_id,json=transactionId,proto3" json:"transaction_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChargeResp) Reset() {
	*x = ChargeResp{}
	mi := &file_v1_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChargeResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChargeResp) ProtoMessage() {}

func (x *ChargeResp) ProtoReflect() protoreflect.Message {
	mi := &file_v1_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChargeResp.ProtoReflect.Descriptor instead.
func (*ChargeResp) Descriptor() ([]byte, []int) {
	return file_v1_service_proto_rawDescGZIP(), []int{2}
}

func (x *ChargeResp) GetTransactionId() string {
	if x != nil {
		return x.TransactionId
	}
	return ""
}

var File_v1_service_proto protoreflect.FileDescriptor

var file_v1_service_proto_rawDesc = string([]byte{
	0x0a, 0x10, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x12, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x22, 0x8e, 0x01, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x64, 0x69,
	0x74, 0x43, 0x61, 0x72, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x12, 0x10, 0x0a, 0x03, 0x63, 0x76, 0x76, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03,
	0x63, 0x76, 0x76, 0x12, 0x27, 0x0a, 0x0f, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x79, 0x65, 0x61, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x65, 0x78,
	0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x59, 0x65, 0x61, 0x72, 0x12, 0x29, 0x0a, 0x10,
	0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6d, 0x6f, 0x6e, 0x74, 0x68,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0f, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x4d, 0x6f, 0x6e, 0x74, 0x68, 0x22, 0x9c, 0x01, 0x0a, 0x09, 0x43, 0x68, 0x61, 0x72,
	0x67, 0x65, 0x52, 0x65, 0x71, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x02, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x43, 0x0a,
	0x0b, 0x63, 0x72, 0x65, 0x64, 0x69, 0x74, 0x5f, 0x63, 0x61, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x22, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x64, 0x69, 0x74, 0x43, 0x61,
	0x72, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x64, 0x69, 0x74, 0x43, 0x61,
	0x72, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x17, 0x0a,
	0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x33, 0x0a, 0x0a, 0x43, 0x68, 0x61, 0x72, 0x67, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x12, 0x25, 0x0a, 0x0e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x32, 0x5b, 0x0a, 0x0e, 0x50,
	0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x49, 0x0a,
	0x06, 0x43, 0x68, 0x61, 0x72, 0x67, 0x65, 0x12, 0x1d, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e,
	0x74, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x68, 0x61,
	0x72, 0x67, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x1e, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x68, 0x61, 0x72,
	0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x22, 0x00, 0x42, 0x18, 0x5a, 0x16, 0x61, 0x70, 0x69, 0x2f,
	0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x3b, 0x70, 0x61, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_v1_service_proto_rawDescOnce sync.Once
	file_v1_service_proto_rawDescData []byte
)

func file_v1_service_proto_rawDescGZIP() []byte {
	file_v1_service_proto_rawDescOnce.Do(func() {
		file_v1_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_v1_service_proto_rawDesc), len(file_v1_service_proto_rawDesc)))
	})
	return file_v1_service_proto_rawDescData
}

var file_v1_service_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_v1_service_proto_goTypes = []any{
	(*CreditCardInfo)(nil), // 0: payment.service.v1.CreditCardInfo
	(*ChargeReq)(nil),      // 1: payment.service.v1.ChargeReq
	(*ChargeResp)(nil),     // 2: payment.service.v1.ChargeResp
}
var file_v1_service_proto_depIdxs = []int32{
	0, // 0: payment.service.v1.ChargeReq.credit_card:type_name -> payment.service.v1.CreditCardInfo
	1, // 1: payment.service.v1.PaymentService.Charge:input_type -> payment.service.v1.ChargeReq
	2, // 2: payment.service.v1.PaymentService.Charge:output_type -> payment.service.v1.ChargeResp
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_v1_service_proto_init() }
func file_v1_service_proto_init() {
	if File_v1_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_v1_service_proto_rawDesc), len(file_v1_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_v1_service_proto_goTypes,
		DependencyIndexes: file_v1_service_proto_depIdxs,
		MessageInfos:      file_v1_service_proto_msgTypes,
	}.Build()
	File_v1_service_proto = out.File
	file_v1_service_proto_goTypes = nil
	file_v1_service_proto_depIdxs = nil
}
