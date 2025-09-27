package mexc

import (
	"context"
	"encoding/json"
)

// TODO: DTOs — будут заполнены позже.

type (
	QueryCurrencyInfoResponse       struct{}
	WithdrawResponse                struct{}
	CancelWithdrawResponse          struct{}
	DepositHistoryResponse          struct{}
	WithdrawHistoryResponse         struct{}
	GenDepositAddressResponse       struct{}
	DepositAddressResponse          struct{}
	WithdrawAddressResponse         struct{}
	TransferResponse                struct{}
	TransferHistoryResponse         struct{}
	TransferHistoryByIDResponse     struct{}
	ConvertListResponse             struct{}
	ConvertResponse                 struct{}
	ConvertHistoryResponse          struct{}
	ETFInfoResponse                 struct{}
	InternalTransferResponse        struct{}
	InternalTransferHistoryResponse struct{}
)

type SpotWalletClient struct {
	log    Logger
	client APIClient
}

func NewSpotWalletClient(log Logger, client APIClient) *SpotWalletClient {
	return &SpotWalletClient{log: log, client: client}
}

// ## 3. Кошелёк (Wallet Endpoints)

// QueryCurrencyInfo 1. Информация о валюте (Query the currency information).
func (s *SpotWalletClient) QueryCurrencyInfo(
	ctx context.Context,
	params map[string]string,
) (*QueryCurrencyInfoResponse, error) {
	const path = "/capital/config/getall"
	s.log.Debug("QueryCurrencyInfo request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в QueryCurrencyInfo", "error", err)
		return nil, err
	}

	var result QueryCurrencyInfoResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal QueryCurrencyInfo response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// Withdraw 2. Вывод средств (Withdraw).
func (s *SpotWalletClient) Withdraw(
	ctx context.Context,
	params map[string]string,
) (*WithdrawResponse, error) {
	const path = "/capital/withdraw/apply"
	s.log.Debug("Withdraw request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivatePost(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в Withdraw", "error", err)
		return nil, err
	}

	var result WithdrawResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal Withdraw response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// CancelWithdraw 3. Отменить вывод (Cancel withdraw).
func (s *SpotWalletClient) CancelWithdraw(
	ctx context.Context,
	params map[string]string,
) (*CancelWithdrawResponse, error) {
	const path = "/capital/withdraw"
	s.log.Debug("CancelWithdraw request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateDelete(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в CancelWithdraw", "error", err)
		return nil, err
	}

	var result CancelWithdrawResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal CancelWithdraw response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// DepositHistory 4. История депозитов (Deposit History).
func (s *SpotWalletClient) DepositHistory(
	ctx context.Context,
	params map[string]string,
) (*DepositHistoryResponse, error) {
	const path = "/capital/deposit/hisrec"
	s.log.Debug("DepositHistory request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в DepositHistory", "error", err)
		return nil, err
	}

	var result DepositHistoryResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal DepositHistory response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// WithdrawHistory 5. История выводов (Withdraw History).
func (s *SpotWalletClient) WithdrawHistory(
	ctx context.Context,
	params map[string]string,
) (*WithdrawHistoryResponse, error) {
	const path = "/capital/withdraw/history"
	s.log.Debug("WithdrawHistory request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в WithdrawHistory", "error", err)
		return nil, err
	}

	var result WithdrawHistoryResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal WithdrawHistory response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// GenDepositAddress 6. Сгенерировать адрес депозита (Generate deposit address).
func (s *SpotWalletClient) GenDepositAddress(
	ctx context.Context,
	params map[string]string,
) (*GenDepositAddressResponse, error) {
	const path = "/capital/deposit/address"
	s.log.Debug("GenDepositAddress request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivatePost(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в GenDepositAddress", "error", err)
		return nil, err
	}

	var result GenDepositAddressResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal GenDepositAddress response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// DepositAddress 7. Получить адрес депозита (Deposit Address).
func (s *SpotWalletClient) DepositAddress(
	ctx context.Context,
	params map[string]string,
) (*DepositAddressResponse, error) {
	const path = "/capital/deposit/address"
	s.log.Debug("DepositAddress request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в DepositAddress", "error", err)
		return nil, err
	}

	var result DepositAddressResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal DepositAddress response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// WithdrawAddress 8. Получить адрес вывода (Withdraw Address).
