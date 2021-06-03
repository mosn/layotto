package apollo

import (
	"mosn.io/pkg/log"
	"testing"
)

func TestNewDefaultLogger(t *testing.T) {
	// logger
	logger := NewDefaultLogger(log.DefaultLogger)
	logger.Debugf("test Debugf")
	logger.Infof("test Infof")
	logger.Warnf("test Warnf")
	logger.Errorf("test Errorf")
	logger.Debug("test Debug")
	logger.Info("test Info")
	logger.Warn("test Warn")
	logger.Error("test Debug")
	logger.Debug()
	logger.Info()
	logger.Warn()
	logger.Error()
}
