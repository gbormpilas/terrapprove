package terrapprove

import (
	"fmt"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-registry-address"
	"reflect"
)

type Rule struct {
	Provider string         `yaml:"provider"`
	Resource string         `yaml:"resource"`
	Actions  tfjson.Actions `yaml:"actions"`
}

func (rule *Rule) validateProvider() {
	provider := tfaddr.MustParseProviderSource(rule.Provider)
	if !provider.HasKnownNamespace() {
		panic(fmt.Sprintf("Could not parse provider source %v: unknown namespace", rule.Provider))
	}
	rule.Provider = provider.String()
}

func (rule Rule) violations(plan *tfjson.Plan) []tfjson.ResourceChange {
	ans := make([]tfjson.ResourceChange, 0)
	for _, c := range plan.ResourceChanges {
		if c.ProviderName == rule.Provider {
			if c.Type == rule.Resource {
				if reflect.DeepEqual(c.Change.Actions, rule.Actions) {
					ans = append(ans, *c)
				}
			}
		}
	}
	return ans
}
