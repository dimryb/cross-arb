package storage

import (
	"sync"

	"github.com/dimryb/cross-arb/internal/types"
)

// TickerEvent — событие обновления тикера.
type TickerEvent struct {
	Ticker types.TickerData
}

// TickerSubscriber — канал для получения событий.
type TickerSubscriber chan TickerEvent

// TickerStore — потокобезопасное хранилище тикеров с поддержкой подписок.
type TickerStore struct {
	mu          sync.RWMutex
	tickers     map[string]types.TickerData
	subscribers []TickerSubscriber // Активные подписчики
}

// NewTickerStore — создаёт новое хранилище.
func NewTickerStore() *TickerStore {
	return &TickerStore{
		tickers:     make(map[string]types.TickerData),
		subscribers: make([]TickerSubscriber, 0),
	}
}

// Set — добавляет или обновляет тикер.
// Если данные изменились — уведомляет подписчиков.
func (s *TickerStore) Set(t types.TickerData) {
	key := t.Symbol + "-" + t.Exchange

	s.mu.Lock()
	old, exists := s.tickers[key]
	s.tickers[key] = t
	s.mu.Unlock()

	if !exists || !old.Equal(t) {
		s.notifySubscribers(t)
	}
}

// GetAll — возвращает копию всех тикеров (используется HTTP).
func (s *TickerStore) GetAll() []types.TickerData {
	s.mu.RLock()
	defer s.mu.RUnlock()

	all := make([]types.TickerData, 0, len(s.tickers))
	for _, v := range s.tickers {
		all = append(all, v)
	}
	return all
}

// AddSubscriber — регистрирует нового подписчика.
// Возвращает канал, из которого он будет читать события.
func (s *TickerStore) AddSubscriber() TickerSubscriber {
	ch := make(TickerSubscriber, 10) // буфер на 10 событий
	s.mu.Lock()
	s.subscribers = append(s.subscribers, ch)
	s.mu.Unlock()
	return ch
}

// notifySubscribers — отправляет событие всем подписчикам.
func (s *TickerStore) notifySubscribers(ticker types.TickerData) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event := TickerEvent{Ticker: ticker}
	for _, ch := range s.subscribers {
		select {
		case ch <- event:
			// успешно отправлено
		default:
			// канал переполнен — клиент слишком медленный
			// можно игнорировать или закрыть (по политике)
		}
	}
}
