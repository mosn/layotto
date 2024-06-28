# File API design documentation

## API

The API definition is mainly based on commonly used file operations and is divided into four additional and deleted interfaces. For the Get/Put interface, file upload and download require support for streaming transfer.Thus the following interface definition is：

```protobuf
  // Get file with stream
  rpc GetFile(GetFileRequest) returns (stream GetFileResponse) {}

  // Put file with stream
  rpc PutFile(stream PuteRequest) returns (google. Rotobuf. mpty) {}

  // List all files
  rpc ListFile(ListFileRequest) returns (ListFileResponse){}

  // Delete specific file
  rpc DelFile(DelFileRequest) returns (google. Rotobuf.Empty {}
```

Discussion of the definition of interfaces can be based on[issue98](https://github.com/mosn/layotto/issues/98)

## Definition of parameters

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

### Get Interface

Get's main entry is three：

\| **Parameter Name** | **Mean** |
\| ---- | -- | -- | -- | -- | -- -- | -- | -- -- | -- -- | -- -- | -- | -- |
\| store_name | corresponding backend components (eg: aliyun. ss, aws.s3) | yes |
\| name | file name | yes|
\| metadata | Metadata where users can specify some of the fields that component needs, (eg:per, username etc.) | yes|

### Put Interface

Put interfaces have three main interfaces and more than one data field is used to transfer file content：

\| **Parameter Name** | **Mean** |
\| ---- | -- | -- | -- | -- | -- -- | -- | -- -- | -- -- | -- -- | -- | -- |
\| store_name | corresponding backend components (eg: aliyun. ss, aws. 3) | yes |
\| name | yes|
\| data | file content | nos (allowing users to upload empty data, each component can be operationalized)|
\| metadata, which users can specify some of the fields that component needs, (eg:permission, username, etc.) | yes|

### Lists and Dels Interfaces

Parameters for both interfaces are the same：

\| **Parameter Name** | **Mean** |
\| ---- | -- | -- | -- | -- | -- -- | -- | -- -- | -- -- | -- -- | -- | -- |
\| store_name | corresponding backend components (eg: aliyun. ss, aws.s3) | yes |
\| name | file name | yes|
\| metadata | Metadata where users can specify some of the fields that component needs, (eg:per, username etc.) | yes|

### Configure Parameters

Configure parameters, different components can be configured in different formats such as aliyun.oss below：

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
