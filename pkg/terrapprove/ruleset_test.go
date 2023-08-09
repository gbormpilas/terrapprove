package terrapprove

import (
	"github.com/go-yaml/yaml"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestPlanAllowed(t *testing.T) {
	plan := ReadPlanFile("../../tests/plan.json")
	rs := ReadRulesFile("../../tests/rules.yaml")
	violations, _ := rs.PlanViolations(&plan)
	allowed := len(violations) == 0

	assert.Equal(t, false, allowed, "Plan should not be approved")

	rs = ReadRulesFile("../../tests/allow_list.yaml")
	violations, _ = rs.PlanViolations(&plan)
	allowed = len(violations) == 0

	assert.Equal(t, true, allowed, "Plan should be approved")

	plan = ReadPlanFile("../../tests/pets/plan.json")
	rs = ReadRulesFile("../../tests/pets/disapprove.yaml")
	violations, _ = rs.PlanViolations(&plan)
	allowed = len(violations) == 0

	assert.Equal(t, false, allowed, "Plan should not be approved")

	rs = ReadRulesFile("../../tests/pets/approve.yaml")
	violations, _ = rs.PlanViolations(&plan)
	allowed = len(violations) == 0

	assert.Equal(t, true, allowed, "Plan should be approved")
}

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

	rs.ValidateRules()

	rsTest := RuleSet{
		Rules: []Rule{{
			Provider: "hashicorp/google",
			Resource: "google_container_cluster",
			Actions:  tfjson.Actions{tfjson.ActionDelete},
		},
			{
				Provider: "foo-cloud/foo-provider",
				Resource: "willnotexist",
				Actions:  tfjson.Actions{tfjson.ActionCreate},
			},
			{
				Provider: "registry.terraform.io/hashicorp/local",
				Resource: "local_file",
				Actions:  tfjson.Actions{tfjson.ActionCreate},
			},
		}}
	rsTest.ValidateRules()

	assert.Equal(t, rsTest, rs, "The yaml rules does not agree with the original")
}
