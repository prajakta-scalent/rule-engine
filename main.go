package main

import (
	"fmt"

	ruleengine "github.com/prajakta-scalent/rule-engine/pkg/rule-engine"
)

func main() {
	user := User{
		Id:  1,
		Age: 18,
	}

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
			Name:        "APICallCheckAgeAllowed",
			Condition:   "callback",
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
		{
			RuleName: "APICallCheckAgeAllowed",
			Value:    user,
		},
	}

	ruleGroup := ruleengine.New()
	ruleGroup.RegisterGroup("UserCondition", rules, false)
	ruleGroup.Execute(ruleInput)
}

type User struct {
	Id  int
	Age int
}

func (u User) RuleCallback() bool {
	fmt.Println("APICallCheckAgeAllowed call callback func called")
	return true
}
