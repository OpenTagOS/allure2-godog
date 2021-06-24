package allure

import "time"

func timestampMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
