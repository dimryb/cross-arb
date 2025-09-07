package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/dimryb/cross-arb/internal/app"
	"github.com/dimryb/cross-arb/internal/config"
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

	app.NewApp(ctx, cancel, cfg).Run()
}
