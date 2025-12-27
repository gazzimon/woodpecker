package polymarket

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"time"
)

// Snapshot is a frozen view of Gamma at time T.
type Snapshot struct {
	SnapshotID string
	Timestamp  time.Time
	Source     string

	Events []EventSnapshot
	Stats  SnapshotStats
}

type EventSnapshot struct {
	EventID string
	Slug    string
	Title   string
	EndDate time.Time

	Liquidity float64
	Volume    float64

	Markets []MarketPoint
}

// MarketPoint is the normalized per-market slice used by feature/signal logic.
type MarketPoint struct {
	MarketID    string
	Slug        string
	ConditionID string

	BestBid  float64
	BestAsk  float64
	MidPrice float64
	Spread   float64

	Liquidity float64
	Volume    float64

	LastTrade float64
	UpdatedAt time.Time
}

type SnapshotStats struct {
	TotalEvents    int
	TotalMarkets   int
	AvgLiquidity   float64
	AvgSpread      float64
	ExtremeMarkets int
}

// BuildSnapshot converts Gamma events into a normalized Snapshot.
func BuildSnapshot(events []Event) Snapshot {
	now := time.Now().UTC()

	var (
		eventSnapshots []EventSnapshot
		totalLiquidity float64
		totalSpread    float64
		totalMarkets   int
		extremeCount   int
	)

	for _, e := range events {
		es := EventSnapshot{
			EventID:   e.ID,
			Liquidity: float64(e.Liquidity),
			Volume:    float64(e.Volume),
		}

		if e.Slug != nil {
			es.Slug = *e.Slug
		}
		if e.Title != nil {
			es.Title = *e.Title
		}

		if e.EndDate != nil && *e.EndDate != "" {
			if t, err := time.Parse(time.RFC3339, *e.EndDate); err == nil {
				es.EndDate = t
			}
		}

		for _, m := range e.Markets {
			mp := MarketPoint{
				MarketID:  m.ID,
				BestBid:   float64(m.BestBid),
				BestAsk:   float64(m.BestAsk),
				Liquidity: float64(m.LiquidityNum),
				Volume:    float64(m.VolumeNum),
				LastTrade: float64(m.LastTradePrice),
			}

			if m.Slug != nil {
				mp.Slug = *m.Slug
			}
			if m.ConditionID != nil {
				mp.ConditionID = *m.ConditionID
			}
			if m.UpdatedAt != nil && *m.UpdatedAt != "" {
				if t, err := time.Parse(time.RFC3339, *m.UpdatedAt); err == nil {
					mp.UpdatedAt = t
				}
			}

			// Mid price + spread
			if mp.BestBid > 0 || mp.BestAsk > 0 {
				// If one side missing, mid collapses to known side
				if mp.BestBid > 0 && mp.BestAsk > 0 {
					mp.MidPrice = (mp.BestBid + mp.BestAsk) / 2
					mp.Spread = mp.BestAsk - mp.BestBid
					totalSpread += mp.Spread
				} else if mp.BestBid > 0 {
					mp.MidPrice = mp.BestBid
				} else {
					mp.MidPrice = mp.BestAsk
				}
			}

			// “Extreme” heuristic
			if (mp.BestBid > 0 && mp.BestBid < 0.05) || (mp.BestAsk > 0.95) {
				extremeCount++
			}

			totalLiquidity += mp.Liquidity
			totalMarkets++

			es.Markets = append(es.Markets, mp)
		}

		eventSnapshots = append(eventSnapshots, es)
	}

	stats := SnapshotStats{
		TotalEvents:    len(eventSnapshots),
		TotalMarkets:   totalMarkets,
		ExtremeMarkets: extremeCount,
	}
	if totalMarkets > 0 {
		stats.AvgLiquidity = totalLiquidity / float64(totalMarkets)
		// AvgSpread should only be computed over markets where spread exists
		if totalSpread > 0 {
			stats.AvgSpread = totalSpread / float64(totalMarkets)
		}
	}

	s := Snapshot{
		Timestamp: now,
		Source:    "polymarket-gamma",
		Events:    eventSnapshots,
		Stats:     stats,
	}
	s.SnapshotID = computeSnapshotID(s)

	return s
}

func computeSnapshotID(s Snapshot) string {
	// Deterministic: sort by event id then market id to keep stable across map ordering.
	type row struct {
		eid string
		mid string
		p   float64
	}
	var rows []row
	for _, e := range s.Events {
		for _, m := range e.Markets {
			rows = append(rows, row{eid: e.EventID, mid: m.MarketID, p: m.MidPrice})
		}
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].eid == rows[j].eid {
			return rows[i].mid < rows[j].mid
		}
		return rows[i].eid < rows[j].eid
	})

	h := sha256.New()
	h.Write([]byte(s.Source))
	h.Write([]byte(s.Timestamp.Format(time.RFC3339Nano)))

	for _, r := range rows {
		h.Write([]byte(r.eid))
		h.Write([]byte(r.mid))
		h.Write([]byte(fmt.Sprintf("%.6f", r.p)))
	}

	return hex.EncodeToString(h.Sum(nil))[:16]
}
