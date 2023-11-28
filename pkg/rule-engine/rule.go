package ruleengine

type Rule struct {
	IsMandatory  bool
	Name         string
	Condition    string
	MatchValue   interface{}
	MatchAgainst interface{}
	Callback     interface{}
}

type RuleGroup struct {
	Name                string
	Rules               []Rule
	ExecuteConcurrently bool
}

type RuleEngine interface {
	RegisterGroup()
	Execute()
	Save()
}

func New() *RuleGroup {
	return &RuleGroup{}
}

func (r *RuleGroup) RegisterGroup(name string, rules []Rule, executeConcurrently bool) {
	r.Rules = append(r.Rules, rules...)
	r.Name = name
	r.ExecuteConcurrently = executeConcurrently
}

func (r *RuleGroup) Execute() {

}

func (r *RuleGroup) Save() {

}
