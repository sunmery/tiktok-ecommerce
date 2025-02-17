// Code generated by protoc-gen-go-errors. DO NOT EDIT.

package category

import (
	fmt "fmt"
	errors "github.com/go-kratos/kratos/v2/errors"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
const _ = errors.SupportPackageIsVersion1

// PARENT_ID不符合业务规则
func IsParentIdUnprocessableEntity(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_PARENT_ID_UNPROCESSABLE_ENTITY.String() && e.Code == 422
}

// PARENT_ID不符合业务规则
func ErrorParentIdUnprocessableEntity(format string, args ...interface{}) *errors.Error {
	return errors.New(422, ErrorReason_PARENT_ID_UNPROCESSABLE_ENTITY.String(), fmt.Sprintf(format, args...))
}

// 找不到该分类名称
func IsCategoryNameNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_CATEGORY_NAME_NOT_FOUND.String() && e.Code == 404
}

// 找不到该分类名称
func ErrorCategoryNameNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(404, ErrorReason_CATEGORY_NAME_NOT_FOUND.String(), fmt.Sprintf(format, args...))
}

// 找不到该分类
func IsCategoryNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_CATEGORY_NOT_FOUND.String() && e.Code == 404
}

// 找不到该分类
func ErrorCategoryNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(404, ErrorReason_CATEGORY_NOT_FOUND.String(), fmt.Sprintf(format, args...))
}
