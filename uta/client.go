package uta

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/khanbekov/go-bitget/common"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
)

// Client represents the UTA API client
type Client struct {
	APIKey      string
	SecretKey   string
	Passphrase  string
	BaseURL     string
	HTTPClient  *fasthttp.Client
	Logger      zerolog.Logger
	json        jsoniter.API
	DemoTrading bool // Enable demo trading mode
}

// NewClient creates a new UTA API client
func NewClient(apiKey, secretKey, passphrase string) *Client {
	return &Client{
		APIKey:     apiKey,
		SecretKey:  secretKey,
		Passphrase: passphrase,
		BaseURL:    BaseURL,
		HTTPClient: &fasthttp.Client{},
		Logger:     zerolog.Nop(),
		json:       jsoniter.ConfigCompatibleWithStandardLibrary,
	}
}

// NewClientWithLogger creates a new UTA API client with custom logger
func NewClientWithLogger(apiKey, secretKey, passphrase string, logger zerolog.Logger) *Client {
	return &Client{
		APIKey:     apiKey,
		SecretKey:  secretKey,
		Passphrase: passphrase,
		BaseURL:    BaseURL,
		HTTPClient: &fasthttp.Client{},
		Logger:     logger,
		json:       jsoniter.ConfigCompatibleWithStandardLibrary,
	}
}

// SetBaseURL sets a custom base URL for the client
func (c *Client) SetBaseURL(baseURL string) *Client {
	c.BaseURL = baseURL
	return c
}

// SetHTTPClient sets a custom HTTP client
func (c *Client) SetHTTPClient(client *fasthttp.Client) *Client {
	c.HTTPClient = client
	return c
}

// SetDemoTrading enables or disables demo trading mode
func (c *Client) SetDemoTrading(demoTrading bool) *Client {
	c.DemoTrading = demoTrading
	return c
}

