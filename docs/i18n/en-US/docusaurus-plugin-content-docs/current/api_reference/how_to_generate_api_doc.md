# How to generate code and documentation from the `.proto` files

Suppose you wrote a new proto file `spec/proto/extension/v1/email/email.proto` and you want to implement this API in Layotto:

```protobuf
// EmailService is used to send emails.
service EmailService {

  // Send an email with template
  rpc SendEmailWithTemplate(SendEmailWithTemplateRequest) returns (SendEmailWithTemplateResponse) {}

  // Send an email with raw content instead of using templates.
  rpc SendEmail(SendEmailRequest) returns (SendEmailResponse) {}

}

// different message types......
```

It's a tedious job because you have to write a lot of code and docs.

Fortunately, Layotto has tools to generate the code/docs/CI configuration automatically. You don't have to do the job yourself!

## step 1. Make sure your proto file meets the following requirements
- The file path should be `spec/proto/extension/v1/{api short name}/{api short name}.proto`
- There should be only one `service` in the proto file. For example, the following file is **WRONG** :

```protobuf
// EmailService is used to send emails.
service EmailService {
  // ...
}

// Wrong: there should be only one `service` in a `.proto` file
service EmailService2 {
  // ...
}

// different message types......
```

- If you don't want to generate the quickstart docs for the proto, add a comment `/* @exclude skip quickstart_generator */` .
- If you don't want to generate the sdk & sidecar code for the proto, add a comment `/* @exclude skip code_generator */` .

You can take the `spec/proto/extension/v1/s3/oss.proto` as an example:

```protobuf
/* @exclude skip quickstart_generator */
/* @exclude skip code_generator */
// ObjectStorageService is an abstraction for blob storage or so called "object storage", such as alibaba cloud OSS, such as AWS S3.
// You invoke ObjectStorageService API to do some CRUD operations on your binary file, e.g. query my file, delete my file, etc.
service ObjectStorageService{
  //......
}
```

These special comments are called "Master's commands". There are many other commands, and you can check [the doc](https://github.com/layotto/protoc-gen-p6#masters-commands) for more details.

## step 2. Check the environment
To run the generator, you need:
- Go version >=1.16
- Start Docker

## step 3. Generate everything
Note: the command below should be executed under layotto directory

```shell
make proto
```

Then you get:
- Generated code
  - `.pb.go` code
  - `_grpc.pb.go` code
  - layotto go-sdk code
  - layotto sidecar code
- Generated documentation
  - API reference docs
  - updated API reference list
  - quickstart document (both chinese and english)
  - updated sidebar (The tool will add the generated quickstart doc into the sidebar of https://mosn.io/layotto )
- Updated CI (The tool will add the generated quickstart doc into the CI script `etc/script/test-quickstart.sh`)

## step 4. Write the rest of the code
Now it's your job to implement:

- Layotto component
- go examples

![image](https://user-images.githubusercontent.com/26001097/188782762-bc1404a8-b891-45d3-a1ac-f86cafdbc0ab.png)

- java examples

![image](https://user-images.githubusercontent.com/26001097/188782989-9aec893f-9d12-4ee6-9a64-940b0ba1ba1b.png)

## Behind the scenes
We have a protoc plugin called [protoc-gen-p6](https://github.com/layotto/protoc-gen-p6) to generate code for Layotto.

## What if I want to generate pb/documentation only?
The steps above generate everything, but what if I only want to generate `.pb.go` code ? What if I only want to generate the docs?

### How to compile the proto files into `.pb.go` code
<!-- tabs:start -->
#### **Make cmmand(recommended)**

```bash
make proto-code
```

This command uses docker to run protoc and generate `.pb.go` code files.

#### **Install protoc**
1. Install protoc version: [v3.17.3](https://github.com/protocolbuffers/protobuf/releases/tag/v3.17.3)

2. Install protoc-gen-go v1.28 and protoc-gen-go-grpc v1.2

3. Generate gRPC `.pb.go` code

```bash
cd spec/proto/runtime/v1
protoc -I. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false,paths=source_relative *.proto
```

<!-- tabs:end -->
### How to generate API reference doc according to the proto files
We can use [protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc) and docker to generate api documents,the command is as follows:  

<!-- tabs:start -->
#### **Make command(recommended)**

```bash
make proto-doc
```

This command uses docker to run protoc-gen-doc and generate docs.

#### **Use docker to run protoc-gen-doc**
`make proto-doc` invokes the script `etc/script/generate-doc.sh`, which uses docker to run protoc-gen-doc.

You can check `etc/script/generate-doc.sh` for more details.

<!-- tabs:end -->
