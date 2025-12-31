package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"woodpecker-kalshi/kalshi"
)

const POLL_INTERVAL = 30 * time.Second

func main() {
	apiKey := os.Getenv("KALSHI_API_KEY")
	if apiKey == "" {
		panic("set KALSHI_API_KEY env var")
	}

	// 🔥 EVENT ticker real del mercado que pasaste
	eventTicker := "KXBTCD-25DEC3117"

	client := kalshi.New(apiKey)

	for {
		fmt.Println("🔍 Consultando mercados…")

		markets, err := client.GetMarketsByEvent(eventTicker)
		if err != nil {
			fmt.Println("❌ error:", err)
			time.Sleep(POLL_INTERVAL)
			continue
		}

		fmt.Printf("📊 markets recibidos: %d\n", len(markets))
		saveSnapshot(eventTicker, markets)

		time.Sleep(POLL_INTERVAL)
	}
}

func saveSnapshot(event string, markets any) {
	ts := time.Now().UTC().Format("2006-01-02T15-04-05Z")
	filename := fmt.Sprintf(
		"snapshots/snapshot_%s_%s.json",
		event,
		ts,
	)

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	enc.Encode(markets)

	fmt.Println("📸 Snapshot guardado:", filename)
}
