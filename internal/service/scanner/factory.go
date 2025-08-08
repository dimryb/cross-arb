package scanner

import (
	"fmt"
	"time"

	"github.com/dimryb/cross-arb/internal/config"
	i "github.com/dimryb/cross-arb/internal/interface"
)

func NewScannerFromConfig(
	logger i.Logger,
	cfg config.ScannerConfig,
	adapters []i.ExchangeAdapter,
) (*Scanner, error) {
	duration, err := time.ParseDuration(cfg.Interval)
	if err != nil {
		return nil, fmt.Errorf("invalid scanner interval %q: %w", cfg.Interval, err)
	}

	if len(cfg.Pairs) == 0 {
		return nil, fmt.Errorf("scanner config: no pairs specified")
	}

	scanner := NewScanner(
		logger,
		WithInterval(duration),
		WithPairs(cfg.Pairs...),
		WithAdapters(adapters...),
	)

	return scanner, nil
}
