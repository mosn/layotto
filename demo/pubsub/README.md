You can run server/client demo with different component names.   
It is worth noting that both server and client demo should set the same store name by param `-s`.  
For example:   
```shell
cd ${project_path}/demo/pubsub/server/
# 1. start subscriber
go build -o subscriber
/.subscriber -s pub_subs_demo

# 2. start layotto
cd ${project_path}/cmd/layotto
go build -o layotto
./layotto start -c ../../configs/config_standalone.json

# 3. start publisher
 cd ${project_path}/demo/pubsub/client/
go build -o publisher
 ./publisher -s pub_subs_demo

```