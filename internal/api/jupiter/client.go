package jupiter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	i "github.com/dimryb/cross-arb/internal/interface"
)

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	logger     i.Logger
}

// NewJupiterClient создает новый клиент для Jupiter Swap API.
// Подробнее: https://dev.jup.ag/docs/api/swap-api/quote
func NewJupiterClient(logger i.Logger, baseURL string) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	return &Client{
		baseURL:    parsedURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		logger:     logger,
	}, nil
}

// QuoteOptions содержит опциональные параметры для запроса котировки.
type QuoteOptions struct {
	SlippageBps                *int  `json:"slippageBps,omitempty"`
	RestrictIntermediateTokens *bool `json:"restrictIntermediateTokens,omitempty"`
	OnlyDirectRoutes           *bool `json:"onlyDirectRoutes,omitempty"`
	AsLegacyTransaction        *bool `json:"asLegacyTransaction,omitempty"`
	PlatformFeeBps             *int  `json:"platformFeeBps,omitempty"`
	MaxAccounts                *int  `json:"maxAccounts,omitempty"`
}

// DefaultQuoteOptions возвращает опции по умолчанию.
func DefaultQuoteOptions() *QuoteOptions {
	return &QuoteOptions{}
}

// Quote получает котировку для обмена токенов.
func (c *Client) Quote(
	ctx context.Context,
	inputMint, outputMint string,
	amount int64,
	opts *QuoteOptions,
) (*QuoteResponse, error) {
	start := time.Now()

	c.logger.Debug("Начало запроса котировки",
		"от_токена", inputMint,
		"к_токену", outputMint,
		"сумма", amount,
	)

	// Валидация входных параметров
	if inputMint == "" {
		return nil, fmt.Errorf("ошибка валидации: inputMint не может быть пустым")
	}
	if outputMint == "" {
		return nil, fmt.Errorf("ошибка валидации: outputMint не может быть пустым")
	}
	if amount <= 0 {
		return nil, fmt.Errorf("ошибка валидации: сумма должна быть положительной, получено: %v", amount)
	}

	if opts == nil {
		opts = DefaultQuoteOptions()
	}

	requestURL := c.buildQuoteURL(inputMint, outputMint, amount, opts)
	c.logger.Debug("Построен URL запроса", "url", requestURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать http-запрос для котировки: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "jupiter-go-client/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить http-запрос для котировки: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.logger.Error("Ошибка при закрытии тела ответа", "ошибка", closeErr)
		}
	}()

	c.logger.Debug("Получен ответ от Jupiter API",
		"статус_код", resp.StatusCode,
		"время_мс", time.Since(start).Milliseconds(),
	)

	quoteResponse, err := c.handleQuoteResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("не удалось обработать ответ с котировкой: %w", err)
	}

	c.logger.Debug("Котировка успешно получена",
		"обмен", fmt.Sprintf("%s → %s", inputMint, outputMint),
		"входная_сумма", amount,
		"выходная_сумма", quoteResponse.OutAmount,
		"время_мс", time.Since(start).Milliseconds(),
	)

	return quoteResponse, nil
}

// buildQuoteURL строит URL для запроса котировки.
func (c *Client) buildQuoteURL(inputMint, outputMint string, amount int64, opts *QuoteOptions) string {
	requestURL := *c.baseURL
	requestURL.Path += "/quote"

	q := url.Values{}
	q.Add("inputMint", inputMint)
	q.Add("outputMint", outputMint)
	q.Add("amount", strconv.FormatInt(amount, 10))

	if opts.SlippageBps != nil {
		q.Add("slippageBps", strconv.Itoa(*opts.SlippageBps))
	}
	if opts.RestrictIntermediateTokens != nil {
		q.Add("restrictIntermediateTokens", strconv.FormatBool(*opts.RestrictIntermediateTokens))
	}
	if opts.OnlyDirectRoutes != nil {
		q.Add("onlyDirectRoutes", strconv.FormatBool(*opts.OnlyDirectRoutes))
	}
	if opts.AsLegacyTransaction != nil {
		q.Add("asLegacyTransaction", strconv.FormatBool(*opts.AsLegacyTransaction))
	}
	if opts.PlatformFeeBps != nil {
		q.Add("platformFeeBps", strconv.Itoa(*opts.PlatformFeeBps))
	}
	if opts.MaxAccounts != nil {
		q.Add("maxAccounts", strconv.Itoa(*opts.MaxAccounts))
	}

	requestURL.RawQuery = q.Encode()
	return requestURL.String()
}

