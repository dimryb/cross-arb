package swap

import (
	"fmt"

	"github.com/dimryb/cross-arb/internal/api/jupiter"
	i "github.com/dimryb/cross-arb/internal/interface"
	blockchain "github.com/dimryb/cross-arb/internal/solana"
)

// Swapper объединяет Jupiter API и Solana блокчейн для полного цикла обмена.
type Swapper struct {
	apiClient    *jupiter.Client
	solanaClient *blockchain.Client
	logger       i.Logger
}

// NewSwapper создает сервис для полного цикла обмена.
func NewSwapper(logger i.Logger, jupiterURL, solanaRPCURL string) (*Swapper, error) {
	apiClient, err := jupiter.NewJupiterClient(logger, jupiterURL)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать Jupiter API клиент: %w", err)
	}

	solanaClient, err := blockchain.NewSolanaClient(logger, solanaRPCURL)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать Solana клиент: %w", err)
	}

	return &Swapper{
		apiClient:    apiClient,
		solanaClient: solanaClient,
		logger:       logger,
	}, nil
}
