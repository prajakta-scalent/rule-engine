package ruleengine

import (
	"errors"
	"reflect"
	"sync"

	"github.com/google/uuid"
	validatecondition "github.com/prajakta-scalent/rule-engine/pkg/rule-engine/validate-condition"
)

const CALLBACK_FUNCTION_NAME = "RuleCallback"
const CONDITION_CALLBACK = "callback"

type RuleEngine interface {
	RegisterGroup(name string, rules []Rule, executeConcurrently bool)
	Execute(data interface{}) (bool, error)
}

type FailedRules struct {
	Rules []RuleResult
}

type RuleEngineImpl struct {
	ruleGroups     []RuleGroup
	resultExporter ResultExporter
}

func New(ruleExporter ResultExporter) (*RuleEngineImpl, error) {
	if ruleExporter == nil {
		return nil, errors.New("rule exporter dependency not met")
	}
	return &RuleEngineImpl{
		resultExporter: ruleExporter,
	}, nil
}

func (r *RuleEngineImpl) RegisterGroup(ruleGroup RuleGroup) {
	r.ruleGroups = append(r.ruleGroups, ruleGroup)
}

func (r *RuleEngineImpl) Execute(data map[string][]Input) (executionID string, result []RuleGroupResult, err error) {
	executionID = uuid.NewString()

	for _, ruleGroup := range r.ruleGroups {
		var pass bool
		if ruleGroup.ExecuteConcurrently {
			ruleResults := r.ExecuteRulesConcurrently(ruleGroup, data[ruleGroup.Name])
			for _, ruleResult := range ruleResults {
				if ruleResult.Error != nil && ruleResult.Outcome {
					pass = true
				}
			}
			result = append(result, RuleGroupResult{
				Name:        ruleGroup.Name,
				RuleResults: ruleResults,
				Status:      pass,
			})
		} else {
			ruleResults := r.ExecuteRulesSequentially(ruleGroup, data[ruleGroup.Name])
			for _, ruleResult := range ruleResults {
				if ruleResult.Error != nil && ruleResult.Outcome {
					pass = true
				}
			}
			result = append(result, RuleGroupResult{
				Name:        ruleGroup.Name,
				RuleResults: ruleResults,
				Status:      pass,
			})
		}
	}

	// Log error
	err = r.resultExporter.Save(executionID, result)
	return executionID, result, err
}

func (r *RuleEngineImpl) ExecuteRulesConcurrently(ruleGroup RuleGroup, data []Input) (result []RuleResult) {
	var wg sync.WaitGroup
	resultCh := make(chan RuleResult)

	for _, rule := range ruleGroup.Rules {
		var response bool
		dataValue := getRuleData(rule.Name, data)
		wg.Add(1)
		go func(rule Rule) {
			defer wg.Done()
			if rule.Condition == CONDITION_CALLBACK {
				callbackMethod := reflect.ValueOf(dataValue).MethodByName(CALLBACK_FUNCTION_NAME)
				callbackRes := callbackMethod.Call(nil)[0]
				response = callbackRes.Interface().(bool)
			} else {
				response = validatecondition.Validate(rule.Condition, rule.MatchValue, dataValue)
			}
			resultCh <- RuleResult{
				Name:        rule.Name,
				Condition:   rule.Condition,
				IsMandatory: rule.IsMandatory,
				MatchValue:  rule.MatchValue,
				Value:       dataValue,
				Outcome:     response,
			}
		}(rule)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	for res := range resultCh {
		result = append(result, res)
	}

	return
}

func (r *RuleEngineImpl) ExecuteRulesSequentially(ruleGroup RuleGroup, data []Input) (result []RuleResult) {
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
		// TODO: Need to save error if any comes during execution
		result = append(result, RuleResult{
			Name:        rule.Name,
			Condition:   rule.Condition,
			IsMandatory: rule.IsMandatory,
			MatchValue:  rule.MatchValue,
			Value:       dataValue,
			Outcome:     response,
		})
	}
	return result
}

func getRuleData(name string, data []Input) interface{} {
	for _, dataVal := range data {
		if name == dataVal.RuleName {
			return dataVal.Value
		}
	}
	return nil
}
