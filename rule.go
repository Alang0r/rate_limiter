package ratelimiter

// Rule - represents individual rule with acvailable actions count for some actionID
type Rule struct {
	ActionID         uint64 // to set the rule for scope specific action
	AvailableActions uint64 // behavior on match
}
