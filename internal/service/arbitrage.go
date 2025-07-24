package service

import (
	"context"
	"sync"
	"time"

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

		for {
			select {
			case <-m.ctx.Done():
				return

			default:
				// Здесь выполняется код в бесконечном цикле сервиса в отдельной горутине
				time.Sleep(1 * time.Second) // Чтобы не грузить за зря процессор, убрать!
			}
		}
	}()

	m.log.Infof("Arbitrage service is running...")

	wg.Wait()

	return nil
}
