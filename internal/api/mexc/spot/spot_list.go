package spotlist

import (
	"fmt"

	"github.com/dimryb/cross-arb/internal/api/mexc/config"
	"github.com/dimryb/cross-arb/internal/api/mexc/utils"
	i "github.com/dimryb/cross-arb/internal/interface"
)

type SpotClient struct {
	log     i.Logger
	BaseURL string
}

func NewSpotClient(log i.Logger, baseURL string) *SpotClient {
	return &SpotClient{
		log:     log,
		BaseURL: baseURL,
	}
}

// # Реализация API-запросов
// ## Эндпоинты для получения рыночных данных (Market Data Endpoints)

// Ping 1. Проверка подключения к серверу (Test Connectivity).
func Ping(jsonParams string) interface{} {
	caseURL := "/ping"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PublicGet(requestURL, jsonParams)
	return response
}

// Time ### 2. Получить серверное время (Check Server Time).
func Time(jsonParams string) interface{} {
	caseURL := "/time"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PublicGet(requestURL, jsonParams)
	return response
}

// ### 3 API交易对 API default symbol.
func APISymbol(jsonParams string) interface{} {
	caseURL := "/defaultSymbols"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PublicGet(requestURL, jsonParams)
	return response
}

// ### 4 交易规范信息 Exchange Information.
func ExchangeInfo(jsonParams string) interface{} {
	caseURL := "/exchangeInfo"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PublicGet(requestURL, jsonParams)
	return response
}

// ### 5 深度信息 Depth.
func Depth(jsonParams string) interface{} {
	caseURL := "/depth"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PublicGet(requestURL, jsonParams)
	return response
}

// ### 6 近期成交列表 Recent Trades List.
func Trades(jsonParams string) interface{} {
	caseURL := "/trades"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PublicGet(requestURL, jsonParams)
	return response
}

// ### 7 近期成交列表（归集） Aggregate Trades List.
func AggTrades(jsonParams string) interface{} {
	caseURL := "/aggTrades"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PublicGet(requestURL, jsonParams)
	return response
}

// ### 8 K线数据 K-line Data.
func Kline(jsonParams string) interface{} {
	caseURL := "/klines"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PublicGet(requestURL, jsonParams)
	return response
}

// ### 9 当前平均价格 Current Average Price.
func AvgPrice(jsonParams string) interface{} {
	caseURL := "/avgPrice"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PublicGet(requestURL, jsonParams)
	return response
}

// ### 10 24小时价格滚动情况 24hr Ticker Price Change Statistics.
func Ticker24hr(jsonParams string) interface{} {
	caseURL := "/ticker/24hr"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PublicGet(requestURL, jsonParams)
	return response
}

// ### 11 最新价格 Symbol Price Ticker.
func Price(jsonParams string) interface{} {
	caseURL := "/ticker/price"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PublicGet(requestURL, jsonParams)
	return response
}

// BookTicker ### 12. Текущие лучшие цены по инструменту (Symbol Order Book Ticker).
func (s *SpotClient) BookTicker(jsonParams string) interface{} {
	caseURL := "/ticker/bookTicker"
	requestURL := s.BaseURL + caseURL
	s.log.Debug("request", "URL", requestURL)
	response := utils.PublicGet(requestURL, jsonParams)
	return response
}

// ## 母子账户接口 Sub-Account Endpoints

// ### 1 创建子账户 Create a Sub-account(For Master Account).
func CreateSub(jsonParams string) interface{} {
	caseURL := "/sub-account/virtualSubAccount"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivatePost(requestURL, jsonParams)
	return response
}

// ### 2 查看子账户列表 Query Sub-account List (For Master Account).
func QuerySub(jsonParams string) interface{} {
	caseURL := "/sub-account/list"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 3 创建子账户的APIkey Create an APIKey for a sub-account (For Master Account).
func CreateSubApikey(jsonParams string) interface{} {
	caseURL := "/sub-account/apiKey" //nolint:goconst
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivatePost(requestURL, jsonParams)
	return response
}

// ### 4 查询子账户的APIKey Query the APIKey of a sub-account (For Master Account).
func QuerySubApikey(jsonParams string) interface{} {
	caseURL := "/sub-account/apiKey"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 5 删除子账户的APIKey Delete the APIKey of a sub-account (For Master Account).
func DeleteSubApikey(jsonParams string) interface{} {
	caseURL := "/sub-account/apiKey"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateDelete(requestURL, jsonParams)
	return response
}

// ### 6 母子用户万向划转 Universal Transfer (For Master Account).
func UniTransfer(jsonParams string) interface{} {
	caseURL := "/capital/sub-account/universalTransfer"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivatePost(requestURL, jsonParams)
	return response
}

// ### 7 查询母子万向划转历史 Query Universal Transfer History (For Master Account).
func QueryUniTransfer(jsonParams string) interface{} {
	caseURL := "/capital/sub-account/universalTransfer"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ## 现货账户和交易接口 Spot Account and Trade

// ### 1 用户API交易对 User API default symbol.
func SelfSymbols(jsonParams string) interface{} {
	caseURL := "/selfSymbols"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 2 测试下单 Test New Order.
func TestOrder(jsonParams string) interface{} {
	caseURL := "/order/test"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivatePost(requestURL, jsonParams)
	return response
}

// ### 3 下单 New Order.
func PlaceOrder(jsonParams string) interface{} {
	caseURL := "/order" //nolint:goconst
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivatePost(requestURL, jsonParams)
	return response
}

// ### 4 批量下单 Batch Orders.
func BatchOrder(jsonParams string) interface{} {
	caseURL := "/batchOrders"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivatePost(requestURL, jsonParams)
	return response
}

// ### 5 撤销订单 Cancel Order.
func CancelOrder(jsonParams string) interface{} {
	caseURL := "/order"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateDelete(requestURL, jsonParams)
	return response
}

// ### 6 撤销单一交易对所有订单 Cancel all Open Orders on a Symbol.
func CancelAllOrders(jsonParams string) interface{} {
	caseURL := "/openOrders"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateDelete(requestURL, jsonParams)
	return response
}

// ### 7 查询订单 Query Order.
func QueryOrder(jsonParams string) interface{} {
	caseURL := "/order"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 8 当前挂单 Current Open Orders.
func OpenOrder(jsonParams string) interface{} {
	caseURL := "/openOrders"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 9 查询所有订单 All Orders.
func AllOrders(jsonParams string) interface{} {
	caseURL := "/allOrders"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 10 账户信息 Account Information.
func SpotAccountInfo(jsonParams string) interface{} {
	caseURL := "/account"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 11 账户成交历史 Account Trade List.
func SpotmyTrade(jsonParams string) interface{} {
	caseURL := "/myTrades"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
}

// ### 12 开启MX抵扣 Enable MX Deduct.
func MxDeduct(jsonParams string) interface{} {
	caseURL := "/mxDeduct/enable"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivatePost(requestURL, jsonParams)
	return response
}

// ### 13 查看MX抵扣状态 Query MX Deduct Status.
func QueryMxDeduct(jsonParams string) interface{} {
	caseURL := "/mxDeduct/enable"
	requestURL := config.BASE_URL + caseURL
	fmt.Println("requestURL:", requestURL)
	response := utils.PrivateGet(requestURL, jsonParams)
	return response
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
