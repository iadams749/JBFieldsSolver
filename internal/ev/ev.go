package ev

import (
	"encoding/json"
	"math"
	"os"

	"github.com/iadams749/JBFieldsSolver/internal/game"
)

// Table holds precomputed expected values for all category subsets.
// ev[cs] = expected total score added across all remaining rounds
// when categories cs remain, under optimal play, before any rolls.
type Table struct {
	ev [512]float64
}

// EV returns the expected value for the given category set.
func (t *Table) EV(cs game.CategorySet) float64 {
	return t.ev[cs]
}

// SetEV sets the expected value for a given category set index.
// Used by the WASM build to load the table from JSON passed via JS.
func (t *Table) SetEV(cs uint16, val float64) {
	t.ev[cs] = val
}

// NewTable returns an empty EV table.
func NewTable() *Table {
	return &Table{}
}

// Compute builds the full EV table using bottom-up dynamic programming.
// It processes subsets of increasing size: all 1-category subsets first,
// then 2-category subsets (using the 1-category results), up to all 9.
//
// onProgress is called after each subset size is completed, with the
// size just finished and total (9). Pass nil to suppress progress.
func Compute(onProgress func(size, total int)) *Table {
	t := &Table{}
	// ev[0] = 0 already (no categories left = no more score)

	allDice := game.AllDice()
	numDice := len(allDice)

	// Precompute score table: scoreTab[cat][diceIdx] = score
	var scoreTab [game.NumCategories][]float64
	for cat := game.Category(0); cat < game.NumCategories; cat++ {
		scoreTab[cat] = make([]float64, numDice)
		for i, d := range allDice {
			scoreTab[cat][i] = float64(game.Score(d, cat))
		}
	}

	// Process category sets bottom-up by size
	for size := 1; size <= int(game.NumCategories); size++ {
		for cs := game.CategorySet(1); cs <= game.AllCategories; cs++ {
			if cs.Count() != size {
				continue
			}

			// Working arrays: v[rollsLeft][diceIdx]
			// We only need two layers at a time (current and previous).
			v0 := make([]float64, numDice) // rollsLeft = 0
			v1 := make([]float64, numDice) // rollsLeft = 1
			v2 := make([]float64, numDice) // rollsLeft = 2

			// Layer 0: rollsLeft = 0, must pick a category to score
			for i := range allDice {
				bestVal := math.Inf(-1)
				cs.ForEach(func(cat game.Category) {
					val := scoreTab[cat][i] + t.ev[cs.Remove(cat)]
					if val > bestVal {
						bestVal = val
					}
				})
				v0[i] = bestVal
			}

			// Layer 1: rollsLeft = 1
			ComputeRerollLayer(allDice, v0, v1)

			// Layer 2: rollsLeft = 2
			ComputeRerollLayer(allDice, v1, v2)

			// EV[cs] = expected value over the first roll (all 5 dice)
			ev := 0.0
			for i := range allDice {
				ev += game.FirstRollProb(i) * v2[i]
			}
			t.ev[cs] = ev
		}

		if onProgress != nil {
			onProgress(size, int(game.NumCategories))
		}
	}

	return t
}

// ComputeRerollLayer computes the optimal value for each dice outcome
// when the player has one more reroll available.
//
// For each dice outcome d:
//
//	V(d, r) = max over all keep decisions k of:
//	  sum over reroll outcomes o of P(o) * prevLayer[index(k + o)]
func ComputeRerollLayer(allDice []game.Dice, prevLayer, curLayer []float64) {
	for i, d := range allDice {
		// Baseline: keep all dice (no reroll)
		bestEV := prevLayer[i]

		EnumerateKeeps(d, func(keep game.Dice, numKept int) {
			if numKept == game.NumDice {
				return // already handled as baseline
			}

			numRerolled := game.NumDice - numKept
			rerolls := game.Rerolls(numRerolled)

			ev := 0.0
			for _, ro := range rerolls {
				resultDice := game.AddDice(keep, ro.Dice)
				resultIdx := game.DiceIndex(resultDice)
				ev += ro.Prob * prevLayer[resultIdx]
			}

			if ev > bestEV {
				bestEV = ev
			}
		})

		curLayer[i] = bestEV
	}
}

// EnumerateKeeps calls fn for every valid keep decision given the current dice.
// A keep decision is how many of each type to hold (0 to d[type]).
func EnumerateKeeps(d game.Dice, fn func(keep game.Dice, numKept int)) {
	var rec func(berry game.Berry, current game.Dice, totalKept int)
	rec = func(berry game.Berry, current game.Dice, totalKept int) {
		if berry == game.NumBerryTypes {
			fn(current, totalKept)
			return
		}
		for k := uint8(0); k <= d[berry]; k++ {
			current[berry] = k
			rec(berry+1, current, totalKept+int(k))
		}
	}
	rec(0, game.Dice{}, 0)
}

// jsonEntry is the structure for each entry in the JSON output.
type jsonEntry struct {
	CategorySet uint16   `json:"category_set"`
	Categories  []string `json:"categories"`
	EV          float64  `json:"ev"`
}

// SaveJSON writes the EV table to a JSON file.
func (t *Table) SaveJSON(path string) error {
	var entries []jsonEntry
	for cs := game.CategorySet(1); cs <= game.AllCategories; cs++ {
		if cs.Count() == 0 {
			continue
		}
		var names []string
		cs.ForEach(func(cat game.Category) {
			names = append(names, cat.String())
		})
		entries = append(entries, jsonEntry{
			CategorySet: uint16(cs),
			Categories:  names,
			EV:          t.ev[cs],
		})
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LoadJSON reads an EV table from a previously saved JSON file.
func LoadJSON(path string) (*Table, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var entries []jsonEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	t := &Table{}
	for _, e := range entries {
		t.ev[e.CategorySet] = e.EV
	}
	return t, nil
}
