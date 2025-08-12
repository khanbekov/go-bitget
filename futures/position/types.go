package position

import (
	"github.com/khanbekov/go-bitget/common/client"
)

// Re-export common types to avoid importing futures package
type (
	ClientInterface = client.ClientInterface
	ApiResponse     = client.ApiResponse
)

// Futures-specific enums (duplicated to avoid import cycle)
type ProductType string

const (
	ProductTypeUSDTFutures ProductType = "USDT-FUTURES"
	ProductTypeCoinFutures ProductType = "COIN-FUTURES"
	ProductTypeUSDCFutures ProductType = "USDC-FUTURES"
)

// API Endpoints for position operations
const (
	EndpointAllPositions     = "/api/v2/mix/position/all-position"     // Get all positions
	EndpointHistoryPositions = "/api/v2/mix/position/history-position" // Get historical positions
	EndpointSinglePosition   = "/api/v2/mix/position/single-position"  // Get single position
	EndpointClosePosition    = "/api/v2/mix/order/close-positions"     // Flash close position
)

// Service Constructor Functions

// NewAllPositionsService creates a new all positions service.
func NewAllPositionsService(client ClientInterface) *AllPositionsService {
	return &AllPositionsService{c: client}
}

// NewSinglePositionService creates a new single position service.
func NewSinglePositionService(client ClientInterface) *SinglePositionService {
	return &SinglePositionService{c: client}
}

// NewHistoryPositionsService creates a new history positions service.
func NewHistoryPositionsService(client ClientInterface) *HistoryPositionsService {
	return &HistoryPositionsService{c: client}
}

// NewClosePositionService creates a new close position service.
func NewClosePositionService(client ClientInterface) *ClosePositionService {
	return &ClosePositionService{c: client}
}