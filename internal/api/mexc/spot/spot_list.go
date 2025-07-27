package spotlist

import (
	"fmt"

	"github.com/dimryb/cross-arb/internal/api/mexc/config"
	"github.com/dimryb/cross-arb/internal/api/mexc/utils"
	i "github.com/dimryb/cross-arb/internal/interface"
	"github.com/go-resty/resty/v2"
)

// SpotClient — клиент для работы с MEXC Spot API.
type SpotClient struct {
	log     i.Logger
	BaseURL string
	client  *utils.Client
}

// NewSpotClient создаёт новый клиент для Spot API.
func NewSpotClient(log i.Logger, baseURL string, client *utils.Client) *SpotClient {
	return &SpotClient{
		log:     log,
		BaseURL: baseURL,
		client:  client,
	}
}

// # Реализация API-запросов
// ## Эндпоинты для получения рыночных данных (Market Data Endpoints)

// Ping 1. Проверка подключения к серверу (Test Connectivity).
func (s *SpotClient) Ping(jsonParams string) (*resty.Response, error) {
	caseURL := "/ping"
	requestURL := s.BaseURL + caseURL
	s.log.Debug("Ping request to MEXC", "url", requestURL, "params", jsonParams)

	resp, err := s.client.PublicGet(requestURL, jsonParams)
	if err != nil {
		s.log.Error("Ошибка в Ping", "error", err)
		return nil, err
	}

	return resp, nil
}

// Time 2. Получить серверное время (Check Server Time).
func (s *SpotClient) Time(jsonParams string) (*resty.Response, error) {
	caseURL := "/time"
	requestURL := s.BaseURL + caseURL
	s.log.Debug("Time request to MEXC", "url", requestURL, "params", jsonParams)

	resp, err := s.client.PublicGet(requestURL, jsonParams)
	if err != nil {
		s.log.Error("Ошибка в Time", "error", err)
		return nil, err
	}

	return resp, nil
}

// APISymbol 3. Список торговых пар по умолчанию (API Default Symbol).
func (s *SpotClient) APISymbol(jsonParams string) (*resty.Response, error) {
	caseURL := "/defaultSymbols"
	requestURL := s.BaseURL + caseURL
	s.log.Debug("API symbol request to MEXC", "url", requestURL, "params", jsonParams)

	resp, err := s.client.PublicGet(requestURL, jsonParams)
	if err != nil {
		s.log.Error("Ошибка в APISymbol", "error", err)
		return nil, err
	}

	return resp, nil
}

// ExchangeInfo 4. Информация о торгах (Exchange Information).
func (s *SpotClient) ExchangeInfo(jsonParams string) (*resty.Response, error) {
	caseURL := "/exchangeInfo"
	requestURL := s.BaseURL + caseURL
	s.log.Debug("Exchange info request to MEXC", "url", requestURL, "params", jsonParams)

	resp, err := s.client.PublicGet(requestURL, jsonParams)
	if err != nil {
		s.log.Error("Ошибка в ExchangeInfo", "error", err)
		return nil, err
	}

	return resp, nil
}

// Depth 5. Глубина стакана (Depth).
func (s *SpotClient) Depth(jsonParams string) (*resty.Response, error) {
	caseURL := "/depth"
	requestURL := s.BaseURL + caseURL
	s.log.Debug("Order book depth request to MEXC", "url", requestURL, "params", jsonParams)

	resp, err := s.client.PublicGet(requestURL, jsonParams)
	if err != nil {
		s.log.Error("Ошибка в Depth", "error", err)
		return nil, err
	}

	return resp, nil
}

// Trades 6. Список последних сделок (Recent Trades List).
func (s *SpotClient) Trades(jsonParams string) (*resty.Response, error) {
	caseURL := "/trades"
	requestURL := s.BaseURL + caseURL
	s.log.Debug("Recent trades request to MEXC", "url", requestURL, "params", jsonParams)

	resp, err := s.client.PublicGet(requestURL, jsonParams)
	if err != nil {
		s.log.Error("Ошибка в Trades", "error", err)
		return nil, err
	}

	return resp, nil
}

