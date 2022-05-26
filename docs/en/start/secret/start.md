# use Secret API to obtain secret
## What is Secret API
The secret API is used to obtain secret from file, env, k8s, etc

Get all API and secret support
## Quick start

This example shows how to obtain the secret in file, env and k8s through the Layotto secret API



### Step 1:  Run Layotto

After downloading the project code to the local, switch the code directory and compile:

```shell
cd ${project_path}/cmd/layotto
```

build:
```shell @if.not.exist layotto
go build -o layotto
```

Once finished, the layotto file will be generated in the directory, run it:

```shell @background
./layotto start -c ../../configs/config_secret_file.json
```

### Step 2: Run the client program and call Layotto to generate a unique id

```shell
 cd ${project_path}/demo/secret/common/
```

```shell @if.not.exist client
 go build -o client
```

```shell
 ./client -s "secret_demo"
```

If the following information is printed, the demo is successful:

```bash
data:{key:"db-user-pass:password" value:"S!S*d$zDsb="}
data:{key:"db-user-pass:password" value:{secrets:{key:"db-user-pass:password" value:"S!S*d$zDsb="}}} data:{key:"db-user-pass:username" value:{secrets:{key:"db-user-pass:username" value:"devuser"}}}
```


## Want to learn more about Secret API?
Layotto reuse Dapr Secret API，learn more：https://docs.dapr.io/operations/components/setup-secret-store/
