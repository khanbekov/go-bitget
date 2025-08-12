package types

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// APIError define API error when response status is 4xx or 5xx
type APIError struct {
	Code        int64  `json:"code"`
	Message     string `json:"msg"`
	RequestTime int64  `json:"requestTime"` // Добавлено поле для requestTime
	Data        any    `json:"data"`        // Добавлено поле для data
	Response    []byte `json:"-"`
}

// Error return error code and message
func (e APIError) Error() string {
	if e.IsValid() {
		return fmt.Sprintf("<APIError> code=%d, msg=%s", e.Code, e.Message)
	}
	return fmt.Sprintf("<APIError> rsp=%s", string(e.Response))
}

func (e APIError) IsValid() bool {
	return e.Code != 0 || e.Message != ""
}

// IsAPIError check if e is an API error
func IsAPIError(e error) bool {
	_, ok := e.(*APIError)
	return ok
}

// UnmarshalJSON для корректного парсинга
func (e *APIError) UnmarshalJSON(data []byte) error {
	type Alias APIError // Избегаем рекурсии

	// Временная структура для парсинга
	var tmp struct {
		Code        string          `json:"code"`
		Message     string          `json:"msg"`
		RequestTime int64           `json:"requestTime"`
		Data        json.RawMessage `json:"data"` // Используем RawMessage для произвольного data
	}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		e.Response = data
		return err
	}

	// Преобразуем код из строки в int64
	code, err := strconv.ParseInt(tmp.Code, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid code format: %s", tmp.Code)
	}

	// Копируем данные в структуру
	e.Code = code
	e.Message = tmp.Message
	e.RequestTime = tmp.RequestTime
	e.Data = tmp.Data // Можно оставить как RawMessage или распарсить дальше

	return nil
}
