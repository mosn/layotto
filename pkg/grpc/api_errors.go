package grpc

import (
	"google.golang.org/grpc/codes"
	"mosn.io/layotto/components/file"
)

var (
	FileErrMap2GrpcErr = map[error]codes.Code{
		file.ErrInvalid:    codes.InvalidArgument,
		file.ErrNotExist:   codes.NotFound,
		file.ErrExist:      codes.AlreadyExists,
		file.ErrExpired:    codes.DataLoss,
		file.ErrPermission: codes.PermissionDenied,
	}
)
