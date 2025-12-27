package reasoner

import "testing"

func TestSimpleReasoner_Evaluate(t *testing.T) {
	r := &SimpleReasoner{Version: "v1"}

	signals := []SignalInput{
		{SignalID: "REGIME_SHIFT", Value: 0.82},
		{SignalID: "PROBABILITY_ACCELERATION", Value: 0.74},
		{SignalID: "CONVICTION_SPIKE", Value: 0.65},
	}

	out, err := r.Evaluate("interpret.regime_state", nil, signals)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.Confidence <= 0 {
		t.Fatalf("expected confidence > 0")
	}

	if out.Status != "strong_signal" {
		t.Fatalf("expected strong_signal, got %s", out.Status)
	}

	if len(out.Signals) != 3 {
		t.Fatalf("expected 3 signals, got %d", len(out.Signals))
	}
}
