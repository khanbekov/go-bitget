package uta

// Product categories for UTA API
const (
	CategorySpot        = "SPOT"
	CategoryMargin      = "MARGIN"
	CategoryUSDTFutures = "USDT-FUTURES"
	CategoryCoinFutures = "COIN-FUTURES"
	CategoryUSDCFutures = "USDC-FUTURES"
)

// Order sides
const (
	SideBuy  = "buy"
	SideSell = "sell"
)

// Order types
const (
	OrderTypeLimit  = "limit"
	OrderTypeMarket = "market"
)

// Time in force options
const (
	TimeInForceGTC      = "gtc"       // Good 'til canceled
	TimeInForceIOC      = "ioc"       // Immediate or cancel
	TimeInForceFOK      = "fok"       // Fill or kill
	TimeInForcePostOnly = "post_only" // Post only (maker only)
)

// Position sides (for futures in hedge mode)
const (
	PositionSideLong  = "long"
	PositionSideShort = "short"
)

// Self Trade Prevention modes
const (
	STPNone        = "none"         // No STP
	STPCancelTaker = "cancel_taker" // Cancel taker order
	STPCancelMaker = "cancel_maker" // Cancel maker order
	STPCancelBoth  = "cancel_both"  // Cancel both orders
)

// Strategy order types
const (
	StrategyTypeTPSL = "tpsl" // Take-Profit and Stop-Loss
)

// TPSL modes
const (
	TPSLModeFull    = "full"    // All position take-profit/stop-loss
	TPSLModePartial = "partial" // Partial position take-profit/stop-loss
)

// Trigger types for strategy orders
const (
	TriggerTypeMarket = "market" // Market price trigger
	TriggerTypeMark   = "mark"   // Mark price trigger
)

// Strategy order execution types
const (
	StrategyOrderTypeLimit  = "limit"  // Limit order execution
	StrategyOrderTypeMarket = "market" // Market order execution
)

// Order status values
const (
	OrderStatusLive            = "live"
	OrderStatusNew             = "new"
	OrderStatusPartiallyFilled = "partially_filled"
	OrderStatusFilled          = "filled"
	OrderStatusCancelled       = "cancelled"
)

// Account types for transfers
const (
	AccountTypeSpot           = "spot"
	AccountTypeP2P            = "p2p"
	AccountTypeUSDTFutures    = "usdt_futures"
	AccountTypeCoinFutures    = "coin_futures"
	AccountTypeUSDCFutures    = "usdc_futures"
	AccountTypeCrossedMargin  = "crossed_margin"
	AccountTypeIsolatedMargin = "isolated_margin"
	AccountTypeUTA            = "uta"
)

// Holding modes
const (
	HoldingModeOneWay = "one_way_mode"
	HoldingModeHedge  = "hedge_mode"
)

// Deduction status
const (
	DeductOn  = "on"
	DeductOff = "off"
)

// Sub-account operations
const (
	SubAccountOpFreeze   = "freeze"
	SubAccountOpUnfreeze = "unfreeze"
)

// Transfer types
const (
	TransferTypeOnChain  = "on_chain"
	TransferTypeInternal = "internal"
)

// Candlestick intervals
const (
	Interval1m  = "1m"
	Interval3m  = "3m"
	Interval5m  = "5m"
	Interval15m = "15m"
	Interval30m = "30m"
	Interval1H  = "1H"
	Interval4H  = "4H"
	Interval6H  = "6H"
	Interval12H = "12H"
	Interval1D  = "1D"
	Interval3D  = "3D"
)

// Candlestick types
const (
	CandlestickTypeMarket = "MARKET"
	CandlestickTypeMark   = "MARK"
	CandlestickTypeIndex  = "INDEX"
)

