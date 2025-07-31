package interfaces

import "github.com/dimryb/cross-arb/internal/types"

type TickerStore interface {
	Set(t types.TickerData)
	GetAll() []types.TickerData
}
