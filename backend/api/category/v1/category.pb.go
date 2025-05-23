// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: v1/category.proto

// 定义包名，用于区分不同的服务模块。

package categoryv1

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

// 批量查询分类请求
type BatchGetCategoriesRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Ids           []int64                `protobuf:"varint,1,rep,packed,name=ids,proto3" json:"ids,omitempty"` // 注意类型与SQL中的bigint匹配
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BatchGetCategoriesRequest) Reset() {
	*x = BatchGetCategoriesRequest{}
	mi := &file_v1_category_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BatchGetCategoriesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchGetCategoriesRequest) ProtoMessage() {}

func (x *BatchGetCategoriesRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use BatchGetCategoriesRequest.ProtoReflect.Descriptor instead.
func (*BatchGetCategoriesRequest) Descriptor() ([]byte, []int) {
	return file_v1_category_proto_rawDescGZIP(), []int{0}
}

func (x *BatchGetCategoriesRequest) GetIds() []int64 {
	if x != nil {
		return x.Ids
	}
	return nil
}

type Categories struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Categories    []*Category            `protobuf:"bytes,1,rep,name=categories,proto3" json:"categories,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Categories) Reset() {
	*x = Categories{}
	mi := &file_v1_category_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Categories) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Categories) ProtoMessage() {}

func (x *Categories) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Categories.ProtoReflect.Descriptor instead.
func (*Categories) Descriptor() ([]byte, []int) {
	return file_v1_category_proto_rawDescGZIP(), []int{1}
}

func (x *Categories) GetCategories() []*Category {
	if x != nil {
		return x.Categories
	}
	return nil
}

type Category struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	ParentId      int64                  `protobuf:"varint,2,opt,name=parent_id,json=parentId,proto3" json:"parent_id,omitempty"` // 0 represents root in proto (stored as NULL in DB)
	Level         int32                  `protobuf:"varint,3,opt,name=level,proto3" json:"level,omitempty"`                       // 0-3
	Path          string                 `protobuf:"bytes,4,opt,name=path,proto3" json:"path,omitempty"`                          // ltree path
	Name          string                 `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	SortOrder     int32                  `protobuf:"varint,6,opt,name=sort_order,json=sortOrder,proto3" json:"sort_order,omitempty"`
	IsLeaf        bool                   `protobuf:"varint,7,opt,name=is_leaf,json=isLeaf,proto3" json:"is_leaf,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Category) Reset() {
	*x = Category{}
	mi := &file_v1_category_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Category) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Category) ProtoMessage() {}

func (x *Category) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Category.ProtoReflect.Descriptor instead.
func (*Category) Descriptor() ([]byte, []int) {
	return file_v1_category_proto_rawDescGZIP(), []int{2}
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

type ClosureRelations struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Relations     []*ClosureRelation     `protobuf:"bytes,1,rep,name=relations,proto3" json:"relations,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ClosureRelations) Reset() {
	*x = ClosureRelations{}
	mi := &file_v1_category_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ClosureRelations) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClosureRelations) ProtoMessage() {}

func (x *ClosureRelations) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use ClosureRelations.ProtoReflect.Descriptor instead.
func (*ClosureRelations) Descriptor() ([]byte, []int) {
	return file_v1_category_proto_rawDescGZIP(), []int{3}
}

func (x *ClosureRelations) GetRelations() []*ClosureRelation {
	if x != nil {
		return x.Relations
	}
	return nil
}

type ClosureRelation struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Ancestor      int64                  `protobuf:"varint,1,opt,name=ancestor,proto3" json:"ancestor,omitempty"`
	Descendant    int64                  `protobuf:"varint,2,opt,name=descendant,proto3" json:"descendant,omitempty"`
	Depth         int32                  `protobuf:"varint,3,opt,name=depth,proto3" json:"depth,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ClosureRelation) Reset() {
	*x = ClosureRelation{}
	mi := &file_v1_category_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ClosureRelation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClosureRelation) ProtoMessage() {}

func (x *ClosureRelation) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use ClosureRelation.ProtoReflect.Descriptor instead.
func (*ClosureRelation) Descriptor() ([]byte, []int) {
	return file_v1_category_proto_rawDescGZIP(), []int{4}
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
	mi := &file_v1_category_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateCategoryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateCategoryRequest) ProtoMessage() {}

func (x *CreateCategoryRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use CreateCategoryRequest.ProtoReflect.Descriptor instead.
func (*CreateCategoryRequest) Descriptor() ([]byte, []int) {
	return file_v1_category_proto_rawDescGZIP(), []int{5}
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
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"` // 分类的唯一标识符。
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetCategoryRequest) Reset() {
	*x = GetCategoryRequest{}
	mi := &file_v1_category_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetCategoryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCategoryRequest) ProtoMessage() {}

