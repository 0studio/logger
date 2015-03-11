package logger

type Logger interface {
	Close()
	Flush()
	Debug(v ...interface{})
	Debugf(format string, params ...interface{})
	LogError(err error)
	Error(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, params ...interface{})
	Infof(format string, params ...interface{})
	Errorf(format string, params ...interface{})
}
