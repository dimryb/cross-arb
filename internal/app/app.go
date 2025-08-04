package app

import (
	"context"

	i "github.com/dimryb/cross-arb/internal/interface"
)

type App struct {
	ctx    context.Context
	logger i.Logger
	store  i.TickerStore
}

func NewApp(
	ctx context.Context,
	logger i.Logger,
	store i.TickerStore,
) *App {
	return &App{
		ctx:    ctx,
		logger: logger,
		store:  store,
	}
}

func (a *App) Context() context.Context {
	return a.ctx
}

func (a *App) Logger() i.Logger {
	return a.logger
}

func (a *App) TickerStore() i.TickerStore {
	return a.store
}
