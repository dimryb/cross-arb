package interfaces

//go:generate mockgen -source=application.go -package=mocks -destination=../../mocks/mock_application.go
type Application interface {
	Logger() Logger
	TickerStore() TickerStore
}
