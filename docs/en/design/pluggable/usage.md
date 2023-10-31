# Pluggable Component Usage Guide

This example demonstrates how users can implement and register their own components through the pluggable component capabilities provided by Layotto.
And verify the correctness of your component writing through Layotto SDK calls.

## step1.Write and run pluggable components

Next, run the already written code

```shell
cd demo/pluggable/hello
go run .
```

Printing the following result indicates successful service startup

```shell
start grpc server
```

>1. Taking the go implementation of the `hello` component as an example, find the corresponding component's proto file in `layotto/spec/proto/pluggable` and generate the corresponding implementation language's grpc file.
The Go language's pb file has been generated and placed under `spec/proto/pluggable/v1`, and users can directly reference it when using it.
>2. In addition to implementing the interfaces defined in the protobuf file, the component also needs to use socket mode to start the file and store the socket file in the default path of `/tmp/runtime/component-sockets`,
You can also use the environment variable `LAYOTTO_COMPONENTS_SOCKETS_FOLDER` modify the sock storage path location.
>3. In addition, users also need to register the reflection service in the GRPC server, which is used to obtain the specific implementation interface spec of the GRPC service during layotto service discovery. For specific code, please refer to `demo/pluggable/hello/main.go`

## step2. Launch Layotto

```shell
cd cmd/layotto
go build -o layotto .
./layotto start -c ../../configs/config_hello_component.json
```

> The type of the component filled in the configuration file is `hello-grpc-demo`, which is determined by the prefix name of the socket file. 
> The configuration items are consistent with registering regular hello components.
> Provide metadata items for users to set custom configuration requirements.

## step3. Component verification

Based on existing component testing code, test the correctness of user implemented pluggable components.

```shell
cd demo/hello/common
go run . -s helloworld
```

The program outputs the following results to indicate successful registration and operation of pluggable components.

```shell
runtime client initializing for: 127.0.0.1:34904
hello
```

## Understand the implementation principle of Layotto pluggable components

If you are interested in the implementation principles or want to extend some functions, you can read the [Design Document for Pluggable Components](en/design/pluggable/design.md)