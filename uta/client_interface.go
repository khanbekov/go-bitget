package uta

import (
	"context"
	"net/url"

	"github.com/valyala/fasthttp"
)

// ClientInterface defines the interface for UTA API client operations
type ClientInterface interface {
	// Core API call method
	CallAPI(ctx context.Context, method string, endpoint string, queryParams url.Values, body []byte, sign bool) (*ApiResponse, *fasthttp.ResponseHeader, error)

	// Account management services
	NewAccountInfoService() *AccountInfoService
	NewAccountAssetsService() *AccountAssetsService
	NewAccountFundingAssetsService() *AccountFundingAssetsService
	NewAccountFeeRateService() *AccountFeeRateService
	NewSetHoldingModeService() *SetHoldingModeService
	NewSetLeverageService() *SetLeverageService
	NewSwitchAccountService() *SwitchAccountService
	NewGetSwitchStatusService() *GetSwitchStatusService

	// Transfer services
	NewTransferService() *TransferService
	NewSubTransferService() *SubTransferService
	NewGetTransferRecordsService() *GetTransferRecordsService
	NewGetTransferableCoinsService() *GetTransferableCoinsService

	// Deposit and withdrawal services
	NewGetDepositAddressService() *GetDepositAddressService
	NewGetDepositRecordsService() *GetDepositRecordsService
	NewGetSubDepositAddressService() *GetSubDepositAddressService
	NewGetSubDepositRecordsService() *GetSubDepositRecordsService
	NewWithdrawalService() *WithdrawalService
	NewGetWithdrawalRecordsService() *GetWithdrawalRecordsService
	NewSetDepositAccountService() *SetDepositAccountService

	// Financial records services
	NewGetFinancialRecordsService() *GetFinancialRecordsService
	NewGetConvertRecordsService() *GetConvertRecordsService
	NewGetDeductInfoService() *GetDeductInfoService
	NewSwitchDeductService() *SwitchDeductService
	NewGetPaymentCoinsService() *GetPaymentCoinsService
	NewGetRepayableCoinsService() *GetRepayableCoinsService
	NewRepayService() *RepayService

	// Sub-account management services
	NewCreateSubAccountService() *CreateSubAccountService
	NewGetSubAccountListService() *GetSubAccountListService
	NewFreezeSubAccountService() *FreezeSubAccountService
	NewCreateSubAccountAPIKeyService() *CreateSubAccountAPIKeyService
	NewGetSubAccountAPIKeysService() *GetSubAccountAPIKeysService
	NewModifySubAccountAPIKeyService() *ModifySubAccountAPIKeyService
	NewDeleteSubAccountAPIKeyService() *DeleteSubAccountAPIKeyService
	NewGetSubAccountAssetsService() *GetSubAccountAssetsService

	// Trading services
	NewPlaceOrderService() *PlaceOrderService
	NewCancelOrderService() *CancelOrderService
	NewModifyOrderService() *ModifyOrderService
	NewBatchPlaceOrdersService() *BatchPlaceOrdersService
	NewBatchCancelOrdersService() *BatchCancelOrdersService
	NewBatchModifyOrdersService() *BatchModifyOrdersService
	NewCancelAllOrdersService() *CancelAllOrdersService
	NewCloseAllPositionsService() *CloseAllPositionsService
	NewCountdownCancelAllService() *CountdownCancelAllService

	// Strategy order services
	NewPlaceStrategyOrderService() *PlaceStrategyOrderService
	NewCancelStrategyOrderService() *CancelStrategyOrderService
	NewModifyStrategyOrderService() *ModifyStrategyOrderService
	NewGetUnfilledStrategyOrdersService() *GetUnfilledStrategyOrdersService
	NewGetStrategyOrderHistoryService() *GetStrategyOrderHistoryService

	// Order and position query services
	NewGetOpenOrdersService() *GetOpenOrdersService
	NewGetOrderDetailsService() *GetOrderDetailsService
	NewGetOrderHistoryService() *GetOrderHistoryService
	NewGetFillHistoryService() *GetFillHistoryService
	NewGetCurrentPositionsService() *GetCurrentPositionsService
	NewGetPositionHistoryService() *GetPositionHistoryService
	NewGetMaxOpenAvailableService() *GetMaxOpenAvailableService
	NewGetLoanOrdersService() *GetLoanOrdersService

	// Market data services
	NewGetTickersService() *GetTickersService
	NewGetCandlesticksService() *GetCandlesticksService
	NewGetHistoryCandlesticksService() *GetHistoryCandlesticksService
	NewGetOrderBookService() *GetOrderBookService

	// General market data services
	NewGetCurrentFundingRateService() *GetCurrentFundingRateService
	NewGetFundingRateHistoryService() *GetFundingRateHistoryService
	NewGetInstrumentsService() *GetInstrumentsService
	NewGetDiscountRateService() *GetDiscountRateService
	NewGetMarginLoansService() *GetMarginLoansService
	NewGetOpenInterestService() *GetOpenInterestService
	NewGetOILimitService() *GetOILimitService
	NewGetProofOfReservesService() *GetProofOfReservesService
	NewGetRiskReserveService() *GetRiskReserveService
	NewGetPositionTierService() *GetPositionTierService
	NewGetRecentPublicFillsService() *GetRecentPublicFillsService
}
