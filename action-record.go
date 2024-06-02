package ratelimiter

import "sync"

type Action struct {
	Type             string
	AvailableActions uint64
	mu               sync.Mutex // parallel access control
}

func (a *Action) Record() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.AvailableActions--
}
