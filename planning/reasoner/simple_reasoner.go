package reasoner

import (
	"errors"
	"time"

	"woodpecker/planning/intents"
)

type SimpleReasoner struct {
	Version string
}

func (r *SimpleReasoner) Evaluate(
	intentID string,
	params map[string]any,
	signals []SignalInput,
) (intents.IntentOutput, error) {

	if intentID == "" {
		return intents.IntentOutput{}, errors.New("intent_id is required")
	}

	if len(signals) == 0 {
		return intents.IntentOutput{}, errors.New("signals are required")
	}

	// ---- naive aggregation (placeholder) ----
	var confidence float64
	usedSignals := make([]intents.SignalUsage, 0, len(signals))

	weight := 1.0 / float64(len(signals))

	for _, s := range signals {
		confidence += s.Value * weight
		usedSignals = append(usedSignals, intents.SignalUsage{
			SignalID: s.SignalID,
			Value:    s.Value,
			Weight:   weight,
		})
	}

	status := intents.StatusWeakSignal
	if confidence > 0.7 {
		status = intents.StatusStrongSignal
	}

	output := intents.IntentOutput{
		Meta: intents.Meta{
			IntentID:  intentID,
			Timestamp: time.Now().UTC(),
			Version:   r.Version,
		},
		Status:     status,
		Confidence: confidence,
		Summary:    "Intent evaluated based on current signal snapshot.",
		Signals:    usedSignals,
		Reasoning: intents.Reasoning{
			Logic: []intents.ReasoningStep{
				{Step: 1, Description: "Aggregated normalized signal values."},
			},
			Explanation: "The intent evaluation is based on the weighted aggregation of provided signals. This is an initial deterministic reasoning step.",
		},
		Guardrails: &intents.Guardrails{
			HumanConfirmationRequired: true,
			ConfidenceCapped:          false,
		},
	}

	return output, nil
}
