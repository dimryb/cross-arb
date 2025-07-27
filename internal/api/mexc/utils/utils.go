package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	i "github.com/dimryb/cross-arb/internal/interface"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	APIKey    string
	SecretKey string
	Logger    i.Logger
}

func NewClient(apiKey, secretKey string, logger i.Logger) *Client {
	return &Client{
		APIKey:    apiKey,
		SecretKey: secretKey,
		Logger:    logger,
	}
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

	c.Logger.Debug("Выполняется публичный GET-запрос", "url", path)

	// Создаём HTTP-клиент
	client := resty.New()

	// Выполняем запрос
	resp, err := client.R().Get(path)
	if err != nil {
		c.Logger.Error("Ошибка при выполнении GET-запроса", "error", err)
		return nil, err
	}

	return resp, nil
}

// PrivateGet выполняет авторизованный GET-запрос.
func (c *Client) PrivateGet(urlStr string, jsonParams string) (*resty.Response, error) { //nolint: dupl
	var path string
	timestamp := time.Now().UnixNano() / 1e6

	c.Logger.Debug("Текущее время (timestamp)", "timestamp", timestamp)

	if jsonParams == "" {
		message := fmt.Sprintf("timestamp=%d", timestamp)
		sign := ComputeHmac256(message, c.SecretKey)
		path = fmt.Sprintf("%s?timestamp=%d&signature=%s", urlStr, timestamp, sign)
		c.Logger.Debug("Подпись запроса", "message", message, "signature", sign)
	} else {
		strParams := JSONToParamStr(jsonParams)
		message := fmt.Sprintf("%s&timestamp=%d", strParams, timestamp)
		sign := ComputeHmac256(message, c.SecretKey)
		path = fmt.Sprintf("%s?%s&timestamp=%d&signature=%s", urlStr, strParams, timestamp, sign)
		c.Logger.Debug("Подпись запроса", "message", message, "signature", sign)
	}

	c.Logger.Debug("Выполняется приватный GET-запрос", "url", path)

	client := resty.New()
	resp, err := client.R().SetHeaders(map[string]string{
		"X-MEXC-APIKEY": c.APIKey,
		"Content-Type":  "application/json",
	}).Get(path)
	if err != nil {
		c.Logger.Error("Ошибка при приватном GET-запросе", "error", err)
		return nil, err
	}

	return resp, nil
}

// PrivatePost выполняет авторизованный POST-запрос.
func (c *Client) PrivatePost(urlStr string, jsonParams string) (*resty.Response, error) { //nolint: dupl
	var path string
	timestamp := time.Now().UnixNano() / 1e6

	c.Logger.Debug("Текущее время (timestamp)", "timestamp", timestamp)

	if jsonParams == "" {
		message := fmt.Sprintf("timestamp=%d", timestamp)
		sign := ComputeHmac256(message, c.SecretKey)
		path = fmt.Sprintf("%s?timestamp=%d&signature=%s", urlStr, timestamp, sign)
		c.Logger.Debug("Подпись запроса", "message", message, "signature", sign)
	} else {
		strParams := JSONToParamStr(jsonParams)
		message := fmt.Sprintf("%s&timestamp=%d", strParams, timestamp)
		sign := ComputeHmac256(message, c.SecretKey)
		path = fmt.Sprintf("%s?%s&timestamp=%d&signature=%s", urlStr, strParams, timestamp, sign)
		c.Logger.Debug("Подпись запроса", "message", message, "signature", sign)
	}

	c.Logger.Debug("Выполняется приватный POST-запрос", "url", path)

	client := resty.New()
	resp, err := client.R().SetHeaders(map[string]string{
		"X-MEXC-APIKEY": c.APIKey,
		"Content-Type":  "application/json",
	}).Post(path)
	if err != nil {
		c.Logger.Error("Ошибка при приватном POST-запросе", "error", err)
		return nil, err
	}

	return resp, nil
}

// PrivateDelete выполняет авторизованный DELETE-запрос.
func (c *Client) PrivateDelete(urlStr string, jsonParams string) (*resty.Response, error) { //nolint: dupl
	var path string
	timestamp := time.Now().UnixNano() / 1e6

	c.Logger.Debug("Текущее время (timestamp)", "timestamp", timestamp)

	if jsonParams == "" {
		message := fmt.Sprintf("timestamp=%d", timestamp)
		sign := ComputeHmac256(message, c.SecretKey)
		path = fmt.Sprintf("%s?timestamp=%d&signature=%s", urlStr, timestamp, sign)
		c.Logger.Debug("Подпись запроса", "message", message, "signature", sign)
	} else {
		strParams := JSONToParamStr(jsonParams)
		message := fmt.Sprintf("%s&timestamp=%d", strParams, timestamp)
		sign := ComputeHmac256(message, c.SecretKey)
		path = fmt.Sprintf("%s?%s&timestamp=%d&signature=%s", urlStr, strParams, timestamp, sign)
		c.Logger.Debug("Подпись запроса", "message", message, "signature", sign)
	}

	c.Logger.Debug("Выполняется приватный DELETE-запрос", "url", path)

	client := resty.New()
	resp, err := client.R().SetHeaders(map[string]string{
		"X-MEXC-APIKEY": c.APIKey,
		"Content-Type":  "application/json",
	}).Delete(path)
	if err != nil {
		c.Logger.Error("Ошибка при приватном DELETE-запросе", "error", err)
		return nil, err
	}

	return resp, nil
}

// PrivatePut выполняет авторизованный PUT-запрос.
func (c *Client) PrivatePut(urlStr string, jsonParams string) (*resty.Response, error) { //nolint: dupl
	var path string
	timestamp := time.Now().UnixNano() / 1e6

	c.Logger.Debug("Текущее время (timestamp)", "timestamp", timestamp)

	if jsonParams == "" {
		message := fmt.Sprintf("timestamp=%d", timestamp)
		sign := ComputeHmac256(message, c.SecretKey)
		path = fmt.Sprintf("%s?timestamp=%d&signature=%s", urlStr, timestamp, sign)
		c.Logger.Debug("Подпись запроса", "message", message, "signature", sign)
	} else {
		strParams := JSONToParamStr(jsonParams)
		message := fmt.Sprintf("%s&timestamp=%d", strParams, timestamp)
		sign := ComputeHmac256(message, c.SecretKey)
		path = fmt.Sprintf("%s?%s&timestamp=%d&signature=%s", urlStr, strParams, timestamp, sign)
		c.Logger.Debug("Подпись запроса", "message", message, "signature", sign)
	}

	c.Logger.Debug("Выполняется приватный PUT-запрос", "url", path)

	client := resty.New()
	resp, err := client.R().SetHeaders(map[string]string{
		"X-MEXC-APIKEY": c.APIKey,
		"Content-Type":  "application/json",
	}).Put(path)
	if err != nil {
		c.Logger.Error("Ошибка при приватном PUT-запросе", "error", err)
		return nil, err
	}

	return resp, nil
}

// JSONToParamStr форматирует строку параметров из JSON.
func JSONToParamStr(jsonParams string) string {
	var paramsarr []string //nolint:prealloc
	m := make(map[string]string)
	err := json.Unmarshal([]byte(jsonParams), &m)
	if err != nil {
		return ""
	}
	for key, value := range m {
		paramsarr = append(paramsarr, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(paramsarr, "&")
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
