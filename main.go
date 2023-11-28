package main

import ruleengine "github.com/prajakta-scalent/rule-engine/pkg/rule-engine"

func main() {
	rules := []ruleengine.Rule{
		{
			Name:       "AgeShouldBeMoreThan",
			Condition:  "lte",
			MatchValue: 18,
		},
		{
			Name: "CheckEvenNumber",
		},
	}
	ruleGroup := ruleengine.New()
	ruleGroup.RegisterGroup("UserCondition", rules, false)
	ruleGroup.Execute()
}
