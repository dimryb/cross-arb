package mexc

import (
	"context"

	"github.com/go-resty/resty/v2"
)

type SpotRebateClient struct {
	log    Logger
	client APIClient
}

func NewSpotRebateClient(log Logger, client APIClient) *SpotRebateClient {
	return &SpotRebateClient{log: log, client: client}
}

// ## 5. Партнёрские и реферальные данные

// RebateHistory 1. История реферальных вознаграждений (Get Rebate History Records).
func (s *SpotRebateClient) RebateHistory(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/rebate/taxQuery"
	s.log.Debug("RebateHistory", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка RebateHistory", "error", err)
		return nil, err
	}
	return resp, nil
}

// RebateDetail 2. Детали реферальных выплат (Get Rebate Records Detail).
func (s *SpotRebateClient) RebateDetail(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/rebate/detail"
	s.log.Debug("RebateDetail", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка RebateDetail", "error", err)
		return nil, err
	}
	return resp, nil
}

// SelfRecordsDetail 3. Детали собственных выплат (Get Self Rebate Records Detail).
func (s *SpotRebateClient) SelfRecordsDetail(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/rebate/detail/kickback"
	s.log.Debug("SelfRecordsDetail", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка SelfRecordsDetail", "error", err)
		return nil, err
	}
	return resp, nil
}

// ReferCode 4. Получить код приглашения (Query ReferCode).
func (s *SpotRebateClient) ReferCode(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/rebate/referCode"
	s.log.Debug("ReferCode", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка ReferCode", "error", err)
		return nil, err
	}
	return resp, nil
}

// AffiliateCommission 5. Комиссии аффилиата (Get Affiliate Commission Record (affiliate only)).
func (s *SpotRebateClient) AffiliateCommission(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/rebate/affiliate/commission"
	s.log.Debug("AffiliateCommission", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка AffiliateCommission", "error", err)
		return nil, err
	}
	return resp, nil
}

// AffiliateWithdraw 6. История выводов аффилиата (Get Affiliate Withdraw Record (affiliate only)).
func (s *SpotRebateClient) AffiliateWithdraw(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/rebate/affiliate/withdraw"
	s.log.Debug("AffiliateWithdraw", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка AffiliateWithdraw", "error", err)
		return nil, err
	}
	return resp, nil
}

// AffiliateCommissionDetail 7. Детали комиссий аффилиата (Get Affiliate Commission Detail Record (affiliate only)).
func (s *SpotRebateClient) AffiliateCommissionDetail(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/rebate/affiliate/commission/detail"
	s.log.Debug("AffiliateCommissionDetail", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка AffiliateCommissionDetail", "error", err)
		return nil, err
	}
	return resp, nil
}

// AffiliateReferral 8. Сводка реферальных данных аффилиата (Get Affiliate Referral Data（affiliate only）).
func (s *SpotRebateClient) AffiliateReferral(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/rebate/affiliate/referral"
	s.log.Debug("AffiliateReferral", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка AffiliateReferral", "error", err)
		return nil, err
	}
	return resp, nil
}

// Subaffiliates 9. Суб‑аффилиаты (Get Subaffiliates Data (affiliate only)).
func (s *SpotRebateClient) Subaffiliates(
	ctx context.Context,
	params map[string]string,
) (*resty.Response, error) {
	casePath := "/rebate/affiliate/subaffiliates"
	s.log.Debug("Subaffiliates", "path", casePath, "params", params)
	resp, err := s.client.PrivateGet(ctx, casePath, params)
	if err != nil {
		s.log.Error("Ошибка Subaffiliates", "error", err)
		return nil, err
	}
	return resp, nil
}
