package ruleengine

import (
	"fmt"
	"reflect"

	validatecondition "github.com/prajakta-scalent/rule-engine/pkg/rule-engine/validate-condition"
)

const CALLBACK_FUNCTION_NAME = "RuleCallback"
const CONDITION_CALLBACK = "callback"

type RuleEngine interface {
	RegisterGroup(name string, rules []Rule, executeConcurrently bool)
	Execute(data interface{})
}

// TO DO: need to implement tha save functionality
type RuleExecutionResult interface {
	Save(result []Result) error
}

type RuleEngineImpl struct {
	ruleGroups []RuleGroup
}

func New() *RuleEngineImpl {
	return &RuleEngineImpl{}
}

func (r *RuleEngineImpl) RegisterGroup(ruleGroup RuleGroup) {
	r.ruleGroups = append(r.ruleGroups, ruleGroup)
}

func (r *RuleEngineImpl) Execute(data map[string][]RuleInput) {
	for _, ruleGroup := range r.ruleGroups {
		if ruleGroup.ExecuteConcurrently {
			r.ExecuteRulesConcurrently(ruleGroup, data[ruleGroup.Name])
		} else {
			r.ExecuteRulesSequentially(ruleGroup, data[ruleGroup.Name])
		}
	}
}

// TODO: implement concurrent rule execution if specifed in value
func (r *RuleEngineImpl) ExecuteRulesConcurrently(ruleGroup RuleGroup, data []RuleInput) {}

func (r *RuleEngineImpl) ExecuteRulesSequentially(ruleGroup RuleGroup, data []RuleInput) {
	result := make([]Result, 0, len(ruleGroup.Rules))

	for _, rule := range ruleGroup.Rules {
		var response bool
		dataValue := getRuleData(rule.Name, data)
		// TODO: Handle the condition of no input found for mandatory rule
		if rule.Condition == CONDITION_CALLBACK {
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
