package file

import "io"

type FileService interface {
	Init(*FileConfig) error
	Put(*PutFileStu) error
	Get(*GetFileStu) (io.ReadCloser, error)
	CompletePut(int64) error
}
