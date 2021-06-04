package apollo

import (
	"mosn.io/pkg/log"
)

// An adapter to implement log.LoggerInterface in agollo package.
type DefaultLogger struct {
	logger log.ErrorLogger
}

func NewDefaultLogger(logger log.ErrorLogger) *DefaultLogger {
	return &DefaultLogger{
		logger: logger,
	}
}
func (d *DefaultLogger) Debugf(format string, params ...interface{}) {
	d.logger.Debugf(format, params)
}

func (d *DefaultLogger) Infof(format string, params ...interface{}) {
	d.logger.Infof(format, params)
}

func (d *DefaultLogger) Warnf(format string, params ...interface{}) {
	d.logger.Warnf(format, params)
}

func (d *DefaultLogger) Errorf(format string, params ...interface{}) {
	d.logger.Errorf(format, params)
}

func (d *DefaultLogger) Debug(v ...interface{}) {
	if len(v) == 0 {
		return
	}
	format := buildForamat(v)
	d.logger.Debugf(format, v)
}

func buildForamat(v []interface{}) string {
	return "%v"
}

func (d *DefaultLogger) Info(v ...interface{}) {
	if len(v) == 0 {
		return
	}
	format := buildForamat(v)
	d.logger.Infof(format, v)
}

func (d *DefaultLogger) Warn(v ...interface{}) {
	if len(v) == 0 {
		return
	}
	format := buildForamat(v)
	d.logger.Warnf(format, v)
}

func (d *DefaultLogger) Error(v ...interface{}) {
	if len(v) == 0 {
		return
	}
	format := buildForamat(v)
	d.logger.Errorf(format, v)
}
