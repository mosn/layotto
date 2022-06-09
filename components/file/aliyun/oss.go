package aliyun

import (
	"context"
	"encoding/json"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"mosn.io/layotto/components/file"
	"mosn.io/layotto/components/file/factory"
)

const (
	DefaultClientInitFunc = "aliyun"
)

func NewAliyunOss() file.Oss {
	return &AliyunOSS{
		client:   make(map[string]*oss.Client),
		metadata: make(map[string]*OssMetadata),
	}
}

func init() {
	factory.RegisterInitFunc(DefaultClientInitFunc, AliyunDefaultInitFunc)
}

func AliyunDefaultInitFunc(staticConf json.RawMessage, DynConf map[string]string) (map[string]interface{}, error) {
	m := make([]*OssMetadata, 0)
	clients := make(map[string]interface{})
	err := json.Unmarshal(staticConf, &m)
	if err != nil {
		return nil, file.ErrInvalid
	}

	for _, v := range m {
		client, err := oss.New(v.Endpoint, v.AccessKeyID, v.AccessKeySecret)
		if err != nil {
			return nil, err
		}
		clients[v.Endpoint] = client
	}
	return clients, nil
}

func (a *AliyunOSS) InitConfig(ctx context.Context, config *file.FileConfig) error {
	a.method = config.Method
	a.rawData = config.Metadata
	return nil
}

func (a *AliyunOSS) InitClient(ctx context.Context, req *file.InitRequest) error {
	if a.method == "" {
		a.method = DefaultClientInitFunc
	}
	initFunc := factory.GetInitFunc(a.method)
	clients, err := initFunc(a.rawData, req.Metadata)
	if err != nil {
		return err
	}
	for k, v := range clients {
		a.client[k] = v.(*oss.Client)
	}
	return nil
}

func (a *AliyunOSS) GetObject(ctx context.Context, req *file.GetObjectInput) (io.ReadCloser, error) {
	client, err := a.selectClient(map[string]string{}, "")
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	body, err := bucket.GetObject(req.Key)
	return body, err
}

func (a *AliyunOSS) PutObject(ctx context.Context, req *file.PutObjectInput) (*file.PutObjectOutput, error) {
	cli, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	err = bucket.PutObject(req.Key, req.DataStream)
	return &file.PutObjectOutput{}, err
}

func (a *AliyunOSS) DeleteObject(ctx context.Context, req *file.DeleteObjectInput) (*file.DeleteObjectOutput, error) {
	cli, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	err = bucket.DeleteObject(req.Key)
	return &file.DeleteObjectOutput{}, err
}
func (a *AliyunOSS) DeleteObjects(ctx context.Context, req *file.DeleteObjectsInput) (*file.DeleteObjectsOutput, error) {
	cli, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	var objects []string
	for _, v := range req.Delete.Objects {
		objects = append(objects, v.Key)
	}
	resp, err := bucket.DeleteObjects(objects)
	if err != nil {
		return nil, err
	}
	out := &file.DeleteObjectsOutput{}
	for _, v := range resp.DeletedObjects {
		object := &file.DeletedObject{Key: v}
		out.Deleted = append(out.Deleted, object)
	}
	return out, err
}

//object标签
func (a *AliyunOSS) PutObjectTagging(ctx context.Context, req *file.PutBucketTaggingInput) (*file.PutBucketTaggingOutput, error) {
	cli, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	tagging := oss.Tagging{}
	for k, v := range req.Tags {
		tag := oss.Tag{Key: k, Value: v}
		tagging.Tags = append(tagging.Tags, tag)
	}
	err = bucket.PutObjectTagging(req.Key, tagging)
	return nil, err
}

func (a *AliyunOSS) DeleteObjectTagging(ctx context.Context, req *file.DeleteObjectTaggingInput) (*file.DeleteObjectTaggingOutput, error) {
	cli, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	err = bucket.DeleteObjectTagging(req.Key)
	return nil, err
}

func (a *AliyunOSS) GetObjectTagging(ctx context.Context, req *file.GetObjectTaggingInput) (*file.GetObjectTaggingOutput, error) {
	cli, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.GetObjectTagging(req.Key)
	if err != nil {
		return nil, err
	}
	out := &file.GetObjectTaggingOutput{Tags: map[string]string{}}
	for _, v := range resp.Tags {
		out.Tags[v.Key] = v.Value
	}
	return out, err
}

