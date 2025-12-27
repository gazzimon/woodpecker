package reasoner

import "testing"

func TestEvaluateConditionBlock_AllTrue(t *testing.T) {
	signals := map[string]float64{
		"REGIME_SHIFT":            0.8,
		"PROBABILITY_ACCELERATION": 0.7,
	}

	block := ConditionBlock{
		All: []Condition{
			{Signal: "REGIME_SHIFT", Op: "gte", Value: 0.7},
			{Signal: "PROBABILITY_ACCELERATION", Op: "gte", Value: 0.6},
		},
	}

	ok := evaluateConditionBlock(block, signals)
	if !ok {
		t.Fatalf("expected ALL block to be true")
	}
}

func TestEvaluateConditionBlock_AllFalse(t *testing.T) {
	signals := map[string]float64{
		"REGIME_SHIFT":            0.8,
		"PROBABILITY_ACCELERATION": 0.4,
	}

	block := ConditionBlock{
		All: []Condition{
			{Signal: "REGIME_SHIFT", Op: "gte", Value: 0.7},
			{Signal: "PROBABILITY_ACCELERATION", Op: "gte", Value: 0.6},
		},
	}

	ok := evaluateConditionBlock(block, signals)
	if ok {
		t.Fatalf("expected ALL block to be false")
	}
}
