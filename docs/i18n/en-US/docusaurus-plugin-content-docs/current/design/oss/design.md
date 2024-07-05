# Layotto object storage (OSS) and collection of API interface definitions and design

## Background

In order for layotto support object storage capability, an abstract interface of oss is required.Abstracts of interfaces need to satisfy the theoretical “portable” and to give the interface a clear terminology.

## Interface Design

The design of the entire interface will follow the following principle：

1. **Semantical, i.e. abstract interfaces have a clear semantic meaning.**
2. **Functional completeness, i.e. abstract interfaces need to satisfy different kinds of abilities as much as possible.**
3. **Maximum portability, i.e. abstract interfaces need to meet the transplantability requirement to the maximum extent possible.**

These principles are designed at a high and low level of priority.In order to meet the above requirements, there may be the following question：

1. The field is redundant, entry and exit may have fields corresponding to specific manufacturers.
2. Some of the interfaces may be supported only by some of the oss manufacturers, i.e. “maximum portability possible”.

## Configure Module Design

The fields of the original configuration module are abstract as shown below in：

```go
/OssMetadata wraps the configuration of loss implementation
type OssMetadata structure LO
	Endpoint string `json:`endpoint`
	AccessKeyID string`json:`accessKeyID``
	AccessKeySecret string `json:"accessKeySecret``
	Region string `json:"region"
}
```

EndPoint, AccessKeyID, AccessKeySecretariat and Regions are all of the concepts available to them and are not explained at length here.

## Interface Design

The definition of the interface is primarily divided into one of the main fields：

1. Generic interfaces are those that are similar to those supported by all of PutObject, GetObject and others.
2. Non-common interfaces are those that are only partially supported by the oss service.For example, ListParts interfaces, aws cannot support.

The design of the interface primarily refers to the aliyun and aws and the minio interface definitions.

> https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/API/API_GetObject.html
> https://help.aliyun.com/document_detail/31980.html
> https://docs.min.io/docs/golang-client-api-reference.html

### PutObject

Object upload interface, used as file upload, is the most basic ability of the goss.

> [https://docs.aws.amazon.com/en_cn/AmazonS3/latest/API/API_PutObject.html](https://docs.aws.amazon.com/en_cn/AmazonS3/latest/API/API_PutObject.html)
> 
> [https://help.aliyun.com/document_detail/31978.html](https://help.aliyun.com/document_detail/31978.html)

### GetObject

Object download interface, used as file downloads, is the most basic ability of oss.

> [https://docs.aws.amazon.com/en_cn/AmazonS3/latest/API/API_GetObject.html](https://docs.aws.amazon.com/en_cn/AmazonS3/latest/API/API_GetObject.html)
> 
> [https://help.aliyun.com/document_detail/31980.html](https://help.aliyun.com/document_detail/31980.html)

### DeleteObject

Object delete interface, used as file deletion, is the most basic ability of the goss.For the definition of the interface please find the link or pb definition above.

### PutObjectTagging

Label an object as the most basic ability of oss.For the definition of the interface please find the link or pb definition above.

### DeleteObjectTagging

Removing the tag of an object is the most basic ability of the oss.For the definition of the interface please find the link or pb definition above.

### GetObjectTagging

Get the tag of the object, which is the most basic ability of the oss.For the definition of the interface please find the link or pb definition above.

### CopyObject

Duplicate an existing object, is the most basic ability of the mass.For the definition of the interface please find the link or pb definition above.

### DeleteObjects

Deleting multiple objects, is the most basic ability of the mass.For the definition of the interface please find the link or pb definition above.

### ListObjects

Query all objects below bucket, supporting pagination queries, is the most basic ability of the mass.For the definition of the interface please find the link or pb definition above.

### GetObjectCannedAcl

Reading the object's announced acl. Users can set an object's action to control the object's access. First issue is the first issue. The following question： needs to be considered when designing the interface.

1. Whether to allow users to set object's action by api.
2. Whether the announced acl is portable

For the first issue, Aliyun is allowed to specify an object's action when uploading and to modify the object dynamically at any time.

---

**Aliyun acl defines the following：**

| Right Limit       | Permissions Description                                                                                                                                          |
| ----------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| public-read-write | Public read/write and write to：anyone (including anonymous visitors).                                                         |
| Public-read       | Public reading：can only be written by the owner of the Object, and anyone (including anonymous visitors) can read the object. |
| Private           | Private：only the owner of an object can read and write the object. No one else can access the object.                            |
| default           | Default：this object follows the write permissions of Bucket, which permissions are for Bucket, what permissions are for the Object.              |

---

**Tencent cloud defined below：**

| Right Limit               | Permissions Description                                                                                                                                                                             |
| ------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| default                   | Empty description, at this time you determine whether requests are allowed (default) based on the explicit settings of the various levels of directories and the bucket settings |
| Private                   | Has FULL_CONTROL permissions for the creator (master account). Others don't have permission                                                 |
| Public-read               | FULL_CONTROL permissions for creator, READ permissions for anonymous user groups                                                                                               |
| authenticated - read      | FULL_CONTROL permissions for creator, READ permissions for authentication user groups                                                                                          |
| bucket-owner-read         | FULL_CONTROL permissions for creator, READ permissions for bucket owners                                                                                                       |
| bucket-owner-full-control | Has both creator and bucket owner FULL_CONTROL permissions                                                                                                                     |

**Description：**
objects do not support granting public-read-write permissions.

---

**aws defined below：**

| **Canned ACL**            | **Applies to**    | **Permissions added to ACL**                                                                                                                                                                                                                             |
| ------------------------- | ----------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Private                   | Bucket and object | Owner gets FULL_CONTROL. No one else has access rights (default).                                                                                                                |
| Public-read               | Bucket and object | Owner gotsFULL_CONTROL. The AllUser group (see [Who is a grante?](https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl-overview.html#specifying-grante))) gets READ access. |
| public-read-write         | Bucket and object | Owner getsFULL_CONTROL. The AllUsers group gets READ and WRITE. Granting this on a bucket is generally not recommended.                                                             |
| aws-exec-read             | Bucket and object | Owner gets FULL_CONTROL. Amazon EC2 gets READ access to GET an Amazon Machine Image (AMI) bundle from Amazon S3.                                                                 |
| authenticated - read      | Bucket and object | Owner gets FULL_CONTROL. The AuthenticatedUser group gets READ access.                                                                                                                              |
| bucket-owner-read         | Object            | Object owner gets FULL_CONTROL. Bucket owner gets READ access. If you specify this ACL when creating a bucket, Amazon S3 ignores it.                                                |
| bucket-owner-full-control | Object            | Both the object owner and the bucket owner get FULL_CONTROL over the object. If you specify this announced ACL when creating a bucket, Amazon S3 ignores it.                                        |
| log-delivery-write        | Bucket            | For more information about logs, see ([Logging requests using server access loging](https://docs.aws.amazon.com/AmazonS3/latest/userLogs.html)).                                                                      |

---

As can be seen from the above list, the definition of acl is differentiated by different oss-manufacturers, but the concept of canned acl exists and therefore the interface belongs to an interface that does not guarantee portability,
this requires a judgement on the value of acl in the implementation of a specific component.For example, during the migration of users from Tencent to Aliyun
to assign ACls to public-read-write and return to Aliyun like “not supported acl”,
as long as it can be done to remind users that the interface does not fit the portable.

**Note: Layotto offers acl, but users need to be cautious about the use of acl, because differences in the service may lead to unanticipated results.**

> [https://help.aliyun.com/document_detail/100676.html] (https://help.aliun.com/document_detail/100676.html) Aliyun object acl type\
> [https://cloud.tencent.com/document/product/436/30752#E9.A2.84.E8.E7.E7.9A.84-acl] (https://cloud.tenent.com/document/product/436/30752#E9.A2.84.E.84.E.8.BE.7.9A.84-acl) Tencast cloud acltype\
> [https://docs.aws.amazon.com/AmazonS3/useruide/acl-overview.html#CannedACL](https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl-overview.html#CannedACL)
> [https://github.com/minio/minio/issues/8195](https://github.com/minio/minio/issues/8195) 对于minio是否应该支持acl的讨论

### PutObjectCannedAcl

This corresponds to the above-mentioned GetObjectCannedAcl, which is used to set an object.

**Note: Layotto offers acl, but users need to be cautious about the use of acl, because differences in the service may lead to unanticipated results.**

### RestoreObject

Launch the RestoreObject interface to unfreeze files of Archive Type (Archive) or Cold Archive.For the definition of the corresponding interface,
please find in the above link or in the reference to the pb interface definition annotation.

### CreateMultipartUpload

Creating a split upload interface, is the most basic ability of oss.For the definition of the corresponding interface, please find in the above link or in the reference to the pb interface definition annotation.

### Upload Part

Snippet upload interface, is the most basic ability of oss.For the definition of the corresponding interface, please find in the above link or in the reference to the pb interface definition annotation.

### Upload PartCopy

The decimal copy interface, is the most basic ability of oss.For the definition of the corresponding interface, please find in the above link or in the reference to the pb interface definition annotation.

### CompleteMultipartUpload

Finish the split upload interface, which is the most basic ability of oss.For the definition of the corresponding interface, please find in the above link or in the reference to the pb interface definition annotation.

### AbortMultipartUpload

Disconnect the particle upload interface, which is the most basic ability of oss.For the definition of the corresponding interface, please find in the above link or in the reference to the pb interface definition annotation.

### ListMultipartUpload

Query for uploaded fragments is the most basic ability of oss.For the definition of the corresponding interface, please find in the above link or in the reference to the pb interface definition annotation.

### ListObjectVersions

Query all version information for the object is the most basic ability of the goss.For the definition of the corresponding interface, please find in the above link or in the reference to the pb interface definition annotation.

### HeadObject

Returns the metadata data of the object, which is the most basic ability of the goss.For the definition of the corresponding interface, please find in the above link or in the reference to the pb interface definition annotation.

### IsObjectExist

The interface is not clearly defined in s3, and users can determine whether an object exists by way of the http's standard error code returned by HeadObject, which is abstracted separately to make the interface more semicolon.

> [http://biercoff.com/how-to-check-with-aws-cli-if-file-exiss-in-s3/](ttp://biercoff.com/how-to-check-with-aws-cli-if-file-exists-in-s3/)
> [https://stackoverflow.com/questions/41871948/aws-s3-how-to-check-if-a-file-exiss-in-a-bucket-using-bash](https://stackoverflow.com/questions/41871948/aws-on-to-the-file-existing-bash)

### SignURL

This interface generates an url for object upload and downloads, mainly for unauthorized users.It is the most basic ability of oss.For the definition of the corresponding interface,
please find in the above link or in the reference to the pb interface definition annotation.

### UpdateDownloadBandwidthRateLimit

This interface provides an interface for Aliyun that can limit the download speed of customers.For specific information, reference may be made to annotations in the pb definition.

### UpdateBandwidthRatteLimit

This interface provides an interface for Aliyun that can limit the upload speed of customers.For specific information, reference may be made to annotations in the pb definition.

### AppendObject

The interface is an additional interface mainly used for the application of the file, which is not supported by aw but is provided by both Aliyun and Tencent clouds and minio.

> https://help.aliyun.com/document_detail/31981.html
> https://github.com/minio/minio-java/issues/980

### ListParts

Query for uploaded fragments is the most basic ability of oss.For the definition of the corresponding interface, please find in the above link or in the reference to the pb interface definition annotation.
