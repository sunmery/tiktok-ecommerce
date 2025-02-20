// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: v1/service.proto

package payment

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

// Validate checks the field values on CreditCardInfo with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *CreditCardInfo) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreditCardInfo with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in CreditCardInfoMultiError,
// or nil if none found.
func (m *CreditCardInfo) ValidateAll() error {
	return m.validate(true)
}

func (m *CreditCardInfo) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetNumber()) != 16 {
		err := CreditCardInfoValidationError{
			field:  "Number",
			reason: "value length must be 16 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)

	}

	if !_CreditCardInfo_Number_Pattern.MatchString(m.GetNumber()) {
		err := CreditCardInfoValidationError{
			field:  "Number",
			reason: "value does not match regex pattern \"^[0-9]{16}$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if val := m.GetCvv(); val < 0 || val > 9999 {
		err := CreditCardInfoValidationError{
			field:  "Cvv",
			reason: "value must be inside range [0, 9999]",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetExpirationYear() < 23 {
		err := CreditCardInfoValidationError{
			field:  "ExpirationYear",
			reason: "value must be greater than or equal to 23",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if val := m.GetExpirationMonth(); val < 1 || val > 12 {
		err := CreditCardInfoValidationError{
			field:  "ExpirationMonth",
			reason: "value must be inside range [1, 12]",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return CreditCardInfoMultiError(errors)
	}

	return nil
}

// CreditCardInfoMultiError is an error wrapping multiple validation errors
// returned by CreditCardInfo.ValidateAll() if the designated constraints
// aren't met.
type CreditCardInfoMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreditCardInfoMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreditCardInfoMultiError) AllErrors() []error { return m }

// CreditCardInfoValidationError is the validation error returned by
// CreditCardInfo.Validate if the designated constraints aren't met.
type CreditCardInfoValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreditCardInfoValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreditCardInfoValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreditCardInfoValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreditCardInfoValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreditCardInfoValidationError) ErrorName() string { return "CreditCardInfoValidationError" }

// Error satisfies the builtin error interface
func (e CreditCardInfoValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreditCardInfo.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreditCardInfoValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreditCardInfoValidationError{}

var _CreditCardInfo_Number_Pattern = regexp.MustCompile("^[0-9]{16}$")

// Validate checks the field values on ChargeReq with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ChargeReq) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ChargeReq with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ChargeReqMultiError, or nil
// if none found.
func (m *ChargeReq) ValidateAll() error {
	return m.validate(true)
}

func (m *ChargeReq) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetAmount() <= 0 {
		err := ChargeReqValidationError{
			field:  "Amount",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetCreditCard() == nil {
		err := ChargeReqValidationError{
			field:  "CreditCard",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetCreditCard()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ChargeReqValidationError{
					field:  "CreditCard",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ChargeReqValidationError{
					field:  "CreditCard",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreditCard()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ChargeReqValidationError{
				field:  "CreditCard",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if utf8.RuneCountInString(m.GetOrderId()) < 1 {
		err := ChargeReqValidationError{
			field:  "OrderId",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetUserId()) < 1 {
		err := ChargeReqValidationError{
			field:  "UserId",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return ChargeReqMultiError(errors)
	}

	return nil
}

// ChargeReqMultiError is an error wrapping multiple validation errors returned
// by ChargeReq.ValidateAll() if the designated constraints aren't met.
type ChargeReqMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ChargeReqMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ChargeReqMultiError) AllErrors() []error { return m }

// ChargeReqValidationError is the validation error returned by
// ChargeReq.Validate if the designated constraints aren't met.
type ChargeReqValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ChargeReqValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ChargeReqValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ChargeReqValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ChargeReqValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ChargeReqValidationError) ErrorName() string { return "ChargeReqValidationError" }

// Error satisfies the builtin error interface
func (e ChargeReqValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sChargeReq.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ChargeReqValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ChargeReqValidationError{}

// Validate checks the field values on ChargeResp with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ChargeResp) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ChargeResp with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ChargeRespMultiError, or
// nil if none found.
func (m *ChargeResp) ValidateAll() error {
	return m.validate(true)
}

func (m *ChargeResp) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for TransactionId

	if len(errors) > 0 {
		return ChargeRespMultiError(errors)
	}

	return nil
}

// ChargeRespMultiError is an error wrapping multiple validation errors
// returned by ChargeResp.ValidateAll() if the designated constraints aren't met.
type ChargeRespMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ChargeRespMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ChargeRespMultiError) AllErrors() []error { return m }

// ChargeRespValidationError is the validation error returned by
// ChargeResp.Validate if the designated constraints aren't met.
type ChargeRespValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ChargeRespValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ChargeRespValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ChargeRespValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ChargeRespValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ChargeRespValidationError) ErrorName() string { return "ChargeRespValidationError" }

// Error satisfies the builtin error interface
func (e ChargeRespValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sChargeResp.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ChargeRespValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ChargeRespValidationError{}
