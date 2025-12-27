package reasoner

import (
	"fmt"

	"woodpecker/planning/intents"
)

// ValidateRules validates a full ruleset before runtime execution.
// It must be called at load time (fail-fast).
func ValidateRules(rules []Rule) error {
	for i, rule := range rules {
		if err := validateRule(rule); err != nil {
			return fmt.Errorf("rule[%d] (%s): %w", i, rule.ID, err)
		}
	}
	return nil
}

func validateRule(rule Rule) error {
	if rule.ID == "" {
		return fmt.Errorf("id must not be empty")
	}

	if rule.Intent == "" {
		return fmt.Errorf("intent must not be empty")
	}

	if rule.Priority < 0 {
		return fmt.Errorf("priority must be >= 0")
	}

	if rule.Then.ConfidenceBoost < 0 || rule.Then.ConfidenceBoost > 1 {
		return fmt.Errorf("confidence_boost must be between 0 and 1")
	}

	if !isValidIntentStatus(rule.Then.Status) {
		return fmt.Errorf("invalid status '%s'", rule.Then.Status)
	}

	if len(rule.When.All) == 0 && len(rule.When.Any) == 0 {
		return fmt.Errorf("rule must define at least one condition in 'all' or 'any'")
	}

	for _, c := range rule.When.All {
		if err := validateCondition(c); err != nil {
			return err
		}
	}

	for _, c := range rule.When.Any {
		if err := validateCondition(c); err != nil {
			return err
		}
	}

	return nil
}

func validateCondition(c Condition) error {
	if c.Signal == "" {
		return fmt.Errorf("condition signal must not be empty")
	}

	switch c.Op {
	case "gte", "lte", "gt", "lt", "eq":
		return nil
	default:
		return fmt.Errorf("invalid operator '%s'", c.Op)
	}
}

func isValidIntentStatus(status string) bool {
	switch intents.IntentStatus(status) {
	case intents.StatusNotTriggered,
		intents.StatusLowConfidence,
		intents.StatusWeakSignal,
		intents.StatusModerateSignal,
		intents.StatusStrongSignal:
		return true
	default:
		return false
	}
}
