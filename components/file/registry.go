/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package file

import (
	"fmt"
	"mosn.io/layotto/components/pkg/info"
)

type Registry interface {
	Register(fs ...*FileFactory)
	Create(compType string) (File, error)
}

type FileFactory struct {
	CompType      string
	FactoryMethod func() File
}

func NewFileFactory(CompType string, f func() File) *FileFactory {
	return &FileFactory{
		CompType:      CompType,
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
		r.files[f.CompType] = f.FactoryMethod
		r.info.RegisterComponent(ServiceName, f.CompType)
	}
}

func (r *FileStoreRegistry) Create(compType string) (File, error) {
	if f, ok := r.files[compType]; ok {
		r.info.LoadComponent(ServiceName, compType)
		return f(), nil
	}
	return nil, fmt.Errorf("service component %s is not regsitered", compType)
}

type OssRegistry interface {
	Register(fs ...*OssFactory)
	Create(name string) (Oss, error)
}

type OssFactory struct {
	Name          string
	FactoryMethod func() Oss
}

func NewOssFactory(name string, f func() Oss) *OssFactory {
	return &OssFactory{
		Name:          name,
		FactoryMethod: f,
	}
}

type OssStoreRegistry struct {
	files map[string]func() Oss
	info  *info.RuntimeInfo
}

func NewOssRegistry(info *info.RuntimeInfo) OssRegistry {
	info.AddService(ServiceName)
	return &OssStoreRegistry{
		files: make(map[string]func() Oss),
		info:  info,
	}
}

func (r *OssStoreRegistry) Register(fs ...*OssFactory) {
	for _, f := range fs {
		r.files[f.Name] = f.FactoryMethod
		r.info.RegisterComponent(ServiceName, f.Name)
	}
}

func (r *OssStoreRegistry) Create(name string) (Oss, error) {
	if f, ok := r.files[name]; ok {
		r.info.LoadComponent(ServiceName, name)
		return f(), nil
	}
	return nil, fmt.Errorf("service component %s is not regsitered", name)
}
