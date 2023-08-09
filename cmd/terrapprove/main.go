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

	if violations, isAllowList := rs.PlanViolations(&plan); len(violations) == 0 {
		fmt.Println("I approve")
	} else {
		fmt.Println("Plan not allowed, the following are the found violations:")
		fmt.Printf("This is an allow list: %v\n", isAllowList)
		for _, v := range violations {
			fmt.Printf("%+v action for address: %+v and Provider: %+v\n", v.Change.Actions, v.Address, v.ProviderName)
		}

		os.Exit(42)
	}
}
