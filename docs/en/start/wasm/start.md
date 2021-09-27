## WASM on Layotto

### What is WASM on Layotto?

Layotto supports load the compiled WASM file, and interacts with it through the API of the `proxy_abi_version_0_2_0` version.

### Quick start

1. start Layotto server

```
go build -tags wasmer -o ./layotto ./cmd/layotto/main.go
./layotto start -c ./demo/wasm/config.json
```

2. send request

```
curl -H 'name:Layotto' -H 'id:id_1' localhost:2045
Hi, Layotto_id_1

curl -H 'name:Layotto' -H 'id:id_2' localhost:2045
Hi, Layotto_id_2
```

### Example description 

In this project, wasm modules with the same functions were developed with golang, rust and assemblyscript. Their process is as follows:

1. Receive HTTP requests through `proxy_on_request_headers`
2. Get the `name` field in the headers through `proxy_get_header_map_pairs`
3. Use `proxy_call_foreign_function` to call the APIs provided by Layotto
4. Return the result to the caller through `proxy_set_buffer_bytes`

golang source code path:

```
layotto/demo/wasm/code/golang/
```

rust source code path:

```
layotto/demo/wasm/code/rust/
```

assemblyscript source code path:

```
layotto/demo/wasm/code/assemblyscript/
```

### Note

This feature is still in the experimental stage, and the implementation of the WASM interactive API in the community is not uniform enough, so if you have any needs for this module, please post it in the issue area, we will build WASM together!