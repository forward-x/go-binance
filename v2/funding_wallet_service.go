package binance

import (
	"context"
	"encoding/json"
	"net/http"
)

// GetFundingAssetService fetches all assets in funding wallet.
type GetFundingAssetService struct {
	c                *Client
	asset            *string
	needBtcValuation bool
}

func (s *GetFundingAssetService) Asset(asset string) *GetFundingAssetService {
	s.asset = &asset
	return s
}

func (s *GetFundingAssetService) NeedBtcValuation(val bool) *GetFundingAssetService {
	s.needBtcValuation = val
	return s
}

type FundingAssetRecord struct {
	Asset        string `json:"asset"`
	Free         string `json:"free"`
	Locked       string `json:"locked"`
	Freeze       string `json:"freeze"`
	Withdrawing  string `json:"withdrawing"`
	BtcValuation string `json:"btcValuation"`
}

func (s *GetFundingAssetService) Do(ctx context.Context) (res []FundingAssetRecord, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v1/asset/get-funding-asset",
		secType:  secTypeSigned,
	}
	if s.asset != nil {
		r.setParam("asset", *s.asset)
	}
	if s.needBtcValuation {
		r.setParam("needBtcValuation", s.needBtcValuation)
	}
	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &res)
	return
}
