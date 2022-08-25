package main

import (
	"os"
	"runtime"

	"mosn.io/api"
	v2 "mosn.io/mosn/pkg/config/v2"
	"mosn.io/mosn/pkg/configmanager"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/metrics"
	"mosn.io/mosn/pkg/mosn"
	"mosn.io/mosn/pkg/protocol"
	"mosn.io/mosn/pkg/protocol/xprotocol"
	"mosn.io/mosn/pkg/protocol/xprotocol/bolt"
	"mosn.io/mosn/pkg/protocol/xprotocol/dubbo"
	"mosn.io/mosn/pkg/stagemanager"
	xstream "mosn.io/mosn/pkg/stream/xprotocol"
	"mosn.io/mosn/pkg/trace"
	mosn_jaeger "mosn.io/mosn/pkg/trace/jaeger"
	"mosn.io/mosn/pkg/trace/skywalking"
	tracehttp "mosn.io/mosn/pkg/trace/sofa/http"
	xtrace "mosn.io/mosn/pkg/trace/sofa/xprotocol"
	tracebolt "mosn.io/mosn/pkg/trace/sofa/xprotocol/bolt"
	mosn_zipkin "mosn.io/mosn/pkg/trace/zipkin"
	"mosn.io/pkg/buffer"

	component_actuators "mosn.io/layotto/components/pkg/actuators"
	"mosn.io/layotto/diagnostics"
	"mosn.io/layotto/diagnostics/jaeger"
	lprotocol "mosn.io/layotto/diagnostics/protocol"
	lsky "mosn.io/layotto/diagnostics/skywalking"
	"mosn.io/layotto/diagnostics/zipkin"

	// Actuator
	"mosn.io/layotto/pkg/actuator/health"
	"mosn.io/layotto/pkg/integrate/actuator"

	"github.com/urfave/cli"
	"mosn.io/mosn/pkg/featuregate"
)

var (
	flagToMosnLogLevel = map[string]string{
		"trace":    "TRACE",
		"debug":    "DEBUG",
		"info":     "INFO",
		"warning":  "WARN",
		"error":    "ERROR",
		"critical": "FATAL",
		"off":      "OFF",
	}

	cmdStart = cli.Command{
		Name:  "start",
		Usage: "start runtime. For example:  ./layotto start -c configs/config_standalone.json",
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
			}, cli.StringFlag{
				Name:   "log-level, l",
				Usage:  "mosn log level, trace|debug|info|warning|error|critical|off",
				EnvVar: "LOG_LEVEL",
			},
		},
		Action: func(c *cli.Context) error {
			app := mosn.NewMosn()
			stm := stagemanager.InitStageManager(c, c.String("config"), app)

			// if needs featuregate init in parameter stage or init stage
			// append a new stage and called featuregate.ExecuteInitFunc(keys...)
			// parameter parsed registered
			stm.AppendParamsParsedStage(ExtensionsRegister)

			stm.AppendParamsParsedStage(DefaultParamsParsed)

			// init Stage
			stm.AppendInitStage(mosn.DefaultInitStage)
			stm.AppendInitStage(func(_ *v2.MOSNConfig) {
				// set version and go version
				metrics.SetVersion(GitVersion)
				metrics.SetGoVersion(runtime.Version())
			})
			// pre-startup
			stm.AppendPreStartStage(mosn.DefaultPreStartStage) // called finally stage by default
			// startup
			stm.AppendStartStage(mosn.DefaultStartStage)
			// after-startup
			stm.AppendAfterStartStage(SetActuatorAfterStart)
			// execute all stages
			stm.RunAll()
			return nil
		},
	}

	cmdStop = cli.Command{
		Name:  "stop",
		Usage: "stop mosn proxy",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "config, c",
				Usage:  "load configuration from `FILE`",
				EnvVar: "MOSN_CONFIG",
				Value:  "configs/mosn_config.json",
			},
		},
		Action: func(c *cli.Context) (err error) {
			app := mosn.NewMosn()
			stm := stagemanager.InitStageManager(c, c.String("config"), app)
			stm.AppendInitStage(mosn.InitDefaultPath)
			return stm.StopMosnProcess()
		},
	}

	cmdReload = cli.Command{
		Name:  "reload",
		Usage: "reconfiguration",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
)

