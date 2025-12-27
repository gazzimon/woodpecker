package reasoner

import "woodpecker/planning/intents"

type IntentReasoner interface {
	Evaluate(
		intentID string,
		params map[string]any,
		signals []SignalInput,
	) (intents.IntentOutput, error)
}

type SignalInput struct {
	SignalID string
	Value    float64
}

