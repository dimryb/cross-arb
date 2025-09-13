package mexc

import (
	"context"

	"github.com/go-resty/resty/v2"
)

type SpotTradeClient struct {
	log    Logger
	client APIClient
}

func NewSpotTradeClient(log Logger, client APIClient) *SpotTradeClient {
	return &SpotTradeClient{log: log, client: client}
}

// ## 2. Торговля (Spot Account & Trade)

// SelfSymbols 1. Список разрешённых символов пользователя (User API default symbol).
func (s *SpotTradeClient) SelfSymbols(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/selfSymbols"
	s.log.Debug("SelfSymbols", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка SelfSymbols", "error", err)
		return nil, err
	}
	return resp, nil
}

// TestOrder 2. Тестовый ордер (Test New Order).
func (s *SpotTradeClient) TestOrder(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/order/test"
	s.log.Debug("TestOrder", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка TestOrder", "error", err)
		return nil, err
	}
	return resp, nil
}

// PlaceOrder 3. Разместить ордер (New Order).
func (s *SpotTradeClient) PlaceOrder(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/order" //nolint:goconst
	s.log.Debug("PlaceOrder", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка PlaceOrder", "error", err)
		return nil, err
	}
	return resp, nil
}

// BatchOrder 4. Пакетное размещение ордеров (Batch Orders).
func (s *SpotTradeClient) BatchOrder(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/batchOrders"
	s.log.Debug("BatchOrder", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка BatchOrder", "error", err)
		return nil, err
	}
	return resp, nil
}

// CancelOrder 5. Отменить ордер (Cancel Order).
func (s *SpotTradeClient) CancelOrder(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/order"
	s.log.Debug("CancelOrder", "path", casePath, "params", params)
	resp, err := s.client.PrivateDelete(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка CancelOrder", "error", err)
		return nil, err
	}
	return resp, nil
}

// CancelAllOrders 6. Отменить все ордера по символу (Cancel all Open Orders on a Symbol).
func (s *SpotTradeClient) CancelAllOrders(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/openOrders"
	s.log.Debug("CancelAllOrders", "path", casePath, "params", params)
	resp, err := s.client.PrivateDelete(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка CancelAllOrders", "error", err)
		return nil, err
	}
	return resp, nil
}

// QueryOrder 7. Информация об ордере (Query Order).
func (s *SpotTradeClient) QueryOrder(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/order"
	s.log.Debug("QueryOrder", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка QueryOrder", "error", err)
		return nil, err
	}
	return resp, nil
}

// OpenOrder 8. Открытые ордера (Current Open Orders).
func (s *SpotTradeClient) OpenOrder(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/openOrders"
	s.log.Debug("OpenOrder", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка OpenOrder", "error", err)
		return nil, err
	}
	return resp, nil
}

// AllOrders 9. Все ордера (All Orders).
func (s *SpotTradeClient) AllOrders(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/allOrders"
	s.log.Debug("AllOrders", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка AllOrders", "error", err)
		return nil, err
	}
	return resp, nil
}

// SpotAccountInfo 10. Информация об аккаунте (Account Information).
func (s *SpotTradeClient) SpotAccountInfo(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/account"
	s.log.Debug("SpotAccountInfo", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка SpotAccountInfo", "error", err)
		return nil, err
	}
	return resp, nil
}

// SpotMyTrade 11. История сделок (Account Trade List).
func (s *SpotTradeClient) SpotMyTrade(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/myTrades"
	s.log.Debug("SpotMyTrade", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка SpotMyTrade", "error", err)
		return nil, err
	}
	return resp, nil
}

// MxDeduct 12. Включить MX‑дедукцию (Enable MX Deduct).
func (s *SpotTradeClient) MxDeduct(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/mxDeduct/enable"
	s.log.Debug("MxDeduct", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка MxDeduct", "error", err)
		return nil, err
	}
	return resp, nil
}

// QueryMxDeduct 13. Статус MX‑дедукции (Query MX Deduct Status).
func (s *SpotTradeClient) QueryMxDeduct(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/mxDeduct/enable"
	s.log.Debug("QueryMxDeduct", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка QueryMxDeduct", "error", err)
		return nil, err
	}
	return resp, nil
}
