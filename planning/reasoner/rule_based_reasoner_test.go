package reasoner

import (
	"testing"

	"woodpecker/planning/intents"
)

func TestRuleBasedReasoner_Evaluate(t *testing.T) {
	rules := []Rule{
		{
			ID:     "regime_weak",
			Intent: "interpret.regime_state",
			When: ConditionBlock{
				All: []Condition{
					{Signal: "REGIME_SHIFT", Op: "gte", Value: 0.5},
				},
			},
			Then: RuleAction{
				Status:          "weak_signal",
				ConfidenceBoost: 0.2,
			},
			Explanation: "Weak regime shift detected",
		},
	}

	r := &RuleBasedReasoner{
		Version: "v1",
		Rules:   rules,
	}

	out, err := r.Evaluate(
		"interpret.regime_state",
		nil,
		[]SignalInput{
			{SignalID: "REGIME_SHIFT", Value: 0.6},
		},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.Status != intents.StatusWeakSignal {
		t.Fatalf("expected weak_signal, got %s", out.Status)
	}

	if out.Confidence <= 0 {
		t.Fatalf("expected confidence > 0")
	}

	if len(out.Reasoning.Logic) != 1 {
		t.Fatalf("expected 1 reasoning step")
	}
}
