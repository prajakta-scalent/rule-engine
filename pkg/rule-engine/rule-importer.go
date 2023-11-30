package ruleengine

type RuleImporter interface {
	Import(filePath string) ([]Rule, error)
}
