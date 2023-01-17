package huaweiyun

import (
	"context"
	"encoding/json"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"mosn.io/layotto/components/oss"
	"mosn.io/layotto/components/pkg/utils"
	"strconv"
)

const connectTimeoutSec = "connectTimeoutSec"

type HuaweiyunOSS struct {
	client   *obs.ObsClient
	metadata utils.OssMetadata
}

func (h *HuaweiyunOSS) Init(ctx context.Context, config *oss.Config) error {
	connectTimeout := 30
	jsonRawMessage := config.Metadata[oss.BasicConfiguration]
	err := json.Unmarshal(jsonRawMessage, &h.metadata)
	if err != nil {
		return oss.ErrInvalid
	}
	if t, ok := config.Metadata[connectTimeoutSec]; ok {
		if v, err := strconv.Atoi(string(t)); err != nil {
			connectTimeout = v
		}
	}

	client, err := obs.New(h.metadata.AccessKeyID, h.metadata.AccessKeySecret, h.metadata.Endpoint, obs.WithConnectTimeout(connectTimeout))
	if err != nil {
		return err
	}
	h.client = client
	return nil
}

func (h *HuaweiyunOSS) GetObject(ctx context.Context, input *oss.GetObjectInput) (*oss.GetObjectOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) PutObject(ctx context.Context, input *oss.PutObjectInput) (*oss.PutObjectOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) DeleteObject(ctx context.Context, input *oss.DeleteObjectInput) (*oss.DeleteObjectOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) PutObjectTagging(ctx context.Context, input *oss.PutObjectTaggingInput) (*oss.PutObjectTaggingOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) DeleteObjectTagging(ctx context.Context, input *oss.DeleteObjectTaggingInput) (*oss.DeleteObjectTaggingOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) GetObjectTagging(ctx context.Context, input *oss.GetObjectTaggingInput) (*oss.GetObjectTaggingOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) CopyObject(ctx context.Context, input *oss.CopyObjectInput) (*oss.CopyObjectOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) DeleteObjects(ctx context.Context, input *oss.DeleteObjectsInput) (*oss.DeleteObjectsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) ListObjects(ctx context.Context, input *oss.ListObjectsInput) (*oss.ListObjectsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) GetObjectCannedAcl(ctx context.Context, input *oss.GetObjectCannedAclInput) (*oss.GetObjectCannedAclOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) PutObjectCannedAcl(ctx context.Context, input *oss.PutObjectCannedAclInput) (*oss.PutObjectCannedAclOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) RestoreObject(ctx context.Context, input *oss.RestoreObjectInput) (*oss.RestoreObjectOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) CreateMultipartUpload(ctx context.Context, input *oss.CreateMultipartUploadInput) (*oss.CreateMultipartUploadOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) UploadPart(ctx context.Context, input *oss.UploadPartInput) (*oss.UploadPartOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) UploadPartCopy(ctx context.Context, input *oss.UploadPartCopyInput) (*oss.UploadPartCopyOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) CompleteMultipartUpload(ctx context.Context, input *oss.CompleteMultipartUploadInput) (*oss.CompleteMultipartUploadOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) AbortMultipartUpload(ctx context.Context, input *oss.AbortMultipartUploadInput) (*oss.AbortMultipartUploadOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) ListMultipartUploads(ctx context.Context, input *oss.ListMultipartUploadsInput) (*oss.ListMultipartUploadsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) ListObjectVersions(ctx context.Context, input *oss.ListObjectVersionsInput) (*oss.ListObjectVersionsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) HeadObject(ctx context.Context, input *oss.HeadObjectInput) (*oss.HeadObjectOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) IsObjectExist(ctx context.Context, input *oss.IsObjectExistInput) (*oss.IsObjectExistOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) SignURL(ctx context.Context, input *oss.SignURLInput) (*oss.SignURLOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) UpdateDownloadBandwidthRateLimit(ctx context.Context, input *oss.UpdateBandwidthRateLimitInput) error {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) UpdateUploadBandwidthRateLimit(ctx context.Context, input *oss.UpdateBandwidthRateLimitInput) error {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) AppendObject(ctx context.Context, input *oss.AppendObjectInput) (*oss.AppendObjectOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) ListParts(ctx context.Context, input *oss.ListPartsInput) (*oss.ListPartsOutput, error) {
	//TODO implement me
	panic("implement me")
}
