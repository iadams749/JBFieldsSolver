// Package main provides a WebAssembly entrypoint for the Jumbleberry Fields solver.
// It registers JavaScript-callable functions for loading the EV table and solving game states.
package main

import (
	"encoding/json"
	"syscall/js"

	"github.com/iadams749/JBFieldsSolver/internal/ev"
	"github.com/iadams749/JBFieldsSolver/internal/solver"
)

var table *ev.Table

func main() {
	js.Global().Set("jbfSolve", js.FuncOf(solve))
	js.Global().Set("jbfLoadEVTable", js.FuncOf(loadEVTable))
	js.Global().Set("jbfComputeEVTable", js.FuncOf(computeEVTable))
	js.Global().Set("jbfReady", js.ValueOf(true))

	// Block forever so the Go runtime stays alive.
	select {}
}

// loadEVTable loads the EV table from a JSON string passed from JavaScript.
// Call: jbfLoadEVTable(jsonString) → returns "" on success, or error string.
func loadEVTable(_ js.Value, args []js.Value) any {
	if len(args) < 1 {
		return "missing JSON argument"
	}

	jsonStr := args[0].String()

	var entries []struct {
		CategorySet uint16  `json:"category_set"`
		EV          float64 `json:"ev"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &entries); err != nil {
		return "invalid JSON: " + err.Error()
	}

	t := ev.NewTable()
	for _, e := range entries {
		t.SetEV(e.CategorySet, e.EV)
	}
	table = t
	return ""
}

// computeEVTable computes the EV table from scratch (fallback).
// Call: jbfComputeEVTable() → returns "" on success.
func computeEVTable(_ js.Value, _ []js.Value) any {
	table = ev.Compute(nil)
	return ""
}

type solveRequest struct {
	Dice       string `json:"dice"`
	RollsLeft  int    `json:"rolls_left"`
	Categories string `json:"categories"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// solve takes a JSON request string and returns a JSON response string.
// Call: jbfSolve(jsonString) → JSON string.
func solve(_ js.Value, args []js.Value) any {
	if table == nil {
		return marshalError("EV table not loaded")
	}

	if len(args) < 1 {
		return marshalError("missing JSON argument")
	}

	var req solveRequest
	if err := json.Unmarshal([]byte(args[0].String()), &req); err != nil {
		return marshalError("invalid JSON: " + err.Error())
	}

	dice, err := solver.ParseDice(req.Dice)
	if err != nil {
		return marshalError("invalid dice: " + err.Error())
	}

	if req.RollsLeft < 0 || req.RollsLeft > 2 {
		return marshalError("rolls_left must be 0, 1, or 2")
	}

	cs, err := solver.ParseCategories(req.Categories)
	if err != nil {
		return marshalError("invalid categories: " + err.Error())
	}

	if cs == 0 {
		return marshalError("at least one category must be selected")
	}

	rec := solver.Solve(dice, req.RollsLeft, cs, table)
	result := solver.RecommendationToJSON(rec)

	data, _ := json.Marshal(result)
	return string(data)
}

func marshalError(msg string) string {
	data, _ := json.Marshal(errorResponse{Error: msg})
	return string(data)
}
