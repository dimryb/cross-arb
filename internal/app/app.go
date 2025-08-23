package app

import (
	"context"
	"errors"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/dimryb/cross-arb/internal/adapter"
	"github.com/dimryb/cross-arb/internal/config"
	i "github.com/dimryb/cross-arb/internal/interface"
	"github.com/dimryb/cross-arb/internal/logger"
	"github.com/dimryb/cross-arb/internal/report"
	"github.com/dimryb/cross-arb/internal/server/grpc"
	"github.com/dimryb/cross-arb/internal/server/http"
	"github.com/dimryb/cross-arb/internal/service"
	"github.com/dimryb/cross-arb/internal/service/scanner"
	"github.com/dimryb/cross-arb/internal/storage"
)

type App struct {
	ctx    context.Context
	cancel context.CancelFunc
	log    i.Logger
	store  i.TickerStore

	cfg *config.CrossArbConfig
}

func NewApp(cfg *config.CrossArbConfig) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Context() context.Context {
	return a.ctx
}

func (a *App) Logger() i.Logger {
	return a.log
}

func (a *App) TickerStore() i.TickerStore {
	return a.store
}

func (a *App) Run() {
	a.ctx, a.cancel = signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer a.cancel()

	a.log = logger.New(a.cfg.Log.Level)
	a.store = storage.NewTickerStore()
	arbitrageService := service.NewArbitrageService(a.ctx, a, a.log, a.cfg, a.store)
	reportSvc := report.NewReportService(a.log, a.store)

	// --- Adapters ---
	mexcAdapter := adapter.NewMexcAdapter(a.log, 3*time.Second)

	// Jupiter: собираем конфиг напрямую (заменяет фабрику)
	jupCfg, ok := a.cfg.Exchanges[config.JupExchange]
	if !ok {
		a.log.Fatalf("exchange %s not found in configuration", config.JupExchange)
	}
	if !jupCfg.Enabled {
		a.log.Fatalf("exchange %s is disabled", config.JupExchange)
	}

	pairMap := make(map[string]adapter.MintPair, len(jupCfg.Pairs))
	for symbol, p := range jupCfg.Pairs {
		if p.Base == "" || p.Quote == "" {
			a.log.Fatalf("missing mint address for Jupiter pair %q", symbol)
		}
		pairMap[symbol] = adapter.MintPair{
			BaseMint:  p.Base,
			QuoteMint: p.Quote,
		}
	}

	jupiterAdapter := adapter.NewJupiterAdapter(a.log, &adapter.JupiterAdapterConfig{
		BaseURL: jupCfg.BaseURL,
		Enabled: true,
		Timeout: jupCfg.Timeout,
		Pairs:   pairMap,
	})

	defer func() {
		err := mexcAdapter.Close()
		if err != nil {
			a.log.Fatalf("failed to close mexc adapter: %v", err)
		}
		err = jupiterAdapter.Close()
		if err != nil {
			a.log.Fatalf("failed to close jupiter adapter: %v", err)
		}
	}()

	adapters := []i.ExchangeAdapter{mexcAdapter, jupiterAdapter}
	scan, err := scanner.NewScannerFromConfig(a.log, a.cfg.Scanner, adapters)
	if err != nil {
		a.log.Error("Failed to create scan", slog.Any("err", err))
		a.cancel()
		return
	}

	// Подписываемся и логируем возможности
	if a.cfg.Scanner.LogOpportunities {
		scanner.SubscribeAndHandle(
			scan,
			a.cfg.Scanner.Pairs,
			scanner.LogOpportunities(a.log),
		)
	}

	// TODO: вынести порты в конфиг
	grpcServer := grpc.NewServer(a, grpc.ServerConfig{Port: "9090"}, a.log)

	go func() {
		httpServer := http.NewHTTPServer(a.store)
		if err := httpServer.Run(":8080"); err != nil {
			a.log.Errorf("HTTP server error: %v", err)
			a.cancel()
		}
	}()

	go func() {
		if err := grpcServer.Run(); err != nil {
			a.log.Errorf("gRPC server failed: %v", err)
			a.cancel()
		}
	}()

	// Запускаем сканер
	go func() {
		if err := scan.Run(a.ctx); err != nil && !errors.Is(err, context.Canceled) {
			a.log.Errorf("Scanner stopped with error: %v", err)
			a.cancel()
		}
	}()

	reportSvc.Start()

	a.log.Info("Starting app...")
	if err = arbitrageService.Run(); err != nil {
		a.log.Errorf("Arbitrage service stopped with error: %v", err)
		a.cancel()
	} else {
		a.log.Infof("Arbitrage service stopped gracefully")
	}
}
