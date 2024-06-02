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
	clients map[string]int // [clientID]ruleSet
	rules   [][]Rule
	mu      sync.Mutex
}

// NewBasicRateLimiter - BasicRateLimiter constructor, pass rules for predeclared set of rules
func NewBasicRateLimiter(clients map[string]int, rules [][]Rule) *BasicRateLimiter {
	// init rules
	rl := &BasicRateLimiter{}
	if clients != nil && rules != nil {
		rl.clients = clients
		rl.rules = rules
	} else {
		rl.clients = make(map[string]int)
		rl.rules = make([][]Rule, 0)
	}

	return rl
}

// AddRule - add rules to the ratelimiter one by one, by clientID
func (rl *BasicRateLimiter) AddRule(clienID string, rule Rule) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if ruleSet, ok := rl.clients[clienID]; !ok {
		rl.rules[ruleSet] = make([]Rule, 1)
	} else {
		rl.rules[ruleSet] = append(rl.rules[ruleSet], rule)
	}

	return nil

}

// DeleteRules - deletes all rules for provided clientID
func (rl *BasicRateLimiter) DeleteRules(clienID string) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if ruleSet, ok := rl.clients[clienID]; !ok {
		return fmt.Errorf(errNoRulesForClient)
	} else {
		rl.rules = append(rl.rules[:ruleSet], rl.rules[ruleSet+1:]...)
	}
	return nil
}

// GetRules - returns all rules for provided clientID
func (rl *BasicRateLimiter) GetRules(clientID string) ([]Rule, error) {
	if ruleSet, ok := rl.clients[clientID]; !ok {
		return nil, fmt.Errorf(errNoRulesForClient)
	} else {
		return rl.rules[ruleSet], nil
	}
}

// SetRules - set rules by rules map
func (rl *BasicRateLimiter) SetRules(clientID string, rules []Rule) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if ruleSet, ok := rl.clients[clientID]; !ok {
		//no rules yet
		rl.rules = append(rl.rules, rules)
		rl.clients[clientID] = len(rl.rules) - 1
	} else {
		rl.rules[ruleSet] = append(rl.rules[ruleSet], rules...)
	}
}

// ResetRules - reset all rules for all clientIDs
func (rl *BasicRateLimiter) ResetRules() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.clients = make(map[string]int)
	rl.rules = make([][]Rule, 0)
}

// CheckLimit - checks for available actions by clientID and actioontype
func (rl *BasicRateLimiter) CheckLimit(clientID string, actionID uint64) (uint64, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if ruleSet, ok := rl.clients[clientID]; ok {
		for _, rule := range rl.rules[ruleSet] {
			if rule.ActionID == actionID {
				if rule.AvailableActions > 0 {
					rule.AvailableActions--
					return rule.AvailableActions, nil
					// decriment available actions and grant accesss
				}
				return 0, fmt.Errorf(errForbidden)
			}
		}
	}
	return 0, nil
}
