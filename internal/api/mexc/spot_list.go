package mexc

// SpotAPI — клиент для работы с MEXC Spot API.
type SpotAPI struct {
	Market    *SpotMarketClient
	Sub       *SpotSubAccountClient
	Wallet    *SpotWalletClient
	ListenKey *ListenKeyClient
	Rebate    *SpotRebateClient
	Trade     *SpotTradeClient
}

// NewSpotAPI создаёт новый клиент для Spot API.
func NewSpotAPI(log Logger, client *Client) *SpotAPI {
	return &SpotAPI{
		Market:    NewSpotMarketClient(log, client),
		Sub:       NewSpotSubAccountClient(log, client),
		Wallet:    NewSpotWalletClient(log, client),
		ListenKey: NewListenKeyClient(log, client),
		Rebate:    NewSpotRebateClient(log, client),
		Trade:     NewSpotTradeClient(log, client),
	}
}
