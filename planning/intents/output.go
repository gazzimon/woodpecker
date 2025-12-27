package intents

import "time"

// IntentOutput is the canonical response produced by the Planning Layer.
// It is intentionally generic: "evaluation" is a free-form map for intent-specific metrics.
type IntentOutput struct {
	Meta       Meta          `json:"meta"`
	Status     IntentStatus  `json:"status"`
	Confidence float64       `json:"confidence"` // 0..1
	Summary    string        `json:"summary"`

	Signals   []SignalUsage `json:"signals"`
	Reasoning Reasoning     `json:"reasoning"`

	// Evaluation is intent-specific metrics (JSON Schema: additionalProperties=true)
	Evaluation map[string]any `json:"evaluation,omitempty"`

	Guardrails *Guardrails `json:"guardrails,omitempty"`
}

type Meta struct {
	IntentID  string    `json:"intent_id"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

// IntentStatus is the single source of truth for intent evaluation states.
// It MUST cover all possible rule outputs.
type IntentStatus string

const (
	StatusNotTriggered  IntentStatus = "not_triggered"
	StatusLowConfidence IntentStatus = "low_confidence"
	StatusWeakSignal    IntentStatus = "weak_signal"
	StatusModerateSignal IntentStatus = "moderate_signal"
	StatusStrongSignal  IntentStatus = "strong_signal"
)

type SignalUsage struct {
	SignalID string  `json:"signal_id"`
	Value    float64 `json:"value"`
	Weight   float64 `json:"weight"` // 0..1
}

type Reasoning struct {
	Logic       []ReasoningStep `json:"logic"`
	Explanation string          `json:"explanation"`
}

type ReasoningStep struct {
	Step        int    `json:"step"`
	Description string `json:"description"`
}

type Guardrails struct {
	HumanConfirmationRequired bool `json:"human_confirmation_required,omitempty"`
	ConfidenceCapped          bool `json:"confidence_capped,omitempty"`
}
