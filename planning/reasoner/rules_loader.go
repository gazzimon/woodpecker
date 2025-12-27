package reasoner

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadRulesFromFile(path string) ([]Rule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var ruleset Ruleset
	if err := yaml.Unmarshal(data, &ruleset); err != nil {
		return nil, err
	}

	return ruleset.Rules, nil
}
