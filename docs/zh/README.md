<div align="center">
  <h1>Layotto (L8): To be the next layer of OSI layer 7</h1>
  <img src="https://gw.alipayobjects.com/zos/bmw-prod/65518bfc-8ba5-4234-a5c5-2bc065e3a5f0.svg" height="120px">

[![Layotto Env Pipeline ð](https://github.com/mosn/layotto/actions/workflows/quickstart-checker.yml/badge.svg)](https://github.com/mosn/layotto/actions/workflows/quickstart-checker.yml)
[![Layotto Dev Pipeline ð](https://github.com/mosn/layotto/actions/workflows/layotto-ci.yml/badge.svg)](https://github.com/mosn/layotto/actions/workflows/layotto-ci.yml)

[![GoDoc](https://godoc.org/mosn.io/layotto?status.svg)](https://godoc.org/mosn.io/layotto)
[![Go Report Card](https://goreportcard.com/badge/github.com/mosn/layotto)](https://goreportcard.com/report/mosn.io/layotto)
[![codecov](https://codecov.io/gh/mosn/layotto/branch/main/graph/badge.svg?token=10RxwSV6Sz)](https://codecov.io/gh/mosn/layotto)
[![Average time to resolve an issue](http://isitmaintained.com/badge/resolution/mosn/layotto.svg)](http://isitmaintained.com/project/mosn/layotto "Average time to resolve an issue")

</div>

Layotto(/leÉªËÉtÉÊ/) æ¯ä¸æ¬¾ä½¿ç¨ Golang å¼åçåºç¨è¿è¡æ¶, æ¨å¨å¸®å©å¼åäººåå¿«éæå»ºäºåçåºç¨ï¼å¸®å©åºç¨ååºç¡è®¾æ½è§£è¦ãå®ä¸ºåºç¨æä¾äºåç§åå¸å¼è½åï¼æ¯å¦ç¶æç®¡çï¼éç½®ç®¡çï¼äºä»¶åå¸è®¢éç­è½åï¼ä»¥ç®ååºç¨çå¼åã

Layotto ä»¥å¼æºç [MOSN](https://github.com/mosn/mosn) ä¸ºåºåº§ï¼å¨æä¾åå¸å¼è½åä»¥å¤ï¼æä¾äº Service Mesh å¯¹äºæµéçç®¡æ§è½åã

## è¯çèæ¯

Layotto å¸æå¯ä»¥æ [Multi-Runtime](https://www.infoq.com/articles/multi-runtime-microservice-architecture/) è· Service
Mesh ä¸¤èçè½åç»åèµ·æ¥ï¼æ è®ºä½ æ¯ä½¿ç¨ MOSN è¿æ¯ Envoy æèå¶ä»äº§åä½ä¸º Service Mesh çæ°æ®é¢ï¼é½å¯ä»¥å¨ä¸å¢å æ°ç sidecar çåæä¸ï¼ä½¿ç¨ Layotto ä¸ºè¿äºæ°æ®é¢è¿½å  Runtime çè½åã

ä¾å¦ï¼éè¿ä¸º MOSN æ·»å  Runtime è½åï¼ä¸ä¸ª Layotto è¿ç¨å¯ä»¥[æ¢ä½ä¸º istio çæ°æ®é¢](zh/start/istio/) åæä¾åç§ Runtime APIï¼ä¾å¦ Configuration API,Pub/Sub API ç­ï¼

æ­¤å¤ï¼éçæ¢ç´¢å®è·µï¼æä»¬åç° sidecar è½åçäºæè¿ä¸æ­¢äºæ­¤ã éè¿å¼å¥[WebAssembly](https://en.wikipedia.org/wiki/WebAssembly) ,æä»¬æ­£å¨å°è¯å° Layotto åæ FaaS (Function as a service) çè¿è¡æ¶å®¹å¨ ã

å¦ææ¨å¯¹è¯çèæ¯æå´è¶£ï¼å¯ä»¥çä¸[è¿ç¯æ¼è®²](https://mosn.io/layotto/#/zh/blog/mosn-subproject-layotto-opening-a-new-chapter-in-service-grid-application-runtime/index)
ã

## åè½

- æå¡éä¿¡
- æå¡æ²»çï¼ä¾å¦æµéçå«æåè§æµï¼æå¡éæµç­
- [ä½ä¸º istio çæ°æ®é¢](zh/start/istio/)
- éç½®ç®¡ç
- ç¶æç®¡ç
- äºä»¶åå¸è®¢é
- å¥åº·æ£æ¥ãæ¥è¯¢è¿è¡æ¶åæ°æ®
- åºäº WASM çå¤è¯­è¨ç¼ç¨

## å·¥ç¨æ¶æ

å¦ä¸å¾æ¶æå¾æç¤ºï¼Layotto ä»¥å¼æº MOSN ä½ä¸ºåºåº§ï¼å¨æä¾äºç½ç»å±ç®¡çè½åçåæ¶æä¾äºåå¸å¼è½åï¼ä¸å¡å¯ä»¥éè¿è½»éçº§ç SDK ç´æ¥ä¸ Layotto è¿è¡äº¤äºï¼èæ éå³æ³¨åç«¯çå·ä½çåºç¡è®¾æ½ã

Layotto æä¾äºå¤ç§è¯­è¨çæ¬ç SDKï¼SDK éè¿ gRPC ä¸ Layotto è¿è¡äº¤äºã

å¦ææ¨æ³æåºç¨é¨ç½²å°ä¸åçäºå¹³å°ï¼ä¾å¦å°é¿éäºä¸çåºç¨é¨ç½²å° AWSï¼ï¼æ¨åªéè¦å¨ Layotto æä¾ç [éç½®æä»¶](https://github.com/mosn/layotto/blob/main/configs/runtime_config.json)
éä¿®æ¹éç½®ãæå®èªå·±æ³ç¨çåºç¡è®¾æ½ç±»åï¼ä¸éè¦ä¿®æ¹åºç¨çä»£ç å°±è½è®©åºç¨æ¥æ"è·¨äºé¨ç½²"è½åï¼å¤§å¤§æé«äºç¨åºçå¯ç§»æ¤æ§ã

<img src="https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*oRkFR63JB7cAAAAAAAAAAAAAARQnAQ" />

## å¿«éå¼å§

### Get started with Layotto

æ¨å¯ä»¥å°è¯ä»¥ä¸ Quickstart demoï¼ä½éª Layotto çåè½ï¼æèä½éª[çº¿ä¸å®éªå®¤](zh/start/lab.md)

### API

| API            | status |                              quick start                              |                               desc                             |
| -------------- | :----: | :-------------------------------------------------------------------: | -------------------------------- |
| State          |   â    |        [demo](https://mosn.io/layotto/#/zh/start/state/start)         |     æä¾è¯»å KV æ¨¡åå­å¨çæ°æ®çè½å |
| Pub/Sub        |   â    |        [demo](https://mosn.io/layotto/#/zh/start/pubsub/start)        |     æä¾æ¶æ¯çåå¸/è®¢éè½å          |
| Service Invoke |   â    |       [demo](https://mosn.io/layotto/#/zh/start/rpc/helloworld)       |      éè¿ MOSN è¿è¡æå¡è°ç¨           |
| Config         |   â    | [demo](https://mosn.io/layotto/#/zh/start/configuration/start-apollo) |   æä¾éç½®å¢å æ¹æ¥åè®¢éçè½å     |
| Lock           |   â    |         [demo](https://mosn.io/layotto/#/zh/start/lock/start)         |    æä¾ lock/unlock åå¸å¼éçå®ç°  |
| Sequencer      |   â    |      [demo](https://mosn.io/layotto/#/zh/start/sequencer/start)       |  æä¾è·ååå¸å¼èªå¢ ID çè½å     |
| File           |   â    |         [demo](https://mosn.io/layotto/#/zh/start/file/start)         |   æä¾è®¿é®æä»¶çè½å               |
| Binding        |   â    |                                 TODO                                  |  æä¾éä¼ æ°æ®çè½å               |

### Service Mesh

| feature | status |                      quick start                       | desc                          |
| ------- | :----: | :----------------------------------------------------: | ----------------------------- |
| Istio   |   â    | [demo](https://mosn.io/layotto/#/zh/start/istio/) | è· Istio éæï¼ä½ä¸º Istio çæ°æ®é¢ |

### å¯æ©å±æ§

| feature  | status |                           quick start                            | desc                        |
| -------- | :----: | :--------------------------------------------------------------: | --------------------------- |
| API æä»¶ |   â    | [demo](https://mosn.io/layotto/#/zh/start/api_plugin/helloworld) | ä¸º Layotto æ·»å æ¨èªå·±ç API |

### å¯è§æµæ§


| feature    | status |                         quick start                         | desc                    |
| ---------- | :----: | :---------------------------------------------------------: | ----------------------- |
| Skywalking |   â    | [demo](https://mosn.io/layotto/#/zh/start/trace/skywalking) | Layotto æ¥å¥ Skywalking |


### Actuator

| feature        | status |                        quick start                        | desc                                  |
| -------------- | :----: | :-------------------------------------------------------: | ------------------------------------- |
| Health Check   |   â    | [demo](https://mosn.io/layotto/#/zh/start/actuator/start) | æ¥è¯¢ Layotto ä¾èµçåç§ç»ä»¶çå¥åº·ç¶æ |
| Metadata Query |   â    | [demo](https://mosn.io/layotto/#/zh/start/actuator/start) | æ¥è¯¢ Layotto æåºç¨å¯¹å¤æ´é²çåä¿¡æ¯   |

### æµéæ§å¶

| feature      | status |                              quick start                              | desc                                       |
| ------------ | :----: | :-------------------------------------------------------------------: | ------------------------------------------ |
| TCP Copy     |   â    |   [demo](https://mosn.io/layotto/#/zh/start/network_filter/tcpcopy)   | æ Layotto æ¶å°ç TCP æ°æ® dump å°æ¬å°æä»¶ |
| Flow Control |   â    | [demo](https://mosn.io/layotto/#/zh/start/stream_filter/flow_control) | éå¶è®¿é® Layotto å¯¹å¤æä¾ç API            |

### å¨ Sidecar ä¸­ç¨ WebAssembly (WASM) åä¸å¡é»è¾

| feature        | status |                      quick start                      | desc                                                             |
| -------------- | :----: | :---------------------------------------------------: | ---------------------------------------------------------------- |
| Go (TinyGo)    |   â   | [demo](https://mosn.io/layotto/#/zh/start/wasm/start) | æç¨ TinyGo å¼åçä»£ç ç¼è¯æ \*.wasm æä»¶è·å¨ Layotto ä¸         |
| Rust           |   â   | [demo](https://mosn.io/layotto/#/zh/start/wasm/start) | æç¨ Rust å¼åçä»£ç ç¼è¯æ \*.wasm æä»¶è·å¨ Layotto ä¸           |
| AssemblyScript |   â   | [demo](https://mosn.io/layotto/#/zh/start/wasm/start) | æç¨  AssemblyScript å¼åçä»£ç ç¼è¯æ \*.wasm æä»¶è·å¨ Layotto ä¸ |

### ä½ä¸º Serverless çè¿è¡æ¶ï¼éè¿ WebAssembly (WASM) å FaaS

| feature        | status |                      quick start                      | desc                                                                                      |
| -------------- | :----: | :---------------------------------------------------: | ----------------------------------------------------------------------------------------- |
| Go (TinyGo)    |   â   | [demo](https://mosn.io/layotto/#/zh/start/faas/start) | æç¨ TinyGo å¼åçä»£ç ç¼è¯æ \*.wasm æä»¶è·å¨ Layotto ä¸ï¼å¹¶ä¸ä½¿ç¨ k8s è¿è¡è°åº¦ã         |
| Rust           |   â   | [demo](https://mosn.io/layotto/#/zh/start/faas/start) | æç¨ Rust å¼åçä»£ç ç¼è¯æ \*.wasm æä»¶è·å¨ Layotto ä¸ï¼å¹¶ä¸ä½¿ç¨ k8s è¿è¡è°åº¦ã           |
| AssemblyScript |   â   | [demo](https://mosn.io/layotto/#/zh/start/faas/start) | æç¨ AssemblyScript å¼åçä»£ç ç¼è¯æ \*.wasm æä»¶è·å¨ Layotto ä¸ï¼å¹¶ä¸ä½¿ç¨ k8s è¿è¡è°åº¦ã |

## Landscapes

<p align="center">
<img src="https://landscape.cncf.io/images/left-logo.svg" width="150"/>&nbsp;&nbsp;<img src="https://landscape.cncf.io/images/right-logo.svg" width="200"/>
<br/><br/>
Layotto enriches the <a href="https://landscape.cncf.io/serverless">CNCF CLOUD NATIVE Landscape.</a>
</p>

## ç¤¾åº

| å¹³å°                                               | èç³»æ¹å¼                                                                                                                                                     |
| :------------------------------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| ð¬ [éé](https://www.dingtalk.com/zh) (ç¨æ·ç¾¤)     | ç¾¤å·: 31912621 æèæ«æä¸æ¹äºç»´ç  <br> <img src="https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*--KAT7yyxXoAAAAAAAAAAAAAARQnAQ" height="200px"> <br> |
| ð¬ [éé](https://www.dingtalk.com/zh) (ç¤¾åºä¼è®®ç¾¤) | ç¾¤å·ï¼41585216 <br> [Layotto å¨æ¯å¨äºæ 8 ç¹è¿è¡ç¤¾åºä¼è®®ï¼æ¬¢è¿ææäºº](zh/community/meeting.md)                                                               |

[comment]: <> (| ð¬ [å¾®ä¿¡]&#40;https://www.wechat.com/&#41; | æ«æä¸æ¹äºç»´ç æ·»å å¥½åï¼å¥¹ä¼éè¯·æ¨å å¥å¾®ä¿¡ç¾¤ <br> <img src="../img/wechat-group.jpg" height="200px">)

## å¦ä½è´¡ç®

[æ°ææ»ç¥ï¼ä»é¶å¼å§æä¸º Layotto è´¡ç®è](zh/development/start-from-zero.md)

[ä»åªä¸æï¼çç"æ°æä»»å¡"åè¡¨](https://github.com/mosn/layotto/issues/108#issuecomment-872779356)

ä½ä¸ºææ¯åå­¦ï¼ä½ æ¯å¦æè¿âæ³åä¸æä¸ªå¼æºé¡¹ç®çå¼åãä½æ¯ä¸ç¥éä»ä½ä¸æâçæè§ï¼
ä¸ºäºå¸®å©å¤§å®¶æ´å¥½çåä¸å¼æºé¡¹ç®ï¼ç¤¾åºä¼å®æåå¸éåæ°æçæ°æå¼åä»»å¡ï¼å¸®å©å¤§å®¶ learning by doing!

[ææ¡£è´¡ç®æå](zh/development/contributing-doc.md)

[ç»ä»¶å¼åæå](zh/development/developing-component.md)

[Layotto Github Workflow æå](zh/development/github-workflows.md)

[Layotto å½ä»¤è¡æå](zh/development/commands.md)

[Layotto è´¡ç®èæå](zh/development/CONTRIBUTING.md)

## è´¡ç®è

æè°¢ææçè´¡ç®èï¼

<a href="https://github.com/mosn/layotto/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=mosn/layotto" />
</a>

## è®¾è®¡ææ¡£

[Actuator è®¾è®¡ææ¡£](zh/design/actuator/actuator-design-doc.md)

[Pubsub API ä¸ Dapr Component çå¼å®¹æ§](zh/design/pubsub/pubsub-api-and-compability-with-dapr-component.md)

[Configuration API with Apollo](zh/design/configuration/configuration-api-with-apollo.md)

[RPC è®¾è®¡ææ¡£](zh/design/rpc/rpcè®¾è®¡ææ¡£.md)

[åå¸å¼é API è®¾è®¡ææ¡£](zh/design/lock/lock-api-design.md)

[FaaS è®¾è®¡ææ¡£](zh/design/faas/faas-poc-design.md)

## FAQ

### è· dapr æä»ä¹å·®å¼ï¼

dapr æ¯ä¸æ¬¾ä¼ç§ç Runtime äº§åï¼ä½å®æ¬èº«ç¼ºå¤±äº Service Mesh çè½åï¼èè¿é¨åè½åå¯¹äºå®éå¨çäº§ç¯å¢è½å°æ¯è³å³éè¦çï¼å æ­¤æä»¬å¸ææ Runtime è· Service Mesh ä¸¤ç§è½åç»åå¨ä¸èµ·ï¼æ»¡è¶³æ´å¤æççäº§è½å°éæ±ã
