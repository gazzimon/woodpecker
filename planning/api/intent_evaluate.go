package api

// IntentEvaluateRequest is the input to the Planning Layer.
// It is intentionally generic and future-proof.
type IntentEvaluateRequest struct {
	IntentID string                 `json:"intent_id"`
	Params   map[string]any         `json:"params,omitempty"`
	Signals  []SignalSnapshot       `json:"signals"`
}

// SignalSnapshot represents a point-in-time signal value
// coming from the Signal Layer (already computed).
type SignalSnapshot struct {
	SignalID string  `json:"signal_id"`
	Value    float64 `json:"value"` // expected normalized 0..1
}
