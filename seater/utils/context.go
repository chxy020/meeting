package utils

import (
	"context"
	"time"
)

// TimeoutContext create a new timeout context for rpc call
func TimeoutContext(a ...time.Duration) (context.Context, context.CancelFunc) {
	var timeout time.Duration
	if len(a) > 0 {
		timeout = a[0]
	} else {
		timeout = 2 * time.Second
	}
	return context.WithTimeout(context.Background(), timeout)
}
