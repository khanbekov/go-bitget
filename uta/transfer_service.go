package uta

import (
	"context"
	"encoding/json"

	"github.com/khanbekov/go-bitget/common"
)

// TransferService handles internal account transfers
type TransferService struct {
	c        ClientInterface
	fromType *string
	toType   *string
	amount   *string
	coin     *string
	symbol   *string
}

// FromType sets the source account type (required)
func (s *TransferService) FromType(fromType string) *TransferService {
	s.fromType = &fromType
	return s
}

// ToType sets the destination account type (required)
func (s *TransferService) ToType(toType string) *TransferService {
	s.toType = &toType
	return s
}

// Amount sets the transfer amount (required)
func (s *TransferService) Amount(amount string) *TransferService {
	s.amount = &amount
	return s
}

// Coin sets the coin to transfer (required)
func (s *TransferService) Coin(coin string) *TransferService {
	s.coin = &coin
	return s
}

// Symbol sets the symbol (optional, for isolated margin)
func (s *TransferService) Symbol(symbol string) *TransferService {
	s.symbol = &symbol
	return s
}

// Do executes the transfer request
func (s *TransferService) Do(ctx context.Context) (*TransferResult, error) {
	if s.fromType == nil {
		return nil, common.NewMissingParameterError("fromType")
	}
	if s.toType == nil {
		return nil, common.NewMissingParameterError("toType")
	}
	if s.amount == nil {
		return nil, common.NewMissingParameterError("amount")
	}
	if s.coin == nil {
		return nil, common.NewMissingParameterError("coin")
	}

	params := map[string]interface{}{
		"fromType": *s.fromType,
		"toType":   *s.toType,
		"amount":   *s.amount,
		"coin":     *s.coin,
	}

	if s.symbol != nil {
		params["symbol"] = *s.symbol
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	res, _, err := s.c.CallAPI(ctx, "POST", EndpointAccountTransfer, nil, body, true)
	if err != nil {
		return nil, err
	}

	var transferResult TransferResult
	if err := common.UnmarshalJSON(res.Data, &transferResult); err != nil {
		return nil, err
	}

	return &transferResult, nil
}
