# File API

## What is File API

File api is used to implement file operations. Applications can perform CRUD operations on files through this interface. The interface supports streaming mode to realize the transmission of large files.

## How to use File API
You can call the File API through grpc. The API is defined in [runtime.proto](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto).

The component needs to be configured before use. Different components should have own configuration.For OSS detail configuration items, see [OSS Component Document](en/component_specs/file/oss.md)

### Example

For examples of using file api, please refer to [File Demo](../../../../demo/file)


## grpc API definition

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

## Research

Refer：

```protobuf
https://github.com/mosn/layotto/issues/98
```

### parameters

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

### Get File

```protobuf
  // Get file with stream
  rpc GetFile(GetFileRequest) returns (stream GetFileResponse) {}
```

To avoid inconsistencies between this document and the code, please refer to [the newest proto file](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto) for detailed input parameters and return values.

### Put File

```protobuf
  // Put file with stream
  rpc PutFile(stream PutFileRequest) returns (google.protobuf.Empty) {}
```

To avoid inconsistencies between this document and the code, please refer to [the newest proto file](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto) for detailed input parameters and return values.

### Delete File

```protobuf
// Delete specific file
rpc DelFile(DelFileRequest) returns (google.protobuf.Empty){}
```

To avoid inconsistencies between this document and the code, please refer to [the newest proto file](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto) for detailed input parameters and return values.

### List File

```protobuf
// List all files
rpc ListFile(ListFileRequest) returns (ListFileResp){}
```

To avoid inconsistencies between this document and the code, please refer to [the newest proto file](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto) for detailed input parameters and return values.