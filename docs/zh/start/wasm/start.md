## WASM概述

### 功能介绍

Layotto支持加载编译好的WASM文件，并通过`proxy_abi_version_0_2_0`版本的API与目标WASM进行交互。

### 快速开始

1. 启动layotto
```
go build -tags wasmer -o ./layotto ./cmd/layotto/main.go
./layotto start -c ./demo/wasm/config.json
```

2. 发送请求
```
curl -H 'name:Layotto' localhost:2045
Hi, Layotto
```

### 示例介绍

工程里分别用golang跟rust开发了两个功能一致的wasm模块，它们的实现思路如下：
1. 通过`proxy_on_request_headers`接收HTTP请求
2. 从`proxy_get_header_map_pairs`中取出header中的name字段
3. 使用`proxy_call_foreign_function`向Layotto发起调用
4. 通过`proxy_set_buffer_bytes`把处理结果返回给调用端

golang源码路径：
```
layotto/demo/wasm/code/golang/
```

rust源码路径：
```
layotto/demo/wasm/code/rust/
```


### 说明

该功能目前仍处于试验阶段，社区里对于WASM跟宿主的交互API也不够统一，因此如果您有该模块的需求欢迎发表在issue区，我们一起建设WASM！