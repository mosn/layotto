   
# PhoneCallService API demo

This example shows how to invoke Layotto PhoneCallService API.

## What is PhoneCallService API used for?

PhoneCallService API is used to send voice messages.

## step 1. Deploy Layotto
<!-- tabs:start -->
### **With Docker**
You can start Layotto with docker

```bash
docker run -v "$(pwd)/configs/config_standalone.json:/runtime/configs/config.json" -d  -p 34904:34904 --name layotto layotto/layotto start
```

### **Compile locally (not for Windows)**
You can compile and run Layotto locally.

> [!TIP|label: Not for Windows users]
> Layotto fails to compile under Windows. Windows users are recommended to deploy using docker

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
./layotto start -c ../../configs/config_standalone.json
```

<!-- tabs:end -->

## step 2. Run the client program to invoke Layotto PhoneCallService API
<!-- tabs:start -->
### **Go**
Build and run the golang demo:

```shell
 cd ${project_path}/demo/phone/common/
 go build -o client
 ./client -s "demo"
```

If the following information is printed, the demo is successful:

```bash
TODO
```

### **Java**

Download java sdk and examples:

```shell @if.not.exist java-sdk
git clone https://github.com/layotto/java-sdk
```

```shell
cd java-sdk
```

Build the demo:

```shell @if.not.exist examples-phone/target/examples-phone-jar-with-dependencies.jar
# build example jar
mvn -f examples-phone/pom.xml clean package
```

Run it:

```shell
java -jar examples-phone/target/examples-phone-jar-with-dependencies.jar
```

If the following information is printed, the demo is successful:

```bash
TODO
```

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
The demo client program uses the SDK provided by Layotto to invoke the Layotto PhoneCallService API.

The golang sdk is located in the `sdk` directory, and the java sdk is in https://github.com/layotto/java-sdk

In addition to using sdk, you can also interact with Layotto directly through grpc in any language you like.

### Details later, let's continue to experience other APIs
Explore other Quickstarts through the navigation bar on the left.

### Reference

[API Reference](https://mosn.io/layotto/api/v1/phone.html)

<!--design_doc_url-->

 <!-- end services -->

