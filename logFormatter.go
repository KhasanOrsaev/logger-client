package logger

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

type logFormatter struct {
	Form string
}

// форматирование логов
func (f *logFormatter) Format(entry *log.Entry) ([]byte, error) {
	dt := time.Now().Local().Format("2006-01-02 15:04:05")
	extra, err := json.Marshal(entry.Data)
	if err != nil {
		return nil, fmt.Errorf(f.Form, err)
	}
	l := fmt.Sprintf(f.Form, dt, entry.Data["module"], entry.Level.String(), entry.Message, "[]", string(extra))
	return append([]byte(l), '\n'), nil
}
