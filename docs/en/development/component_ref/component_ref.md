## layotto component reference

### Background

When a component starts, it may need to use another component's skill. For example, when the `sequencer` component `A` starts, it needs to read its settings from the `config` component `B`.

To make this happen, layotto offers the "component reference" feature. This feature lets component A use the features of component B.
### Related Designs

Currently, other components can only reference two types of components: ConfigStore and SecretStore. These are used to get configuration and secret keys. 

The "referenced" components must implement the interface :

```go
type SetComponent interface {
    SetConfigStore(cs configstores.Store) (err error)
    SetSecretStore(ss secretstores.SecretStore) (err error)
}
```

### How to configure

You can refer to the configuration file: `configs/config_ref_example.json`, configure the components to be used in the component configuration, and then inject them into the component when the component is initialized.

### How to use
Suppose we are developing a helloword component, it needs to read the secret key from the secret store (for example, to obtain the key to connect to the database) and read the configuration from the config store (for example, to read the IP address of the database to connect to the database) when it starts, then how should we develop the helloword component?

Take the `helloword` component as an example. First, the `helloword` component needs to implement the `SetConfigStore` and `SetSecretStore` interfaces. The interface implementation is the user's own logic, for example:

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

Then configure other components that need to be injected to the helloworld component in the configuration file, for example:

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

The `helloword` component can use the `config_demo` and `local.file` components during initialization
