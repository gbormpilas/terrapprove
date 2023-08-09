package terrapprove

import (
	tfjson "github.com/hashicorp/terraform-json"
	"reflect"
)

type RuleSet struct {
	Rules       []Rule `yaml:"rules"`
	IsAllowList bool   `yaml:"isAllowList"`
}

func (rs *RuleSet) ValidateRules() {
	for _, rule := range rs.Rules {
		rule.validateProvider()
	}
}

func (rs RuleSet) PlanViolations(plan *tfjson.Plan) ([]tfjson.ResourceChange, bool) {
	ans := make([]tfjson.ResourceChange, 0)
	if !rs.IsAllowList {
		for _, rule := range rs.Rules {
			if violations := rule.violations(plan); len(violations) > 0 {
				ans = append(ans, violations...)
			}
		}
	}
	if rs.IsAllowList {
		for _, c := range plan.ResourceChanges {
			if !rs.changeAllowed(c) {
				ans = append(ans, *c)
			}
		}
	}

	return ans, rs.IsAllowList
}

func (rs RuleSet) changeAllowed(c *tfjson.ResourceChange) bool {
	for _, rule := range rs.Rules {
		if c.ProviderName == rule.Provider {
			if c.Type == rule.Resource {
				if reflect.DeepEqual(c.Change.Actions, rule.Actions) {
					return true
				}
			}
		}
	}
	return false
}
