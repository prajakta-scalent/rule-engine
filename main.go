package main

import (
	"fmt"

	dbexporter "github.com/prajakta-scalent/rule-engine/pkg/db-exporter"
	ruleengine "github.com/prajakta-scalent/rule-engine/pkg/rule-engine"
	"github.com/prajakta-scalent/rule-engine/services"
)

func main() {
	user := services.User{
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
			MatchValue:  "john",
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

	ruleGroup := ruleengine.RuleGroup{
		Name:                "userRulesGroup",
		Rules:               rules,
		ExecuteConcurrently: true,
	}

	ruleGroup2 := ruleengine.RuleGroup{
		Name:                "userTestRulesGroup",
		Rules:               rules,
		ExecuteConcurrently: true,
	}

	ruleInput := map[string][]ruleengine.Input{
		"userRulesGroup": {
			{
				RuleName: "AgeShouldBeMoreThan",
				Value:    20,
			},
			{
				RuleName: "NameEqualTo",
				Value:    "john",
			},
			{
				RuleName: "BalanceMoreThan",
				Value:    0.00,
			},
			{
				RuleName: "APICallCheckAgeAllowed",
				Value:    user,
			},
		},
		"userTestRulesGroup": {
			{
				RuleName: "AgeShouldBeMoreThan",
				Value:    20,
			},
			{
				RuleName: "NameEqualTo",
				Value:    "john",
			},
			{
				RuleName: "BalanceMoreThan",
				Value:    0.00,
			},
			{
				RuleName: "APICallCheckAgeAllowed",
				Value:    user,
			},
		},
	}

	dbExporter := dbexporter.DBExporter{}

	ruleEngine, err := ruleengine.New(dbExporter)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ruleEngine.RegisterGroup(ruleGroup)
	ruleEngine.RegisterGroup(ruleGroup2)
	executionID, result, err := ruleEngine.Execute(ruleInput)
	fmt.Println("\n\n###########################################", executionID, result, err)
}
