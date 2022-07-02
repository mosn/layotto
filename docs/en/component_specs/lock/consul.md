# Consul

## metadata fields

Example：configs/config_consul.json

| Field | Required | Description |
| --- | --- | --- |
| address | Y | consul server address, such as localhost:8500 |
| scheme | Y | client connection scheme,HTTP/HTTPS |
| username | N | specify username |
| password | N | specify password |

## How to start Consul

If you want to run the Consul demo, you need to start a Consul server with Docker first.

command:

```shell
docker run --name consul -d -p 8500:8500 consul
```