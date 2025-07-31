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
	// QuoteResponse - ответ, полученный от /quote. Обязательно.
	QuoteResponse QuoteResponse `json:"quoteResponse"`

	// UserPublicKey - публичный ключ кошелька пользователя. Обязательно.
	UserPublicKey string `json:"userPublicKey"`

	// Payer - публичный ключ кошелька для оплаты комиссий.
	// По умолчанию используется UserPublicKey.
	Payer *string `json:"payer,omitempty"`

	// WrapAndUnwrapSol - автоматически оборачивать/разворачивать SOL в WSOL при необходимости.
	// По умолчанию: true. Рекомендуется оставлять включенным для корректной работы с SOL.
	WrapAndUnwrapSol *bool `json:"wrapAndUnwrapSol,omitempty"`

	// UseSharedAccounts - использовать общие токен-аккаунты для экономии на комиссиях за создание аккаунтов.
	// По умолчанию: true. Помогает снизить стоимость транзакций.
	UseSharedAccounts *bool `json:"useSharedAccounts,omitempty"`

	// FeeAccount - токен-аккаунт для получения комиссии платформы (integrator fee).
	// Должен быть существующим и принадлежать указанному в feeAccountOwner.
	FeeAccount *string `json:"feeAccount,omitempty"`

	// TrackingAccount - публичный ключ для отслеживания источника трафика/рефералов.
	// Используется для аналитики и потенциальных программ партнерства.
	TrackingAccount *string `json:"trackingAccount,omitempty"`

	// PrioritizationFeeLamports - структурированная комиссия за приоритет транзакции.
	// Взаимоисключающий с ComputeUnitPriceMicroLamports. Используйте только один из них.
	PrioritizationFeeLamports *PrioritizationFeeLamports `json:"prioritizationFeeLamports,omitempty"`

	// AsLegacyTransaction - создать legacy-транзакцию вместо versioned transaction.
	// Должно совпадать со значением, переданным в запрос котировки /quote.
	AsLegacyTransaction *bool `json:"asLegacyTransaction,omitempty"`

	// DestinationTokenAccount - указать конкретный токен-аккаунт для получения выходных токенов.
	// Если не указан, будет использован или создан associated token account пользователя.
	DestinationTokenAccount *string `json:"destinationTokenAccount,omitempty"`

	// DynamicComputeUnitLimit - позволить Jupiter автоматически оптимизировать лимит Compute Units.
	// По умолчанию: true. Помогает избежать переплаты за неиспользованные compute units.
	DynamicComputeUnitLimit *bool `json:"dynamicComputeUnitLimit,omitempty"`

	// SkipUserAccountsRPCCalls - пропустить предварительные RPC-запросы для проверки аккаунтов пользователя.
	// Ускоряет создание транзакции, но может привести к ошибкам при выполнении.
	SkipUserAccountsRPCCalls *bool `json:"skipUserAccountsRpcCalls,omitempty"`

	// DynamicSlippage - использовать динамическое проскальзывание на основе текущих рыночных условий.
	// Может помочь при высокой волатильности, но делает результат менее предсказуемым.
	DynamicSlippage *bool `json:"dynamicSlippage,omitempty"`

	// ComputeUnitPriceMicroLamports - цена за единицу вычислений (CU) в микро-лампортах.
	// Взаимоисключающий с PrioritizationFeeLamports. Используйте только один из них.
	// 1 лампорт = 1,000,000 микро-лампортов. Чем выше значение, тем выше приоритет.
	ComputeUnitPriceMicroLamports *uint64 `json:"computeUnitPriceMicroLamports,omitempty"`

	// BlockhashSlotsToExpiry - количество слотов, в течение которых blockhash транзакции будет действителен.
	// Максимальное значение ограничено сетью Solana.
	// Меньшие значения повышают шансы на быстрое исполнение.
	BlockhashSlotsToExpiry *int `json:"blockhashSlotsToExpiry,omitempty"`
}

// PrioritizationFeeLamports определяет структурированную комиссию для приоритизации транзакции.
// Предоставляет более гибкие возможности по сравнению с ComputeUnitPriceMicroLamports.
type PrioritizationFeeLamports struct {
	// PriorityLevelWithMaxLamports - автоматический выбор приоритета с лимитом стоимости.
	*PriorityLevelWithMaxLamports `json:"priorityLevelWithMaxLamports,omitempty"`

	// JitoTipLamports - дополнительные "чаевые" валидаторам Jito для MEV-защиты и ускорения.
	// Рекомендуется для транзакций, требующих быстрого исполнения.
	JitoTipLamports *uint64 `json:"jitoTipLamports,omitempty"`
}

// PriorityLevelWithMaxLamports задает автоматический уровень приоритета с максимальным порогом стоимости.
type PriorityLevelWithMaxLamports struct {
	// PriorityLevel - предустановленный уровень приоритета.
	PriorityLevel PriorityLevel `json:"priorityLevel"`

	// MaxLamports - максимальная сумма в лампортах, которую готовы заплатить за приоритет.
	MaxLamports uint64 `json:"maxLamports"`
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
	// SwapTransaction - сериализованная транзакция в формате base64.
	// Готова для подписи и отправки в сеть Solana.
	SwapTransaction string `json:"swapTransaction"`

	// LastValidBlockHeight - номер последнего блока, в котором транзакция останется действительной.
	// После этого блока транзакция будет отклонена сетью.
	LastValidBlockHeight uint64 `json:"lastValidBlockHeight"`

	// PrioritizationFeeLamports - итоговая рассчитанная комиссия за приоритет в лампортах.
	// Включает все настройки приоритизации из запроса.
	PrioritizationFeeLamports uint64 `json:"prioritizationFeeLamports"`
}
