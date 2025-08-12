package account

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/khanbekov/go-bitget/futures"
)

// SetLeverageService provides methods to set leverage for a trading pair
type SetLeverageService struct {
	c             futures.ClientInterface
	symbol        string
	productType   futures.ProductType
	marginCoin    string
	leverage      string // Optional parameter
	longLeverage  string // Optional parameter
	shortLeverage string // Optional parameter
	holdSide      string // Optional parameter
}

// Symbol sets the trading pair (required)
func (s *SetLeverageService) Symbol(symbol string) *SetLeverageService {
	s.symbol = symbol
	return s
}

// ProductType sets the product type (required)
func (s *SetLeverageService) ProductType(productType futures.ProductType) *SetLeverageService {
	s.productType = productType
	return s
}

// MarginCoin sets the margin currency (required)
func (s *SetLeverageService) MarginCoin(marginCoin string) *SetLeverageService {
	s.marginCoin = marginCoin
	return s
}

// Leverage sets the leverage value (optional)
func (s *SetLeverageService) Leverage(leverage string) *SetLeverageService {
	s.leverage = leverage
	return s
}

// LongLeverage sets the leverage value (optional)
func (s *SetLeverageService) LongLeverage(leverage string) *SetLeverageService {
	s.longLeverage = leverage
	return s
}

// ShortLeverage sets the leverage value (optional)
func (s *SetLeverageService) ShortLeverage(leverage string) *SetLeverageService {
	s.shortLeverage = leverage
	return s
}

// TradeSide sets the position direction for setting leverage (optional)
func (s *SetLeverageService) HoldSide(holdSide string) *SetLeverageService {
	s.holdSide = holdSide
	return s
}

// checkRequiredParams validates required parameters
func (s *SetLeverageService) checkRequiredParams() error {
	if s.symbol == "" {
		return fmt.Errorf("symbol is required")
	}
	if s.productType == "" {
		return fmt.Errorf("productType is required")
	}
	if s.marginCoin == "" {
		return fmt.Errorf("marginCoin is required")
	}

	return nil
}

// Do sends the set leverage request
func (s *SetLeverageService) Do(ctx context.Context) error {
	if err := s.checkRequiredParams(); err != nil {
		return err
	}

	body := s.setLeverageRequestBody()
	bodyBytes, err := jsoniter.Marshal(body)
	if err != nil {
		return err
	}

	res, _, err := s.c.CallAPI(ctx, "POST", futures.EndpointSetLeverage, nil, bodyBytes, true)
	if err != nil {
		return err
	}

	// Check for API success
	if res.Code != "00000" {
		return fmt.Errorf("API error: %s (code %s)", res.Msg, res.Code)
	}
	return nil
}

// setLeverageRequestBody constructs the request payload
func (s *SetLeverageService) setLeverageRequestBody() map[string]string {
	body := make(map[string]string)
	body["symbol"] = s.symbol
	body["productType"] = string(s.productType)
	body["marginCoin"] = s.marginCoin

	if s.holdSide != "" {
		body["holdSide"] = s.holdSide
	}
	if s.leverage != "" {
		body["leverage"] = s.leverage
	}
	if s.longLeverage != "" {
		body["longLeverage"] = s.longLeverage
	}
	if s.shortLeverage != "" {
		body["shortLeverage"] = s.shortLeverage
	}

	return body
}