func (s *SpotWalletClient) WithdrawAddress(
	ctx context.Context,
	params map[string]string,
) (*WithdrawAddressResponse, error) {
	const path = "/capital/withdraw/address"
	s.log.Debug("WithdrawAddress request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в WithdrawAddress", "error", err)
		return nil, err
	}

	var result WithdrawAddressResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal WithdrawAddress response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// Transfer 9. Универсальный перевод (User Universal Transfer).
func (s *SpotWalletClient) Transfer(
	ctx context.Context,
	params map[string]string,
) (*TransferResponse, error) {
	const path = "/capital/transfer"
	s.log.Debug("Transfer request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivatePost(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в Transfer", "error", err)
		return nil, err
	}

	var result TransferResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal Transfer response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// TransferHistory 10. История переводов (Query User Universal Transfer History).
func (s *SpotWalletClient) TransferHistory(
	ctx context.Context,
	params map[string]string,
) (*TransferHistoryResponse, error) {
	const path = "/capital/transfer"
	s.log.Debug("TransferHistory request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в TransferHistory", "error", err)
		return nil, err
	}

	var result TransferHistoryResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal TransferHistory response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// TransferHistoryByID 11. История перевода по tranId (Query User Universal Transfer History （by tranId）).
func (s *SpotWalletClient) TransferHistoryByID(
	ctx context.Context,
	params map[string]string,
) (*TransferHistoryByIDResponse, error) {
	const path = "/capital/transfer/tranId"
	s.log.Debug("TransferHistoryByID request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в TransferHistoryByID", "error", err)
		return nil, err
	}

	var result TransferHistoryByIDResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal TransferHistoryByID response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// ConvertList 12. Список активов для конвертации (Get Assets That Can Be Converted Into MX).
func (s *SpotWalletClient) ConvertList(
	ctx context.Context,
	params map[string]string,
) (*ConvertListResponse, error) {
	const path = "/capital/convert/list"
	s.log.Debug("ConvertList request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в ConvertList", "error", err)
		return nil, err
	}

	var result ConvertListResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal ConvertList response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// Convert 13. Конвертация мелких активов (Dust Transfer).
func (s *SpotWalletClient) Convert(
	ctx context.Context,
	params map[string]string,
) (*ConvertResponse, error) {
	const path = "/capital/convert"
	s.log.Debug("Convert request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivatePost(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в Convert", "error", err)
		return nil, err
	}

	var result ConvertResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal Convert response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// ConvertHistory 14. История конвертаций (DustLog).
func (s *SpotWalletClient) ConvertHistory(
	ctx context.Context,
	params map[string]string,
) (*ConvertHistoryResponse, error) {
	const path = "/capital/convert"
	s.log.Debug("ConvertHistory request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в ConvertHistory", "error", err)
		return nil, err
	}

	var result ConvertHistoryResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal ConvertHistory response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// ETFInfo 15. Информация об ETF (Get ETF info).
func (s *SpotWalletClient) ETFInfo(
	ctx context.Context,
	params map[string]string,
) (*ETFInfoResponse, error) {
	const path = "/etf/info"
	s.log.Debug("ETFInfo request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в ETFInfo", "error", err)
		return nil, err
	}

	var result ETFInfoResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal ETFInfo response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// InternalTransfer 16. Внутренний перевод (Internal Transfer).
func (s *SpotWalletClient) InternalTransfer(
	ctx context.Context,
	params map[string]string,
) (*InternalTransferResponse, error) {
	const path = "/capital/transfer/internal"
	s.log.Debug("InternalTransfer request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivatePost(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в InternalTransfer", "error", err)
		return nil, err
	}

	var result InternalTransferResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal InternalTransfer response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// InternalTransferHistory 17. История внутренних переводов (Internal Transfer History).
func (s *SpotWalletClient) InternalTransferHistory(
	ctx context.Context,
	params map[string]string,
) (*InternalTransferHistoryResponse, error) {
	const path = "/capital/transfer/internal"
	s.log.Debug("InternalTransferHistory request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в InternalTransferHistory", "error", err)
		return nil, err
	}

	var result InternalTransferHistoryResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal InternalTransferHistory response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}
