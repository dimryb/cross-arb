package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/dimryb/cross-arb/internal/adapter/jupiter"
	"github.com/dimryb/cross-arb/internal/adapter/mexc"
	"github.com/dimryb/cross-arb/internal/config"
	"github.com/dimryb/cross-arb/internal/controller/grpc"
	"github.com/dimryb/cross-arb/internal/controller/http"
	"github.com/dimryb/cross-arb/internal/entity"
	i "github.com/dimryb/cross-arb/internal/interface"
	"github.com/dimryb/cross-arb/internal/logger"
	"github.com/dimryb/cross-arb/internal/report"
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

func NewApp(ctx context.Context, cancel context.CancelFunc, cfg *config.CrossArbConfig) *App {
	return &App{
		ctx:    ctx,
		cancel: cancel,
		cfg:    cfg,
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
	a.log = logger.New(a.cfg.Log.Level)
	a.store = storage.NewTickerStore()
	reportSvc := report.NewReportService(a.log, a.store)

	// --- Adapters ---
	mexcAdapter := mexc.NewAdapter(a.log, 3*time.Second)

	// Jupiter: собираем конфиг напрямую (заменяет фабрику)
	jupCfg, ok := a.cfg.Exchanges[config.JupExchange]
	if !ok {
		a.log.Fatalf("exchange %s not found in configuration", config.JupExchange)
	}
	if !jupCfg.Enabled {
		a.log.Fatalf("exchange %s is disabled", config.JupExchange)
	}

	pairMap := make(map[string]jupiter.MintPair, len(jupCfg.Pairs))
	for symbol, p := range jupCfg.Pairs {
		if p.Base == "" || p.Quote == "" {
			a.log.Fatalf("missing mint address for Jupiter pair %q", symbol)
		}
		pairMap[symbol] = jupiter.MintPair{
			BaseMint:  p.Base,
			QuoteMint: p.Quote,
		}
	}

	jupiterAdapter := jupiter.NewAdapter(a.log, &jupiter.AdapterConfig{
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

	adapters := []i.EXAdapter{mexcAdapter, jupiterAdapter}

	pricesCh := make(chan entity.ExecutableQuote, a.cfg.Scanner.Buffers.Prices)
	orderBooksCh := make(chan entity.OrderBookResult, a.cfg.Scanner.Buffers.OrderBooks)
	oppCh := make(chan entity.ArbOpportunity, a.cfg.Scanner.Buffers.Opportunities)

	interval, err := time.ParseDuration(a.cfg.Scanner.Interval)
	if err != nil {
		a.log.Error("invalid scanner interval",
			slog.String("value", a.cfg.Scanner.Interval),
			slog.Any("err", err))
		a.cancel()
		return
	}
	_, err = scanner.NewService(
		a.log,
		interval,
		1.0, // placeholder: объём сделки в BASE для квотирования DEX
		a.cfg.Scanner.Pairs,
		adapters,
		pricesCh,
		orderBooksCh,
		oppCh,
		nil, // используем DEXPriceUseCase по умолчанию
		nil, // используем CEXOrderBookUseCase по умолчанию
		nil, // используем OpportunityUseCase по умолчанию
	)
	if err != nil {
		a.log.Error("Failed to create scanner service", slog.Any("err", err))
		a.cancel()
		return
	}

	// Заглушки консьюмеры
	go func() {
		for pp := range pricesCh {
			a.log.Info("price point",
				slog.String("pair", pp.Pair),
				slog.String("exchange", pp.Exchange),
				slog.Float64("bid", pp.Bid),
				slog.Float64("ask", pp.Ask),
				slog.Time("ts", pp.Timestamp),
			)
		}
	}()

	go func() {
		for opp := range oppCh {
			a.log.Info("opportunity",
				slog.String("pair", opp.Pair),
				slog.String("buy_on", opp.BuyOn),
				slog.Float64("buy_price", opp.BuyPrice),
				slog.String("sell_on", opp.SellOn),
				slog.Float64("sell_price", opp.SellPrice),
				slog.Float64("net", opp.NetPnl),
				slog.Float64("spread_pct", opp.SpreadPct),
				slog.Time("ts", opp.DetectedAt),
			)
		}
	}()

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

	reportSvc.Start()

	a.log.Info("App started")

	<-a.ctx.Done()

	a.log.Info("Shutting down...")

	// Останавливаем фоновые сервисы
	reportSvc.Stop()

	a.log.Info("App stopped gracefully")
}
