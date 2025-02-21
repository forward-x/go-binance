package portfolio

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/adshao/go-binance/v2/common"
	"github.com/bitly/go-simplejson"
)

// SideType define side type of order
type SideType string

// PositionSideType define position side type of order
type PositionSideType string

// TimeInForceType define time in force type of order
type TimeInForceType string

// NewOrderRespType define response JSON verbosity
type NewOrderRespType string

// OrderType define order type
type OrderType string

// StrategyType define conditional order type of order
type StrategyType string

// WorkingType define working type d
type WorkingType string

// OrderStatusType define order status type of order
type OrderStatusType string

// StrategyStatusType define conditional order status type
type StrategyStatusType string

// ContractType define contract type
type ContractType string

// ContractStatusType define contract status type
type ContractStatusType string

// RateLimitType define rate limit type
type RateLimitType string

// RateLimitInterval define rate limit interval
type RateLimitIntervalType string

// SymbolFilterType define symbol filter type
type SymbolFilterType string

// Endpoints
var (
	BaseApiMainUrl = "https://papi.binance.com"
)

/*
	baseasseet refers to the asset that is the quantity of a symbol.
	quoteAsset refers to the asset that is the price of a symbol.
	Margin refers to Cross Margin
	UM refers to USD-M Futures
	CM refers to Coin-M Futures
*/

// Global enums
const (
	SideTypeBuy  SideType = "BUY"
	SideTypeSell SideType = "SELL"

	PositionSideTypeBoth  PositionSideType = "BOTH"
	PositionSideTypeLong  PositionSideType = "LONG"
	PositionSideTypeShort PositionSideType = "SHORT"

	TimeInForceTypeGTC TimeInForceType = "GTC" // Good Till Cancel
	TimeInForceTypeIOC TimeInForceType = "IOC" // Immediate or Cancel
	TimeInForceTypeFOK TimeInForceType = "FOK" // Fill or Kill
	TimeInForceTypeGTX TimeInForceType = "GTX" // Good Till Crossing (Post Only)

	NewOrderRespTypeACK    NewOrderRespType = "ACK"
	NewOrderRespTypeRESULT NewOrderRespType = "RESULT"

	OrderTypeLimit  OrderType = "LIMIT"
	OrderTypeMarket OrderType = "MARKET"

	StrategyTypeStop               StrategyType = "STOP"
	StrategyTypeStopMarket         StrategyType = "STOP_MARKET"
	StrategyTypeTakeProfit         StrategyType = "TAKE_PROFIT"
	StrategyTypeTakeProfitMarket   StrategyType = "TAKE_PROFIT_MARKET"
	StrategyTypeTrailingStopMarket StrategyType = "TRAILING_STOP_MARKET"

	WorkingTypeMarkPrice WorkingType = "MARK_PRICE"

	OrderStatusTypeNew             OrderStatusType = "NEW"
	OrderStatusTypeCanceled        OrderStatusType = "CANCELED"
	OrderStatusTypeRejected        OrderStatusType = "REJECTED"
	OrderStatusTypePartiallyFilled OrderStatusType = "PARTIALLY_FILLED"
	OrderStatusTypeFilled          OrderStatusType = "FILLED"
	OrderStatusTypeExpired         OrderStatusType = "EXPIRED"

	StrategyStatusTypeNew       StrategyStatusType = "NEW"
	StrategyStatusTypeCanceled  StrategyStatusType = "CANCELED"
	StrategyStatusTypeTriggered StrategyStatusType = "TRIGGERED"
	StrategyStatusTypeFinished  StrategyStatusType = "FINISHED"
	StrategyStatusTypeExpired   StrategyStatusType = "EXPIRED"

	SymbolFilterTypeLotSize          SymbolFilterType = "LOT_SIZE"
	SymbolFilterTypePrice            SymbolFilterType = "PRICE_FILTER"
	SymbolFilterTypePercentPrice     SymbolFilterType = "PERCENT_PRICE"
	SymbolFilterTypeMarketLotSize    SymbolFilterType = "MARKET_LOT_SIZE"
	SymbolFilterTypeMaxNumOrders     SymbolFilterType = "MAX_NUM_ORDERS"
	SymbolFilterTypeMaxNumAlgoOrders SymbolFilterType = "MAX_NUM_ALGO_ORDERS"
	SymbolFilterTypeMinNotional      SymbolFilterType = "MIN_NOTIONAL"

	ContractTypePerpetual           ContractType = "PERPETUAL"
	ContractTypePerpetualDelivering ContractType = "PERPETUAL_DELIVERING"
	ContractTypeCurrentMonth        ContractType = "CURRENT_MONTH"
	ContractTypeNextMonth           ContractType = "NEXT_MONTH"
	ContractTypeCurrentQuarter      ContractType = "CURRENT_QUARTER"
	ContractTypeNextQuarter         ContractType = "NEXT_QUARTER"

	ContractStatusTypePendingTrading ContractStatusType = "PENDING_TRADING"
	ContractStatusTypeTrading        ContractStatusType = "TRADING"
	ContractStatusTypePreDelivering  ContractStatusType = "PRE_DELIVERING"
	ContractStatusTypeDelivering     ContractStatusType = "DELIVERING"
	ContractStatusTypeDelivered      ContractStatusType = "DELIVERED"
	ContractStatusTypePreSettle      ContractStatusType = "PRE_SETTLE"
	ContractStatusTypeSettling       ContractStatusType = "SETTLING"
	ContractStatusTypeClose          ContractStatusType = "CLOSE"

	RateLimitTypeRequestWeight RateLimitType = "REQUEST_WEIGHT"
	RateLimitTypeOrders        RateLimitType = "ORDERS"

	RateLimitIntervalTypeMinute RateLimitIntervalType = "MINUTE"

	timestampKey  = "timestamp"
	signatureKey  = "signature"
	recvWindowKey = "recvWindow"
)

