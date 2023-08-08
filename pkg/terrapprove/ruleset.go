package terrapprove

import (
	tfjson "github.com/hashicorp/terraform-json"
)

type RuleSet struct {
	Rules []Rule `yaml:"rules"`
}

func (rs *RuleSet) ValidateRules() {
	for _, rule := range rs.Rules {
		rule.validateProvider()
	}
}

func (rs RuleSet) PlanAllowed(plan *tfjson.Plan) bool {
	for _, rule := range rs.Rules {
		if !rule.evaluate(plan) {
			return false
		}
	}
	return true
}
