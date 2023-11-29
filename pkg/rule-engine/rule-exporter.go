package ruleengine

type ResultExporter interface {
	Save(result []RuleGroupResult) error
}
