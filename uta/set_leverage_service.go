package uta

import (
	"context"
	"encoding/json"

	"github.com/khanbekov/go-bitget/common"
)

// SetLeverageService sets leverage for trading
type SetLeverageService struct {
	c        ClientInterface
	category *string
	leverage *string
	symbol   *string
	coin     *string
	posSide  *string
}

// Category sets the product category (required)
func (s *SetLeverageService) Category(category string) *SetLeverageService {
	s.category = &category
	return s
}

// Leverage sets the leverage value (required)
func (s *SetLeverageService) Leverage(leverage string) *SetLeverageService {
	s.leverage = &leverage
	return s
}

// Symbol sets the symbol (optional, for futures)
func (s *SetLeverageService) Symbol(symbol string) *SetLeverageService {
	s.symbol = &symbol
	return s
}

// Coin sets the coin (optional, for margin)
func (s *SetLeverageService) Coin(coin string) *SetLeverageService {
	s.coin = &coin
	return s
}

// PositionSide sets the position side (optional, for isolated margin)
func (s *SetLeverageService) PositionSide(posSide string) *SetLeverageService {
	s.posSide = &posSide
	return s
}

// Do executes the set leverage request
func (s *SetLeverageService) Do(ctx context.Context) error {
	if s.category == nil {
		return common.NewMissingParameterError("category")
	}
	if s.leverage == nil {
		return common.NewMissingParameterError("leverage")
	}

	params := map[string]interface{}{
		"category": *s.category,
		"leverage": *s.leverage,
	}

	if s.symbol != nil {
		params["symbol"] = *s.symbol
	}
	if s.coin != nil {
		params["coin"] = *s.coin
	}
	if s.posSide != nil {
		params["posSide"] = *s.posSide
	}

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	_, _, err = s.c.CallAPI(ctx, "POST", EndpointAccountSetLeverage, nil, body, true)
	return err
}
