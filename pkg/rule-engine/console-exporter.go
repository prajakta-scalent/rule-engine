package ruleengine

import "fmt"

type ConsoleExporter struct {
}

func (ce ConsoleExporter) Save(result []RuleGroupResult) error {
	for _, ruleGroupResult := range result {
		fmt.Println("Execution ID: ", ruleGroupResult.ExecutionID)
		for _, ruleResult := range ruleGroupResult.RuleResults {
			fmt.Printf("Name: %s\n Condition: %s\n IsMandatory: %v\n MatchValue: %v\n Value: %v\n Outcome: %v\n Error: %v\n", ruleResult.Name,
				ruleResult.Condition,
				ruleResult.IsMandatory,
				ruleResult.MatchValue,
				ruleResult.Value,
				ruleResult.Outcome,
				ruleResult.Error,
			)
		}
	}
	return nil
}
