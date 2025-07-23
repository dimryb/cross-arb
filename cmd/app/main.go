package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dimryb/cross-arb/internal/app"
	"github.com/dimryb/cross-arb/internal/config"
	"github.com/dimryb/cross-arb/internal/logger"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	fmt.Println("config", cfg)

	logg := logger.New(cfg.Log.Level)
	application := app.NewApp(logg)
	_ = application

	logg.Info("Starting app...")
	logg.Info("Stopped gracefully")
}
