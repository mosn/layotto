package main

import (
	"encoding/json"
	"fmt"
	"github.com/layotto/layotto/pkg/integrate/actuator"
	"github.com/layotto/layotto/pkg/services/configstores"
	"github.com/layotto/layotto/pkg/services/configstores/apollo"
	"github.com/layotto/layotto/pkg/services/configstores/etcdv3"
	"os"
	"strconv"
	"time"

	_ "github.com/layotto/layotto/pkg/actuator"
	health "github.com/layotto/layotto/pkg/actuator/health"
	actuatorInfo "github.com/layotto/layotto/pkg/actuator/info"
	_ "github.com/layotto/layotto/pkg/filter/network/tcpcopy"
	_ "github.com/layotto/layotto/pkg/filter/stream/actuator/http"
	"github.com/layotto/layotto/pkg/runtime"
	"github.com/layotto/layotto/pkg/services/hello"
	"github.com/layotto/layotto/pkg/services/hello/helloworld"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
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
	_ "mosn.io/pkg/buffer"
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
		runtime.WithHelloFactory(
			hello.NewHelloFactory("helloworld", helloworld.NewHelloWorld),
		),
		runtime.WithConfigStoresFactory(
			configstores.NewStoreFactory("etcd", etcdv3.NewStore),
			configstores.NewStoreFactory("apollo", apollo.NewStore),
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
	app.Name = "LayOtto"
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
