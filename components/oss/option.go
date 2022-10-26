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
	"bytes"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/jinzhu/copier"
)

var (
	Int64ToTime = copier.TypeConverter{
		SrcType: int64(0),
		DstType: &time.Time{},
		Fn: func(src interface{}) (interface{}, error) {
			s, _ := src.(int64)
			t := time.Unix(s, 0)
			return &t, nil
		},
	}
	TimeToInt64 = copier.TypeConverter{
		SrcType: &time.Time{},
		DstType: int64(0),
		Fn: func(src interface{}) (interface{}, error) {
			s, _ := src.(*time.Time)
			return s.Unix(), nil
		},
	}
)

func GetGetObjectOutput(ob *s3.GetObjectOutput) (*GetObjectOutput, error) {
	out := &GetObjectOutput{}
	err := copier.Copy(out, ob)
	if err != nil {
		return nil, err
	}
	out.DataStream = ob.Body
	return out, nil
}

func GetPutObjectOutput(resp *manager.UploadOutput) (*PutObjectOutput, error) {
	out := &PutObjectOutput{}
	err := copier.Copy(out, resp)
	if err != nil {
		return nil, err
	}
	return out, err
}

func GetDeleteObjectOutput(resp *s3.DeleteObjectOutput) (*DeleteObjectOutput, error) {
	versionId := ""
	if resp.VersionId != nil {
		versionId = *resp.VersionId
	}
	return &DeleteObjectOutput{DeleteMarker: resp.DeleteMarker, RequestCharged: string(resp.RequestCharged), VersionId: versionId}, nil
}

func GetDeleteObjectTaggingOutput(resp *s3.DeleteObjectTaggingOutput) (*DeleteObjectTaggingOutput, error) {
	versionId := ""
	if resp.VersionId != nil {
		versionId = *resp.VersionId
	}
	return &DeleteObjectTaggingOutput{VersionId: versionId}, nil
}

func GetGetObjectTaggingOutput(resp *s3.GetObjectTaggingOutput) (*GetObjectTaggingOutput, error) {
	output := &GetObjectTaggingOutput{Tags: map[string]string{}}
	for _, tags := range resp.TagSet {
		output.Tags[*tags.Key] = *tags.Value
	}
	return output, nil
}

func GetCopyObjectOutput(resp *s3.CopyObjectOutput) (*CopyObjectOutput, error) {
	out := &CopyObjectOutput{}
	err := copier.CopyWithOption(out, resp, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func GetListObjectsOutput(resp *s3.ListObjectsOutput) (*ListObjectsOutput, error) {
	output := &ListObjectsOutput{}
	err := copier.CopyWithOption(output, resp, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{TimeToInt64}})
	// if not return NextMarker, use the value of the last Key in the response as the marker
	if output.IsTruncated && output.NextMarker == "" {
		index := len(output.Contents) - 1
		output.NextMarker = output.Contents[index].Key
	}
	return output, err
}

func GetGetObjectCannedAclOutput(resp *s3.GetObjectAclOutput) (*GetObjectCannedAclOutput, error) {
	out := &GetObjectCannedAclOutput{}
	err := copier.CopyWithOption(out, resp, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
	if err != nil {
		return nil, err
	}
	bs, _ := json.Marshal(resp.Grants)
	var bf bytes.Buffer
	err = json.Indent(&bf, bs, "", "\t")
	if err != nil {
		return nil, err
	}
	out.CannedAcl = bf.String()
	return out, nil
}

func GetUploadPartOutput(resp *s3.UploadPartOutput) (*UploadPartOutput, error) {
	output := &UploadPartOutput{}
	err := copier.Copy(output, resp)
	if err != nil {
		return nil, err
	}
	return output, err
}

func GetUploadPartCopyOutput(resp *s3.UploadPartCopyOutput) (*UploadPartCopyOutput, error) {
	out := &UploadPartCopyOutput{}
	err := copier.CopyWithOption(out, resp, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
	if err != nil {
		return nil, err
	}
	return out, err
}

func GetListPartsOutput(resp *s3.ListPartsOutput) (*ListPartsOutput, error) {
	output := &ListPartsOutput{}
	err := copier.CopyWithOption(output, resp, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
	if err != nil {
		return nil, err
	}
	return output, err
}

func GetListMultipartUploadsOutput(resp *s3.ListMultipartUploadsOutput) (*ListMultipartUploadsOutput, error) {
	output := &ListMultipartUploadsOutput{CommonPrefixes: []string{}, Uploads: []*MultipartUpload{}}
	err := copier.Copy(output, resp)
	if err != nil {
		return nil, err
	}
	for _, v := range resp.CommonPrefixes {
		output.CommonPrefixes = append(output.CommonPrefixes, *v.Prefix)
	}
	for _, v := range resp.Uploads {
		upload := &MultipartUpload{}
		copier.CopyWithOption(upload, v, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		output.Uploads = append(output.Uploads, upload)
	}
	return output, err
}

func GetListObjectVersionsOutput(resp *s3.ListObjectVersionsOutput) (*ListObjectVersionsOutput, error) {
	output := &ListObjectVersionsOutput{}
	err := copier.Copy(output, resp)
	if err != nil {
		return nil, err
	}
	for _, v := range resp.CommonPrefixes {
		output.CommonPrefixes = append(output.CommonPrefixes, *v.Prefix)
	}
	for _, v := range resp.DeleteMarkers {
		entry := &DeleteMarkerEntry{IsLatest: v.IsLatest, Key: *v.Key, Owner: &Owner{DisplayName: *v.Owner.DisplayName, ID: *v.Owner.ID}, VersionId: *v.VersionId}
		output.DeleteMarkers = append(output.DeleteMarkers, entry)
	}
	for _, v := range resp.Versions {
		version := &ObjectVersion{}
		copier.CopyWithOption(version, v, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{TimeToInt64}})
		output.Versions = append(output.Versions, version)
	}
	return output, err
}
