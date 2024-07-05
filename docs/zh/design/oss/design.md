# Layotto对象存储（OSS）并集API接口定义及设计

## 背景

为了让layotto支持对象存储能力，需要对oss的接口进行抽象。抽象出的接口需要满足理论上的”可移植性“以及让接口具有明确的语义。

## 接口设计

整个接口的设计将遵循以下原则：

1. **语义性，即抽象出的接口具有明确的语义。**
2. **功能完整性，即抽象出的接口需要尽可能满足不同oss的能力。**
3. **最大可移植性，即抽象出的接口需要尽最大可能的满足可移植性的要求。**

上述原则设计的时候考量的优先级从高到低。为了满足上述的要求，可能会存在以下问题：

1. 字段的冗余，入参和出参可能会存在对应特定厂商的字段。 
2. 部分接口可能只在部分oss厂商上可以支持，即“最大可能移植性”。

## 配置模块设计

oss原始配置模块的字段抽象如下所示：

```go
// OssMetadata wraps the configuration of oss implementation
type OssMetadata struct {
	Endpoint        string   `json:"endpoint"`
	AccessKeyID     string   `json:"accessKeyID"`
	AccessKeySecret string   `json:"accessKeySecret"`
	Region          string   `json:"region"`
}
```

Endpoint、AccessKeyID、AccessKeySecret、Region是现有的oss都有的概念，本文不做过多解释。


## 接口设计

接口的定义主要分为两类：

1. 通用接口，即类似于PutObject、GetObject等所有的oss服务都支持的接口。 
2. 非通用接口，即只有部分oss服务支持的接口。比如ListParts接口，aws就无法支持。

该接口的设计主要参考aliyun和aws以及minio的接口定义。
> https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/API/API_GetObject.html     
> https://help.aliyun.com/document_detail/31980.html    
> https://docs.min.io/docs/golang-client-api-reference.html    

### PutObject

对象上传接口，用作上传文件，是oss最基本能力。

