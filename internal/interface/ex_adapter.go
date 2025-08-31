package interfaces

import (
	"context"

	"github.com/dimryb/cross-arb/internal/entity"
)

// EXAdapter описывает минимальный набор функций,
// который должен реализовать адаптер конкретной биржи,
// чтобы сканер мог получать рыночные данные и рассчитывать спреды.
//
// Все методы обязаны быть потокобезопасными и уважать контекст.
// Ошибки возвращаются вызывающему для дальнейшей обработки.
type EXAdapter interface {
	// Name возвращает короткое идентификатор биржи в snake_case,
	// используемый в логах и метриках.
	Name() string

	// TradingFee возвращает комиссию мейкера и тейкера для пары.
	// Значения указаны как доли (0.001 == 0.1 %).
	TradingFee(pair string) (maker, taker float64)

	// Close освобождает ресурсы (например, закрывает WebSocket-соединения).
	// Должен быть идемпотентным.
	Close() error
}

// DEXAdapter — адаптер для DEX-бирж (агрегаторы маршрутов, AMM и т.п.).
// Объем-зависимое квотирование: возвращает эффективные котировки для заданного объёма baseAmount.
type DEXAdapter interface {
	EXAdapter
	// Quote возвращает эффективные котировки bid/ask в QUOTE за BASE для указанного объёма baseAmount (в BASE).
	// Реализация должна учесть маршрутизацию/слиппедж. Для малых объёмов результат может совпадать с top-of-book.
	Quote(ctx context.Context, pair string, baseAmount float64) (bid, ask float64, err error)
}

// CEXAdapter — адаптер для централизованных бирж, предоставляющий глубину стакана.
// Реализация может быть по REST или WS; метод обязан уважать ctx.
type CEXAdapter interface {
	EXAdapter
	// OrderBookDepth возвращает полную (или ограниченную limit) глубину стакана для пары.
	// Предполагается, что bids упорядочены по убыванию, а asks — по возрастанию цены.
	OrderBookDepth(ctx context.Context, pair string, limit int) (entity.OrderBook, error)
}
