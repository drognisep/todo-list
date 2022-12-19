package eventlog

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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

func (l *EventLog) InfoEvent(message string, args ...any) {
	l.levelLog(linfo, message, args)
}

func (l *EventLog) WarnEvent(message string, args ...any) {
	l.levelLog(lwarn, message, args)
}

func (l *EventLog) ErrorEvent(message string, args ...any) {
	l.levelLog(lerror, message, args)
}

func (l *EventLog) DebugEvent(message string, args ...any) {
	l.levelLog(ldebug, message, args)
}

func (l *EventLog) levelLog(level, message string, args ...any) {
	if level != ldebug || (level == ldebug && l.DebugEnabled) {
		now := time.Now()
		values := l.assembleArgs(args)
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
	values := map[string]any{}
	for i := 0; (i + 1) < len(args); i += 2 {
		if s, ok := args[i].(string); ok {
			values[s] = args[i+1]
		} else {
			key := fmt.Sprintf("%v", args[i])
			values[key] = args[i+1]
		}
	}

	return values
}

func (l *EventLog) LogEventName() string {
	return LogEventName
}
