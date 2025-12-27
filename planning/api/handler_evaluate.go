package api

import (
	"encoding/json"
	"net/http"

	"woodpecker/planning/reasoner"
)

type PlanningHandler struct {
	Reasoner reasoner.IntentReasoner
}

func (h *PlanningHandler) EvaluateIntent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req IntentEvaluateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	signalInputs := make([]reasoner.SignalInput, 0, len(req.Signals))
	for _, s := range req.Signals {
		signalInputs = append(signalInputs, reasoner.SignalInput{
			SignalID: s.SignalID,
			Value:    s.Value,
		})
	}

	result, err := h.Reasoner.Evaluate(req.IntentID, req.Params, signalInputs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}