func (a *AliyunOSS) GetObjectAcl(ctx context.Context, req *file.GetObjectAclInput) (*file.GetObjectAclOutput, error) {
	cli, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.GetObjectACL(req.Key)
	if err != nil {
		return nil, err
	}
	output := &file.GetObjectAclOutput{Owner: &file.Owner{DisplayName: resp.Owner.DisplayName, ID: resp.Owner.ID}}
	grant := &file.Grant{&file.Grantee{DisplayName: resp.Owner.DisplayName, ID: resp.Owner.ID}, resp.ACL}
	output.Grants = append(output.Grants, grant)
	return output, err
}
func (a *AliyunOSS) PutObjectAcl(ctx context.Context, req *file.PutObjectAclInput) (*file.PutObjectAclOutput, error) {
	cli, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	err = bucket.SetObjectACL(req.Key, oss.ACLType(req.Acl))
	output := &file.PutObjectAclOutput{}
	return output, err
}
func (a *AliyunOSS) ListObjects(ctx context.Context, req *file.ListObjectsInput) (*file.ListObjectsOutput, error) {
	cli, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.ListObjects()
	if err != nil {
		return nil, err
	}
	out := &file.ListObjectsOutput{
		CommonPrefixes: resp.CommonPrefixes,
		Delimiter:      resp.Delimiter,
		IsTruncated:    resp.IsTruncated,
		Marker:         resp.Marker,
		MaxKeys:        int32(resp.MaxKeys),
		NextMarker:     resp.NextMarker,
		Prefix:         resp.Prefix,
	}
	for _, v := range resp.Objects {
		object := &file.Object{
			ETag:         v.ETag,
			Key:          v.Key,
			LastModified: v.LastModified.Unix(),
			Owner:        &file.Owner{ID: v.Owner.ID, DisplayName: v.Owner.DisplayName},
			Size:         v.Size,
			StorageClass: v.StorageClass,
		}
		out.Contents = append(out.Contents, object)
	}
	return out, nil
}
func (a *AliyunOSS) CopyObject(ctx context.Context, req *file.CopyObjectInput) (*file.CopyObjectOutput, error) {
	cli, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.CopyObject(req.CopySource.CopySourceKey, req.Key, VersionId(req.CopySource.CopySourceVersionId))
	if err != nil {
		return nil, err
	}
	out := &file.CopyObjectOutput{CopyObjectResult: &file.CopyObjectResult{ETag: resp.ETag, LastModified: resp.LastModified.Unix()}}
	return out, err
}

