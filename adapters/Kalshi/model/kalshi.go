package model

type MarketsResponse struct {
	Markets []Market `json:"markets"`
}

type Market struct {
	Ticker         string  `json:"ticker"`
	EventTicker    string  `json:"event_ticker"`
	SeriesTicker   string  `json:"series_ticker,omitempty"`
	Status         string  `json:"status"`
	MarketType     string  `json:"market_type"`

	StrikeType     string   `json:"strike_type"`
	FloorStrike    *float64 `json:"floor_strike,omitempty"`
	CapStrike      *float64 `json:"cap_strike,omitempty"`

	YesBid         int      `json:"yes_bid"`
	YesAsk         int      `json:"yes_ask"`
	NoBid          int      `json:"no_bid"`
	NoAsk          int      `json:"no_ask"`

	Liquidity      int64    `json:"liquidity"`
	Volume24h      int64    `json:"volume_24h"`
	OpenInterest   int64    `json:"open_interest"`

	OpenTime       string   `json:"open_time"`
	CloseTime      string   `json:"close_time"`
	ExpirationTime string   `json:"expiration_time"`
}
