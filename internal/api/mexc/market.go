package mexc

import (
	"context"
	"encoding/json"
)

// DTOs (будут заполнены позже).

type PingResponse struct {
	// TODO: определить поля ответа /ping
}

type TimeResponse struct {
	// TODO: определить поля ответа /time
}

type APISymbolResponse struct {
	// TODO: определить поля ответа /defaultSymbols
}

type ExchangeInfoResponse struct {
	// TODO: определить поля ответа /exchangeInfo
}

type DepthResponse struct {
	Bids [][]string `json:"bids"`
	Asks [][]string `json:"asks"`
	// TODO: определить поля ответа /depth
}

type TradesResponse struct {
	// TODO: определить поля ответа /trades
}

type AggTradesResponse struct {
	// TODO: определить поля ответа /aggTrades
}

type KlineResponse struct {
	// TODO: определить поля ответа /klines
}

type AvgPriceResponse struct {
	// TODO: определить поля ответа /avgPrice
}

type Ticker24hrResponse struct {
	// TODO: определить поля ответа /ticker/24hr
}

type PriceResponse struct {
	// TODO: определить поля ответа /ticker/price
}

type BookTickerResponse struct {
	Symbol   string `json:"symbol"`
	BidPrice string `json:"bidPrice"`
	BidQty   string `json:"bidQty"`
	AskPrice string `json:"askPrice"`
	AskQty   string `json:"askQty"`
}

type SpotMarketClient struct {
	log    Logger
	client APIClient
}

func NewSpotMarketClient(log Logger, client APIClient) *SpotMarketClient {
	return &SpotMarketClient{log: log, client: client}
}

// # Реализация API-запросов
// ## Эндпоинты для получения рыночных данных (Market Data Endpoints)

// Ping 1. Проверка подключения к серверу (Test Connectivity).
func (s *SpotMarketClient) Ping(
	ctx context.Context,
	params map[string]string,
) (*PingResponse, error) {
	casePath := "/ping"
	s.log.Debug("Ping request to MEXC", "path", casePath, "params", params)
	resp, err := s.client.PublicGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в Ping", "error", err)
		return nil, err
	}

	var result PingResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal Ping response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// Time 2. Получить серверное время (Check Server Time).
func (s *SpotMarketClient) Time(
	ctx context.Context,
	params map[string]string,
) (*TimeResponse, error) {
	casePath := "/time"
	s.log.Debug("Time request to MEXC", "path", casePath, "params", params)
	resp, err := s.client.PublicGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в Time", "error", err)
		return nil, err
	}

	var result TimeResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal Time response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// APISymbol 3. Список торговых пар по умолчанию (API Default Symbol).
func (s *SpotMarketClient) APISymbol(
	ctx context.Context,
	params map[string]string,
) (*APISymbolResponse, error) {
	casePath := "/defaultSymbols"
	s.log.Debug("API symbol request to MEXC", "path", casePath, "params", params)
	resp, err := s.client.PublicGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в APISymbol", "error", err)
		return nil, err
	}

	var result APISymbolResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal APISymbol response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// ExchangeInfo 4. Информация о торгах (Exchange Information).
func (s *SpotMarketClient) ExchangeInfo(
	ctx context.Context,
	params map[string]string,
) (*ExchangeInfoResponse, error) {
	casePath := "/exchangeInfo"
	s.log.Debug("Exchange info request to MEXC", "path", casePath, "params", params)
	resp, err := s.client.PublicGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в ExchangeInfo", "error", err)
		return nil, err
	}

	var result ExchangeInfoResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal ExchangeInfo response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// Depth 5. Глубина стакана (Depth).
func (s *SpotMarketClient) Depth(
	ctx context.Context,
	params map[string]string,
) (*DepthResponse, error) {
	casePath := "/depth"
	s.log.Debug("Order book depth request to MEXC", "path", casePath, "params", params)
	resp, err := s.client.PublicGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в Depth", "error", err)
		return nil, err
	}

	var result DepthResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal Depth response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// Trades 6. Список последних сделок (Recent Trades List).
func (s *SpotMarketClient) Trades(
	ctx context.Context,
	params map[string]string,
) (*TradesResponse, error) {
	casePath := "/trades"
	s.log.Debug("Recent trades request to MEXC", "path", casePath, "params", params)
	resp, err := s.client.PublicGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в Trades", "error", err)
		return nil, err
	}

	var result TradesResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal Trades response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// AggTrades 7. Агрегированный список сделок (Aggregate Trades List).
func (s *SpotMarketClient) AggTrades(
	ctx context.Context,
	params map[string]string,
) (*AggTradesResponse, error) {
	casePath := "/aggTrades"
	s.log.Debug("Aggregate trades request to MEXC", "path", casePath, "params", params)
	resp, err := s.client.PublicGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в AggTrades", "error", err)
		return nil, err
	}

	var result AggTradesResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal AggTrades response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// Kline 8. Данные свечей (K-line Data).
func (s *SpotMarketClient) Kline(
	ctx context.Context,
	params map[string]string,
) (*KlineResponse, error) {
	casePath := "/klines"
	s.log.Debug("K-line data request to MEXC", "path", casePath, "params", params)
	resp, err := s.client.PublicGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в Kline", "error", err)
		return nil, err
	}

	var result KlineResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal Kline response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// AvgPrice 9. Средняя цена за период (Current Average Price).
func (s *SpotMarketClient) AvgPrice(
	ctx context.Context,
	params map[string]string,
) (*AvgPriceResponse, error) {
	casePath := "/avgPrice"
	s.log.Debug("Average price request to MEXC", "path", casePath, "params", params)
	resp, err := s.client.PublicGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в AvgPrice", "error", err)
		return nil, err
	}

	var result AvgPriceResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal AvgPrice response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// Ticker24hr 10. Статистика изменения цены за 24 часа (24hr Ticker Price Change Statistics).
func (s *SpotMarketClient) Ticker24hr(
	ctx context.Context,
	params map[string]string,
) (*Ticker24hrResponse, error) {
	casePath := "/ticker/24hr"
	s.log.Debug("24hr ticker stats request to MEXC", "path", casePath, "params", params)
	resp, err := s.client.PublicGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в Ticker24hr", "error", err)
		return nil, err
	}

	var result Ticker24hrResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal Ticker24hr response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// Price 11. Текущая цена символа (Symbol Price Ticker).
func (s *SpotMarketClient) Price(
	ctx context.Context,
	params map[string]string,
) (*PriceResponse, error) {
	casePath := "/ticker/price"
	s.log.Debug("Symbol price request to MEXC", "path", casePath, "params", params)
	resp, err := s.client.PublicGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в Price", "error", err)
		return nil, err
	}

	var result PriceResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal Price response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// BookTicker 12. Лучшие цены в стакане (Symbol Order Book Ticker).
func (s *SpotMarketClient) BookTicker(
	ctx context.Context,
	params map[string]string,
) (*BookTickerResponse, error) {
	casePath := "/ticker/bookTicker"
	s.log.Debug("Order book ticker request to MEXC", "path", casePath, "params", params)
	resp, err := s.client.PublicGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в BookTicker", "error", err)
		return nil, err
	}

	var result BookTickerResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal BookTicker response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}
