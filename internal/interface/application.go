package interfaces

import (
	"context"
)

//go:generate mockgen -source=application.go -package=mocks -destination=../../mocks/mock_application.go
type Application interface {
	Context() context.Context
	Logger() Logger
	TickerStore() TickerStore
}
