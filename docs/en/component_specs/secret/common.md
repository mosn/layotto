# Secret Store component
This component can access secrets from local files, environment variables, k8s, etc.,  Layotto use dapr's secret API, learn more: https://docs.dapr.io/operations/components/setup-secret-store/
**Configuration file structure**

The json configuration file has the following structure:

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

Configuration examples of local file keys, local environment variables, and k8s keys:

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