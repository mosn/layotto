# Source Parsing 7 Layer Traffic Governance, Interface Limit

> Author Profile：
> was a fester of an open source community committed to embracing open sources and hoping to interact with each other’s open-source enthusiasts for progress and growth.
>
> Writing Time: 20 April 2022

## Overview

The purpose of this document is to analyze the implementation of the interface flow

## Prerequisite：

Document content refers to the following version of the code

[https://github.com/mosn/mosn](https://github.com/mosn/mosn)

Mosn d11b5a638a137045c2fb03d9d8ca36ecc0def11 (Division Develop)

## Source analysis

### Overall analysis

Reference to <br />[https://mosn.io/docs/concept/extensions/](https://mosn.io/docs/concept/extensions/)

Mosn Stream Filter Extension

![01.png](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*tSn4SpIkAa4AAAAAAAAAAAAAARQnAQ)

### Code in： [flowcontrol代码](https://github.com/mosn/mosn/tree/master/pkg/filter/stream/flowcontrol)

### stream_filter_factory.go analysis

This class is a factory class to create StreamFilter.

Some constant values are defined for default values

![02.png](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*PAWCTL6MS40AAAAAAAAAAAAAARQnAQ)

Defines the restricted stream config class to load yaml definition and parse production corresponding functions

![03.png](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*Ua32SokhILEAAAAAAAAAAAAAARQnAQ)

init() Inner initialization is the storage of name and corresponding constructor to the filter blocking plant map

![04.png](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*kb3qRqWnqxYAAAAAAAAAAAAAARQnAQ)

Highlight createRpcFlowControlFilterFactory Production rpc Current Factory

![05.png](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*u5rkS54zkgAAAAAAAAAAAAAAARQnAQ)

Before looking at streamfilter, we see how factory classes are producing restricted streamers

![06.png](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*cj0nT5O69OYAAAAAAAAAAAAAARQnAQ)

Limit the streaming to the restricted stream chain structure to take effect in sequential order.

CreateFilterChain method adds multiple filters to the link structure

![07.png](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*a8ClQ76odpEAAAAAAAAAAAAAARQnAQ)

We can see that this interface is achieved by a wide variety of plant types, including those that we are studying today.

![08.png](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*sBDbT44r2vgAAAAAAAAAAAAAARQnAQ)

### Stream_filter.go Analysis

![09.png](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*wsw3RKe1GH8AAAAAAAAAAAAAARQnAQ)

## Overall process：

Finally, we look back at the overall process progress:

1. Starting from the initialization function of stream_filter_factory.go, the program inserted createRpcFlowControlFilterFactory.

2. Mosn created a filter chain (code position[factory.go](https://github.com/mosn/mosn/tree/master/pkg/streamfilter/factory.go)) by circulating CreateFilterChain to include all filters in the chain structure, including our master restricted streaming today.

3. Create Limiter NewStreamFilter().

4. OnReceive() and eventually by sentinel (whether the threshold has been reached, whether to release traffic or stop traffic, StreamFilterStop or StreamFilterContinue).
