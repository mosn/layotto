# Distributed lock component
**Configuration file structure**

The json configuration file has the following structure:

```json
"lock": {
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
"lock": {
  "lock_demo": {
    "type": "redis",
    "metadata": {
      "redisHost": "localhost:6380",
      "redisPassword": ""
    }
  }
}
```

**Common configuration item description**

| Field | Required | Description |
| --- | --- | --- |
| keyPrefix | N | Key prefix strategy |


the `keyPrefix` field supports the following key prefix strategies:

* **`appid`** - This is the default policy. The resource_id passed in by the user will eventually be saved as `lock|||current appid||resource_id`

* **`name`** - This setting uses the name of the component as a prefix. For example, the redis component will store the resource_id passed in by the user as `lock|||redis||resource_id`

* **`none`** - The resource_id passed in by the user will eventually be saved as `lock|||resource_id`.

* Any other string that does not contain `||`. For example, if the keyPrefix is configured as "abc", the resource_id passed in by the user will eventually be saved as `lock|||abc||resource_id`


**Other configuration items**

In addition to the above general configuration items, each distributed lock component has its own special configuration items. Please refer to the documentation for each component.