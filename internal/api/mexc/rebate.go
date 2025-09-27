package mexc

import (
	"context"
	"encoding/json"
)

// TODO: DTOs — будут заполнены позже.

type (
	RebateHistoryResponse             struct{}
	RebateDetailResponse              struct{}
	SelfRecordsDetailResponse         struct{}
	ReferCodeResponse                 struct{}
	AffiliateCommissionResponse       struct{}
	AffiliateWithdrawResponse         struct{}
	AffiliateCommissionDetailResponse struct{}
	AffiliateReferralResponse         struct{}
	SubaffiliatesResponse             struct{}
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
) (*RebateHistoryResponse, error) {
	const path = "/rebate/taxQuery"
	s.log.Debug("RebateHistory request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в RebateHistory", "error", err)
		return nil, err
	}

	var result RebateHistoryResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal RebateHistory response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// RebateDetail 2. Детали реферальных выплат (Get Rebate Records Detail).
func (s *SpotRebateClient) RebateDetail(
	ctx context.Context,
	params map[string]string,
) (*RebateDetailResponse, error) {
	const path = "/rebate/detail"
	s.log.Debug("RebateDetail request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в RebateDetail", "error", err)
		return nil, err
	}

	var result RebateDetailResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal RebateDetail response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// SelfRecordsDetail 3. Детали собственных выплат (Get Self Rebate Records Detail).
func (s *SpotRebateClient) SelfRecordsDetail(
	ctx context.Context,
	params map[string]string,
) (*SelfRecordsDetailResponse, error) {
	const path = "/rebate/detail/kickback"
	s.log.Debug("SelfRecordsDetail request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в SelfRecordsDetail", "error", err)
		return nil, err
	}

	var result SelfRecordsDetailResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal SelfRecordsDetail response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// ReferCode 4. Получить код приглашения (Query ReferCode).
func (s *SpotRebateClient) ReferCode(
	ctx context.Context,
	params map[string]string,
) (*ReferCodeResponse, error) {
	const path = "/rebate/referCode"
	s.log.Debug("ReferCode request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в ReferCode", "error", err)
		return nil, err
	}

	var result ReferCodeResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal ReferCode response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// AffiliateCommission 5. Комиссии аффилиата (Get Affiliate Commission Record (affiliate only)).
func (s *SpotRebateClient) AffiliateCommission(
	ctx context.Context,
	params map[string]string,
) (*AffiliateCommissionResponse, error) {
	const path = "/rebate/affiliate/commission"
	s.log.Debug("AffiliateCommission request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в AffiliateCommission", "error", err)
		return nil, err
	}

	var result AffiliateCommissionResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal AffiliateCommission response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// AffiliateWithdraw 6. История выводов аффилиата (Get Affiliate Withdraw Record (affiliate only)).
func (s *SpotRebateClient) AffiliateWithdraw(
	ctx context.Context,
	params map[string]string,
) (*AffiliateWithdrawResponse, error) {
	const path = "/rebate/affiliate/withdraw"
	s.log.Debug("AffiliateWithdraw request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в AffiliateWithdraw", "error", err)
		return nil, err
	}

	var result AffiliateWithdrawResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal AffiliateWithdraw response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// AffiliateCommissionDetail 7. Детали комиссий аффилиата (Get Affiliate Commission Detail Record (affiliate only)).
func (s *SpotRebateClient) AffiliateCommissionDetail(
	ctx context.Context,
	params map[string]string,
) (*AffiliateCommissionDetailResponse, error) {
	const path = "/rebate/affiliate/commission/detail"
	s.log.Debug("AffiliateCommissionDetail request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в AffiliateCommissionDetail", "error", err)
		return nil, err
	}

	var result AffiliateCommissionDetailResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal AffiliateCommissionDetail response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// AffiliateReferral 8. Сводка реферальных данных аффилиата (Get Affiliate Referral Data（affiliate only）).
func (s *SpotRebateClient) AffiliateReferral(
	ctx context.Context,
	params map[string]string,
) (*AffiliateReferralResponse, error) {
	const path = "/rebate/affiliate/referral"
	s.log.Debug("AffiliateReferral request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в AffiliateReferral", "error", err)
		return nil, err
	}

	var result AffiliateReferralResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal AffiliateReferral response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}

// Subaffiliates 9. Суб‑аффилиаты (Get Subaffiliates Data (affiliate only)).
func (s *SpotRebateClient) Subaffiliates(
	ctx context.Context,
	params map[string]string,
) (*SubaffiliatesResponse, error) {
	const path = "/rebate/affiliate/subaffiliates"
	s.log.Debug("Subaffiliates request to MEXC", "path", path, "params", params)

	resp, err := s.client.PrivateGet(ctx, path, params)
	if err != nil {
		s.log.Error("Ошибка в Subaffiliates", "error", err)
		return nil, err
	}

	var result SubaffiliatesResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		s.log.Error("Failed to unmarshal Subaffiliates response", "error", err, "body", string(resp.Body()))
		return nil, err
	}

	return &result, nil
}
