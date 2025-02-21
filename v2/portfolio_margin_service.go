package binance

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/adshao/go-binance/v2/common"
)

// GetPortfolioMarginAssetIndexPriceService get asset index price of an asset or all assets.
type GetPortfolioMarginAssetIndexPriceService struct {
	c     *Client
	asset *string
}

// Asset set asset.
func (s *GetPortfolioMarginAssetIndexPriceService) Asset(asset string) *GetPortfolioMarginAssetIndexPriceService {
	s.asset = &asset
	return s
}

// Do send request.
func (s *GetPortfolioMarginAssetIndexPriceService) Do(ctx context.Context, opts ...RequestOption) (res []*AssetIndexPrice, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/portfolio/asset-index-price",
		secType:  secTypeSigned,
	}
	if s.asset != nil {
		r.setParam("asset", *s.asset)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	data = common.ToJSONList(data)
	if err != nil {
		return []*AssetIndexPrice{}, err
	}
	res = make([]*AssetIndexPrice, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*AssetIndexPrice{}, err
	}
	return res, nil
}

// AssetIndexPrice define asset index price.
type AssetIndexPrice struct {
	Asset           string `json:"asset"`
	AssetIndexPrice string `json:"assetIndexPrice"`
	Time            int64  `json:"time"`
}