func SetActuatorAfterStart(_ stagemanager.Application) {
	// register component actuator
	component_actuators.RangeAllIndicators(
		func(name string, v *component_actuators.ComponentsIndicator) bool {
			if v != nil {
				health.AddLivenessIndicator(name, v.LivenessIndicator)
				health.AddReadinessIndicator(name, v.ReadinessIndicator)
			}
			return true
		})
	// set started
	actuator.GetRuntimeReadinessIndicator().SetStarted()
	actuator.GetRuntimeLivenessIndicator().SetStarted()
}

func DefaultParamsParsed(c *cli.Context) {
	// log level control
	flagLogLevel := c.String("log-level")
	if mosnLogLevel, ok := flagToMosnLogLevel[flagLogLevel]; ok {
		if mosnLogLevel == "OFF" {
			log.GetErrorLoggerManagerInstance().Disable()
		} else {
			log.GetErrorLoggerManagerInstance().SetLogLevelControl(configmanager.ParseLogLevel(mosnLogLevel))
		}
	}
	// set feature gates
	err := featuregate.Set(c.String("feature-gates"))
	if err != nil {
		log.StartLogger.Infof("[mosn] [start] parse feature-gates flag fail : %+v", err)
		os.Exit(1)
	}
}

// ExtensionsRegister for register mosn rpc extensions
func ExtensionsRegister(_ *cli.Context) {
	// 1. tracer driver register
	// Q: What is a tracer driver ?
	// A: MOSN implement a group of trace drivers, but only a configured driver will be loaded.
	//	A tracer driver can create different tracer by different protocol.
	//	When MOSN receive a request stream, MOSN will try to start a tracer according to the request protocol.
	// 	For more details,see https://mosn.io/blog/posts/skywalking-support/
	trace.RegisterDriver("SOFATracer", trace.NewDefaultDriverImpl())

	// 2. xprotocol action register
	// RegisterXProtocolAction is MOSN's xprotocol framework's extensions.
	// when a xprotocol implementation (defined by api.XProtocolCodec) registered, the registered action will be called.
	xprotocol.RegisterXProtocolAction(xstream.NewConnPool, xstream.NewStreamFactory, func(codec api.XProtocolCodec) {
		name := codec.ProtocolName()
		trace.RegisterTracerBuilder("SOFATracer", name, xtrace.NewTracer)
	})

	// 3. register protocols that are used by layotto.
	// RegisterXProtocolCodec add a new xprotocol implementation, which is a wrapper for protocol register
	_ = xprotocol.RegisterXProtocolCodec(&bolt.XCodec{})
	_ = xprotocol.RegisterXProtocolCodec(&dubbo.XCodec{})

	// 4. register tracer
	xtrace.RegisterDelegate(bolt.ProtocolName, tracebolt.Boltv1Delegate)
	trace.RegisterTracerBuilder("SOFATracer", protocol.HTTP1, tracehttp.NewTracer)
	trace.RegisterTracerBuilder("SOFATracer", lprotocol.Layotto, diagnostics.NewTracer)
	trace.RegisterTracerBuilder(skywalking.SkyDriverName, lprotocol.Layotto, lsky.NewGrpcSkyTracer)
	trace.RegisterTracerBuilder(mosn_jaeger.DriverName, lprotocol.Layotto, jaeger.NewGrpcJaegerTracer)
	trace.RegisterTracerBuilder(mosn_zipkin.DriverName, lprotocol.Layotto, zipkin.NewGrpcZipTracer)

	// register buffer logger
	buffer.SetLogFunc(func(msg string) {
		log.DefaultLogger.Errorf("[iobuffer] iobuffer error log info: %s", msg)
	})
}
