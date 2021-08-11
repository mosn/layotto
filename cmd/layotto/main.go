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
	"os"
	"strconv"
	"time"

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
	lock_redis "mosn.io/layotto/components/lock/redis"
	lock_zookeeper "mosn.io/layotto/components/lock/zookeeper"
	runtime_lock "mosn.io/layotto/pkg/runtime/lock"

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

var (
	// loggerForDaprComp is constructed for reusing dapr's components.
	loggerForDaprComp = logger.NewLogger("reuse.dapr.component")
)

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
		// Hello
		runtime.WithHelloFactory(
			hello.NewHelloFactory("helloworld", helloworld.NewHelloWorld),
		),
		// Configuration
		runtime.WithConfigStoresFactory(
			configstores.NewStoreFactory("apollo", apollo.NewStore),
		),
		// RPC
		runtime.WithRpcFactory(
			rpc.NewRpcFactory("mosn", mosninvoker.NewMosnInvoker),
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
		),
		// State
		runtime.WithStateFactory(
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
			runtime_lock.NewFactory("redis", func() lock.LockStore {
				return lock_redis.NewStandaloneRedisLock(log.DefaultLogger)
			}),
			runtime_lock.NewFactory("zookeeper", func() lock.LockStore {
				return lock_zookeeper.NewZookeeperLock(log.DefaultLogger)
			}),
		),
	)

	// 4. check if unhealthy
	if err != nil {
		actuator.GetRuntimeReadinessIndicator().SetUnhealthy(err.Error())
		actuator.GetRuntimeLivenessIndicator().SetUnhealthy(err.Error())
	}
	return server, err
}

var (
	cmdStart = cli.Command{
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
)

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

	//commands
	app.Commands = []cli.Command{
		cmdStart,
	}
	//action
	app.Action = func(c *cli.Context) error {
		if c.NumFlags() == 0 {
			return cli.ShowAppHelp(c)
		}

		return startCmd.Action.(func(c *cli.Context) error)(c)
	}

	return app
}
