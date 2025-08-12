package uta

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
)

// MockClient is a mock implementation of the ClientInterface for testing
type MockClient struct {
	mock.Mock
}

func (m *MockClient) CallAPI(ctx context.Context, method string, endpoint string, queryParams url.Values, body []byte, sign bool) (*ApiResponse, *fasthttp.ResponseHeader, error) {
	args := m.Called(ctx, method, endpoint, queryParams, body, sign)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*fasthttp.ResponseHeader), args.Error(2)
	}
	return args.Get(0).(*ApiResponse), args.Get(1).(*fasthttp.ResponseHeader), args.Error(2)
}

// Service factory method implementations for MockClient
func (m *MockClient) NewAccountInfoService() *AccountInfoService { return &AccountInfoService{c: m} }
func (m *MockClient) NewAccountAssetsService() *AccountAssetsService {
	return &AccountAssetsService{c: m}
}
func (m *MockClient) NewAccountFundingAssetsService() *AccountFundingAssetsService {
	return &AccountFundingAssetsService{c: m}
}
func (m *MockClient) NewAccountFeeRateService() *AccountFeeRateService {
	return &AccountFeeRateService{c: m}
}
func (m *MockClient) NewSetHoldingModeService() *SetHoldingModeService {
	return &SetHoldingModeService{c: m}
}
func (m *MockClient) NewSetLeverageService() *SetLeverageService { return &SetLeverageService{c: m} }
func (m *MockClient) NewSwitchAccountService() *SwitchAccountService {
	return &SwitchAccountService{c: m}
}
func (m *MockClient) NewGetSwitchStatusService() *GetSwitchStatusService {
	return &GetSwitchStatusService{c: m}
}
func (m *MockClient) NewTransferService() *TransferService       { return &TransferService{c: m} }
func (m *MockClient) NewSubTransferService() *SubTransferService { return &SubTransferService{c: m} }
func (m *MockClient) NewGetTransferRecordsService() *GetTransferRecordsService {
	return &GetTransferRecordsService{c: m}
}
func (m *MockClient) NewGetTransferableCoinsService() *GetTransferableCoinsService {
	return &GetTransferableCoinsService{c: m}
}
func (m *MockClient) NewGetDepositAddressService() *GetDepositAddressService {
	return &GetDepositAddressService{c: m}
}
func (m *MockClient) NewGetDepositRecordsService() *GetDepositRecordsService {
	return &GetDepositRecordsService{c: m}
}
func (m *MockClient) NewGetSubDepositAddressService() *GetSubDepositAddressService {
	return &GetSubDepositAddressService{c: m}
}
func (m *MockClient) NewGetSubDepositRecordsService() *GetSubDepositRecordsService {
	return &GetSubDepositRecordsService{c: m}
}
func (m *MockClient) NewWithdrawalService() *WithdrawalService { return &WithdrawalService{c: m} }
func (m *MockClient) NewGetWithdrawalRecordsService() *GetWithdrawalRecordsService {
	return &GetWithdrawalRecordsService{c: m}
}
func (m *MockClient) NewSetDepositAccountService() *SetDepositAccountService {
	return &SetDepositAccountService{c: m}
}
func (m *MockClient) NewGetFinancialRecordsService() *GetFinancialRecordsService {
	return &GetFinancialRecordsService{c: m}
}
func (m *MockClient) NewGetConvertRecordsService() *GetConvertRecordsService {
	return &GetConvertRecordsService{c: m}
}
func (m *MockClient) NewGetDeductInfoService() *GetDeductInfoService {
	return &GetDeductInfoService{c: m}
}
func (m *MockClient) NewSwitchDeductService() *SwitchDeductService { return &SwitchDeductService{c: m} }
func (m *MockClient) NewGetPaymentCoinsService() *GetPaymentCoinsService {
	return &GetPaymentCoinsService{c: m}
}
func (m *MockClient) NewGetRepayableCoinsService() *GetRepayableCoinsService {
	return &GetRepayableCoinsService{c: m}
}
func (m *MockClient) NewRepayService() *RepayService { return &RepayService{c: m} }
func (m *MockClient) NewCreateSubAccountService() *CreateSubAccountService {
	return &CreateSubAccountService{c: m}
}
func (m *MockClient) NewGetSubAccountListService() *GetSubAccountListService {
	return &GetSubAccountListService{c: m}
}
func (m *MockClient) NewFreezeSubAccountService() *FreezeSubAccountService {
	return &FreezeSubAccountService{c: m}
}
func (m *MockClient) NewCreateSubAccountAPIKeyService() *CreateSubAccountAPIKeyService {
	return &CreateSubAccountAPIKeyService{c: m}
}
func (m *MockClient) NewGetSubAccountAPIKeysService() *GetSubAccountAPIKeysService {
	return &GetSubAccountAPIKeysService{c: m}
}
func (m *MockClient) NewModifySubAccountAPIKeyService() *ModifySubAccountAPIKeyService {
	return &ModifySubAccountAPIKeyService{c: m}
}
func (m *MockClient) NewDeleteSubAccountAPIKeyService() *DeleteSubAccountAPIKeyService {
	return &DeleteSubAccountAPIKeyService{c: m}
}
func (m *MockClient) NewGetSubAccountAssetsService() *GetSubAccountAssetsService {
	return &GetSubAccountAssetsService{c: m}
}
func (m *MockClient) NewPlaceOrderService() *PlaceOrderService   { return &PlaceOrderService{c: m} }
func (m *MockClient) NewCancelOrderService() *CancelOrderService { return &CancelOrderService{c: m} }
func (m *MockClient) NewModifyOrderService() *ModifyOrderService { return &ModifyOrderService{c: m} }
func (m *MockClient) NewBatchPlaceOrdersService() *BatchPlaceOrdersService {
	return &BatchPlaceOrdersService{c: m}
}
func (m *MockClient) NewBatchCancelOrdersService() *BatchCancelOrdersService {
	return &BatchCancelOrdersService{c: m}
}
func (m *MockClient) NewBatchModifyOrdersService() *BatchModifyOrdersService {
	return &BatchModifyOrdersService{c: m}
}
func (m *MockClient) NewCancelAllOrdersService() *CancelAllOrdersService {
	return &CancelAllOrdersService{c: m}
}
func (m *MockClient) NewCloseAllPositionsService() *CloseAllPositionsService {
	return &CloseAllPositionsService{c: m}
}
func (m *MockClient) NewCountdownCancelAllService() *CountdownCancelAllService {
	return &CountdownCancelAllService{c: m}
}
func (m *MockClient) NewPlaceStrategyOrderService() *PlaceStrategyOrderService {
	return &PlaceStrategyOrderService{c: m}
}
func (m *MockClient) NewCancelStrategyOrderService() *CancelStrategyOrderService {
	return &CancelStrategyOrderService{c: m}
}
func (m *MockClient) NewModifyStrategyOrderService() *ModifyStrategyOrderService {
	return &ModifyStrategyOrderService{c: m}
}
func (m *MockClient) NewGetUnfilledStrategyOrdersService() *GetUnfilledStrategyOrdersService {
	return &GetUnfilledStrategyOrdersService{c: m}
}
func (m *MockClient) NewGetStrategyOrderHistoryService() *GetStrategyOrderHistoryService {
	return &GetStrategyOrderHistoryService{c: m}
}
func (m *MockClient) NewGetOpenOrdersService() *GetOpenOrdersService {
	return &GetOpenOrdersService{c: m}
}
func (m *MockClient) NewGetOrderDetailsService() *GetOrderDetailsService {
	return &GetOrderDetailsService{c: m}
}
func (m *MockClient) NewGetOrderHistoryService() *GetOrderHistoryService {
	return &GetOrderHistoryService{c: m}
}
func (m *MockClient) NewGetFillHistoryService() *GetFillHistoryService {
	return &GetFillHistoryService{c: m}
}
func (m *MockClient) NewGetCurrentPositionsService() *GetCurrentPositionsService {
	return &GetCurrentPositionsService{c: m}
}
func (m *MockClient) NewGetPositionHistoryService() *GetPositionHistoryService {
	return &GetPositionHistoryService{c: m}
}
func (m *MockClient) NewGetMaxOpenAvailableService() *GetMaxOpenAvailableService {
	return &GetMaxOpenAvailableService{c: m}
}
func (m *MockClient) NewGetLoanOrdersService() *GetLoanOrdersService {
	return &GetLoanOrdersService{c: m}
}
func (m *MockClient) NewGetTickersService() *GetTickersService { return &GetTickersService{c: m} }
func (m *MockClient) NewGetCandlesticksService() *GetCandlesticksService {
	return &GetCandlesticksService{c: m}
}
func (m *MockClient) NewGetHistoryCandlesticksService() *GetHistoryCandlesticksService {
	return &GetHistoryCandlesticksService{c: m}
}
func (m *MockClient) NewGetOrderBookService() *GetOrderBookService { return &GetOrderBookService{c: m} }
func (m *MockClient) NewGetCurrentFundingRateService() *GetCurrentFundingRateService {
	return &GetCurrentFundingRateService{c: m}
}
func (m *MockClient) NewGetFundingRateHistoryService() *GetFundingRateHistoryService {
	return &GetFundingRateHistoryService{c: m}
}
func (m *MockClient) NewGetInstrumentsService() *GetInstrumentsService {
	return &GetInstrumentsService{c: m}
}
func (m *MockClient) NewGetDiscountRateService() *GetDiscountRateService {
	return &GetDiscountRateService{c: m}
}
func (m *MockClient) NewGetMarginLoansService() *GetMarginLoansService {
	return &GetMarginLoansService{c: m}
}
func (m *MockClient) NewGetOpenInterestService() *GetOpenInterestService {
	return &GetOpenInterestService{c: m}
}
func (m *MockClient) NewGetOILimitService() *GetOILimitService { return &GetOILimitService{c: m} }
func (m *MockClient) NewGetProofOfReservesService() *GetProofOfReservesService {
	return &GetProofOfReservesService{c: m}
}
func (m *MockClient) NewGetRiskReserveService() *GetRiskReserveService {
	return &GetRiskReserveService{c: m}
}
func (m *MockClient) NewGetPositionTierService() *GetPositionTierService {
	return &GetPositionTierService{c: m}
}
func (m *MockClient) NewGetRecentPublicFillsService() *GetRecentPublicFillsService {
	return &GetRecentPublicFillsService{c: m}
}

