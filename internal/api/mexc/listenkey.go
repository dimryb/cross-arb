package mexc

import (
	"context"
	"encoding/json"
)

// DTOs для управления ListenKey (WebSocket подписки).

type CreateListenKeyResponse struct {
	// TODO: определить поля ответа POST /userDataStream
}

type KeepListenKeyResponse struct {
	// TODO: определить поля ответа PUT /userDataStream
}

type CloseListenKeyResponse struct {
	// TODO: определить поля ответа DELETE /userDataStream
}

type ListenKeyClient struct {
	log    Logger
	client APIClient
}

func NewListenKeyClient(log Logger, client APIClient) *ListenKeyClient {
	return &ListenKeyClient{log: log, client: client}
}

const casePath = "/userDataStream"

// ## 4. WS ListenKey

// CreateListenKey 1. Создать ListenKey (Listen Key  Create a ListenKey).
func (s *ListenKeyClient) CreateListenKey(
	ctx context.Context,
	params map[string]string,
) (*CreateListenKeyResponse, error) {
	s.log.Debug("CreateListenKey request to MEXC", "path", casePath, "params", params)

	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в CreateListenKey", "error", err)
		return nil, err
	}

	var result CreateListenKeyResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal CreateListenKey response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// KeepListenKey 2. Продлить ListenKey (Keep-alive a ListenKey).
func (s *ListenKeyClient) KeepListenKey(
	ctx context.Context,
	params map[string]string,
) (*KeepListenKeyResponse, error) {
	s.log.Debug("KeepListenKey request to MEXC", "path", casePath, "params", params)

	resp, err := s.client.PrivatePut(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в KeepListenKey", "error", err)
		return nil, err
	}

	var result KeepListenKeyResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal KeepListenKey response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// CloseListenKey 3. Закрыть ListenKey ().
func (s *ListenKeyClient) CloseListenKey(
	ctx context.Context,
	params map[string]string,
) (*CloseListenKeyResponse, error) {
	s.log.Debug("CloseListenKey request to MEXC", "path", casePath, "params", params)

	resp, err := s.client.PrivateDelete(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка в CloseListenKey", "error", err)
		return nil, err
	}

	var result CloseListenKeyResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal CloseListenKey response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}
