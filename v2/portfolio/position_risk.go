package portfolio

import (
	"context"
	"encoding/json"
	"net/http"
)

// GetUmPositionRiskService get current USD-M position information.
type GetUmPositionRiskService struct {
	c      *Client
	symbol *string
}

// Symbol set symbol
func (s *GetUmPositionRiskService) Symbol(symbol string) *GetUmPositionRiskService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *GetUmPositionRiskService) Do(ctx context.Context, opts ...RequestOption) (res []*UmPositionRisk, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/um/positionRisk",
		secType:  secTypeSigned,
	}
	if s.symbol != nil {
		r.setParam("symbol", *s.symbol)
	}

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*UmPositionRisk{}, err
	}
	res = make([]*UmPositionRisk, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*UmPositionRisk{}, err
	}
	return res, nil
}

// UmPositionRisk define USD-M position risk info
type UmPositionRisk struct {
	EntryPrice       string `json:"entryPrice"`
	Leverage         string `json:"leverage"`
	MarkPrice        string `json:"markPrice"`
	MaxNotionalValue string `json:"maxNotionalValue"`
	PositionAmt      string `json:"positionAmt"`
	Notional         string `json:"notional"`
	Symbol           string `json:"symbol"`
	UnRealizedProfit string `json:"unRealizedProfit"`
	LiquidationPrice string `json:"liquidationPrice"`
	PositionSide     string `json:"positionSide"`
	UpdateTime       int64  `json:"updateTime"`
}

// GetCmPositionRiskService get current COIN-M position information.
type GetCmPositionRiskService struct {
	c           *Client
	marginAsset *string
	pair        *string
}

// MarginAsset set marginAsset
func (s *GetCmPositionRiskService) MarginAsset(marginAsset string) *GetCmPositionRiskService {
	s.marginAsset = &marginAsset
	return s
}

// Pair set pair
func (s *GetCmPositionRiskService) Pair(pair string) *GetCmPositionRiskService {
	s.pair = &pair
	return s
}

// Do send request
func (s *GetCmPositionRiskService) Do(ctx context.Context, opts ...RequestOption) (res []*CmPositionRisk, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/cm/positionRisk",
		secType:  secTypeSigned,
	}
	if s.marginAsset != nil {
		r.setParam("marginAsset", *s.marginAsset)
	}
	if s.pair != nil {
		r.setParam("pair", *s.pair)
	}

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*CmPositionRisk{}, err
	}
	res = make([]*CmPositionRisk, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*CmPositionRisk{}, err
	}
	return res, nil
}

// CmPositionRisk define COIN-M position risk info
type CmPositionRisk struct {
	Symbol           string `json:"symbol"`
	PositionAmt      string `json:"positionAmt"`
	EntryPrice       string `json:"entryPrice"`
	MarkPrice        string `json:"markPrice"`
	UnRealizedProfit string `json:"unRealizedProfit"`
	LiquidationPrice string `json:"liquidationPrice"`
	Leverage         string `json:"leverage"`
	PositionSide     string `json:"positionSide"`
	UpdateTime       int64  `json:"updateTime"`
	MaxQty           string `json:"maxQty"`
	NotionalValue    string `json:"notionalValue"`
}
