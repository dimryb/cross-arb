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

// SwapRequest представляет тело запроса для эндпоинта /swap.
// Создает транзакцию для обмена на основе котировки.
type SwapRequest struct {
	QuoteResponse                 QuoteResponse              `json:"quoteResponse"`
	UserPublicKey                 string                     `json:"userPublicKey"`
	Payer                         *string                    `json:"payer,omitempty"`
	WrapAndUnwrapSol              *bool                      `json:"wrapAndUnwrapSol,omitempty"`
	UseSharedAccounts             *bool                      `json:"useSharedAccounts,omitempty"`
	FeeAccount                    *string                    `json:"feeAccount,omitempty"`
	TrackingAccount               *string                    `json:"trackingAccount,omitempty"`
	PrioritizationFeeLamports     *PrioritizationFeeLamports `json:"prioritizationFeeLamports,omitempty"`
	AsLegacyTransaction           *bool                      `json:"asLegacyTransaction,omitempty"`
	DestinationTokenAccount       *string                    `json:"destinationTokenAccount,omitempty"`
	DynamicComputeUnitLimit       *bool                      `json:"dynamicComputeUnitLimit,omitempty"`
	SkipUserAccountsRPCCalls      *bool                      `json:"skipUserAccountsRpcCalls,omitempty"`
	DynamicSlippage               *bool                      `json:"dynamicSlippage,omitempty"`
	ComputeUnitPriceMicroLamports *uint64                    `json:"computeUnitPriceMicroLamports,omitempty"`
	BlockhashSlotsToExpiry        *int                       `json:"blockhashSlotsToExpiry,omitempty"`
}

// PrioritizationFeeLamports определяет структурированную комиссию для приоритизации транзакции.
// Предоставляет более гибкие возможности по сравнению с ComputeUnitPriceMicroLamports.
type PrioritizationFeeLamports struct {
	*PriorityLevelWithMaxLamports `json:"priorityLevelWithMaxLamports,omitempty"`
	JitoTipLamports               *uint64 `json:"jitoTipLamports,omitempty"`
}

// PriorityLevelWithMaxLamports задает автоматический уровень приоритета с максимальным порогом стоимости.
type PriorityLevelWithMaxLamports struct {
	PriorityLevel PriorityLevel `json:"priorityLevel"`
	MaxLamports   uint64        `json:"maxLamports"`
}

// PriorityLevel определяет предустановленные уровни приоритета транзакции.
type PriorityLevel string

const (
	PriorityLevelMedium   PriorityLevel = "medium"
	PriorityLevelHigh     PriorityLevel = "high"
	PriorityLevelVeryHigh PriorityLevel = "veryHigh"
)

// SwapResponse представляет ответ от эндпоинта /swap.
type SwapResponse struct {
	SwapTransaction           string `json:"swapTransaction"`
	LastValidBlockHeight      uint64 `json:"lastValidBlockHeight"`
	PrioritizationFeeLamports uint64 `json:"prioritizationFeeLamports"`
}
