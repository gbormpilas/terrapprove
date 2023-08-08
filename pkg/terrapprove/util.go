package terrapprove

import (
	"encoding/json"
	"fmt"
	"github.com/go-yaml/yaml"
	tfjson "github.com/hashicorp/terraform-json"
	"io/ioutil"
)

func ReadPlanFile(path string) tfjson.Plan {
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

func ReadRulesFile(path string) RuleSet {
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
