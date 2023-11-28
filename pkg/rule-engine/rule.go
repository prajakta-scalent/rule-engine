package ruleengine

import (
	"fmt"
	"reflect"

	validatecondition "github.com/prajakta-scalent/rule-engine/pkg/rule-engine/validate-condition"
)

const CALLBACK_FUNCTION_NAME = "RuleCallback"

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

type RuleEngine interface {
	RegisterGroup()
	Execute(data interface{})
}

// TO DO: need to implement tha save functionality
type RuleExecutionResult interface {
	Save(result []Result) error
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
	var result []Result

	for _, rule := range r.Rules {
		var response bool
		dataValue := getRuleData(rule.Name, data)

		if rule.Condition == "callback" {
			callbackMethod := reflect.ValueOf(dataValue).MethodByName(CALLBACK_FUNCTION_NAME)
			callbackRes := callbackMethod.Call(nil)[0]
			response = callbackRes.Interface().(bool)
		} else {
			response = validatecondition.Validate(rule.Condition, rule.MatchValue, dataValue)
		}
		result = append(result, Result{
			Rule:       rule,
			InputValue: dataValue,
			Outcome:    response,
		})
	}
	//TO DO need to create interface to save result
	fmt.Println("Final Result", result)
}

func getRuleData(name string, data []RuleInput) interface{} {
	for _, dataVal := range data {
		if name == dataVal.RuleName {
			return dataVal.Value
		}
	}
	return nil
}