func currentTimestamp() int64 {
	return int64(time.Nanosecond) * time.Now().UnixNano() / int64(time.Millisecond)
}

func newJSON(data []byte) (j *simplejson.Json, err error) {
	j, err = simplejson.NewJson(data)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// getApiEndpoint return the base endpoint of the Rest API according the UseTestnet flag
func getApiEndpoint() string {
	return BaseApiMainUrl
}

// NewClient initialize an API client instance with API key and secret key.
// You should always call this function before using this SDK.
// Services will be created by the form client.NewXXXService().
func NewClient(apiKey, secretKey string) *Client {
	return &Client{
		APIKey:     apiKey,
		SecretKey:  secretKey,
		KeyType:    common.KeyTypeHmac,
		BaseURL:    getApiEndpoint(),
		UserAgent:  "Binance/golang",
		HTTPClient: http.DefaultClient,
		Logger:     log.New(os.Stderr, "Binance-golang ", log.LstdFlags),
	}
}

// NewProxiedClient passing a proxy url
func NewProxiedClient(apiKey, secretKey, proxyUrl string) *Client {
	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		log.Fatal(err)
	}
	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &Client{
		APIKey:    apiKey,
		SecretKey: secretKey,
		KeyType:   common.KeyTypeHmac,
		BaseURL:   getApiEndpoint(),
		UserAgent: "Binance/golang",
		HTTPClient: &http.Client{
			Transport: tr,
		},
		Logger: log.New(os.Stderr, "Binance-golang ", log.LstdFlags),
	}
}

type doFunc func(req *http.Request) (*http.Response, error)

// Client define API client
type Client struct {
	APIKey     string
	SecretKey  string
	KeyType    string
	BaseURL    string
	UserAgent  string
	HTTPClient *http.Client
	Debug      bool
	Logger     *log.Logger
	TimeOffset int64
	do         doFunc
}

func (c *Client) debug(format string, v ...interface{}) {
	if c.Debug {
		c.Logger.Printf(format, v...)
	}
}

