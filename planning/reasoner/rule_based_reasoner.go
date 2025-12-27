package reasoner

import (
	"sort"
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

	// 1️⃣ Normalizar señales
	signalMap := make(map[string]float64)
	for _, s := range signals {
		signalMap[s.SignalID] = s.Value
	}

	// 2️⃣ Evaluar reglas
	matchedRules, err := EvaluateRules(intentID, signalMap, r.Rules)
	if err != nil {
		return intents.IntentOutput{}, err
	}

	// 3️⃣ Resolver precedencia (priority DESC)
	sort.Slice(matchedRules, func(i, j int) bool {
		return matchedRules[i].Priority > matchedRules[j].Priority
	})

	status := intents.StatusNotTriggered
	confidence := 0.0
	reasonSteps := []intents.ReasoningStep{}

	// 4️⃣ Aplicar reglas
	for i, rule := range matchedRules {
		// status siempre viene de la regla de mayor prioridad
		if i == 0 {
			status = intents.IntentStatus(rule.Then.Status)
		}

		confidence += rule.Then.ConfidenceBoost

		reasonSteps = append(reasonSteps, intents.ReasoningStep{
			Step:        i + 1,
			Description: rule.Explanation,
		})
	}

	if confidence > 1 {
		confidence = 1
	}

	// 5️⃣ Fallback explícito si no matchea nada
	if len(matchedRules) == 0 {
		status = intents.StatusLowConfidence
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
			Explanation: "Declarative rules matched the current signal snapshot.",
		},
		Guardrails: &intents.Guardrails{
			HumanConfirmationRequired: true,
		},
	}, nil
}

func mapSignals(inputs []SignalInput) []intents.SignalUsage {
	if len(inputs) == 0 {
		return nil
	}

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
