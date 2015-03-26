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
	logger seelog.LoggerInterface
	mode   string
}

func NewSeeLogLogger(mode, logPath string, platform, server, process uint64, appName string) Logger {
	var seeLogLogger SeeLogLogger
	seeLogLogger.mode = mode
	byteContent, err := ioutil.ReadFile(logPath)
	if err != nil {
		fmt.Println("read seelog config errr:", err)
		return &seeLogLogger
	}
	content := strings.Replace(string(byteContent), "%appName", appName, -1)
	content = strings.Replace(content, "%platformid", strconv.FormatUint(platform, 10), -1)
	content = strings.Replace(content, "%serverid", strconv.FormatUint(server, 10), -1)
	content = strings.Replace(content, "%processidx", strconv.FormatUint(process, 10), -1)

	seeLogLogger.logger, err = seelog.LoggerFromConfigAsString(content)
	// logger, err := seelog.LoggerFromConfigAsFile(logPath)
	if err != nil {
		fmt.Println(err)
		panic("init_seelog_fail")
		// os.Exit(defs.EXIT_CODE_SEELOG_INIT_FAIL)
	}

	seelog.ReplaceLogger(seeLogLogger.logger)
	defer seelog.Flush()
	return &seeLogLogger
}

// func Init() {
// 	logPath := fmt.Sprintf("%sseelog_%s.xml", conf.GetConfigDir(), conf.GetMode())
// }
// func ReloadLogger() {
// 	Close()
// 	Init()
// }
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
		seelog.Debug(seelogger.GetPathLine(), v)
	}
}
func (seelogger *SeeLogLogger) Debugf(format string, params ...interface{}) {
	if seelogger.IsModeDev() {
		seelog.Debugf(fmt.Sprintf("%s:%s", seelogger.GetPathLine(), format), params...)
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
		fmt.Println("[Warn]:", v)
	} else {
		seelog.Warn(v)
		if seelogger.IsModeDev() {
			seelogger.Flush()
		}
	}
}
func (seelogger *SeeLogLogger) Warnf(format string, params ...interface{}) {
	if seelogger.IsModeTest() {
		fmt.Printf("[Warn]:"+format, params...)
	} else {
		seelog.Warnf(format, params...)
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
