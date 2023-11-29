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
		if len(ruleGroup.Rules) != len(data[ruleGroup.Name]) {
			return executionID, nil, errors.New("rules and input data not matching for group:" + ruleGroup.Name)
		}
		pass := true
		var ruleResults []RuleResult
		if ruleGroup.ExecuteConcurrently {
			ruleResults = r.ExecuteRulesConcurrently(ruleGroup, data[ruleGroup.Name])
			for _, ruleResult := range ruleResults {
				if !ruleResult.Outcome && ruleResult.IsMandatory {
					pass = false
				}
			}
		} else {
			ruleResults = r.ExecuteRulesSequentially(ruleGroup, data[ruleGroup.Name])
			for _, ruleResult := range ruleResults {
				if !ruleResult.Outcome && ruleResult.IsMandatory {
					pass = false
				}
			}
		}
		result = append(result, RuleGroupResult{
			Name:        ruleGroup.Name,
			RuleResults: ruleResults,
			Status:      pass,
		})
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
		var err error

		dataValue, err := getRuleData(rule.Name, data)
		if err != nil {
			result = append(result, RuleResult{
				Name:        rule.Name,
				Condition:   rule.Condition,
				IsMandatory: rule.IsMandatory,
				MatchValue:  rule.MatchValue,
				Value:       dataValue,
				Outcome:     response,
				Error:       err,
			})
			break
		}

		wg.Add(1)
		go func(rule Rule) {
			defer wg.Done()
			if rule.Condition == CONDITION_CALLBACK {
				callbackMethod := reflect.ValueOf(dataValue).MethodByName(CALLBACK_FUNCTION_NAME)
				callbackRes := callbackMethod.Call(nil)

				if len(callbackRes) > 0 {
					response = callbackRes[0].Interface().(bool)
					if callbackRes[1].Interface() != nil {
						err = callbackRes[1].Interface().(error)
					}
				} else {
					err = errors.New("error while fetching call back response")
				}
			} else {
				response, err = validatecondition.Validate(rule.Condition, rule.MatchValue, dataValue)
			}

			resultCh <- RuleResult{
				Name:        rule.Name,
				Condition:   rule.Condition,
				IsMandatory: rule.IsMandatory,
				MatchValue:  rule.MatchValue,
				Value:       dataValue,
				Outcome:     response,
				Error:       err,
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
		var err error

		dataValue, err := getRuleData(rule.Name, data)
		if err != nil {
			result = append(result, RuleResult{
				Name:        rule.Name,
				Condition:   rule.Condition,
				IsMandatory: rule.IsMandatory,
				MatchValue:  rule.MatchValue,
				Value:       dataValue,
				Outcome:     response,
				Error:       err,
			})
			break
		}

		if rule.Condition == CONDITION_CALLBACK {
			callbackMethod := reflect.ValueOf(dataValue).MethodByName(CALLBACK_FUNCTION_NAME)
			callbackRes := callbackMethod.Call(nil)
			if len(callbackRes) > 0 {
				response = callbackRes[0].Interface().(bool)
				if callbackRes[1].Interface() != nil {
					err = callbackRes[1].Interface().(error)
				}
			} else {
				err = errors.New("error while fetching call back response")
			}
		} else {
			response, err = validatecondition.Validate(rule.Condition, rule.MatchValue, dataValue)
		}

		result = append(result, RuleResult{
			Name:        rule.Name,
			Condition:   rule.Condition,
			IsMandatory: rule.IsMandatory,
			MatchValue:  rule.MatchValue,
			Value:       dataValue,
			Outcome:     response,
			Error:       err,
		})
	}

	return result
}

func getRuleData(name string, data []Input) (interface{}, error) {
	for _, dataVal := range data {
		if name == dataVal.RuleName {
			return dataVal.Value, nil
		}
	}

	return nil, errors.New("no data found for " + name)
}
