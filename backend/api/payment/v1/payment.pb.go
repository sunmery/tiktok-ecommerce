// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: v1/payment.proto

package paymentv1

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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

type StringList struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Values        []string               `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty"` // 表示 []string
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StringList) Reset() {
	*x = StringList{}
	mi := &file_v1_payment_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StringList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StringList) ProtoMessage() {}

func (x *StringList) ProtoReflect() protoreflect.Message {
	mi := &file_v1_payment_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StringList.ProtoReflect.Descriptor instead.
func (*StringList) Descriptor() ([]byte, []int) {
	return file_v1_payment_proto_rawDescGZIP(), []int{0}
}

func (x *StringList) GetValues() []string {
	if x != nil {
		return x.Values
	}
	return nil
}

type PaymentNotifyReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Values        map[string]*StringList `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PaymentNotifyReq) Reset() {
	*x = PaymentNotifyReq{}
	mi := &file_v1_payment_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PaymentNotifyReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaymentNotifyReq) ProtoMessage() {}

func (x *PaymentNotifyReq) ProtoReflect() protoreflect.Message {
	mi := &file_v1_payment_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PaymentNotifyReq.ProtoReflect.Descriptor instead.
func (*PaymentNotifyReq) Descriptor() ([]byte, []int) {
	return file_v1_payment_proto_rawDescGZIP(), []int{1}
}

func (x *PaymentNotifyReq) GetValues() map[string]*StringList {
	if x != nil {
		return x.Values
	}
	return nil
}

type PaymentNotifyResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Code          string                 `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg           string                 `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PaymentNotifyResp) Reset() {
	*x = PaymentNotifyResp{}
	mi := &file_v1_payment_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PaymentNotifyResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaymentNotifyResp) ProtoMessage() {}

func (x *PaymentNotifyResp) ProtoReflect() protoreflect.Message {
	mi := &file_v1_payment_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PaymentNotifyResp.ProtoReflect.Descriptor instead.
func (*PaymentNotifyResp) Descriptor() ([]byte, []int) {
	return file_v1_payment_proto_rawDescGZIP(), []int{2}
}

func (x *PaymentNotifyResp) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *PaymentNotifyResp) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type CreatePaymentReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OrderId       string                 `protobuf:"bytes,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`                   // 主订单ID
	Currency      string                 `protobuf:"bytes,2,opt,name=currency,proto3" json:"currency,omitempty"`                                // 支付货币
	Amount        string                 `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount,omitempty"`                                    // 支付金额
	PaymentMethod string                 `protobuf:"bytes,4,opt,name=payment_method,json=paymentMethod,proto3" json:"payment_method,omitempty"` // 支付方式
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreatePaymentReq) Reset() {
	*x = CreatePaymentReq{}
	mi := &file_v1_payment_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreatePaymentReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreatePaymentReq) ProtoMessage() {}

func (x *CreatePaymentReq) ProtoReflect() protoreflect.Message {
	mi := &file_v1_payment_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreatePaymentReq.ProtoReflect.Descriptor instead.
func (*CreatePaymentReq) Descriptor() ([]byte, []int) {
	return file_v1_payment_proto_rawDescGZIP(), []int{3}
}

func (x *CreatePaymentReq) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

func (x *CreatePaymentReq) GetCurrency() string {
	if x != nil {
		return x.Currency
	}
	return ""
}

func (x *CreatePaymentReq) GetAmount() string {
	if x != nil {
		return x.Amount
	}
	return ""
}

func (x *CreatePaymentReq) GetPaymentMethod() string {
	if x != nil {
		return x.PaymentMethod
	}
	return ""
}

type PaymentResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	PaymentId     string                 `protobuf:"bytes,1,opt,name=payment_id,json=paymentId,proto3" json:"payment_id,omitempty"`
	Status        string                 `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`                           // CREATED/PENDING/SUCCEEDED/FAILED
	PaymentUrl    string                 `protobuf:"bytes,3,opt,name=payment_url,json=paymentUrl,proto3" json:"payment_url,omitempty"` // 第三方支付链接
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PaymentResp) Reset() {
	*x = PaymentResp{}
	mi := &file_v1_payment_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PaymentResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaymentResp) ProtoMessage() {}