// handleQuoteResponse обрабатывает ответ от API.
func (c *Client) handleQuoteResponse(resp *http.Response) (*QuoteResponse, error) {
	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		//c.logger.Warn("Получен некорректный статус от Jupiter API",
		//	"статус_код", resp.StatusCode,
		//	"статус", resp.Status,
		//)

		var errorResp struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}

		// Пытаемся декодировать ошибку
		if err := decoder.Decode(&errorResp); err == nil && errorResp.Error != "" {
			c.logger.Error("Jupiter API вернул структурированную ошибку",
				"статус_код", resp.StatusCode,
				"api_ошибка", errorResp.Error,
				"api_сообщение", errorResp.Message,
			)
			return nil, fmt.Errorf("API error (status %d): %s - %s", resp.StatusCode, errorResp.Error, errorResp.Message)
		}

		//c.logger.Error("Jupiter API вернул неструктурированную ошибку", "статус_код", resp.StatusCode)
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var quoteResponse QuoteResponse
	if err := decoder.Decode(&quoteResponse); err != nil {
		c.logger.Error("Ошибка декодирования JSON-ответа от Jupiter API", "ошибка", err)
		return nil, fmt.Errorf("failed to decode quote response: %w", err)
	}

	return &quoteResponse, nil
}

// Swap создает транзакцию для обмена на основе полученной котировки.
// Подробнее: https://dev.jup.ag/docs/api/swap-api/swap
func (c *Client) Swap(ctx context.Context, swapReq *SwapRequest) (*SwapResponse, error) {
	start := time.Now()
	c.logger.Debug("Начало запроса на обмен", "пользователь", swapReq.UserPublicKey)

	if swapReq.UserPublicKey == "" {
		return nil, fmt.Errorf("ошибка валидации: UserPublicKey не может быть пустым")
	}

	requestURL := *c.baseURL
	requestURL.Path += "/swap"

	body, err := json.Marshal(swapReq)
	if err != nil {
		return nil, fmt.Errorf("не удалось сериализовать тело запроса для обмена: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("не удалось создать http-запрос для обмена: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "jupiter-go-client/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить http-запрос для обмена: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.logger.Error("Ошибка при закрытии тела ответа", "ошибка", closeErr)
		}
	}()

	c.logger.Debug("Получен ответ от Jupiter Swap API",
		"статус_код", resp.StatusCode,
		"время_мс", time.Since(start).Milliseconds(),
	)

	swapResponse, err := c.handleSwapResponse(resp)
	if err != nil {
		return nil, err
	}

	c.logger.Info("Транзакция для обмена успешно создана",
		"пользователь", swapReq.UserPublicKey,
		"время_мс", time.Since(start).Milliseconds(),
	)

	return swapResponse, nil
}

// handleSwapResponse обрабатывает HTTP-ответ для /swap.
func (c *Client) handleSwapResponse(resp *http.Response) (*SwapResponse, error) {
	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("ошибка от Jupiter API: статус %d, не удалось прочитать тело ответа: %w",
				resp.StatusCode, err)
		}
		return nil, fmt.Errorf("ошибка от Jupiter API: статус %d, тело: %s", resp.StatusCode, string(bodyBytes))
	}

	var swapResponse SwapResponse
	if err := json.NewDecoder(resp.Body).Decode(&swapResponse); err != nil {
		return nil, fmt.Errorf("не удалось десериализовать ответ от /swap: %w", err)
	}

	return &swapResponse, nil
}
