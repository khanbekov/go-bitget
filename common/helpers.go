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
