package intents

import (
	"errors"
	"fmt"
)

func (o IntentOutput) ValidateBasic() error {
	if o.Meta.IntentID == "" {
		return errors.New("meta.intent_id is required")
	}
	if o.Meta.Version == "" {
		return errors.New("meta.version is required")
	}
	switch o.Status {
	case StatusNotTriggered, StatusWeakSignal, StatusStrongSignal:
		// ok
	default:
		return fmt.Errorf("status is invalid: %q", o.Status)
	}
	if o.Confidence < 0 || o.Confidence > 1 {
		return fmt.Errorf("confidence out of range: %v", o.Confidence)
	}
	if o.Summary == "" {
		return errors.New("summary is required")
	}
	if len(o.Signals) == 0 {
		return errors.New("signals must be non-empty")
	}
	for i, s := range o.Signals {
		if s.SignalID == "" {
			return fmt.Errorf("signals[%d].signal_id is required", i)
		}
		if s.Weight < 0 || s.Weight > 1 {
			return fmt.Errorf("signals[%d].weight out of range: %v", i, s.Weight)
		}
	}
	if o.Reasoning.Explanation == "" {
		return errors.New("reasoning.explanation is required")
	}
	// reasoning.logic can be empty in some modes, but schema requires it; keep strict:
	if len(o.Reasoning.Logic) == 0 {
		return errors.New("reasoning.logic must be non-empty")
	}
	return nil
}
