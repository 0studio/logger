package logger

import (
	"fmt"
)

func NewStdoutLogger() Logger {
	return &StdoutLogger{}
}

type StdoutLogger struct {
}

func (out *StdoutLogger) Close() {
}

func (out *StdoutLogger) Debug(v ...interface{}) {
	fmt.Println(v)
}
func (out *StdoutLogger) Debugf(format string, params ...interface{}) {
	fmt.Printf(format+"\n", params)
}

func (out *StdoutLogger) LogError(err error) {
	fmt.Println(err)
}
func (out *StdoutLogger) Error(v ...interface{}) {
	fmt.Println(v)
}
func (out *StdoutLogger) Flush() {
}

func (out *StdoutLogger) Info(v ...interface{}) {
	fmt.Println(v)
}
func (out *StdoutLogger) Warn(v ...interface{}) {
	fmt.Println(v)
}
func (out *StdoutLogger) Warnf(format string, params ...interface{}) {
	fmt.Println(format, params)
}

func (out *StdoutLogger) Infof(format string, params ...interface{}) {
	fmt.Println(format, params)
}

func (out *StdoutLogger) Errorf(format string, params ...interface{}) {
	fmt.Println(format, params)
}
