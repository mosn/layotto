# MongoDB

## 配置项说明     

实例：configs/config_lock_mongo.json

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| host | Y | MongoDB的服务地址，例如localhost:27017 |
| username | N | MongoDB用户名 |
| password | N | MongoDB密码 |
| database | N | MongoDB数据库 |
| params | N | 自定义参数 |


## 怎么启动 MongoDB

如果想启动MongoDB的demo，需要先用Docker启动一个MongoDB 命令：

```shell 
docker run --name mongoDB -d -p 27017:27017 mongo
```