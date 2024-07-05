# 文件管理组件

**配置文件结构**

json配置文件有如下结构：

```json
"file": {
    "file_demo": {
        "type": "aliyun.oss"
        "metadata":[
                {
                    "endpoint": "endpoint_address",
                    "accessKeyID": "accessKey",
                    "accessKeySecret": "secret",
                    "bucket": ["bucket1", "bucket2"]
                }
            ]
    }
}
```


**配置项说明**

配置项定义如下：

```golang
    type FileConfig struct {
	    Type     string
	    Metadata json.RawMessage
    }

    Files   map[string]file.FileConfig          `json:"file"`
```

上面的Files是一个map,key为component的名字，比如上述json的aliyun.oss，component的配置没有具体的格式限制，不同component可以根据需求自己定义，比如:

```json
"file": {
    "file_demo": {
      "type": "localFile",
      "group":{
        "name": "group1"
        "permisson":"rwx",
        "users":[
        "layotto","mosn"
        ]       
      }   
    }
}
```



