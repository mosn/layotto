// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
	code int
	msg  string
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
