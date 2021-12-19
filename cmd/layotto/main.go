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
	"fmt"
	"mosn.io/layotto/pkg/grpc/default_api"
	"os"
	"strconv"
	"time"

	mock_state "mosn.io/layotto/pkg/mock/components/state"

	"mosn.io/layotto/components/file/local"

	"mosn.io/layotto/components/file/s3/alicloud"
	"mosn.io/layotto/components/file/s3/aws"
	"mosn.io/layotto/components/file/s3/minio"

	dbindings "github.com/dapr/components-contrib/bindings"
	"github.com/dapr/components-contrib/bindings/http"
	"mosn.io/layotto/components/configstores/etcdv3"
	"mosn.io/layotto/components/file"
	"mosn.io/layotto/components/sequencer"
	"mosn.io/layotto/pkg/runtime/bindings"
	runtime_sequencer "mosn.io/layotto/pkg/runtime/sequencer"
	"mosn.io/pkg/log"

	// Hello
	"mosn.io/layotto/components/hello"
	"mosn.io/layotto/components/hello/helloworld"

	// Configuration
	"mosn.io/layotto/components/configstores"
	"mosn.io/layotto/components/configstores/apollo"

	// Pub/Sub
	dapr_comp_pubsub "github.com/dapr/components-contrib/pubsub"
	pubsub_snssqs "github.com/dapr/components-contrib/pubsub/aws/snssqs"
	pubsub_eventhubs "github.com/dapr/components-contrib/pubsub/azure/eventhubs"
	"github.com/dapr/components-contrib/pubsub/azure/servicebus"
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
	lock_mongo "mosn.io/layotto/components/lock/mongo"
	lock_redis "mosn.io/layotto/components/lock/redis"
	lock_zookeeper "mosn.io/layotto/components/lock/zookeeper"
	runtime_lock "mosn.io/layotto/pkg/runtime/lock"

	// Sequencer
	sequencer_etcd "mosn.io/layotto/components/sequencer/etcd"
	sequencer_redis "mosn.io/layotto/components/sequencer/redis"
	sequencer_zookeeper "mosn.io/layotto/components/sequencer/zookeeper"


	// Secret stores.
	"github.com/dapr/components-contrib/secretstores"
	"github.com/dapr/components-contrib/secretstores/aws/parameterstore"
	"github.com/dapr/components-contrib/secretstores/aws/secretmanager"
	"github.com/dapr/components-contrib/secretstores/azure/keyvault"
	gcp_secretmanager "github.com/dapr/components-contrib/secretstores/gcp/secretmanager"
	"github.com/dapr/components-contrib/secretstores/hashicorp/vault"
	sercetstores_kubernetes "github.com/dapr/components-contrib/secretstores/kubernetes"
	secretstore_env "github.com/dapr/components-contrib/secretstores/local/env"
	secretstore_file "github.com/dapr/components-contrib/secretstores/local/file"
	secretstores_loader "mosn.io/layotto/pkg/runtime/secretstores"

	// Actuator
	_ "mosn.io/layotto/pkg/actuator"
	"mosn.io/layotto/pkg/actuator/health"
	actuatorInfo "mosn.io/layotto/pkg/actuator/info"
	_ "mosn.io/layotto/pkg/filter/stream/actuator/http"
	"mosn.io/layotto/pkg/integrate/actuator"

	"github.com/urfave/cli"
	"google.golang.org/grpc"
	_ "mosn.io/layotto/pkg/filter/network/tcpcopy"
	"mosn.io/layotto/pkg/runtime"
	_ "mosn.io/layotto/pkg/wasm"
	"mosn.io/mosn/pkg/featuregate"
	_ "mosn.io/mosn/pkg/filter/network/grpc"
	mgrpc "mosn.io/mosn/pkg/filter/network/grpc"
	_ "mosn.io/mosn/pkg/filter/network/proxy"
	_ "mosn.io/mosn/pkg/filter/stream/flowcontrol"
	_ "mosn.io/mosn/pkg/metrics/sink"
	_ "mosn.io/mosn/pkg/metrics/sink/prometheus"
	"mosn.io/mosn/pkg/mosn"
	_ "mosn.io/mosn/pkg/network"
	_ "mosn.io/mosn/pkg/stream/http"
	_ "mosn.io/mosn/pkg/wasm/runtime/wasmer"
	_ "mosn.io/pkg/buffer"
)

