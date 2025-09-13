package mexc

import (
	"context"

	"github.com/go-resty/resty/v2"
)

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
) (*resty.Response, error) {
	casePath := "/ping"
	s.log.Debug("Ping request to MEXC", "path", casePath, "params", params)
	resp, err := s.client.PublicGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в Ping", "error", err)
		return nil, err
	}

	return resp, nil
}

// Time 2. Получить серверное время (Check Server Time).
func (s *SpotMarketClient) Time(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/time"
	s.log.Debug("Time request to MEXC", "path", casePath, "params", params)
	resp, err := s.client.PublicGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в Time", "error", err)
		return nil, err
	}

	return resp, nil
}

// APISymbol 3. Список торговых пар по умолчанию (API Default Symbol).
func (s *SpotMarketClient) APISymbol(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/defaultSymbols"
	s.log.Debug("API symbol request to MEXC", "path", casePath, "params", params)
	resp, err := s.client.PublicGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в APISymbol", "error", err)
		return nil, err
	}

	return resp, nil
}

// ExchangeInfo 4. Информация о торгах (Exchange Information).
func (s *SpotMarketClient) ExchangeInfo(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/exchangeInfo"
	s.log.Debug("Exchange info request to MEXC", "path", casePath, "params", params)
	resp, err := s.client.PublicGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в ExchangeInfo", "error", err)
		return nil, err
	}

	return resp, nil
}

// Depth 5. Глубина стакана (Depth).
func (s *SpotMarketClient) Depth(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/depth"
	s.log.Debug("Order book depth request to MEXC", "path", casePath, "params", params)
	resp, err := s.client.PublicGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в Depth", "error", err)
		return nil, err
	}

	return resp, nil
}

// Trades 6. Список последних сделок (Recent Trades List).
func (s *SpotMarketClient) Trades(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/trades"
	s.log.Debug("Recent trades request to MEXC", "path", casePath, "params", params)

	resp, err := s.client.PublicGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в Trades", "error", err)
		return nil, err
	}

	return resp, nil
}

// AggTrades 7. Агрегированный список сделок (Aggregate Trades List).
func (s *SpotMarketClient) AggTrades(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/aggTrades"
	s.log.Debug("Aggregate trades request to MEXC", "path", casePath, "params", params)
	resp, err := s.client.PublicGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в AggTrades", "error", err)
		return nil, err
	}

	return resp, nil
}

// Kline 8. Данные свечей (K-line Data).
func (s *SpotMarketClient) Kline(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/klines"
	s.log.Debug("K-line data request to MEXC", "path", casePath, "params", params)
	resp, err := s.client.PublicGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в Kline", "error", err)
		return nil, err
	}

	return resp, nil
}

// AvgPrice 9. Средняя цена за период (Current Average Price).
func (s *SpotMarketClient) AvgPrice(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/avgPrice"
	s.log.Debug("Average price request to MEXC", "path", casePath, "params", params)
	resp, err := s.client.PublicGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в AvgPrice", "error", err)
		return nil, err
	}

	return resp, nil
}

// Ticker24hr 10. Статистика изменения цены за 24 часа (24hr Ticker Price Change Statistics).
func (s *SpotMarketClient) Ticker24hr(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/ticker/24hr"
	s.log.Debug("24hr ticker stats request to MEXC", "path", casePath, "params", params)
	resp, err := s.client.PublicGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в Ticker24hr", "error", err)
		return nil, err
	}

	return resp, nil
}

// Price 11. Текущая цена символа (Symbol Price Ticker).
func (s *SpotMarketClient) Price(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/ticker/price"
	s.log.Debug("Symbol price request to MEXC", "path", casePath, "params", params)
	resp, err := s.client.PublicGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в Price", "error", err)
		return nil, err
	}

	return resp, nil
}

// BookTicker 12. Лучшие цены в стакане (Symbol Order Book Ticker).
func (s *SpotMarketClient) BookTicker(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/ticker/bookTicker"
	s.log.Debug("Order book ticker request to MEXC", "path", casePath, "params", params)
	resp, err := s.client.PublicGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в BookTicker", "error", err)
		return nil, err
	}

	return resp, nil
}
