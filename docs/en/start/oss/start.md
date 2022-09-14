   
# ObjectStorageService API demo

This example shows how to invoke Layotto ObjectStorageService API.

ObjectStorageService is an abstraction for blob storage or so called "object storage", such as alibaba cloud OSS, such as AWS S3.
You invoke ObjectStorageService API to do some CRUD operations on your binary file, e.g. query my file, delete my file, etc.

## step 0. modify the configuration
Please modify the OSS configuration in the `configs/config_oss.json`

```json
"grpc_config": {
  "oss": {
    "oss_demo": {
      "type": "aws.oss",
      "metadata": {
        "basic_config":{
          "region": "your-oss-resource-region",
          "endpoint": "your-oss-resource-endpoint",
          "accessKeyID": "your-oss-resource-accessKeyID",
          "accessKeySecret": "your-oss-resource-accessKeySecret"
        }
      }
    }
  }
}
```

## step 1. Deploy Layotto
<!-- tabs:start -->
### **With Docker**
You can start Layotto with docker

```bash
docker run -v "$(pwd)/configs/config_oss.json:/runtime/configs/config.json" -d  -p 34904:34904 --name layotto layotto/layotto start
```

### **Compile locally (not for Windows)**
You can compile and run Layotto locally.

> [!TIP|label: Not for Windows users]
> Layotto fails to compile under Windows. Windows users are recommended to deploy using docker-compose

After downloading the project code to the local, switch the code directory and compile:

```shell
cd ${project_path}/cmd/layotto
```

```shell @if.not.exist layotto
go build
```

Once finished, the layotto binary will be generated in the directory.

Run it:

```shell @background
./layotto start -c ../../configs/config_oss.json
```

<!-- tabs:end -->

## step 2. Run the client program to invoke Layotto ObjectStorageService API
<!-- tabs:start -->
### **Go**
Build and run the golang [demo](https://github.com/mosn/layotto/blob/main/demo/oss/client.go) :

```shell
cd ${project_path}/demo/oss/
go build client.go

# upload test3.txt with content "hello" to the bucket named `antsys-wenxuwan`
./client put antsys-wenxuwan test3.txt "hello"

# get test3.txt in the bucket antsys-wenxuwan
./client get antsys-wenxuwan test3.txt

# delete test3.txt
./client del antsys-wenxuwan test3.txt

# list the files in the bucket antsys-wenxuwan
./client list antsys-wenxuwan

```

### **Java**
<!-- 

Download java sdk and examples:

```shell @if.not.exist java-sdk
git clone https://github.com/layotto/java-sdk
```

```shell
cd java-sdk
```

Build the demo:

```shell @if.not.exist examples-oss/target/examples-oss-jar-with-dependencies.jar
# build example jar
mvn -f examples-oss/pom.xml clean package
```

Run it:

```shell
java -jar examples-oss/target/examples-oss-jar-with-dependencies.jar
```

If the following information is printed, the demo is successful:

```bash
TODO
```

-->

<!-- tabs:end -->

## step 3. Stop containers and release resources
<!-- tabs:start -->
### **Destroy the Docker container**
If you started Layotto with docker, you can destroy the container as follows:

```bash
docker rm -f layotto
```

<!-- tabs:end -->

## Next step
### What does this client program do?
The demo client program uses the SDK provided by Layotto to invoke the Layotto ObjectStorageService API.

The golang sdk is located in the `sdk` directory, and the java sdk is in https://github.com/layotto/java-sdk

In addition to using sdk, you can also interact with Layotto directly through grpc in any language you like.

### Details later, let's continue to experience other APIs
Explore other Quickstarts through the navigation bar on the left.

### Reference

[API reference](https://mosn.io/layotto/api/v1/s3.html)

[Design doc of ObjectStorageService API ](zh/design/oss/design)

 <!-- end services -->

