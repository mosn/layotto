# PhoneCall API design

## What would you like to be added
IVR API, or PhoneCall API

Developers can invoke this API to send voice messages to specific people.

## Why is this needed
In the monitoring scenarios, monitor systems need to send alarm messages to people on-call.
The messages might be in different forms, including IM,SMS, Email and phone calls, depending on the level of urgency.

## Product research
| IVR product |Docs|
|---|---|
|Aliyun VMS| https://www.aliyun.com/product/vms |
|AWS Pinpoint | https://aws.amazon.com/cn/pinpoint/ |


## Detailed Design

We need to consider the following factors:
- Portability
  For example, a monitor system might be deployed on alibaba cloud(using [VMS](https://www.aliyun.com/product/vms) to send voice message) or AWS (using [AWS Pinpoint](https://aws.amazon.com/cn/pinpoint/)  to send voice message). So portability is important here.

```proto
// PhoneCallService is one of Notify APIs. It's used to send voice messages
service PhoneCallService {

  // Send voice using the specific template
  rpc SendVoiceWithTemplate(SendVoiceWithTemplateRequest) returns (SendVoiceWithTemplateResponse) {}

}

// The request of SendVoiceWithTemplate method
message SendVoiceWithTemplateRequest{

  // If your system uses multiple IVR services at the same time,
  // you can specify which service to use with this field.
  string service_name = 1;

  // Required
  VoiceTemplate template = 2;

  // Required
  repeated string to_mobile = 3;

  // This field is required by some cloud providers.
  string from_mobile = 4;

}

// VoiceTemplate
message VoiceTemplate{

  // Required
  string template_id = 1;

  // Required
  map<string, string>  template_params = 2;

}

// The response of `SendVoiceWithTemplate` method
message SendVoiceWithTemplateResponse{

  // Id of this request.
  string request_id = 1;

}

```
