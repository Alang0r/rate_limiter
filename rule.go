package rateLimiter

type Rule struct {
	ActionType string // to set the rule for scope specific action
	ClientID   string // to set the rule for scope specific client
	AvailableActions uint64   // behavior on match
}


