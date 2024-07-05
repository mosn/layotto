

const sidebars = {

  mySidebar: [
    {
      type: 'doc',
      label: '首页',
      id: 'README'
    },
    {
      type: 'category',
      label: '快速开始',
      link: {
        type: 'doc',
        id: 'start/README',
      },
      items: [
        {
          type: "doc",

          id: "start/state/start",
        },
        {
          type: "category",
          label: "使用Configuration API",
          items:[
            "start/configuration/start-apollo",
            {
              type: "doc",
              id: "start/configuration/start",
            },
            {
              type: "doc",
              id: "start/configuration/start-nacos",
            },
          ]
        },
        {
          type: "doc",
          id: "start/pubsub/start",
        },
        {
          type: "doc",
          id: "start/delay_queue/start",
        },
        {
          type: "doc",
          id: "start/lock/start",
        },
        {
          type: "doc",
          id: "start/sequencer/start",
        },
        {
          type: "doc",
          id: "start/secret/start",
        },
        {
          type: "doc",
          id: "start/rpc/helloworld",
        },
        {
          type: "doc",
          id: "start/file/minio",
        },
        {
          type: "doc",
          id: "start/oss/oss",
        },
        {
          type: "doc",
          id: "start/uds/start",
        },
        {
          type: "doc",
          id: "start/sms/start",
        },
        {
          type: "doc",
          id: "start/cryption/start",
        },
        {
          type: "doc",
          id: "start/phone/start",
        },
        {
          type: "doc",
          id: "start/email/start",
        },
        { type: "doc", id: "start/lifecycle/start", },
        {
          type: 'category',
          label: 'API插件',
          items: [
            {
              type: 'doc',
              id: 'start/api_plugin/helloworld',
            },
            {
              type: 'doc',
              id: 'start/api_plugin/generate',
            }
            ]},
            {
              type: 'category',
              label: '作为 Istio 的数据面',
              items: [
                {
                  type: 'doc',
                  id: 'start/istio/README',
                },
                {
                  type: 'doc',
                  id: 'start/istio/start',
                },
              ],
            },
            {
              type: 'category',
              label: '在四层网络进行流量干预',
              items: [
                {
                  type: 'doc',
                  id: 'start/network_filter/tcpcopy',
                },
              ],
            },
            {
              type: 'category',
              label: '在七层网络进行流量干预',
              items: [
                {
                  type: 'doc',
                  id: 'start/stream_filter/flow_control',
                },
              ],
            },
            {
              type: 'doc',
              id: 'start/actuator/start',
            },
            {
              type: 'category',
              label: '可观测性',
              items: [
                {
                  type: 'doc',
                  id: 'start/trace/trace',
                },
                {
                  type: 'doc',
                  id: 'start/trace/skywalking',
                },
                {
                  type: 'doc',
                  id: 'start/trace/zipkin',
                },
                {
                  type: 'doc',
                  id: 'start/trace/jaeger',
                },
                {
                  type: 'doc',
                  id: 'start/trace/prometheus',
                },
              ],
            },
            {
              type: 'doc',
              id: 'start/wasm/start',
            },
            {
              type: 'doc',
              id: 'start/faas/start',
            }]
    },
    {
      type: 'doc',
      label: '线上实验室',
      id: 'start/lab'
    },
        {
          type: 'category',
          label: '用户手册',
          items: [
            {
              type: 'category',
              label: '功能介绍',
              items: [
                {
                  type: 'doc',
                  id: 'building_blocks/file/file',
                },
                {
                  type: 'doc',
                  id: 'building_blocks/actuator/actuator',
                },
                {
                  type: 'doc',
                  id: 'building_blocks/state/reference',
                },
                {
                  type: 'doc',
                  id: 'building_blocks/sequencer/reference',
                },
                {
                  type: 'doc',
                  id: 'building_blocks/lock/reference',
                },
                {
                  type: 'doc',
                  id: 'building_blocks/pubsub/reference',
                },
                {
                  type: 'doc',
                  id: 'building_blocks/rpc/reference',
                },
                {
                  type: 'doc',
                  id: 'building_blocks/configuration/reference',
                },
                {
                  type: 'category',
                  label: '可扩展性',
                  items: [
                    {
                      type: 'doc',
                      id: 'design/api_plugin/design',
                    },
                    {
                      type: 'doc',
                      id: 'design/pluggable/usage',
                    },
                  ],
                },
              ],
            },
            {
              "type": "doc",
              "id": "api_reference/README"
            },
            {
              "type": "link",
              label: 'java sdk',
              "href": "https://github.com/layotto/java-sdk"
            },
            {
              "type": "link",
              label: '.net sdk',

              "href": "https://github.com/layotto/dotnet-sdk"
            },
            {
              "type": "link",
              label: 'js sdk',
              "href": "https://github.com/layotto/js-sdk"
            },
            {
              "type": "doc",
              "id": "sdk_reference/go/start"
            }

          ],
        },

    //

{
  type: 'category',
  label: '运维手册',
  items: [
    {
      type: 'category',
      label: '如何配置 Layotto',
      items: [
        {
          type: 'doc',
          id: 'configuration/overview',
        },
        {
          type: 'category',
          label: '组件配置说明',
          link:{
            type:"doc",
            id:"component_specs/overview"
          },
          items: [
            {
              type: 'category',
              label: 'State',
              link:{
                type:"doc",
                id:"component_specs/state/common"
              },
              items: [
                {
                  type: 'doc',
                  id: 'component_specs/state/redis',
                },
                {
                  type: 'doc',
                  id: 'component_specs/state/others',
                },
              ],
            },
            {
              type: 'category',
              label: 'Pub/Sub',
              link:{
                type:"doc",
                id:"component_specs/pubsub/common"
              },
              items: [
                {
                  type: 'doc',
                  id: 'component_specs/pubsub/redis',
                },
                {
                  type: 'doc',
                  id: 'component_specs/pubsub/others',
                },
              ],
            },
            {
              type: 'category',
              label: 'Distributed Lock',
              link:{
                type:"doc",
                id:"component_specs/lock/common"
              },
              items: [
                {
                  type: 'doc',
                  id: 'component_specs/lock/redis',
                },
                {
                  type: 'doc',
                  id: 'component_specs/lock/etcd',
                },
                {
                  type: 'doc',
                  id: 'component_specs/lock/zookeeper',
                },
                {
                  type: 'doc',
                  id: 'component_specs/lock/consul',
                },
                {
                  type: 'doc',
                  id: 'component_specs/lock/mongo',
                },
              ],
            },
            {
              type: 'doc',
              id: 'component_specs/configuration/etcd',
            },
            {
              type: 'doc',
              id: 'component_specs/file/oss',
            },
            {
              type: 'category',
              label: 'Sequencer',
              items: [{
                "type": "doc",
                "id": "component_specs/sequencer/etcd"
              },
                {
                  "type": "doc",
                  "id": "component_specs/sequencer/redis"
                },
                {
                  "type": "doc",
                  "id": "component_specs/sequencer/zookeeper"
                },
                {
                  "type": "doc",
                  "id": "component_specs/sequencer/mongo"
                },
                {
                  "type": "doc",
                  "id": "component_specs/sequencer/mysql"
                },
                {
                  "type": "doc",
                  "id": "component_specs/sequencer/snowflake"
                }

              ],
            },
            {
              type: 'doc',
              id: 'component_specs/secret/common',
            },
            {
              type: 'doc',
              id: 'component_specs/custom/common',
            },
          ],
        },
      ],
    },
    {
      type: 'doc',
      id: 'operation/README',
    },
    {
      type: 'doc',
      id: 'operation/sidecar_injector',
    },
    {
      type: 'doc',
      id: 'operation/local',
    },
  ]
},


      //
    {
      type: 'category',
      label: '设计文档',
      items: [
        {
          type: 'doc',
          id: 'design/lifecycle/apply_configuration',
        },
        {
          type: 'doc',
          id: 'design/actuator/actuator-design-doc',
        },
        {
          type: 'doc',
          id: 'design/actuator/grpc-design-doc',
        },
        {
          type: 'doc',
          id: 'design/configuration/configuration-api-with-apollo',
        },
        {
          type: 'doc',
          id: 'design/pubsub/pubsub-api-and-compability-with-dapr-component',
        },
        {
          type: 'doc',
          id: 'design/rpc/rpc_design_document',
        },
        {
          type: 'doc',
          id: 'design/lock/lock-api-design',
        },
        {
          type: 'doc',
          id: 'design/sequencer/design',
        },
        {
          type: 'doc',
          id: 'design/file/file-design',
        },
        {
          type: 'doc',
          id: 'design/faas/faas-poc-design',
        },
        {
          type: 'doc',
          id: 'design/api_plugin/design',
        },
        {
          type: 'doc',
          id: 'design/api_plugin/dapr_api',
        },
        {
          type: 'doc',
          id: 'design/oss/design',
        },
        {
          type: 'doc',
          id: 'design/pluggable/design',
        },
      ],
    },
    {
      type: "category",
      label: "贡献指南",
      items: [
        "development/CONTRIBUTING",
        "development/start-from-zero",
        {
          type: "category",
          label: "想要贡献文档?",
          items: ["development/contributing-doc", "development/test-quickstart"]
        },
        "development/developing-component",
        "development/component_ref/component_ref",
        {
          type: "category",
          label: "想要修改proto文件或API定义？",
          items: [
            "api_reference/how_to_generate_api_doc",
            "api_reference/comment_spec_of_proto",
            "development/developing-api"
          ]
        },
        {
          type: "doc",
          id: "development/github-workflows"
        },
        {
          type: "doc",
          id: "development/commands"
        },
        {
          type: "category",
          label: "如何给 issue 打 label",
          items: [
            {
              type: "doc",
              id: "development/label-spec"
            }
          ]
        },
        {
          type: "doc",
          id: "development/release-guide"
        },
        {
          type: "doc",
          id: "development/problems-to-solve"
        }
      ]
    }
,
    {
      type: "category",
      label: "社区",
      items: [
        "community/meeting",
        "community/governance",
        "community/promote",
        "community/people"
      ]
    },

    {
      type: "doc",
      id: "video/README"
    }
  ]
};

export default sidebars;
