package ev

import (
	"math"
	"testing"

	"github.com/iadams749/JBFieldsSolver/internal/game"
)

func TestTableEV(t *testing.T) {
	t.Parallel()

	// Create a simple table with known values
	table := &Table{}
	table.ev[0] = 0.0
	table.ev[1] = 10.5
	table.ev[511] = 121.8

	tests := []struct {
		name string
		cs   game.CategorySet
		want float64
	}{
		{
			name: "no categories",
			cs:   0,
			want: 0.0,
		},
		{
			name: "single category",
			cs:   1,
			want: 10.5,
		},
		{
			name: "all categories",
			cs:   511,
			want: 121.8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := table.EV(tt.cs)
			if math.Abs(got-tt.want) > 1e-6 {
				t.Errorf("Table.EV(%v) = %v, want %v", tt.cs, got, tt.want)
			}
		})
	}
}

func TestEnumerateKeeps(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		dice          game.Dice
		wantMinCount  int // minimum number of keep decisions
		checkSpecific []game.Dice
	}{
		{
			name:         "all same berry",
			dice:         game.Dice{5, 0, 0, 0, 0},
			wantMinCount: 6, // keep 0, 1, 2, 3, 4, or 5
			checkSpecific: []game.Dice{
				{0, 0, 0, 0, 0}, // keep nothing
				{5, 0, 0, 0, 0}, // keep all
				{3, 0, 0, 0, 0}, // keep 3
			},
		},
		{
			name:         "two types",
			dice:         game.Dice{3, 2, 0, 0, 0},
			wantMinCount: 12, // (3+1)*(2+1) = 12 combinations
			checkSpecific: []game.Dice{
				{0, 0, 0, 0, 0}, // keep nothing
				{3, 2, 0, 0, 0}, // keep all
				{3, 0, 0, 0, 0}, // keep only J
				{0, 2, 0, 0, 0}, // keep only S
				{1, 1, 0, 0, 0}, // keep one of each
			},
		},
		{
			name:         "mixed dice",
			dice:         game.Dice{2, 1, 1, 1, 0},
			wantMinCount: 24, // (2+1)*(1+1)*(1+1)*(1+1) = 24
			checkSpecific: []game.Dice{
				{0, 0, 0, 0, 0}, // keep nothing
				{2, 1, 1, 1, 0}, // keep all
				{0, 0, 0, 1, 0}, // keep only moonberry
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var keeps []game.Dice
			EnumerateKeeps(tt.dice, func(keep game.Dice, numKept int) {
				keeps = append(keeps, keep)

				// Verify numKept matches actual keep total
				if keep.Total() != numKept {
					t.Errorf("EnumerateKeeps callback: keep.Total() = %d, numKept = %d", keep.Total(), numKept)
				}

				// Verify keep doesn't exceed original dice
				for b := game.Berry(0); b < game.NumBerryTypes; b++ {
					if keep[b] > tt.dice[b] {
						t.Errorf("EnumerateKeeps: keep[%v]=%d > dice[%v]=%d", b, keep[b], b, tt.dice[b])
					}
				}
			})

			if len(keeps) < tt.wantMinCount {
				t.Errorf("EnumerateKeeps(%v) generated %d keeps, want at least %d", tt.dice, len(keeps), tt.wantMinCount)
			}

			// Check that specific keep decisions are present
			for _, expected := range tt.checkSpecific {
				found := false
				for _, k := range keeps {
					if k == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("EnumerateKeeps(%v) did not generate expected keep: %v", tt.dice, expected)
				}
			}

			// Verify all keeps are unique
			seen := make(map[game.Dice]bool)
			for _, k := range keeps {
				if seen[k] {
					t.Errorf("EnumerateKeeps(%v) generated duplicate keep: %v", tt.dice, k)
				}
				seen[k] = true
			}
		})
	}
}

func TestComputeRerollLayer(t *testing.T) {
	// This is more of an integration test
	t.Parallel()

	allDice := game.AllDice()
	numDice := len(allDice)

	// Create a simple previous layer: value = sum of dice points
	prevLayer := make([]float64, numDice)
	for i, d := range allDice {
		prevLayer[i] = float64(d.Points())
	}

	curLayer := make([]float64, numDice)
	ComputeRerollLayer(allDice, prevLayer, curLayer)

	// Check basic properties
	for i := range curLayer {
		// Current layer value should be at least as good as keeping all dice (no reroll)
		if curLayer[i] < prevLayer[i] {
			t.Errorf("ComputeRerollLayer: curLayer[%d]=%v < prevLayer[%d]=%v (reroll worse than keeping all)",
				i, curLayer[i], i, prevLayer[i])
		}

		// Value should be finite and non-negative
		if math.IsNaN(curLayer[i]) || math.IsInf(curLayer[i], 0) {
			t.Errorf("ComputeRerollLayer: curLayer[%d] is not finite: %v", i, curLayer[i])
		}
		if curLayer[i] < 0 {
			t.Errorf("ComputeRerollLayer: curLayer[%d]=%v is negative", i, curLayer[i])
		}
	}
}

func TestComputeBasic(t *testing.T) {
	// Test that Compute runs without panicking and produces reasonable values
	t.Parallel()

	table := Compute(nil) // no progress callback

	// Test some basic properties
	if table == nil {
		t.Fatal("Compute() returned nil")
	}

	// EV for no categories should be 0
	if ev := table.EV(0); ev != 0.0 {
		t.Errorf("table.EV(0) = %v, want 0.0", ev)
	}

	// EV for all categories should be positive and reasonable (around 121.8)
	evAll := table.EV(game.AllCategories)
	if evAll < 100.0 || evAll > 150.0 {
		t.Errorf("table.EV(AllCategories) = %v, expected to be in range [100, 150]", evAll)
	}

	// EV should increase (or stay same) as more categories are added
	ev1 := table.EV(1 << game.CatMoonberry)
	ev2 := table.EV((1 << game.CatMoonberry) | (1 << game.CatFreeRoll))
	if ev2 < ev1 {
		t.Errorf("Adding category decreased EV: ev1=%v, ev2=%v", ev1, ev2)
	}
}

func TestSaveAndLoadJSON(t *testing.T) {
	// Create a temporary file for testing
	t.Parallel()

	// Create a simple table
	table := &Table{}
	table.ev[0] = 0.0
	table.ev[1] = 10.5
	table.ev[511] = 121.8

	// Save to temp file
	tmpFile := t.TempDir() + "/test_ev_table.json"
	if err := table.SaveJSON(tmpFile); err != nil {
		t.Fatalf("SaveJSON() error = %v", err)
	}

	// Load it back
	loaded, err := LoadJSON(tmpFile)
	if err != nil {
		t.Fatalf("LoadJSON() error = %v", err)
	}

	// Verify values match
	tests := []game.CategorySet{0, 1, 511}
	for _, cs := range tests {
		want := table.EV(cs)
		got := loaded.EV(cs)
		if math.Abs(got-want) > 1e-6 {
			t.Errorf("After save/load: EV(%v) = %v, want %v", cs, got, want)
		}
	}
}
