# File API 设计文档

## API

API定义主要依据常用的文件操作来定义的，分为增删改查四个接口，对于Get/Put接口来说，文件的上传和下载需要支持流传输。因此接口定义如下：

```protobuf
  // Get file with stream
  rpc GetFile(GetFileRequest) returns (stream GetFileResponse) {}

  // Put file with stream
  rpc PutFile(stream PutFileRequest) returns (google.protobuf.Empty) {}

  // List all files
  rpc ListFile(ListFileRequest) returns (ListFileResp){}

  // Delete specific file
  rpc DelFile(DelFileRequest) returns (google.protobuf.Empty){}
```

关于接口的定义的讨论可以参照[issue98](https://github.com/mosn/layotto/issues/98)


## 参数定义


```protobuf
message GetFileRequest {
  //
  string store_name = 1;
  // The name of the file or object want to get.
  string name = 2;
  // The metadata for user extension.
  map<string,string> metadata = 3;
}

message GetFileResponse {
  bytes data = 1;
}

message PutFileRequest {
  string store_name = 1;
  // The name of the file or object want to put.
  string name = 2;
  // The data will be store.
  bytes data = 3;
  // The metadata for user extension.
  map<string,string> metadata = 4;
}

message FileRequest {
  string store_name = 1;
  // The name of the directory
  string name = 2;
  // The metadata for user extension.
  map<string,string> metadata = 3;
}

message ListFileRequest {
  FileRequest request = 1;
}

message ListFileResp {
  repeated string file_name = 1;
}

message DelFileRequest {
  FileRequest request = 1;
}
```

### Get接口

Get的入参主要有三个：

| **参数名** | **意义** | **是否必传** |
| --- | --- | --- | --- | --- | --- | --- |
| store_name | 后端对应的components（eg: aliyun.oss, aws.s3） | yes |
| name | 文件名字 | yes|
| metadata | 元数据，该字段用户可以用来指定component需要的一些字段，（eg:权限，用户名等） | yes|

### Put接口

Put接口入参主要有三个，多了一个data字段用来传输文件内容：

| **参数名** | **意义** | **是否必传** |
| --- | --- | --- | --- | --- | --- | --- |
| store_name | 后端对应的components（eg: aliyun.oss, aws.s3） | yes |
| name | 文件名字 | yes|
| data | 文件内容 | no（允许用户上传空数据，每个component可以做具体实现）|
| metadata | 元数据，该字段用户可以用来指定component需要的一些字段，（eg:权限，用户名等） | yes|


### List和Del接口

两个接口的参数是一样的：

| **参数名** | **意义** | **是否必传** |
| --- | --- | --- | --- | --- | --- | --- |
| store_name | 后端对应的components（eg: aliyun.oss, aws.s3） | yes |
| name | 文件名字 | yes|
| metadata | 元数据，该字段用户可以用来指定component需要的一些字段，（eg:权限，用户名等） | yes|

### 配置参数

配置参数，不同的component可以配置不同格式，比如aliyun.oss的配置如下：

```protobuf

{
    "file": {
      "file_demo": {
        "type": "aliyun.oss",
        "metadata":[
          {
            "endpoint": "endpoint_address",
            "accessKeyID": "accessKey",
            "accessKeySecret": "secret",
            "bucket": ["bucket1", "bucket2"]
          }
        ]
      }
    }
}

```