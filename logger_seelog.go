package logger

import (
	"fmt"
	"github.com/cihub/seelog"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
)

type SeeLogLogger struct {
	logger     seelog.LoggerInterface
	mode       string
	debugLevel int
}

const (
	DebugLevelDefault = 0
	DebugLevelVV      = 1
	DebugLevelVVV     = 2
)

func NewSeeLogLoggerC(debugLevel int, mode, logPath string, argStack ArgStack) Logger {
	var seeLogLogger SeeLogLogger
	seeLogLogger.debugLevel = debugLevel
	seeLogLogger.mode = mode
	byteContent, err := ioutil.ReadFile(logPath)
	if err != nil {
		fmt.Println("read seelog config errr:", err)
		return &seeLogLogger
	}
	
	content := string(byteContent)
	for key, value := range argStack {
		content = strings.Replace(content, fmt.Sprintf("%%%s", key), value, -1)
	}
	seeLogLogger.logger, err = seelog.LoggerFromConfigAsString(content)
	if err != nil {
		fmt.Println(err)
		panic("init_seelog_fail")
	}

	seelog.ReplaceLogger(seeLogLogger.logger)
	defer seelog.Flush()
	return &seeLogLogger
}

func NewSeeLogLogger(debugLevel int, mode, logPath string, platform, server, process uint64, appName string) Logger {
	argStack := NewArgStack()
	argStack.Add("appName", appName)
	argStack.Add("platformid", strconv.FormatUint(platform, 10))
	argStack.Add("serverid", strconv.FormatUint(server, 10))
	argStack.Add("processidx", strconv.FormatUint(process, 10))
	
	return NewSeeLogLoggerC(debugLevel, mode, logPath, argStack)
}

func (seelogger *SeeLogLogger) Close() {
	if seelogger.logger != nil && !seelogger.logger.Closed() {
		seelog.Flush()
		seelogger.logger.Close()
		seelogger.logger = nil
	}
}

func (seelogger *SeeLogLogger) Debug(v ...interface{}) {
	// 如果debug 的日志太多，会导致如下提示
	// Seelog queue overflow: more than 10000 messages in the queue. Flushing.

	if seelogger.IsModeDev() {
		// 虽然可以在seelog.xml 里配置 Debug的内容不打印， 但是 似乎 seelog  在处理的时候
		seelog.Info(seelogger.GetPathLine(), v)
	}
}
func (seelogger *SeeLogLogger) Debugf(format string, params ...interface{}) {
	if seelogger.IsModeDev() {
		seelog.Infof(fmt.Sprintf("%s:%s", seelogger.GetPathLine(), format), params...)
	}
}

func (seelogger *SeeLogLogger) DebugVV(v ...interface{}) {
	// 如果debug 的日志太多，会导致如下提示
	// Seelog queue overflow: more than 10000 messages in the queue. Flushing.

	if seelogger.IsModeDev() && seelogger.debugLevel >= DebugLevelVV {
		// 虽然可以在seelog.xml 里配置 Debug的内容不打印， 但是 似乎 seelog  在处理的时候
		seelog.Info(seelogger.GetPathLine(), v)
	}
}
func (seelogger *SeeLogLogger) DebugVVf(format string, params ...interface{}) {
	if seelogger.IsModeDev() && seelogger.debugLevel >= DebugLevelVV {
		seelog.Infof(fmt.Sprintf("%s:%s", seelogger.GetPathLine(), format), params...)
	}
}
func (seelogger *SeeLogLogger) DebugVVV(v ...interface{}) {
	// 如果debug 的日志太多，会导致如下提示
	// Seelog queue overflow: more than 10000 messages in the queue. Flushing.

	if seelogger.IsModeDev() && seelogger.debugLevel >= DebugLevelVVV {
		// 虽然可以在seelog.xml 里配置 Debug的内容不打印， 但是 似乎 seelog  在处理的时候
		seelog.Info(seelogger.GetPathLine(), v)
	}
}
func (seelogger *SeeLogLogger) DebugVVVf(format string, params ...interface{}) {
	if seelogger.IsModeDev() && seelogger.debugLevel >= DebugLevelVVV {
		seelog.Infof(fmt.Sprintf("%s:%s", seelogger.GetPathLine(), format), params...)
	}
}

func (seelogger *SeeLogLogger) LogError(err error) {
	seelogger.Error(err)
}
func (seelogger *SeeLogLogger) Error(v ...interface{}) {
	if !seelogger.IsModePro() {
		fmt.Println(seelogger.GetPathLine(), v)
	}
	if !seelogger.IsModeTest() {
		seelog.Error(seelogger.GetPathLine(), v)
	}
}
func (seelogger *SeeLogLogger) Flush() {
	seelog.Flush()
}

func (seelogger *SeeLogLogger) Info(v ...interface{}) {
	if seelogger.IsModeTest() {
		fmt.Println("[Info]:", v)
	} else {
		seelog.Info(v)
		if seelogger.IsModeDev() {
			seelogger.Flush()
		}
	}
}
func (seelogger *SeeLogLogger) Warn(v ...interface{}) {
	if seelogger.IsModeTest() {
		fmt.Println("[Warn]", fmt.Sprintf("%s:", seelogger.GetPathLine), v)
	} else {
		seelog.Warn(v)
		if seelogger.IsModeDev() {
			seelogger.Flush()
		}
	}
}
func (seelogger *SeeLogLogger) Warnf(format string, params ...interface{}) {
	if seelogger.IsModeTest() {
		fmt.Printf("[Warn]"+fmt.Sprintf("%s:%s", seelogger.GetPathLine(), format), params...)
	} else {
		seelog.Warnf(fmt.Sprintf("%s:%s", seelogger.GetPathLine(), format), params...)
		if seelogger.IsModeDev() {
			seelogger.Flush()
		}
	}
}

func (seelogger *SeeLogLogger) Infof(format string, params ...interface{}) {
	if seelogger.IsModeTest() {
		fmt.Printf("[Info]:"+format, params...)
	} else {
		seelog.Infof(format, params...)
		if seelogger.IsModeDev() {
			seelogger.Flush()
		}
	}
}

func (seelogger *SeeLogLogger) Errorf(format string, params ...interface{}) {
	if seelogger.IsModeTest() {
		fmt.Printf(fmt.Sprintf("[Error]:%s:%s", seelogger.GetPathLine(), format), params...)
	} else {
		seelog.Errorf(fmt.Sprintf("%s:%s", seelogger.GetPathLine(), format), params...)
		if seelogger.IsModeDev() {
			seelogger.Flush()
		}
	}
}
func (seelogger *SeeLogLogger) GetPathLine() string {
	_, path, line, ok := runtime.Caller(2)
	if ok {
		return fmt.Sprintf("%s:%d", path, line)
	}
	return ""

}
func (seelogger *SeeLogLogger) IsModeTest() bool {
	return seelogger.mode == "test"
}
func (seelogger *SeeLogLogger) IsModeDev() bool {
	return seelogger.mode == "dev"
}
func (seelogger *SeeLogLogger) IsModePro() bool {
	return seelogger.mode == "pro"
}


type ArgStack map[string]string

func NewArgStack() ArgStack {
	return make(map[string]string)
}

func (this ArgStack) Add(key, value string) {
	this[key] = value
}