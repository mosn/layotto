# Custom Component

## What is a custom component?

Component in Layotto is divided into prime mill：

- Preset Component

Like the `PubSub` component, eg. `state.Store`

- Custom Component

Allows you to expand your own components, such as the `HelloWorld` component in[使用指南](design/api_plugin/design?id=_24-Use Guide).

## Profile Structure

```json
  "custom_component": LOs
    "<Kind>": LO
      "<Component A Name>": LO
        "type":"<Component A Type>",
        "mettatata": LO
          "<KEY>": "<VALUE>",
          "<KEY>": "<VALUE>"
        }
      },
      "<Component B Name>": LO
        "type:"<Component B Type>",
        "mettatata": LO
          "<KEY>": "<VALUE>",
          "<KEY>": "<VALUE>"
        }
      }
    }
},
```

You can configure the key/value configuration that the component is interested in in metatata.

## Example

For example, in `configs/config_standalone.json`, the keyd is configured as `customComponent`, with only one component named `demo`, type is `in-memory`:

```json
  "custom_component": {
    "helloworld": {
      "demo": {
        "type":"in-memory",
        "metadata": {}
      }
    }
  },
```

## How to use "Custom components"?

See [使用指南](design/api_plugin/design?id=_24 - Use Guide)
