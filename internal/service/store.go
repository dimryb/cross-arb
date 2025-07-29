package service

import (
	"sync"
)

type TickerData struct {
	Symbol   string  `json:"symbol"`
	Exchange string  `json:"exchange"`
	BidPrice float64 `json:"bidPrice"`
	BidQty   float64 `json:"bidQty"`
	AskPrice float64 `json:"askPrice"`
	AskQty   float64 `json:"askQty"`
}

type TickerStore struct {
	mu      sync.RWMutex
	tickers map[string]TickerData
}

func NewTickerStore() *TickerStore {
	return &TickerStore{
		tickers: make(map[string]TickerData),
	}
}

func (s *TickerStore) Set(t TickerData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := t.Symbol + "-" + t.Exchange
	s.tickers[key] = t
}

func (s *TickerStore) GetAll() []TickerData {
	s.mu.RLock()
	defer s.mu.RUnlock()

	all := make([]TickerData, 0, len(s.tickers))
	for _, v := range s.tickers {
		all = append(all, v)
	}
	return all
}
