// Логгер
// обертка логера логрус
package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

const (
	FormatDefault = "[%s] %s.%s message: %s context: %s extra: %s"
	ModuleDefault = "default_module"
	OutputDefault = "stdout"
)

type (

	// Экспортируемый из модуля тип логгера
	Logger struct {
		Client                *logrus.Logger
		eventAttributesCommon map[string]interface{}
	}
)

var (
	// Логгер по умолчанию
	loggerDefault Logger
	// Функция аварийного завершения работы
	doExit = func(code int) { os.Exit(code) }
)

func init() {
	_,err := NewLoggerDefault(map[string]interface{}{})
	if err != nil {
		logrus.Fatal(err.Error())
	}
}
// initFileDirectory инициализацирует атрибуты логгера по умолчаню, не зависящие от дальнейших настроек логгера
func initFileDirectory() {
	// если директории нет то создаем
	if _, err := os.Stat("./var/log"); os.IsNotExist(err) {
		err = os.MkdirAll("./var/log", os.ModePerm)
		if err != nil {
			logrus.Fatal(err.Error())
		}
	}
}

// getResource получает новый экземпляр отправителя (транспорта) сообщений
func getResource(output, module, format string, level logrus.Level) (resource *logrus.Logger, err error) {
	var outputResource io.Writer
	switch output {
	case "file":
		initFileDirectory()
		outputResource, err = os.OpenFile("./var/log/" + module  + "_" + time.Now().Format("2006-01-02") + ".log",
			os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
		if err != nil {
			return
		}
	case OutputDefault:
		outputResource = os.Stdout
	}
	resource = &logrus.Logger{
		Out:          outputResource,
		Hooks:        nil,
		Formatter:    &logFormatter{format},
		ReportCaller: false,
		Level:        level,
		ExitFunc:     doExit,
	}
	return
}
/*
NewLogger создаёт новый экземпляр логгера, если не задан модуль, формат и уровень, то используются
дефолтные значения
*/
func NewLogger(
	add map[string]interface{}, // Допонительные настройки логгера
) (*Logger, error) {
	var err error
    var module, format,output string
	var level logrus.Level
	if m,ok := add["module"];!ok {
		module = ModuleDefault
	} else {
		module = m.(string)
	}
	if m,ok := add["format"];!ok {
		format = FormatDefault
	} else {
		format = m.(string)
	}
	if m,ok := add["level"];!ok {
		level = logrus.WarnLevel
	} else {
		level = logrus.Level(m.(int))
	}
	if m,ok := add["output"];ok {
		output = m.(string)
	} else {
		output = OutputDefault
	}
	client, err := getResource(output, module, format, level)
	if err != nil {
		return nil, err
	}

	return &Logger{Client: client, eventAttributesCommon: map[string]interface{}{
		"_pid":         os.Getpid(),
		"module":      module,
	}}, nil
}

// NewLoggerDefault инициализирует логгер по умолчанию
// Если логгер не инициализирован, используется стандартный Go-логгер по умолчанию
func NewLoggerDefault(
	add map[string]interface{}, // Допонительные настройки логгера
) (*Logger, error) {
	var err error
	var module, format, output string
	var level logrus.Level
	if m,ok := add["module"];!ok {
		module = ModuleDefault
	} else {
		module = m.(string)
	}
	if m,ok := add["format"];!ok {
		format = FormatDefault
	} else {
		format = m.(string)
	}
	if m,ok := add["level"];!ok {
		level = logrus.WarnLevel
	} else {
		level = logrus.Level(m.(int))
	}
	if m,ok := add["output"];ok {
		output = m.(string)
	} else {
		output = OutputDefault
	}
	loggerDefault.Client, err = getResource(output, module, format, level)
	loggerDefault.eventAttributesCommon = map[string]interface{}{
		"module": module,
		"_pid":         os.Getpid(),
	}
	if err != nil {
		return nil, err
	}

	return &loggerDefault, nil
}

// mergeMaps производит объединение произвольных ассоциированных масивов со строковыми ключами
func mergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// Log поизводит запись события с произвольным уровнем
func (logger *Logger) Log(
	level logrus.Level, // Уровень события
	eventShortMessage string, // Краткая информация о событии
	eventFullMessage string, // Полная информация о событии
	eventAttributes *map[string]interface{}, // Дополнительные атрибуты события
	eventError error, // Логируемая ошибка. Если есть, замещает собой eventFullMessage
) {
	fields := logrus.Fields{}
	if eventError != nil {
		eventFullMessage = eventError.Error()
	}
	fields["host"],_ = os.Hostname()
	fields["full_message"] = eventFullMessage
	var attr map[string]interface{}
	if eventAttributes != nil {
		attr = mergeMaps(logger.eventAttributesCommon, *eventAttributes)
	} else {
		attr = logger.eventAttributesCommon
	}
	for i,v := range attr {
		fields[i] = v
	}
	l := logger.Client.WithFields(fields)
	l.Log(level, eventShortMessage)
}

// Log поизводит запись события с произвольным уровнем в логгер по умолчанию
func Log(
	level logrus.Level, // Уровень события
	eventShortMessage string, // Краткая информация о событии
	eventFullMessage string, // Полная информация о событии
	eventAttributes *map[string]interface{}, // Дополнительные атрибуты события
	eventError error, // Ошибка события. Если есть, замещает собой eventFullMessage
) {
	loggerDefault.Log(level, eventShortMessage, eventFullMessage, eventAttributes, eventError)
}

// Debug поизводит запись события с уровнем "отладка"
func (logger *Logger) Debug(
	eventShortMessage string, // Краткая информация о событии
	eventFullMessage string, // Полная информация о событии
	eventAttributes *map[string]interface{}, // Дополнительные атрибуты события
	eventError error, // Логируемая ошибка. Если есть, замещает собой eventFullMessage
) {
	logger.Log(logrus.DebugLevel, eventShortMessage, eventFullMessage, eventAttributes, eventError)
}

// Debug поизводит запись события с уровнем "отладка" в логгер по умолчанию
func Debug(
	eventShortMessage string, // Краткая информация о событии
	eventFullMessage string, // Полная информация о событии
	eventAttributes *map[string]interface{}, // Дополнительные атрибуты события
	eventError error, // Логируемая ошибка. Если есть, замещает собой eventFullMessage
) {
	loggerDefault.Debug(eventShortMessage, eventFullMessage, eventAttributes, eventError)
}

// Info поизводит запись события с уровнем "информирование"
func (logger *Logger) Info(
	eventShortMessage string, // Краткая информация о событии
	eventFullMessage string, // Полная информация о событии
	eventAttributes *map[string]interface{}, // Дополнительные атрибуты события
	eventError error, // Логируемая ошибка. Если есть, замещает собой eventFullMessage
) {
	logger.Log(logrus.InfoLevel, eventShortMessage, eventFullMessage, eventAttributes, eventError)
}

// Info поизводит запись события с уровнем "информирование" в логгер по умолчанию
func Info(
	eventShortMessage string, // Краткая информация о событии
	eventFullMessage string, // Полная информация о событии
	eventAttributes *map[string]interface{}, // Дополнительные атрибуты события
	eventError error, // Логируемая ошибка. Если есть, замещает собой eventFullMessage
) {
	loggerDefault.Info(eventShortMessage, eventFullMessage, eventAttributes, eventError)
}

// Warning поизводит запись события с уровнем "предупреждение"
func (logger *Logger) Warning(
	eventShortMessage string, // Краткая информация о событии
	eventFullMessage string, // Полная информация о событии
	eventAttributes *map[string]interface{}, // Дополнительные атрибуты события
	eventError error, // Логируемая ошибка. Если есть, замещает собой eventFullMessage
) {
	logger.Log(logrus.WarnLevel, eventShortMessage, eventFullMessage, eventAttributes, eventError)
}

// Warning поизводит запись события с уровнем "предупреждение" в логгер по умолчанию
func Warning(
	eventShortMessage string, // Краткая информация о событии
	eventFullMessage string, // Полная информация о событии
	eventAttributes *map[string]interface{}, // Дополнительные атрибуты события
	eventError error, // Логируемая ошибка. Если есть, замещает собой eventFullMessage
) {
	loggerDefault.Warning(eventShortMessage, eventFullMessage, eventAttributes, eventError)
}

// Error поизводит запись события с уровнем "ошибка"
func (logger *Logger) Error(
	eventShortMessage string, // Краткая информация о событии
	eventFullMessage string, // Полная информация о событии
	eventAttributes *map[string]interface{}, // Дополнительные атрибуты события
	eventError error, // Логируемая ошибка. Если есть, замещает собой eventFullMessage
) {
	logger.Log(logrus.ErrorLevel, eventShortMessage, eventFullMessage, eventAttributes, eventError)
}

// Error поизводит запись события с уровнем "ошибка" в логгер по умолчанию
func Error(
	eventShortMessage string, // Краткая информация о событии
	eventFullMessage string, // Полная информация о событии
	eventAttributes *map[string]interface{}, // Дополнительные атрибуты события
	eventError error, // Логируемая ошибка. Если есть, замещает собой eventFullMessage
) {
	loggerDefault.Error(eventShortMessage, eventFullMessage, eventAttributes, eventError)
}

// Fatal поизводит запись события с уровнем "критическая ошибка" и завершает работу с кодом 1
func (logger *Logger) Fatal(
	eventShortMessage string, // Краткая информация о событии
	eventFullMessage string, // Полная информация о событии
	eventAttributes *map[string]interface{}, // Дополнительные атрибуты события
	eventError error, // Логируемая ошибка. Если есть, замещает собой eventFullMessage
) {
	logger.Log(logrus.FatalLevel, eventShortMessage, eventFullMessage, eventAttributes, eventError)
	doExit(1)
}

// Fatal поизводит запись события с уровнем "критическая ошибка" в логгер по умолчанию и завершает работу с кодом 1
func Fatal(
	eventShortMessage string, // Краткая информация о событии
	eventFullMessage string, // Полная информация о событии
	eventAttributes *map[string]interface{}, // Дополнительные атрибуты события
	eventError error, // Логируемая ошибка. Если есть, замещает собой eventFullMessage
) {
	loggerDefault.Fatal(eventShortMessage, eventFullMessage, eventAttributes, eventError)
}

// Alert поизводит запись события с уровнем "Тревога" и завершает работу с кодом 2
func (logger *Logger) Alert(
	eventShortMessage string, // Краткая информация о событии
	eventFullMessage string, // Полная информация о событии
	eventAttributes *map[string]interface{}, // Дополнительные атрибуты события
	eventError error, // Логируемая ошибка. Если есть, замещает собой eventFullMessage
) {
	logger.Log(logrus.FatalLevel, eventShortMessage, eventFullMessage, eventAttributes, eventError)
	doExit(2)
}

// Alert поизводит запись события с уровнем "Тревога" в логгер по умолчанию и завершает работу с кодом 2
func Alert(
	eventShortMessage string, // Краткая информация о событии
	eventFullMessage string, // Полная информация о событии
	eventAttributes *map[string]interface{}, // Дополнительные атрибуты события
	eventError error, // Логируемая ошибка. Если есть, замещает собой eventFullMessage
) {
	loggerDefault.Alert(eventShortMessage, eventFullMessage, eventAttributes, eventError)
}

// Emergency поизводит запись события с уровнем "Катастрофа" и завершает работу с кодом 3
func (logger *Logger) Emergency(
	eventShortMessage string, // Краткая информация о событии
	eventFullMessage string, // Полная информация о событии
	eventAttributes *map[string]interface{}, // Дополнительные атрибуты события
	eventError error, // Логируемая ошибка. Если есть, замещает собой eventFullMessage
) {
	logger.Log(logrus.FatalLevel, eventShortMessage, eventFullMessage, eventAttributes, eventError)
	doExit(3)
}

// Emergency поизводит запись события с уровнем "Катастрофа" в логгер по умолчанию и завершает работу с кодом 3
func Emergency(
	eventShortMessage string, // Краткая информация о событии
	eventFullMessage string, // Полная информация о событии
	eventAttributes *map[string]interface{}, // Дополнительные атрибуты события
	eventError error, // Логируемая ошибка. Если есть, замещает собой eventFullMessage
) {
	loggerDefault.Emergency(eventShortMessage, eventFullMessage, eventAttributes, eventError)
}
