# Apollo

## Configuration item description
Example: configs/config_apollo.json
![img.png](../../../img/configuration/apollo/img.png)

| Field | Required | Description |
| --- | --- | --- |
| address | Y | apollo server address, array type |
| metadata.app_id | Y | Corresponds to 'application' in the apollo data model |
| metadata.cluster | Y | Corresponding to 'cluster' in the apollo data model |
| metadata.namespace_name | Y | Corresponding to the 'namespace' in the apollo data model. You can join multiple namespaces with commas, such as "dubbo,product.joe,application" |
| metadata.is_backup_config | N | Whether to back up the configuration to a local file, corresponding to [agollo sdk](https://github.com/apolloconfig/agollo/wiki/%E4%BD%BF%E7%94%A8%E6%8C%87%E5%8D%97) is_backup_config configuration item. The default value is true |
| metadata.secret | N | secret for access to apollo |
| metadata.open_api_address | Y | The address of apollo open-api (open-api is used to modify configuration items, which is different from the apollo server address used for query) |
| metadata.open_api_token | Y | Token needed to access apollo open-api |
| metadata.open_api_user | Y | User accessing apollo open-api |

## How to start Apollo
There is no need to deploy the apollo server yourself to use the demo in the project. The demo will use the demo environment provided by apollo http://81.68.181.139

If you want to deploy apollo yourself, you can refer to [apollo official document](https://www.apolloconfig.com/#/zh/deployment/quick-start)

After deployment, you need to modify Layotto's [config file](https://github.com/mosn/layotto/blob/main/configs/config_apollo.json) to change the apollo server address and other information to your own.