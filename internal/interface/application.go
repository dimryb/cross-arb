package interfaces

import (
	"context"
)

type Application interface {
	Context() context.Context
	Logger() Logger
	TickerStore() TickerStore
}
