package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"time"
	"todo-list/eventlog"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	logger := &eventlog.EventLog{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	loadedCh := make(chan error)
	var tasks *TaskController

	go func() {
		var err error
		tasks, err = NewTaskController(logger)
		if err != nil {
			loadedCh <- err
		}
		close(loadedCh)
	}()

	select {
	case err := <-loadedCh:
		if err != nil {
			time.Sleep(200 * time.Millisecond)
			os.Exit(1)
		}
	case <-ctx.Done():
		fmt.Println("Timed out, the database may be open elsewhere")
		time.Sleep(200 * time.Millisecond)
		os.Exit(0)
	}
	cancel()
	app := NewApp(logger, tasks)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "ToDo List",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 34, G: 34, B: 34, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			logger,
			tasks,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