// API endpoints
const (
	BaseURL = "https://api.bitget.com"

	// Account management endpoints
	EndpointAccountSettings          = "/api/v3/account/settings"
	EndpointAccountAssets            = "/api/v3/account/assets"
	EndpointAccountFundingAssets     = "/api/v3/account/funding-assets"
	EndpointAccountFeeRate           = "/api/v3/account/fee-rate"
	EndpointAccountSetHoldingMode    = "/api/v3/account/set-hold-mode"
	EndpointAccountSetLeverage       = "/api/v3/account/set-leverage"
	EndpointAccountSwitch            = "/api/v3/account/switch"
	EndpointAccountSwitchStatus      = "/api/v3/account/switch-status"
	EndpointAccountTransfer          = "/api/v3/account/transfer"
	EndpointAccountSubTransfer       = "/api/v3/account/sub-transfer"
	EndpointAccountSubTransferRecord = "/api/v3/account/sub-transfer-record"
	EndpointAccountTransferableCoins = "/api/v3/account/transferable-coins"
	EndpointAccountDepositAddress    = "/api/v3/account/deposit-address"
	EndpointAccountDepositRecords    = "/api/v3/account/deposit-records"
	EndpointAccountSubDepositAddress = "/api/v3/account/sub-deposit-address"
	EndpointAccountSubDepositRecords = "/api/v3/account/sub-deposit-records"
	EndpointAccountWithdrawal        = "/api/v3/account/withdrawal"
	EndpointAccountWithdrawalRecords = "/api/v3/account/withdrawal-records"
	EndpointAccountDepositAccount    = "/api/v3/account/deposit-account"
	EndpointAccountFinancialRecords  = "/api/v3/account/financial-records"
	EndpointAccountConvertRecords    = "/api/v3/account/convert-records"
	EndpointAccountDeductInfo        = "/api/v3/account/deduct-info"
	EndpointAccountSwitchDeduct      = "/api/v3/account/switch-deduct"
	EndpointAccountPaymentCoins      = "/api/v3/account/payment-coins"
	EndpointAccountRepayableCoins    = "/api/v3/account/repayable-coins"
	EndpointAccountRepay             = "/api/v3/account/repay"

	// Sub-account management endpoints
	EndpointUserCreateSub    = "/api/v3/user/create-sub"
	EndpointUserSubList      = "/api/v3/user/sub-list"
	EndpointUserFreezeSub    = "/api/v3/user/freeze-sub"
	EndpointUserCreateSubAPI = "/api/v3/user/create-sub-api"
	EndpointUserSubAPIList   = "/api/v3/user/sub-api-list"
	EndpointUserUpdateSubAPI = "/api/v3/user/update-sub-api"
	EndpointUserDeleteSubAPI = "/api/v3/user/delete-sub-api"
	EndpointAccountSubAssets = "/api/v3/account/sub-unified-assets"

	// Trading endpoints
	EndpointTradePlaceOrder          = "/api/v3/trade/place-order"
	EndpointTradeCancelOrder         = "/api/v3/trade/cancel-order"
	EndpointTradeModifyOrder         = "/api/v3/trade/modify-order"
	EndpointTradePlaceBatch          = "/api/v3/trade/place-batch"
	EndpointTradeCancelBatch         = "/api/v3/trade/cancel-batch"
	EndpointTradeBatchModifyOrder    = "/api/v3/trade/batch-modify-order"
	EndpointTradeCancelSymbolOrder   = "/api/v3/trade/cancel-symbol-order"
	EndpointTradeClosePositions      = "/api/v3/trade/close-positions"
	EndpointTradeCountdownCancelAll  = "/api/v3/trade/countdown-cancel-all"
	EndpointTradePlaceStrategyOrder  = "/api/v3/trade/place-strategy-order"
	EndpointTradeCancelStrategyOrder = "/api/v3/trade/cancel-strategy-order"
	EndpointTradeModifyStrategyOrder = "/api/v3/trade/modify-strategy-order"
	EndpointTradeUnfilledOrders      = "/api/v3/trade/unfilled-orders"
	EndpointTradeOrderInfo           = "/api/v3/trade/order-info"
	EndpointTradeHistoryOrders       = "/api/v3/trade/history-orders"
	EndpointTradeFills               = "/api/v3/trade/fills"
	EndpointTradeUnfilledStrategy    = "/api/v3/trade/unfilled-strategy-orders"
	EndpointTradeHistoryStrategy     = "/api/v3/trade/history-strategy-orders"
	EndpointAccountMaxOpenAvailable  = "/api/v3/account/max-open-available"
	EndpointInsLoanLoanOrder         = "/api/v3/ins-loan/loan-order"

	// Position endpoints
	EndpointPositionCurrentPosition = "/api/v3/position/current-position"
	EndpointPositionHistoryPosition = "/api/v3/position/history-position"

	// Market data endpoints
	EndpointMarketTickers        = "/api/v3/market/tickers"
	EndpointMarketCandles        = "/api/v3/market/candles"
	EndpointMarketHistoryCandles = "/api/v3/market/history-candles"
	EndpointMarketOrderbook      = "/api/v3/market/orderbook"

	// General endpoints
	EndpointMarketCurrentFundRate = "/api/v3/market/current-fund-rate"
	EndpointMarketHistoryFundRate = "/api/v3/market/history-fund-rate"
	EndpointMarketInstruments     = "/api/v3/market/instruments"
	EndpointMarketDiscountRate    = "/api/v3/market/discount-rate"
	EndpointMarketMarginLoans     = "/api/v3/market/margin-loans"
	EndpointMarketOpenInterest    = "/api/v3/market/open-interest"
	EndpointMarketOILimit         = "/api/v3/market/oi-limit"
	EndpointMarketProofOfReserves = "/api/v3/market/proof-of-reserves"
	EndpointMarketRiskReserve     = "/api/v3/market/risk-reserve"
	EndpointMarketPositionTier    = "/api/v3/market/position-tier"
	EndpointMarketFills           = "/api/v3/market/fills"

	// Institutional loan endpoints
	EndpointInsLoanTransfered    = "/api/v3/ins-loan/transfered"
	EndpointInsLoanBindUID       = "/api/v3/ins-loan/bind-uid"
	EndpointInsLoanLTV           = "/api/v3/ins-loan/ltv-convert"
	EndpointInsLoanEnsureCoins   = "/api/v3/ins-loan/ensure-coins-convert"
	EndpointInsLoanProductInfos  = "/api/v3/ins-loan/product-infos"
	EndpointInsLoanRepaidHistory = "/api/v3/ins-loan/repaid-history"
	EndpointInsLoanRiskUnit      = "/api/v3/ins-loan/risk-unit"
	EndpointInsLoanSymbols       = "/api/v3/ins-loan/symbols"
)