// AggTrades 7. Агрегированный список сделок (Aggregate Trades List).
func (s *SpotClient) AggTrades(jsonParams string) (*resty.Response, error) {
	caseURL := "/aggTrades"
	requestURL := s.BaseURL + caseURL
	s.log.Debug("Aggregate trades request to MEXC", "url", requestURL, "params", jsonParams)

	resp, err := s.client.PublicGet(requestURL, jsonParams)
	if err != nil {
		s.log.Error("Ошибка в AggTrades", "error", err)
		return nil, err
	}

	return resp, nil
}

// Kline 8. Данные свечей (K-line Data).
func (s *SpotClient) Kline(jsonParams string) (*resty.Response, error) {
	caseURL := "/klines"
	requestURL := s.BaseURL + caseURL
	s.log.Debug("K-line data request to MEXC", "url", requestURL, "params", jsonParams)

	resp, err := s.client.PublicGet(requestURL, jsonParams)
	if err != nil {
		s.log.Error("Ошибка в Kline", "error", err)
		return nil, err
	}

	return resp, nil
}

// AvgPrice 9. Средняя цена за период (Current Average Price).
func (s *SpotClient) AvgPrice(jsonParams string) (*resty.Response, error) {
	caseURL := "/avgPrice"
	requestURL := s.BaseURL + caseURL
	s.log.Debug("Average price request to MEXC", "url", requestURL, "params", jsonParams)

	resp, err := s.client.PublicGet(requestURL, jsonParams)
	if err != nil {
		s.log.Error("Ошибка в AvgPrice", "error", err)
		return nil, err
	}

	return resp, nil
}

// Ticker24hr 10. Статистика изменения цены за 24 часа (24hr Ticker Price Change Statistics).
func (s *SpotClient) Ticker24hr(jsonParams string) (*resty.Response, error) {
	caseURL := "/ticker/24hr"
	requestURL := s.BaseURL + caseURL
	s.log.Debug("24hr ticker stats request to MEXC", "url", requestURL, "params", jsonParams)

	resp, err := s.client.PublicGet(requestURL, jsonParams)
	if err != nil {
		s.log.Error("Ошибка в Ticker24hr", "error", err)
		return nil, err
	}

	return resp, nil
}

// Price 11. Текущая цена символа (Symbol Price Ticker).
func (s *SpotClient) Price(jsonParams string) (*resty.Response, error) {
	caseURL := "/ticker/price"
	requestURL := s.BaseURL + caseURL
	s.log.Debug("Symbol price request to MEXC", "url", requestURL, "params", jsonParams)

	resp, err := s.client.PublicGet(requestURL, jsonParams)
	if err != nil {
		s.log.Error("Ошибка в Price", "error", err)
		return nil, err
	}

	return resp, nil
}

// BookTicker 12. Лучшие цены в стакане (Symbol Order Book Ticker).
func (s *SpotClient) BookTicker(jsonParams string) (*resty.Response, error) {
	caseURL := "/ticker/bookTicker"
	requestURL := s.BaseURL + caseURL
	s.log.Debug("Order book ticker request to MEXC", "url", requestURL, "params", jsonParams)

	resp, err := s.client.PublicGet(requestURL, jsonParams)
	if err != nil {
		s.log.Error("Ошибка в BookTicker", "error", err)
		return nil, err
	}

	return resp, nil
}

// Суб‑аккаунты (Sub‑Account Endpoints).

// CreateSub создаёт виртуальный суб‑аккаунт.
func (s *SpotClient) CreateSub(jsonParams string) (*resty.Response, error) {
	caseURL := "/sub-account/virtualSubAccount"
	url := s.BaseURL + caseURL
	s.log.Debug("CreateSub", "url", url, "params", jsonParams)

	resp, err := s.client.PrivatePost(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка CreateSub", "error", err)
		return nil, err
	}
	return resp, nil
}

// QuerySub возвращает список суб‑аккаунтов.
func (s *SpotClient) QuerySub(jsonParams string) (*resty.Response, error) {
	caseURL := "/sub-account/list"
	url := s.BaseURL + caseURL
	s.log.Debug("QuerySub", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка QuerySub", "error", err)
		return nil, err
	}
	return resp, nil
}

