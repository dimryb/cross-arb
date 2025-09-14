package mexc

import (
	"context"
	"encoding/json"
)

// TODO: DTOs — будут заполнены позже.

type (
	SelfSymbolsResponse     struct{}
	TestOrderResponse       struct{}
	PlaceOrderResponse      struct{}
	BatchOrderResponse      struct{}
	CancelOrderResponse     struct{}
	CancelAllOrdersResponse struct{}
	QueryOrderResponse      struct{}
	OpenOrderResponse       struct{}
	AllOrdersResponse       struct{}
	SpotAccountInfoResponse struct{}
	SpotMyTradeResponse     struct{}
	MxDeductResponse        struct{}
	QueryMxDeductResponse   struct{}
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
) (*SelfSymbolsResponse, error) {
	const path = "/selfSymbols"
	s.log.Debug("SelfSymbols request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в SelfSymbols", "error", err)
		return nil, err
	}

	var result SelfSymbolsResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal SelfSymbols response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// TestOrder 2. Тестовый ордер (Test New Order).
func (s *SpotTradeClient) TestOrder(
	ctx context.Context,
	params map[string]string,
) (*TestOrderResponse, error) {
	const path = "/order/test"
	s.log.Debug("TestOrder request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivatePost(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в TestOrder", "error", err)
		return nil, err
	}

	var result TestOrderResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal TestOrder response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// PlaceOrder 3. Разместить ордер (New Order).
func (s *SpotTradeClient) PlaceOrder(
	ctx context.Context,
	params map[string]string,
) (*PlaceOrderResponse, error) {
	const path = "/order"
	s.log.Debug("PlaceOrder request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivatePost(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в PlaceOrder", "error", err)
		return nil, err
	}

	var result PlaceOrderResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal PlaceOrder response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// BatchOrder 4. Пакетное размещение ордеров (Batch Orders).
func (s *SpotTradeClient) BatchOrder(
	ctx context.Context,
	params map[string]string,
) (*BatchOrderResponse, error) {
	const path = "/batchOrders"
	s.log.Debug("BatchOrder request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivatePost(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в BatchOrder", "error", err)
		return nil, err
	}

	var result BatchOrderResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal BatchOrder response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// CancelOrder 5. Отменить ордер (Cancel Order).
func (s *SpotTradeClient) CancelOrder(
	ctx context.Context,
	params map[string]string,
) (*CancelOrderResponse, error) {
	const path = "/order"
	s.log.Debug("CancelOrder request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateDelete(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в CancelOrder", "error", err)
		return nil, err
	}

	var result CancelOrderResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal CancelOrder response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// CancelAllOrders 6. Отменить все ордера по символу (Cancel all Open Orders on a Symbol).
func (s *SpotTradeClient) CancelAllOrders(
	ctx context.Context,
	params map[string]string,
) (*CancelAllOrdersResponse, error) {
	const path = "/openOrders"
	s.log.Debug("CancelAllOrders request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateDelete(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в CancelAllOrders", "error", err)
		return nil, err
	}

	var result CancelAllOrdersResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal CancelAllOrders response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// QueryOrder 7. Информация об ордере (Query Order).
func (s *SpotTradeClient) QueryOrder(
	ctx context.Context,
	params map[string]string,
) (*QueryOrderResponse, error) {
	const path = "/order"
	s.log.Debug("QueryOrder request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в QueryOrder", "error", err)
		return nil, err
	}

	var result QueryOrderResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal QueryOrder response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// OpenOrder 8. Открытые ордера (Current Open Orders).
func (s *SpotTradeClient) OpenOrder(
	ctx context.Context,
	params map[string]string,
) (*OpenOrderResponse, error) {
	const path = "/openOrders"
	s.log.Debug("OpenOrder request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в OpenOrder", "error", err)
		return nil, err
	}

	var result OpenOrderResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal OpenOrder response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// AllOrders 9. Все ордера (All Orders).
func (s *SpotTradeClient) AllOrders(
	ctx context.Context,
	params map[string]string,
) (*AllOrdersResponse, error) {
	const path = "/allOrders"
	s.log.Debug("AllOrders request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в AllOrders", "error", err)
		return nil, err
	}

	var result AllOrdersResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal AllOrders response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// SpotAccountInfo 10. Информация об аккаунте (Account Information).
func (s *SpotTradeClient) SpotAccountInfo(
	ctx context.Context,
	params map[string]string,
) (*SpotAccountInfoResponse, error) {
	const path = "/account"
	s.log.Debug("SpotAccountInfo request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в SpotAccountInfo", "error", err)
		return nil, err
	}

	var result SpotAccountInfoResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal SpotAccountInfo response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// SpotMyTrade 11. История сделок (Account Trade List).
func (s *SpotTradeClient) SpotMyTrade(
	ctx context.Context,
	params map[string]string,
) (*SpotMyTradeResponse, error) {
	const path = "/myTrades"
	s.log.Debug("SpotMyTrade request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в SpotMyTrade", "error", err)
		return nil, err
	}

	var result SpotMyTradeResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal SpotMyTrade response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// MxDeduct 12. Включить MX‑дедукцию (Enable MX Deduct).
func (s *SpotTradeClient) MxDeduct(
	ctx context.Context,
	params map[string]string,
) (*MxDeductResponse, error) {
	const path = "/mxDeduct/enable"
	s.log.Debug("MxDeduct request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivatePost(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в MxDeduct", "error", err)
		return nil, err
	}

	var result MxDeductResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal MxDeduct response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// QueryMxDeduct 13. Статус MX‑дедукции (Query MX Deduct Status).
func (s *SpotTradeClient) QueryMxDeduct(
	ctx context.Context,
	params map[string]string,
) (*QueryMxDeductResponse, error) {
	const path = "/mxDeduct/enable"
	s.log.Debug("QueryMxDeduct request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в QueryMxDeduct", "error", err)
		return nil, err
	}

	var result QueryMxDeductResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal QueryMxDeduct response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}
