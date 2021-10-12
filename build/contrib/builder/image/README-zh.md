# 通过docker安装layotto

## 拉取镜像
```shell
docker pull mosn/layotto
```

## 命令行启动
```shell
docker run -it -p 34904:34904  mosn/layotto
```

## 配置
layotto的配置目录为``/home/admin/layotto/configs``，默认加载``config.json``文件，可以通过以下两种方式修改配置：
### 方法一：挂载配置文件
可以自定义配置文件，然后将该文件挂载到配置目录中，通过``-c``命令指定该配置文件：
```json
{
  "servers":[
    {
      "default_log_path":"stdout",
      "default_log_level": "INFO",
      "listeners":[
        {
          "name":"grpc",
          "address": "0.0.0.0:34904",
          "bind_port": true,
          "filter_chains": [{
            "filters": [
              {
                "type": "grpc",
                "config": {
                  "server_name":"runtime",
                  "grpc_config": {
                    "hellos": {
                      "helloworld": {
                        "hello": "greeting"
                      }
                    }
                  }
                }
              }
            ]
          }],
          "stream_filters": [
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
        }
      ]
    }
  ]
}
```
```shell
touch custom_config.json
docker run -it -p 34904:34904 -v full_path_to/custom_config.json:/home/admin/layotto/configs/custom_config.json mosn/layotto -c configs/custom_config.json
```

### 方法二：自定义镜像
如果在某些情况下想通过layotto基础镜像来自定义自己的镜像，也自定义可以编写dockerfile：
```dockerfile
FROM mosn/layotto
COPY configs/custom_config.json /home/admin/layotto/configs/custom_config.json
CMD ["-c", "configs/custom_config.json"]
```
然后以该dockerfile构建镜像运行：
```shell
docker build -f Dockerfile -t layotto:demo . 
docker run -it -p 34904:34904  layotto:demo
```