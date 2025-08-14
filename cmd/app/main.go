package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/dimryb/cross-arb/internal/adapter"
	"github.com/dimryb/cross-arb/internal/app"
	"github.com/dimryb/cross-arb/internal/config"
	i "github.com/dimryb/cross-arb/internal/interface"
	"github.com/dimryb/cross-arb/internal/logger"
	reportpkg "github.com/dimryb/cross-arb/internal/report"
	"github.com/dimryb/cross-arb/internal/server/grpc"
	"github.com/dimryb/cross-arb/internal/server/http"
	"github.com/dimryb/cross-arb/internal/service"
	scan "github.com/dimryb/cross-arb/internal/service/scanner"
	"github.com/dimryb/cross-arb/internal/storage"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	printVersion()
	if flag.Arg(0) == "version" {
		return
	}

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg := logger.New(cfg.Log.Level)
	store := storage.NewTickerStore()
	application := app.NewApp(ctx, logg, store)
	arbitrageService := service.NewArbitrageService(application, cfg)
	reportSvc := reportpkg.NewReportService(logg, store)

	// --- Adapters ---
	mexcAdapter := adapter.NewMexcAdapter(logg, 3*time.Second)

	// Jupiter: собираем конфиг напрямую (заменяет фабрику)
	jupCfg, ok := cfg.Exchanges[config.JupExchange]
	if !ok {
		logg.Fatalf("exchange %s not found in configuration", config.JupExchange)
	}
	if !jupCfg.Enabled {
		logg.Fatalf("exchange %s is disabled", config.JupExchange)
	}

	pairMap := make(map[string]adapter.MintPair, len(jupCfg.Pairs))
	for symbol, p := range jupCfg.Pairs {
		if p.Base == "" || p.Quote == "" {
			logg.Fatalf("missing mint address for Jupiter pair %q", symbol)
		}
		pairMap[symbol] = adapter.MintPair{
			BaseMint:  p.Base,
			QuoteMint: p.Quote,
		}
	}

	jupiterAdapter := adapter.NewJupiterAdapter(logg, &adapter.JupiterAdapterConfig{
		BaseURL: jupCfg.BaseURL,
		Enabled: true,
		Timeout: jupCfg.Timeout,
		Pairs:   pairMap,
	})

	defer mexcAdapter.Close()
	defer jupiterAdapter.Close()

	adapters := []i.ExchangeAdapter{mexcAdapter, jupiterAdapter}
	scanner, err := scan.NewScannerFromConfig(logg, cfg.Scanner, adapters)
	if err != nil {
		logg.Error("Failed to create scanner", slog.Any("err", err))
		cancel()
		return
	}

	// Подписываемся и логируем возможности
	if cfg.Scanner.LogOpportunities {
		scan.SubscribeAndHandle(
			scanner,
			cfg.Scanner.Pairs,
			scan.LogOpportunities(logg),
		)
	}

	// TODO: вынести порты в конфиг
	grpcServer := grpc.NewServer(application, grpc.ServerConfig{Port: "9090"}, logg)

	go func() {
		httpServer := http.NewHTTPServer(store)
		if err := httpServer.Run(":8080"); err != nil {
			logg.Errorf("HTTP server error: %v", err)
			cancel()
		}
	}()

	go func() {
		if err := grpcServer.Run(); err != nil {
			logg.Errorf("gRPC server failed: %v", err)
			cancel()
		}
	}()

	// Запускаем сканер
	go func() {
		if err := scanner.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			logg.Errorf("Scanner stopped with error: %v", err)
			cancel()
		}
	}()

	reportSvc.Start()

	logg.Info("Starting app...")
	if err = arbitrageService.Run(); err != nil {
		logg.Errorf("Arbitrage service stopped with error: %v", err)
		cancel()
	} else {
		logg.Infof("Arbitrage service stopped gracefully")
	}
}
