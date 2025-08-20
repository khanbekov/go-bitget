package common

import (
	"strconv"
	"time"
)

func TimestampSec() string {
	timesStamp := time.Now().Unix()
	return strconv.FormatInt(timesStamp, 10)
}

func TimestampMs() string {
	timesStamp := time.Now().Unix() * 1000
	return strconv.FormatInt(timesStamp, 10)
}

func SafeStringCast(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}

func SafeFloat64Cast(value interface{}) float64 {
	if f, ok := value.(float64); ok {
		return f
	}
	return 0.0
}
