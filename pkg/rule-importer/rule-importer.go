package ruleimporter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	ruleengine "github.com/prajakta-scalent/rule-engine/pkg/rule-engine"
)

type RuleImporter struct {
}

func (input RuleImporter) Import(filePath string) (result []ruleengine.Rule, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	// Read the file content
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	err = json.Unmarshal(content, &result)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %v", err)
		return nil, err
	}

	return result, nil

}
