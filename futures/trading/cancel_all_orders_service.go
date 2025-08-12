package trading

import (
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/net/context"
)

// CancelAllOrdersService account info
type CancelAllOrdersService struct {
	c             ClientInterface
	productType   ProductType
	marginCoin    string
	requestTime   string
	receiveWindow string
}

func (s *CancelAllOrdersService) ProductType(productType ProductType) *CancelAllOrdersService {
	s.productType = productType
	return s
}

func (s *CancelAllOrdersService) MarginCoin(marginCoin string) *CancelAllOrdersService {
	s.marginCoin = marginCoin
	return s
}

func (s *CancelAllOrdersService) RequestTime(requestTime string) *CancelAllOrdersService {
	s.requestTime = requestTime
	return s
}

func (s *CancelAllOrdersService) ReceiveWindow(receiveWindow string) *CancelAllOrdersService {
	s.receiveWindow = receiveWindow
	return s
}

func (s *CancelAllOrdersService) Do(ctx context.Context) (cancelResponse *CancelAllOrdersResponse, err error) {
	body := make(map[string]string)
	// Set params of request
	body["productType"] = string(s.productType)
	body["marginCoin"] = s.marginCoin

	if s.requestTime != "" {
		body["requestTime"] = s.requestTime
	}
	if s.receiveWindow != "" {
		body["receiveWindow"] = s.receiveWindow
	}

	// Marshal body to JSON
	bodyBytes, err := jsoniter.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Make request to API
	var res *ApiResponse

	res, _, err = s.c.CallAPI(ctx, "POST", EndpointCancelAllOrders, nil, bodyBytes, true)

	if err != nil {
		return nil, err
	}

	// Unmarshal json from response
	err = jsoniter.Unmarshal(res.Data, &cancelResponse)

	if err != nil {
		return nil, err
	}

	return cancelResponse, nil
}

type CancelAllOrdersResponse struct {
	SuccessList []OrderInfo       `json:"successList"`
	FailureList []OrderInfoFailed `json:"failureList"`
}

// UnmarshalJSON realization interface json.Unmarshaler for CancelAllOrdersResponse
func (c *CancelAllOrdersResponse) UnmarshalJSON(data []byte) error {
	// Создаем временную структуру для парсинга
	type Alias CancelAllOrdersResponse // Для избежания рекурсии
	var tmp struct {
		SuccessList []OrderInfo       `json:"successList"`
		FailureList []OrderInfoFailed `json:"failureList"`
	}

	// Парсим входной JSON в временную структуру
	if err := json.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf("ошибка парсинга data: %w", err)
	}

	// Присваиваем значения основной структуре
	*c = CancelAllOrdersResponse(tmp)
	return nil
}
