package report

import (
	"fmt"
	"time"

	"github.com/dimryb/cross-arb/internal/types"
)

func PrintTickersReport(results []types.Result) {
	fmt.Printf("=== Обновление цен (%s) ===\n", time.Now().Format("15:04:05.000"))
	for _, r := range results {
		PrintTicker(r.Data)
	}
	fmt.Println()
}

func PrintTicker(t types.BookTicker) {
	fmt.Printf(
		"  [%s] -> покупка: %s (%s) | продажа: %s (%s)\n",
		t.Symbol,
		t.BidPrice, t.BidQty,
		t.AskPrice, t.AskQty,
	)
}
