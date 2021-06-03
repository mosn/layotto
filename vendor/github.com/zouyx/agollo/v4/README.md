Agollo - Go Client for Apollo
================

[![golang](https://img.shields.io/badge/Language-Go-green.svg?style=flat)](https://golang.org)
[![Build Status](https://travis-ci.org/zouyx/agollo.svg?branch=master)](https://travis-ci.org/zouyx/agollo)
[![Go Report Card](https://goreportcard.com/badge/github.com/zouyx/agollo)](https://goreportcard.com/report/github.com/zouyx/agollo)
[![codebeat badge](https://codebeat.co/badges/bc2009d6-84f1-4f11-803e-fc571a12a1c0)](https://codebeat.co/projects/github-com-zouyx-agollo-master)
[![Coverage Status](https://coveralls.io/repos/github/zouyx/agollo/badge.svg?branch=master)](https://coveralls.io/github/zouyx/agollo?branch=master)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc](http://godoc.org/github.com/zouyx/agollo?status.svg)](http://godoc.org/github.com/zouyx/agollo)
[![GitHub release](https://img.shields.io/github/release/zouyx/agollo.svg)](https://github.com/zouyx/agollo/releases)
[![996.icu](https://img.shields.io/badge/link-996.icu-red.svg)](https://996.icu)

方便Golang接入配置中心框架 [Apollo](https://github.com/ctripcorp/apollo) 所开发的Golang版本客户端。

# Features

* 支持多 IP、AppID、namespace
* 实时同步配置
* 灰度配置
* 延迟加载（运行时）namespace
* 客户端，配置文件容灾
* 自定义日志，缓存组件
* 支持配置访问秘钥

# Usage

## 快速入门

### 导入 agollo

```
go get -u github.com/zouyx/agollo/v4@latest
```

### 启动 agollo

```
package main

import (
	"fmt"
	"github.com/zouyx/agollo/v4"
	"github.com/zouyx/agollo/v4/env/config"
)

func main() {
	c := &config.AppConfig{
		AppID:          "testApplication_yang",
		Cluster:        "dev",
		IP:             "http://106.54.227.205:8080",
		NamespaceName:  "dubbo",
		IsBackupConfig: true,
		Secret:         "6ce3ff7e96a24335a9634fe9abca6d51",
	}

	agollo.SetLogger(&DefaultLogger{})

	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
	fmt.Println("初始化Apollo配置成功")

	//Use your apollo key to test
	cache := client.GetConfigCache(c.NamespaceName)
	value, _ := client.Get("key")
	fmt.Println(value)
}
```

## 更多用法

***使用Demo*** ：[agollo_demo](https://github.com/zouyx/agollo_demo)

***其他语言*** ： [agollo-agent](https://github.com/zouyx/agollo-agent.git) 做本地agent接入，如：PHP

欢迎查阅 [Wiki](https://github.com/zouyx/agollo/wiki) 或者 [godoc](http://godoc.org/github.com/zouyx/agollo) 获取更多有用的信息

如果你觉得该工具还不错或者有问题，一定要让我知道，可以发邮件或者[留言](https://github.com/zouyx/agollo/issues)。

# User

* [使用者名单](https://github.com/zouyx/agollo/issues/20)

# Contribution

* Source Code: https://github.com/zouyx/agollo/
* Issue Tracker: https://github.com/zouyx/agollo/issues

# License

The project is licensed under the [Apache 2 license](https://github.com/zouyx/agollo/blob/master/LICENSE).

# Reference

Apollo : https://github.com/ctripcorp/apollo
