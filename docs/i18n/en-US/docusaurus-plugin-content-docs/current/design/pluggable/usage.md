# Pluggable Component Use Documents

This example shows how users can implement and register their components with the ability provided by Layotto to plug components.and use Layotto sdk to verify the correctness of their components.

## Step1. Write and Run Plug Components

Next, run code already written

```shell
cd demo/pluggable/hello
go run.
```

Print the results below for successful service startup

```shell
start grpc server
```

> 1. Use go to implement hello components as an example, find proto file for corresponding components in `layotto/spec/proto/pluggable`, and generate grpc files for the corresponding language.
>    The pb file in the go language has been generated and placed under `spec/proto/plugable/v1` to use it directly when referenced by the user.
> 2. In addition to implementing the interface defined in the protocoluf file, the component needs to start the file using socket and store the sock file in the `/tmp/runtime/component-sockets` default path.
>    can also modify the location of the sock store path through the environmental variable `LAYOTTO_COMPONETS_SOCKETS_FOLDER`.
> 3. In addition to this, the user needs to register a reflection service in grpserver for the specific implementation interface of the grpc service when the layotto service is found. Specific code can be referenced to `demo/plugble/hello/main.go`

## Step2. Launch Layotto

```shell
cd cmd/layotto
go build -o layotto.
./layotto start -c ../../configs/config_hello_component.json
```

> The type of the component in the configuration file is `hello-grpc-demo`, determined by the prefix of the socket file. The configuration item is consistent with the registered normal hello component.Provides metadata items to allow users to set custom configuration requirements.

## Step 3. Validate component

Based on existing component test code, test the validity of plug-in plugins implemented by the user.

```shell
cd demo/hello/common
go run
```

The following result of the program output indicates that the unpluggable component registration is running successfully.

```shell
runtime client initializing for: 127.0.0.1:34904
hello
```

## Learn how the Layotto Plug Components can be implemented

If you are interested in implementing the rationale or want to expand some features, you can read[可插拔组件的设计文档](design/pluggable/design.md)
