package file

import (
	"fmt"

	"mosn.io/layotto/components/pkg/info"
)

type Registry interface {
	Register(fs ...*FileFactory)
	Create(name string) (File, error)
}

type FileFactory struct {
	Name          string
	FactoryMethod func() File
}

func NewFileFactory(name string, f func() File) *FileFactory {
	return &FileFactory{
		Name:          name,
		FactoryMethod: f,
	}
}

type FileStoreRegistry struct {
	files map[string]func() File
	info  *info.RuntimeInfo
}

func NewRegistry(info *info.RuntimeInfo) Registry {
	info.AddService(ServiceName)
	return &FileStoreRegistry{
		files: make(map[string]func() File),
		info:  info,
	}
}

func (r *FileStoreRegistry) Register(fs ...*FileFactory) {
	for _, f := range fs {
		r.files[f.Name] = f.FactoryMethod
		r.info.RegisterComponent(ServiceName, f.Name)
	}
}

func (r *FileStoreRegistry) Create(name string) (File, error) {
	if f, ok := r.files[name]; ok {
		r.info.LoadComponent(ServiceName, name)
		return f(), nil
	}
	return nil, fmt.Errorf("service component %s is not regsitered", name)
}
