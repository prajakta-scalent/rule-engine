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

type ruleEngineImpl struct {
	ruleGroups     []RuleGroup
	resultExporter ResultExporter
}

func New(ruleExporter ResultExporter) (*ruleEngineImpl, error) {
	if ruleExporter == nil {
		return nil, ErrRuleExporterEmpty
	}
	return &ruleEngineImpl{
		resultExporter: ruleExporter,
	}, nil
}

func (r *ruleEngineImpl) RegisterGroup(ruleGroup RuleGroup) error {
	if ruleGroup.Name == "" {
		return ErrRuleGroupNameEmpty
	} else if len(ruleGroup.Rules) == 0 {
		return ErrRulesEmpty
	}

	r.ruleGroups = append(r.ruleGroups, ruleGroup)

	return nil
}

func (r *ruleEngineImpl) Execute(data map[string][]Input) (executionID string, result []RuleGroupResult, err error) {
	executionID = uuid.NewString()

	if data == nil {
		return executionID, nil, ErrRuleInputDataEmpty
	}

	for _, ruleGroup := range r.ruleGroups {
		inputData := data[ruleGroup.Name]

		if inputData == nil {
			return executionID, nil, ErrRuleInputValueEmpty
		}
		if len(ruleGroup.Rules) != len(inputData) {
			return executionID, nil, ErrRuleInputLengthNotEqual
		}

		pass := true
		var ruleResults []RuleResult
		if ruleGroup.ExecuteConcurrently {
			ruleResults = r.ExecuteRulesConcurrently(ruleGroup, inputData)
			for _, ruleResult := range ruleResults {
				if !ruleResult.Outcome && ruleResult.IsMandatory {
					pass = false
				}
			}
		} else {
			ruleResults = r.ExecuteRulesSequentially(ruleGroup, inputData)
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

func (r *ruleEngineImpl) ExecuteRulesConcurrently(ruleGroup RuleGroup, data []Input) (result []RuleResult) {
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

func (r *ruleEngineImpl) ExecuteRulesSequentially(ruleGroup RuleGroup, data []Input) (result []RuleResult) {
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
