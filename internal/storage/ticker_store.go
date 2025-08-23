package storage

import (
	"context"
	"sync"

	"github.com/dimryb/cross-arb/internal/entity"
	i "github.com/dimryb/cross-arb/internal/interface"
)

// TickerStore — потокобезопасное хранилище тикеров с поддержкой подписок.
type TickerStore struct {
	mu          sync.RWMutex
	tickers     map[string]entity.TickerData
	subscribers []chan entity.TickerEvent // Каналы для рассылки
}

// NewTickerStore — создаёт новое хранилище.
func NewTickerStore() *TickerStore {
	return &TickerStore{
		tickers:     make(map[string]entity.TickerData),
		subscribers: make([]chan entity.TickerEvent, 0),
	}
}

// Set — добавляет или обновляет тикер.
// Если данные изменились — уведомляет подписчиков.
func (s *TickerStore) Set(t entity.TickerData) {
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
func (s *TickerStore) GetAll() []entity.TickerData {
	s.mu.RLock()
	defer s.mu.RUnlock()

	all := make([]entity.TickerData, 0, len(s.tickers))
	for _, v := range s.tickers {
		all = append(all, v)
	}
	return all
}

// AddSubscriber — регистрирует нового подписчика.
// Возвращает обёртку i.TickerSubscriber для безопасного чтения и управления.
func (s *TickerStore) AddSubscriber() i.TickerSubscriber {
	ch := make(chan entity.TickerEvent, 10)
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
func (s *TickerStore) notifySubscribers(ticker entity.TickerData) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event := entity.TickerEvent{Ticker: ticker}
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
func (s *TickerStore) removeSubscriber(ch chan entity.TickerEvent) {
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
