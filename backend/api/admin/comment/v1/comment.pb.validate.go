// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: comment/v1/comment.proto

package admincommentv1

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

// Validate checks the field values on SensitiveWord with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *SensitiveWord) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SensitiveWord with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in SensitiveWordMultiError, or
// nil if none found.
func (m *SensitiveWord) ValidateAll() error {
	return m.validate(true)
}

func (m *SensitiveWord) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for CreatedBy

	// no validation rules for Category

	// no validation rules for Word

	// no validation rules for Level

	// no validation rules for IsActive

	if all {
		switch v := interface{}(m.GetCreatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, SensitiveWordValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, SensitiveWordValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return SensitiveWordValidationError{
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
				errors = append(errors, SensitiveWordValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, SensitiveWordValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetUpdatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return SensitiveWordValidationError{
				field:  "UpdatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if m.Id != nil {
		// no validation rules for Id
	}

	if len(errors) > 0 {
		return SensitiveWordMultiError(errors)
	}

	return nil
}

// SensitiveWordMultiError is an error wrapping multiple validation errors
// returned by SensitiveWord.ValidateAll() if the designated constraints
// aren't met.
type SensitiveWordMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SensitiveWordMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SensitiveWordMultiError) AllErrors() []error { return m }

// SensitiveWordValidationError is the validation error returned by
// SensitiveWord.Validate if the designated constraints aren't met.
type SensitiveWordValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SensitiveWordValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SensitiveWordValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SensitiveWordValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SensitiveWordValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SensitiveWordValidationError) ErrorName() string { return "SensitiveWordValidationError" }

// Error satisfies the builtin error interface
func (e SensitiveWordValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSensitiveWord.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SensitiveWordValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SensitiveWordValidationError{}

// Validate checks the field values on SetSensitiveWordsReq with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *SetSensitiveWordsReq) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SetSensitiveWordsReq with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// SetSensitiveWordsReqMultiError, or nil if none found.
func (m *SetSensitiveWordsReq) ValidateAll() error {
	return m.validate(true)
}

func (m *SetSensitiveWordsReq) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetSensitiveWords() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, SetSensitiveWordsReqValidationError{
						field:  fmt.Sprintf("SensitiveWords[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, SetSensitiveWordsReqValidationError{
						field:  fmt.Sprintf("SensitiveWords[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return SetSensitiveWordsReqValidationError{
					field:  fmt.Sprintf("SensitiveWords[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return SetSensitiveWordsReqMultiError(errors)
	}

	return nil
}

// SetSensitiveWordsReqMultiError is an error wrapping multiple validation
// errors returned by SetSensitiveWordsReq.ValidateAll() if the designated
// constraints aren't met.
type SetSensitiveWordsReqMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SetSensitiveWordsReqMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SetSensitiveWordsReqMultiError) AllErrors() []error { return m }

// SetSensitiveWordsReqValidationError is the validation error returned by
// SetSensitiveWordsReq.Validate if the designated constraints aren't met.
type SetSensitiveWordsReqValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SetSensitiveWordsReqValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SetSensitiveWordsReqValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SetSensitiveWordsReqValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SetSensitiveWordsReqValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SetSensitiveWordsReqValidationError) ErrorName() string {
	return "SetSensitiveWordsReqValidationError"
}

// Error satisfies the builtin error interface
func (e SetSensitiveWordsReqValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSetSensitiveWordsReq.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SetSensitiveWordsReqValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SetSensitiveWordsReqValidationError{}

// Validate checks the field values on SetSensitiveWordsReply with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *SetSensitiveWordsReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SetSensitiveWordsReply with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// SetSensitiveWordsReplyMultiError, or nil if none found.
func (m *SetSensitiveWordsReply) ValidateAll() error {
	return m.validate(true)
}

func (m *SetSensitiveWordsReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Rows

	if len(errors) > 0 {
		return SetSensitiveWordsReplyMultiError(errors)
	}

	return nil
}

// SetSensitiveWordsReplyMultiError is an error wrapping multiple validation
// errors returned by SetSensitiveWordsReply.ValidateAll() if the designated
// constraints aren't met.
type SetSensitiveWordsReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SetSensitiveWordsReplyMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SetSensitiveWordsReplyMultiError) AllErrors() []error { return m }

// SetSensitiveWordsReplyValidationError is the validation error returned by
// SetSensitiveWordsReply.Validate if the designated constraints aren't met.
type SetSensitiveWordsReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SetSensitiveWordsReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SetSensitiveWordsReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SetSensitiveWordsReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SetSensitiveWordsReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SetSensitiveWordsReplyValidationError) ErrorName() string {
	return "SetSensitiveWordsReplyValidationError"
}

// Error satisfies the builtin error interface
func (e SetSensitiveWordsReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSetSensitiveWordsReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SetSensitiveWordsReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SetSensitiveWordsReplyValidationError{}

// Validate checks the field values on GetSensitiveWordsReq with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetSensitiveWordsReq) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetSensitiveWordsReq with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetSensitiveWordsReqMultiError, or nil if none found.
func (m *GetSensitiveWordsReq) ValidateAll() error {
	return m.validate(true)
}

func (m *GetSensitiveWordsReq) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Page

	// no validation rules for PageSize

	if len(errors) > 0 {
		return GetSensitiveWordsReqMultiError(errors)
	}

	return nil
}

// GetSensitiveWordsReqMultiError is an error wrapping multiple validation
// errors returned by GetSensitiveWordsReq.ValidateAll() if the designated
// constraints aren't met.
type GetSensitiveWordsReqMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetSensitiveWordsReqMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetSensitiveWordsReqMultiError) AllErrors() []error { return m }

