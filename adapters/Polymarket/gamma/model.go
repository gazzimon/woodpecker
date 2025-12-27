package polymarket

import (
	"encoding/json"
	"fmt"
)

// Float64 is a helper that unmarshals from:
// - JSON number (e.g. 0.12)
// - JSON string (e.g. "0.12")
// - null
type Float64 float64

func (f *Float64) UnmarshalJSON(b []byte) error {
	// null
	if string(b) == "null" {
		*f = 0
		return nil
	}

	// number
	var num json.Number
	if err := json.Unmarshal(b, &num); err == nil {
		v, err := num.Float64()
		if err != nil {
			return fmt.Errorf("Float64: invalid number %q: %w", string(b), err)
		}
		*f = Float64(v)
		return nil
	}

	// string
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		if s == "" {
			*f = 0
			return nil
		}
		n := json.Number(s)
		v, err := n.Float64()
		if err != nil {
			return fmt.Errorf("Float64: invalid string-number %q: %w", s, err)
		}
		*f = Float64(v)
		return nil
	}

	return fmt.Errorf("Float64: unsupported json value %q", string(b))
}

// Raw models mirroring Gamma API responses.
// These MUST NOT contain business logic.

type Event struct {
	ID            string   `json:"id"`
	Ticker        *string  `json:"ticker,omitempty"`
	Slug          *string  `json:"slug,omitempty"`
	Title         *string  `json:"title,omitempty"`
	CreatedAt     *string  `json:"createdAt,omitempty"`
	EndDate       *string  `json:"endDate,omitempty"`
	Active        *bool    `json:"active,omitempty"`
	Closed        *bool    `json:"closed,omitempty"`
	Restricted    *bool    `json:"restricted,omitempty"`
	Liquidity     Float64  `json:"liquidity"`
	Volume        Float64  `json:"volume"`
	Volume24hr    Float64  `json:"volume24hr"`
	OpenInterest  Float64  `json:"openInterest"`
	LiquidityAmm  Float64  `json:"liquidityAmm"`
	LiquidityClob Float64  `json:"liquidityClob"`
	Markets       []Market `json:"markets"`
}

type Market struct {
	ID          string  `json:"id"`
	Slug        *string `json:"slug,omitempty"`
	ConditionID *string `json:"conditionId,omitempty"`

	// These two sometimes come as numbers; keep as Float64 for safety.
	BestBid Float64 `json:"bestBid"`
	BestAsk Float64 `json:"bestAsk"`

	LastTradePrice Float64 `json:"lastTradePrice"`

	// Gamma often exposes both numeric and string forms; we prefer numeric ones if present.
	VolumeNum    Float64 `json:"volumeNum"`
	LiquidityNum Float64 `json:"liquidityNum"`

	// Timestamps (strings or null)
	StartDate *string `json:"startDate,omitempty"`
	EndDate   *string `json:"endDate,omitempty"`
	UpdatedAt *string `json:"updatedAt,omitempty"`

	// JSON-encoded arrays in a string (Gamma does this on some fields)
	Outcomes *string `json:"outcomes,omitempty"`
}
