package reasoner

import "testing"

func TestValidateRules_InvalidStatus(t *testing.T) {
	rules := []Rule{
		{
			ID:       "bad-rule",
			Intent:   "test.intent",
			Priority: 1,
			When: ConditionBlock{
				All: []Condition{
					{Signal: "X", Op: "gte", Value: 0.5},
				},
			},
			Then: RuleAction{
				Status:          "strong-singal", // typo
				ConfidenceBoost: 0.3,
			},
		},
	}

	err := ValidateRules(rules)
	if err == nil {
		t.Fatal("expected validation error for invalid status")
	}
}
