// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: v1/category.proto

// 定义包名，用于区分不同的服务模块。

package category

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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

type Categorys struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Categorys     []*Category            `protobuf:"bytes,1,rep,name=categorys,proto3" json:"categorys,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Categorys) Reset() {
	*x = Categorys{}
	mi := &file_v1_category_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Categorys) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Categorys) ProtoMessage() {}

func (x *Categorys) ProtoReflect() protoreflect.Message {
	mi := &file_v1_category_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Categorys.ProtoReflect.Descriptor instead.
func (*Categorys) Descriptor() ([]byte, []int) {
	return file_v1_category_proto_rawDescGZIP(), []int{0}
}

func (x *Categorys) GetCategorys() []*Category {
	if x != nil {
		return x.Categorys
	}
	return nil
}

// 基础数据结构：分类（Category）
type Category struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`                                // 分类的唯一标识符。
	ParentId      int64                  `protobuf:"varint,2,opt,name=parent_id,json=parentId,proto3" json:"parent_id,omitempty"`    // 父分类的 ID，根节点的 parent_id 为 0。
	Level         int32                  `protobuf:"varint,3,opt,name=level,proto3" json:"level,omitempty"`                          // 分类的层级（例如：1 表示一级分类，2 表示二级分类）。
	Path          string                 `protobuf:"bytes,4,opt,name=path,proto3" json:"path,omitempty"`                             // 分类路径（ltree 序列化为字符串，用于快速查询层级关系）。
	Name          string                 `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`                             // 分类名称。
	SortOrder     int32                  `protobuf:"varint,6,opt,name=sort_order,json=sortOrder,proto3" json:"sort_order,omitempty"` // 排序值，用于控制分类的显示顺序。
	IsLeaf        bool                   `protobuf:"varint,7,opt,name=is_leaf,json=isLeaf,proto3" json:"is_leaf,omitempty"`          // 是否为叶子节点（没有子分类）。
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`  // 分类创建时间。
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`  // 分类最后更新时间。
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Category) Reset() {
	*x = Category{}
	mi := &file_v1_category_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Category) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Category) ProtoMessage() {}

func (x *Category) ProtoReflect() protoreflect.Message {
	mi := &file_v1_category_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Category.ProtoReflect.Descriptor instead.
func (*Category) Descriptor() ([]byte, []int) {
	return file_v1_category_proto_rawDescGZIP(), []int{1}
}

func (x *Category) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Category) GetParentId() int64 {
	if x != nil {
		return x.ParentId
	}
	return 0
}

func (x *Category) GetLevel() int32 {
	if x != nil {
		return x.Level
	}
	return 0
}

func (x *Category) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *Category) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Category) GetSortOrder() int32 {
	if x != nil {
		return x.SortOrder
	}
	return 0
}

func (x *Category) GetIsLeaf() bool {
	if x != nil {
		return x.IsLeaf
	}
	return false
}

func (x *Category) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Category) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

// 数据结构：闭包关系（ClosureRelation）
type ClosureRelation struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Ancestor      int64                  `protobuf:"varint,1,opt,name=ancestor,proto3" json:"ancestor,omitempty"`     // 祖先分类的 ID。
	Descendant    int64                  `protobuf:"varint,2,opt,name=descendant,proto3" json:"descendant,omitempty"` // 后代分类的 ID。
	Depth         int32                  `protobuf:"varint,3,opt,name=depth,proto3" json:"depth,omitempty"`           // 祖先与后代之间的层级深度（0 表示自身，1 表示直接子节点）。
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ClosureRelation) Reset() {
	*x = ClosureRelation{}
	mi := &file_v1_category_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ClosureRelation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClosureRelation) ProtoMessage() {}

