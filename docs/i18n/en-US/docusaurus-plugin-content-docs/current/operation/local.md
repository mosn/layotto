# How to Debugging Locally

## 1. How the app developers can develop and debug the app

There are generally the following types of programmes:

### Local commencement of sidecar

Start Layotto+other storage systems with a docker to sidecar, or docker-compose. Start Layotto+Other storage (e.g. Reddis)

### The company provides remote Layotto sidecar

For example, test environments in remote areas, pull up a Pod, running Layotto sidecar inside.

- If you have direct access to the remote test environment with ip, pod：
  - You can change Layp to Pod Pip, local ip link pod
- 如果不能以ip直接访问远端测试环境pod：
  - This pod service type can be set to `NodePort` or `LoadBalancer`, local direct service
  - You can register the pod to gateway, directly to gateway

When debugging locally, the "Local app process is connected to remote Layotto sidecar" is implemented in the above way.

To go further, the team in the company responsible for research and development can automate the above actions by providing the "One click to apply for remote sidecar" feature.

### Cloud R&D environment

If a company has a cloud research environment similar to the github codespace, it can send sidecar in a research environment

## How Layotto developers locally, debugging Layotto

The local compilation will run Layotto.

For example, when running Layotto with Goland IDE, the configuration is as follows:

![](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*CHFYQK6kMEgAAAAAAAAAAAAAARQnAQ)
