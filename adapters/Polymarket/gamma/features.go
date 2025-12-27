package polymarket

import "math"

// FeatureVector contains continuous, numerical features.
// NO thresholds. NO decisions.
type FeatureVector struct {
	PEvent              float64
	LogOdds             float64
	ProbabilityMomentum float64
	BeliefVolatility    float64
	ImpliedConfidence   float64
	Dispersion          float64
}

func ComputeFeatures(
	current MarketPoint,
	previous *MarketPoint,
	history []MarketPoint,
	peers []MarketPoint,
) FeatureVector {

	p := clampProb(current.MidPrice)
	logOdds := logit(p)

	momentum := 0.0
	if previous != nil {
		momentum = logOdds - logit(clampProb(previous.MidPrice))
	}

	vol := logOddsVolatility(history)

	conf := impliedConfidence(
		current.Liquidity,
		current.Volume,
		current.Spread,
	)

	disp := crossMarketDispersion(current, peers)

	return FeatureVector{
		PEvent:              p,
		LogOdds:             logOdds,
		ProbabilityMomentum: momentum,
		BeliefVolatility:    vol,
		ImpliedConfidence:   conf,
		Dispersion:          disp,
	}
}

/* ---------- math helpers ---------- */

func logit(p float64) float64 { return math.Log(p / (1 - p)) }

func clampProb(p float64) float64 {
	if p < 1e-6 {
		return 1e-6
	}
	if p > 1-1e-6 {
		return 1 - 1e-6
	}
	return p
}

func logOddsVolatility(history []MarketPoint) float64 {
	if len(history) < 2 {
		return 0
	}
	vals := make([]float64, 0, len(history))
	for _, h := range history {
		vals = append(vals, logit(clampProb(h.MidPrice)))
	}
	μ := mean(vals)
	var sum float64
	for _, v := range vals {
		d := v - μ
		sum += d * d
	}
	return math.Sqrt(sum / float64(len(vals)))
}

// A lightweight confidence proxy: more liquidity/volume, tighter spread => higher confidence.
func impliedConfidence(liquidity, volume, spread float64) float64 {
	eps := 1e-6
	// If spread missing, penalize a bit
	if spread <= 0 {
		spread = 0.02
	}
	raw := (math.Log1p(liquidity) + math.Log1p(volume)) / (spread + eps)
	return math.Tanh(raw / 10)
}

func crossMarketDispersion(current MarketPoint, peers []MarketPoint) float64 {
	if len(peers) == 0 {
		return 0
	}
	vals := []float64{logit(clampProb(current.MidPrice))}
	for _, p := range peers {
		vals = append(vals, logit(clampProb(p.MidPrice)))
	}
	μ := mean(vals)
	var sum float64
	for _, v := range vals {
		d := v - μ
		sum += d * d
	}
	return math.Sqrt(sum / float64(len(vals)))
}

func mean(xs []float64) float64 {
	var s float64
	for _, x := range xs {
		s += x
	}
	return s / float64(len(xs))
}