func (x *ClosureRelation) ProtoReflect() protoreflect.Message {
	mi := &file_v1_category_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClosureRelation.ProtoReflect.Descriptor instead.
func (*ClosureRelation) Descriptor() ([]byte, []int) {
	return file_v1_category_proto_rawDescGZIP(), []int{2}
}

func (x *ClosureRelation) GetAncestor() int64 {
	if x != nil {
		return x.Ancestor
	}
	return 0
}

func (x *ClosureRelation) GetDescendant() int64 {
	if x != nil {
		return x.Descendant
	}
	return 0
}

func (x *ClosureRelation) GetDepth() int32 {
	if x != nil {
		return x.Depth
	}
	return 0
}

// 创建分类请求
type CreateCategoryRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ParentId      int64                  `protobuf:"varint,1,opt,name=parent_id,json=parentId,proto3" json:"parent_id,omitempty"`    // 父分类的 ID，根节点的 parent_id 为 0。
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`                             // 分类名称。
	SortOrder     int32                  `protobuf:"varint,3,opt,name=sort_order,json=sortOrder,proto3" json:"sort_order,omitempty"` // 排序值。
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateCategoryRequest) Reset() {
	*x = CreateCategoryRequest{}
	mi := &file_v1_category_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateCategoryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateCategoryRequest) ProtoMessage() {}

func (x *CreateCategoryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_category_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateCategoryRequest.ProtoReflect.Descriptor instead.
func (*CreateCategoryRequest) Descriptor() ([]byte, []int) {
	return file_v1_category_proto_rawDescGZIP(), []int{3}
}

func (x *CreateCategoryRequest) GetParentId() int64 {
	if x != nil {
		return x.ParentId
	}
	return 0
}

func (x *CreateCategoryRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateCategoryRequest) GetSortOrder() int32 {
	if x != nil {
		return x.SortOrder
	}
	return 0
}

// 获取单个分类请求
type GetCategoryRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"` // 分类的唯一标识符。
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetCategoryRequest) Reset() {
	*x = GetCategoryRequest{}
	mi := &file_v1_category_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetCategoryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCategoryRequest) ProtoMessage() {}

