package reasoner

import (
	"time"

	"woodpecker/planning/intents"
)

type RuleBasedReasoner struct {
	Version string
	Rules   []Rule
}

func (r *RuleBasedReasoner) Evaluate(
	intentID string,
	params map[string]any,
	signals []SignalInput,
) (intents.IntentOutput, error) {

	signalMap := make(map[string]float64)
	for _, s := range signals {
		signalMap[s.SignalID] = s.Value
	}

	matchedRules, _ := EvaluateRules(intentID, signalMap, r.Rules)

	status := intents.StatusNotTriggered
	confidence := 0.0
	reasonSteps := []intents.ReasoningStep{}

	for i, rule := range matchedRules {
		status = intents.IntentStatus(rule.Then.Status)
		confidence += rule.Then.ConfidenceBoost

		reasonSteps = append(reasonSteps, intents.ReasoningStep{
			Step:        i + 1,
			Description: rule.Explanation,
		})
	}

	if confidence > 1 {
		confidence = 1
	}

	return intents.IntentOutput{
		Meta: intents.Meta{
			IntentID:  intentID,
			Timestamp: time.Now().UTC(),
			Version:   r.Version,
		},
		Status:     status,
		Confidence: confidence,
		Summary:    "Intent evaluated using declarative rule engine.",
		Signals:    mapSignals(signals),
		Reasoning: intents.Reasoning{
			Logic:       reasonSteps,
			Explanation: "One or more declarative rules matched the current signal snapshot.",
		},
		Guardrails: &intents.Guardrails{
			HumanConfirmationRequired: true,
		},
	}, nil
}

func mapSignals(inputs []SignalInput) []intents.SignalUsage {
	w := 1.0 / float64(len(inputs))
	out := make([]intents.SignalUsage, 0, len(inputs))
	for _, s := range inputs {
		out = append(out, intents.SignalUsage{
			SignalID: s.SignalID,
			Value:    s.Value,
			Weight:   w,
		})
	}
	return out
}
