# MongoDB

## metadata fields

Example：configs/config_lock_mongo.json

|  Field | Required | Description |
| --- | --- | --- |
| host | Y | MongoDB server address, such as localhost:27017 |
| username | N |specify username username |
| password | N | specify password |
| params | N | custom params |


## How to start MongoDB

If you want to run the mongoDB demo, you need to start a mongoDB server with Docker first.

```shell 
docker run --name mongoDB -d -p 27017:27017 mongo
```