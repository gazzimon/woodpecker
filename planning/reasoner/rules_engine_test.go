package reasoner

import "testing"

func TestEvaluateRules_Match(t *testing.T) {
	signals := map[string]float64{
		"REGIME_SHIFT":             0.8,
		"PROBABILITY_ACCELERATION": 0.7,
	}

	rules := []Rule{
		{
			ID:     "regime_strong",
			Intent: "interpret.regime_state",
			When: ConditionBlock{
				All: []Condition{
					{Signal: "REGIME_SHIFT", Op: "gte", Value: 0.7},
					{Signal: "PROBABILITY_ACCELERATION", Op: "gte", Value: 0.6},
				},
			},
			Then: RuleAction{
				Status:          "strong_signal",
				ConfidenceBoost: 0.3,
			},
			Explanation: "Strong regime shift detected",
		},
	}

	matched, err := EvaluateRules("interpret.regime_state", signals, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(matched) != 1 {
		t.Fatalf("expected 1 matched rule, got %d", len(matched))
	}
}
