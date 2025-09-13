package mexc

import (
	"context"

	"github.com/go-resty/resty/v2"
)

type ListenKeyClient struct {
	log    Logger
	client APIClient
}

func NewListenKeyClient(log Logger, client APIClient) *ListenKeyClient {
	return &ListenKeyClient{log: log, client: client}
}

// ## 4. WS ListenKey

// CreateListenKey 1. Создать ListenKey (Listen Key  Create a ListenKey).
func (s *ListenKeyClient) CreateListenKey(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/userDataStream" //nolint:goconst
	s.log.Debug("CreateListenKey", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка CreateListenKey", "error", err)
		return nil, err
	}
	return resp, nil
}

// KeepListenKey 2. Продлить ListenKey (Keep-alive a ListenKey).
func (s *ListenKeyClient) KeepListenKey(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/userDataStream"
	s.log.Debug("KeepListenKey", "path", casePath, "params", params)
	resp, err := s.client.PrivatePut(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка KeepListenKey", "error", err)
		return nil, err
	}
	return resp, nil
}

// CloseListenKey 3. Закрыть ListenKey ().
func (s *ListenKeyClient) CloseListenKey(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/userDataStream"
	s.log.Debug("CloseListenKey", "path", casePath, "params", params)
	resp, err := s.client.PrivateDelete(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка CloseListenKey", "error", err)
		return nil, err
	}
	return resp, nil
}
