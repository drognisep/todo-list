package main

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"strings"
	"time"
	"todo-list/eventlog"
)

// App struct
type App struct {
	ctx    context.Context
	tasks  *TaskController
	logger *eventlog.EventLog
}

// NewApp creates a new App application struct
func NewApp(logger *eventlog.EventLog, tasks *TaskController) *App {
	return &App{
		tasks:  tasks,
		logger: logger,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if err := a.tasks.onStartup(ctx); err != nil {
		// TODO: Add some more helpful error reporting mechanism.
		log.Fatalln("Failed to start task controller:", err)
	}
	a.logger.Ctx = ctx

	valueFormatter := func(values map[string]any) string {
		var buf strings.Builder
		i := 0
		for k, v := range values {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(fmt.Sprintf("%s=%v", k, v))
			i++
		}
		return buf.String()
	}

	runtime.EventsOn(ctx, eventlog.LogEventName, func(ievent ...interface{}) {
		if len(ievent) != 1 {
			a.logger.WarnEvent("Received log event with unexpected data count", "len", len(ievent))
			return
		}
		if event, ok := ievent[0].(*eventlog.LogEvent); ok {
			valStr := valueFormatter(event.Values)
			if len(valStr) > 0 {
				valStr = ": " + valStr
			}
			fmt.Printf("%s [%s] %s%s\n", event.Time.Format(time.RFC3339), event.Level, event.Message, valStr)
		} else {
			a.logger.WarnEvent("Received non-LogEvent from logger", "type", fmt.Sprintf("%T", ievent[0]))
		}
	})
}
