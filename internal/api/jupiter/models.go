package jupiter

// QuoteResponse представляет ответ от эндпоинта /quote.
type QuoteResponse struct {
	InputMint            string       `json:"inputMint"`
	InAmount             string       `json:"inAmount"`
	OutputMint           string       `json:"outputMint"`
	OutAmount            string       `json:"outAmount"`
	OtherAmountThreshold string       `json:"otherAmountThreshold"`
	SwapMode             SwapMode     `json:"swapMode"`
	SlippageBps          int          `json:"slippageBps"`
	PlatformFee          *PlatformFee `json:"platformFee"`
	PriceImpactPct       string       `json:"priceImpactPct"`
	RoutePlan            []RoutePlan  `json:"routePlan"`
	ContextSlot          uint64       `json:"contextSlot"`
	TimeTaken            float64      `json:"timeTaken"`
}

// SwapMode фиксируем как enum вместо свободной строки.
type SwapMode string

const (
	SwapModeExactIn  SwapMode = "ExactIn"
	SwapModeExactOut SwapMode = "ExactOut"
)

// PlatformFee представляет структуру комиссии платформы.
type PlatformFee struct {
	Amount string `json:"amountToExchange"`
	FeeBps string `json:"feeBps"`
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
