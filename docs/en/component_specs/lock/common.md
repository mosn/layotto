# Distributed lock component

**Common configuration item description**

| Field | Required | Description |
| --- | --- | --- |
| keyPrefix | N | Key prefix strategy |


the `keyPrefix` field supports the following key prefix strategies:

* **`appid`** - This is the default policy. The resource_id passed in by the user will eventually be saved as `current appid||resource_id`

* **`name`** - This setting uses the name of the component as a prefix. For example, the redis component will store the resource_id passed in by the user as `redis||resource_id`

* **`none`** - No prefix will be added.

* Any other string that does not contain `||`. For example, if the keyPrefix is configured as "abc", the resource_id passed in by the user will eventually be saved as `abc||resource_id`


**Other configuration items**

In addition to the above general configuration items, each distributed lock component has its own special configuration items. Please refer to the documentation for each component.