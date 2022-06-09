package aliyun

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

// Prefix is an option to set prefix parameter
func Prefix(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.Prefix(value)
}

func KeyMarker(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.KeyMarker(value)
}

func MaxUploads(value int) oss.Option {
	if value <= 0 {
		return nil
	}
	return oss.MaxUploads(value)
}

func Delimiter(value string) oss.Option {
	if value == "" {
		return nil
	}

	return oss.Delimiter(value)
}
func UploadIDMarker(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.UploadIDMarker(value)
}

func VersionId(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.VersionId(value)
}
