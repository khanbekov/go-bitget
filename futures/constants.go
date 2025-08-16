package futures

// API Endpoints - All Bitget Futures API v2 endpoints centralized
const (
	// Account Management Endpoints
	EndpointAccountInfo      = "/api/v2/mix/account/account"            // Get single account
	EndpointAccountList      = "/api/v2/mix/account/accounts"           // Get all accounts
	EndpointAccountBills     = "/api/v2/mix/account/bill"               // Get account bills
	EndpointSetMargin        = "/api/v2/mix/account/set-margin"         // Adjust position margin
	EndpointSetLeverage      = "/api/v2/mix/account/set-leverage"       // Change leverage
	EndpointSetMarginMode    = "/api/v2/mix/account/set-margin-mode"    // Change margin mode
	EndpointSetPositionMode  = "/api/v2/mix/account/set-position-mode"  // Change position mode
	EndpointSetAllLeverage   = "/api/v2/mix/account/set-all-leverage"   // Change product line leverage
	EndpointSubAccountAssets = "/api/v2/mix/account/sub-account-assets" // Get subaccount assets
	EndpointInterestHistory  = "/api/v2/mix/account/interest-history"   // Get USDT-M interest history
	EndpointOpenCount        = "/api/v2/mix/account/open-count"         // Get estimated open count
	EndpointSetAutoMargin    = "/api/v2/mix/account/set-auto-margin"    // Set isolated auto margin
	EndpointSetAssetMode     = "/api/v2/mix/account/set-asset-mode"     // Set USDT-M asset mode

	// Market Data Endpoints
	EndpointAllTickers          = "/api/v2/mix/market/tickers"                     // Get all tickers
	EndpointCandlesticks        = "/api/v2/mix/market/candles"                     // Get candlesticks
	EndpointHistoryCandles      = "/api/v2/mix/market/history-candles"             // Get historical candles
	EndpointContracts           = "/api/v2/mix/market/contracts"                   // Get contract config
	EndpointCurrentFundingRate  = "/api/v2/mix/market/current-fund-rate"           // Get current funding rate
	EndpointHistoryFundingRate  = "/api/v2/mix/market/history-fund-rate"           // Get historical funding rates
	EndpointOpenInterest        = "/api/v2/mix/market/open-interest"               // Get open interest
	EndpointTicker              = "/api/v2/mix/market/ticker"                      // Get single ticker
	EndpointSymbolPrice         = "/api/v2/mix/market/symbol-price"                // Get mark/index/market prices
	EndpointRecentTrades        = "/api/v2/mix/market/fills"                       // Get recent transactions
	EndpointOrderbook           = "/api/v2/mix/market/orderbook"                   // Get orderbook
	EndpointOILimit             = "/api/v2/mix/market/oi-limit"                    // Get contract OI limit
	EndpointDiscountRate        = "/api/v2/mix/market/discount-rate"               // Get discount rate
	EndpointHistoryIndexCandles = "/api/v2/mix/market/history-index-candles"       // Get historical index candles
	EndpointHistoryMarkCandles  = "/api/v2/mix/market/history-mark-candles"        // Get historical mark candles
	EndpointExchangeRate        = "/api/v2/mix/market/exchange-rate"               // Get interest exchange rate
	EndpointInterestRateHistory = "/api/v2/mix/market/union-interest-rate-history" // Get interest rate history
	EndpointMergeDepth          = "/api/v2/mix/market/merge-depth"                 // Get merge market depth
	EndpointFundingTime         = "/api/v2/mix/market/funding-time"                // Get next funding time
	EndpointPositionTier        = "/api/v2/mix/market/query-position-lever"        // Get position tier
	EndpointFillsHistory        = "/api/v2/mix/market/fills-history"               // Get history transactions
	EndpointVIPFeeRate          = "/api/v2/mix/market/vip-fee-rate"                // Get VIP fee rate

	// Position Management Endpoints
	EndpointAllPositions     = "/api/v2/mix/position/all-position"     // Get all positions
	EndpointHistoryPositions = "/api/v2/mix/position/history-position" // Get historical positions
	EndpointSinglePosition   = "/api/v2/mix/position/single-position"  // Get single position
	EndpointClosePosition    = "/api/v2/mix/order/close-positions"     // Flash close position
	EndpointADLRank          = "/api/v2/mix/position/adlRank"          // Get position ADL rank

	// Trading Endpoints
	EndpointPlaceOrder        = "/api/v2/mix/order/place-order"         // Place order
	EndpointModifyOrder       = "/api/v2/mix/order/modify-order"        // Modify order
	EndpointCancelOrder       = "/api/v2/mix/order/cancel-order"        // Cancel order
	EndpointCancelAllOrders   = "/api/v2/mix/order/cancel-all-orders"   // Cancel all orders
	EndpointBatchPlaceOrder   = "/api/v2/mix/order/batch-place-order"   // Batch place orders
	EndpointBatchCancelOrders = "/api/v2/mix/order/batch-cancel-orders" // Batch cancel orders
	EndpointOrderDetails      = "/api/v2/mix/order/detail"              // Get order details
	EndpointOrderHistory      = "/api/v2/mix/order/orders-history"      // Get order history
	EndpointPendingOrders     = "/api/v2/mix/order/orders-pending"      // Get pending orders
	EndpointFillHistory       = "/api/v2/mix/order/fill-history"        // Get fill history
	EndpointOrderFills        = "/api/v2/mix/order/fills"               // Get order fill details
	EndpointPlacePlanOrder    = "/api/v2/mix/order/place-plan-order"    // Place trigger order
	EndpointModifyPlanOrder   = "/api/v2/mix/order/modify-plan-order"   // Modify trigger order
	EndpointCancelPlanOrder   = "/api/v2/mix/order/cancel-plan-order"   // Cancel trigger order
	EndpointPendingPlanOrders = "/api/v2/mix/order/orders-plan-pending" // Get pending trigger orders
	EndpointHistoryPlanOrders = "/api/v2/mix/order/orders-plan-history" // Get history trigger orders
	EndpointModifyTPSL        = "/api/v2/mix/order/modify-tpsl-order"   // Modify TP/SL plan order
	EndpointReversal          = "/api/v2/mix/order/click-backhand"      // Reversal order
	EndpointPlacePosTPSL      = "/api/v2/mix/order/place-pos-tpsl"      // Place simultaneous TP/SL
	EndpointPlaceTPSL         = "/api/v2/mix/order/place-tpsl-order"    // Place TP/SL plan order
	EndpointPlanSubOrder      = "/api/v2/mix/order/plan-sub-order"      // Get trigger sub orders
)

