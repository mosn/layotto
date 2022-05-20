# Secret Store component
This component can access secrets from local files, environment variables, k8s, etc.,  Layotto use dapr's secret API, learn more: https://docs.dapr.io/operations/components/setup-secret-store/
**Configuration file structure**

The json configuration file has the following structure:
```json
"secret_store": {
  "<STORE NAME>": {
    "metadata": {
      "<KEY>": "<VALUE>",
      "<KEY>": "<VALUE>"
    }
  }
}
```
Configuration examples of local file keys, local environment variables, and k8s keys:
```
       "secret_store": {
                        "local.file": {
                          "metadata": {
                            "secretsFile": "../../configs/config_secret_local_file.json"
                          }
                        },
                        "local.env": {
                          "metadata": {
                          }
                        },
                        "kubernetes": {
                          "metadata": {
                          }
                        }
                      }
```