package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	polymarket "woodpecker/adapters/Polymarket/gamma"
)

func main() {
	limit := flag.Int("limit", 10, "cantidad de eventos a pedir a Gamma")
	offset := flag.Int("offset", 0, "offset para paginación")
	maxMarkets := flag.Int("maxMarkets", 50, "máximo de markets a imprimir/procesar (total)")
	perEvent := flag.Int("perEvent", 10, "máximo de markets por evento a imprimir/procesar")
	verbose := flag.Bool("v", false, "modo verbose")
	flag.Parse()

	client := polymarket.NewClient()

	start := time.Now()
	events, err := client.FetchActiveEvents(*limit, *offset)
	if err != nil {
		log.Fatalf("FetchActiveEvents failed: %v", err)
	}
	if len(events) == 0 {
		log.Printf("Gamma devolvió 0 eventos (limit=%d offset=%d). Esto puede ser normal si cambió el filtro o si Gamma respondió vacío.",
			*limit, *offset,
		)
		return
	}

	if *verbose {
		fmt.Printf("Fetched %d events from Gamma in %s\n", len(events), time.Since(start))
		// mini-resumen de eventos
		for i, e := range events {
			if i >= 5 {
				fmt.Printf("... (%d más)\n", len(events)-5)
				break
			}
			title := "(no title)"
			if e.Title != nil && *e.Title != "" {
				title = *e.Title
			}
			slug := "(no slug)"
			if e.Slug != nil && *e.Slug != "" {
				slug = *e.Slug
			}
			fmt.Printf("Event[%d] id=%s slug=%s title=%s markets=%d liquidity=%.2f volume=%.2f closed=%v\n",
				i, e.ID, slug, title, len(e.Markets), float64(e.Liquidity), float64(e.Volume), boolPtr(e.Closed),
			)
		}
	}

	snapshot := polymarket.BuildSnapshot(events)

	fmt.Printf(
		"SnapshotID=%s source=%s ts=%s events=%d markets=%d avg_liq=%.2f avg_spread=%.6f extreme=%d\n",
		snapshot.SnapshotID,
		snapshot.Source,
		snapshot.Timestamp.Format(time.RFC3339),
		snapshot.Stats.TotalEvents,
		snapshot.Stats.TotalMarkets,
		snapshot.Stats.AvgLiquidity,
		snapshot.Stats.AvgSpread,
		snapshot.Stats.ExtremeMarkets,
	)

	// Recolectamos todos los MarketPoint para poder calcular peers (opcional)
	all := make([]polymarket.MarketPoint, 0, snapshot.Stats.TotalMarkets)
	for _, es := range snapshot.Events {
		all = append(all, es.Markets...)
	}

	// Orden determinístico para debug
	sort.Slice(all, func(i, j int) bool {
		if all[i].MarketID == all[j].MarketID {
			return all[i].UpdatedAt.Before(all[j].UpdatedAt)
		}
		return all[i].MarketID < all[j].MarketID
	})

	processed := 0
	perEventProcessed := map[string]int{}

	// Un map para “previous” simplificado (en este probe no tenemos histórico real; guardamos último visto por MarketID).
	prevByMarket := map[string]polymarket.MarketPoint{}

	for _, es := range snapshot.Events {
		for _, mp := range es.Markets {
			if processed >= *maxMarkets {
				fmt.Printf("Reached maxMarkets=%d, stopping.\n", *maxMarkets)
				return
			}
			if perEventProcessed[es.EventID] >= *perEvent {
				continue
			}

			// previous (si ya vimos ese market en este mismo run)
			var prev *polymarket.MarketPoint
			if p, ok := prevByMarket[mp.MarketID]; ok {
				pp := p // copiar para tomar address estable
				prev = &pp
			}

			// peers (muy básico): 5 peers “cercanos” por liquidez, excluyendo el market actual
			peers := pickPeers(all, mp.MarketID, 5)

			features := polymarket.ComputeFeatures(
				mp,
				prev,
				nil,   // history (en este probe no la tenemos; para volatilidad real necesitás series temporales)
				peers, // peers cross-market
			)

			signals := polymarket.BuildSignals(features)

			// Diagnóstico de datos crudos que suelen romper features:
			// MidPrice=0 puede pasar si bid/ask vienen 0 o vacíos.
			dataWarn := ""
			if mp.MidPrice <= 0 {
				dataWarn = " [WARN: midPrice=0 (bid/ask missing?)]"
			}

			fmt.Printf(
				"event=%s market=%s p=%.4f logOdds=%.4f mom=%.4f vol=%.4f conf=%.4f disp=%.4f bid=%.4f ask=%.4f spread=%.4f liq=%.2f vol=%.2f%s\n",
				es.EventID,
				mp.MarketID,
				features.PEvent,
				features.LogOdds,
				features.ProbabilityMomentum,
				features.BeliefVolatility,
				features.ImpliedConfidence,
				features.Dispersion,
				mp.BestBid,
				mp.BestAsk,
				mp.Spread,
				mp.Liquidity,
				mp.Volume,
				dataWarn,
			)

			if *verbose {
				if len(signals) == 0 {
					fmt.Printf("  signals: (none)\n")
				} else {
					fmt.Printf("  signals:\n")
					for _, s := range signals {
						fmt.Printf("    - %s = %.4f\n", s.SignalID, s.Value)
					}
				}
			}

			// guardar previous para el próximo market igual (en este run)
			prevByMarket[mp.MarketID] = mp

			processed++
			perEventProcessed[es.EventID]++
		}
	}

	if processed == 0 {
		fmt.Fprintln(os.Stderr, "No se procesó ningún market. Revisa si los eventos vienen sin markets o si perEvent/maxMarkets es 0.")
	}
}

func boolPtr(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// pickPeers elige hasta n peers con mayor liquidez, excluyendo marketID.
// Esto es SOLO para debug de crossMarketDispersion; en producción vas a definir peers por “mismo evento/tema”.
func pickPeers(all []polymarket.MarketPoint, excludeMarketID string, n int) []polymarket.MarketPoint {
	type row struct {
		mp  polymarket.MarketPoint
		liq float64
	}
	rows := make([]row, 0, len(all))
	for _, m := range all {
		if m.MarketID == excludeMarketID {
			continue
		}
		rows = append(rows, row{mp: m, liq: m.Liquidity})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].liq > rows[j].liq })
	if n > len(rows) {
		n = len(rows)
	}
	out := make([]polymarket.MarketPoint, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, rows[i].mp)
	}
	return out
}
