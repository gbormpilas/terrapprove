package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/go-yaml/yaml"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
)

func TestRulesYamlParsing(t *testing.T) {
	//read rules
	rulefile, err := ioutil.ReadFile("../../tests/rules.yaml")
	if err != nil {
		t.Fatalf("could not read rule file")
	}

	var rs RuleSet

	err = yaml.Unmarshal(rulefile, &rs)

	if err != nil {
		t.Errorf("failed to unmarshal yaml")
	}

	rules := []Rule{{
		Provider: "google",
		Resource: "google_container_cluster",
		Actions:  tfjson.Actions{tfjson.ActionDelete},
	},
		{
			Provider: "foo-provider",
			Resource: "willnotexist",
			Actions:  tfjson.Actions{tfjson.ActionCreate},
		},
		{
			Provider: "registry.terraform.io/hashicorp/local",
			Resource: "local_file",
			Actions:  tfjson.Actions{tfjson.ActionCreate},
		},
	}

	assert.Equal(t, rules, rs.Rules, "The yaml rules does not agree with the original")
}
func TestPlanJsonLoad(t *testing.T) {
	// read plan
	file, err := ioutil.ReadFile("../../tests/plan.json")
	if err != nil {
		t.Fatalf("could not read plan file")
	}

	plan := tfjson.Plan{}

	err = json.Unmarshal([]byte(file), &plan)
	if err != nil {
		t.Errorf("Unmarshal failed")
	}
}

func TestPlanAllowed(t *testing.T) {
	plan := readPlanFile("../../tests/plan.json")
	rs := readRulesFile("../../tests/rules.yaml")

	assert.Equal(t, false, rs.planAllowed(&plan), "Plan should not be approved")
}
