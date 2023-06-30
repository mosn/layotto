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

package nacos

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"mosn.io/pkg/log"
)

func TestNewDefaultLogger(t *testing.T) {
	mosnLogger, err := log.GetOrCreateLogger("stdout", nil)
	assert.Nil(t, err)
	errorLog := &log.SimpleErrorLog{
		Logger: mosnLogger,
		Level:  log.DEBUG,
	}

	logger := NewDefaultLogger(errorLog)
	logger.Debugf("test Debugf %d", 100)
	logger.Debugf("test Debugf", 100)
	logger.Infof("test Infof")
	logger.Warnf("test Warnf")
	logger.Errorf("test Errorf")
	logger.Debug("test Debug")
	logger.Info("test Info")
	logger.Warn("test Warn")
	logger.Error("test Error")
	logger.Debug()
	logger.Info()
	logger.Warn()
	logger.Error()
}
