[查看中文版本](zh/start/stream_filter/flow_control.md)

## Method Level Flow Control

### Introduction

There is a config of flow control in [runtime_config.json](https://github.com/mosn/layotto/blob/main/configs/runtime_config.json):

```json
[
    {
        "type": "flowControlFilter",
        "config": {
            "global_switch": true,
            "limit_key_type": "PATH",
            "rules": [
                {
                    "resource": "/spec.proto.runtime.v1.Runtime/SayHello",
                    "grade": 1,
                    "threshold": 5
                }
            ]
        }
    }
]
```

this can help `/spec.proto.runtime.v1.Runtime/SayHello` method has a flow control feature, which means we can only access this method below 5 times in 1 second.

this code of the client is here [client.go](https://github.com/mosn/layotto/blob/main/demo/flowcontrol/client.go)，the logic is very simple, send 10 times request to the server，and the result is below:

![img.png](../../../img/flow_control.png)

the previous 5 times request access is succeed while the last 5 times request is under control.

### Configuration Description

config:

| field_name | field_type | desc |
|  ----  | ----  | ---- |
| global_switch  | boolean | switch of this feature, true is open, false is close |
| limit_key_type  | String | Unique identifier，Currently a fixed value `PATH` |
| rules  | Array | rules of flow control |

rules of flow control:

| field_name | field_type | desc |
|  ----  | ----  | ---- |
| resource  | String | limit resource identification，here is the method URL |
| grade  | int | Threshold update cycle，unit:second |
| threshold  | int | be limited if exceed this threshold |