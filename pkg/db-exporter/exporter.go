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
	fmt.Printf("\n Rule execution result:\n %+v", result)
	return nil
}
