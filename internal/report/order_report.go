package report

import (
	"fmt"
	"time"

	"github.com/dimryb/cross-arb/internal/types"
)

func PrintOrderBookReport(results []types.OrderBookResult) {
	enable := false
	if enable {
		fmt.Printf("=== Стакан (обновлено: %s) ===\n", time.Now().Format("15:04:05.000"))
		for _, r := range results {
			if r.Error != nil {
				fmt.Printf("  [%s] Error: %v\n", r.Symbol, r.Error)
				continue
			}
			baseAsset := extractBaseAsset(r.Symbol)
			fmt.Printf("[%s] ASK (Можно купить):\n", baseAsset)
			for _, ask := range r.Data.Asks {
				fmt.Printf("Купить по цене: %.2f %s | Доступное количество: %.3f %s\n",
					ask.Price, "USDT", ask.Quantity, baseAsset)
			}
			fmt.Printf("[%s] BID (Можно Продать):\n", baseAsset)
			for _, bid := range r.Data.Bids {
				fmt.Printf("Продать по цене: %.2f %s | Доступное количество: %.3f %s\n",
					bid.Price, "USDT", bid.Quantity, baseAsset)
			}
		}
	}
}

func extractBaseAsset(symbol string) string {
	if len(symbol) > 4 && symbol[len(symbol)-4:] == "USDT" {
		return symbol[:len(symbol)-4]
	}
	return symbol
}
