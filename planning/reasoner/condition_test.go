package reasoner

import "testing"

func TestEvaluateCondition_GreaterEqual(t *testing.T) {
	signals := map[string]float64{
		"REGIME_SHIFT": 0.75,
	}

	cond := Condition{
		Signal: "REGIME_SHIFT",
		Op:     "gte",
		Value:  0.7,
	}

	ok := evaluateCondition(cond, signals)
	if !ok {
		t.Fatalf("expected condition to be true")
	}
}

func TestEvaluateCondition_False(t *testing.T) {
	signals := map[string]float64{
		"REGIME_SHIFT": 0.4,
	}

	cond := Condition{
		Signal: "REGIME_SHIFT",
		Op:     "gte",
		Value:  0.7,
	}

	ok := evaluateCondition(cond, signals)
	if ok {
		t.Fatalf("expected condition to be false")
	}
}
