package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"woodpecker/planning/reasoner"
)

func TestEvaluateIntentHandler_OK(t *testing.T) {
	r := &reasoner.SimpleReasoner{Version: "v1"}
	h := &PlanningHandler{Reasoner: r}

	body := IntentEvaluateRequest{
		IntentID: "interpret.regime_state",
		Signals: []SignalSnapshot{
			{SignalID: "REGIME_SHIFT", Value: 0.82},
			{SignalID: "PROBABILITY_ACCELERATION", Value: 0.74},
			{SignalID: "CONVICTION_SPIKE", Value: 0.65},
		},
	}

	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/planning/intent/evaluate", bytes.NewReader(b))
	w := httptest.NewRecorder()

	h.EvaluateIntent(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}
