package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-yaml/yaml"
	tfjson "github.com/hashicorp/terraform-json"
	"io/ioutil"
	"os"
	"reflect"
)

type Rule struct {
	Provider string         `yaml:"provider"`
	Resource string         `yaml:"resource"`
	Actions  tfjson.Actions `yaml:"actions"`
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

func main() {
	var planFile = flag.String("plan", "plan.json", "Path to terraform json plan file")
	var rulesFile = flag.String("rules", "rules.yaml", "Path to rules yaml file")

	flag.Parse()

	// read plan
	file, err := ioutil.ReadFile(*planFile)
	if err != nil {
		panic(fmt.Errorf("Could not read plan file: %v", err))
	}

	plan := tfjson.Plan{}
	err = json.Unmarshal([]byte(file), &plan)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal plan json: %v", err))
	}

	//read rules
	rulefile, err := ioutil.ReadFile(*rulesFile)
	if err != nil {
		panic(fmt.Errorf("Could not read rules file: %v", err))
	}

	var rs RuleSet
	err = yaml.Unmarshal(rulefile, &rs)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal rules yaml: %v", err))
	}

	for _, rule := range rs.Rules {
		if !rule.evaluate(&plan) {
			fmt.Println("Plan not allowed")
			os.Exit(42)
		}
	}
	fmt.Println("I approve")
}
