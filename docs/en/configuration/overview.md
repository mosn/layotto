# Introduction to Layotto configuration file
Example: configs/config_apollo.json

Currently, Layotto uses a MOSN layer 4 filter to integrate with MOSN and run on MOSN, so the configuration file used by Layotto is actually a MOSN configuration file.

![img.png](../../img/configuration/layotto/img.png)

As shown in the example above, most of the configurations are MOSN configuration items, please refer to [MOSN configuration instructions](https://mosn.io/docs/configuration/);

Among them, the filter corresponding to `"type":"grpc"` is a layer 4 filter of MOSN, which is used to integrate Layotto and MOSN.

The configuration item in `grpc_config` is Layotto's component configuration, the structure is:

```json
"grpc_config": {
  "<API NAME>": {
    "<COMPONENT A NAME>": {
      "type": "<COMPONENT A Type>"
      "<KEY>": "<VALUE>",
      "metadata": {
        "<KEY>": "<VALUE>",
        "<KEY>": "<VALUE>"
      }
    }
  },
  "<API NAME>": {
    "<COMPONENT B NAME>": {
      "type": "<COMPONENT B Type>"
      "<KEY>": "<VALUE>",
      "metadata": {
        "<KEY>": "<VALUE>",
        "<KEY>": "<VALUE>"
      }
    }
  },
}
```

As for what to fill in each `<API NAME>`, what is each `<COMPONENT NAME>`, and which `"<KEY>": "<VALUE>"` configuration items can be configured with the components, you can refer to [Component specs](en/component_specs/overview) .