// GetSensitiveWordsReqValidationError is the validation error returned by
// GetSensitiveWordsReq.Validate if the designated constraints aren't met.
type GetSensitiveWordsReqValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetSensitiveWordsReqValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetSensitiveWordsReqValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetSensitiveWordsReqValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetSensitiveWordsReqValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetSensitiveWordsReqValidationError) ErrorName() string {
	return "GetSensitiveWordsReqValidationError"
}

// Error satisfies the builtin error interface
func (e GetSensitiveWordsReqValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetSensitiveWordsReq.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetSensitiveWordsReqValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetSensitiveWordsReqValidationError{}

// Validate checks the field values on GetSensitiveWordsReply with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetSensitiveWordsReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetSensitiveWordsReply with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetSensitiveWordsReplyMultiError, or nil if none found.
func (m *GetSensitiveWordsReply) ValidateAll() error {
	return m.validate(true)
}

func (m *GetSensitiveWordsReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetWords() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, GetSensitiveWordsReplyValidationError{
						field:  fmt.Sprintf("Words[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, GetSensitiveWordsReplyValidationError{
						field:  fmt.Sprintf("Words[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GetSensitiveWordsReplyValidationError{
					field:  fmt.Sprintf("Words[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return GetSensitiveWordsReplyMultiError(errors)
	}

	return nil
}

// GetSensitiveWordsReplyMultiError is an error wrapping multiple validation
// errors returned by GetSensitiveWordsReply.ValidateAll() if the designated
// constraints aren't met.
type GetSensitiveWordsReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetSensitiveWordsReplyMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetSensitiveWordsReplyMultiError) AllErrors() []error { return m }

// GetSensitiveWordsReplyValidationError is the validation error returned by
// GetSensitiveWordsReply.Validate if the designated constraints aren't met.
type GetSensitiveWordsReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetSensitiveWordsReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetSensitiveWordsReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetSensitiveWordsReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetSensitiveWordsReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetSensitiveWordsReplyValidationError) ErrorName() string {
	return "GetSensitiveWordsReplyValidationError"
}

// Error satisfies the builtin error interface
func (e GetSensitiveWordsReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetSensitiveWordsReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetSensitiveWordsReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetSensitiveWordsReplyValidationError{}

// Validate checks the field values on DeleteSensitiveWordReq with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DeleteSensitiveWordReq) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteSensitiveWordReq with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeleteSensitiveWordReqMultiError, or nil if none found.
func (m *DeleteSensitiveWordReq) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteSensitiveWordReq) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if len(errors) > 0 {
		return DeleteSensitiveWordReqMultiError(errors)
	}

	return nil
}

// DeleteSensitiveWordReqMultiError is an error wrapping multiple validation
// errors returned by DeleteSensitiveWordReq.ValidateAll() if the designated
// constraints aren't met.
type DeleteSensitiveWordReqMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteSensitiveWordReqMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteSensitiveWordReqMultiError) AllErrors() []error { return m }

// DeleteSensitiveWordReqValidationError is the validation error returned by
// DeleteSensitiveWordReq.Validate if the designated constraints aren't met.
type DeleteSensitiveWordReqValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteSensitiveWordReqValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteSensitiveWordReqValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteSensitiveWordReqValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteSensitiveWordReqValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteSensitiveWordReqValidationError) ErrorName() string {
	return "DeleteSensitiveWordReqValidationError"
}

