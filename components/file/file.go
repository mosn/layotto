package file

import "io"

const ServiceName = "file"

type File interface {
	Init(*FileConfig) error
	Put(*PutFileStu) error
	Get(*GetFileStu) (io.ReadCloser, error)
	CompletePut(int64) error
}
