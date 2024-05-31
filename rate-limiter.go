package rateLimiter

import "fmt"

type RateLimiter interface {
	AddRuleAddRule(name string, rule Rule) // add a single rule
	DeleteRule(name string)                // delete a single rule
	SetRules(ruleSet map[string]Rule)      // define a set of rules
	ResetRules()                           // delete all rules
	GetRules() map[string]Rule             // get all rules

	CheckLimit() // check request with rules
	ResetLimit() // reset limit for request
}

type BasicRateLimiter struct {
	rules   map[string]Rule
	clients map[string][]Action
}

func NewBasicRateLimiter(rules map[string]Rule) *BasicRateLimiter {
	rl := &BasicRateLimiter{
		rules: rules,
	}
	return rl
}

func (rl *BasicRateLimiter) AddRule(name string, rule Rule) {
	rl.rules[name] = rule
}

func (rl *BasicRateLimiter) DeleteRule(name string) {
	delete(rl.rules, name)
}

func (rl *BasicRateLimiter) GetRules() map[string]Rule {
	return rl.rules
}

func (rl *BasicRateLimiter) SetRules(ruleSet map[string]Rule) {
	rl.rules = ruleSet
}

func (rl *BasicRateLimiter) ResetRules() {
	rl.rules = nil
}

func (rl *BasicRateLimiter) CheckLimit(cID string, actionType string) error {
	if clientActions, ok := rl.clients[cID]; !ok {
		// check rules slice for action, count avilable for cllient, record it
		for _, rule := range rl.rules {
			if rule.ActionType == actionType {
				rl.clients[cID] = append(rl.clients[cID], Action{
					Type:             actionType,
					AvailableActions: rule.AvailableActions,
				})
			}
		}

	} else {
		// clientID found in slice, check if action is available, change available actions count
		for _, action := range clientActions {

			if action.Type == actionType {
				if action.AvailableActions < 1 {
					return fmt.Errorf(errForbidden)
				} else {
					action.Record()
					return nil
				}
			}
		}
	}

	return nil
}
