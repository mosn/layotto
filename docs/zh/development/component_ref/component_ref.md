## layotto组件引用

### 背景

组件初始化的时候，需要引用其他组件的能力，比如sequencer组件初始化的时候需要从config组件读取相关的配置，以此实现组件之间的引用。

###  相关设计

目前只支持最需要被应用的两个组件类型:ConfigStore和SecretStore,即配置组件和秘钥组件。  
Ref接口设计，组件需实现此接口才能实现注入。  

```go
type SetComponent interface {
	SetConfigStore(cs configstores.Store) (err error)
	SetSecretStore(ss secretstores.SecretStore) (err error)
}
```

旧组件不去适配此接口，用户有注入需求可以通过实现Ref接口来实现注入。

### 如何配置

可参考配置文件:`configs/config_ref_example.json`, 在组件配置中配置需要使用的组件，便可以在组件初始化时注入给组件。

### 如何使用
假如我们要想开发一个helloword组件，它需要在启动的时候，从 secret store 中读取秘钥(比如用来获取连接数据库的秘钥)、从 config store 读取配置(例如读取数据库的 ip 地址，以便连接数据库), 那么 helloword 组件应该如何开发呢？
以`helloword`组件为例,首先`helloword`组件需要实现`SetConfigStore`和`SetSecretStore`接口,接口实现里是用户自己的逻辑，例如：

```go
func (hw *HelloWorld) SetConfigStore(cs configstores.Store) (err error) {
//save for use
hw.config=cs
return nil
}
func (hw *HelloWorld) SetSecretStore(ss secretstores.SecretStore) (err error) {
//save for use
hw.secretStore = ss
return nil
}
//fetch secret/config when component init
func (hw *HelloWorld) Init(config *hello.HelloConfig) error {
  hw.secretStore.GetSecret(secretstores.GetSecretRequest{
     Name:     "dbPassword",
  })
  hw.config.Get(context.Background(),&configstores.GetRequest{
     Keys:     []string{"dbAddress"},
  })
  return nil
}
```

然后再配置文件里给helloworld组件配置需要注入的其他组件，例如：

```json
        {
  "helloworld": {
    "type": "helloworld",
    "hello": "greeting",
    "secret_ref": [
      {
        "store_name": "local.file",
        "key": "db-user-pass:password",
        "sub_key": "db-user-pass:password",
        "inject_as": "redisPassword"
      }
    ],
    "component_ref": {
      "config_store": "config_demo",
      "secret_store": "local.file"
    }
  }
}
```

`helloword`组件在初始化的时候便可以使用`config_demo`和`local.file`组件
