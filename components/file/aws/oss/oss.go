package oss

import (
	"io"

	"mosn.io/layotto/components/file"
)

// AwsOss is a binding for
type AwsOss struct {
}

func NewAwsOss() file.File {
	return &AwsOss{}
}

func (a *AwsOss) Init(*file.FileConfig) error {
	return nil
}
func (a *AwsOss) Put(*file.PutFileStu) error {
	return nil
}
func (a *AwsOss) Get(*file.GetFileStu) (io.ReadCloser, error) {
	return nil, nil
}
func (a *AwsOss) List(*file.ListRequest) (*file.ListResp, error) {
	return nil, nil
}
func (a *AwsOss) Del(*file.DelRequest) error {
	return nil
}
