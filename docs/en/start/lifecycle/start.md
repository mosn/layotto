   
# Lifecycle API demo

This example shows how to invoke Layotto Lifecycle API.

Lifecycle API is used to manage the sidecar lifecycle.
For example, by invoking the lifecycle API, you can modify the components' configuration during runtime

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

## step 2. Run the client program to invoke Layotto Lifecycle API
<!-- tabs:start -->
### **Go**
Build and run the golang demo:

```shell
 cd ${project_path}/demo/lifecycle/common/
 go build -o client
 ./client -s "demo"
```

If the following information is printed, the demo is successful:

```bash
TODO
```

### **Java**

[comment]: <> (Download java sdk and examples:)

[comment]: <> (```shell @if.not.exist java-sdk)

[comment]: <> (git clone https://github.com/layotto/java-sdk)

[comment]: <> (```)

[comment]: <> (```shell)

[comment]: <> (cd java-sdk)

[comment]: <> (```)

[comment]: <> (Build the demo:)

[comment]: <> (```shell @if.not.exist examples-lifecycle/target/examples-lifecycle-jar-with-dependencies.jar)

[comment]: <> (# build example jar)

[comment]: <> (mvn -f examples-lifecycle/pom.xml clean package)

[comment]: <> (```)

[comment]: <> (Run it:)

[comment]: <> (```shell)

[comment]: <> (java -jar examples-lifecycle/target/examples-lifecycle-jar-with-dependencies.jar)

[comment]: <> (```)

[comment]: <> (If the following information is printed, the demo is successful:)

[comment]: <> (```bash)

[comment]: <> (TODO)

[comment]: <> (```)

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
The demo client program uses the SDK provided by Layotto to invoke the Layotto Lifecycle API.

The golang sdk is located in the `sdk` directory, and the java sdk is in https://github.com/layotto/java-sdk

In addition to using sdk, you can also interact with Layotto directly through grpc in any language you like.

### Details later, let's continue to experience other APIs
Explore other Quickstarts through the navigation bar on the left.

### Reference

[API Reference](https://mosn.io/layotto/api/v1/runtime.html)

[Design doc](zh/design/lifecycle/apply_configuration)

 <!-- end services -->

