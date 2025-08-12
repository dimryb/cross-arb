package main

import (
	"context"
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
	"github.com/dimryb/cross-arb/internal/report"
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
	report := report.NewReportService(logg, store)

	mexcAdapter := adapter.NewMexcAdapter(logg, 3*time.Second)
	jupiterAdapter, err := adapter.NewJupiterAdapterFromConfig(logg, cfg)
	if err != nil {
		logg.Fatalf("Failed to create Jupiter adapter: %v", err)
	}

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

	// TODO: Вынести параметр порта grpc в конфиг файл
	grpcServer := grpc.NewServer(application, grpc.ServerConfig{Port: "9090"}, logg)

	go func() {
		httpServer := http.NewHTTPServer(store)
		// TODO: Вынести параметр порта http в конфиг файл
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

	go func() {
		err = scanner.Run(ctx)
		if err != nil {
			return
		}
	}()

	report.Start()

	logg.Info("Starting app...")
	if err = arbitrageService.Run(); err != nil {
		logg.Errorf("Arbitrage service stopped with error: %v", err)
		cancel()
	} else {
		logg.Infof("Arbitrage service stopped gracefully")
	}
}
