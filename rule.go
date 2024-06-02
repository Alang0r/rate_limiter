package ratelimiter

type Rule struct {
	ActionType       string // to set the rule for scope specific action
	AvailableActions uint64 // behavior on match
}