// loggerForDaprComp is constructed for reusing dapr's components.
var loggerForDaprComp = logger.NewLogger("reuse.dapr.component")

func init() {
	mgrpc.RegisterServerHandler("runtime", NewRuntimeGrpcServer)
	// Register default actuator implementations
	actuatorInfo.AddInfoContributor("app", actuator.GetAppContributor())
	health.AddReadinessIndicator("runtime_startup", actuator.GetRuntimeReadinessIndicator())
	health.AddLivenessIndicator("runtime_startup", actuator.GetRuntimeLivenessIndicator())
}

func NewRuntimeGrpcServer(data json.RawMessage, opts ...grpc.ServerOption) (mgrpc.RegisteredServer, error) {
	// 1. parse config
	cfg, err := runtime.ParseRuntimeConfig(data)
	if err != nil {
		actuator.GetRuntimeReadinessIndicator().SetUnhealthy(fmt.Sprintf("parse config error.%v", err))
		return nil, err
	}
	// 2. new instance
	rt := runtime.NewMosnRuntime(cfg)
	// 3. run
	server, err := rt.Run(
		runtime.WithGrpcOptions(opts...),
		// register your grpc API here
		runtime.WithGrpcAPI(
			default_api.NewGrpcAPI,
		),
		// Hello
		runtime.WithHelloFactory(
			hello.NewHelloFactory("helloworld", helloworld.NewHelloWorld),
		),
		// Configuration
		runtime.WithConfigStoresFactory(
			configstores.NewStoreFactory("apollo", apollo.NewStore),
			configstores.NewStoreFactory("etcd", etcdv3.NewStore),
		),

		// RPC
		runtime.WithRpcFactory(
			rpc.NewRpcFactory("mosn", mosninvoker.NewMosnInvoker),
		),

		// File
		runtime.WithFileFactory(
			file.NewFileFactory("aliOSS", alicloud.NewAliCloudOSS),
			file.NewFileFactory("minioOSS", minio.NewMinioOss),
			file.NewFileFactory("awsOSS", aws.NewAwsOss),
			file.NewFileFactory("local", local.NewLocalStore),
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
		))

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
		)

	// 4. check if unhealthy
	if err != nil {
		actuator.GetRuntimeReadinessIndicator().SetUnhealthy(err.Error())
		actuator.GetRuntimeLivenessIndicator().SetUnhealthy(err.Error())
	}
	return server, err
}

var cmdStart = cli.Command{
	Name:  "start",
	Usage: "start runtime",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "Load configuration from `FILE`",
			EnvVar: "RUNTIME_CONFIG",
			Value:  "configs/config.json",
		}, cli.StringFlag{
			Name:   "feature-gates, f",
			Usage:  "config feature gates",
			EnvVar: "FEATURE_GATES",
		},
	},
	Action: func(c *cli.Context) error {
		stm := mosn.NewStageManager(c, c.String("config"))

		stm.AppendParamsParsedStage(func(c *cli.Context) {
			err := featuregate.Set(c.String("feature-gates"))
			if err != nil {
				os.Exit(1)
			}
		})

		stm.AppendInitStage(mosn.DefaultInitStage)

		stm.AppendPreStartStage(mosn.DefaultPreStartStage) // called finally stage by default

		stm.AppendStartStage(mosn.DefaultStartStage)

		stm.Run()

		actuator.GetRuntimeReadinessIndicator().SetStarted()
		actuator.GetRuntimeLivenessIndicator().SetStarted()
		// wait mosn finished
		stm.WaitFinish()
		return nil
	},
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
}

func newRuntimeApp(startCmd *cli.Command) *cli.App {
	app := cli.NewApp()
	app.Name = "Layotto"
	app.Version = "0.1.0"
	app.Compiled = time.Now()
	app.Copyright = "(c) " + strconv.Itoa(time.Now().Year()) + " Ant Group"
	app.Usage = "A fast and efficient cloud native application runtime based on MOSN."
	app.Flags = cmdStart.Flags

	// commands
	app.Commands = []cli.Command{
		cmdStart,
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