func (x *GetCategoryRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use GetCategoryRequest.ProtoReflect.Descriptor instead.
func (*GetCategoryRequest) Descriptor() ([]byte, []int) {
	return file_v1_category_proto_rawDescGZIP(), []int{6}
}

func (x *GetCategoryRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

// 更新分类请求
type UpdateCategoryRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`    // 分类的唯一标识符。
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"` // 新的分类名称。
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateCategoryRequest) Reset() {
	*x = UpdateCategoryRequest{}
	mi := &file_v1_category_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateCategoryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateCategoryRequest) ProtoMessage() {}

func (x *UpdateCategoryRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use UpdateCategoryRequest.ProtoReflect.Descriptor instead.
func (*UpdateCategoryRequest) Descriptor() ([]byte, []int) {
	return file_v1_category_proto_rawDescGZIP(), []int{7}
}

func (x *UpdateCategoryRequest) GetId() uint64 {
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
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"` // 分类的唯一标识符。
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteCategoryRequest) Reset() {
	*x = DeleteCategoryRequest{}
	mi := &file_v1_category_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteCategoryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteCategoryRequest) ProtoMessage() {}

func (x *DeleteCategoryRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use DeleteCategoryRequest.ProtoReflect.Descriptor instead.
func (*DeleteCategoryRequest) Descriptor() ([]byte, []int) {
	return file_v1_category_proto_rawDescGZIP(), []int{8}
}

func (x *DeleteCategoryRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

// 获取子树请求
type GetSubTreeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RootId        uint64                 `protobuf:"varint,1,opt,name=root_id,json=rootId,proto3" json:"root_id,omitempty"` // 根节点的 ID，用于获取其子树。
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetSubTreeRequest) Reset() {
	*x = GetSubTreeRequest{}
	mi := &file_v1_category_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetSubTreeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSubTreeRequest) ProtoMessage() {}

func (x *GetSubTreeRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use GetSubTreeRequest.ProtoReflect.Descriptor instead.
func (*GetSubTreeRequest) Descriptor() ([]byte, []int) {
	return file_v1_category_proto_rawDescGZIP(), []int{9}
}

func (x *GetSubTreeRequest) GetRootId() uint64 {
	if x != nil {
		return x.RootId
	}
	return 0
}

// 获取分类路径请求
type GetCategoryPathRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CategoryId    uint64                 `protobuf:"varint,1,opt,name=category_id,json=categoryId,proto3" json:"category_id,omitempty"` // 分类的唯一标识符，用于获取其完整路径。
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetCategoryPathRequest) Reset() {
	*x = GetCategoryPathRequest{}
	mi := &file_v1_category_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetCategoryPathRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCategoryPathRequest) ProtoMessage() {}

func (x *GetCategoryPathRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use GetCategoryPathRequest.ProtoReflect.Descriptor instead.
func (*GetCategoryPathRequest) Descriptor() ([]byte, []int) {
	return file_v1_category_proto_rawDescGZIP(), []int{10}
}

func (x *GetCategoryPathRequest) GetCategoryId() uint64 {
	if x != nil {
		return x.CategoryId
	}
	return 0
}

// 获取直接子分类请求
type GetDirectSubCategoriesRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ParentId      uint64                 `protobuf:"varint,1,opt,name=parent_id,json=parentId,proto3" json:"parent_id,omitempty"` // 父分类的唯一标识符，用于获取其直接子分类。
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetDirectSubCategoriesRequest) Reset() {
	*x = GetDirectSubCategoriesRequest{}
	mi := &file_v1_category_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetDirectSubCategoriesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDirectSubCategoriesRequest) ProtoMessage() {}

func (x *GetDirectSubCategoriesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_category_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDirectSubCategoriesRequest.ProtoReflect.Descriptor instead.
func (*GetDirectSubCategoriesRequest) Descriptor() ([]byte, []int) {
	return file_v1_category_proto_rawDescGZIP(), []int{11}
}

func (x *GetDirectSubCategoriesRequest) GetParentId() uint64 {
	if x != nil {
		return x.ParentId
	}
	return 0
}

// 获取闭包关系请求
type GetClosureRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CategoryId    uint64                 `protobuf:"varint,1,opt,name=category_id,json=categoryId,proto3" json:"category_id,omitempty"` // 分类的唯一标识符，用于获取其闭包关系。
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetClosureRequest) Reset() {
	*x = GetClosureRequest{}
	mi := &file_v1_category_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetClosureRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetClosureRequest) ProtoMessage() {}

func (x *GetClosureRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_category_proto_msgTypes[12]
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
	return file_v1_category_proto_rawDescGZIP(), []int{12}
}

func (x *GetClosureRequest) GetCategoryId() uint64 {
	if x != nil {
		return x.CategoryId
	}
	return 0
}

// 更新闭包关系深度请求
type UpdateClosureDepthRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CategoryId    int64                  `protobuf:"varint,1,opt,name=category_id,json=categoryId,proto3" json:"category_id,omitempty"` // 分类的唯一标识符。
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateClosureDepthRequest) Reset() {
	*x = UpdateClosureDepthRequest{}
	mi := &file_v1_category_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateClosureDepthRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateClosureDepthRequest) ProtoMessage() {}

func (x *UpdateClosureDepthRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_category_proto_msgTypes[13]
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
	return file_v1_category_proto_rawDescGZIP(), []int{13}
}

func (x *UpdateClosureDepthRequest) GetCategoryId() int64 {
	if x != nil {
		return x.CategoryId
	}
	return 0
}

func (x *UpdateClosureDepthRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_v1_category_proto protoreflect.FileDescriptor

const file_v1_category_proto_rawDesc = "" +
	"\n" +
	"\x11v1/category.proto\x12\x15ecommerce.category.v1\x1a\x1cgoogle/api/annotations.proto\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x1bgoogle/protobuf/empty.proto\"-\n" +
	"\x19BatchGetCategoriesRequest\x12\x10\n" +
	"\x03ids\x18\x01 \x03(\x03R\x03ids\"M\n" +
	"\n" +
	"Categories\x12?\n" +
	"\n" +
	"categories\x18\x01 \x03(\v2\x1f.ecommerce.category.v1.CategoryR\n" +
	"categories\"\xa3\x02\n" +
	"\bCategory\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x1b\n" +
	"\tparent_id\x18\x02 \x01(\x03R\bparentId\x12\x14\n" +
	"\x05level\x18\x03 \x01(\x05R\x05level\x12\x12\n" +
	"\x04path\x18\x04 \x01(\tR\x04path\x12\x12\n" +
	"\x04name\x18\x05 \x01(\tR\x04name\x12\x1d\n" +
	"\n" +
	"sort_order\x18\x06 \x01(\x05R\tsortOrder\x12\x17\n" +
	"\ais_leaf\x18\a \x01(\bR\x06isLeaf\x129\n" +
	"\n" +
	"created_at\x18\b \x01(\v2\x1a.google.protobuf.TimestampR\tcreatedAt\x129\n" +
	"\n" +
	"updated_at\x18\t \x01(\v2\x1a.google.protobuf.TimestampR\tupdatedAt\"X\n" +
	"\x10ClosureRelations\x12D\n" +
	"\trelations\x18\x01 \x03(\v2&.ecommerce.category.v1.ClosureRelationR\trelations\"c\n" +
	"\x0fClosureRelation\x12\x1a\n" +
	"\bancestor\x18\x01 \x01(\x03R\bancestor\x12\x1e\n" +
	"\n" +
	"descendant\x18\x02 \x01(\x03R\n" +
	"descendant\x12\x14\n" +
	"\x05depth\x18\x03 \x01(\x05R\x05depth\"g\n" +
	"\x15CreateCategoryRequest\x12\x1b\n" +
	"\tparent_id\x18\x01 \x01(\x03R\bparentId\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12\x1d\n" +
	"\n" +
	"sort_order\x18\x03 \x01(\x05R\tsortOrder\"$\n" +
	"\x12GetCategoryRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\";\n" +
	"\x15UpdateCategoryRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\"'\n" +
	"\x15DeleteCategoryRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\",\n" +
	"\x11GetSubTreeRequest\x12\x17\n" +
	"\aroot_id\x18\x01 \x01(\x04R\x06rootId\"9\n" +
	"\x16GetCategoryPathRequest\x12\x1f\n" +
	"\vcategory_id\x18\x01 \x01(\x04R\n" +
	"categoryId\"<\n" +
	"\x1dGetDirectSubCategoriesRequest\x12\x1b\n" +
	"\tparent_id\x18\x01 \x01(\x04R\bparentId\"4\n" +
	"\x11GetClosureRequest\x12\x1f\n" +
	"\vcategory_id\x18\x01 \x01(\x04R\n" +
	"categoryId\"P\n" +
	"\x19UpdateClosureDepthRequest\x12\x1f\n" +
	"\vcategory_id\x18\x01 \x01(\x03R\n" +
	"categoryId\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name2\xce\v\n" +
	"\x0fCategoryService\x12z\n" +
	"\x0eCreateCategory\x12,.ecommerce.category.v1.CreateCategoryRequest\x1a\x1f.ecommerce.category.v1.Category\"\x19\x82\xd3\xe4\x93\x02\x13:\x01*\"\x0e/v1/categories\x12m\n" +
	"\x11GetLeafCategories\x12\x16.google.protobuf.Empty\x1a!.ecommerce.category.v1.Categories\"\x1d\x82\xd3\xe4\x93\x02\x17\x12\x15/v1/categories/leaves\x12\x87\x01\n" +
	"\x12BatchGetCategories\x120.ecommerce.category.v1.BatchGetCategoriesRequest\x1a!.ecommerce.category.v1.Categories\"\x1c\x82\xd3\xe4\x93\x02\x16\x12\x14/v1/categories/batch\x12v\n" +
	"\vGetCategory\x12).ecommerce.category.v1.GetCategoryRequest\x1a\x1f.ecommerce.category.v1.Category\"\x1b\x82\xd3\xe4\x93\x02\x15\x12\x13/v1/categories/{id}\x12v\n" +
	"\x0eUpdateCategory\x12,.ecommerce.category.v1.UpdateCategoryRequest\x1a\x16.google.protobuf.Empty\"\x1e\x82\xd3\xe4\x93\x02\x18:\x01*\x1a\x13/v1/categories/{id}\x12s\n" +
	"\x0eDeleteCategory\x12,.ecommerce.category.v1.DeleteCategoryRequest\x1a\x16.google.protobuf.Empty\"\x1b\x82\xd3\xe4\x93\x02\x15*\x13/v1/categories/{id}\x12\x83\x01\n" +
	"\n" +
	"GetSubTree\x12(.ecommerce.category.v1.GetSubTreeRequest\x1a!.ecommerce.category.v1.Categories\"(\x82\xd3\xe4\x93\x02\"\x12 /v1/categories/{root_id}/subtree\x12\x9e\x01\n" +
	"\x16GetDirectSubCategories\x124.ecommerce.category.v1.GetDirectSubCategoriesRequest\x1a!.ecommerce.category.v1.Categories\"+\x82\xd3\xe4\x93\x02%\x12#/v1/categories/{parent_id}/children\x12\x8e\x01\n" +
	"\x0fGetCategoryPath\x12-.ecommerce.category.v1.GetCategoryPathRequest\x1a!.ecommerce.category.v1.Categories\")\x82\xd3\xe4\x93\x02#\x12!/v1/categories/{category_id}/path\x12\x96\x01\n" +
	"\x13GetClosureRelations\x12(.ecommerce.category.v1.GetClosureRequest\x1a'.ecommerce.category.v1.ClosureRelations\",\x82\xd3\xe4\x93\x02&\x12$/v1/categories/{category_id}/closure\x12\x8f\x01\n" +
	"\x12UpdateClosureDepth\x120.ecommerce.category.v1.UpdateClosureDepthRequest\x1a\x16.google.protobuf.Empty\"/\x82\xd3\xe4\x93\x02):\x01*2$/v1/categories/{category_id}/closureB$Z\"backend/api/category/v1;categoryv1b\x06proto3"

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

var file_v1_category_proto_msgTypes = make([]protoimpl.MessageInfo, 14)
var file_v1_category_proto_goTypes = []any{
	(*BatchGetCategoriesRequest)(nil),     // 0: ecommerce.category.v1.BatchGetCategoriesRequest
	(*Categories)(nil),                    // 1: ecommerce.category.v1.Categories
	(*Category)(nil),                      // 2: ecommerce.category.v1.Category
	(*ClosureRelations)(nil),              // 3: ecommerce.category.v1.ClosureRelations
	(*ClosureRelation)(nil),               // 4: ecommerce.category.v1.ClosureRelation
	(*CreateCategoryRequest)(nil),         // 5: ecommerce.category.v1.CreateCategoryRequest
	(*GetCategoryRequest)(nil),            // 6: ecommerce.category.v1.GetCategoryRequest
	(*UpdateCategoryRequest)(nil),         // 7: ecommerce.category.v1.UpdateCategoryRequest
	(*DeleteCategoryRequest)(nil),         // 8: ecommerce.category.v1.DeleteCategoryRequest
	(*GetSubTreeRequest)(nil),             // 9: ecommerce.category.v1.GetSubTreeRequest
	(*GetCategoryPathRequest)(nil),        // 10: ecommerce.category.v1.GetCategoryPathRequest
	(*GetDirectSubCategoriesRequest)(nil), // 11: ecommerce.category.v1.GetDirectSubCategoriesRequest
	(*GetClosureRequest)(nil),             // 12: ecommerce.category.v1.GetClosureRequest
	(*UpdateClosureDepthRequest)(nil),     // 13: ecommerce.category.v1.UpdateClosureDepthRequest
	(*timestamppb.Timestamp)(nil),         // 14: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),                 // 15: google.protobuf.Empty
}
var file_v1_category_proto_depIdxs = []int32{
	2,  // 0: ecommerce.category.v1.Categories.categories:type_name -> ecommerce.category.v1.Category
	14, // 1: ecommerce.category.v1.Category.created_at:type_name -> google.protobuf.Timestamp
	14, // 2: ecommerce.category.v1.Category.updated_at:type_name -> google.protobuf.Timestamp
	4,  // 3: ecommerce.category.v1.ClosureRelations.relations:type_name -> ecommerce.category.v1.ClosureRelation
	5,  // 4: ecommerce.category.v1.CategoryService.CreateCategory:input_type -> ecommerce.category.v1.CreateCategoryRequest
	15, // 5: ecommerce.category.v1.CategoryService.GetLeafCategories:input_type -> google.protobuf.Empty
	0,  // 6: ecommerce.category.v1.CategoryService.BatchGetCategories:input_type -> ecommerce.category.v1.BatchGetCategoriesRequest
	6,  // 7: ecommerce.category.v1.CategoryService.GetCategory:input_type -> ecommerce.category.v1.GetCategoryRequest
	7,  // 8: ecommerce.category.v1.CategoryService.UpdateCategory:input_type -> ecommerce.category.v1.UpdateCategoryRequest
	8,  // 9: ecommerce.category.v1.CategoryService.DeleteCategory:input_type -> ecommerce.category.v1.DeleteCategoryRequest
	9,  // 10: ecommerce.category.v1.CategoryService.GetSubTree:input_type -> ecommerce.category.v1.GetSubTreeRequest
	11, // 11: ecommerce.category.v1.CategoryService.GetDirectSubCategories:input_type -> ecommerce.category.v1.GetDirectSubCategoriesRequest
	10, // 12: ecommerce.category.v1.CategoryService.GetCategoryPath:input_type -> ecommerce.category.v1.GetCategoryPathRequest
	12, // 13: ecommerce.category.v1.CategoryService.GetClosureRelations:input_type -> ecommerce.category.v1.GetClosureRequest
	13, // 14: ecommerce.category.v1.CategoryService.UpdateClosureDepth:input_type -> ecommerce.category.v1.UpdateClosureDepthRequest
	2,  // 15: ecommerce.category.v1.CategoryService.CreateCategory:output_type -> ecommerce.category.v1.Category
	1,  // 16: ecommerce.category.v1.CategoryService.GetLeafCategories:output_type -> ecommerce.category.v1.Categories
	1,  // 17: ecommerce.category.v1.CategoryService.BatchGetCategories:output_type -> ecommerce.category.v1.Categories
	2,  // 18: ecommerce.category.v1.CategoryService.GetCategory:output_type -> ecommerce.category.v1.Category
	15, // 19: ecommerce.category.v1.CategoryService.UpdateCategory:output_type -> google.protobuf.Empty
	15, // 20: ecommerce.category.v1.CategoryService.DeleteCategory:output_type -> google.protobuf.Empty
	1,  // 21: ecommerce.category.v1.CategoryService.GetSubTree:output_type -> ecommerce.category.v1.Categories
	1,  // 22: ecommerce.category.v1.CategoryService.GetDirectSubCategories:output_type -> ecommerce.category.v1.Categories
	1,  // 23: ecommerce.category.v1.CategoryService.GetCategoryPath:output_type -> ecommerce.category.v1.Categories
	3,  // 24: ecommerce.category.v1.CategoryService.GetClosureRelations:output_type -> ecommerce.category.v1.ClosureRelations
	15, // 25: ecommerce.category.v1.CategoryService.UpdateClosureDepth:output_type -> google.protobuf.Empty
	15, // [15:26] is the sub-list for method output_type
	4,  // [4:15] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
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
			NumMessages:   14,
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
