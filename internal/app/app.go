package app

import (
	i "github.com/dimryb/cross-arb/internal/interface"
)

type App struct {
	Logger i.Logger
}

func NewApp(logger i.Logger) *App {
	return &App{Logger: logger}
}
