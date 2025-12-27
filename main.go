package main

import (
	"log"
	"net/http"

	"woodpecker/planning/api"
	"woodpecker/planning/reasoner"
)

func main() {
	// 1️⃣ Load rules
	rules, err := reasoner.LoadRulesFromFile("planning/reasoner/rules.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// 2️⃣ Initialize Rule-Based Reasoner
	r := &reasoner.RuleBasedReasoner{
		Version: "v1",
		Rules:   rules,
	}

	// 3️⃣ Wire handler
	handler := &api.PlanningHandler{
		Reasoner: r,
	}

	// 4️⃣ Routes
	http.HandleFunc("/planning/intent/evaluate", handler.EvaluateIntent)

	// 5️⃣ Start server
	log.Println("🪵🐦 Woodpecker Planning Layer listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
