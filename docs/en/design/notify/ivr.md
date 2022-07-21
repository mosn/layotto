# IVR API design
<!-- Please only use this template for submitting enhancement requests -->

## What would you like to be added
IVR API, or PhoneCall API

Developers can invoke this API to send voice messages to specific people.

## Why is this needed
In the monitoring scenarios, monitor systems need to send alarm messages to people on-call.
The messages might be in different forms, including IM,SMS, Email and phone calls, depending on the level of urgency.

## Detailed Design

We need to consider the following factors:
- Portability
  
  For example, a monitor system might be deployed on alibaba cloud(using [VMS](https://www.aliyun.com/product/vms) to send voice message) or AWS (using [AWS Pinpoint](https://aws.amazon.com/cn/pinpoint/)  to send voice message). So portability is important here.

```proto

service IvrService {

  //Send voice using the specific template
  rpc SendVoiceWithTemplate(SendVoiceWithTemplateRequest) returns (SendVoiceWithTemplateResponse) {}

}

message SendVoiceWithTemplateRequest{
  // Required
  Template template = 1;
  // Required
  string to_mobile = 2;
  // This field is required by some cloud providers.
  string from_mobile = 3;
}

message Template{
  // Required
  string template_id = 1;
  // Required
  map<string, string>  template_params = 2;
}

message SendVoiceWithTemplateResponse{
  string message_id = 1;
}


```
