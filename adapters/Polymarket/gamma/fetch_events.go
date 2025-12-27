package polymarket

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) FetchActiveEvents(limit, offset int) ([]Event, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	url := fmt.Sprintf(
		"%s/events?order=id&ascending=false&closed=false&limit=%d&offset=%d",
		c.BaseURL,
		limit,
		offset,
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "woodpecker/0.1 (gamma-adapter)")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gamma api error: status %d", resp.StatusCode)
	}

	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, err
	}

	return events, nil
}
