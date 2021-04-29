package main

import (
	"encoding/json"
	"gitlab.alipay-inc.com/ant-mesh/runtime/pkg/services/configstores"
	"gitlab.alipay-inc.com/ant-mesh/runtime/pkg/services/configstores/apollo"
	"gitlab.alipay-inc.com/ant-mesh/runtime/pkg/services/configstores/etcdv3"
	"os"
	"strconv"
	"time"

	"github.com/urfave/cli"
	_ "gitlab.alipay-inc.com/ant-mesh/runtime/pkg/filter/network/tcpcopy"
	"gitlab.alipay-inc.com/ant-mesh/runtime/pkg/runtime"
	"gitlab.alipay-inc.com/ant-mesh/runtime/pkg/services/hello"
	"gitlab.alipay-inc.com/ant-mesh/runtime/pkg/services/hello/helloworld"
	"google.golang.org/grpc"
	"mosn.io/mosn/pkg/featuregate"
	_ "mosn.io/mosn/pkg/filter/network/grpc"
	mgrpc "mosn.io/mosn/pkg/filter/network/grpc"
	_ "mosn.io/mosn/pkg/filter/stream/flowcontrol"
	_ "mosn.io/mosn/pkg/metrics/sink"
	_ "mosn.io/mosn/pkg/metrics/sink/prometheus"
	"mosn.io/mosn/pkg/mosn"
	_ "mosn.io/mosn/pkg/network"
	_ "mosn.io/pkg/buffer"
)

// TODO: 开源以后完善,目前就做一个example
func init() {
	mgrpc.RegisterServerHandler("runtime", NewRuntimeGrpcServer)
}

func NewRuntimeGrpcServer(data json.RawMessage, opts ...grpc.ServerOption) (mgrpc.RegisteredServer, error) {
	cfg, err := runtime.ParseRuntimeConfig(data)
	if err != nil {
		return nil, err
	}
	rt := runtime.NewMosnRuntime(cfg)
	return rt.Run(
		runtime.WithGrpcOptions(opts...),
		runtime.WithHelloFactory(
			hello.NewHelloFactory("helloworld", helloworld.NewHelloWorld),
		),
		runtime.WithConfigStoresFactory(
			configstores.NewStoreFactory("etcd", etcdv3.NewStore),
			configstores.NewStoreFactory("apollo", apollo.NewStore),
		),
	)
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

			// wait mosn finished
			stm.WaitFinish()
			return nil
		},
	}
)

func main() {
	app := newRuntimeApp(&cmdStart)
	_ = app.Run(os.Args)
}

func newRuntimeApp(startCmd *cli.Command) *cli.App {
	app := cli.NewApp()
	app.Name = "runtime"
	app.Version = "0.1.0"
	app.Compiled = time.Now()
	app.Copyright = "(c) " + strconv.Itoa(time.Now().Year()) + " Ant Group"
	app.Usage = "mosn based runtime"
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