// CreateSubApikey создаёт API‑ключ для суб‑аккаунта.
func (s *SpotClient) CreateSubApikey(jsonParams string) (*resty.Response, error) {
	caseURL := "/sub-account/apiKey"
	url := s.BaseURL + caseURL
	s.log.Debug("CreateSubApikey", "url", url, "params", jsonParams)

	resp, err := s.client.PrivatePost(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка CreateSubApikey", "error", err)
		return nil, err
	}
	return resp, nil
}

// QuerySubApikey возвращает API‑ключи суб‑аккаунта.
func (s *SpotClient) QuerySubApikey(jsonParams string) (*resty.Response, error) {
	caseURL := "/sub-account/apiKey"
	url := s.BaseURL + caseURL
	s.log.Debug("QuerySubApikey", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка QuerySubApikey", "error", err)
		return nil, err
	}
	return resp, nil
}

// DeleteSubApikey удаляет API‑ключ суб‑аккаунта.
func (s *SpotClient) DeleteSubApikey(jsonParams string) (*resty.Response, error) {
	caseURL := "/sub-account/apiKey"
	url := s.BaseURL + caseURL
	s.log.Debug("DeleteSubApikey", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateDelete(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка DeleteSubApikey", "error", err)
		return nil, err
	}
	return resp, nil
}

// UniTransfer выполняет межаккаунтный перевод.
func (s *SpotClient) UniTransfer(jsonParams string) (*resty.Response, error) {
	caseURL := "/capital/sub-account/universalTransfer"
	url := s.BaseURL + caseURL
	s.log.Debug("UniTransfer", "url", url, "params", jsonParams)

	resp, err := s.client.PrivatePost(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка UniTransfer", "error", err)
		return nil, err
	}
	return resp, nil
}

// QueryUniTransfer возвращает историю универсальных переводов.
func (s *SpotClient) QueryUniTransfer(jsonParams string) (*resty.Response, error) {
	caseURL := "/capital/sub-account/universalTransfer"
	url := s.BaseURL + caseURL
	s.log.Debug("QueryUniTransfer", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка QueryUniTransfer", "error", err)
		return nil, err
	}
	return resp, nil
}

// ## Трейдинг (Spot Account & Trade)

// SelfSymbols возвращает список разрешённых торговых пар пользователя.
func (s *SpotClient) SelfSymbols(jsonParams string) (*resty.Response, error) {
	caseURL := "/selfSymbols"
	url := s.BaseURL + caseURL
	s.log.Debug("SelfSymbols", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка SelfSymbols", "error", err)
		return nil, err
	}
	return resp, nil
}

// TestOrder выполняет тестовый ордер (без размещения в книге).
func (s *SpotClient) TestOrder(jsonParams string) (*resty.Response, error) {
	caseURL := "/order/test"
	url := s.BaseURL + caseURL
	s.log.Debug("TestOrder", "url", url, "params", jsonParams)

	resp, err := s.client.PrivatePost(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка TestOrder", "error", err)
		return nil, err
	}
	return resp, nil
}

// PlaceOrder размещает новый ордер.
func (s *SpotClient) PlaceOrder(jsonParams string) (*resty.Response, error) {
	caseURL := "/order"
	url := s.BaseURL + caseURL
	s.log.Debug("PlaceOrder", "url", url, "params", jsonParams)

	resp, err := s.client.PrivatePost(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка PlaceOrder", "error", err)
		return nil, err
	}
	return resp, nil
}

// BatchOrder размещает пакет ордеров.
func (s *SpotClient) BatchOrder(jsonParams string) (*resty.Response, error) {
	caseURL := "/batchOrders"
	url := s.BaseURL + caseURL
	s.log.Debug("BatchOrder", "url", url, "params", jsonParams)

	resp, err := s.client.PrivatePost(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка BatchOrder", "error", err)
		return nil, err
	}
	return resp, nil
}

// CancelOrder отменяет один ордер.
func (s *SpotClient) CancelOrder(jsonParams string) (*resty.Response, error) {
	caseURL := "/order"
	url := s.BaseURL + caseURL
	s.log.Debug("CancelOrder", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateDelete(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка CancelOrder", "error", err)
		return nil, err
	}
	return resp, nil
}

// CancelAllOrders отменяет все ордера по символу.
func (s *SpotClient) CancelAllOrders(jsonParams string) (*resty.Response, error) {
	caseURL := "/openOrders"
	url := s.BaseURL + caseURL
	s.log.Debug("CancelAllOrders", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateDelete(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка CancelAllOrders", "error", err)
		return nil, err
	}
	return resp, nil
}

// QueryOrder возвращает информацию об ордере.
func (s *SpotClient) QueryOrder(jsonParams string) (*resty.Response, error) {
	caseURL := "/order"
	url := s.BaseURL + caseURL
	s.log.Debug("QueryOrder", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка QueryOrder", "error", err)
		return nil, err
	}
	return resp, nil
}

// OpenOrder возвращает текущие открытые ордера.
func (s *SpotClient) OpenOrder(jsonParams string) (*resty.Response, error) {
	caseURL := "/openOrders"
	url := s.BaseURL + caseURL
	s.log.Debug("OpenOrder", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка OpenOrder", "error", err)
		return nil, err
	}
	return resp, nil
}

// AllOrders возвращает историю всех ордеров.
func (s *SpotClient) AllOrders(jsonParams string) (*resty.Response, error) {
	caseURL := "/allOrders"
	url := s.BaseURL + caseURL
	s.log.Debug("AllOrders", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка AllOrders", "error", err)
		return nil, err
	}
	return resp, nil
}

// SpotAccountInfo возвращает информацию об аккаунте.
func (s *SpotClient) SpotAccountInfo(jsonParams string) (*resty.Response, error) {
	caseURL := "/account"
	url := s.BaseURL + caseURL
	s.log.Debug("SpotAccountInfo", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка SpotAccountInfo", "error", err)
		return nil, err
	}
	return resp, nil
}

// SpotMyTrade возвращает историю сделок.
func (s *SpotClient) SpotMyTrade(jsonParams string) (*resty.Response, error) {
	caseURL := "/myTrades"
	url := s.BaseURL + caseURL
	s.log.Debug("SpotMyTrade", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка SpotMyTrade", "error", err)
		return nil, err
	}
	return resp, nil
}

// MxDeduct включает использование MX для оплаты комиссий.
func (s *SpotClient) MxDeduct(jsonParams string) (*resty.Response, error) {
	caseURL := "/mxDeduct/enable"
	url := s.BaseURL + caseURL
	s.log.Debug("MxDeduct", "url", url, "params", jsonParams)

	resp, err := s.client.PrivatePost(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка MxDeduct", "error", err)
		return nil, err
	}
	return resp, nil
}

// QueryMxDeduct возвращает статус MX‑дедукции.
func (s *SpotClient) QueryMxDeduct(jsonParams string) (*resty.Response, error) {
	caseURL := "/mxDeduct/enable"
	url := s.BaseURL + caseURL
	s.log.Debug("QueryMxDeduct", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка QueryMxDeduct", "error", err)
		return nil, err
	}
	return resp, nil
}

// ## Кошелёк (Wallet Endpoints)

// QueryCurrencyInfo возвращает информацию о валюте.
func (s *SpotClient) QueryCurrencyInfo(jsonParams string) (*resty.Response, error) {
	caseURL := "/capital/config/getall"
	url := s.BaseURL + caseURL
	s.log.Debug("QueryCurrencyInfo", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка QueryCurrencyInfo", "error", err)
		return nil, err
	}
	return resp, nil
}

// Withdraw выполняет вывод средств.
func (s *SpotClient) Withdraw(jsonParams string) (*resty.Response, error) {
	caseURL := "/capital/withdraw/apply"
	url := s.BaseURL + caseURL
	s.log.Debug("Withdraw", "url", url, "params", jsonParams)

	resp, err := s.client.PrivatePost(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка Withdraw", "error", err)
		return nil, err
	}
	return resp, nil
}

// CancelWithdraw отменяет заявку на вывод.
func (s *SpotClient) CancelWithdraw(jsonParams string) (*resty.Response, error) {
	caseURL := "/capital/withdraw"
	url := s.BaseURL + caseURL
	s.log.Debug("CancelWithdraw", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateDelete(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка CancelWithdraw", "error", err)
		return nil, err
	}
	return resp, nil
}

// DepositHistory возвращает историю депозитов.
func (s *SpotClient) DepositHistory(jsonParams string) (*resty.Response, error) {
	caseURL := "/capital/deposit/hisrec"
	url := s.BaseURL + caseURL
	s.log.Debug("DepositHistory", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка DepositHistory", "error", err)
		return nil, err
	}
	return resp, nil
}

// WithdrawHistory возвращает историю выводов.
func (s *SpotClient) WithdrawHistory(jsonParams string) (*resty.Response, error) {
	caseURL := "/capital/withdraw/historyl"
	url := s.BaseURL + caseURL
	s.log.Debug("WithdrawHistory", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка WithdrawHistory", "error", err)
		return nil, err
	}
	return resp, nil
}

// GenDepositAddress генерирует адрес депозита.
func (s *SpotClient) GenDepositAddress(jsonParams string) (*resty.Response, error) {
	caseURL := "/capital/deposit/address"
	url := s.BaseURL + caseURL
	s.log.Debug("GenDepositAddress", "url", url, "params", jsonParams)

	resp, err := s.client.PrivatePost(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка GenDepositAddress", "error", err)
		return nil, err
	}
	return resp, nil
}

// DepositAddress возвращает адрес депозита.
func (s *SpotClient) DepositAddress(jsonParams string) (*resty.Response, error) {
	caseURL := "/capital/deposit/address"
	url := s.BaseURL + caseURL
	s.log.Debug("DepositAddress", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка DepositAddress", "error", err)
		return nil, err
	}
	return resp, nil
}

// WithdrawAddress возвращает адрес вывода.
func (s *SpotClient) WithdrawAddress(jsonParams string) (*resty.Response, error) {
	caseURL := "/capital/withdraw/address"
	url := s.BaseURL + caseURL
	s.log.Debug("WithdrawAddress", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка WithdrawAddress", "error", err)
		return nil, err
	}
	return resp, nil
}

// Transfer выполняет универсальный перевод.
func (s *SpotClient) Transfer(jsonParams string) (*resty.Response, error) {
	caseURL := "/capital/transfer"
	url := s.BaseURL + caseURL
	s.log.Debug("Transfer", "url", url, "params", jsonParams)

	resp, err := s.client.PrivatePost(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка Transfer", "error", err)
		return nil, err
	}
	return resp, nil
}

// TransferHistory возвращает историю переводов.
func (s *SpotClient) TransferHistory(jsonParams string) (*resty.Response, error) {
	caseURL := "/capital/transfer"
	url := s.BaseURL + caseURL
	s.log.Debug("TransferHistory", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка TransferHistory", "error", err)
		return nil, err
	}
	return resp, nil
}

// TransferHistoryByID возвращает перевод по tranId.
func (s *SpotClient) TransferHistoryByID(jsonParams string) (*resty.Response, error) {
	caseURL := "/capital/transfer/tranId"
	url := s.BaseURL + caseURL
	s.log.Debug("TransferHistoryByID", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка TransferHistoryByID", "error", err)
		return nil, err
	}
	return resp, nil
}

// ConvertList возвращает список активов для конвертации.
func (s *SpotClient) ConvertList(jsonParams string) (*resty.Response, error) {
	caseURL := "/capital/convert/list"
	url := s.BaseURL + caseURL
	s.log.Debug("ConvertList", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка ConvertList", "error", err)
		return nil, err
	}
	return resp, nil
}

// Convert выполняет конвертацию мелких активов.
func (s *SpotClient) Convert(jsonParams string) (*resty.Response, error) {
	caseURL := "/capital/convert"
	url := s.BaseURL + caseURL
	s.log.Debug("Convert", "url", url, "params", jsonParams)

	resp, err := s.client.PrivatePost(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка Convert", "error", err)
		return nil, err
	}
	return resp, nil
}

// ConvertHistory возвращает историю конвертаций.
func (s *SpotClient) ConvertHistory(jsonParams string) (*resty.Response, error) {
	caseURL := "/capital/convert"
	url := s.BaseURL + caseURL
	s.log.Debug("ConvertHistory", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка ConvertHistory", "error", err)
		return nil, err
	}
	return resp, nil
}

// ETFInfo возвращает информацию об ETF.
func (s *SpotClient) ETFInfo(jsonParams string) (*resty.Response, error) {
	caseURL := "/etf/info"
	url := s.BaseURL + caseURL
	s.log.Debug("ETFInfo", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка ETFInfo", "error", err)
		return nil, err
	}
	return resp, nil
}

// InternalTransfer выполняет внутренний перевод между пользователями.
func (s *SpotClient) InternalTransfer(jsonParams string) (*resty.Response, error) {
	caseURL := "/capital/transfer/internal"
	url := s.BaseURL + caseURL
	s.log.Debug("InternalTransfer", "url", url, "params", jsonParams)

	resp, err := s.client.PrivatePost(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка InternalTransfer", "error", err)
		return nil, err
	}
	return resp, nil
}

// InternalTransferHistory возвращает историю внутренних переводов.
func (s *SpotClient) InternalTransferHistory(jsonParams string) (*resty.Response, error) {
	caseURL := "/capital/transfer/internal"
	url := s.BaseURL + caseURL
	s.log.Debug("InternalTransferHistory", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка InternalTransferHistory", "error", err)
		return nil, err
	}
	return resp, nil
}

// ## WS ListenKey

// ### 1 生成 Listen Key  Create a ListenKey.
func CreateListenKey(jsonParams string) interface{} {
	caseURL := "/userDataStream" //nolint:goconst
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivatePost(requestURL, jsonParams)
	return response
}

// ### 2 延长 Listen Key 有效期  Keep-alive a ListenKey.
func KeepListenKey(jsonParams string) interface{} {
	caseURL := "/userDataStream"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivatePut(requestURL, jsonParams)
	return response
}

// ### 3 关闭 Listen Key  Close a ListenKey.
func CloseListenKey(jsonParams string) interface{} {
	caseURL := "/userDataStream"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateDelete(requestURL, jsonParams)
	return response
}

// ## 邀请返佣接口

// ### 1 获取邀请返佣记录 Get Rebate History Records.
func RebateHistory(jsonParams string) interface{} {
	caseURL := "/rebate/taxQuery"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 2 获取返佣记录明细 Get Rebate Records Detail.
func RebateDetail(jsonParams string) interface{} {
	caseURL := "/rebate/detail"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 3 获取自返记录明细 Get Self Rebate Records Detail.
func SelfRecordsDetail(jsonParams string) interface{} {
	caseURL := "/rebate/detail/kickback"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 4 获取邀请人 Query ReferCode.
func ReferCode(jsonParams string) interface{} {
	caseURL := "/rebate/referCode"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 5 获取代理邀请返佣记录 （代理账户）Get Affiliate Commission Record (affiliate only).
func AffiliateCommission(jsonParams string) interface{} {
	caseURL := "/rebate/affiliate/commission"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 6 获取代理提现记录 （代理账户）Get Affiliate Withdraw Record (affiliate only).
func AffiliateWithdraw(jsonParams string) interface{} {
	caseURL := "/rebate/affiliate/withdraw"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 7 获取代理返佣明细 （代理账户）Get Affiliate Commission Detail Record (affiliate only).
func AffiliateCommissionDetail(jsonParams string) interface{} {
	caseURL := "/rebate/affiliate/commission/detail"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 8 查询直客页面数据 （代理账户）Get Affiliate Referral Data（affiliate only）.
func AffiliateReferral(jsonParams string) interface{} {
	caseURL := "/rebate/affiliate/referral"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 9 查询子代理页面数据 （代理账户）Get Subaffiliates Data (affiliate only).
func Subaffiliates(jsonParams string) interface{} {
	caseURL := "/rebate/affiliate/subaffiliates"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}
