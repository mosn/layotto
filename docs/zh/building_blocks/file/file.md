# File API

## 什么是File API

File api 用于实现文件操作。应用程序可以通过该接口对文件进行CRUD操作。该接口支持大文件的流式传输。

## 如何使用File API

File api定义在 [runtime.proto](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto) 。应用可以通过grpc调用
对应的File api，实现文件的操作。

该接口在使用前需要进行配置。不同的文件系统可能配置不同，用户可根据自己的文件系统进行配置。比如OSS详细配置项可参考 [OSS组件文档](en/component_specs/file/oss.md)

### 例子

当前提供了基于layotto访问aliyun oss文件系统的示例，具体可参照 [文件演示](../../../../demo/file)

## 接口定义

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

## 调研和讨论

参考：

api的设计请参考下面的issue：

```protobuf
https://github.com/mosn/layotto/issues/98
```

### 接口参数

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

### 读文件

```protobuf
  // Get file with stream
  rpc GetFile(GetFileRequest) returns (stream GetFileResponse) {}
```

为避免文档和代码不一致，详细入参和返回值请参考 [the newest proto file](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto).

### 写文件

```protobuf
  // Put file with stream
  rpc PutFile(stream PutFileRequest) returns (google.protobuf.Empty) {}
```

为避免文档和代码不一致，详细入参和返回值请参考 [the newest proto file](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto).

### 删文件

```protobuf
// Delete specific file
rpc DelFile(DelFileRequest) returns (google.protobuf.Empty){}
```

为避免文档和代码不一致，详细入参和返回值请参考 [the newest proto file](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto).

### 查文件

```protobuf
// List all files
rpc ListFile(ListFileRequest) returns (ListFileResp){}
```

为避免文档和代码不一致，详细入参和返回值请参考 [the newest proto file](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto).
