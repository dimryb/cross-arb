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

// ## 钱包接口 Wallet Endpoints

// ### 1 查询币种信息 Query the currency information.
func QueryCurrencyInfo(jsonParams string) interface{} {
	caseURL := "/capital/config/getall"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 2 提币 Withdraw.
func Withdraw(jsonParams string) interface{} {
	caseURL := "/capital/withdraw/apply"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivatePost(requestURL, jsonParams)
	return response
}

// ### 3 取消提币 Cancel withdraw.
func CancelWithdraw(jsonParams string) interface{} {
	caseURL := "/capital/withdraw"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateDelete(requestURL, jsonParams)
	return response
}

// ### 4 获取充值历史 Deposit History.
func DepositHistory(jsonParams string) interface{} {
	caseURL := "/capital/deposit/hisrec"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 5 获取提币历史 Withdraw History.
func WithdrawHistory(jsonParams string) interface{} {
	caseURL := "/capital/withdraw/historyl"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 6 生成充值地址 Generate deposit address.
func GenDepositAddress(jsonParams string) interface{} {
	caseURL := "/capital/deposit/address"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivatePost(requestURL, jsonParams)
	return response
}

// ### 7 获取充值地址 Deposit Address.
func DepositAddress(jsonParams string) interface{} {
	caseURL := "/capital/deposit/address"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 8 获取提币地址 Withdraw Address.
func WithdrawAddress(jsonParams string) interface{} {
	caseURL := "/capital/withdraw/address"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 9 用户万向划转 User Universal Transfer.
func Transfer(jsonParams string) interface{} {
	caseURL := "/capital/transfer"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivatePost(requestURL, jsonParams)
	return response
}

// ### 10 查询用户万向划转历史 Query User Universal Transfer History.
func TransferHistory(jsonParams string) interface{} {
	caseURL := "/capital/transfer"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 11 查询用户万向划转历史（根据tranId） Query User Universal Transfer History （by tranId）.
func TransferHistoryByID(jsonParams string) interface{} {
	caseURL := "/capital/transfer/tranId"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 12 获取小额资产可兑换列表 Get Assets That Can Be Converted Into MX.
func ConvertList(jsonParams string) interface{} {
	caseURL := "/capital/convert/list"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 13 小额资产兑换 Dust Transfer.
func Convert(jsonParams string) interface{} {
	caseURL := "/capital/convert"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivatePost(requestURL, jsonParams)
	return response
}

// ### 14 查询小额资产兑换历史 DustLog.
func ConvertHistory(jsonParams string) interface{} {
	caseURL := "/capital/convert"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 15 获取ETF基础信息 Get ETF info.
func ETFInfo(jsonParams string) interface{} {
	caseURL := "/etf/info"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 16 用户站内转账 Internal Transfer.
func InternalTransfer(jsonParams string) interface{} {
	caseURL := "/capital/transfer/internal"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivatePost(requestURL, jsonParams)
	return response
}

// ### 17 用户站内转账历史 Internal Transfer History.
func InternalTransferHistory(jsonParams string) interface{} {
	caseURL := "/capital/transfer/internal"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
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
