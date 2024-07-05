## TLDR

This proposal seeks to allow the open source Layotto to support both the Layotto API and the Dapr API.Similar to "Minio supports both Minio API, and AWS S3 API"

## Issues to solve

1. At present, we try to ensure that the various fields in the Layotto API are defined in the same way as Dapr, but the real concern of the users is whether they can be reused.While we are trying to ensure consistency in the proto field, as long as sdk cannot be reused, we do not solve users' problems and increase maintenance costs for ourselves.
   e.g.：
   ![image](https://user-images.githubusercontent.com/26001097/145837477-00fc5cd8-32eb-4ce9-bbfb-6e590172fce8.png)

So we want Layotto directly support the grpc API in Dapr (like a package, including package), which he can freely switch between with Dapr sdk without fear of being bound by the manufacturer.

2. On the other hand, some expansion is required.We found that the current Dapr API was not able to meet the needs completely, and some extension of the API was unavoidable.The extended API has been added to the current Layotto API. The proposal has been submitted to the Dapr community but is still waiting for it to be accepted by the community, such as config Is, such as the Lock API.

## Programmes

### Layotto API on Dapr API

![image](https://user-images.githubusercontent.com/26001097/145838604-e3a0caad-9473-4092-a2c6-0cc46c972790.png)

1. Layotto will start a grpc server with the frontline just adding an API plugin in the form of an API plugin.
2. On the other hand, the Layotto API is retained.Layotto received the Layotto API request and translated into Dapr API, then proceeded to Dapr API.
   Such benefits are：

- Reuse Code
- The Layotto API can be expanded according to production needs, such as support to the Lock API, configuration API, etc.; it can be extended to the Dapr community and then slowly discuss, even if the outcome of the final discussion differs from the original proposal, it only affects the resulting Dapr API, and does not affect users already using Layotto API.

### User Value

For users Issal：

- If users worry about the manufacturer's binding, they can only use Dapr API to migrate between Dapr and Layotto with the same set of Dapr sdk.
- If users believe our landed experience and are willing to use Layotto API, they can use Layotto API, at the cost of not migrating between two sidecar with the same sdk

### Q&A

#### How to add a field to Dapr API

##### Want to add a field (field)

For example, if you want to add a abc field to your layotto api, you can pass this field on to the dap API
dapr API implementation and then pass it over to the component, which parses this field

##### Not only fields but also logic, mechanisms (mechanism)

For example, layotto api adds a abc field if abc=true, then runtime takes a special logic

This will change the implementation of the Dapr API, plus an if else

#### Want to add new API

Add to the layotto API, the new API does not need to reuse Dapr API; wait for Dapr to receive the proposal before modifying it, layotto APIs.
