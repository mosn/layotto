// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: envoy/config/accesslog/v2/file.proto

package envoy_config_accesslog_v2

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang/protobuf/ptypes"
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
	_ = ptypes.DynamicAny{}
)

// Validate checks the field values on FileAccessLog with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *FileAccessLog) Validate() error {
	if m == nil {
		return nil
	}

	if len(m.GetPath()) < 1 {
		return FileAccessLogValidationError{
			field:  "Path",
			reason: "value length must be at least 1 bytes",
		}
	}

	switch m.AccessLogFormat.(type) {

	case *FileAccessLog_Format:
		// no validation rules for Format

	case *FileAccessLog_JsonFormat:

		if v, ok := interface{}(m.GetJsonFormat()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return FileAccessLogValidationError{
					field:  "JsonFormat",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *FileAccessLog_TypedJsonFormat:

		if v, ok := interface{}(m.GetTypedJsonFormat()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return FileAccessLogValidationError{
					field:  "TypedJsonFormat",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// FileAccessLogValidationError is the validation error returned by
// FileAccessLog.Validate if the designated constraints aren't met.
type FileAccessLogValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e FileAccessLogValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e FileAccessLogValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e FileAccessLogValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e FileAccessLogValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e FileAccessLogValidationError) ErrorName() string { return "FileAccessLogValidationError" }

// Error satisfies the builtin error interface
func (e FileAccessLogValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sFileAccessLog.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = FileAccessLogValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = FileAccessLogValidationError{}
