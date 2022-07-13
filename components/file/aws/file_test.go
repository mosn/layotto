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

const cfgwithuid = `[
				{	
					"uid": "1",
					"buckets":["bucket3"],
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

func TestAwsOss_SelectClient(t *testing.T) {
	oss := &AwsOss{
		client: make(map[string]*s3.Client),
		meta:   make(map[string]*utils.OssMetadata),
	}
	err := oss.Init(context.TODO(), &file.FileConfig{Metadata: []byte(cfg)})
	assert.Equal(t, nil, err)

	// not specify endpoint, select default client
	_, err = oss.selectClient("", "bucket1")
	assert.Nil(t, err)

	// specify endpoint equal config
	client, _ := oss.selectClient("", "bucket1")
	assert.NotNil(t, client)

	// specicy not exist endpoint, select default one
	client, err = oss.selectClient("", "bucket1")
	assert.Nil(t, err)
	assert.NotNil(t, client)
	// new client with endpoint
	oss.client["bucket2"] = &s3.Client{}
	client, _ = oss.selectClient("", "bucket2")
	assert.NotNil(t, client)

	err = oss.Init(context.TODO(), &file.FileConfig{Metadata: []byte(cfgwithuid)})
	assert.Equal(t, nil, err)

	// specify uid
	client, _ = oss.selectClient("1", "bucket1")
	assert.Equal(t, client, oss.client["1"])
	assert.NotNil(t, client)
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
	assert.Equal(t, err.Error(), "awsoss put file[] fail,err: invalid fileName format")

	req.FileName = "/a.txt"
	err = oss.Put(context.Background(), req)
	assert.Equal(t, err.Error(), "awsoss put file[/a.txt] fail,err: invalid fileName format")
}

func TestAwsOss_Get(t *testing.T) {
	oss := NewAwsFile()
	err := oss.Init(context.TODO(), &file.FileConfig{Metadata: []byte(cfg)})
	assert.Equal(t, nil, err)

	putReq := &file.PutFileStu{
		FileName: "/a.txt",
	}
	err = oss.Put(context.Background(), putReq)

	assert.Equal(t, err.Error(), "awsoss put file[/a.txt] fail,err: invalid fileName format")

	req := &file.GetFileStu{
		FileName: "",
	}
	_, err = oss.Get(context.Background(), req)
	assert.Equal(t, err.Error(), "awsoss get file[] fail,err: invalid fileName format")

	req.FileName = "/a.txt"
	_, err = oss.Get(context.Background(), req)
	assert.Equal(t, err.Error(), "awsoss get file[/a.txt] fail,err: invalid fileName format")
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
