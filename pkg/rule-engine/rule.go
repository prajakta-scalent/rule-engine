package ruleengine

type Rule struct {
	IsMandatory bool
	Name        string
	Condition   string
	MatchValue  interface{}
}

type RuleGroup struct {
	Name                string
	Rules               []Rule
	ExecuteConcurrently bool
}

type Input struct {
	RuleName string
	Value    interface{}
}

type RuleResult struct {
	Name        string
	Condition   string
	IsMandatory bool
	MatchValue  interface{}
	Value       interface{}
	Outcome     bool
	Error       error
}

type RuleGroupResult struct {
	Name        string
	Status      bool
	RuleResults []RuleResult
}
