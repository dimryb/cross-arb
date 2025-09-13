package mexc

import (
	"context"

	"github.com/go-resty/resty/v2"
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
) (*resty.Response, error) {
	casePath := "/sub-account/virtualSubAccount"
	s.log.Debug("CreateSub", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка CreateSub", "error", err)
		return nil, err
	}
	return resp, nil
}

// QuerySub 2. Получить список суб‑аккаунтов (Query Sub-account List).
func (s *SpotSubAccountClient) QuerySub(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/sub-account/list"
	s.log.Debug("QuerySub", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка QuerySub", "error", err)
		return nil, err
	}
	return resp, nil
}

// CreateSubApikey 3. Создать API‑ключ для суб‑аккаунта (Create an apiKey for a sub-account).
func (s *SpotSubAccountClient) CreateSubApikey(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/sub-account/apiKey" //nolint:goconst
	s.log.Debug("CreateSubApikey", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка CreateSubApikey", "error", err)
		return nil, err
	}
	return resp, nil
}

// QuerySubApikey 4. Получить API‑ключи суб‑аккаунта (apiKey Query the apiKey of a sub-account).
func (s *SpotSubAccountClient) QuerySubApikey(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/sub-account/apiKey"
	s.log.Debug("QuerySubApikey", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка QuerySubApikey", "error", err)
		return nil, err
	}
	return resp, nil
}

// DeleteSubApikey 5. Удалить API‑ключ суб‑аккаунта (apiKey Delete the apiKey of a sub-account).
func (s *SpotSubAccountClient) DeleteSubApikey(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/sub-account/apiKey"
	s.log.Debug("DeleteSubApikey", "path", casePath, "params", params)
	resp, err := s.client.PrivateDelete(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка DeleteSubApikey", "error", err)
		return nil, err
	}
	return resp, nil
}

// UniTransfer 6. Универсальный перевод между аккаунтами (Universal Transfer).
func (s *SpotSubAccountClient) UniTransfer(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/capital/sub-account/universalTransfer"
	s.log.Debug("UniTransfer", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка UniTransfer", "error", err)
		return nil, err
	}
	return resp, nil
}

// QueryUniTransfer 7. История универсальных переводов (Query Universal Transfer History).
func (s *SpotSubAccountClient) QueryUniTransfer(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/capital/sub-account/universalTransfer"
	s.log.Debug("QueryUniTransfer", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка QueryUniTransfer", "error", err)
		return nil, err
	}
	return resp, nil
}
