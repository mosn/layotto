/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"encoding/json"
	_ "net/http/pprof"
	"os"
	"strconv"
	"time"

	"mosn.io/layotto/pkg/grpc/lifecycle"

	huaweicloud_oss "mosn.io/layotto/components/oss/huaweicloud"

	"mosn.io/layotto/components/oss"

	aws_oss "mosn.io/layotto/components/oss/aws"

	aliyun_oss "mosn.io/layotto/components/oss/aliyun"

	ceph_oss "mosn.io/layotto/components/oss/ceph"

	aliyun_file "mosn.io/layotto/components/file/aliyun"
	"mosn.io/layotto/components/file/local"

	"mosn.io/mosn/pkg/istio"

	"github.com/dapr/components-contrib/secretstores"
	"github.com/dapr/components-contrib/secretstores/aws/parameterstore"
	"github.com/dapr/components-contrib/secretstores/aws/secretmanager"
	"github.com/dapr/components-contrib/secretstores/azure/keyvault"
	gcp_secretmanager "github.com/dapr/components-contrib/secretstores/gcp/secretmanager"
	"github.com/dapr/components-contrib/secretstores/hashicorp/vault"
	sercetstores_kubernetes "github.com/dapr/components-contrib/secretstores/kubernetes"
	secretstore_env "github.com/dapr/components-contrib/secretstores/local/env"
	secretstore_file "github.com/dapr/components-contrib/secretstores/local/file"

	_ "mosn.io/layotto/pkg/filter/stream/wasm/http"
	"mosn.io/layotto/pkg/grpc/default_api"
	mock_state "mosn.io/layotto/pkg/mock/components/state"
	secretstores_loader "mosn.io/layotto/pkg/runtime/secretstores"
	_ "mosn.io/layotto/pkg/wasm"
	_ "mosn.io/layotto/pkg/wasm/install"
	_ "mosn.io/layotto/pkg/wasm/uninstall"
	_ "mosn.io/layotto/pkg/wasm/update"

	_ "mosn.io/mosn/pkg/filter/stream/grpcmetric"

	dbindings "github.com/dapr/components-contrib/bindings"
	"github.com/dapr/components-contrib/bindings/http"
	"mosn.io/pkg/log"

	"mosn.io/layotto/cmd/layotto_multiple_api/helloworld/component"
	"mosn.io/layotto/components/configstores/etcdv3"
	"mosn.io/layotto/components/custom"
	"mosn.io/layotto/components/file"
	aws_file "mosn.io/layotto/components/file/aws"
	"mosn.io/layotto/components/file/minio"
	"mosn.io/layotto/components/file/qiniu"
	"mosn.io/layotto/components/file/tencentcloud"
	"mosn.io/layotto/components/sequencer"
	"mosn.io/layotto/pkg/grpc/dapr"
	"mosn.io/layotto/pkg/runtime/bindings"
	runtime_sequencer "mosn.io/layotto/pkg/runtime/sequencer"

	// Hello
	"mosn.io/layotto/components/hello"
	"mosn.io/layotto/components/hello/helloworld"

	// Configuration
	"mosn.io/layotto/components/configstores"
	"mosn.io/layotto/components/configstores/apollo"
	"mosn.io/layotto/components/configstores/nacos"

	// Pub/Sub
	dapr_comp_pubsub "github.com/dapr/components-contrib/pubsub"
	pubsub_snssqs "github.com/dapr/components-contrib/pubsub/aws/snssqs"
	pubsub_eventhubs "github.com/dapr/components-contrib/pubsub/azure/eventhubs"
	pubsub_gcp "github.com/dapr/components-contrib/pubsub/gcp/pubsub"
	pubsub_hazelcast "github.com/dapr/components-contrib/pubsub/hazelcast"
	pubsub_inmemory "github.com/dapr/components-contrib/pubsub/in-memory"
	pubsub_kafka "github.com/dapr/components-contrib/pubsub/kafka"
	pubsub_mqtt "github.com/dapr/components-contrib/pubsub/mqtt"
	"github.com/dapr/components-contrib/pubsub/natsstreaming"
	pubsub_pulsar "github.com/dapr/components-contrib/pubsub/pulsar"
	"github.com/dapr/components-contrib/pubsub/rabbitmq"
	pubsub_redis "github.com/dapr/components-contrib/pubsub/redis"
	"github.com/dapr/kit/logger"

	"mosn.io/layotto/components/delay_queue/azure/servicebus"

	"mosn.io/layotto/pkg/runtime/pubsub"

	// RPC
	"mosn.io/layotto/components/rpc"
	mosninvoker "mosn.io/layotto/components/rpc/invoker/mosn"

	// State Stores
	"github.com/dapr/components-contrib/state"
	"github.com/dapr/components-contrib/state/aerospike"
	state_dynamodb "github.com/dapr/components-contrib/state/aws/dynamodb"
	state_azure_blobstorage "github.com/dapr/components-contrib/state/azure/blobstorage"
	state_cosmosdb "github.com/dapr/components-contrib/state/azure/cosmosdb"
	state_azure_tablestorage "github.com/dapr/components-contrib/state/azure/tablestorage"
	"github.com/dapr/components-contrib/state/cassandra"
	"github.com/dapr/components-contrib/state/cloudstate"
	"github.com/dapr/components-contrib/state/couchbase"
	"github.com/dapr/components-contrib/state/gcp/firestore"
	"github.com/dapr/components-contrib/state/hashicorp/consul"
	"github.com/dapr/components-contrib/state/hazelcast"
	"github.com/dapr/components-contrib/state/memcached"
	"github.com/dapr/components-contrib/state/mongodb"
	state_mysql "github.com/dapr/components-contrib/state/mysql"
	"github.com/dapr/components-contrib/state/postgresql"
	state_redis "github.com/dapr/components-contrib/state/redis"
	"github.com/dapr/components-contrib/state/rethinkdb"
	"github.com/dapr/components-contrib/state/sqlserver"
	"github.com/dapr/components-contrib/state/zookeeper"

	runtime_state "mosn.io/layotto/pkg/runtime/state"

	// Lock
	"mosn.io/layotto/components/lock"
	lock_consul "mosn.io/layotto/components/lock/consul"
	lock_etcd "mosn.io/layotto/components/lock/etcd"
	lock_inmemory "mosn.io/layotto/components/lock/in-memory"
	lock_mongo "mosn.io/layotto/components/lock/mongo"
	lock_redis "mosn.io/layotto/components/lock/redis"
	lock_zookeeper "mosn.io/layotto/components/lock/zookeeper"
	runtime_lock "mosn.io/layotto/pkg/runtime/lock"

	// Sequencer
	sequencer_etcd "mosn.io/layotto/components/sequencer/etcd"
	sequencer_inmemory "mosn.io/layotto/components/sequencer/in-memory"
	sequencer_mongo "mosn.io/layotto/components/sequencer/mongo"
	sequencer_mysql "mosn.io/layotto/components/sequencer/mysql"
	sequencer_redis "mosn.io/layotto/components/sequencer/redis"
	sequencer_snowflake "mosn.io/layotto/components/sequencer/snowflake"
	sequencer_zookeeper "mosn.io/layotto/components/sequencer/zookeeper"

	// Actuator
	_ "mosn.io/layotto/pkg/actuator"
	"mosn.io/layotto/pkg/actuator/health"
	actuatorInfo "mosn.io/layotto/pkg/actuator/info"
	_ "mosn.io/layotto/pkg/filter/stream/actuator/http"
	"mosn.io/layotto/pkg/integrate/actuator"

	"github.com/urfave/cli"
	"google.golang.org/grpc"
	_ "mosn.io/mosn/pkg/filter/network/grpc"
	mgrpc "mosn.io/mosn/pkg/filter/network/grpc"
	_ "mosn.io/mosn/pkg/filter/network/proxy"
	_ "mosn.io/mosn/pkg/filter/stream/flowcontrol"
	_ "mosn.io/mosn/pkg/metrics/sink"
	_ "mosn.io/mosn/pkg/metrics/sink/prometheus"
	_ "mosn.io/mosn/pkg/network"
	_ "mosn.io/mosn/pkg/stream/http"
	_ "mosn.io/mosn/pkg/wasm/runtime/wasmer"
	_ "mosn.io/pkg/buffer"

	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	_ "mosn.io/mosn/istio/istio1106"
	_ "mosn.io/mosn/istio/istio1106/filter/stream/jwtauthn"
	_ "mosn.io/mosn/istio/istio1106/filter/stream/mixer"
	_ "mosn.io/mosn/istio/istio1106/filter/stream/stats"
	_ "mosn.io/mosn/istio/istio1106/sds"
	_ "mosn.io/mosn/istio/istio1106/xds"
	_ "mosn.io/mosn/pkg/filter/listener/originaldst"
	_ "mosn.io/mosn/pkg/filter/network/connectionmanager"
	_ "mosn.io/mosn/pkg/filter/network/streamproxy"
	_ "mosn.io/mosn/pkg/filter/network/tunnel"
	_ "mosn.io/mosn/pkg/filter/stream/dsl"
	_ "mosn.io/mosn/pkg/filter/stream/dubbo"
	_ "mosn.io/mosn/pkg/filter/stream/faultinject"
	_ "mosn.io/mosn/pkg/filter/stream/faulttolerance"
	_ "mosn.io/mosn/pkg/filter/stream/gzip"
	_ "mosn.io/mosn/pkg/filter/stream/headertometadata"
	_ "mosn.io/mosn/pkg/filter/stream/ipaccess"
	_ "mosn.io/mosn/pkg/filter/stream/mirror"
	_ "mosn.io/mosn/pkg/filter/stream/payloadlimit"
	_ "mosn.io/mosn/pkg/filter/stream/proxywasm"
	_ "mosn.io/mosn/pkg/filter/stream/seata"
	_ "mosn.io/mosn/pkg/filter/stream/transcoder/http2bolt"
	_ "mosn.io/mosn/pkg/filter/stream/transcoder/httpconv"
	_ "mosn.io/mosn/pkg/protocol"
	_ "mosn.io/mosn/pkg/protocol/xprotocol"
	_ "mosn.io/mosn/pkg/router"
	_ "mosn.io/mosn/pkg/server/keeper"
	_ "mosn.io/mosn/pkg/stream/http2"
	_ "mosn.io/mosn/pkg/stream/xprotocol"
	_ "mosn.io/mosn/pkg/trace/jaeger"
	_ "mosn.io/mosn/pkg/trace/skywalking"
	_ "mosn.io/mosn/pkg/trace/skywalking/http"
	_ "mosn.io/mosn/pkg/trace/sofa/http"
	_ "mosn.io/mosn/pkg/trace/sofa/xprotocol"
	_ "mosn.io/mosn/pkg/trace/sofa/xprotocol/bolt"
	_ "mosn.io/mosn/pkg/upstream/healthcheck"
	_ "mosn.io/mosn/pkg/upstream/servicediscovery/dubbod"

	_ "mosn.io/layotto/pkg/filter/network/tcpcopy"
	l8_grpc "mosn.io/layotto/pkg/grpc"
	"mosn.io/layotto/pkg/runtime"

	helloworld_api "mosn.io/layotto/cmd/layotto_multiple_api/helloworld"
	_ "mosn.io/layotto/diagnostics/exporter_iml"
)

