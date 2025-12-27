package intents

import (
	"testing"
	"time"
)

func TestIntentOutput_ValidateBasic_OK(t *testing.T) {
	out := IntentOutput{
		Meta: Meta{
			IntentID:  "interpret.regime_state",
			Timestamp: time.Now().UTC(),
			Version:   "v1",
		},
		Status:     StatusStrongSignal,
		Confidence: 0.73,
		Summary:    "Test summary",
		Signals: []SignalUsage{
			{SignalID: "REGIME_SHIFT", Value: 0.8, Weight: 0.5},
			{SignalID: "PROBABILITY_ACCELERATION", Value: 0.7, Weight: 0.5},
		},
		Reasoning: Reasoning{
			Logic: []ReasoningStep{
				{Step: 1, Description: "Test logic"},
			},
			Explanation: "Test explanation",
		},
	}

	if err := out.ValidateBasic(); err != nil {
		t.Fatalf("expected valid intent output, got error: %v", err)
	}
}
