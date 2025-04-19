// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: v1/comment.proto

package commentv1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on CommentType with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *CommentType) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CommentType with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in CommentTypeMultiError, or
// nil if none found.
func (m *CommentType) ValidateAll() error {
	return m.validate(true)
}

func (m *CommentType) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for ProductId

	// no validation rules for MerchantId

	// no validation rules for UserId

	// no validation rules for Score

	// no validation rules for Content

	if all {
		switch v := interface{}(m.GetCreatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CommentTypeValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CommentTypeValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CommentTypeValidationError{
				field:  "CreatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetUpdatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CommentTypeValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CommentTypeValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetUpdatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CommentTypeValidationError{
				field:  "UpdatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return CommentTypeMultiError(errors)
	}

	return nil
}

// CommentTypeMultiError is an error wrapping multiple validation errors
// returned by CommentType.ValidateAll() if the designated constraints aren't met.
type CommentTypeMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CommentTypeMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CommentTypeMultiError) AllErrors() []error { return m }

// CommentTypeValidationError is the validation error returned by
// CommentType.Validate if the designated constraints aren't met.
type CommentTypeValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CommentTypeValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CommentTypeValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CommentTypeValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CommentTypeValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CommentTypeValidationError) ErrorName() string { return "CommentTypeValidationError" }

// Error satisfies the builtin error interface
func (e CommentTypeValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCommentType.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CommentTypeValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CommentTypeValidationError{}

// Validate checks the field values on CreateCommentRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateCommentRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateCommentRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateCommentRequestMultiError, or nil if none found.
func (m *CreateCommentRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateCommentRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for ProductId

	// no validation rules for MerchantId

	// no validation rules for UserId

	// no validation rules for Score

	// no validation rules for Content

	if len(errors) > 0 {
		return CreateCommentRequestMultiError(errors)
	}

	return nil
}

// CreateCommentRequestMultiError is an error wrapping multiple validation
// errors returned by CreateCommentRequest.ValidateAll() if the designated
// constraints aren't met.
type CreateCommentRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateCommentRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateCommentRequestMultiError) AllErrors() []error { return m }

// CreateCommentRequestValidationError is the validation error returned by
// CreateCommentRequest.Validate if the designated constraints aren't met.
type CreateCommentRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateCommentRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateCommentRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateCommentRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateCommentRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateCommentRequestValidationError) ErrorName() string {
	return "CreateCommentRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreateCommentRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateCommentRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateCommentRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateCommentRequestValidationError{}

// Validate checks the field values on GetCommentsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetCommentsRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetCommentsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetCommentsRequestMultiError, or nil if none found.
func (m *GetCommentsRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetCommentsRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for ProductId

	// no validation rules for MerchantId

	// no validation rules for Page

	// no validation rules for PageSize

	if len(errors) > 0 {
		return GetCommentsRequestMultiError(errors)
	}

	return nil
}

// GetCommentsRequestMultiError is an error wrapping multiple validation errors
// returned by GetCommentsRequest.ValidateAll() if the designated constraints
// aren't met.
type GetCommentsRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetCommentsRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetCommentsRequestMultiError) AllErrors() []error { return m }

