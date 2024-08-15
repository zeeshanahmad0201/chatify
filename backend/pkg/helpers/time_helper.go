package helpers

import "time"

// GetCurrentTimeInMillis returns millis
func GetCurrentTimeInMillis() int64 {
	return time.Now().UnixMilli()
}
