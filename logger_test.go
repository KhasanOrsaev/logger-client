package logger

import (
	"errors"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

type testingData struct {
	short string
	full  string
	add   *map[string]interface{}
	err   error
	level logrus.Level
}

var (
	logger        *Logger
	exitCode      *int
	message       *string
)

var dataTesting = []testingData{
	{short: "aaaa", full: "27", add: nil, err: nil, level: 5},
	{short: "cccc ddd", full: "e ff", add: &map[string]interface{}{
		"p": interface{}("17"),
		"q": interface{}(71),
	}, err: nil, level: 2},
	{short: "cccc ddd", full: "e ff", add: &map[string]interface{}{
		"p": interface{}("17"),
		"q": interface{}(71),
	}, err: errors.New("bad command"), level: 1},
	{short: "cccc ddd", full: "e ff", add: &map[string]interface{}{
		"p": interface{}("17"),
		"q": interface{}(71),
	}, err: errors.New("bad command"), level: 2},
}

func (_ *Logger) WriteMessage(msg string) error {
	message = &msg
	return nil
}

func (_ *Logger) Write(data []byte) (int, error) {
	return 0, nil
}

func (_ *Logger) Close() error {
	return nil
}

func init() {
	doExit = func(code int) { exitCode = &code }
	time.Sleep(time.Second / 2)
}

func TestNewLoggerDefault(t *testing.T) {
	if logger, err := NewLoggerDefault(map[string]interface{}{
		"module": "test",
		"format": "[%s] %s.%s message: %s context: %s extra: %s",
		"level": 3,
	}); err != nil {
		t.Fatal(err)
	} else if logger == nil {
		t.Fatal("no default logger created")
	}
}

func TestLogDefault(t *testing.T) {
	for _, v := range dataTesting {
		exitCode = nil
		message = nil
		Log(v.level, v.short, v.full, v.add, v.err)
	}
}

func TestDebugDefault(t *testing.T) {
	for _, v := range dataTesting {
		exitCode = nil
		message = nil
		Debug(v.short, v.full, v.add, v.err)
		//v.level = gelf.LOG_DEBUG
		//analyzeResult(t, &loggerDefault, v, message, 0, timeStart)
	}
}

func TestInfoDefault(t *testing.T) {
	for _, v := range dataTesting {
		exitCode = nil
		message = nil
		NewLoggerDefault(map[string]interface{}{"level":5, "format": "json"})
		Info(v.short, v.full, v.add, v.err)
		//v.level = gelf.LOG_INFO
		//analyzeResult(t, &loggerDefault, v, message, 0, timeStart)
	}
}

func TestWarningDefault(t *testing.T) {
	for _, v := range dataTesting {
		exitCode = nil
		message = nil
		Warning(v.short, v.full, v.add, v.err)
		//v.level = gelf.LOG_WARNING
		//analyzeResult(t, &loggerDefault, v, message, 0, timeStart)
	}
}

func TestErrorDefault(t *testing.T) {
	for _, v := range dataTesting {
		exitCode = nil
		message = nil
		Error(v.short, v.full, v.add, v.err)
		//v.level = gelf.LOG_ERR
		//analyzeResult(t, &loggerDefault, v, message, 0, timeStart)
	}
}

func TestAlertDefault(t *testing.T) {
	for _, v := range dataTesting {
		exitCode = nil
		message = nil
		Alert(v.short, v.full, v.add, v.err)
		//v.level = gelf.LOG_ALERT
		//analyzeResult(t, &loggerDefault, v, message, 2, timeStart)
	}
}

func TestEmergencyDefault(t *testing.T) {
	for _, v := range dataTesting {
		exitCode = nil
		message = nil
		Emergency(v.short, v.full, v.add, v.err)
		//v.level = gelf.LOG_EMERG
		//analyzeResult(t, &loggerDefault, v, message, 3, timeStart)
	}
}

func TestNewLogger(t *testing.T) {
	var err error
	logger, err = NewLogger(map[string]interface{}{
		"module": "new_test",
		"action_type": "test_typr",
		"format": "",
		"level": 2,
	})
	if err != nil {
		t.Fatal(err)
	} else if logger == nil {
		t.Fatal("no logger created")
	}
}

func TestLog(t *testing.T) {
	for _, v := range dataTesting {
		exitCode = nil
		message = nil
		logger.Log(v.level, v.short, v.full, v.add, v.err)
		//analyzeResult(t, logger, v, message, 0, timeStart)
	}
}

func TestDebug(t *testing.T) {
	for _, v := range dataTesting {
		exitCode = nil
		message = nil
		logger.Debug(v.short, v.full, v.add, v.err)
		//v.level = gelf.LOG_DEBUG
		//analyzeResult(t, logger, v, message, 0, timeStart)
	}
}

func TestInfo(t *testing.T) {
	for _, v := range dataTesting {
		exitCode = nil
		message = nil
		logger.Info(v.short, v.full, v.add, v.err)
		//v.level = gelf.LOG_INFO
		//analyzeResult(t, logger, v, message, 0, timeStart)
	}
}

func TestWarning(t *testing.T) {
	for _, v := range dataTesting {
		exitCode = nil
		message = nil
		logger.Warning(v.short, v.full, v.add, v.err)
		//v.level = gelf.LOG_WARNING
		//analyzeResult(t, logger, v, message, 0, timeStart)
	}
}

func TestError(t *testing.T) {
	for _, v := range dataTesting {
		exitCode = nil
		message = nil
		logger.Error(v.short, v.full, v.add, v.err)
		//v.level = gelf.LOG_ERR
		//analyzeResult(t, logger, v, message, 0, timeStart)
	}
}

func TestAlert(t *testing.T) {
	for _, v := range dataTesting {
		exitCode = nil
		message = nil
		logger.Alert(v.short, v.full, v.add, v.err)
		//v.level = gelf.LOG_ALERT
		//analyzeResult(t, logger, v, message, 2, timeStart)
	}
}

func TestEmergency(t *testing.T) {
	for _, v := range dataTesting {
		exitCode = nil
		message = nil
		logger.Emergency(v.short, v.full, v.add, v.err)
		//v.level = gelf.LOG_EMERG
		//analyzeResult(t, logger, v, message, 3, timeStart)
	}
}
