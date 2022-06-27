# State management component
**Configuration file structure**

The json configuration file has the following structure:

```json
"state": {
    "<Component A Name>": {
    "type":"<Component A Type>",
      "metadata": {
        "<KEY>": "<VALUE>",
        "<KEY>": "<VALUE>"
      }
    },
    "<Component B Name>": {
      "type":"<Component B Type>",
      "metadata": {
        "<KEY>": "<VALUE>",
        "<KEY>": "<VALUE>"
      }
    }
}
```

You can configure the key/value configuration items that the component cares about in the metadata. For example, [redis component configuration](https://github.com/mosn/layotto/blob/main/configs/config_redis.json) is as follows:

```json
"state": {
  "state_demo": {
    "type": "redis",
    "metadata": {
      "redisHost": "localhost:6380",
      "redisPassword": ""
    }
  }
}
```


**Common configuration item description**

The common configuration items of different State components are:

| Field | Required | Description |
| --- | --- | --- |
| keyPrefix | N | Key prefix strategy |


the `keyPrefix` field supports the following key prefix strategies:

* **`appid`** - The key passed in by the user will eventually be saved as `current appid||key`

* **`name`** - This setting uses the name of the component as a prefix. For example, the redis component will store the key passed in by the user as `redis||key`

* **`none`** - No prefix will be added. **This is the default policy.**

* Any other string that does not contain `||`. For example, if the keyPrefix is configured as "abc", the key passed in by the user will eventually be saved as `abc||key`


**Other configuration items**

In addition to the above general configuration items, each component has its own special configuration items. Please refer to the documentation for each component.