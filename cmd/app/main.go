package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/dimryb/cross-arb/internal/app"
	"github.com/dimryb/cross-arb/internal/config"
	"github.com/dimryb/cross-arb/internal/logger"
	"github.com/dimryb/cross-arb/internal/server/grpc"
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
	store := storage.NewTickerStore()
	application := app.NewApp(ctx, logg, store)
	arbitrageService := service.NewArbitrageService(application, cfg)

	zapLog, _ := zap.NewProduction() // нужен адаптерам/сканеру
	defer zapLog.Sync()

	mexcAdapter := service.NewMexcAdapter(zapLog, 3*time.Second)
	jupPairs := map[string][2]string{
		"SOL/USDT": {
			"So11111111111111111111111111111111111111112",  // SOL (wSOL)
			"Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB", // USDT
		},
	}

	jupiterAdapter := service.NewJupiterAdapter(zapLog, jupPairs, 3*time.Second)
	defer mexcAdapter.Close()
	defer jupiterAdapter.Close()

	scanner := service.NewScanner(
		zapLog,
		service.WithInterval(2*time.Second),
		service.WithPairs("SOL/USDT"),
		service.WithAdapters(mexcAdapter, jupiterAdapter),
	)

	// Подписываемся и логируем возможности
	for _, pair := range []string{"SOL/USDT"} {
		ch, _ := scanner.Subscribe(pair, 10)
		go func(_ string, c <-chan service.Opportunity) {
			for opp := range c {
				logg.Infof("Арбитраж %s: BUY %s @ %.4f → SELL %s @ %.4f  (%.4f %%)",
					opp.Pair, opp.BuyOn, opp.BuyPrice, opp.SellOn, opp.SellPrice, opp.SpreadPct)
			}
		}(pair, ch)
	}

	grpcServer := grpc.NewServer(application, grpc.ServerConfig{Port: "9090"}, logg)

	go func() {
		httpServer := http.NewHTTPServer(store)
		if err := httpServer.Run(":8080"); err != nil {
			logg.Errorf("HTTP server error: %v", err)
			cancel()
		}
	}()

	go func() {
		logg.Info("Starting gRPC server on :9090")
		if err := grpcServer.Run(); err != nil {
			logg.Errorf("gRPC server failed: %v", err)
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