> [https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/API/API_PutObject.html](https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/API/API_PutObject.html)    
> [https://help.aliyun.com/document_detail/31978.html](https://help.aliyun.com/document_detail/31978.html)

### GetObject

对象下载接口，用作文件下载，是oss最基本能力。

> [https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/API/API_GetObject.html](https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/API/API_GetObject.html)    
> [https://help.aliyun.com/document_detail/31980.html](https://help.aliyun.com/document_detail/31980.html)


### DeleteObject

对象删除接口，用作文件删除，是oss最基本能力。对应接口的定义，请在上述的链接或者pb定义中查询。

### PutObjectTagging

给对象添加标签，是oss最基本能力。对应接口的定义，请在上述的链接或者pb定义中查询。

### DeleteObjectTagging

删除对象的标签，是oss最基本能力。对应接口的定义，请在上述的链接或者pb定义中查询。

### GetObjectTagging

获取对象的标签，是oss最基本能力。对应接口的定义，请在上述的链接或者pb定义中查询。

### CopyObject

复制已经存在的object，是oss最基本能力。对应接口的定义，请在上述的链接或者pb定义中查询。

### DeleteObjects

删除多个object，是oss最基本能力。对应接口的定义，请在上述的链接或者pb定义中查询。

### ListObjects

查询bucket下面的所有的objects，支持分页查询，是oss最基本能力。对应接口的定义，请在上述的链接或者pb定义中查询。

### GetObjectCannedAcl

读取对象的canned acl，用户可以设置object的acl来控制object的访问权限，首先第一个问题，设计该接口时需要考虑以下问题：

1. 是否允许用户通过api来设置object的acl。 
2. oss的canned acl是否是具有可移植性的

对于第一个问题，阿里云是允许用户上传时指定object的acl的，同时也允许随时动态的修改object的acl，腾讯云也允许此操作。

---

**阿里云的acl定义如下：**

| 权限值 | 权限描述 |
| --- | --- |
| public-read-write | 公共读写：任何人（包括匿名访问者）都可以对该Object进行读写操作。 |
| public-read | 公共读：只有该Object的拥有者可以对该Object进行写操作，任何人（包括匿名访问者）都可以对该Object进行读操作。 |
| private | 私有：只有Object的拥有者可以对该Object进行读写操作，其他人无法访问该Object。 |
| default | 默认：该Object遵循Bucket的读写权限，即Bucket是什么权限，Object就是什么权限。 |


---

**腾讯云的定义如下：**

| 权限值 | 权限描述 |
| --- | --- |
| default | 空描述，此时根据各级目录的显式设置及存储桶的设置来确定是否允许请求（默认） |
| private | 创建者（主账号）具备 FULL_CONTROL 权限，其他人没有权限 |
| public-read | 创建者具备 FULL_CONTROL 权限，匿名用户组具备 READ 权限 |
| authenticated-read | 创建者具备 FULL_CONTROL 权限，认证用户组具备 READ 权限 |
| bucket-owner-read | 创建者具备 FULL_CONTROL 权限，存储桶拥有者具备 READ 权限 |
| bucket-owner-full-control | 创建者和存储桶拥有者都具备 FULL_CONTROL 权限 |

**说明：**
对象不支持授予 public-read-write 权限。

---

**aws定义如下：**

| **Canned ACL** | **Applies to** | **Permissions added to ACL** |
| --- | --- | --- |
| private | Bucket and object | Owner gets FULL_CONTROL. No one else has access rights (default). |
| public-read | Bucket and object | Owner gets FULL_CONTROL. The AllUsers group (see [Who is a grantee?](https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl-overview.html#specifying-grantee)) gets READ access. |
| public-read-write | Bucket and object | Owner gets FULL_CONTROL. The AllUsers group gets READ and WRITE access. Granting this on a bucket is generally not recommended. |
| aws-exec-read | Bucket and object | Owner gets FULL_CONTROL. Amazon EC2 gets READ access to GET an Amazon Machine Image (AMI) bundle from Amazon S3. |
| authenticated-read | Bucket and object | Owner gets FULL_CONTROL. The AuthenticatedUsers group gets READ access. |
| bucket-owner-read | Object | Object owner gets FULL_CONTROL. Bucket owner gets READ access. If you specify this canned ACL when creating a bucket, Amazon S3 ignores it. |
| bucket-owner-full-control | Object | Both the object owner and the bucket owner get FULL_CONTROL over the object. If you specify this canned ACL when creating a bucket, Amazon S3 ignores it. |
| log-delivery-write | Bucket | The LogDelivery group gets WRITE and READ_ACP permissions on the bucket. For more information about logs, see ([Logging requests using server access logging](https://docs.aws.amazon.com/AmazonS3/latest/userguide/ServerLogs.html)). |

---

从上述列表可以看出，不同的oss厂商对于acl的定义是有区别的，但canned acl的概念是都存在的，因此该接口属于不保证可移植性的接口，
这就需要在具体的component的实现中对acl的值进行判断。例如用户从腾讯云迁移到阿里云的过程中，
如果指定了ACl为public-read-write，那么在迁移到阿里云的时候，需要返回类似于“not supported acl”，
只要可以做到提醒用户该接口不能满足可移植即可。

**Note: Layotto虽然提供了acl的操作，但用户对于acl的使用需要谨慎，因为服务端不同的差异可能会导致不可预期的结果。**



> [https://help.aliyun.com/document_detail/100676.html](https://help.aliyun.com/document_detail/100676.html)  阿里云object acl类型    
> [https://cloud.tencent.com/document/product/436/30752#.E9.A2.84.E8.AE.BE.E7.9A.84-acl](https://cloud.tencent.com/document/product/436/30752#.E9.A2.84.E8.AE.BE.E7.9A.84-acl) 腾讯云acl类型    
> [https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl-overview.html#CannedACL](https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl-overview.html#CannedACL)    
> [https://github.com/minio/minio/issues/8195](https://github.com/minio/minio/issues/8195) 对于minio是否应该支持acl的讨论    


### PutObjectCannedAcl

这个和上述的GetObjectCannedAcl相对应，用来设置object的canned acl。

**Note: Layotto虽然提供了acl的操作，但用户对于acl的使用需要谨慎，因为服务端不同的差异可能会导致不可预期的结果。**

### RestoreObject

调用RestoreObject接口解冻归档类型（Archive）或冷归档（Cold Archive）的文件（Object）。对应接口的定义，
请在上述的链接中或者pb接口定义注释的引用中查询。

### CreateMultipartUpload

创建分片上传接口，是oss最基本能力。对应接口的定义，请在上述的链接中或者pb接口定义注释的引用中查询。

### UploadPart

分片上传接口，是oss最基本能力。对应接口的定义，请在上述的链接中或者pb接口定义注释的引用中查询。

### UploadPartCopy

分片copy接口，是oss最基本能力。对应接口的定义，请在上述的链接中或者pb接口定义注释的引用中查询。

### CompleteMultipartUpload

完成分片上传接口，是oss最基本能力。对应接口的定义，请在上述的链接中或者pb接口定义注释的引用中查询。

### AbortMultipartUpload

中断分片上传接口，是oss最基本能力。对应接口的定义，请在上述的链接中或者pb接口定义注释的引用中查询。

### ListMultipartUploads

查询已经上传的分片，是oss最基本能力。对应接口的定义，请在上述的链接中或者pb接口定义注释的引用中查询。

### ListObjectVersions

查询对象所有的版本信息，是oss最基本能力。对应接口的定义，请在上述的链接中或者pb接口定义注释的引用中查询。

### HeadObject

返回object的metadata数据，是oss最基本能力。对应接口的定义，请在上述的链接中或者pb接口定义注释的引用中查询。

### IsObjectExist

该接口在s3中没有明确的定义，用户可以通过HeadObject返回的http标准错误码是不是404来判断object是不是存在，这里单独抽象出来是为了让接口更加具有语义信息。

> [http://biercoff.com/how-to-check-with-aws-cli-if-file-exists-in-s3/](http://biercoff.com/how-to-check-with-aws-cli-if-file-exists-in-s3/)      
> [https://stackoverflow.com/questions/41871948/aws-s3-how-to-check-if-a-file-exists-in-a-bucket-using-bash](https://stackoverflow.com/questions/41871948/aws-s3-how-to-check-if-a-file-exists-in-a-bucket-using-bash)

### SignURL

该接口会生成一个url用作object的上传和下载，主要用于未经授权的用户。是oss最基本能力。对应接口的定义，
请在上述的链接中或者pb接口定义注释的引用中查询。

### UpdateDownloadBandwidthRateLimit

该接口为阿里云提供的接口，可以限制client的下载速度。具体信息可参照pb定义里面的注释信息。

### UpdateUploadBandwidthRateLimit

该接口为阿里云提供的接口，可以限制client的上传速度。具体信息可参照pb定义里面的注释信息。

### AppendObject

该接口为追加接口，主要用于对文件进行append操作，aws不支持该操作，但阿里云和腾讯云以及minio都提供了对应的方式来实现。

> https://help.aliyun.com/document_detail/31981.html
> https://github.com/minio/minio-java/issues/980

### ListParts

查询已经上传的分片，是oss最基本能力。对应接口的定义，请在上述的链接中或者pb接口定义注释的引用中查询。
