// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: order/v1/order.proto

package adminorderv1

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

	orderv1 "backend/api/order/v1"
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

	_ = orderv1.PaymentStatus(0)
)

// Validate checks the field values on GetAllOrdersReq with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetAllOrdersReq) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetAllOrdersReq with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetAllOrdersReqMultiError, or nil if none found.
func (m *GetAllOrdersReq) ValidateAll() error {
	return m.validate(true)
}

func (m *GetAllOrdersReq) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Page

	// no validation rules for PageSize

	if len(errors) > 0 {
		return GetAllOrdersReqMultiError(errors)
	}

	return nil
}

// GetAllOrdersReqMultiError is an error wrapping multiple validation errors
// returned by GetAllOrdersReq.ValidateAll() if the designated constraints
// aren't met.
type GetAllOrdersReqMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetAllOrdersReqMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetAllOrdersReqMultiError) AllErrors() []error { return m }

// GetAllOrdersReqValidationError is the validation error returned by
// GetAllOrdersReq.Validate if the designated constraints aren't met.
type GetAllOrdersReqValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetAllOrdersReqValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetAllOrdersReqValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetAllOrdersReqValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetAllOrdersReqValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetAllOrdersReqValidationError) ErrorName() string { return "GetAllOrdersReqValidationError" }

// Error satisfies the builtin error interface
func (e GetAllOrdersReqValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetAllOrdersReq.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetAllOrdersReqValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetAllOrdersReqValidationError{}

// Validate checks the field values on SubOrderItem with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *SubOrderItem) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SubOrderItem with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in SubOrderItemMultiError, or
// nil if none found.
func (m *SubOrderItem) ValidateAll() error {
	return m.validate(true)
}

func (m *SubOrderItem) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetItem()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, SubOrderItemValidationError{
					field:  "Item",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, SubOrderItemValidationError{
					field:  "Item",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetItem()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return SubOrderItemValidationError{
				field:  "Item",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Cost

	if len(errors) > 0 {
		return SubOrderItemMultiError(errors)
	}

	return nil
}

// SubOrderItemMultiError is an error wrapping multiple validation errors
// returned by SubOrderItem.ValidateAll() if the designated constraints aren't met.
type SubOrderItemMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SubOrderItemMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SubOrderItemMultiError) AllErrors() []error { return m }

// SubOrderItemValidationError is the validation error returned by
// SubOrderItem.Validate if the designated constraints aren't met.
type SubOrderItemValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SubOrderItemValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SubOrderItemValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SubOrderItemValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SubOrderItemValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SubOrderItemValidationError) ErrorName() string { return "SubOrderItemValidationError" }

// Error satisfies the builtin error interface
func (e SubOrderItemValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSubOrderItem.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SubOrderItemValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SubOrderItemValidationError{}

// Validate checks the field values on SubOrder with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *SubOrder) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SubOrder with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in SubOrderMultiError, or nil
// if none found.
func (m *SubOrder) ValidateAll() error {
	return m.validate(true)
}

func (m *SubOrder) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for OrderId

	// no validation rules for SubOrderId

	// no validation rules for TotalAmount

	// no validation rules for ConsumerId

	if all {
		switch v := interface{}(m.GetAddress()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, SubOrderValidationError{
					field:  "Address",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, SubOrderValidationError{
					field:  "Address",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetAddress()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return SubOrderValidationError{
				field:  "Address",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for ConsumerEmail

	if utf8.RuneCountInString(m.GetCurrency()) != 3 {
		err := SubOrderValidationError{
			field:  "Currency",
			reason: "value length must be 3 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)

	}

	for idx, item := range m.GetSubOrderItems() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, SubOrderValidationError{
						field:  fmt.Sprintf("SubOrderItems[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, SubOrderValidationError{
						field:  fmt.Sprintf("SubOrderItems[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return SubOrderValidationError{
					field:  fmt.Sprintf("SubOrderItems[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for PaymentStatus

	// no validation rules for ShippingStatus

	if all {
		switch v := interface{}(m.GetCreatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, SubOrderValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, SubOrderValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return SubOrderValidationError{
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
				errors = append(errors, SubOrderValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, SubOrderValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetUpdatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return SubOrderValidationError{
				field:  "UpdatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return SubOrderMultiError(errors)
	}

	return nil
}

// SubOrderMultiError is an error wrapping multiple validation errors returned
// by SubOrder.ValidateAll() if the designated constraints aren't met.
type SubOrderMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SubOrderMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SubOrderMultiError) AllErrors() []error { return m }

// SubOrderValidationError is the validation error returned by
// SubOrder.Validate if the designated constraints aren't met.
type SubOrderValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SubOrderValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SubOrderValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SubOrderValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SubOrderValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SubOrderValidationError) ErrorName() string { return "SubOrderValidationError" }

// Error satisfies the builtin error interface
func (e SubOrderValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSubOrder.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SubOrderValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SubOrderValidationError{}

// Validate checks the field values on AdminOrderReply with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *AdminOrderReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on AdminOrderReply with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// AdminOrderReplyMultiError, or nil if none found.
func (m *AdminOrderReply) ValidateAll() error {
	return m.validate(true)
}

func (m *AdminOrderReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetOrders() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, AdminOrderReplyValidationError{
						field:  fmt.Sprintf("Orders[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, AdminOrderReplyValidationError{
						field:  fmt.Sprintf("Orders[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return AdminOrderReplyValidationError{
					field:  fmt.Sprintf("Orders[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return AdminOrderReplyMultiError(errors)
	}

	return nil
}

// AdminOrderReplyMultiError is an error wrapping multiple validation errors
// returned by AdminOrderReply.ValidateAll() if the designated constraints
// aren't met.
type AdminOrderReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m AdminOrderReplyMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m AdminOrderReplyMultiError) AllErrors() []error { return m }

// AdminOrderReplyValidationError is the validation error returned by
// AdminOrderReply.Validate if the designated constraints aren't met.
type AdminOrderReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AdminOrderReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AdminOrderReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AdminOrderReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AdminOrderReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AdminOrderReplyValidationError) ErrorName() string { return "AdminOrderReplyValidationError" }

// Error satisfies the builtin error interface
func (e AdminOrderReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAdminOrderReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AdminOrderReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AdminOrderReplyValidationError{}
