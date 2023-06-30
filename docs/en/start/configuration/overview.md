# Use Configuration API to call the configuration center

## What is Configuration API

When the application is started and running, it will read some "configuration information", such as: database connection parameters, startup parameters, RPC timeout, application port, etc. "Configuration" basically accompanies the entire life cycle of the application.

After the application evolves to the microservice architecture, it will be deployed on many machines, and the configuration will be scattered on each machine in the cluster, which is difficult to manage. So there is a "configuration center", which centrally manages the configuration, and also solves some new problems, such as: version management (in order to support rollback), authority management, etc.

There are many commonly used configuration centers, such as Spring Cloud Config, Apollo, Nacos, and cloud vendors often provide their own configuration management services, such as AWS Parameter Store, Google RuntimeConfig

Unfortunately, the APIs of these configuration centers are different. When developers want to deploy their apps across clouds, or want their apps to be portable (for example, easily moving from Alibaba Cloud to Tencent Cloud), they have to refactor their code.

The design goal of Layotto Configuration API is to define a unified configuration center API. Applications only need to care about the API, not which configuration center is used, so that the application can be transplanted at will, and the application is sufficiently "cloud native".

## What is the difference between Configuration API and State API?
Q: Why did we design the Configuration API?  What is the main difference with State API? I feel the two are almost the same

A: Configuration has some special capabilities, such as sidecar caching, such as app subscription to configuration change messages, such as configuration with some special schemas (tag, version, namespace, etc.)

This is like the difference between the configuration center and the database, both are storage, but the former is domain-specific and has special functions

## Quick start
- [Use Apollo as Configuration Center](en/start/configuration/start-apollo.md)
- [Use Etcd as Configuration Center](en/start/configuration/start.md)
- [Use Nacos as Configuration Center](en/start/configuration/start-nacos.md)