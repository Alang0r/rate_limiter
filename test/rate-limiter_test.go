package test

import (
	"fmt"
	ratelimiter "rate_limiter"
	"testing"
)

func TestRateLimiter(t *testing.T) {
	r1 := ratelimiter.Rule{
		ActionID:         1,
		AvailableActions: 10,
	}

	r2 := ratelimiter.Rule{
		ActionID:         2,
		AvailableActions: 0,
	}

	limiter := ratelimiter.NewBasicRateLimiter(nil, nil)

	limiter.AddRule("client1", r1)
	limiter.AddRule("client2", r2)

	if actionsLeft, err := limiter.CheckLimit("client1", 1); err != nil {
		fmt.Printf("error check limit: %s", err.Error())
	} else {
		fmt.Printf("action is approved, available actions: %d", actionsLeft)
	}
}
