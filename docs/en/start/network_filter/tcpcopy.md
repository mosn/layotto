# Dump TCP traffic

## Introduction

When you run the demo according to the quick-start document [Configuration demo with apollo](en/start/configuration/start-apollo.md), you may notice that there is such a configuration in the configuration file config_apollo.json:

```json
                {
                  "type": "tcpcopy",
                  "config": {
                    "strategy": {
                      "switch": "ON",
                      "interval": 30,
                      "duration": 10,
                      "cpu_max_rate": 80,
                      "mem_max_rate": 80
                    }
                  }
```

The meaning of this configuration is to load the tcpcopy plug-in at startup to dump the tcp traffic.

After enabling this configuration, when Layotto receives a request and the conditions for traffic dump are met, it will write the binary request data to the local file system.

The "dumped" binary traffic data will be stored in the ${user's home directory}/logs/mosn directory, or under the /home/admin/logs/mosn directory:

![img.png](../../../img/tcp_dump.png)

You can use these data in combination with other tools and infrastructure to do something cool, such as traffic playback, bypass verification, etc.

## Configuration description

In the above config_apollo.json, the strategy configuration item is mainly used to configure the sampling strategy. The specific configuration descriptions are as follows:

```go
type DumpConfig struct {
	Switch     string  `json:"switch"`       // dump switch.'ON' or 'OFF'
	Interval   int     `json:"interval"`     // dump sampling interval, unit: second
	Duration   int     `json:"duration"`     // Single sampling duration,unit: second
	CpuMaxRate float64 `json:"cpu_max_rate"` // cpu max rate.When cpu rate bigger than this threshold,dump function will be fused
	MemMaxRate float64 `json:"mem_max_rate"` // mem max rate.When memory rate bigger than this threshold,dump function will be fused
}
```

## Principle of work

The Layotto server runs on MOSN and uses MOSN's filter expansion capabilities, so the tcpcopy above is actually a network filter plug-in of MOSN.

You can refer to [MOSN source code analysis-filter extension mechanism](https://mosn.io/blog/code/mosn-filters/) to implement your own 4-layer filter plug-in