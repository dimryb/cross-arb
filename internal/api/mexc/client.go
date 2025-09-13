package mexc

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
func (c *Client) PublicGet(urlStr string, jsonParams string) (*resty.Response, error) {
	var path string
	if jsonParams == "" {
		path = urlStr
	} else {
		strParams := JSONToParamStr(jsonParams)
		path = urlStr + "?" + strParams
	}

	c.logger.Debug("Выполняется публичный GET-запрос", "url", path)

	// Создаём HTTP-клиент
	client := resty.New()

	// Выполняем запрос
	resp, err := client.R().Get(path)
	if err != nil {
		c.logger.Error("Ошибка при выполнении GET-запроса", "error", err)
		return nil, err
	}

	return resp, nil
}

// PrivateGet выполняет авторизованный GET-запрос.
func (c *Client) PrivateGet(urlStr string, jsonParams string) (*resty.Response, error) { //nolint: dupl
	var path string
	timestamp := time.Now().UnixNano() / 1e6

	c.logger.Debug("Текущее время (timestamp)", "timestamp", timestamp)

	if jsonParams == "" {
		message := fmt.Sprintf("timestamp=%d", timestamp)
		sign := ComputeHmac256(message, c.secretKey)
		path = fmt.Sprintf("%s?timestamp=%d&signature=%s", urlStr, timestamp, sign)
		c.logger.Debug("Подпись запроса", "message", message, "signature", sign)
	} else {
		strParams := JSONToParamStr(jsonParams)
		message := fmt.Sprintf("%s&timestamp=%d", strParams, timestamp)
		sign := ComputeHmac256(message, c.secretKey)
		path = fmt.Sprintf("%s?%s&timestamp=%d&signature=%s", urlStr, strParams, timestamp, sign)
		c.logger.Debug("Подпись запроса", "message", message, "signature", sign)
	}

	c.logger.Debug("Выполняется приватный GET-запрос", "url", path)

	client := resty.New()
	resp, err := client.R().SetHeaders(map[string]string{
		"X-MEXC-APIKEY": c.apiKey,
		"Content-Type":  "application/json",
	}).Get(path)
	if err != nil {
		c.logger.Error("Ошибка при приватном GET-запросе", "error", err)
		return nil, err
	}

	return resp, nil
}

// PrivatePost выполняет авторизованный POST-запрос.
func (c *Client) PrivatePost(urlStr string, jsonParams string) (*resty.Response, error) { //nolint: dupl
	var path string
	timestamp := time.Now().UnixNano() / 1e6

	c.logger.Debug("Текущее время (timestamp)", "timestamp", timestamp)

	if jsonParams == "" {
		message := fmt.Sprintf("timestamp=%d", timestamp)
		sign := ComputeHmac256(message, c.secretKey)
		path = fmt.Sprintf("%s?timestamp=%d&signature=%s", urlStr, timestamp, sign)
		c.logger.Debug("Подпись запроса", "message", message, "signature", sign)
	} else {
		strParams := JSONToParamStr(jsonParams)
		message := fmt.Sprintf("%s&timestamp=%d", strParams, timestamp)
		sign := ComputeHmac256(message, c.secretKey)
		path = fmt.Sprintf("%s?%s&timestamp=%d&signature=%s", urlStr, strParams, timestamp, sign)
		c.logger.Debug("Подпись запроса", "message", message, "signature", sign)
	}

	c.logger.Debug("Выполняется приватный POST-запрос", "url", path)

	client := resty.New()
	resp, err := client.R().SetHeaders(map[string]string{
		"X-MEXC-APIKEY": c.apiKey,
		"Content-Type":  "application/json",
	}).Post(path)
	if err != nil {
		c.logger.Error("Ошибка при приватном POST-запросе", "error", err)
		return nil, err
	}

	return resp, nil
}

// PrivateDelete выполняет авторизованный DELETE-запрос.
func (c *Client) PrivateDelete(urlStr string, jsonParams string) (*resty.Response, error) { //nolint: dupl
	var path string
	timestamp := time.Now().UnixNano() / 1e6

	c.logger.Debug("Текущее время (timestamp)", "timestamp", timestamp)

	if jsonParams == "" {
		message := fmt.Sprintf("timestamp=%d", timestamp)
		sign := ComputeHmac256(message, c.secretKey)
		path = fmt.Sprintf("%s?timestamp=%d&signature=%s", urlStr, timestamp, sign)
		c.logger.Debug("Подпись запроса", "message", message, "signature", sign)
	} else {
		strParams := JSONToParamStr(jsonParams)
		message := fmt.Sprintf("%s&timestamp=%d", strParams, timestamp)
		sign := ComputeHmac256(message, c.secretKey)
		path = fmt.Sprintf("%s?%s&timestamp=%d&signature=%s", urlStr, strParams, timestamp, sign)
		c.logger.Debug("Подпись запроса", "message", message, "signature", sign)
	}

	c.logger.Debug("Выполняется приватный DELETE-запрос", "url", path)

	client := resty.New()
	resp, err := client.R().SetHeaders(map[string]string{
		"X-MEXC-APIKEY": c.apiKey,
		"Content-Type":  "application/json",
	}).Delete(path)
	if err != nil {
		c.logger.Error("Ошибка при приватном DELETE-запросе", "error", err)
		return nil, err
	}

	return resp, nil
}

// PrivatePut выполняет авторизованный PUT-запрос.
func (c *Client) PrivatePut(urlStr string, jsonParams string) (*resty.Response, error) { //nolint: dupl
	var path string
	timestamp := time.Now().UnixNano() / 1e6

	c.logger.Debug("Текущее время (timestamp)", "timestamp", timestamp)

	if jsonParams == "" {
		message := fmt.Sprintf("timestamp=%d", timestamp)
		sign := ComputeHmac256(message, c.secretKey)
		path = fmt.Sprintf("%s?timestamp=%d&signature=%s", urlStr, timestamp, sign)
		c.logger.Debug("Подпись запроса", "message", message, "signature", sign)
	} else {
		strParams := JSONToParamStr(jsonParams)
		message := fmt.Sprintf("%s&timestamp=%d", strParams, timestamp)
		sign := ComputeHmac256(message, c.secretKey)
		path = fmt.Sprintf("%s?%s&timestamp=%d&signature=%s", urlStr, strParams, timestamp, sign)
		c.logger.Debug("Подпись запроса", "message", message, "signature", sign)
	}

	c.logger.Debug("Выполняется приватный PUT-запрос", "url", path)

	client := resty.New()
	resp, err := client.R().SetHeaders(map[string]string{
		"X-MEXC-APIKEY": c.apiKey,
		"Content-Type":  "application/json",
	}).Put(path)
	if err != nil {
		c.logger.Error("Ошибка при приватном PUT-запросе", "error", err)
		return nil, err
	}

	return resp, nil
}

// JSONToParamStr форматирует строку параметров из JSON.
func JSONToParamStr(jsonParams string) string {
	m := make(map[string]string)
	err := json.Unmarshal([]byte(jsonParams), &m)
	if err != nil {
		return ""
	}
	params := make([]string, 0, len(m))
	for key, value := range m {
		params = append(params, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(params, "&")
}

// ParamsEncode кодирует строку как URL-параметры.
func ParamsEncode(paramStr string) string {
	return url.QueryEscape(paramStr)
}

// ComputeHmac256 рассчитывает подпись HMAC SHA256.
func ComputeHmac256(message string, secKey string) string {
	key := []byte(secKey)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}
