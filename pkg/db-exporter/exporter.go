package dbexporter

import (
	"fmt"

	ruleengine "github.com/prajakta-scalent/rule-engine/pkg/rule-engine"
)

type DBExporter struct {
}

func (de DBExporter) Save(executionID string, result []ruleengine.RuleGroupResult) error {
	// TODO: implement db saving functionality
	fmt.Println("Implement db exporter to save result of rule engine here")
	fmt.Println("\n Rule Engine execution ID: ", executionID)
	fmt.Printf("\n Rule execution result:\n ")
	for _, ruleGroupResult := range result {
		fmt.Println("\n Rule Group Name ", ruleGroupResult.Name)
		for _, ruleResult := range ruleGroupResult.RuleResults {
			fmt.Printf("\nName: %s\n Condition: %s\n IsMandatory: %v\n MatchValue: %v\n Value: %v\n Outcome: %v\n Error: %v\n", ruleResult.Name,
				ruleResult.Condition,
				ruleResult.IsMandatory,
				ruleResult.MatchValue,
				ruleResult.Value,
				ruleResult.Outcome,
				ruleResult.Error,
			)
		}
		fmt.Println("rule execution status: ", ruleGroupResult.Status)
	}
	return nil
}
