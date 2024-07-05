# mysql

## 配置项说明
| 字段        | 必填  | 说明         |
|-----------|-----|------------|
| mysqlUrl  | Y   | mysql的服务地址 |
| username  | N   | mysql用户名   |
 | password  | N   | mysql密码    |
 | dataBaseName| N   | mysql数据库名称 |
| connectionString| Y   | 数据库连接串 例如 "root:123456@tcp(127.0.0.1:3306)/test" |

## 怎么启动 mysql

如果想启动mysql的demo，需要先用Docker启动一个mysql命令：

```shell 
docker run --name mysql-test -d -p 3306:3306 mysql
```