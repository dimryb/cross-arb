package spotlist

import (
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

// ## 1. Суб‑аккаунты (Sub‑Account Endpoints)

// 1. Создать суб‑аккаунт.
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

// 2. Получить список суб‑аккаунтов.
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

// 3. Создать API‑ключ для суб‑аккаунта.
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

// 4. Получить API‑ключи суб‑аккаунта.
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

// 5. Удалить API‑ключ суб‑аккаунта.
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

// 6. Универсальный перевод между аккаунтами.
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

// 7. История универсальных переводов.
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

// ## 2. Торговля (Spot Account & Trade)

// 1. Список разрешённых символов пользователя.
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

// 2. Тестовый ордер.
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

// 3. Разместить ордер.
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

// 4. Пакетное размещение ордеров.
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

// 5. Отменить ордер.
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

// 6. Отменить все ордера по символу.
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

// 7. Информация об ордере.
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

// 8. Открытые ордера.
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

// 9. Все ордера.
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

// 10. Информация об аккаунте.
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

// 11. История сделок.
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

// 12. Включить MX‑дедукцию.
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

// 13. Статус MX‑дедукции.
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

// ## 3. Кошелёк (Wallet Endpoints)

// 1. Информация о валюте.
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

// 2. Вывод средств.
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

// 3. Отменить вывод.
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

// 4. История депозитов.
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

// 5. История выводов.
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

// 6. Сгенерировать адрес депозита.
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

// 7. Получить адрес депозита.
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

// 8. Получить адрес вывода.
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

// 9. Универсальный перевод.
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

// 10. История переводов.
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

// 11. История перевода по tranId.
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

// 12. Список активов для конвертации.
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

// 13. Конвертация мелких активов.
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

// 14. История конвертаций.
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

// 15. Информация об ETF.
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

// 16. Внутренний перевод.
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

// 17. История внутренних переводов.
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

// ## 4. WS ListenKey

// 1. Создать ListenKey.
func (s *SpotClient) CreateListenKey(jsonParams string) (*resty.Response, error) {
	caseURL := "/userDataStream"
	url := s.BaseURL + caseURL
	s.log.Debug("CreateListenKey", "url", url, "params", jsonParams)

	resp, err := s.client.PrivatePost(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка CreateListenKey", "error", err)
		return nil, err
	}
	return resp, nil
}

// 2. Продлить ListenKey.
func (s *SpotClient) KeepListenKey(jsonParams string) (*resty.Response, error) {
	caseURL := "/userDataStream"
	url := s.BaseURL + caseURL
	s.log.Debug("KeepListenKey", "url", url, "params", jsonParams)

	resp, err := s.client.PrivatePut(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка KeepListenKey", "error", err)
		return nil, err
	}
	return resp, nil
}

// 3. Закрыть ListenKey.
func (s *SpotClient) CloseListenKey(jsonParams string) (*resty.Response, error) {
	caseURL := "/userDataStream"
	url := s.BaseURL + caseURL
	s.log.Debug("CloseListenKey", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateDelete(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка CloseListenKey", "error", err)
		return nil, err
	}
	return resp, nil
}

// ## 5. Партнёрские и реферальные данные

// 1. История реферальных вознаграждений.
func (s *SpotClient) RebateHistory(jsonParams string) (*resty.Response, error) {
	caseURL := "/rebate/taxQuery"
	url := s.BaseURL + caseURL
	s.log.Debug("RebateHistory", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка RebateHistory", "error", err)
		return nil, err
	}
	return resp, nil
}

// 2. Детали реферальных выплат.
func (s *SpotClient) RebateDetail(jsonParams string) (*resty.Response, error) {
	caseURL := "/rebate/detail"
	url := s.BaseURL + caseURL
	s.log.Debug("RebateDetail", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка RebateDetail", "error", err)
		return nil, err
	}
	return resp, nil
}

// 3. Детали собственных выплат.
func (s *SpotClient) SelfRecordsDetail(jsonParams string) (*resty.Response, error) {
	caseURL := "/rebate/detail/kickback"
	url := s.BaseURL + caseURL
	s.log.Debug("SelfRecordsDetail", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка SelfRecordsDetail", "error", err)
		return nil, err
	}
	return resp, nil
}

// 4. Получить код приглашения.
func (s *SpotClient) ReferCode(jsonParams string) (*resty.Response, error) {
	caseURL := "/rebate/referCode"
	url := s.BaseURL + caseURL
	s.log.Debug("ReferCode", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка ReferCode", "error", err)
		return nil, err
	}
	return resp, nil
}

// 5. Комиссии аффилиата.
func (s *SpotClient) AffiliateCommission(jsonParams string) (*resty.Response, error) {
	caseURL := "/rebate/affiliate/commission"
	url := s.BaseURL + caseURL
	s.log.Debug("AffiliateCommission", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка AffiliateCommission", "error", err)
		return nil, err
	}
	return resp, nil
}

// 6. История выводов аффилиата.
func (s *SpotClient) AffiliateWithdraw(jsonParams string) (*resty.Response, error) {
	caseURL := "/rebate/affiliate/withdraw"
	url := s.BaseURL + caseURL
	s.log.Debug("AffiliateWithdraw", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка AffiliateWithdraw", "error", err)
		return nil, err
	}
	return resp, nil
}

// 7. Детали комиссий аффилиата.
func (s *SpotClient) AffiliateCommissionDetail(jsonParams string) (*resty.Response, error) {
	caseURL := "/rebate/affiliate/commission/detail"
	url := s.BaseURL + caseURL
	s.log.Debug("AffiliateCommissionDetail", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка AffiliateCommissionDetail", "error", err)
		return nil, err
	}
	return resp, nil
}

// 8. Сводка реферальных данных аффилиата.
func (s *SpotClient) AffiliateReferral(jsonParams string) (*resty.Response, error) {
	caseURL := "/rebate/affiliate/referral"
	url := s.BaseURL + caseURL
	s.log.Debug("AffiliateReferral", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка AffiliateReferral", "error", err)
		return nil, err
	}
	return resp, nil
}

// 9. Суб‑аффилиаты.
func (s *SpotClient) Subaffiliates(jsonParams string) (*resty.Response, error) {
	caseURL := "/rebate/affiliate/subaffiliates"
	url := s.BaseURL + caseURL
	s.log.Debug("Subaffiliates", "url", url, "params", jsonParams)

	resp, err := s.client.PrivateGet(url, jsonParams)
	if err != nil {
		s.log.Error("Ошибка Subaffiliates", "error", err)
		return nil, err
	}
	return resp, nil
}
