# Dynamically configure the sending, component heated reload

## Issues addressed

- Now the producer user has a customized set of initialization configurations：with some configurations in the app, after app starts, sidecar, allowing sidecar to be initialized based on these configurations.Programmes are not common enough and are more common to do

![image](https://user-images.githubusercontent.com/26001097/168947177-6a26397e-4648-47f0-a8df-e898285cd8f9.png)

- Support "Dynamically Down Configuration".
  - One idea is to couple the configuration file and mirror to the container via disk.For example, Dapr config item is released into Configuration CRD, CRD changes will require the carrier to reboot the cluster via k8s scroll.
  - Another line of thought is [inject config_store components into other components] (https://github.com/mosn/layotto/issues/500#issuomment-1119390497) but has some disadvantages:
    - Users who want to use the "Dynamic Configuration" function have no means of extending hands, no community ready-made components, and have to develop their components.
      The runtime level is best used for some common features, "empowerment" all components, community maintenance of ready-made components, support for dynamic configuration, and user-friendliness and open boxes.
  - Another line of thought is to split the configuration between the first expert：bootstrap configuration (static configuration), which is not placed in mirrors, and the dynamic configuration which supports the release of the configuration and heated reload based on the configuration.

## Product design

### User story

1. The user changed the disaster switching configuration to Redis on the apollo page, so that the Redis component received the new configuration and switches the traffic to the Disaster Preparedness Cluster
2. There are already producing users who can migrate the initialization process to a new model that will be compatible downward.

### Programming UI

For example, the start configuration for state.redis now has the following (screenshot taken from [dapr documents](https://docs.dapr.io/reference/components-reference/supported-state-stores/setup-redis/)
![image](https://user-images.githubusercontent.com/26001097/168946975-9804d792-8851-463f-80ee-26231468f0aa.png)

The status quo is that these configurations kv are initialized when the：redis component startup; all configurations are static configurations, only once, and no subsequent configuration changes are listened.

But we can change：

- These kv can be dynamically employed
- layotto listen to these kv changes, except for changes that **Reinitialize component with the latest configuration**
- Dynamic update interface can be implemented if the component feels too small to reinitialize

Advantage and disadvantages analysis：

- pros
  - runtime layer can perform some general functionality, "empowerment" all components; easy access for users, community maintenance of ready-made components, and dynamic configuration support, user open boxes
- Cons
  - Realization is complex.How do you ensure that traffic does not damage during reinitialization?
  - I am not aware that this will not meet users' production needs, worries about early design, over-design

## High-level design

![image](https://user-images.githubusercontent.com/26001097/168949648-3f440a84-45d3-45c1-89ef-79cb25d49713.png)

### Exposes UpdateConfiguration API after startup

Sidear starts up or uses json files to make a new API for configuration thermal changes once it is passed:

```protobuf
rpc UpdateConfiguration( RuntimeConfig) returns (UpdateResponse)
```

### Agent is responsible for interacting and calling the UpdateConfiguration API

That is, Sidecar is just opening an interface and waiting for others to configure it.And things that interact with the face and subscribe to configuration changes can be done with an agent 2 on the map, which subscribes to the apollo configuration change, which changes the interface with Sidecar to make Sidecar hot updates.

For existing production users, you can listen to the app feeding configuration, dump configuration, load configuration on reboot and push configuration to Sidecar.

For example, you can write a File Agent issue, listen to file changes, read the new configuration and notify Sidecar to reload.

Agent does not have to be a separate process but also start a separate process in the main one.

### Component Hot Reload

When Sidecar is brought to an UpdateConfiguration API, it will:

1. No "Increment Update" interface has been implemented by the judging component:

```go
UpdateConfig(ctx context.Context, metadata map[string]string) (err error, needReload bool)
```

2. runtime tries to update components if they have an interface, runtime
3. runtime **Reinitialize component based on full configuration** if incremental update fails, or if the interface is not implemented
4. After the new component has been reinitialized (check through readability), take over traffic from the original component

## Detailed design

### GRPC API design

```protobuf
Service Lifecycle L/

  rpc ApplyConfiguration(DynamicConfiguration) returnns (ApplyConfigurationResponse) {}

}

message Dynamic Configuration{

  Component_config = 1;

} } }

message ApplyConfigurationResponse
}
```

#### ComponentConfig Field Design

##### Design of a common updating interface

```protobuf
message Composiconfig Fact

  // For example, `lock`, `state`
  string type = 1;

  // The component name. For example, `state_demo`
  string name = 2;

  map<string, string> metadata = 3;
}
```

~With google/protobuf/struct.proto describe dynamic json see https://stackoverflow.com/questions/5296644/is-google-proto-the-best-way-to-send-dynamic-json-over-grpc~!

Upload configuration with `map<string, string>`.

- Advantages
  Do not change the sdk of each language for each new API or configuration structure to allow users to pass through and sidecar side to deserialize

- The
  field format does not show definitions, are not clear, and does not have enough personnel

##### b. Structured definitions of each class configuration

```protobuf
// Component configuration
message ComponentConfigure
  // For example, `lock`, `state`
  string kind = 1;
  // The component name. For example, `state_demo`
  string name = 2;

  google. Rotobuf. Cart metadata = 3;

  one of common_config LO
    LockCommonConfigurationlock_config = 4;

    StateCommonConfiguration state_config = 5;

    // .
  }
}
```

Advantage and disadvantages above

##### Conclusion

Select A, reduce the cost of SDK maintenance

#### Q: A separate API plugin or an existing API plugin

Write a separate API plugin

#### Q: Was to configure vs active pull configuration vs push back

Wait Push Configuration

#### Q: API accepts full or incremental configuration

Additions, sequencing issues are guaranteed by stream.

```protobuf
service Lifecycle {

  rpc UpdateComponentConfiguration(stream ComponentConfig) returns (UpdateResponse){}

}
```

b. Full amount

Conclusion: b, simplerYou can add an additional interface to make incremental changes via stream, if needed.

### Component API design

```go
Type Dynamic interface of Jean-Marie
    ApplyConfig(ctx context.Context, metadata map[string]string) (err error, needReload bool)
}
```

## Future work

### pubsub subscriber service

Requires some more structured configuration data to be released

### Component Hot Reload

//TODO

- How to ensure the loss of data during reinitialization
- Config Priority：has some configurations that are customized for a single app and some are common configurations for all app utilities, with priority for both parties
- Configure transaction reading and writing to avoid dirty reading
