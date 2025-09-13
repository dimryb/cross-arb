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
	"strings"
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

// NewMexcClient создаёт новый клиент для MEXC Spot API.
func NewMexcClient(apiKey, secretKey, baseURL string, logger Logger) (*Client, error) {
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
func (c *Client) PublicGet(ctx context.Context, path string, params map[string]string) (*resty.Response, error) {
	fullURL := joinURLPath(c.baseURL, path)
	query := fullURL.Query()
	for k, v := range params {
		query.Set(k, v)
	}
	fullURL.RawQuery = query.Encode()
	finalURL := fullURL.String()

	c.logger.Debug("Выполняется публичный GET-запрос", "url", finalURL)

	client := resty.New().SetTimeout(defaultTimeout)
	resp, err := client.R().
		SetContext(ctx).
		Get(finalURL)
	if err != nil {
		c.logger.Error("Ошибка при выполнении GET-запроса", "error", err)
		return nil, err
	}

	return resp, nil
}

// PrivateGet выполняет авторизованный GET-запрос.
func (c *Client) PrivateGet(ctx context.Context, path string, params map[string]string) (*resty.Response, error) {
	return c.privateRequest(ctx, http.MethodGet, path, params)
}

// PrivatePost выполняет авторизованный POST-запрос.
func (c *Client) PrivatePost(ctx context.Context, path string, params map[string]string) (*resty.Response, error) {
	return c.privateRequest(ctx, http.MethodPost, path, params)
}

// PrivateDelete выполняет авторизованный DELETE-запрос.
func (c *Client) PrivateDelete(ctx context.Context, path string, params map[string]string) (*resty.Response, error) {
	return c.privateRequest(ctx, http.MethodDelete, path, params)
}

// PrivatePut выполняет авторизованный PUT-запрос.
func (c *Client) PrivatePut(ctx context.Context, path string, params map[string]string) (*resty.Response, error) {
	return c.privateRequest(ctx, http.MethodPut, path, params)
}

// privateRequest выполняет подписанный приватный запрос.
func (c *Client) privateRequest(
	ctx context.Context,
	method, path string,
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

	fullURL := joinURLPath(c.baseURL, path)
	fullURL.RawQuery = query.Encode()
	finalURL := fullURL.String()

	c.logger.Debug("Выполняется приватный запрос",
		"method", method,
		"url", finalURL,
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
		resp, err = req.Get(finalURL)
	case http.MethodPost:
		resp, err = req.Post(finalURL)
	case http.MethodDelete:
		resp, err = req.Delete(finalURL)
	case http.MethodPut:
		resp, err = req.Put(finalURL)
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		c.logger.Error("Ошибка при выполнении приватного запроса", "error", err, "method", method)
		return nil, err
	}

	return resp, nil
}

func joinURLPath(base *url.URL, path string) *url.URL {
	baseCopy := *base
	basePath := strings.TrimRight(base.Path, "/")
	relPath := strings.TrimLeft(path, "/")
	baseCopy.Path = basePath + "/" + relPath
	return &baseCopy
}

// ComputeHmac256 рассчитывает подпись HMAC SHA256.
func ComputeHmac256(message string, secKey string) string {
	key := []byte(secKey)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}
