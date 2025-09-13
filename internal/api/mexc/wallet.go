package mexc

import (
	"context"

	"github.com/go-resty/resty/v2"
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
) (*resty.Response, error) {
	casePath := "/capital/config/getall"
	s.log.Debug("QueryCurrencyInfo", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка QueryCurrencyInfo", "error", err)
		return nil, err
	}
	return resp, nil
}

// Withdraw 2. Вывод средств (Withdraw).
func (s *SpotWalletClient) Withdraw(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/capital/withdraw/apply"
	s.log.Debug("Withdraw", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка Withdraw", "error", err)
		return nil, err
	}
	return resp, nil
}

// CancelWithdraw 3. Отменить вывод (Cancel withdraw).
func (s *SpotWalletClient) CancelWithdraw(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/capital/withdraw"
	s.log.Debug("CancelWithdraw", "path", casePath, "params", params)
	resp, err := s.client.PrivateDelete(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка CancelWithdraw", "error", err)
		return nil, err
	}
	return resp, nil
}

// DepositHistory 4. История депозитов (Deposit History).
func (s *SpotWalletClient) DepositHistory(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/capital/deposit/hisrec"
	s.log.Debug("DepositHistory", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка DepositHistory", "error", err)
		return nil, err
	}
	return resp, nil
}

// WithdrawHistory 5. История выводов (Withdraw History).
func (s *SpotWalletClient) WithdrawHistory(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/capital/withdraw/history"
	s.log.Debug("WithdrawHistory", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка WithdrawHistory", "error", err)
		return nil, err
	}
	return resp, nil
}

// GenDepositAddress 6. Сгенерировать адрес депозита (Generate deposit address).
func (s *SpotWalletClient) GenDepositAddress(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/capital/deposit/address"
	s.log.Debug("GenDepositAddress", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка GenDepositAddress", "error", err)
		return nil, err
	}
	return resp, nil
}

// DepositAddress 7. Получить адрес депозита (Deposit Address).
func (s *SpotWalletClient) DepositAddress(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/capital/deposit/address"
	s.log.Debug("DepositAddress", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка DepositAddress", "error", err)
		return nil, err
	}
	return resp, nil
}

// WithdrawAddress 8. Получить адрес вывода (Withdraw Address).
func (s *SpotWalletClient) WithdrawAddress(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/capital/withdraw/address"
	s.log.Debug("WithdrawAddress", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка WithdrawAddress", "error", err)
		return nil, err
	}
	return resp, nil
}

// Transfer 9. Универсальный перевод (User Universal Transfer).
func (s *SpotWalletClient) Transfer(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/capital/transfer"
	s.log.Debug("Transfer", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка Transfer", "error", err)
		return nil, err
	}
	return resp, nil
}

// TransferHistory 10. История переводов (Query User Universal Transfer History).
func (s *SpotWalletClient) TransferHistory(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/capital/transfer"
	s.log.Debug("TransferHistory", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка TransferHistory", "error", err)
		return nil, err
	}
	return resp, nil
}

// TransferHistoryByID 11. История перевода по tranId (Query User Universal Transfer History （by tranId）).
func (s *SpotWalletClient) TransferHistoryByID(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/capital/transfer/tranId"
	s.log.Debug("TransferHistoryByID", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка TransferHistoryByID", "error", err)
		return nil, err
	}
	return resp, nil
}

// ConvertList 12. Список активов для конвертации (Get Assets That Can Be Converted Into MX).
func (s *SpotWalletClient) ConvertList(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/capital/convert/list"
	s.log.Debug("ConvertList", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка ConvertList", "error", err)
		return nil, err
	}
	return resp, nil
}

// Convert 13. Конвертация мелких активов (Dust Transfer).
func (s *SpotWalletClient) Convert(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/capital/convert"
	s.log.Debug("Convert", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка Convert", "error", err)
		return nil, err
	}
	return resp, nil
}

// ConvertHistory 14. История конвертаций (DustLog).
func (s *SpotWalletClient) ConvertHistory(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/capital/convert"
	s.log.Debug("ConvertHistory", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка ConvertHistory", "error", err)
		return nil, err
	}
	return resp, nil
}

// ETFInfo 15. Информация об ETF (Get ETF info).
func (s *SpotWalletClient) ETFInfo(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/etf/info"
	s.log.Debug("ETFInfo", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка ETFInfo", "error", err)
		return nil, err
	}
	return resp, nil
}

// InternalTransfer 16. Внутренний перевод (Internal Transfer).
func (s *SpotWalletClient) InternalTransfer(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/capital/transfer/internal"
	s.log.Debug("InternalTransfer", "path", casePath, "params", params)
	resp, err := s.client.PrivatePost(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка InternalTransfer", "error", err)
		return nil, err
	}
	return resp, nil
}

// InternalTransferHistory 17. История внутренних переводов (Internal Transfer History).
func (s *SpotWalletClient) InternalTransferHistory(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/capital/transfer/internal"
	s.log.Debug("InternalTransferHistory", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка InternalTransferHistory", "error", err)
		return nil, err
	}
	return resp, nil
}
