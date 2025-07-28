package jupiter

// QuoteResponse представляет ответ от эндпоинта /quote.
type QuoteResponse struct {
	InputMint            string      `json:"inputMint"`
	InAmount             string      `json:"inAmount"`
	OutputMint           string      `json:"outputMint"`
	OutAmount            string      `json:"outAmount"`
	OtherAmountThreshold string      `json:"otherAmountThreshold"`
	SwapMode             string      `json:"swapMode"`
	SlippageBps          int         `json:"slippageBps"`
	PlatformFee          interface{} `json:"platformFee"`
	PriceImpactPct       string      `json:"priceImpactPct"`
	RoutePlan            []RoutePlan `json:"routePlan"`
	ContextSlot          int         `json:"contextSlot"`
	TimeTaken            float64     `json:"timeTaken"`
}

// RoutePlan представляет единичный маршрут в обмене.
type RoutePlan struct {
	SwapInfo SwapInfo `json:"swapInfo"`
	Percent  int      `json:"percent"`
}

// SwapInfo содержит детали об обмене.
type SwapInfo struct {
	AmmKey     string `json:"ammKey"`
	Label      string `json:"label"`
	InputMint  string `json:"inputMint"`
	OutputMint string `json:"outputMint"`
	InAmount   string `json:"inAmount"`
	OutAmount  string `json:"outAmount"`
	FeeAmount  string `json:"feeAmount"`
	FeeMint    string `json:"feeMint"`
}

// SwapRequest представляет запрос для эндпоинта /swap.
type SwapRequest struct {
	UserPublicKey                 string        `json:"userPublicKey"`
	WrapAndUnwrapSol              bool          `json:"wrapAndUnwrapSol"`
	UseSharedAccounts             bool          `json:"useSharedAccounts"`
	FeeAccount                    string        `json:"feeAccount"`
	ComputeUnitPriceMicroLamports interface{}   `json:"computeUnitPriceMicroLamports"`
	PrioritizationFeeLamports     interface{}   `json:"prioritizationFeeLamports"`
	QuoteResponse                 QuoteResponse `json:"quoteResponse"`
}

// SwapResponse представляет ответ от эндпоинта /swap.
type SwapResponse struct {
	SwapTransaction      string `json:"swapTransaction"`
	LastValidBlockHeight int    `json:"lastValidBlockHeight"`
}
