package ruleengine

type Result struct {
	Rule       Rule
	InputValue interface{}
	Outcome    bool
}

type Rule struct {
	IsMandatory bool
	Name        string
	Condition   string
	MatchValue  interface{}
}

type RuleInput struct {
	RuleName string
	Value    interface{}
}

type RuleGroup struct {
	Name                string
	Rules               []Rule
	ExecuteConcurrently bool
}
