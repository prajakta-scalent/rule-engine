package ruleengine

import "errors"

var (
	ErrRuleGroupNameEmpty      = errors.New("group name cannot be empty")
	ErrRulesEmpty              = errors.New("rules cannot be empty")
	ErrRuleExporterEmpty       = errors.New("rule exporter is not provided")
	ErrRuleInputDataEmpty      = errors.New("no rule input provided")
	ErrRuleInputValueEmpty     = errors.New("provided name with empty rule input")
	ErrRuleInputLengthNotEqual = errors.New("length of rules and input data not matching")
)
