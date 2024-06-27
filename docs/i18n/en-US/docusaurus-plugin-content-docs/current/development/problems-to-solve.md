# Multi-Runtime 2022：pending issue

## API Standard-Building

Building API standards based on the production needs of landed users continues to be presented to the Dapr community.Like：

- Distributed LockAPI
- Configure API
- Delay Message API

## Ecological development

How to smooth migration to Multi-Runtime for users who have fallen Service Mesh?One thing currently being done is Layotto on Envoy's support;

Can the Runtime API be better integrated into K8S ecology?What is being done is Layotto integration into k8s ecology;

## Early productive users of services

Open source functions that are universal and that address production problems.Watch early production users, currently facing the following problems：

### Expansion

Let the entire project be extended, for example, a company wants to expand some of its features with layotto, either by starting up a project, by importing the layoto's source layoutto, or by expanding the layoutto binary files by connecting dynamically to the library.Neither of these options, dapr nor layotto cannot, want to extend to fork to change code

### Stability risk

After an important open source of Layotto, the panic is at great risk because it relies on all Dapr components, which use a wide variety of libraries and may be panic, and may rely on conflict.Can panic risks be reduced by customizing, isolating designs?

There are currently too few open source project test inputs relative to the testing process in the firm, how to build an open source test system;

### Detectability

> There were problems with my service mesh. Then there was the service messh. I could only look for someone else to check on
> — a test goes on.

Service Mesh in the productive environment makes troubleshooting, and multi-Runtime is more functional and more difficult.
Multi-Runtime observability needs to be built to avoid making problems more difficult for productive users.

## New research and development models

sidecar supports serverless land;
