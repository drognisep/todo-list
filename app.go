package main

import (
	"context"
)

// App struct
type App struct {
	ctx   context.Context
	tasks *TaskController
}

// NewApp creates a new App application struct
func NewApp(tasks *TaskController) *App {
	return &App{
		tasks: tasks,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.tasks.ctx = ctx
}