func (x *GetCategoryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_category_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCategoryRequest.ProtoReflect.Descriptor instead.
func (*GetCategoryRequest) Descriptor() ([]byte, []int) {
	return file_v1_category_proto_rawDescGZIP(), []int{4}
}

func (x *GetCategoryRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

// 更新分类请求
type UpdateCategoryRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`    // 分类的唯一标识符。
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"` // 新的分类名称。
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateCategoryRequest) Reset() {
	*x = UpdateCategoryRequest{}
	mi := &file_v1_category_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateCategoryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateCategoryRequest) ProtoMessage() {}

func (x *UpdateCategoryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_category_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateCategoryRequest.ProtoReflect.Descriptor instead.
func (*UpdateCategoryRequest) Descriptor() ([]byte, []int) {
	return file_v1_category_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateCategoryRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateCategoryRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// 删除分类请求
type DeleteCategoryRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"` // 分类的唯一标识符。
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteCategoryRequest) Reset() {
	*x = DeleteCategoryRequest{}
	mi := &file_v1_category_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteCategoryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteCategoryRequest) ProtoMessage() {}

func (x *DeleteCategoryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_category_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteCategoryRequest.ProtoReflect.Descriptor instead.
func (*DeleteCategoryRequest) Descriptor() ([]byte, []int) {
	return file_v1_category_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteCategoryRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

// 获取子树请求
type GetSubTreeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RootId        int64                  `protobuf:"varint,1,opt,name=root_id,json=rootId,proto3" json:"root_id,omitempty"` // 根节点的 ID，用于获取其子树。
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetSubTreeRequest) Reset() {
	*x = GetSubTreeRequest{}
	mi := &file_v1_category_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetSubTreeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSubTreeRequest) ProtoMessage() {}

func (x *GetSubTreeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_category_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSubTreeRequest.ProtoReflect.Descriptor instead.
func (*GetSubTreeRequest) Descriptor() ([]byte, []int) {
	return file_v1_category_proto_rawDescGZIP(), []int{7}
}

func (x *GetSubTreeRequest) GetRootId() int64 {
	if x != nil {
		return x.RootId
	}
	return 0
}

// 获取分类路径请求
type GetCategoryPathRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CategoryId    int64                  `protobuf:"varint,1,opt,name=category_id,json=categoryId,proto3" json:"category_id,omitempty"` // 分类的唯一标识符，用于获取其完整路径。
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetCategoryPathRequest) Reset() {
	*x = GetCategoryPathRequest{}
	mi := &file_v1_category_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetCategoryPathRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCategoryPathRequest) ProtoMessage() {}

func (x *GetCategoryPathRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_category_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCategoryPathRequest.ProtoReflect.Descriptor instead.
func (*GetCategoryPathRequest) Descriptor() ([]byte, []int) {
	return file_v1_category_proto_rawDescGZIP(), []int{8}
}

func (x *GetCategoryPathRequest) GetCategoryId() int64 {
	if x != nil {
		return x.CategoryId
	}
	return 0
}

// 获取闭包关系请求
type GetClosureRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CategoryId    int64                  `protobuf:"varint,1,opt,name=category_id,json=categoryId,proto3" json:"category_id,omitempty"` // 分类的唯一标识符，用于获取其闭包关系。
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetClosureRequest) Reset() {
	*x = GetClosureRequest{}
	mi := &file_v1_category_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetClosureRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetClosureRequest) ProtoMessage() {}

func (x *GetClosureRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_category_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetClosureRequest.ProtoReflect.Descriptor instead.
func (*GetClosureRequest) Descriptor() ([]byte, []int) {
	return file_v1_category_proto_rawDescGZIP(), []int{9}
}

func (x *GetClosureRequest) GetCategoryId() int64 {
	if x != nil {
		return x.CategoryId
	}
	return 0
}

// 更新闭包关系深度请求
type UpdateClosureDepthRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CategoryId    int64                  `protobuf:"varint,1,opt,name=category_id,json=categoryId,proto3" json:"category_id,omitempty"` // 分类的唯一标识符。
	Delta         int32                  `protobuf:"varint,2,opt,name=delta,proto3" json:"delta,omitempty"`                             // 深度变化值（正数表示增加深度，负数表示减少深度）。
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateClosureDepthRequest) Reset() {
	*x = UpdateClosureDepthRequest{}
	mi := &file_v1_category_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateClosureDepthRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateClosureDepthRequest) ProtoMessage() {}

func (x *UpdateClosureDepthRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_category_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateClosureDepthRequest.ProtoReflect.Descriptor instead.
func (*UpdateClosureDepthRequest) Descriptor() ([]byte, []int) {
	return file_v1_category_proto_rawDescGZIP(), []int{10}
}

func (x *UpdateClosureDepthRequest) GetCategoryId() int64 {
	if x != nil {
		return x.CategoryId
	}
	return 0
}

func (x *UpdateClosureDepthRequest) GetDelta() int32 {
	if x != nil {
		return x.Delta
	}
	return 0
}

var File_v1_category_proto protoreflect.FileDescriptor

var file_v1_category_proto_rawDesc = string([]byte{
	0x0a, 0x11, 0x76, 0x31, 0x2f, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x61, 0x70, 0x69, 0x2e, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72,
	0x79, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x44, 0x0a, 0x09, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x73, 0x12, 0x37, 0x0a,
	0x09, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x2e,
	0x76, 0x31, 0x2e, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x52, 0x09, 0x63, 0x61, 0x74,
	0x65, 0x67, 0x6f, 0x72, 0x79, 0x73, 0x22, 0xa3, 0x02, 0x0a, 0x08, 0x43, 0x61, 0x74, 0x65, 0x67,
	0x6f, 0x72, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x49, 0x64,
	0x12, 0x14, 0x0a, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1d,
	0x0a, 0x0a, 0x73, 0x6f, 0x72, 0x74, 0x5f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x09, 0x73, 0x6f, 0x72, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x17, 0x0a,
	0x07, 0x69, 0x73, 0x5f, 0x6c, 0x65, 0x61, 0x66, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06,
	0x69, 0x73, 0x4c, 0x65, 0x61, 0x66, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x5f, 0x61, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x63, 0x0a, 0x0f,
	0x43, 0x6c, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x1a, 0x0a, 0x08, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x08, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x64,
	0x65, 0x73, 0x63, 0x65, 0x6e, 0x64, 0x61, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0a, 0x64, 0x65, 0x73, 0x63, 0x65, 0x6e, 0x64, 0x61, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x64,
	0x65, 0x70, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x64, 0x65, 0x70, 0x74,
	0x68, 0x22, 0x67, 0x0a, 0x15, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x61, 0x74, 0x65, 0x67,
	0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61,
	0x72, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x70,
	0x61, 0x72, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73,
	0x6f, 0x72, 0x74, 0x5f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x09, 0x73, 0x6f, 0x72, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x22, 0x24, 0x0a, 0x12, 0x47, 0x65,
	0x74, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64,
	0x22, 0x3b, 0x0a, 0x15, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f,
	0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x27, 0x0a,
	0x15, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x22, 0x2c, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x53, 0x75, 0x62,
	0x54, 0x72, 0x65, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x72,
	0x6f, 0x6f, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x72, 0x6f,
	0x6f, 0x74, 0x49, 0x64, 0x22, 0x39, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x43, 0x61, 0x74, 0x65, 0x67,
	0x6f, 0x72, 0x79, 0x50, 0x61, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f,
	0x0a, 0x0b, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0a, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x49, 0x64, 0x22,
	0x34, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x63, 0x61, 0x74, 0x65, 0x67,
	0x6f, 0x72, 0x79, 0x49, 0x64, 0x22, 0x52, 0x0a, 0x19, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43,
	0x6c, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x44, 0x65, 0x70, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72,
	0x79, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x32, 0xce, 0x08, 0x0a, 0x0f, 0x43, 0x61,
	0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x6e, 0x0a,
	0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12,
	0x26, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x2e, 0x76,
	0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x63, 0x61,
	0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f,
	0x72, 0x79, 0x22, 0x19, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x13, 0x3a, 0x01, 0x2a, 0x22, 0x0e, 0x2f,
	0x76, 0x31, 0x2f, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x12, 0x6a, 0x0a,
	0x0b, 0x47, 0x65, 0x74, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x23, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x47,
	0x65, 0x74, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x19, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79,
	0x2e, 0x76, 0x31, 0x2e, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x22, 0x1b, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x15, 0x12, 0x13, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f,
	0x72, 0x69, 0x65, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x70, 0x0a, 0x0e, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x26, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x1e, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x18, 0x3a, 0x01, 0x2a, 0x1a, 0x13, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x61, 0x74, 0x65,
	0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x6d, 0x0a, 0x0e, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x26, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x1b, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x15, 0x2a, 0x13, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x61, 0x74, 0x65, 0x67,
	0x6f, 0x72, 0x69, 0x65, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x77, 0x0a, 0x0a, 0x47, 0x65,
	0x74, 0x53, 0x75, 0x62, 0x54, 0x72, 0x65, 0x65, 0x12, 0x22, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x63,
	0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x75,
	0x62, 0x54, 0x72, 0x65, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x43,
	0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x22, 0x28, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x22, 0x12,
	0x20, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x2f,
	0x7b, 0x72, 0x6f, 0x6f, 0x74, 0x5f, 0x69, 0x64, 0x7d, 0x2f, 0x73, 0x75, 0x62, 0x74, 0x72, 0x65,
	0x65, 0x30, 0x01, 0x12, 0x82, 0x01, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x43, 0x61, 0x74, 0x65, 0x67,
	0x6f, 0x72, 0x79, 0x50, 0x61, 0x74, 0x68, 0x12, 0x27, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x63, 0x61,
	0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x61, 0x74,
	0x65, 0x67, 0x6f, 0x72, 0x79, 0x50, 0x61, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x19, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x2e,
	0x76, 0x31, 0x2e, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x22, 0x29, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x23, 0x12, 0x21, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72,
	0x69, 0x65, 0x73, 0x2f, 0x7b, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x5f, 0x69, 0x64,
	0x7d, 0x2f, 0x70, 0x61, 0x74, 0x68, 0x30, 0x01, 0x12, 0x66, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x4c,
	0x65, 0x61, 0x66, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x12, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x63, 0x61, 0x74, 0x65,
	0x67, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79,
	0x73, 0x22, 0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17, 0x12, 0x15, 0x2f, 0x76, 0x31, 0x2f, 0x63,
	0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x2f, 0x6c, 0x65, 0x61, 0x76, 0x65, 0x73,
	0x12, 0x8b, 0x01, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x52,
	0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x22, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x63,
	0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6c,
	0x6f, 0x73, 0x75, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x43,
	0x6c, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x2c,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x26, 0x12, 0x24, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x61, 0x74, 0x65,
	0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x2f, 0x7b, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79,
	0x5f, 0x69, 0x64, 0x7d, 0x2f, 0x63, 0x6c, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x30, 0x01, 0x12, 0x89,
	0x01, 0x0a, 0x12, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x6f, 0x73, 0x75, 0x72, 0x65,
	0x44, 0x65, 0x70, 0x74, 0x68, 0x12, 0x2a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x63, 0x61, 0x74, 0x65,
	0x67, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6c,
	0x6f, 0x73, 0x75, 0x72, 0x65, 0x44, 0x65, 0x70, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x2f, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x29, 0x3a, 0x01, 0x2a, 0x32, 0x24, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f,
	0x72, 0x69, 0x65, 0x73, 0x2f, 0x7b, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x5f, 0x69,
	0x64, 0x7d, 0x2f, 0x63, 0x6c, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x42, 0x1a, 0x5a, 0x18, 0x61, 0x70,
	0x69, 0x2f, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x2f, 0x76, 0x31, 0x3b, 0x63, 0x61,
	0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_v1_category_proto_rawDescOnce sync.Once
	file_v1_category_proto_rawDescData []byte
)

func file_v1_category_proto_rawDescGZIP() []byte {
	file_v1_category_proto_rawDescOnce.Do(func() {
		file_v1_category_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_v1_category_proto_rawDesc), len(file_v1_category_proto_rawDesc)))
	})
	return file_v1_category_proto_rawDescData
}

var file_v1_category_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_v1_category_proto_goTypes = []any{
	(*Categorys)(nil),                 // 0: api.category.v1.Categorys
	(*Category)(nil),                  // 1: api.category.v1.Category
	(*ClosureRelation)(nil),           // 2: api.category.v1.ClosureRelation
	(*CreateCategoryRequest)(nil),     // 3: api.category.v1.CreateCategoryRequest
	(*GetCategoryRequest)(nil),        // 4: api.category.v1.GetCategoryRequest
	(*UpdateCategoryRequest)(nil),     // 5: api.category.v1.UpdateCategoryRequest
	(*DeleteCategoryRequest)(nil),     // 6: api.category.v1.DeleteCategoryRequest
	(*GetSubTreeRequest)(nil),         // 7: api.category.v1.GetSubTreeRequest
	(*GetCategoryPathRequest)(nil),    // 8: api.category.v1.GetCategoryPathRequest
	(*GetClosureRequest)(nil),         // 9: api.category.v1.GetClosureRequest
	(*UpdateClosureDepthRequest)(nil), // 10: api.category.v1.UpdateClosureDepthRequest
	(*timestamppb.Timestamp)(nil),     // 11: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),             // 12: google.protobuf.Empty
}
var file_v1_category_proto_depIdxs = []int32{
	1,  // 0: api.category.v1.Categorys.categorys:type_name -> api.category.v1.Category
	11, // 1: api.category.v1.Category.created_at:type_name -> google.protobuf.Timestamp
	11, // 2: api.category.v1.Category.updated_at:type_name -> google.protobuf.Timestamp
	3,  // 3: api.category.v1.CategoryService.CreateCategory:input_type -> api.category.v1.CreateCategoryRequest
	4,  // 4: api.category.v1.CategoryService.GetCategory:input_type -> api.category.v1.GetCategoryRequest
	5,  // 5: api.category.v1.CategoryService.UpdateCategory:input_type -> api.category.v1.UpdateCategoryRequest
	6,  // 6: api.category.v1.CategoryService.DeleteCategory:input_type -> api.category.v1.DeleteCategoryRequest
	7,  // 7: api.category.v1.CategoryService.GetSubTree:input_type -> api.category.v1.GetSubTreeRequest
	8,  // 8: api.category.v1.CategoryService.GetCategoryPath:input_type -> api.category.v1.GetCategoryPathRequest
	12, // 9: api.category.v1.CategoryService.GetLeafCategories:input_type -> google.protobuf.Empty
	9,  // 10: api.category.v1.CategoryService.GetClosureRelations:input_type -> api.category.v1.GetClosureRequest
	10, // 11: api.category.v1.CategoryService.UpdateClosureDepth:input_type -> api.category.v1.UpdateClosureDepthRequest
	1,  // 12: api.category.v1.CategoryService.CreateCategory:output_type -> api.category.v1.Category
	1,  // 13: api.category.v1.CategoryService.GetCategory:output_type -> api.category.v1.Category
	12, // 14: api.category.v1.CategoryService.UpdateCategory:output_type -> google.protobuf.Empty
	12, // 15: api.category.v1.CategoryService.DeleteCategory:output_type -> google.protobuf.Empty
	1,  // 16: api.category.v1.CategoryService.GetSubTree:output_type -> api.category.v1.Category
	1,  // 17: api.category.v1.CategoryService.GetCategoryPath:output_type -> api.category.v1.Category
	0,  // 18: api.category.v1.CategoryService.GetLeafCategories:output_type -> api.category.v1.Categorys
	2,  // 19: api.category.v1.CategoryService.GetClosureRelations:output_type -> api.category.v1.ClosureRelation
	12, // 20: api.category.v1.CategoryService.UpdateClosureDepth:output_type -> google.protobuf.Empty
	12, // [12:21] is the sub-list for method output_type
	3,  // [3:12] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_v1_category_proto_init() }
func file_v1_category_proto_init() {
	if File_v1_category_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_v1_category_proto_rawDesc), len(file_v1_category_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_v1_category_proto_goTypes,
		DependencyIndexes: file_v1_category_proto_depIdxs,
		MessageInfos:      file_v1_category_proto_msgTypes,
	}.Build()
	File_v1_category_proto = out.File
	file_v1_category_proto_goTypes = nil
	file_v1_category_proto_depIdxs = nil
}
