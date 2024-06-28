# Pluggable Component Design Document

## Background

The current Layotto component is actually working in Layotto, which requires users to develop in golang language in order to use the new component, and to do so in Layotto project, and then compile uniformly.
Being highly unfriendly to multilingual users, Layotto therefore needs to provide the ability to plugable components to allow users to communicate in any language using their components, Layotto via the grpc protocol and external components.

## Programmes

- Local cross-language component service discovery based on unix domain socket reduces communication costs.
- Based on proto achieve cross-language implementation of components.

## Stream Architecture

![](/img/pluggable/layotto_datatflow.png)

This is the current user's stream of data starting with sdk.The dotted section is the main data flow involved with pluggable components.

### Component Found

![](/img/pluggable/layotto.png)

As shown in the graph above, the user-defined component starts the socket service and places the socket file in the specified directory. When layotto starts, all socket files in this directory are read (skipped folder) and socket connected.

At present, layotto is aligned to the dapr and is not responsible for the lifecycle of the user component, which cannot be used if the user component is offline during the service period.
Later, depending on community use, it was decided whether layotto needed to support process management modules or to use a separate service to manage them.

As windows support for uds is not yet perfect, and layotto itself eliminates compatibility with windows, the uds discovery mode used for new features is not compatible with windows.

## Component Registration

The user registered component needs to implement plugable proto defined grpc services, as shown in the data stream architecture chart above. layotto will implement the forward interface based on the grpc interface, here
corresponds to the wrapp component in the data flow.prap component does not differ from build-in component for layotto runtime and has no special perception for users.

layotto gets what components have been implemented by the user service through grpc reflected library registered in the global component registration centre for user use.
