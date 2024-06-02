package ratelimiter

import (
	"fmt"
	"sync"
)

// RateLimiter - ...
type RateLimiter interface {
	AddRuleAddRule(name string, rule Rule) // add a single rule
	DeleteRules(name string)               // delete all rule
	SetRules(ruleSet map[string]Rule)      // define a set of rules
	ResetRules()                           // delete all rules
	GetRules() map[string]Rule             // get all rules

	CheckLimit() // check request with rules
	ResetLimit() // reset limit for request
}

// BasicRateLimiter - basic ratelimiter implementation
type BasicRateLimiter struct {
	rules map[string][]Rule // [clientID][]rule
	mu    sync.Mutex
}

// NewBasicRateLimiter - BasicRateLimiter constructor, pass rules for predeclared set of rules
func NewBasicRateLimiter(rules map[string][]Rule) *BasicRateLimiter {
	// init rules
	rl := &BasicRateLimiter{}
	if rules != nil {
		rl.rules = rules
	} else {
		rules = make(map[string][]Rule, 0)
	}

	return rl
}

// AddRule - add rules to the ratelimiter one by one, by clientID
func (rl *BasicRateLimiter) AddRule(clienID string, rule Rule) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if _, ok := rl.rules[clienID]; ok {
		return fmt.Errorf(errRuleAlreadyExists)
	}

	rl.rules[clienID] = append(rl.rules[clienID], rule)
	return nil

}

// DeleteRules - deletes all rules for provided clientID
func (rl *BasicRateLimiter) DeleteRules(clienID string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	delete(rl.rules, clienID)
}

// GetRules - returns all rules for provided clientID
func (rl *BasicRateLimiter) GetRules(clientID string) []Rule {
	rules, ok := rl.rules[clientID]
	if !ok {
		return nil
	}
	return rules

}

// SetRules - set rules by rules map
func (rl *BasicRateLimiter) SetRules(ruleSet map[string][]Rule) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.rules = ruleSet
}

// ResetRules - reset all rules for all clientIDs
func (rl *BasicRateLimiter) ResetRules() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.rules = make(map[string][]Rule)
}

// CheckLimit - checks for available actions by clientID and actioontype
func (rl *BasicRateLimiter) CheckLimit(clientID string, actionID uint64) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if ruleSet, ok := rl.rules[clientID]; !ok {
		// no rules for client - access granted
		return nil
	} else {
		for _, rule := range ruleSet {
			if rule.ActionID == actionID {
				if rule.AvailableActions > 0 {
					rule.AvailableActions--
					return nil
					// decriment available actions and grant accesss
				} else {
					return fmt.Errorf(errForbidden)
				}
			}
		}
	}
	return nil
}
