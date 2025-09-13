package mexc

import (
	"context"

	"github.com/go-resty/resty/v2"
)

// SpotList — клиент для работы с MEXC Spot API.
type SpotList struct {
	log    Logger
	client *Client
}

// NewSpotList создаёт новый клиент для Spot API.
func NewSpotList(log Logger, client *Client) *SpotList {
	return &SpotList{
		log:    log,
		client: client,
	}
}

// # Реализация API-запросов
// ## Эндпоинты для получения рыночных данных (Market Data Endpoints)

// Ping 1. Проверка подключения к серверу (Test Connectivity).
func (s *SpotList) Ping(ctx context.Context, params map[string]string) (*resty.Response, error) {
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
func (s *SpotList) Time(ctx context.Context, params map[string]string) (*resty.Response, error) {
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
func (s *SpotList) APISymbol(ctx context.Context, params map[string]string) (*resty.Response, error) {
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
func (s *SpotList) ExchangeInfo(ctx context.Context, params map[string]string) (*resty.Response, error) {
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
func (s *SpotList) Depth(ctx context.Context, params map[string]string) (*resty.Response, error) {
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
func (s *SpotList) Trades(ctx context.Context, params map[string]string) (*resty.Response, error) {
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
func (s *SpotList) AggTrades(ctx context.Context, params map[string]string) (*resty.Response, error) {
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
func (s *SpotList) Kline(ctx context.Context, params map[string]string) (*resty.Response, error) {
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
func (s *SpotList) AvgPrice(ctx context.Context, params map[string]string) (*resty.Response, error) {
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
func (s *SpotList) Ticker24hr(ctx context.Context, params map[string]string) (*resty.Response, error) {
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
func (s *SpotList) Price(ctx context.Context, params map[string]string) (*resty.Response, error) {
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
func (s *SpotList) BookTicker(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/ticker/bookTicker"
	s.log.Debug("Order book ticker request to MEXC", "path", casePath, "params", params)
	resp, err := s.client.PublicGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в BookTicker", "error", err)
		return nil, err
	}

	return resp, nil
}

// ## 1. Суб‑аккаунты (Sub‑Account Endpoints)

// CreateSub 1. Создать суб‑аккаунт (Create a Sub-account).
func (s *SpotList) CreateSub(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/sub-account/virtualSubAccount"
	s.log.Debug("CreateSub", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка CreateSub", "error", err)
		return nil, err
	}
	return resp, nil
}

// QuerySub 2. Получить список суб‑аккаунтов (Query Sub-account List).
func (s *SpotList) QuerySub(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/sub-account/list"
	s.log.Debug("QuerySub", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка QuerySub", "error", err)
		return nil, err
	}
	return resp, nil
}

// CreateSubApikey 3. Создать API‑ключ для суб‑аккаунта (Create an apiKey for a sub-account).
func (s *SpotList) CreateSubApikey(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/sub-account/apiKey" //nolint:goconst
	s.log.Debug("CreateSubApikey", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка CreateSubApikey", "error", err)
		return nil, err
	}
	return resp, nil
}

// QuerySubApikey 4. Получить API‑ключи суб‑аккаунта (apiKey Query the apiKey of a sub-account).
func (s *SpotList) QuerySubApikey(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/sub-account/apiKey"
	s.log.Debug("QuerySubApikey", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка QuerySubApikey", "error", err)
		return nil, err
	}
	return resp, nil
}

// DeleteSubApikey 5. Удалить API‑ключ суб‑аккаунта (apiKey Delete the apiKey of a sub-account).
func (s *SpotList) DeleteSubApikey(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/sub-account/apiKey"
	s.log.Debug("DeleteSubApikey", "path", casePath, "params", params)
	resp, err := s.client.PrivateDelete(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка DeleteSubApikey", "error", err)
		return nil, err
	}
	return resp, nil
}

// UniTransfer 6. Универсальный перевод между аккаунтами (Universal Transfer).
func (s *SpotList) UniTransfer(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/capital/sub-account/universalTransfer"
	s.log.Debug("UniTransfer", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка UniTransfer", "error", err)
		return nil, err
	}
	return resp, nil
}

// QueryUniTransfer 7. История универсальных переводов (Query Universal Transfer History).
func (s *SpotList) QueryUniTransfer(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/capital/sub-account/universalTransfer"
	s.log.Debug("QueryUniTransfer", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка QueryUniTransfer", "error", err)
		return nil, err
	}
	return resp, nil
}

// ## 2. Торговля (Spot Account & Trade)

// SelfSymbols 1. Список разрешённых символов пользователя (User API default symbol).
func (s *SpotList) SelfSymbols(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/selfSymbols"
	s.log.Debug("SelfSymbols", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка SelfSymbols", "error", err)
		return nil, err
	}
	return resp, nil
}

// TestOrder 2. Тестовый ордер (Test New Order).
func (s *SpotList) TestOrder(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/order/test"
	s.log.Debug("TestOrder", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка TestOrder", "error", err)
		return nil, err
	}
	return resp, nil
}

// PlaceOrder 3. Разместить ордер (New Order).
func (s *SpotList) PlaceOrder(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/order" //nolint:goconst
	s.log.Debug("PlaceOrder", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка PlaceOrder", "error", err)
		return nil, err
	}
	return resp, nil
}

// BatchOrder 4. Пакетное размещение ордеров (Batch Orders).
func (s *SpotList) BatchOrder(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/batchOrders"
	s.log.Debug("BatchOrder", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка BatchOrder", "error", err)
		return nil, err
	}
	return resp, nil
}

// CancelOrder 5. Отменить ордер (Cancel Order).
func (s *SpotList) CancelOrder(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/order"
	s.log.Debug("CancelOrder", "path", casePath, "params", params)
	resp, err := s.client.PrivateDelete(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка CancelOrder", "error", err)
		return nil, err
	}
	return resp, nil
}

// CancelAllOrders 6. Отменить все ордера по символу (Cancel all Open Orders on a Symbol).
func (s *SpotList) CancelAllOrders(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/openOrders"
	s.log.Debug("CancelAllOrders", "path", casePath, "params", params)
	resp, err := s.client.PrivateDelete(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка CancelAllOrders", "error", err)
		return nil, err
	}
	return resp, nil
}

// QueryOrder 7. Информация об ордере (Query Order).
func (s *SpotList) QueryOrder(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/order"
	s.log.Debug("QueryOrder", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка QueryOrder", "error", err)
		return nil, err
	}
	return resp, nil
}

// OpenOrder 8. Открытые ордера (Current Open Orders).
func (s *SpotList) OpenOrder(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/openOrders"
	s.log.Debug("OpenOrder", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка OpenOrder", "error", err)
		return nil, err
	}
	return resp, nil
}

// AllOrders 9. Все ордера (All Orders).
func (s *SpotList) AllOrders(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/allOrders"
	s.log.Debug("AllOrders", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка AllOrders", "error", err)
		return nil, err
	}
	return resp, nil
}

// SpotAccountInfo 10. Информация об аккаунте (Account Information).
func (s *SpotList) SpotAccountInfo(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/account"
	s.log.Debug("SpotAccountInfo", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка SpotAccountInfo", "error", err)
		return nil, err
	}
	return resp, nil
}

// SpotMyTrade 11. История сделок (Account Trade List).
func (s *SpotList) SpotMyTrade(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/myTrades"
	s.log.Debug("SpotMyTrade", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка SpotMyTrade", "error", err)
		return nil, err
	}
	return resp, nil
}

// MxDeduct 12. Включить MX‑дедукцию (Enable MX Deduct).
func (s *SpotList) MxDeduct(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/mxDeduct/enable"
	s.log.Debug("MxDeduct", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка MxDeduct", "error", err)
		return nil, err
	}
	return resp, nil
}

// QueryMxDeduct 13. Статус MX‑дедукции (Query MX Deduct Status).
func (s *SpotList) QueryMxDeduct(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/mxDeduct/enable"
	s.log.Debug("QueryMxDeduct", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка QueryMxDeduct", "error", err)
		return nil, err
	}
	return resp, nil
}

// ## 3. Кошелёк (Wallet Endpoints)

// QueryCurrencyInfo 1. Информация о валюте (Query the currency information).
func (s *SpotList) QueryCurrencyInfo(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/capital/config/getall"
	s.log.Debug("QueryCurrencyInfo", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка QueryCurrencyInfo", "error", err)
		return nil, err
	}
	return resp, nil
}

// Withdraw 2. Вывод средств (Withdraw).
func (s *SpotList) Withdraw(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/capital/withdraw/apply"
	s.log.Debug("Withdraw", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка Withdraw", "error", err)
		return nil, err
	}
	return resp, nil
}

// CancelWithdraw 3. Отменить вывод (Cancel withdraw).
func (s *SpotList) CancelWithdraw(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/capital/withdraw"
	s.log.Debug("CancelWithdraw", "path", casePath, "params", params)
	resp, err := s.client.PrivateDelete(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка CancelWithdraw", "error", err)
		return nil, err
	}
	return resp, nil
}

// DepositHistory 4. История депозитов (Deposit History).
func (s *SpotList) DepositHistory(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/capital/deposit/hisrec"
	s.log.Debug("DepositHistory", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка DepositHistory", "error", err)
		return nil, err
	}
	return resp, nil
}

// WithdrawHistory 5. История выводов (Withdraw History).
func (s *SpotList) WithdrawHistory(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/capital/withdraw/history"
	s.log.Debug("WithdrawHistory", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка WithdrawHistory", "error", err)
		return nil, err
	}
	return resp, nil
}

// GenDepositAddress 6. Сгенерировать адрес депозита (Generate deposit address).
func (s *SpotList) GenDepositAddress(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/capital/deposit/address"
	s.log.Debug("GenDepositAddress", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка GenDepositAddress", "error", err)
		return nil, err
	}
	return resp, nil
}

// DepositAddress 7. Получить адрес депозита (Deposit Address).
func (s *SpotList) DepositAddress(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/capital/deposit/address"
	s.log.Debug("DepositAddress", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка DepositAddress", "error", err)
		return nil, err
	}
	return resp, nil
}

// WithdrawAddress 8. Получить адрес вывода (Withdraw Address).
func (s *SpotList) WithdrawAddress(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/capital/withdraw/address"
	s.log.Debug("WithdrawAddress", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка WithdrawAddress", "error", err)
		return nil, err
	}
	return resp, nil
}

// Transfer 9. Универсальный перевод (User Universal Transfer).
func (s *SpotList) Transfer(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/capital/transfer"
	s.log.Debug("Transfer", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка Transfer", "error", err)
		return nil, err
	}
	return resp, nil
}

// TransferHistory 10. История переводов (Query User Universal Transfer History).
func (s *SpotList) TransferHistory(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/capital/transfer"
	s.log.Debug("TransferHistory", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка TransferHistory", "error", err)
		return nil, err
	}
	return resp, nil
}

// TransferHistoryByID 11. История перевода по tranId (Query User Universal Transfer History （by tranId）).
func (s *SpotList) TransferHistoryByID(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/capital/transfer/tranId"
	s.log.Debug("TransferHistoryByID", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка TransferHistoryByID", "error", err)
		return nil, err
	}
	return resp, nil
}

// ConvertList 12. Список активов для конвертации (Get Assets That Can Be Converted Into MX).
func (s *SpotList) ConvertList(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/capital/convert/list"
	s.log.Debug("ConvertList", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка ConvertList", "error", err)
		return nil, err
	}
	return resp, nil
}

// Convert 13. Конвертация мелких активов (Dust Transfer).
func (s *SpotList) Convert(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/capital/convert"
	s.log.Debug("Convert", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка Convert", "error", err)
		return nil, err
	}
	return resp, nil
}

// ConvertHistory 14. История конвертаций (DustLog).
func (s *SpotList) ConvertHistory(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/capital/convert"
	s.log.Debug("ConvertHistory", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка ConvertHistory", "error", err)
		return nil, err
	}
	return resp, nil
}

// ETFInfo 15. Информация об ETF (Get ETF info).
func (s *SpotList) ETFInfo(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/etf/info"
	s.log.Debug("ETFInfo", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка ETFInfo", "error", err)
		return nil, err
	}
	return resp, nil
}

// InternalTransfer 16. Внутренний перевод (Internal Transfer).
func (s *SpotList) InternalTransfer(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/capital/transfer/internal"
	s.log.Debug("InternalTransfer", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка InternalTransfer", "error", err)
		return nil, err
	}
	return resp, nil
}

// InternalTransferHistory 17. История внутренних переводов (Internal Transfer History).
func (s *SpotList) InternalTransferHistory(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/capital/transfer/internal"
	s.log.Debug("InternalTransferHistory", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка InternalTransferHistory", "error", err)
		return nil, err
	}
	return resp, nil
}

// ## 4. WS ListenKey

// CreateListenKey 1. Создать ListenKey (Listen Key  Create a ListenKey).
func (s *SpotList) CreateListenKey(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/userDataStream" //nolint:goconst
	s.log.Debug("CreateListenKey", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка CreateListenKey", "error", err)
		return nil, err
	}
	return resp, nil
}

// KeepListenKey 2. Продлить ListenKey (Keep-alive a ListenKey).
func (s *SpotList) KeepListenKey(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/userDataStream"
	s.log.Debug("KeepListenKey", "path", casePath, "params", params)
	resp, err := s.client.PrivatePut(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка KeepListenKey", "error", err)
		return nil, err
	}
	return resp, nil
}

// CloseListenKey 3. Закрыть ListenKey ().
func (s *SpotList) CloseListenKey(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/userDataStream"
	s.log.Debug("CloseListenKey", "path", casePath, "params", params)
	resp, err := s.client.PrivateDelete(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка CloseListenKey", "error", err)
		return nil, err
	}
	return resp, nil
}

// ## 5. Партнёрские и реферальные данные

// RebateHistory 1. История реферальных вознаграждений (Get Rebate History Records).
func (s *SpotList) RebateHistory(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/rebate/taxQuery"
	s.log.Debug("RebateHistory", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка RebateHistory", "error", err)
		return nil, err
	}
	return resp, nil
}

// RebateDetail 2. Детали реферальных выплат (Get Rebate Records Detail).
func (s *SpotList) RebateDetail(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/rebate/detail"
	s.log.Debug("RebateDetail", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка RebateDetail", "error", err)
		return nil, err
	}
	return resp, nil
}

// SelfRecordsDetail 3. Детали собственных выплат (Get Self Rebate Records Detail).
func (s *SpotList) SelfRecordsDetail(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/rebate/detail/kickback"
	s.log.Debug("SelfRecordsDetail", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка SelfRecordsDetail", "error", err)
		return nil, err
	}
	return resp, nil
}

// ReferCode 4. Получить код приглашения (Query ReferCode).
func (s *SpotList) ReferCode(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/rebate/referCode"
	s.log.Debug("ReferCode", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка ReferCode", "error", err)
		return nil, err
	}
	return resp, nil
}

// AffiliateCommission 5. Комиссии аффилиата (Get Affiliate Commission Record (affiliate only)).
func (s *SpotList) AffiliateCommission(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/rebate/affiliate/commission"
	s.log.Debug("AffiliateCommission", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка AffiliateCommission", "error", err)
		return nil, err
	}
	return resp, nil
}

// AffiliateWithdraw 6. История выводов аффилиата (Get Affiliate Withdraw Record (affiliate only)).
func (s *SpotList) AffiliateWithdraw(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/rebate/affiliate/withdraw"
	s.log.Debug("AffiliateWithdraw", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка AffiliateWithdraw", "error", err)
		return nil, err
	}
	return resp, nil
}

// AffiliateCommissionDetail 7. Детали комиссий аффилиата (Get Affiliate Commission Detail Record (affiliate only)).
func (s *SpotList) AffiliateCommissionDetail(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/rebate/affiliate/commission/detail"
	s.log.Debug("AffiliateCommissionDetail", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка AffiliateCommissionDetail", "error", err)
		return nil, err
	}
	return resp, nil
}

// AffiliateReferral 8. Сводка реферальных данных аффилиата (Get Affiliate Referral Data（affiliate only）).
func (s *SpotList) AffiliateReferral(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/rebate/affiliate/referral"
	s.log.Debug("AffiliateReferral", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка AffiliateReferral", "error", err)
		return nil, err
	}
	return resp, nil
}

// Subaffiliates 9. Суб‑аффилиаты (Get Subaffiliates Data (affiliate only)).
func (s *SpotList) Subaffiliates(ctx context.Context, params map[string]string) (*resty.Response, error) {
	casePath := "/rebate/affiliate/subaffiliates"
	s.log.Debug("Subaffiliates", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка Subaffiliates", "error", err)
		return nil, err
	}
	return resp, nil
}
