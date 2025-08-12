package position

import (
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/net/context"
	"net/url"

	"github.com/khanbekov/go-bitget/futures"
)

// SinglePositionService retrieves detailed information for a specific position
type SinglePositionService struct {
	c           futures.ClientInterface
	symbol      string
	productType futures.ProductType
	marginCoin  string
}

func (s *SinglePositionService) Symbol(symbol string) *SinglePositionService {
	s.symbol = symbol
	return s
}

func (s *SinglePositionService) ProductType(productType futures.ProductType) *SinglePositionService {
	s.productType = productType
	return s
}

func (s *SinglePositionService) MarginCoin(marginCoin string) *SinglePositionService {
	s.marginCoin = marginCoin
	return s
}

func (s *SinglePositionService) Do(ctx context.Context) ([]*Position, error) {
	queryParams := url.Values{}

	// Set params of request
	queryParams.Set("symbol", s.symbol)
	queryParams.Set("productType", string(s.productType))
	queryParams.Set("marginCoin", s.marginCoin)

	// Make request to API
	var res *futures.ApiResponse

	res, _, err := s.c.CallAPI(ctx, "GET", futures.EndpointSinglePosition, queryParams, nil, true)

	if err != nil {
		return nil, err
	}

	// Unmarshal json from response
	var positions []*Position
	err = jsoniter.Unmarshal(res.Data, &positions)

	if err != nil {
		return nil, err
	}

	return positions, nil
}