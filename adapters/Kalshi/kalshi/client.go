package kalshi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"woodpecker-kalshi/model"
)

type Client struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

func New(apiKey string) *Client {
	return &Client{
		BaseURL: "https://api.elections.kalshi.com/trade-api/v2",
		APIKey:  apiKey,
		Client: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

// 🔑 USAR EVENT_TICKER para above/below
func (c *Client) GetMarketsByEvent(eventTicker string) ([]model.Market, error) {
	url := fmt.Sprintf(
		"%s/markets?event_ticker=%s&limit=1000",
		c.BaseURL,
		eventTicker,
	)

	fmt.Println("📌 GET", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kalshi HTTP %s", resp.Status)
	}

	var out model.MarketsResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	fmt.Printf("↳ Kalshi respondió %d markets\n", len(out.Markets))
	return out.Markets, nil
}
