package common

import (
	"fmt"
	"strconv"
	"strings"
)

func ConvertToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case int64:
		return float64(v), nil
	case string:
		if v == "" {
			return 0.0, nil
		}
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, &ParseError{v, "float64", fmt.Sprintf("Unable to convert string: %s", err.Error())}
		}
		return val, nil
	case nil:
		return 0.0, nil
	default:
		return 0, &ParseError{value, "float64", "Unsupported type"}
	}
}

func ConvertToInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case string:
		if v == "" {
			return 0, nil
		}
		val, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, &ParseError{v, "int64", fmt.Sprintf("Unable to convert string: %s", err.Error())}
		}
		return val, nil
	case nil:
		return 0, nil
	default:
		return 0, &ParseError{value, "int64", "Unsupported type"}
	}
}

func ConvertToBool(value interface{}) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		if v == "" {
			return false, nil
		}
		vLower := strings.ToLower(v)
		switch vLower {
		case "true":
			return true, nil
		case "false":
			return false, nil
		default:
			return false, &ParseError{v, "bool", "Invalid string value"}
		}
	case nil:
		return false, nil
	default:
		return false, &ParseError{value, "bool", "Unsupported type"}
	}
}