func (c *Client) parseRequest(r *request, opts ...RequestOption) (err error) {
	// set request options from user
	for _, opt := range opts {
		opt(r)
	}
	err = r.validate()
	if err != nil {
		return err
	}

	fullURL := fmt.Sprintf("%s%s", c.BaseURL, r.endpoint)
	if r.recvWindow > 0 {
		r.setParam(recvWindowKey, r.recvWindow)
	}
	if r.secType == secTypeSigned {
		r.setParam(timestampKey, currentTimestamp()-c.TimeOffset)
	}
	queryString := r.query.Encode()
	body := &bytes.Buffer{}
	bodyString := r.form.Encode()
	header := http.Header{}
	if r.header != nil {
		header = r.header.Clone()
	}
	if bodyString != "" {
		header.Set("Content-Type", "application/x-www-form-urlencoded")
		body = bytes.NewBufferString(bodyString)
	}
	if r.secType == secTypeAPIKey || r.secType == secTypeSigned {
		header.Set("X-MBX-APIKEY", c.APIKey)
	}
	kt := c.KeyType
	if kt == "" {
		kt = common.KeyTypeHmac
	}
	sf, err := common.SignFunc(kt)
	if err != nil {
		return err
	}
	if r.secType == secTypeSigned {
		raw := fmt.Sprintf("%s%s", queryString, bodyString)
		sign, err := sf(c.SecretKey, raw)
		if err != nil {
			return err
		}
		v := url.Values{}
		v.Set(signatureKey, *sign)
		if queryString == "" {
			queryString = v.Encode()
		} else {
			queryString = fmt.Sprintf("%s&%s", queryString, v.Encode())
		}
	}
	if queryString != "" {
		fullURL = fmt.Sprintf("%s?%s", fullURL, queryString)
	}
	c.debug("full url: %s, body: %s\n", fullURL, bodyString)

	r.fullURL = fullURL
	r.header = header
	r.body = body
	return nil
}

func (c *Client) callAPI(ctx context.Context, r *request, opts ...RequestOption) (data []byte, header *http.Header, err error) {
	err = c.parseRequest(r, opts...)
	if err != nil {
		return []byte{}, &http.Header{}, err
	}
	req, err := http.NewRequest(r.method, r.fullURL, r.body)
	if err != nil {
		return []byte{}, &http.Header{}, err
	}
	req = req.WithContext(ctx)
	req.Header = r.header
	c.debug("request: %#v\n", req)
	f := c.do
	if f == nil {
		f = c.HTTPClient.Do
	}
	res, err := f(req)
	if err != nil {
		return []byte{}, &http.Header{}, err
	}
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, &http.Header{}, err
	}
	defer func() {
		cerr := res.Body.Close()
		// Only overwrite the returned error if the original error was nil and an
		// error occurred while closing the body.
		if err == nil && cerr != nil {
			err = cerr
		}
	}()
	c.debug("response: %#v\n", res)
	c.debug("response body: %s\n", string(data))
	c.debug("response status code: %d\n", res.StatusCode)

	if res.StatusCode >= http.StatusBadRequest {
		apiErr := new(common.APIError)
		e := json.Unmarshal(data, apiErr)
		if e != nil {
			c.debug("failed to unmarshal json: %s\n", e)
		}
		if !apiErr.IsValid() {
			apiErr.Response = data
		}
		return nil, &res.Header, apiErr
	}
	return data, &res.Header, nil
}

// SetApiEndpoint set api Endpoint
func (c *Client) SetApiEndpoint(url string) *Client {
	c.BaseURL = url
	return c
}

// NewPingService init ping service
func (c *Client) NewPingService() *PingService {
	return &PingService{c: c}
}

// NewGetBalanceService init getting balance service
func (c *Client) NewGetBalanceService() *GetBalanceService {
	return &GetBalanceService{c: c}
}

// NewGetUmPositionRiskService init getting Um position risk service
func (c *Client) NewGetUmPositionRiskService() *GetUmPositionRiskService {
	return &GetUmPositionRiskService{c: c}
}

// NewGetCmPositionRiskService init getting Cm position risk service
func (c *Client) NewGetCmPositionRiskService() *GetCmPositionRiskService {
	return &GetCmPositionRiskService{c: c}
}
