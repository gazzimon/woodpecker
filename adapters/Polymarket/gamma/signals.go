package polymarket

import (
	"math"

	"woodpecker/planning/reasoner"
)

// BuildSignals maps continuous features into logical signals.
// Thresholds live here for now; later you can move them to YAML/config.
func BuildSignals(f FeatureVector) []reasoner.SignalInput {
	var out []reasoner.SignalInput

	// 1) PROBABILITY_ACCELERATION
	// Momentum is in log-odds space. We squash it into [0..1].
	// Positive large momentum => close to 1.
	accel := squashSigned(f.ProbabilityMomentum, 0.35) // scale factor
	if accel > 0.60 {
		out = append(out, reasoner.SignalInput{
			SignalID: "PROBABILITY_ACCELERATION",
			Value:    accel,
		})
	}

	// 2) CONVICTION_SPIKE
	// f.ImpliedConfidence is already ~[0..1]
	if f.ImpliedConfidence > 0.60 {
		out = append(out, reasoner.SignalInput{
			SignalID: "CONVICTION_SPIKE",
			Value:    clamp01(f.ImpliedConfidence),
		})
	}

	// 3) DIVERGENCE_ALERT
	// Dispersion is stdev of log-odds across peers. Convert to [0..1].
	div := squashPositive(f.Dispersion, 0.8)
	if div > 0.55 {
		out = append(out, reasoner.SignalInput{
			SignalID: "DIVERGENCE_ALERT",
			Value:    div,
		})
	}

	// 4) LOW_CONFIDENCE_MOVE
	// “Move” without confidence: high acceleration while confidence is low.
	// This is your LOW_CONFIDENCE_MOVE definition.
	lowConfMove := accel * (1.0 - clamp01(f.ImpliedConfidence))
	if lowConfMove > 0.55 {
		out = append(out, reasoner.SignalInput{
			SignalID: "LOW_CONFIDENCE_MOVE",
			Value:    clamp01(lowConfMove),
		})
	}

	// 5) REGIME_SHIFT (composite)
	// A regime shift is: strong acceleration + decent confidence + low volatility (stable belief)
	// Here “low volatility” means BeliefVolatility is small.
	volPenalty := 1.0 - squashPositive(f.BeliefVolatility, 1.2)
	regime := clamp01(0.45*accel + 0.35*clamp01(f.ImpliedConfidence) + 0.20*volPenalty)
	if regime > 0.60 {
		out = append(out, reasoner.SignalInput{
			SignalID: "REGIME_SHIFT",
			Value:    regime,
		})
	}

	return out
}

/* ---- helpers ---- */

func clamp01(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}

// squashPositive maps x>=0 roughly into [0..1] using tanh(x/scale).
func squashPositive(x, scale float64) float64 {
	if x <= 0 {
		return 0
	}
	if scale <= 0 {
		scale = 1
	}
	return clamp01(math.Tanh(x / scale))
}

// squashSigned maps signed values into [0..1] centered at 0.
// >0 => above 0.5; <0 => below 0.5.
func squashSigned(x, scale float64) float64 {
	if scale <= 0 {
		scale = 1
	}
	// tanh in [-1..1] -> map to [0..1]
	return clamp01(0.5 * (1.0 + math.Tanh(x/scale)))
}