// ProductType define product type
type ProductType string

const (
	// ProductTypeUSDTFutures USDT-M Futures, Futures settled in USDT
	ProductTypeUSDTFutures ProductType = "USDT-FUTURES"
	// ProductTypeUSDCFutures USDC-M Futures, Futures settled in USDC
	ProductTypeUSDCFutures ProductType = "USDC-FUTURES"
	// ProductTypeCOINFutures Coin-M Futures, Futures settled in cryptocurrencies
	ProductTypeCOINFutures ProductType = "COIN-FUTURES"
)

type MarginModeType string

const (
	MarginModeIsolated MarginModeType = "ISOLATED"
	MarginModeCrossed  MarginModeType = "CROSSED"
)

type PositionModeType string

const (
	PositionModeOneWay PositionModeType = "one_way_mode"
	PositionModeHedge  PositionModeType = "hedge"
)

type SideType string

const (
	SideTypeBuy  SideType = "BUY"
	SideTypeSell SideType = "SELL"
)

type PositionSideType string

const (
	PositionSideOpen   PositionSideType = "OPEN"
	PositionSideTClose PositionSideType = "CLOSE"
)

type HoldSideType string

const (
	HoldSideLong  HoldSideType = "long"
	HoldSideShort HoldSideType = "short"
)

type TimeInForceType string

const (
	TimeInForceTypeGTC TimeInForceType = "GTC"       // Good Till Cancel
	TimeInForceTypeIOC TimeInForceType = "IOC"       // Immediate or Cancel
	TimeInForceTypeFOK TimeInForceType = "FOK"       // Fill or Kill
	TimeInForceTypeGTX TimeInForceType = "post_only" // Good Till Crossing (Post Only)
)

type ReduceOnlyType string

const (
	ReduceOnlyTypeYes ReduceOnlyType = "YES"
	ReduceOnlyTypeNo  ReduceOnlyType = "NO"
)

type SelfTradePreventionType string

const (
	SelfTradePreventionNone        SelfTradePreventionType = "NONE"
	SelfTradePreventionCancelTaker SelfTradePreventionType = "CANCEL_TAKER"
	SelfTradePreventionCancelMaker SelfTradePreventionType = "CANCEL_MAKER"
	SelfTradePreventionCancelBoth  SelfTradePreventionType = "CANCEL_BOTH"
)

type OrderType string

const (
	OrderTypeMarket OrderType = "MARKET"
	OrderTypeLimit  OrderType = "LIMIT"
)

// Type aliases for services in subdirectories to avoid import cycles
// These will be updated during the client integration phase
