package file

import "io"

const ServiceName = "file"

type File interface {
	Init(*FileConfig) error
	Put(*PutFileStu) error
	Get(*GetFileStu) (io.ReadCloser, error)
	List(*ListRequest) (*ListResp, error)
	Del(*DelRequest) error
	Complete(int64, bool) error
}
