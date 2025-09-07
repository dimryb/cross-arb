package scan

import (
	"context"

	"github.com/dimryb/cross-arb/internal/entity"
)

// ArbOpportunityUseCase потребляет поток цен и публикует найденные арбитражные возможности.
// Должен корректно завершаться по ctx; владение выходным каналом out — у вызывающей стороны.
type ArbOpportunityUseCase interface {
	Detect(
		ctx context.Context,
		in <-chan entity.ExecutableQuote,
		out chan<- entity.ArbOpportunity,
	) error
}

// NoopOpportunityUseCase — заглушка выходного юзкейса.
type NoopOpportunityUseCase struct{}

func NewNoopOpportunityUseCase() *NoopOpportunityUseCase { return &NoopOpportunityUseCase{} }

func (n *NoopOpportunityUseCase) Detect(
	ctx context.Context,
	_ <-chan entity.ExecutableQuote,
	_ chan<- entity.ArbOpportunity,
) error {
	<-ctx.Done()
	return ctx.Err()
}
