package eventlog

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"strconv"
	"time"
)

const (
	ldebug       = "DEBUG"
	linfo        = "INFO"
	lwarn        = "WARN"
	lerror       = "ERROR"
	LogEventName = "log-event"
)

type LogEvent struct {
	Time    time.Time      `json:"time"`
	Level   string         `json:"level"`
	Message string         `json:"message"`
	Values  map[string]any `json:"values"`
}

type EventLog struct {
	Ctx          context.Context `json:"-"`
	DebugEnabled bool            `json:"debugEnabled"`
}

func (l *EventLog) SetDebug(debug bool) {
	l.DebugEnabled = debug
}

func (l *EventLog) InfoEvent(message string, args ...any) {
	l.levelLog(linfo, message, args...)
}

func (l *EventLog) WarnEvent(message string, args ...any) {
	l.levelLog(lwarn, message, args...)
}

func (l *EventLog) ErrorEvent(message string, args ...any) {
	l.levelLog(lerror, message, args...)
}

func (l *EventLog) DebugEvent(message string, args ...any) {
	l.levelLog(ldebug, message, args...)
}

func (l *EventLog) levelLog(level, message string, args ...any) {
	if level != ldebug || (level == ldebug && l.DebugEnabled) {
		now := time.Now()
		values := l.assembleArgs(args...)
		event := LogEvent{
			Time:    now,
			Level:   level,
			Message: message,
			Values:  values,
		}
		runtime.EventsEmit(l.Ctx, LogEventName, &event)
	}
}

func (l *EventLog) assembleArgs(args ...any) map[string]any {
	if len(args) == 1 {
		if aa, ok := args[0].([]any); ok {
			return l.assembleArgs(aa...)
		}
	}

	values := map[string]any{}
	for i := 0; (i + 1) < len(args); i += 2 {
		k, v := coerceArgs(args[i], args[i+1])
		values[k] = v
	}

	return values
}

func coerceArgs(a any, b any) (string, any) {
	var (
		key string
		val = fmt.Sprintf("%v", b)
	)

	if a == nil {
		return "<nil>", val
	}

	switch k := a.(type) {
	case string:
		key = k
	case float64:
		key = strconv.FormatFloat(k, 'f', -1, 64)
	case int64:
		key = strconv.FormatInt(k, 64)
	case bool:
		key = strconv.FormatBool(k)
	default:
		key = fmt.Sprintf("%v", a)
	}

	return key, val
}

func (l *EventLog) LogEventName() string {
	return LogEventName
}
