## layotto component reference

### Background

When a component is initialized, it needs to refer to the ability of other components. For example, when the sequencer component is initialized, it needs to read the relevant configuration from the config component, so as to realize the reference between components.

### Related Designs

Currently, only two component types that need to be applied most are supported: ConfigStore and SecretStore, that is, configuration components and secret key components.
Ref interface design, components need to implement this interface to achieve injection.
````go
type SetComponent interface {
SetConfigStore(cs configstores.Store) (err error)
SetSecretStore(ss secretstores.SecretStore) (err error)
}
````
The old components do not adapt to this interface. Users who have injection requirements can implement the injection by implementing the Ref interface.

### How to configure

You can refer to the configuration file: `configs/config_ref_example.json`, configure the components to be used in the component configuration, and then inject them into the component when the component is initialized.

### How to use
Take the `helloword` component as an example. First, the `helloword` component needs to implement the `SetConfigStore` and `SetSecretStore` interfaces. The interface implementation is the user's own logic, for example:
````go
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
````
Then configure other components that need to be injected to the helloworld component in the configuration file, for example:
````json
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
````
The `helloword` component can use the `config_demo` and `local.file` components during initialization