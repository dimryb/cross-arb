package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/dimryb/cross-arb/internal/app"
	"github.com/dimryb/cross-arb/internal/config"
	"github.com/dimryb/cross-arb/internal/logger"
	"github.com/dimryb/cross-arb/internal/server/http"
	"github.com/dimryb/cross-arb/internal/service"
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
	application := app.NewApp(logg)
	store := storage.NewTickerStore()
	arbitrageService := service.NewArbitrageService(ctx, application, logg, cfg, store)

	go func() {
		httpServer := http.NewHTTPServer(store)
		if err := httpServer.Run(":8080"); err != nil {
			logg.Errorf("HTTP server error: %v", err)
			cancel()
		}
	}()

	logg.Info("Starting app...")
	if err = arbitrageService.Run(); err != nil {
		logg.Errorf("Arbitrage service stopped with error: %v", err)
		cancel()
	} else {
		logg.Infof("Arbitrage service stopped gracefully")
	}
}
