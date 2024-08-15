package helpers

import (
	"context"
	"time"
)

// Creates a context with timeout of 10 seconds by default
func CreateContext(timeout ...time.Duration) (context.Context, context.CancelFunc) {
	var defaultTimeout = time.Second * 10
	if len(timeout) > 0 {
		defaultTimeout = timeout[0]
	}
	return context.WithTimeout(context.Background(), defaultTimeout)
}