func (x *PaymentResp) ProtoReflect() protoreflect.Message {
	mi := &file_v1_payment_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PaymentResp.ProtoReflect.Descriptor instead.
func (*PaymentResp) Descriptor() ([]byte, []int) {
	return file_v1_payment_proto_rawDescGZIP(), []int{4}
}

func (x *PaymentResp) GetPaymentId() string {
	if x != nil {
		return x.PaymentId
	}
	return ""
}

func (x *PaymentResp) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *PaymentResp) GetPaymentUrl() string {
	if x != nil {
		return x.PaymentUrl
	}
	return ""
}

func (x *PaymentResp) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

type GetPaymentReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	PaymentId     string                 `protobuf:"bytes,1,opt,name=payment_id,json=paymentId,proto3" json:"payment_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetPaymentReq) Reset() {
	*x = GetPaymentReq{}
	mi := &file_v1_payment_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPaymentReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPaymentReq) ProtoMessage() {}

func (x *GetPaymentReq) ProtoReflect() protoreflect.Message {
	mi := &file_v1_payment_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPaymentReq.ProtoReflect.Descriptor instead.
func (*GetPaymentReq) Descriptor() ([]byte, []int) {
	return file_v1_payment_proto_rawDescGZIP(), []int{5}
}

func (x *GetPaymentReq) GetPaymentId() string {
	if x != nil {
		return x.PaymentId
	}
	return ""
}

type PaymentCallbackReq struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	PaymentId       string                 `protobuf:"bytes,1,opt,name=payment_id,json=paymentId,proto3" json:"payment_id,omitempty"`
	Status          string                 `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	GatewayResponse string                 `protobuf:"bytes,3,opt,name=gateway_response,json=gatewayResponse,proto3" json:"gateway_response,omitempty"` // 原始网关响应
	ProcessedAt     *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=processed_at,json=processedAt,proto3" json:"processed_at,omitempty"`
	RequestForm     map[string]*StringList `protobuf:"bytes,5,rep,name=request_form,json=requestForm,proto3" json:"request_form,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *PaymentCallbackReq) Reset() {
	*x = PaymentCallbackReq{}
	mi := &file_v1_payment_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PaymentCallbackReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaymentCallbackReq) ProtoMessage() {}

func (x *PaymentCallbackReq) ProtoReflect() protoreflect.Message {
	mi := &file_v1_payment_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PaymentCallbackReq.ProtoReflect.Descriptor instead.
func (*PaymentCallbackReq) Descriptor() ([]byte, []int) {
	return file_v1_payment_proto_rawDescGZIP(), []int{6}
}

func (x *PaymentCallbackReq) GetPaymentId() string {
	if x != nil {
		return x.PaymentId
	}
	return ""
}

func (x *PaymentCallbackReq) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *PaymentCallbackReq) GetGatewayResponse() string {
	if x != nil {
		return x.GatewayResponse
	}
	return ""
}

func (x *PaymentCallbackReq) GetProcessedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.ProcessedAt
	}
	return nil
}

func (x *PaymentCallbackReq) GetRequestForm() map[string]*StringList {
	if x != nil {
		return x.RequestForm
	}
	return nil
}

type PaymentCallbackResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PaymentCallbackResp) Reset() {
	*x = PaymentCallbackResp{}
	mi := &file_v1_payment_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PaymentCallbackResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaymentCallbackResp) ProtoMessage() {}

func (x *PaymentCallbackResp) ProtoReflect() protoreflect.Message {
	mi := &file_v1_payment_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PaymentCallbackResp.ProtoReflect.Descriptor instead.
func (*PaymentCallbackResp) Descriptor() ([]byte, []int) {
	return file_v1_payment_proto_rawDescGZIP(), []int{7}
}

var File_v1_payment_proto protoreflect.FileDescriptor

var file_v1_payment_proto_rawDesc = string([]byte{
	0x0a, 0x10, 0x76, 0x31, 0x2f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x14, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x61,
	0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x24, 0x0a, 0x0a, 0x53, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x22, 0xbb, 0x01,
	0x0a, 0x10, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x52,
	0x65, 0x71, 0x12, 0x4a, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x32, 0x2e, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x72, 0x63, 0x65, 0x2e, 0x70,
	0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e,
	0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x52, 0x65, 0x71, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x1a, 0x5b,
	0x0a, 0x0b, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x36, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20,
	0x2e, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73, 0x74,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x39, 0x0a, 0x11, 0x50,
	0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x52, 0x65, 0x73, 0x70,
	0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x88, 0x01, 0x0a, 0x10, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x12, 0x19, 0x0a, 0x08, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x63, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x63, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x70, 0x61,
	0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x22, 0xa0, 0x01, 0x0a, 0x0b, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64,
	0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x61, 0x79, 0x6d,
	0x65, 0x6e, 0x74, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70,
	0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x55, 0x72, 0x6c, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x22, 0x2e, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x50, 0x61, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x71, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x61, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x49, 0x64, 0x22, 0xf5, 0x02, 0x0a, 0x12, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x43, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x12, 0x1d, 0x0a, 0x0a, 0x70,
	0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x29, 0x0a, 0x10, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x5f, 0x72, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x67, 0x61,
	0x74, 0x65, 0x77, 0x61, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3d, 0x0a,
	0x0c, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x0b, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x64, 0x41, 0x74, 0x12, 0x5c, 0x0a, 0x0c,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x66, 0x6f, 0x72, 0x6d, 0x18, 0x05, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x39, 0x2e, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x72, 0x63, 0x65, 0x2e, 0x70,
	0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e,
	0x74, 0x43, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x2e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x46, 0x6f, 0x72, 0x6d, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0b, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x46, 0x6f, 0x72, 0x6d, 0x1a, 0x60, 0x0a, 0x10, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x46, 0x6f, 0x72, 0x6d, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x36, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x20, 0x2e, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x61, 0x79, 0x6d,
	0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73,
	0x74, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x15, 0x0a, 0x13,
	0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x52,
	0x65, 0x73, 0x70, 0x32, 0x93, 0x04, 0x0a, 0x0e, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x73, 0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x26, 0x2e, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65,
	0x72, 0x63, 0x65, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x1a,
	0x21, 0x2e, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x61, 0x79, 0x6d,
	0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x22, 0x17, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x11, 0x3a, 0x01, 0x2a, 0x22, 0x0c, 0x2f,
	0x76, 0x31, 0x2f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x80, 0x01, 0x0a, 0x0d,
	0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x12, 0x26, 0x2e,
	0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e,
	0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x4e, 0x6f, 0x74, 0x69,
	0x66, 0x79, 0x52, 0x65, 0x71, 0x1a, 0x27, 0x2e, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x72, 0x63,
	0x65, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x52, 0x65, 0x73, 0x70, 0x22, 0x1e,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x3a, 0x01, 0x2a, 0x22, 0x13, 0x2f, 0x76, 0x31, 0x2f, 0x70,
	0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x12, 0x8f,
	0x01, 0x0a, 0x16, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e,
	0x74, 0x43, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x12, 0x28, 0x2e, 0x65, 0x63, 0x6f, 0x6d,
	0x6d, 0x65, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31,
	0x2e, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b,
	0x52, 0x65, 0x71, 0x1a, 0x29, 0x2e, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x72, 0x63, 0x65, 0x2e,
	0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x43, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x22, 0x20,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1a, 0x3a, 0x01, 0x2a, 0x22, 0x15, 0x2f, 0x76, 0x31, 0x2f, 0x70,
	0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x63, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b,
	0x12, 0x77, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x23,
	0x2e, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x71, 0x1a, 0x21, 0x2e, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x72, 0x63, 0x65, 0x2e,
	0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x22, 0x21, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1b, 0x12, 0x19,
	0x2f, 0x76, 0x31, 0x2f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x7b, 0x70, 0x61,
	0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x7d, 0x42, 0x22, 0x5a, 0x20, 0x62, 0x61, 0x63,
	0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x2f, 0x76, 0x31, 0x3b, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x76, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_v1_payment_proto_rawDescOnce sync.Once
	file_v1_payment_proto_rawDescData []byte
)

func file_v1_payment_proto_rawDescGZIP() []byte {
	file_v1_payment_proto_rawDescOnce.Do(func() {
		file_v1_payment_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_v1_payment_proto_rawDesc), len(file_v1_payment_proto_rawDesc)))
	})
	return file_v1_payment_proto_rawDescData
}

var file_v1_payment_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_v1_payment_proto_goTypes = []any{
	(*StringList)(nil),            // 0: ecommerce.payment.v1.StringList
	(*PaymentNotifyReq)(nil),      // 1: ecommerce.payment.v1.PaymentNotifyReq
	(*PaymentNotifyResp)(nil),     // 2: ecommerce.payment.v1.PaymentNotifyResp
	(*CreatePaymentReq)(nil),      // 3: ecommerce.payment.v1.CreatePaymentReq
	(*PaymentResp)(nil),           // 4: ecommerce.payment.v1.PaymentResp
	(*GetPaymentReq)(nil),         // 5: ecommerce.payment.v1.GetPaymentReq
	(*PaymentCallbackReq)(nil),    // 6: ecommerce.payment.v1.PaymentCallbackReq
	(*PaymentCallbackResp)(nil),   // 7: ecommerce.payment.v1.PaymentCallbackResp
	nil,                           // 8: ecommerce.payment.v1.PaymentNotifyReq.ValuesEntry
	nil,                           // 9: ecommerce.payment.v1.PaymentCallbackReq.RequestFormEntry
	(*timestamppb.Timestamp)(nil), // 10: google.protobuf.Timestamp
}
var file_v1_payment_proto_depIdxs = []int32{
	8,  // 0: ecommerce.payment.v1.PaymentNotifyReq.values:type_name -> ecommerce.payment.v1.PaymentNotifyReq.ValuesEntry
	10, // 1: ecommerce.payment.v1.PaymentResp.created_at:type_name -> google.protobuf.Timestamp
	10, // 2: ecommerce.payment.v1.PaymentCallbackReq.processed_at:type_name -> google.protobuf.Timestamp
	9,  // 3: ecommerce.payment.v1.PaymentCallbackReq.request_form:type_name -> ecommerce.payment.v1.PaymentCallbackReq.RequestFormEntry
	0,  // 4: ecommerce.payment.v1.PaymentNotifyReq.ValuesEntry.value:type_name -> ecommerce.payment.v1.StringList
	0,  // 5: ecommerce.payment.v1.PaymentCallbackReq.RequestFormEntry.value:type_name -> ecommerce.payment.v1.StringList
	3,  // 6: ecommerce.payment.v1.PaymentService.CreatePayment:input_type -> ecommerce.payment.v1.CreatePaymentReq
	1,  // 7: ecommerce.payment.v1.PaymentService.PaymentNotify:input_type -> ecommerce.payment.v1.PaymentNotifyReq
	6,  // 8: ecommerce.payment.v1.PaymentService.ProcessPaymentCallback:input_type -> ecommerce.payment.v1.PaymentCallbackReq
	5,  // 9: ecommerce.payment.v1.PaymentService.GetPayment:input_type -> ecommerce.payment.v1.GetPaymentReq
	4,  // 10: ecommerce.payment.v1.PaymentService.CreatePayment:output_type -> ecommerce.payment.v1.PaymentResp
	2,  // 11: ecommerce.payment.v1.PaymentService.PaymentNotify:output_type -> ecommerce.payment.v1.PaymentNotifyResp
	7,  // 12: ecommerce.payment.v1.PaymentService.ProcessPaymentCallback:output_type -> ecommerce.payment.v1.PaymentCallbackResp
	4,  // 13: ecommerce.payment.v1.PaymentService.GetPayment:output_type -> ecommerce.payment.v1.PaymentResp
	10, // [10:14] is the sub-list for method output_type
	6,  // [6:10] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_v1_payment_proto_init() }
func file_v1_payment_proto_init() {
	if File_v1_payment_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_v1_payment_proto_rawDesc), len(file_v1_payment_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_v1_payment_proto_goTypes,
		DependencyIndexes: file_v1_payment_proto_depIdxs,
		MessageInfos:      file_v1_payment_proto_msgTypes,
	}.Build()
	File_v1_payment_proto = out.File
	file_v1_payment_proto_goTypes = nil
	file_v1_payment_proto_depIdxs = nil
}
