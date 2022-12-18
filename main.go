package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	tasks, err := NewTaskController()
	// Create an instance of the app structure
	app := NewApp(tasks)
	if err != nil {
		fmt.Printf("Failed to start task controller: %v\n", err)
		os.Exit(1)
	}

	// Create application with options
	err = wails.Run(&options.App{
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
			tasks,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