// loggerForDaprComp is constructed for reusing dapr's components.
var loggerForDaprComp = logger.NewLogger("reuse.dapr.component")

// GitVersion mosn version is specified by latest tag
var GitVersion = ""
var IstioVersion = "1.10.6"

func init() {
	mgrpc.RegisterServerHandler("runtime", NewRuntimeGrpcServer)
	// Register default actuator implementations
	actuatorInfo.AddInfoContributor("app", actuator.GetAppContributor())
	health.AddReadinessIndicator("runtime_startup", actuator.GetRuntimeReadinessIndicator())
	health.AddLivenessIndicator("runtime_startup", actuator.GetRuntimeLivenessIndicator())
}

func NewRuntimeGrpcServer(data json.RawMessage, opts ...grpc.ServerOption) (mgrpc.RegisteredServer, error) {
	var err error
	defer func() {
		if err != nil {
			// fail fast if error occurs during startup.
			// The reason we panic in a new goroutine is to prevent mosn from recovering.
			go func() {
				log.DefaultLogger.Errorf("An error occurred during startup : %v", err)
				panic(err)
			}()
		}
	}()
	// 1. parse config
	cfg, err := runtime.ParseRuntimeConfig(data)
	if err != nil {
		return nil, err
	}
	// 2. new instance
	rt := runtime.NewMosnRuntime(cfg)
	rt.AppendInitRuntimeStage(runtime.DefaultInitRuntimeStage)
	// 3. run
	server, err := rt.Run(
		runtime.WithGrpcOptions(opts...),
		// wrap the grpc server with actuator
		runtime.WithNewServer(func(apis []l8_grpc.GrpcAPI, opts ...grpc.ServerOption) (mgrpc.RegisteredServer, error) {
			server, err := l8_grpc.NewDefaultServer(apis, opts...)
			if err != nil {
				return nil, err
			}
			return actuator.NewGrpcServerWithActuator(server)
		}),
		// register your gRPC API here
		runtime.WithGrpcAPI(
			// default GrpcAPI
			default_api.NewGrpcAPI,
			lifecycle.NewLifecycleAPI,

			// a demo to show how to register your own gRPC API
			helloworld_api.NewHelloWorldAPI,

			// support Dapr API
			// Currently it only support Dapr's InvokeService,secret API,state API and InvokeBinding API.
			// Note: this feature is still in Alpha state and we don't recommend that you use it in your production environment.
			dapr.NewDaprAPI_Alpha,
		),
		runtime.WithExtensionGrpcAPI(),
		// Hello
		runtime.WithHelloFactory(
			hello.NewHelloFactory("helloworld", helloworld.NewHelloWorld),
		),
		// Configuration
		runtime.WithConfigStoresFactory(
			configstores.NewStoreFactory("apollo", apollo.NewStore),
			configstores.NewStoreFactory("etcd", etcdv3.NewStore),
			configstores.NewStoreFactory("nacos", nacos.NewStore),
		),

		// RPC
		runtime.WithRpcFactory(
			rpc.NewRpcFactory("mosn", mosninvoker.NewMosnInvoker),
		),

		// File
		runtime.WithFileFactory(
			file.NewFileFactory("aliyun.oss", aliyun_file.NewAliyunFile),
			file.NewFileFactory("minio", minio.NewMinioOss),
			file.NewFileFactory("aws.s3", aws_file.NewAwsFile),
			file.NewFileFactory("tencent.oss", tencentcloud.NewTencentCloudOSS),
			file.NewFileFactory("local", local.NewLocalStore),
			file.NewFileFactory("qiniu.oss", qiniu.NewQiniuOSS),
		),

		//OSS
		runtime.WithOssFactory(
			oss.NewFactory("aws.oss", aws_oss.NewAwsOss),
			oss.NewFactory("aliyun.oss", aliyun_oss.NewAliyunOss),
			oss.NewFactory("ceph", ceph_oss.NewCephOss),
			oss.NewFactory("huaweicloud.oss", huaweicloud_oss.NewHuaweicloudOSS),
		),

		// PubSub
		runtime.WithPubSubFactory(
			pubsub.NewFactory("redis", func() dapr_comp_pubsub.PubSub {
				return pubsub_redis.NewRedisStreams(loggerForDaprComp)
			}),
			pubsub.NewFactory("natsstreaming", func() dapr_comp_pubsub.PubSub {
				return natsstreaming.NewNATSStreamingPubSub(loggerForDaprComp)
			}),
			pubsub.NewFactory("azure.eventhubs", func() dapr_comp_pubsub.PubSub {
				return pubsub_eventhubs.NewAzureEventHubs(loggerForDaprComp)
			}),
			pubsub.NewFactory("azure.servicebus", func() dapr_comp_pubsub.PubSub {
				return servicebus.NewAzureServiceBus(loggerForDaprComp)
			}),
			pubsub.NewFactory("rabbitmq", func() dapr_comp_pubsub.PubSub {
				return rabbitmq.NewRabbitMQ(loggerForDaprComp)
			}),
			pubsub.NewFactory("hazelcast", func() dapr_comp_pubsub.PubSub {
				return pubsub_hazelcast.NewHazelcastPubSub(loggerForDaprComp)
			}),
			pubsub.NewFactory("gcp.pubsub", func() dapr_comp_pubsub.PubSub {
				return pubsub_gcp.NewGCPPubSub(loggerForDaprComp)
			}),
			pubsub.NewFactory("kafka", func() dapr_comp_pubsub.PubSub {
				return pubsub_kafka.NewKafka(loggerForDaprComp)
			}),
			pubsub.NewFactory("snssqs", func() dapr_comp_pubsub.PubSub {
				return pubsub_snssqs.NewSnsSqs(loggerForDaprComp)
			}),
			pubsub.NewFactory("mqtt", func() dapr_comp_pubsub.PubSub {
				return pubsub_mqtt.NewMQTTPubSub(loggerForDaprComp)
			}),
			pubsub.NewFactory("pulsar", func() dapr_comp_pubsub.PubSub {
				return pubsub_pulsar.NewPulsar(loggerForDaprComp)
			}),
			pubsub.NewFactory("in-memory", func() dapr_comp_pubsub.PubSub {
				return pubsub_inmemory.New(loggerForDaprComp)
			}),
		),
		// State
		runtime.WithStateFactory(
			runtime_state.NewFactory("in-memory", func() state.Store {
				return mock_state.New(loggerForDaprComp)
			}),
			runtime_state.NewFactory("redis", func() state.Store {
				return state_redis.NewRedisStateStore(loggerForDaprComp)
			}),
			runtime_state.NewFactory("consul", func() state.Store {
				return consul.NewConsulStateStore(loggerForDaprComp)
			}),
			runtime_state.NewFactory("azure.blobstorage", func() state.Store {
				return state_azure_blobstorage.NewAzureBlobStorageStore(loggerForDaprComp)
			}),
			runtime_state.NewFactory("azure.cosmosdb", func() state.Store {
				return state_cosmosdb.NewCosmosDBStateStore(loggerForDaprComp)
			}),
			runtime_state.NewFactory("azure.tablestorage", func() state.Store {
				return state_azure_tablestorage.NewAzureTablesStateStore(loggerForDaprComp)
			}),
			runtime_state.NewFactory("cassandra", func() state.Store {
				return cassandra.NewCassandraStateStore(loggerForDaprComp)
			}),
			runtime_state.NewFactory("memcached", func() state.Store {
				return memcached.NewMemCacheStateStore(loggerForDaprComp)
			}),
			runtime_state.NewFactory("mongodb", func() state.Store {
				return mongodb.NewMongoDB(loggerForDaprComp)
			}),
			runtime_state.NewFactory("zookeeper", func() state.Store {
				return zookeeper.NewZookeeperStateStore(loggerForDaprComp)
			}),
			runtime_state.NewFactory("gcp.firestore", func() state.Store {
				return firestore.NewFirestoreStateStore(loggerForDaprComp)
			}),
			runtime_state.NewFactory("postgresql", func() state.Store {
				return postgresql.NewPostgreSQLStateStore(loggerForDaprComp)
			}),
			runtime_state.NewFactory("sqlserver", func() state.Store {
				return sqlserver.NewSQLServerStateStore(loggerForDaprComp)
			}),
			runtime_state.NewFactory("hazelcast", func() state.Store {
				return hazelcast.NewHazelcastStore(loggerForDaprComp)
			}),
			runtime_state.NewFactory("cloudstate.crdt", func() state.Store {
				return cloudstate.NewCRDT(loggerForDaprComp)
			}),
			runtime_state.NewFactory("couchbase", func() state.Store {
				return couchbase.NewCouchbaseStateStore(loggerForDaprComp)
			}),
			runtime_state.NewFactory("aerospike", func() state.Store {
				return aerospike.NewAerospikeStateStore(loggerForDaprComp)
			}),
			runtime_state.NewFactory("rethinkdb", func() state.Store {
				return rethinkdb.NewRethinkDBStateStore(loggerForDaprComp)
			}),
			runtime_state.NewFactory("aws.dynamodb", state_dynamodb.NewDynamoDBStateStore),
			runtime_state.NewFactory("mysql", func() state.Store {
				return state_mysql.NewMySQLStateStore(loggerForDaprComp)
			}),
		),
		// Lock
		runtime.WithLockFactory(
			runtime_lock.NewFactory("redis_cluster", func() lock.LockStore {
				return lock_redis.NewClusterRedisLock(log.DefaultLogger)
			}),
			runtime_lock.NewFactory("redis", func() lock.LockStore {
				return lock_redis.NewStandaloneRedisLock(log.DefaultLogger)
			}),
			runtime_lock.NewFactory("zookeeper", func() lock.LockStore {
				return lock_zookeeper.NewZookeeperLock(log.DefaultLogger)
			}),
			runtime_lock.NewFactory("etcd", func() lock.LockStore {
				return lock_etcd.NewEtcdLock(log.DefaultLogger)
			}),
			runtime_lock.NewFactory("consul", func() lock.LockStore {
				return lock_consul.NewConsulLock(log.DefaultLogger)
			}),
			runtime_lock.NewFactory("mongo", func() lock.LockStore {
				return lock_mongo.NewMongoLock(log.DefaultLogger)
			}),
			runtime_lock.NewFactory("in-memory", func() lock.LockStore {
				return lock_inmemory.NewInMemoryLock()
			}),
		),

		// bindings
		runtime.WithOutputBindings(
			bindings.NewOutputBindingFactory("http", func() dbindings.OutputBinding {
				return http.NewHTTP(loggerForDaprComp)
			}),
		),

		// Sequencer
		runtime.WithSequencerFactory(
			runtime_sequencer.NewFactory("etcd", func() sequencer.Store {
				return sequencer_etcd.NewEtcdSequencer(log.DefaultLogger)
			}),
			runtime_sequencer.NewFactory("redis", func() sequencer.Store {
				return sequencer_redis.NewStandaloneRedisSequencer(log.DefaultLogger)
			}),
			runtime_sequencer.NewFactory("zookeeper", func() sequencer.Store {
				return sequencer_zookeeper.NewZookeeperSequencer(log.DefaultLogger)
			}),
			runtime_sequencer.NewFactory("mongo", func() sequencer.Store {
				return sequencer_mongo.NewMongoSequencer(log.DefaultLogger)
			}),
			runtime_sequencer.NewFactory("in-memory", func() sequencer.Store {
				return sequencer_inmemory.NewInMemorySequencer()
			}),
			runtime_sequencer.NewFactory("mysql", func() sequencer.Store {
				return sequencer_mysql.NewMySQLSequencer(log.DefaultLogger)
			}),
			runtime_sequencer.NewFactory("snowflake", func() sequencer.Store {
				return sequencer_snowflake.NewSnowFlakeSequencer(log.DefaultLogger)
			}),
		),
		// secretstores
		runtime.WithSecretStoresFactory(
			secretstores_loader.NewFactory("kubernetes", func() secretstores.SecretStore {
				return sercetstores_kubernetes.NewKubernetesSecretStore(loggerForDaprComp)
			}),
			secretstores_loader.NewFactory("azure.keyvault", func() secretstores.SecretStore {
				return keyvault.NewAzureKeyvaultSecretStore(loggerForDaprComp)
			}),
			secretstores_loader.NewFactory("hashicorp.vault", func() secretstores.SecretStore {
				return vault.NewHashiCorpVaultSecretStore(loggerForDaprComp)
			}),
			secretstores_loader.NewFactory("aws.secretmanager", func() secretstores.SecretStore {
				return secretmanager.NewSecretManager(loggerForDaprComp)
			}),
			secretstores_loader.NewFactory("aws.parameterstore", func() secretstores.SecretStore {
				return parameterstore.NewParameterStore(loggerForDaprComp)
			}),
			secretstores_loader.NewFactory("gcp.secretmanager", func() secretstores.SecretStore {
				return gcp_secretmanager.NewSecreteManager(loggerForDaprComp)
			}),
			secretstores_loader.NewFactory("local.file", func() secretstores.SecretStore {
				return secretstore_file.NewLocalSecretStore(loggerForDaprComp)
			}),
			secretstores_loader.NewFactory("local.env", func() secretstores.SecretStore {
				return secretstore_env.NewEnvSecretStore(loggerForDaprComp)
			}),
		), // Custom components
		runtime.WithCustomComponentFactory("helloworld",
			custom.NewComponentFactory("in-memory", component.NewInMemoryHelloWorld),
			custom.NewComponentFactory("goodbye", component.NewSayGoodbyeHelloWorld),
		),
	)
	return server, err
}

func main() {
	app := newRuntimeApp(&cmdStart)
	registerAppInfo(app)
	_ = app.Run(os.Args)
}

func registerAppInfo(app *cli.App) {
	appInfo := actuator.NewAppInfo()
	appInfo.Name = app.Name
	appInfo.Version = app.Version
	appInfo.Compiled = app.Compiled
	actuator.SetAppInfoSingleton(appInfo)
	// set istio version
	istio.IstioVersion = IstioVersion
}

func newRuntimeApp(startCmd *cli.Command) *cli.App {
	app := cli.NewApp()
	app.Name = "Layotto"
	app.Version = GitVersion
	app.Compiled = time.Now()
	app.Copyright = "(c) " + strconv.Itoa(time.Now().Year()) + " Layotto Authors"
	app.Usage = "A fast and efficient cloud native application runtime based on MOSN."
	app.Flags = cmdStart.Flags

	// commands
	app.Commands = []cli.Command{
		cmdStart,
		cmdStop,
		cmdReload,
	}
	// action
	app.Action = func(c *cli.Context) error {
		if c.NumFlags() == 0 {
			return cli.ShowAppHelp(c)
		}

		return startCmd.Action.(func(c *cli.Context) error)(c)
	}

	return app
}
