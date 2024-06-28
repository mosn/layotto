# mysql

## metadata fields
| Field        | Required  | Description                                                               |
|-----------|-----|---------------------------------------------------------------------------|
| mysqlUrl  | Y   | mysql's service address                                                   |
| username  | N   | specify username                                                          |
| password  | N   | specify password                                                          |
| dataBaseName| N   | mysql dataBaseName                                                        |
| connectionString| Y   | Database connection string such as "root:123456@tcp(127.0.0.1:3306)/test" |

## How to start mysql

If you want to run the mysql demo, you need to start a mysql server with Docker first.

```shell 
docker run --name mysql-test -d -p 3306:3306 mysql
```