// Error satisfies the builtin error interface
func (e DeleteSensitiveWordReqValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteSensitiveWordReq.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteSensitiveWordReqValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteSensitiveWordReqValidationError{}

// Validate checks the field values on DeleteSensitiveWordReply with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DeleteSensitiveWordReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteSensitiveWordReply with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeleteSensitiveWordReplyMultiError, or nil if none found.
func (m *DeleteSensitiveWordReply) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteSensitiveWordReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Success

	if len(errors) > 0 {
		return DeleteSensitiveWordReplyMultiError(errors)
	}

	return nil
}

// DeleteSensitiveWordReplyMultiError is an error wrapping multiple validation
// errors returned by DeleteSensitiveWordReply.ValidateAll() if the designated
// constraints aren't met.
type DeleteSensitiveWordReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteSensitiveWordReplyMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteSensitiveWordReplyMultiError) AllErrors() []error { return m }

// DeleteSensitiveWordReplyValidationError is the validation error returned by
// DeleteSensitiveWordReply.Validate if the designated constraints aren't met.
type DeleteSensitiveWordReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteSensitiveWordReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteSensitiveWordReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteSensitiveWordReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteSensitiveWordReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteSensitiveWordReplyValidationError) ErrorName() string {
	return "DeleteSensitiveWordReplyValidationError"
}

// Error satisfies the builtin error interface
func (e DeleteSensitiveWordReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteSensitiveWordReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteSensitiveWordReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteSensitiveWordReplyValidationError{}

// Validate checks the field values on UpdateSensitiveWordReq with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateSensitiveWordReq) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateSensitiveWordReq with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UpdateSensitiveWordReqMultiError, or nil if none found.
func (m *UpdateSensitiveWordReq) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateSensitiveWordReq) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for CreatedBy

	// no validation rules for Level

	// no validation rules for IsActive

	// no validation rules for Category

	// no validation rules for Word

	if len(errors) > 0 {
		return UpdateSensitiveWordReqMultiError(errors)
	}

	return nil
}

// UpdateSensitiveWordReqMultiError is an error wrapping multiple validation
// errors returned by UpdateSensitiveWordReq.ValidateAll() if the designated
// constraints aren't met.
type UpdateSensitiveWordReqMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateSensitiveWordReqMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateSensitiveWordReqMultiError) AllErrors() []error { return m }

// UpdateSensitiveWordReqValidationError is the validation error returned by
// UpdateSensitiveWordReq.Validate if the designated constraints aren't met.
type UpdateSensitiveWordReqValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateSensitiveWordReqValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateSensitiveWordReqValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateSensitiveWordReqValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateSensitiveWordReqValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateSensitiveWordReqValidationError) ErrorName() string {
	return "UpdateSensitiveWordReqValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateSensitiveWordReqValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateSensitiveWordReq.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateSensitiveWordReqValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateSensitiveWordReqValidationError{}

// Validate checks the field values on UpdateSensitiveWordReply with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateSensitiveWordReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateSensitiveWordReply with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UpdateSensitiveWordReplyMultiError, or nil if none found.
func (m *UpdateSensitiveWordReply) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateSensitiveWordReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Success

	if len(errors) > 0 {
		return UpdateSensitiveWordReplyMultiError(errors)
	}

	return nil
}

// UpdateSensitiveWordReplyMultiError is an error wrapping multiple validation
// errors returned by UpdateSensitiveWordReply.ValidateAll() if the designated
// constraints aren't met.
type UpdateSensitiveWordReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateSensitiveWordReplyMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateSensitiveWordReplyMultiError) AllErrors() []error { return m }

// UpdateSensitiveWordReplyValidationError is the validation error returned by
// UpdateSensitiveWordReply.Validate if the designated constraints aren't met.
type UpdateSensitiveWordReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateSensitiveWordReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateSensitiveWordReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateSensitiveWordReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateSensitiveWordReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateSensitiveWordReplyValidationError) ErrorName() string {
	return "UpdateSensitiveWordReplyValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateSensitiveWordReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateSensitiveWordReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateSensitiveWordReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateSensitiveWordReplyValidationError{}
