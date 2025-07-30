package storage

import (
	"sync"
)

// TickerData — данные одного тикера.
type TickerData struct {
	Symbol   string  `json:"symbol" proto:"symbol"`
	Exchange string  `json:"exchange" proto:"exchange"`
	BidPrice float64 `json:"bidPrice" proto:"bid_price"`
	BidQty   float64 `json:"bidQty" proto:"bid_qty"`
	AskPrice float64 `json:"askPrice" proto:"ask_price"`
	AskQty   float64 `json:"askQty" proto:"ask_qty"`
}

// Equal Сравнение двух тикеров (для оптимизации уведомлений).
func (t TickerData) Equal(other TickerData) bool {
	return t.Symbol == other.Symbol &&
		t.Exchange == other.Exchange &&
		t.BidPrice == other.BidPrice &&
		t.BidQty == other.BidQty &&
		t.AskPrice == other.AskPrice &&
		t.AskQty == other.AskQty
}

// TickerEvent — событие обновления тикера.
type TickerEvent struct {
	Ticker TickerData
}

// TickerSubscriber — канал для получения событий.
type TickerSubscriber chan TickerEvent

// TickerStore — потокобезопасное хранилище тикеров с поддержкой подписок.
type TickerStore struct {
	mu          sync.RWMutex
	tickers     map[string]TickerData
	subscribers []TickerSubscriber // Активные подписчики
}

// NewTickerStore — создаёт новое хранилище.
func NewTickerStore() *TickerStore {
	return &TickerStore{
		tickers:     make(map[string]TickerData),
		subscribers: make([]TickerSubscriber, 0),
	}
}

// Set — добавляет или обновляет тикер.
// Если данные изменились — уведомляет подписчиков.
func (s *TickerStore) Set(t TickerData) {
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
func (s *TickerStore) GetAll() []TickerData {
	s.mu.RLock()
	defer s.mu.RUnlock()

	all := make([]TickerData, 0, len(s.tickers))
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
func (s *TickerStore) notifySubscribers(ticker TickerData) {
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
