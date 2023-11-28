package ruleengine

import (
	"fmt"

	validatecondition "github.com/prajakta-scalent/rule-engine/pkg/rule-engine/validate-condition"
)

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

type RuleEngine interface {
	RegisterGroup()
	Execute(data interface{})
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

func (r *RuleGroup) Execute(data []RuleInput) {
	if r.ExecuteConcurrently {
		r.ExecuteRulesConcurrently(data)
	} else {
		r.ExecuteRulesSequentially(data)
	}
}

func (r *RuleGroup) ExecuteRulesConcurrently(data []RuleInput) {
	// for _, rule := range r.Rules {

	// }
}

func (r *RuleGroup) ExecuteRulesSequentially(data []RuleInput) {
	for _, rule := range r.Rules {
		if rule.Condition != "" {
			dataValue := getRuleData(rule.Name, data)
			fmt.Println(rule.Name, dataValue)
			// if dataValue == nil {
			// 	fmt.Println(rule.Name, dataValue)
			// }
			result := validatecondition.Validate(rule.Condition, rule.MatchValue, dataValue)
			fmt.Println("rule response", result)
		}
	}
}

func getRuleData(name string, data []RuleInput) interface{} {
	for _, dataVal := range data {
		if name == dataVal.RuleName {
			return dataVal
		}
	}
	return nil
}

func (r *RuleGroup) Save() {

}
