package account

import (
	"context"
	"fmt"
	"net/url"

	"github.com/khanbekov/go-bitget/futures"
	jsoniter "github.com/json-iterator/go"
)

// BillResponse represents the account bill data
type BillResponse struct {
	Symbol    string `json:"symbol"`
	StartUnit string `json:"startUnit"`
	EndUnit   string `json:"endUnit"`
}

// GetAccountBillService provides methods to retrieve account bill information
type GetAccountBillService struct {
	c         futures.ClientInterface
	symbol    string
	startUnit string
	endUnit   string
}

// Symbol sets the trading pair (required)
func (s *GetAccountBillService) Symbol(symbol string) *GetAccountBillService {
	s.symbol = symbol
	return s
}

// StartUnit sets the start of query range (optional)
func (s *GetAccountBillService) StartUnit(startUnit string) *GetAccountBillService {
	s.startUnit = startUnit
	return s
}

// EndUnit sets the end of query range (optional)
func (s *GetAccountBillService) EndUnit(endUnit string) *GetAccountBillService {
	s.endUnit = endUnit
	return s
}

// checkRequiredParams validates required parameters
func (s *GetAccountBillService) checkRequiredParams() error {
	if s.symbol == "" {
		return fmt.Errorf("symbol is required")
	}
	return nil
}

// Do sends the account bill request
func (s *GetAccountBillService) Do(ctx context.Context) (bill *BillResponse, err error) {
	if err := s.checkRequiredParams(); err != nil {
		return nil, err
	}

	queryParams := url.Values{}
	queryParams.Add("symbol", s.symbol)
	if s.startUnit != "" {
		queryParams.Add("startUnit", s.startUnit)
	}
	if s.endUnit != "" {
		queryParams.Add("endUnit", s.endUnit)
	}

	res, _, err := s.c.CallAPI(ctx, "GET", futures.EndpointAccountBills, queryParams, nil, true)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Data BillResponse `json:"data"`
	}
	if err := jsoniter.Unmarshal(res.Data, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}