func (a *AliyunOSS) CreateMultipartUpload(ctx context.Context, req *file.CreateMultipartUploadInput) (*file.CreateMultipartUploadOutput, error) {
	cli, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.InitiateMultipartUpload(req.Key)
	output := &file.CreateMultipartUploadOutput{Bucket: resp.Bucket, Key: resp.Key, UploadId: resp.UploadID}
	return output, err
}
func (a *AliyunOSS) UploadPart(ctx context.Context, req *file.UploadPartInput) (*file.UploadPartOutput, error) {
	cli, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.UploadPart(
		oss.InitiateMultipartUploadResult{Bucket: req.Bucket, Key: req.Key, UploadID: req.UploadId},
		req.DataStream,
		req.ContentLength,
		int(req.PartNumber))
	output := &file.UploadPartOutput{ETag: resp.ETag}
	return output, err
}
func (a *AliyunOSS) UploadPartCopy(ctx context.Context, req *file.UploadPartCopyInput) (*file.UploadPartCopyOutput, error) {
	cli, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.UploadPartCopy(
		oss.InitiateMultipartUploadResult{Bucket: req.Bucket, Key: req.Key, UploadID: req.UploadId},
		req.CopySource.CopySourceBucket,
		req.CopySource.CopySourceKey,
		req.StartPosition,
		req.PartSize,
		int(req.PartNumber),
		VersionId(req.CopySource.CopySourceVersionId),
	)
	output := &file.UploadPartCopyOutput{CopyPartResult: &file.CopyPartResult{ETag: resp.ETag}}
	return output, err
}
func (a *AliyunOSS) CompleteMultipartUpload(ctx context.Context, req *file.CompleteMultipartUploadInput) (*file.CompleteMultipartUploadOutput, error) {
	cli, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}

	parts := make([]oss.UploadPart, 0)
	if req.MultipartUpload != nil {
		for _, v := range req.MultipartUpload.Parts {
			part := oss.UploadPart{PartNumber: int(v.PartNumber), ETag: v.ETag}
			parts = append(parts, part)
		}
	}
	resp, err := bucket.CompleteMultipartUpload(
		oss.InitiateMultipartUploadResult{Bucket: req.Bucket, Key: req.Key, UploadID: req.UploadId},
		parts,
	)
	output := &file.CompleteMultipartUploadOutput{Location: resp.Location, Bucket: resp.Bucket, Key: resp.Key, ETag: resp.ETag}
	return output, err
}
func (a *AliyunOSS) AbortMultipartUpload(ctx context.Context, req *file.AbortMultipartUploadInput) (*file.AbortMultipartUploadOutput, error) {
	cli, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}

	err = bucket.AbortMultipartUpload(
		oss.InitiateMultipartUploadResult{Bucket: req.Bucket, Key: req.Key, UploadID: req.UploadId},
	)
	output := &file.AbortMultipartUploadOutput{}
	return output, err
}
func (a *AliyunOSS) ListMultipartUploads(ctx context.Context, req *file.ListMultipartUploadsInput) (*file.ListMultipartUploadsOutput, error) {
	cli, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.ListMultipartUploads(Prefix(req.Prefix), KeyMarker(req.KeyMarker), MaxUploads(int(req.MaxUploads)), Delimiter(req.Delimiter), UploadIDMarker(req.UploadIdMarker))
	output := &file.ListMultipartUploadsOutput{
		Bucket:             resp.Bucket,
		Delimiter:          resp.Delimiter,
		Prefix:             resp.Prefix,
		KeyMarker:          resp.KeyMarker,
		UploadIdMarker:     resp.UploadIDMarker,
		NextKeyMarker:      resp.NextKeyMarker,
		NextUploadIdMarker: resp.NextUploadIDMarker,
		MaxUploads:         int32(resp.MaxUploads),
		IsTruncated:        resp.IsTruncated,
		CommonPrefixes:     resp.CommonPrefixes,
	}
	for _, v := range resp.Uploads {
		upload := &file.MultipartUpload{Initiated: v.Initiated.Unix(), UploadId: v.UploadID, Key: v.Key}
		output.Uploads = append(output.Uploads, upload)
	}
	return output, err
}

//未测试
func (a *AliyunOSS) RestoreObject(ctx context.Context, req *file.RestoreObjectInput) (*file.RestoreObjectOutput, error) {
	cli, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	err = bucket.RestoreObject(req.Key)
	output := &file.RestoreObjectOutput{}
	return output, err
}

//当前不允许该操作
func (a *AliyunOSS) ListObjectVersions(ctx context.Context, req *file.ListObjectVersionsInput) (*file.ListObjectVersionsOutput, error) {
	cli, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.ListObjectVersions()
	output := &file.ListObjectVersionsOutput{
		Name:                resp.Name,
		Prefix:              resp.Prefix,
		KeyMarker:           resp.KeyMarker,
		VersionIdMarker:     resp.VersionIdMarker,
		MaxKeys:             int32(resp.MaxKeys),
		Delimiter:           resp.Delimiter,
		IsTruncated:         resp.IsTruncated,
		NextKeyMarker:       resp.NextKeyMarker,
		NextVersionIdMarker: resp.NextVersionIdMarker,
		CommonPrefixes:      resp.CommonPrefixes,
	}
	for _, v := range resp.ObjectDeleteMarkers {
		marker := &file.DeleteMarkerEntry{
			IsLatest:     v.IsLatest,
			Key:          v.Key,
			LastModified: v.LastModified.Unix(),
			Owner: &file.Owner{
				ID:          v.Owner.ID,
				DisplayName: v.Owner.DisplayName,
			},
			VersionId: v.VersionId,
		}
		output.DeleteMarkers = append(output.DeleteMarkers, marker)
	}

	for _, v := range resp.ObjectVersions {
		version := &file.ObjectVersion{
			ETag:         v.ETag,
			IsLatest:     v.IsLatest,
			Key:          v.Key,
			LastModified: v.LastModified.Unix(),
			Owner: &file.Owner{
				ID:          v.Owner.ID,
				DisplayName: v.Owner.DisplayName,
			},
			Size:         v.Size,
			StorageClass: v.StorageClass,
			VersionId:    v.VersionId,
		}
		output.Versions = append(output.Versions, version)
	}

	return output, err
}