// GetCommentsRequestValidationError is the validation error returned by
// GetCommentsRequest.Validate if the designated constraints aren't met.
type GetCommentsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetCommentsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetCommentsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetCommentsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetCommentsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetCommentsRequestValidationError) ErrorName() string {
	return "GetCommentsRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetCommentsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetCommentsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetCommentsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetCommentsRequestValidationError{}

// Validate checks the field values on GetCommentsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetCommentsResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetCommentsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetCommentsResponseMultiError, or nil if none found.
func (m *GetCommentsResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetCommentsResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetComments() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, GetCommentsResponseValidationError{
						field:  fmt.Sprintf("Comments[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, GetCommentsResponseValidationError{
						field:  fmt.Sprintf("Comments[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GetCommentsResponseValidationError{
					field:  fmt.Sprintf("Comments[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for Total

	if len(errors) > 0 {
		return GetCommentsResponseMultiError(errors)
	}

	return nil
}

// GetCommentsResponseMultiError is an error wrapping multiple validation
// errors returned by GetCommentsResponse.ValidateAll() if the designated
// constraints aren't met.
type GetCommentsResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetCommentsResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetCommentsResponseMultiError) AllErrors() []error { return m }

// GetCommentsResponseValidationError is the validation error returned by
// GetCommentsResponse.Validate if the designated constraints aren't met.
type GetCommentsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetCommentsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetCommentsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetCommentsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetCommentsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetCommentsResponseValidationError) ErrorName() string {
	return "GetCommentsResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetCommentsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetCommentsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetCommentsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetCommentsResponseValidationError{}

// Validate checks the field values on UpdateCommentRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateCommentRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateCommentRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UpdateCommentRequestMultiError, or nil if none found.
func (m *UpdateCommentRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateCommentRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for CommentId

	// no validation rules for UserId

	// no validation rules for Score

	// no validation rules for Content

	if len(errors) > 0 {
		return UpdateCommentRequestMultiError(errors)
	}

	return nil
}

// UpdateCommentRequestMultiError is an error wrapping multiple validation
// errors returned by UpdateCommentRequest.ValidateAll() if the designated
// constraints aren't met.
type UpdateCommentRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateCommentRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateCommentRequestMultiError) AllErrors() []error { return m }

// UpdateCommentRequestValidationError is the validation error returned by
// UpdateCommentRequest.Validate if the designated constraints aren't met.
type UpdateCommentRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateCommentRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateCommentRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateCommentRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateCommentRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateCommentRequestValidationError) ErrorName() string {
	return "UpdateCommentRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateCommentRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateCommentRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateCommentRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateCommentRequestValidationError{}

// Validate checks the field values on DeleteCommentRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DeleteCommentRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteCommentRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeleteCommentRequestMultiError, or nil if none found.
func (m *DeleteCommentRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteCommentRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for CommentId

	// no validation rules for UserId

	if len(errors) > 0 {
		return DeleteCommentRequestMultiError(errors)
	}

	return nil
}

// DeleteCommentRequestMultiError is an error wrapping multiple validation
// errors returned by DeleteCommentRequest.ValidateAll() if the designated
// constraints aren't met.
type DeleteCommentRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteCommentRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteCommentRequestMultiError) AllErrors() []error { return m }

// DeleteCommentRequestValidationError is the validation error returned by
// DeleteCommentRequest.Validate if the designated constraints aren't met.
type DeleteCommentRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteCommentRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteCommentRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteCommentRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteCommentRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteCommentRequestValidationError) ErrorName() string {
	return "DeleteCommentRequestValidationError"
}

// Error satisfies the builtin error interface
func (e DeleteCommentRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteCommentRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteCommentRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteCommentRequestValidationError{}

// Validate checks the field values on DeleteCommentResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DeleteCommentResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteCommentResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeleteCommentResponseMultiError, or nil if none found.
func (m *DeleteCommentResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteCommentResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Success

	if len(errors) > 0 {
		return DeleteCommentResponseMultiError(errors)
	}

	return nil
}

// DeleteCommentResponseMultiError is an error wrapping multiple validation
// errors returned by DeleteCommentResponse.ValidateAll() if the designated
// constraints aren't met.
type DeleteCommentResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteCommentResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteCommentResponseMultiError) AllErrors() []error { return m }

// DeleteCommentResponseValidationError is the validation error returned by
// DeleteCommentResponse.Validate if the designated constraints aren't met.
type DeleteCommentResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteCommentResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteCommentResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteCommentResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteCommentResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteCommentResponseValidationError) ErrorName() string {
	return "DeleteCommentResponseValidationError"
}

// Error satisfies the builtin error interface
func (e DeleteCommentResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteCommentResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteCommentResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteCommentResponseValidationError{}
