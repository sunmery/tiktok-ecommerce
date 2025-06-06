// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: v1/order.proto

package orderv1

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

// Validate checks the field values on MarkOrderPaidReq with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *MarkOrderPaidReq) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on MarkOrderPaidReq with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// MarkOrderPaidReqMultiError, or nil if none found.
func (m *MarkOrderPaidReq) ValidateAll() error {
	return m.validate(true)
}

func (m *MarkOrderPaidReq) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for OrderId

	if m.UserId != nil {
		// no validation rules for UserId
	}

	if len(errors) > 0 {
		return MarkOrderPaidReqMultiError(errors)
	}

	return nil
}

// MarkOrderPaidReqMultiError is an error wrapping multiple validation errors
// returned by MarkOrderPaidReq.ValidateAll() if the designated constraints
// aren't met.
type MarkOrderPaidReqMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m MarkOrderPaidReqMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m MarkOrderPaidReqMultiError) AllErrors() []error { return m }

// MarkOrderPaidReqValidationError is the validation error returned by
// MarkOrderPaidReq.Validate if the designated constraints aren't met.
type MarkOrderPaidReqValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MarkOrderPaidReqValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MarkOrderPaidReqValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MarkOrderPaidReqValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MarkOrderPaidReqValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MarkOrderPaidReqValidationError) ErrorName() string { return "MarkOrderPaidReqValidationError" }

// Error satisfies the builtin error interface
func (e MarkOrderPaidReqValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMarkOrderPaidReq.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MarkOrderPaidReqValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MarkOrderPaidReqValidationError{}

// Validate checks the field values on MarkOrderPaidResp with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *MarkOrderPaidResp) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on MarkOrderPaidResp with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// MarkOrderPaidRespMultiError, or nil if none found.
func (m *MarkOrderPaidResp) ValidateAll() error {
	return m.validate(true)
}

func (m *MarkOrderPaidResp) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return MarkOrderPaidRespMultiError(errors)
	}

	return nil
}

// MarkOrderPaidRespMultiError is an error wrapping multiple validation errors
// returned by MarkOrderPaidResp.ValidateAll() if the designated constraints
// aren't met.
type MarkOrderPaidRespMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m MarkOrderPaidRespMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m MarkOrderPaidRespMultiError) AllErrors() []error { return m }

// MarkOrderPaidRespValidationError is the validation error returned by
// MarkOrderPaidResp.Validate if the designated constraints aren't met.
type MarkOrderPaidRespValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MarkOrderPaidRespValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MarkOrderPaidRespValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MarkOrderPaidRespValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MarkOrderPaidRespValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MarkOrderPaidRespValidationError) ErrorName() string {
	return "MarkOrderPaidRespValidationError"
}

// Error satisfies the builtin error interface
func (e MarkOrderPaidRespValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMarkOrderPaidResp.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MarkOrderPaidRespValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MarkOrderPaidRespValidationError{}
