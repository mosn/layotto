# 秘钥组件
该组件可以从本地文件、环境变量、k8s等获取秘钥，复用了dapr的secret API，了解更多：https://docs.dapr.io/operations/components/setup-secret-store/

**配置文件结构**

json配置文件有如下结构：

```json
"secret_store": {
  "<Component A Name>": {
    "type":"<Component A Type>",
      "metadata": {
        "<KEY>": "<VALUE>",
        "<KEY>": "<VALUE>"
      }
    },
    "<Component B Name>": {
      "type":"<Component B Type>",
      "metadata": {
        "<KEY>": "<VALUE>",
        "<KEY>": "<VALUE>"
      }
    }
}
```

本地文件秘钥、本地环境变量、k8s秘钥的配置例子：

```json
       "secret_store": {
                        "secret_demo": {
                          "type": "local.file",
                          "metadata": {
                            "secretsFile": "../../configs/config_secret_local_file.json"
                          }
                        },
                        "secret_demo1": {
                          "type": "local.env",
                          "metadata": {
                          }
                        },
                        "secret_demo2": {
                          "type": "kubernetes",
                          "metadata": {
                          }
                        }
                      }
```