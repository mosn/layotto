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

package aws

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"mosn.io/layotto/components/pkg/utils"

	"mosn.io/layotto/components/oss"

	"github.com/jinzhu/copier"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/file"
)

const cfg = `[
				{
					"buckets":["bucket1"],
					"endpoint": "protocol://service-code.region-code.amazonaws.com",
					"accessKeyID": "accessKey",
					"accessKeySecret": "secret",
					"region": "us-west-2"
				}
			]`

func TestAwsOss_Init(t *testing.T) {
	oss := NewAwsFile()
	err := oss.Init(context.TODO(), &file.FileConfig{})
	assert.Equal(t, err.Error(), "invalid config for aws oss")
	err = oss.Init(context.TODO(), &file.FileConfig{Metadata: []byte(cfg)})
	assert.Equal(t, nil, err)
}

func TestAwsOss_IsAwsMetaValid(t *testing.T) {
	mt := &utils.OssMetadata{}
	a := AwsOss{}
	assert.False(t, a.isAwsMetaValid(mt))
	mt.AccessKeyID = "a"
	assert.False(t, a.isAwsMetaValid(mt))
	mt.Endpoint = "a"
	assert.False(t, a.isAwsMetaValid(mt))
	mt.AccessKeySecret = "a"
	assert.True(t, a.isAwsMetaValid(mt))

}

func TestAwsOss_Put(t *testing.T) {
	oss := NewAwsFile()
	err := oss.Init(context.TODO(), &file.FileConfig{Metadata: []byte(cfg)})
	assert.Equal(t, nil, err)

	req := &file.PutFileStu{
		FileName: "",
	}
	err = oss.Put(context.Background(), req)
	assert.Equal(t, err.Error(), "aws.s3 put file[] fail,err: invalid fileName format")

	req.FileName = "/a.txt"
	err = oss.Put(context.Background(), req)
	assert.Equal(t, err.Error(), "aws.s3 put file[/a.txt] fail,err: invalid fileName format")
}

func TestAwsOss_Get(t *testing.T) {
	oss := NewAwsFile()
	err := oss.Init(context.TODO(), &file.FileConfig{Metadata: []byte(cfg)})
	assert.Equal(t, nil, err)

	putReq := &file.PutFileStu{
		FileName: "/a.txt",
	}
	err = oss.Put(context.Background(), putReq)

	assert.Equal(t, err.Error(), "aws.s3 put file[/a.txt] fail,err: invalid fileName format")

	req := &file.GetFileStu{
		FileName: "",
	}
	_, err = oss.Get(context.Background(), req)
	assert.Equal(t, err.Error(), "aws.s3 get file[] fail,err: invalid fileName format")

	req.FileName = "/a.txt"
	_, err = oss.Get(context.Background(), req)
	assert.Equal(t, err.Error(), "aws.s3 get file[/a.txt] fail,err: invalid fileName format")
}

type fun = func() (string, error)

func TestCopier(t *testing.T) {
	hello := "hello"
	target := &oss.ListObjectsOutput{}
	source := &s3.ListObjectsOutput{Delimiter: &hello, EncodingType: "encoding type"}
	re := reflect.TypeOf(source)
	h, _ := re.Elem().FieldByName("EncodingType")
	fmt.Println(h.Type.Name(), h.Type.Kind())
	err := copier.Copy(target, source)
	if err != nil {
		t.Fail()
	}
	var s fun
	if s == nil {
		fmt.Printf("s is nil \n")
	}
	fmt.Println(target)

}
