package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/dimryb/cross-arb/internal/api/mexc/config"
	"github.com/go-resty/resty/v2"
)

// PublicGet публичный GET-запрос.
func PublicGet(urlStr string, jsonParams string) *resty.Response {
	var path string
	if jsonParams == "" {
		path = urlStr
	} else {
		strParams := JSONToParamStr(jsonParams)
		path = urlStr + "?" + strParams
		// fmt.Println("Путь:", path)
	}
	// Создаём запрос
	client := resty.New()
	// Отправляем запрос
	resp, err := client.R().Get(path)
	if err != nil {
		log.Fatal("Ошибка запроса:", err)
	}

	// fmt.Println("Response Info:", resp)
	return resp
}

// 私有get请求.
func PrivateGet(urlStr string, jsonParams string) interface{} { //nolint: dupl
	var path string
	timestamp := time.Now().UnixNano() / 1e6
	fmt.Println(timestamp)
	if jsonParams == "" {
		message := fmt.Sprintf("timestamp=%d", timestamp)
		sign := ComputeHmac256(message, config.SEC_KEY)
		path = fmt.Sprintf("%s?timestamp=%d&signature=%s", urlStr, timestamp, sign)
		fmt.Println("message:", message)
		fmt.Println("sign:", sign)
		fmt.Println("path:", path)
	} else {
		strParams := JSONToParamStr(jsonParams)
		message := fmt.Sprintf("%s&timestamp=%d", strParams, timestamp)
		sign := ComputeHmac256(message, config.SEC_KEY)
		path = fmt.Sprintf("%s?%s&timestamp=%d&signature=%s", urlStr, strParams, timestamp, sign)
		fmt.Println("message:", ParamsEncode(message))
		fmt.Println("sign:", sign)
		fmt.Println("path:", path)
	}
	// 创建请求
	client := resty.New()
	// 发送请求
	resp, err := client.R().SetHeaders(map[string]string{
		"X-MEXC-APIKEY": config.API_KEY,
		"Content-Type":  "application/json",
	}).Get(path)
	if err != nil {
		log.Fatal("请求报错：", err)
	}

	// fmt.Println("Response Info:", resp)
	return resp
}

// 私有post请求.
func PrivatePost(urlStr string, jsonParams string) interface{} { //nolint: dupl
	var path string
	timestamp := time.Now().UnixNano() / 1e6
	fmt.Println(timestamp)
	if jsonParams == "" {
		message := fmt.Sprintf("timestamp=%d", timestamp)
		sign := ComputeHmac256(message, config.SEC_KEY)
		path = fmt.Sprintf("%s?timestamp=%d&signature=%s", urlStr, timestamp, sign)
		fmt.Println("message:", message)
		fmt.Println("sign:", sign)
		fmt.Println("path:", path)
	} else {
		strParams := JSONToParamStr(jsonParams)
		message := fmt.Sprintf("%s&timestamp=%d", strParams, timestamp)
		sign := ComputeHmac256(message, config.SEC_KEY)
		path = fmt.Sprintf("%s?%s&timestamp=%d&signature=%s", urlStr, strParams, timestamp, sign)
		fmt.Println("message:", ParamsEncode(message))
		fmt.Println("sign:", sign)
		fmt.Println("path:", path)
	}
	// 创建请求
	client := resty.New()
	// 发送请求
	resp, err := client.R().SetHeaders(map[string]string{
		"X-MEXC-APIKEY": config.API_KEY,
		"Content-Type":  "application/json",
	}).Post(path)
	if err != nil {
		log.Fatal("请求报错：", err)
	}

	// fmt.Println("Response Info:", resp)
	return resp
}

// 私有delete请求.
func PrivateDelete(urlStr string, jsonParams string) interface{} { //nolint: dupl
	var path string
	timestamp := time.Now().UnixNano() / 1e6
	fmt.Println(timestamp)
	if jsonParams == "" {
		message := fmt.Sprintf("timestamp=%d", timestamp)
		sign := ComputeHmac256(message, config.SEC_KEY)
		path = fmt.Sprintf("%s?timestamp=%d&signature=%s", urlStr, timestamp, sign)
		fmt.Println("message:", message)
		fmt.Println("sign:", sign)
		fmt.Println("path:", path)
	} else {
		strParams := JSONToParamStr(jsonParams)
		message := fmt.Sprintf("%s&timestamp=%d", strParams, timestamp)
		sign := ComputeHmac256(message, config.SEC_KEY)
		path = fmt.Sprintf("%s?%s&timestamp=%d&signature=%s", urlStr, strParams, timestamp, sign)
		fmt.Println("message:", ParamsEncode(message))
		fmt.Println("sign:", sign)
		fmt.Println("path:", path)
	}
	// 创建请求
	client := resty.New()
	// 发送请求
	resp, err := client.R().SetHeaders(map[string]string{
		"X-MEXC-APIKEY": config.API_KEY,
		"Content-Type":  "application/json",
	}).Delete(path)
	if err != nil {
		log.Fatal("请求报错：", err)
	}

	// fmt.Println("Response Info:", resp)
	return resp
}

// 私有put请求.
func PrivatePut(urlStr string, jsonParams string) interface{} { //nolint: dupl
	var path string
	timestamp := time.Now().UnixNano() / 1e6
	fmt.Println(timestamp)
	if jsonParams == "" {
		message := fmt.Sprintf("timestamp=%d", timestamp)
		sign := ComputeHmac256(message, config.SEC_KEY)
		path = fmt.Sprintf("%s?timestamp=%d&signature=%s", urlStr, timestamp, sign)
		fmt.Println("message:", message)
		fmt.Println("sign:", sign)
		fmt.Println("path:", path)
	} else {
		strParams := JSONToParamStr(jsonParams)
		message := fmt.Sprintf("%s&timestamp=%d", strParams, timestamp)
		sign := ComputeHmac256(message, config.SEC_KEY)
		path = fmt.Sprintf("%s?%s&timestamp=%d&signature=%s", urlStr, strParams, timestamp, sign)
		fmt.Println("message:", ParamsEncode(message))
		fmt.Println("sign:", sign)
		fmt.Println("path:", path)
	}
	// 创建请求
	client := resty.New()
	// 发送请求
	resp, err := client.R().SetHeaders(map[string]string{
		"X-MEXC-APIKEY": config.API_KEY,
		"Content-Type":  "application/json",
	}).Put(path)
	if err != nil {
		log.Fatal("请求报错：", err)
	}

	// fmt.Println("Response Info:", resp)
	return resp
}

// JSONToParamStr форматирует строку параметров из JSON.
func JSONToParamStr(jsonParams string) string {
	// Преобразуем параметры из JSON в строку параметров
	var paramsarr []string //nolint:prealloc
	var arritem string
	m := make(map[string]string)
	err := json.Unmarshal([]byte(jsonParams), &m)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Printf("map:%v\n", m)
	i := 0
	for key, value := range m {
		arritem = fmt.Sprintf("%s=%s", key, value)
		paramsarr = append(paramsarr, arritem)
		i++
		// fmt.Println("Итерация: ", i, "всего", len(m))
		if i > len(m) {
			break
		}
	}
	paramsstr := strings.Join(paramsarr, "&")
	// fmt.Println("Строка параметров:", paramsstr)
	return paramsstr
}

// urlencode.
func ParamsEncode(paramStr string) string {
	return url.QueryEscape(paramStr)
}

// 加密.
func ComputeHmac256(message string, secKey string) string {
	key := []byte(secKey)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}
