package terrapprove

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
)

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
