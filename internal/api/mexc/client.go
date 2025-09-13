package mexc

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

const defaultTimeout = 30 * time.Second

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	apiKey     string
	secretKey  string
	logger     Logger
}

// NewClient создаёт новый клиент для MEXC Spot API.
func NewClient(apiKey, secretKey, baseURL string, logger Logger) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL '%s': %w", baseURL, err)
	}

	return &Client{
		baseURL:   parsedURL,
		apiKey:    apiKey,
		secretKey: secretKey,
		logger:    logger,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}, nil
}

// PublicGet публичный GET-запрос.
func (c *Client) PublicGet(ctx context.Context, urlStr string, params map[string]string) (*resty.Response, error) {
	query := make(url.Values)
	for k, v := range params {
		query.Set(k, v)
	}
	fullURL := c.baseURL.String() + urlStr
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	c.logger.Debug("Выполняется публичный GET-запрос", "url", fullURL)

	client := resty.New().SetTimeout(defaultTimeout)
	resp, err := client.R().
		SetContext(ctx).
		Get(fullURL)
	if err != nil {
		c.logger.Error("Ошибка при выполнении GET-запроса", "error", err)
		return nil, err
	}

	return resp, nil
}

// PrivateGet выполняет авторизованный GET-запрос.
func (c *Client) PrivateGet(ctx context.Context, urlStr string, params map[string]string) (*resty.Response, error) {
	return c.privateRequest(ctx, http.MethodGet, urlStr, params)
}

// PrivatePost выполняет авторизованный POST-запрос.
func (c *Client) PrivatePost(ctx context.Context, urlStr string, params map[string]string) (*resty.Response, error) {
	return c.privateRequest(ctx, http.MethodPost, urlStr, params)
}

// PrivateDelete выполняет авторизованный DELETE-запрос.
func (c *Client) PrivateDelete(ctx context.Context, urlStr string, params map[string]string) (*resty.Response, error) {
	return c.privateRequest(ctx, http.MethodDelete, urlStr, params)
}

// PrivatePut выполняет авторизованный PUT-запрос.
func (c *Client) PrivatePut(ctx context.Context, urlStr string, params map[string]string) (*resty.Response, error) {
	return c.privateRequest(ctx, http.MethodPut, urlStr, params)
}

// privateRequest выполняет подписанный приватный запрос.
func (c *Client) privateRequest(
	ctx context.Context,
	method, urlStr string,
	params map[string]string,
) (*resty.Response, error) {
	timestamp := time.Now().UnixMilli() // MEXC использует миллисекунды

	// Формируем строку запроса
	query := make(url.Values)
	for k, v := range params {
		query.Set(k, v)
	}
	query.Set("timestamp", strconv.FormatInt(timestamp, 10))

	toSign := query.Encode()
	signature := ComputeHmac256(toSign, c.secretKey)
	query.Set("signature", signature)

	fullURL := c.baseURL.String() + urlStr + "?" + query.Encode()

	c.logger.Debug("Выполняется приватный запрос",
		"method", method,
		"url", fullURL,
		"params", params,
	)

	client := resty.New().SetTimeout(defaultTimeout)
	req := client.R().
		SetContext(ctx).
		SetHeaders(map[string]string{
			"X-MEXC-APIKEY": c.apiKey,
			"Content-Type":  "application/json",
		})

	var resp *resty.Response
	var err error

	switch method {
	case http.MethodGet:
		resp, err = req.Get(fullURL)
	case http.MethodPost:
		resp, err = req.Post(fullURL)
	case http.MethodDelete:
		resp, err = req.Delete(fullURL)
	case http.MethodPut:
		resp, err = req.Put(fullURL)
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		c.logger.Error("Ошибка при выполнении приватного запроса", "error", err, "method", method)
		return nil, err
	}

	return resp, nil
}

// ComputeHmac256 рассчитывает подпись HMAC SHA256.
func ComputeHmac256(message string, secKey string) string {
	key := []byte(secKey)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}
