/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package log

//Logger logger 对象
var Logger LoggerInterface

func init() {
	Logger = &DefaultLogger{}
}

//InitLogger 初始化logger对象
func InitLogger(ILogger LoggerInterface) {
	Logger = ILogger
}

//LoggerInterface 日志接口
type LoggerInterface interface {
	Debugf(format string, params ...interface{})

	Infof(format string, params ...interface{})

	Warnf(format string, params ...interface{})

	Errorf(format string, params ...interface{})

	Debug(v ...interface{})

	Info(v ...interface{})

	Warn(v ...interface{})

	Error(v ...interface{})
}

//Debugf debug 格式化
func Debugf(format string, params ...interface{}) {
	Logger.Debugf(format, params)
}

//Infof 打印info
func Infof(format string, params ...interface{}) {
	Logger.Infof(format, params)
}

//Warnf warn格式化
func Warnf(format string, params ...interface{}) {
	Logger.Warnf(format, params)
}

//Errorf error格式化
func Errorf(format string, params ...interface{}) {
	Logger.Errorf(format, params)
}

//Debug 打印debug
func Debug(v ...interface{}) {
	Logger.Debug(v)
}

//Info 打印Info
func Info(v ...interface{}) {
	Logger.Info(v)
}

//Warn 打印Warn
func Warn(v ...interface{}) {
	Logger.Warn(v)
}

//Error 打印Error
func Error(v ...interface{}) {
	Logger.Error(v)
}

//DefaultLogger 默认日志实现
type DefaultLogger struct {
}

//Debugf debug 格式化
func (d *DefaultLogger) Debugf(format string, params ...interface{}) {

}

//Infof 打印info
func (d *DefaultLogger) Infof(format string, params ...interface{}) {

}

//Warnf warn格式化
func (d *DefaultLogger) Warnf(format string, params ...interface{}) {
}

//Errorf error格式化
func (d *DefaultLogger) Errorf(format string, params ...interface{}) {
}

//Debug 打印debug
func (d *DefaultLogger) Debug(v ...interface{}) {

}

//Info 打印Info
func (d *DefaultLogger) Info(v ...interface{}) {

}

//Warn 打印Warn
func (d *DefaultLogger) Warn(v ...interface{}) {
}

//Error 打印Error
func (d *DefaultLogger) Error(v ...interface{}) {
}
