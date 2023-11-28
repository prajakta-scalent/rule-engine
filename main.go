package main

import (
	ruleengine "github.com/prajakta-scalent/rule-engine/pkg/rule-engine"
)

type UserService struct {
	ruleEnginge *ruleengine.RuleEngine
}

func main() {
	rules := []ruleengine.Rule{
		{
			Name:        "AgeShouldBeMoreThan",
			Condition:   "lte",
			MatchValue:  18,
			IsMandatory: true,
		},
		{
			Name:        "NameEqualTo",
			Condition:   "eq",
			MatchValue:  "prajakta",
			IsMandatory: true,
		},
		{
			Name:        "BalanceMoreThan",
			Condition:   "gt",
			MatchValue:  0.0,
			IsMandatory: true,
		},
		{
			Name:        "CheckEvenNumber",
			IsMandatory: true,
		},
	}

	ruleInput := []ruleengine.RuleInput{
		{
			RuleName: "AgeShouldBeMoreThan",
			Value:    10,
		},
		{
			RuleName: "NameEqualTo",
			Value:    "prajakta",
		},
		{
			RuleName: "BalanceMoreThan",
			Value:    0.00,
		},
	}

	ruleGroup := ruleengine.New()
	ruleGroup.RegisterGroup("UserCondition", rules, false)
	ruleGroup.Execute(ruleInput)
}