// CallAPI makes an API call to the UTA API
func (c *Client) CallAPI(ctx context.Context, method string, endpoint string, queryParams url.Values, body []byte, sign bool) (*ApiResponse, *fasthttp.ResponseHeader, error) {
	// Build URL
	fullURL := c.BaseURL + endpoint
	if queryParams != nil && len(queryParams) > 0 {
		fullURL += "?" + queryParams.Encode()
	}

	// Create request
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(fullURL)
	req.Header.SetMethod(method)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "go-bitget-uta/1.0")

	if body != nil {
		req.SetBody(body)
	}

	// Add authentication headers if required
	if sign {
		timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

		// Build signature string
		var signString strings.Builder
		signString.WriteString(timestamp)
		signString.WriteString(method)
		signString.WriteString(endpoint)
		if queryParams != nil && len(queryParams) > 0 {
			signString.WriteString("?")
			signString.WriteString(queryParams.Encode())
		}
		if body != nil {
			signString.WriteString(string(body))
		}

		// Create signature
		signature := c.createSignature(signString.String())

		// Set headers
		req.Header.Set("ACCESS-KEY", c.APIKey)
		req.Header.Set("ACCESS-SIGN", signature)
		req.Header.Set("ACCESS-TIMESTAMP", timestamp)
		req.Header.Set("ACCESS-PASSPHRASE", c.Passphrase)
	}

	// Add demo trading header if enabled
	if c.DemoTrading {
		req.Header.Set("paptrading", "1")
	}

	c.Logger.Debug().
		Str("method", method).
		Str("url", fullURL).
		Str("body", string(body)).
		Bool("signed", sign).
		Msg("Making UTA API request")

	// Make request with context
	err := c.HTTPClient.DoTimeout(req, resp, 30*time.Second)
	if err != nil {
		c.Logger.Error().Err(err).Msg("HTTP request failed")
		return nil, nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	// Check status code
	statusCode := resp.StatusCode()
	if statusCode != fasthttp.StatusOK {
		c.Logger.Error().
			Int("status_code", statusCode).
			Str("response", string(resp.Body())).
			Msg("API request failed with non-200 status")
		return nil, nil, fmt.Errorf("API request failed with status %d: %s", statusCode, string(resp.Body()))
	}

	// Parse response
	var apiResp ApiResponse
	if err := c.json.Unmarshal(resp.Body(), &apiResp); err != nil {
		c.Logger.Error().
			Err(err).
			Str("response_body", string(resp.Body())).
			Msg("Failed to unmarshal API response")
		return nil, nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	c.Logger.Debug().
		Str("code", apiResp.Code).
		Str("msg", apiResp.Msg).
		Int64("request_time", apiResp.RequestTime).
		Msg("Received UTA API response")

	// Check for API errors
	if apiResp.Code != "00000" {
		apiError := &common.APIError{
			Code:    apiResp.Code,
			Message: apiResp.Msg,
		}
		c.Logger.Error().
			Str("error_code", apiError.Code).
			Str("error_message", apiError.Message).
			Msg("API returned error")
		return &apiResp, &resp.Header, apiError
	}

	return &apiResp, &resp.Header, nil
}

// createSignature creates HMAC SHA256 signature for request authentication
func (c *Client) createSignature(message string) string {
	mac := hmac.New(sha256.New, []byte(c.SecretKey))
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// Service factory methods

// Account management services
func (c *Client) NewAccountInfoService() *AccountInfoService {
	return &AccountInfoService{c: c}
}

func (c *Client) NewAccountAssetsService() *AccountAssetsService {
	return &AccountAssetsService{c: c}
}

func (c *Client) NewAccountFundingAssetsService() *AccountFundingAssetsService {
	return &AccountFundingAssetsService{c: c}
}

func (c *Client) NewAccountFeeRateService() *AccountFeeRateService {
	return &AccountFeeRateService{c: c}
}

func (c *Client) NewSetHoldingModeService() *SetHoldingModeService {
	return &SetHoldingModeService{c: c}
}

func (c *Client) NewSetLeverageService() *SetLeverageService {
	return &SetLeverageService{c: c}
}

func (c *Client) NewSwitchAccountService() *SwitchAccountService {
	return &SwitchAccountService{c: c}
}

func (c *Client) NewGetSwitchStatusService() *GetSwitchStatusService {
	return &GetSwitchStatusService{c: c}
}

// Transfer services
func (c *Client) NewTransferService() *TransferService {
	return &TransferService{c: c}
}

func (c *Client) NewSubTransferService() *SubTransferService {
	return &SubTransferService{c: c}
}

func (c *Client) NewGetTransferRecordsService() *GetTransferRecordsService {
	return &GetTransferRecordsService{c: c}
}

func (c *Client) NewGetTransferableCoinsService() *GetTransferableCoinsService {
	return &GetTransferableCoinsService{c: c}
}

// Deposit and withdrawal services
func (c *Client) NewGetDepositAddressService() *GetDepositAddressService {
	return &GetDepositAddressService{c: c}
}

func (c *Client) NewGetDepositRecordsService() *GetDepositRecordsService {
	return &GetDepositRecordsService{c: c}
}

func (c *Client) NewGetSubDepositAddressService() *GetSubDepositAddressService {
	return &GetSubDepositAddressService{c: c}
}

func (c *Client) NewGetSubDepositRecordsService() *GetSubDepositRecordsService {
	return &GetSubDepositRecordsService{c: c}
}

func (c *Client) NewWithdrawalService() *WithdrawalService {
	return &WithdrawalService{c: c}
}

func (c *Client) NewGetWithdrawalRecordsService() *GetWithdrawalRecordsService {
	return &GetWithdrawalRecordsService{c: c}
}

func (c *Client) NewSetDepositAccountService() *SetDepositAccountService {
	return &SetDepositAccountService{c: c}
}

// Financial records services
func (c *Client) NewGetFinancialRecordsService() *GetFinancialRecordsService {
	return &GetFinancialRecordsService{c: c}
}

func (c *Client) NewGetConvertRecordsService() *GetConvertRecordsService {
	return &GetConvertRecordsService{c: c}
}

func (c *Client) NewGetDeductInfoService() *GetDeductInfoService {
	return &GetDeductInfoService{c: c}
}

func (c *Client) NewSwitchDeductService() *SwitchDeductService {
	return &SwitchDeductService{c: c}
}

func (c *Client) NewGetPaymentCoinsService() *GetPaymentCoinsService {
	return &GetPaymentCoinsService{c: c}
}

func (c *Client) NewGetRepayableCoinsService() *GetRepayableCoinsService {
	return &GetRepayableCoinsService{c: c}
}

func (c *Client) NewRepayService() *RepayService {
	return &RepayService{c: c}
}

// Sub-account management services
func (c *Client) NewCreateSubAccountService() *CreateSubAccountService {
	return &CreateSubAccountService{c: c}
}

func (c *Client) NewGetSubAccountListService() *GetSubAccountListService {
	return &GetSubAccountListService{c: c}
}

func (c *Client) NewFreezeSubAccountService() *FreezeSubAccountService {
	return &FreezeSubAccountService{c: c}
}

func (c *Client) NewCreateSubAccountAPIKeyService() *CreateSubAccountAPIKeyService {
	return &CreateSubAccountAPIKeyService{c: c}
}

func (c *Client) NewGetSubAccountAPIKeysService() *GetSubAccountAPIKeysService {
	return &GetSubAccountAPIKeysService{c: c}
}

func (c *Client) NewModifySubAccountAPIKeyService() *ModifySubAccountAPIKeyService {
	return &ModifySubAccountAPIKeyService{c: c}
}

func (c *Client) NewDeleteSubAccountAPIKeyService() *DeleteSubAccountAPIKeyService {
	return &DeleteSubAccountAPIKeyService{c: c}
}

func (c *Client) NewGetSubAccountAssetsService() *GetSubAccountAssetsService {
	return &GetSubAccountAssetsService{c: c}
}

// Trading services
func (c *Client) NewPlaceOrderService() *PlaceOrderService {
	return &PlaceOrderService{c: c}
}

func (c *Client) NewCancelOrderService() *CancelOrderService {
	return &CancelOrderService{c: c}
}

func (c *Client) NewModifyOrderService() *ModifyOrderService {
	return &ModifyOrderService{c: c}
}

func (c *Client) NewBatchPlaceOrdersService() *BatchPlaceOrdersService {
	return &BatchPlaceOrdersService{c: c}
}

func (c *Client) NewBatchCancelOrdersService() *BatchCancelOrdersService {
	return &BatchCancelOrdersService{c: c}
}

func (c *Client) NewBatchModifyOrdersService() *BatchModifyOrdersService {
	return &BatchModifyOrdersService{c: c}
}

func (c *Client) NewCancelAllOrdersService() *CancelAllOrdersService {
	return &CancelAllOrdersService{c: c}
}

func (c *Client) NewCloseAllPositionsService() *CloseAllPositionsService {
	return &CloseAllPositionsService{c: c}
}

func (c *Client) NewCountdownCancelAllService() *CountdownCancelAllService {
	return &CountdownCancelAllService{c: c}
}

// Strategy order services
func (c *Client) NewPlaceStrategyOrderService() *PlaceStrategyOrderService {
	return &PlaceStrategyOrderService{c: c}
}

func (c *Client) NewCancelStrategyOrderService() *CancelStrategyOrderService {
	return &CancelStrategyOrderService{c: c}
}

func (c *Client) NewModifyStrategyOrderService() *ModifyStrategyOrderService {
	return &ModifyStrategyOrderService{c: c}
}

func (c *Client) NewGetUnfilledStrategyOrdersService() *GetUnfilledStrategyOrdersService {
	return &GetUnfilledStrategyOrdersService{c: c}
}

func (c *Client) NewGetStrategyOrderHistoryService() *GetStrategyOrderHistoryService {
	return &GetStrategyOrderHistoryService{c: c}
}

// Order and position query services
func (c *Client) NewGetOpenOrdersService() *GetOpenOrdersService {
	return &GetOpenOrdersService{c: c}
}

func (c *Client) NewGetOrderDetailsService() *GetOrderDetailsService {
	return &GetOrderDetailsService{c: c}
}

func (c *Client) NewGetOrderHistoryService() *GetOrderHistoryService {
	return &GetOrderHistoryService{c: c}
}

func (c *Client) NewGetFillHistoryService() *GetFillHistoryService {
	return &GetFillHistoryService{c: c}
}

func (c *Client) NewGetCurrentPositionsService() *GetCurrentPositionsService {
	return &GetCurrentPositionsService{c: c}
}

func (c *Client) NewGetPositionHistoryService() *GetPositionHistoryService {
	return &GetPositionHistoryService{c: c}
}

func (c *Client) NewGetMaxOpenAvailableService() *GetMaxOpenAvailableService {
	return &GetMaxOpenAvailableService{c: c}
}

func (c *Client) NewGetLoanOrdersService() *GetLoanOrdersService {
	return &GetLoanOrdersService{c: c}
}

// Market data services
func (c *Client) NewGetTickersService() *GetTickersService {
	return &GetTickersService{c: c}
}

func (c *Client) NewGetCandlesticksService() *GetCandlesticksService {
	return &GetCandlesticksService{c: c}
}

func (c *Client) NewGetHistoryCandlesticksService() *GetHistoryCandlesticksService {
	return &GetHistoryCandlesticksService{c: c}
}

func (c *Client) NewGetOrderBookService() *GetOrderBookService {
	return &GetOrderBookService{c: c}
}

// General market data services
func (c *Client) NewGetCurrentFundingRateService() *GetCurrentFundingRateService {
	return &GetCurrentFundingRateService{c: c}
}

func (c *Client) NewGetFundingRateHistoryService() *GetFundingRateHistoryService {
	return &GetFundingRateHistoryService{c: c}
}

func (c *Client) NewGetInstrumentsService() *GetInstrumentsService {
	return &GetInstrumentsService{c: c}
}

func (c *Client) NewGetDiscountRateService() *GetDiscountRateService {
	return &GetDiscountRateService{c: c}
}

func (c *Client) NewGetMarginLoansService() *GetMarginLoansService {
	return &GetMarginLoansService{c: c}
}

func (c *Client) NewGetOpenInterestService() *GetOpenInterestService {
	return &GetOpenInterestService{c: c}
}

func (c *Client) NewGetOILimitService() *GetOILimitService {
	return &GetOILimitService{c: c}
}

func (c *Client) NewGetProofOfReservesService() *GetProofOfReservesService {
	return &GetProofOfReservesService{c: c}
}

func (c *Client) NewGetRiskReserveService() *GetRiskReserveService {
	return &GetRiskReserveService{c: c}
}

func (c *Client) NewGetPositionTierService() *GetPositionTierService {
	return &GetPositionTierService{c: c}
}

func (c *Client) NewGetRecentPublicFillsService() *GetRecentPublicFillsService {
	return &GetRecentPublicFillsService{c: c}
}

// Ensure Client implements ClientInterface
var _ ClientInterface = (*Client)(nil)
