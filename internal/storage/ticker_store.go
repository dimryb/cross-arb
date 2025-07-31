package storage

import (
	"context"
	"sync"

	i "github.com/dimryb/cross-arb/internal/interface"
	"github.com/dimryb/cross-arb/internal/types"
)

// TickerStore — потокобезопасное хранилище тикеров с поддержкой подписок.
type TickerStore struct {
	mu          sync.RWMutex
	tickers     map[string]types.TickerData
	subscribers []chan types.TickerEvent // Каналы для рассылки
}

// NewTickerStore — создаёт новое хранилище.
func NewTickerStore() *TickerStore {
	return &TickerStore{
		tickers:     make(map[string]types.TickerData),
		subscribers: make([]chan types.TickerEvent, 0),
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
// Возвращает обёртку i.TickerSubscriber для безопасного чтения и управления.
func (s *TickerStore) AddSubscriber() i.TickerSubscriber {
	ch := make(chan types.TickerEvent, 10)
	ctx, cancel := context.WithCancel(context.Background())
	sub := newSubscriber(ctx, ch, cancel)

	s.mu.Lock()
	s.subscribers = append(s.subscribers, ch)
	s.mu.Unlock()

	go func() {
		<-ctx.Done()
		s.removeSubscriber(ch)
	}()

	return sub
}

// notifySubscribers — отправляет событие всем активным подписчикам (каналам).
func (s *TickerStore) notifySubscribers(ticker types.TickerData) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event := types.TickerEvent{Ticker: ticker}
	for _, ch := range s.subscribers {
		select {
		case ch <- event:
			// Успешно отправлено
		default:
			// Клиент слишком медленный, пропускаем (или можно закрыть в продвинутой версии)
		}
	}
}

// Вызывается при завершении контекста подписки.
func (s *TickerStore) removeSubscriber(ch chan types.TickerEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for ind, subscriberChan := range s.subscribers {
		if subscriberChan == ch {
			// Удаляем из слайса
			s.subscribers = append(s.subscribers[:ind], s.subscribers[ind+1:]...)
			// Закрываем канал
			close(ch)
			return
		}
	}
}
