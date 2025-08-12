package uta

import (
	"context"
	"encoding/json"

	"github.com/khanbekov/go-bitget/common"
)

// SetHoldingModeService sets the position holding mode
type SetHoldingModeService struct {
	c           ClientInterface
	holdingMode *string
}

// HoldingMode sets the holding mode (required)
// Values: "one_way_mode" or "hedge_mode"
func (s *SetHoldingModeService) HoldingMode(holdingMode string) *SetHoldingModeService {
	s.holdingMode = &holdingMode
	return s
}

// Do executes the set holding mode request
func (s *SetHoldingModeService) Do(ctx context.Context) error {
	if s.holdingMode == nil {
		return common.NewMissingParameterError("holdingMode")
	}

	params := map[string]interface{}{
		"holdMode": *s.holdingMode,
	}

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	_, _, err = s.c.CallAPI(ctx, "POST", EndpointAccountSetHoldingMode, nil, body, true)
	return err
}