// Ensure MockClient implements ClientInterface
var _ ClientInterface = (*MockClient)(nil)

func TestNewClient(t *testing.T) {
	apiKey := "test_api_key"
	secretKey := "test_secret_key"
	passphrase := "test_passphrase"

	client := NewClient(apiKey, secretKey, passphrase)

	assert.NotNil(t, client)
	assert.Equal(t, apiKey, client.APIKey)
	assert.Equal(t, secretKey, client.SecretKey)
	assert.Equal(t, passphrase, client.Passphrase)
	assert.Equal(t, BaseURL, client.BaseURL)
	assert.NotNil(t, client.HTTPClient)
}

func TestClient_SetBaseURL(t *testing.T) {
	client := NewClient("", "", "")
	customURL := "https://custom.api.com"

	result := client.SetBaseURL(customURL)

	assert.Equal(t, customURL, client.BaseURL)
	assert.Equal(t, client, result) // Should return self for chaining
}

func TestClient_ServiceFactoryMethods(t *testing.T) {
	client := NewClient("test", "test", "test")

	// Test a few key service factory methods
	assert.NotNil(t, client.NewAccountInfoService())
	assert.NotNil(t, client.NewAccountAssetsService())
	assert.NotNil(t, client.NewPlaceOrderService())
	assert.NotNil(t, client.NewGetTickersService())
	assert.NotNil(t, client.NewTransferService())
}

func TestClient_CreateSignature(t *testing.T) {
	client := NewClient("test_api_key", "test_secret_key", "test_passphrase")

	message := "1234567890GETapi/v3/account/settings"
	signature := client.createSignature(message)

	assert.NotEmpty(t, signature)
	// Should be base64 encoded
	assert.Regexp(t, `^[A-Za-z0-9+/]+=*$`, signature)
}
