// TODO: DELETE
package adapter

import (
	"fmt"

	"github.com/dimryb/cross-arb/internal/config"
	i "github.com/dimryb/cross-arb/internal/interface"
)

func NewJupiterAdapterFromConfig(logg i.Logger, cfg *config.CrossArbConfig) (*JupiterAdapter, error) {
	jupConfig, ok := cfg.Exchanges[config.JupExchange]
	if !ok {
		return nil, fmt.Errorf("exchange %s not found in configuration", config.JupExchange)
	}
	if !jupConfig.Enabled {
		return nil, fmt.Errorf("exchange %s is disabled", config.JupExchange)
	}

	pairMap := make(map[string]MintPair, len(jupConfig.Pairs))
	for symbol, pair := range jupConfig.Pairs {
		if pair.Base == "" || pair.Quote == "" {
			return nil, fmt.Errorf("missing mint address for Jupiter pair %q", symbol)
		}
		pairMap[symbol] = MintPair{
			BaseMint:  pair.Base,
			QuoteMint: pair.Quote,
		}
	}

	jupiterAdapter := NewJupiterAdapter(logg, &JupiterAdapterConfig{
		BaseURL: jupConfig.BaseURL,
		Timeout: jupConfig.Timeout,
		Pairs:   pairMap,
		Enabled: true,
	})

	return jupiterAdapter, nil
}
