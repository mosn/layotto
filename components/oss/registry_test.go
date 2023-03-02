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

package oss

import (
	"context"
	"testing"

	"mosn.io/layotto/components/pkg/info"
)

type mockOss struct{}

func (m *mockOss) Init(ctx context.Context, cfg *Config) error {
	return nil
}

func (m *mockOss) GetObject(ctx context.Context, input *GetObjectInput) (*GetObjectOutput, error) {
	return nil, nil
}

func (m *mockOss) PutObject(ctx context.Context, input *PutObjectInput) (*PutObjectOutput, error) {
	return nil, nil
}

func (m *mockOss) DeleteObject(ctx context.Context, input *DeleteObjectInput) (*DeleteObjectOutput, error) {
	return nil, nil
}

func (m *mockOss) PutObjectTagging(ctx context.Context, input *PutObjectTaggingInput) (*PutObjectTaggingOutput, error) {
	return nil, nil
}

func (m *mockOss) DeleteObjectTagging(ctx context.Context, input *DeleteObjectTaggingInput) (*DeleteObjectTaggingOutput, error) {
	return nil, nil
}

func (m *mockOss) GetObjectTagging(ctx context.Context, input *GetObjectTaggingInput) (*GetObjectTaggingOutput, error) {
	return nil, nil
}

func (m *mockOss) CopyObject(ctx context.Context, input *CopyObjectInput) (*CopyObjectOutput, error) {
	return nil, nil
}

func (m *mockOss) DeleteObjects(ctx context.Context, input *DeleteObjectsInput) (*DeleteObjectsOutput, error) {
	return nil, nil
}

func (m *mockOss) ListObjects(ctx context.Context, input *ListObjectsInput) (*ListObjectsOutput, error) {
	return nil, nil
}

func (m *mockOss) GetObjectCannedAcl(ctx context.Context, input *GetObjectCannedAclInput) (*GetObjectCannedAclOutput, error) {
	return nil, nil
}

func (m *mockOss) PutObjectCannedAcl(ctx context.Context, input *PutObjectCannedAclInput) (*PutObjectCannedAclOutput, error) {
	return nil, nil
}

func (m *mockOss) RestoreObject(ctx context.Context, input *RestoreObjectInput) (*RestoreObjectOutput, error) {
	return nil, nil
}

func (m *mockOss) CreateMultipartUpload(ctx context.Context, input *CreateMultipartUploadInput) (*CreateMultipartUploadOutput, error) {
	return nil, nil
}

func (m *mockOss) UploadPart(ctx context.Context, input *UploadPartInput) (*UploadPartOutput, error) {
	return nil, nil
}

func (m *mockOss) UploadPartCopy(ctx context.Context, input *UploadPartCopyInput) (*UploadPartCopyOutput, error) {
	return nil, nil
}

func (m *mockOss) CompleteMultipartUpload(ctx context.Context, input *CompleteMultipartUploadInput) (*CompleteMultipartUploadOutput, error) {
	return nil, nil
}

func (m *mockOss) AbortMultipartUpload(ctx context.Context, input *AbortMultipartUploadInput) (*AbortMultipartUploadOutput, error) {
	return nil, nil
}

func (m *mockOss) ListMultipartUploads(ctx context.Context, input *ListMultipartUploadsInput) (*ListMultipartUploadsOutput, error) {
	return nil, nil
}

func TestRegistry(t *testing.T) {
	r := NewRegistry(info.NewRuntimeInfo())

	// Register a factory
	f := NewFactory("mock", func() Oss {
		return &mockOss{}
	})
	r.Register(f)

	// Test that the factory was registered
	if _, ok := r.(*registry).oss[f.CompType]; !ok {
		t.Errorf("Factory with component type %q was not registered", f.CompType)
	}

	// Test that a registered component can be created
	oss, err := r.Create(f.CompType)
	if err != nil {
		t.Errorf("Unexpected error creating component: %v", err)
	}
	if _, ok := oss.(*mockOss); !ok {
		t.Errorf("Expected component type *mockOss, but got %T", oss)
	}

	// Test that creating an unregistered component returns an error
	_, err = r.Create("unregistered")
	if err == nil {
		t.Error("Expected an error creating unregistered component, but got nil")
	}
	expectedErrMsg := "service component unregistered is not registered"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message %q, but got %q", expectedErrMsg, err.Error())
	}
}

func (m *mockOss) ListObjectVersions(ctx context.Context, input *ListObjectVersionsInput) (*ListObjectVersionsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockOss) HeadObject(ctx context.Context, input *HeadObjectInput) (*HeadObjectOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockOss) IsObjectExist(ctx context.Context, input *IsObjectExistInput) (*IsObjectExistOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockOss) SignURL(ctx context.Context, input *SignURLInput) (*SignURLOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockOss) UpdateDownloadBandwidthRateLimit(ctx context.Context, input *UpdateBandwidthRateLimitInput) error {
	//TODO implement me
	panic("implement me")
}

func (m *mockOss) UpdateUploadBandwidthRateLimit(ctx context.Context, input *UpdateBandwidthRateLimitInput) error {
	//TODO implement me
	panic("implement me")
}

func (m *mockOss) AppendObject(ctx context.Context, input *AppendObjectInput) (*AppendObjectOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockOss) ListParts(ctx context.Context, input *ListPartsInput) (*ListPartsOutput, error) {
	//TODO implement me
	panic("implement me")
}
