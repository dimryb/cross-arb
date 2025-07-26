package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	spotlist "github.com/dimryb/cross-arb/internal/api/mexc/spot"
	"github.com/dimryb/cross-arb/internal/config"
	i "github.com/dimryb/cross-arb/internal/interface"
)

type Arbitrage struct {
	ctx context.Context
	app i.Application
	log i.Logger
	cfg *config.CrossArbConfig
}

func NewArbitrageService(
	ctx context.Context,
	app i.Application,
	logger i.Logger,
	cfg *config.CrossArbConfig,
) *Arbitrage {
	return &Arbitrage{
		ctx: ctx,
		app: app,
		log: logger,
		cfg: cfg,
	}
}

func (m *Arbitrage) Run() error {
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Second)
		params := `{"symbol":"SOLUSDT"}`

		for {
			select {
			case <-m.ctx.Done():
				return
			case <-ticker.C:
				BookTicker := spotlist.BookTicker(params)
				fmt.Println("Получили:", BookTicker)
			}
		}
	}()

	m.log.Infof("Arbitrage service is running...")

	wg.Wait()

	return nil
}
