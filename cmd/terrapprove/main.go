package main

import (
	"flag"
	"fmt"
	"github.com/gbormpilas/terrapprove/pkg/terrapprove"
	"os"
)

func main() {
	var planFile = flag.String("plan", "plan.json", "Path to terraform json plan file")
	var rulesFile = flag.String("rules", "rules.yaml", "Path to rules yaml file")

	flag.Parse()

	plan := terrapprove.ReadPlanFile(*planFile)
	rs := terrapprove.ReadRulesFile(*rulesFile)

	rs.ValidateRules()

	if rs.PlanAllowed(&plan) {
		fmt.Println("I approve")
	} else {
		fmt.Println("Plan not allowed")
		os.Exit(42)
	}
}
