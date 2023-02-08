package huaweiyun

import "errors"

var (
	ErrHaveNotTag                = errors.New("huaweiyun obs object haven't tagging feature")
	ErrDownloadNotBandwidthLimit = errors.New("huaweiyun obs download haven't bandwidth limit feature")
	ErrUploadNotBandwidthLimit   = errors.New("huaweiyun obs download haven't bandwidth limit feature")
	ErrNotSupportAclGet          = errors.New("huaweiyun obs haven't acl-get support")
)
