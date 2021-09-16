package common

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	TimeoutCode int = iota
	UnavailebleCode
	InternalCode
	InvalidArgsCode
)

type CommonError interface {
	Code() int
	Msg() string
	Error() string
}

type commonError struct {
	code int    `json:"code"`
	msg  string `json:"msg"`
}

func (le *commonError) Code() int {
	return le.code
}

func (le *commonError) Msg() string {
	return le.msg
}

func (le *commonError) Error() string {
	return fmt.Sprintf("code %d, msg: %s", le.code, le.msg)
}

func Error(code int, msg string) CommonError {
	return &commonError{code: code, msg: msg}
}

func Errorf(code int, format string, a ...interface{}) CommonError {
	return &commonError{code: code, msg: fmt.Sprintf(format, a...)}
}

func ToGrpcError(err error) error {
	switch v := err.(type) {
	case CommonError:
		var code codes.Code
		switch v.Code() {
		case TimeoutCode:
			code = codes.DeadlineExceeded
		case UnavailebleCode:
			code = codes.Unavailable
		case InternalCode:
			code = codes.Internal
		case InvalidArgsCode:
			code = codes.InvalidArgument
		default:
			code = codes.Unknown
		}
		return status.Error(code, v.Msg())
	default:
		return status.Error(codes.Unknown, err.Error())
	}
}
