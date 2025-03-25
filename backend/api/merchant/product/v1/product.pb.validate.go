// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: product/v1/product.proto

package productv1

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

// Validate checks the field values on GetMerchantProductRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetMerchantProductRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetMerchantProductRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetMerchantProductRequestMultiError, or nil if none found.
func (m *GetMerchantProductRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetMerchantProductRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return GetMerchantProductRequestMultiError(errors)
	}

	return nil
}

// GetMerchantProductRequestMultiError is an error wrapping multiple validation
// errors returned by GetMerchantProductRequest.ValidateAll() if the
// designated constraints aren't met.
type GetMerchantProductRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetMerchantProductRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetMerchantProductRequestMultiError) AllErrors() []error { return m }

// GetMerchantProductRequestValidationError is the validation error returned by
// GetMerchantProductRequest.Validate if the designated constraints aren't met.
type GetMerchantProductRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetMerchantProductRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetMerchantProductRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetMerchantProductRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetMerchantProductRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetMerchantProductRequestValidationError) ErrorName() string {
	return "GetMerchantProductRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetMerchantProductRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetMerchantProductRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetMerchantProductRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetMerchantProductRequestValidationError{}

// Validate checks the field values on UpdateProductRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateProductRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateProductRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UpdateProductRequestMultiError, or nil if none found.
func (m *UpdateProductRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateProductRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for MerchantId

	// no validation rules for Name

	// no validation rules for Description

	// no validation rules for Price

	if len(errors) > 0 {
		return UpdateProductRequestMultiError(errors)
	}

	return nil
}

// UpdateProductRequestMultiError is an error wrapping multiple validation
// errors returned by UpdateProductRequest.ValidateAll() if the designated
// constraints aren't met.
type UpdateProductRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateProductRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateProductRequestMultiError) AllErrors() []error { return m }

// UpdateProductRequestValidationError is the validation error returned by
// UpdateProductRequest.Validate if the designated constraints aren't met.
type UpdateProductRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateProductRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateProductRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateProductRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateProductRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateProductRequestValidationError) ErrorName() string {
	return "UpdateProductRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateProductRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateProductRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateProductRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateProductRequestValidationError{}

// Validate checks the field values on UpdateProductReply with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateProductReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateProductReply with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UpdateProductReplyMultiError, or nil if none found.
func (m *UpdateProductReply) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateProductReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Message

	// no validation rules for Code

	if len(errors) > 0 {
		return UpdateProductReplyMultiError(errors)
	}

	return nil
}

// UpdateProductReplyMultiError is an error wrapping multiple validation errors
// returned by UpdateProductReply.ValidateAll() if the designated constraints
// aren't met.
type UpdateProductReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateProductReplyMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateProductReplyMultiError) AllErrors() []error { return m }

// UpdateProductReplyValidationError is the validation error returned by
// UpdateProductReply.Validate if the designated constraints aren't met.
type UpdateProductReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateProductReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateProductReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateProductReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateProductReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateProductReplyValidationError) ErrorName() string {
	return "UpdateProductReplyValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateProductReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateProductReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateProductReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateProductReplyValidationError{}
