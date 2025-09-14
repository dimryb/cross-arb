package mexc

import (
	"context"
	"encoding/json"
)

// TODO: DTOs — будут заполнены позже.

type (
	CreateSubResponse        struct{}
	QuerySubResponse         struct{}
	CreateSubApikeyResponse  struct{}
	QuerySubApikeyResponse   struct{}
	DeleteSubApikeyResponse  struct{}
	UniTransferResponse      struct{}
	QueryUniTransferResponse struct{}
)

type SpotSubAccountClient struct {
	log    Logger
	client APIClient
}

func NewSpotSubAccountClient(log Logger, client APIClient) *SpotSubAccountClient {
	return &SpotSubAccountClient{log: log, client: client}
}

// ## 1. Суб‑аккаунты (Sub‑Account Endpoints)

// CreateSub 1. Создать суб‑аккаунт (Create a Sub-account).
func (s *SpotSubAccountClient) CreateSub(
	ctx context.Context,
	params map[string]string,
) (*CreateSubResponse, error) {
	const path = "/sub-account/virtualSubAccount"
	s.log.Debug("CreateSub request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivatePost(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в CreateSub", "error", err)
		return nil, err
	}

	var result CreateSubResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal CreateSub response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// QuerySub 2. Получить список суб‑аккаунтов (Query Sub-account List).
func (s *SpotSubAccountClient) QuerySub(
	ctx context.Context,
	params map[string]string,
) (*QuerySubResponse, error) {
	const path = "/sub-account/list"
	s.log.Debug("QuerySub request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в QuerySub", "error", err)
		return nil, err
	}

	var result QuerySubResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal QuerySub response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// CreateSubApikey 3. Создать API‑ключ для суб‑аккаунта (Create an apiKey for a sub-account).
func (s *SpotSubAccountClient) CreateSubApikey(
	ctx context.Context,
	params map[string]string,
) (*CreateSubApikeyResponse, error) {
	const path = "/sub-account/apiKey"
	s.log.Debug("CreateSubApikey request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivatePost(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в CreateSubApikey", "error", err)
		return nil, err
	}

	var result CreateSubApikeyResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal CreateSubApikey response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// QuerySubApikey 4. Получить API‑ключи суб‑аккаунта (apiKey Query the apiKey of a sub-account).
func (s *SpotSubAccountClient) QuerySubApikey(
	ctx context.Context,
	params map[string]string,
) (*QuerySubApikeyResponse, error) {
	const path = "/sub-account/apiKey"
	s.log.Debug("QuerySubApikey request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в QuerySubApikey", "error", err)
		return nil, err
	}

	var result QuerySubApikeyResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal QuerySubApikey response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// DeleteSubApikey 5. Удалить API‑ключ суб‑аккаунта (apiKey Delete the apiKey of a sub-account).
func (s *SpotSubAccountClient) DeleteSubApikey(
	ctx context.Context,
	params map[string]string,
) (*DeleteSubApikeyResponse, error) {
	const path = "/sub-account/apiKey"
	s.log.Debug("DeleteSubApikey request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateDelete(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в DeleteSubApikey", "error", err)
		return nil, err
	}

	var result DeleteSubApikeyResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal DeleteSubApikey response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// UniTransfer 6. Универсальный перевод между аккаунтами (Universal Transfer).
func (s *SpotSubAccountClient) UniTransfer(
	ctx context.Context,
	params map[string]string,
) (*UniTransferResponse, error) {
	const path = "/capital/sub-account/universalTransfer"
	s.log.Debug("UniTransfer request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivatePost(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в UniTransfer", "error", err)
		return nil, err
	}

	var result UniTransferResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal UniTransfer response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// QueryUniTransfer 7. История универсальных переводов (Query Universal Transfer History).
func (s *SpotSubAccountClient) QueryUniTransfer(
	ctx context.Context,
	params map[string]string,
) (*QueryUniTransferResponse, error) {
	const path = "/capital/sub-account/universalTransfer"
	s.log.Debug("QueryUniTransfer request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в QueryUniTransfer", "error", err)
		return nil, err
	}

	var result QueryUniTransferResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal QueryUniTransfer response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}
