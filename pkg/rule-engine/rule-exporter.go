package ruleengine

type ResultExporter interface {
	Save(executionID string, result []RuleGroupResult) error
}
