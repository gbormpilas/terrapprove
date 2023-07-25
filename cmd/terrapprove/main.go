package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-yaml/yaml"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-registry-address"
	"io/ioutil"
	"os"
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

func (rule Rule) evaluate(plan *tfjson.Plan) bool {
	for _, c := range plan.ResourceChanges {
		if c.ProviderName == rule.Provider {
			if c.Type == rule.Resource {
				if reflect.DeepEqual(c.Change.Actions, rule.Actions) {
					return false
				}
			}
		}
	}
	return true
}

type RuleSet struct {
	Rules []Rule `yaml:"rules"`
}

func (rs *RuleSet) validateRules() {
	for _, rule := range rs.Rules {
		rule.validateProvider()
	}
}

func main() {
	var planFile = flag.String("plan", "plan.json", "Path to terraform json plan file")
	var rulesFile = flag.String("rules", "rules.yaml", "Path to rules yaml file")

	flag.Parse()

	plan := readPlanFile(*planFile)
	rs := readRulesFile(*rulesFile)

	rs.validateRules()

	if rs.planAllowed(&plan) {
		fmt.Println("I approve")
	} else {
		fmt.Println("Plan not allowed")
		os.Exit(42)
	}
}

func (rs RuleSet) planAllowed(plan *tfjson.Plan) bool {
	for _, rule := range rs.Rules {
		if !rule.evaluate(plan) {
			return false
		}
	}
	return true
}

func readPlanFile(path string) tfjson.Plan {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("Could not read plan file: %v", err))
	}

	plan := tfjson.Plan{}
	err = json.Unmarshal([]byte(file), &plan)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal plan json: %v", err))
	}
	return plan
}

func readRulesFile(path string) RuleSet {
	rulefile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("Could not read rules file: %v", err))
	}

	var rs RuleSet
	err = yaml.Unmarshal(rulefile, &rs)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal rules yaml: %v", err))
	}
	return rs
}
