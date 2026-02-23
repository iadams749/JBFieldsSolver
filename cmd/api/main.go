// Package main provides an HTTP API server for the Jumbleberry Fields solver.
// POST /solve accepts JSON requests with dice, rolls_left, and categories,
// returning optimal action recommendations.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/iadams749/JBFieldsSolver/internal/ev"
	"github.com/iadams749/JBFieldsSolver/internal/evloader"
	"github.com/iadams749/JBFieldsSolver/internal/solver"
)

var table *ev.Table

type solveRequest struct {
	Dice       string `json:"dice"`
	RollsLeft  int    `json:"rolls_left"`
	Categories string `json:"categories"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	evPath := flag.String("ev", "ev_table.json", "path to EV table JSON")
	flag.Parse()

	var err error
	table, err = evloader.Load(*evPath)
	if err != nil {
		fmt.Printf("Fatal: %v\n", err)
		os.Exit(1)
	}

	http.HandleFunc("POST /solve", handleSolve)

	log.Printf("Listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func handleSolve(w http.ResponseWriter, r *http.Request) {
	var req solveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	dice, err := solver.ParseDice(req.Dice)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid dice: "+err.Error())
		return
	}

	if req.RollsLeft < 0 || req.RollsLeft > 2 {
		writeError(w, http.StatusBadRequest, "rolls_left must be 0, 1, or 2")
		return
	}

	cs, err := solver.ParseCategories(req.Categories)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid categories: "+err.Error())
		return
	}

	rec := solver.Solve(dice, req.RollsLeft, cs, table)
	result := solver.RecommendationToJSON(rec)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{Error: msg})
}
