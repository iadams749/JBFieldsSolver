package game

import (
	"fmt"
	"math"
	"testing"
)

func TestFirstRollProbSumsToOne(t *testing.T) {
	t.Parallel()

	allDice := AllDice()
	sum := 0.0
	for i := range allDice {
		sum += FirstRollProb(i)
	}
	if math.Abs(sum-1.0) > 1e-10 {
		t.Errorf("FirstRollProb sum = %v, want 1.0", sum)
	}
}

func TestRerollProbSumsToOne(t *testing.T) {
	t.Parallel()

	for n := 0; n <= NumDice; n++ {
		n := n // capture for parallel
		t.Run(fmt.Sprintf("n=%d", n), func(t *testing.T) {
			t.Parallel()
			rerolls := Rerolls(n)
			sum := 0.0
			for _, ro := range rerolls {
				sum += ro.Prob
			}
			if math.Abs(sum-1.0) > 1e-10 {
				t.Errorf("Rerolls(%d) prob sum = %v, want 1.0", n, sum)
			}
		})
	}
}

func TestDiceIndexRoundTrip(t *testing.T) {
	t.Parallel()

	allDice := AllDice()
	for i, d := range allDice {
		i, d := i, d // capture for parallel
		t.Run(fmt.Sprintf("dice_%d", i), func(t *testing.T) {
			t.Parallel()
			idx := DiceIndex(d)
			if idx != i {
				t.Errorf("DiceIndex(%v) = %d, want %d", d, idx, i)
			}
		})
	}
}

func TestDiceProb(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		dice Dice
		n    int
		want float64
	}{
		{
			name: "all same berry with 5 dice",
			dice: Dice{5, 0, 0, 0, 0},
			n:    5,
			want: math.Pow(0.3, 5), // (3/10)^5
		},
		{
			name: "no dice rolled",
			dice: Dice{0, 0, 0, 0, 0},
			n:    0,
			want: 1.0,
		},
		{
			name: "single die - jumbleberry",
			dice: Dice{1, 0, 0, 0, 0},
			n:    1,
			want: 0.3,
		},
		{
			name: "single die - moonberry",
			dice: Dice{0, 0, 0, 1, 0},
			n:    1,
			want: 0.1,
		},
		{
			name: "two different berries",
			dice: Dice{1, 1, 0, 0, 0},
			n:    2,
			want: 2 * 0.3 * 0.3, // 2! / (1!*1!) * 0.3 * 0.3
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := DiceProb(tt.dice, tt.n)
			if math.Abs(got-tt.want) > 1e-10 {
				t.Errorf("DiceProb(%v, %d) = %v, want %v", tt.dice, tt.n, got, tt.want)
			}
		})
	}
}

func TestAddDice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		a    Dice
		b    Dice
		want Dice
	}{
		{
			name: "both empty",
			a:    Dice{0, 0, 0, 0, 0},
			b:    Dice{0, 0, 0, 0, 0},
			want: Dice{0, 0, 0, 0, 0},
		},
		{
			name: "add to empty",
			a:    Dice{2, 1, 1, 1, 0},
			b:    Dice{0, 0, 0, 0, 0},
			want: Dice{2, 1, 1, 1, 0},
		},
		{
			name: "add two different",
			a:    Dice{2, 0, 0, 1, 0},
			b:    Dice{0, 1, 1, 0, 0},
			want: Dice{2, 1, 1, 1, 0},
		},
		{
			name: "add same types",
			a:    Dice{2, 1, 1, 1, 0},
			b:    Dice{1, 1, 1, 1, 0},
			want: Dice{3, 2, 2, 2, 0},
		},
		{
			name: "with pests",
			a:    Dice{1, 1, 1, 1, 1},
			b:    Dice{1, 1, 1, 1, 1},
			want: Dice{2, 2, 2, 2, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := AddDice(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("AddDice(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestAllDiceCount(t *testing.T) {
	t.Parallel()

	allDice := AllDice()
	expectedCount := 126 // C(5+5-1, 5-1) = C(9, 4)

	if got := len(allDice); got != expectedCount {
		t.Errorf("len(AllDice()) = %d, want %d", got, expectedCount)
	}
	if got := NumAllDice(); got != expectedCount {
		t.Errorf("NumAllDice() = %d, want %d", got, expectedCount)
	}
}

func TestRerollsCaching(t *testing.T) {
	t.Parallel()

	// Test that repeated calls return the same cached slice
	for n := 0; n <= NumDice; n++ {
		n := n
		t.Run(fmt.Sprintf("n=%d", n), func(t *testing.T) {
			t.Parallel()
			first := Rerolls(n)
			second := Rerolls(n)

			if len(first) != len(second) {
				t.Errorf("Rerolls(%d) returned different lengths: %d vs %d", n, len(first), len(second))
			}

			// Verify they're the same cached data by comparing addresses
			if len(first) > 0 && len(second) > 0 {
				if &first[0] != &second[0] {
					t.Errorf("Rerolls(%d) did not return cached slice", n)
				}
			}
		})
	}
}

func TestFaceProb(t *testing.T) {
	t.Parallel()

	tests := []struct {
		berry Berry
		want  float64
	}{
		{Jumbleberry, 0.3},
		{Sugarberry, 0.3},
		{Pickleberry, 0.2},
		{Moonberry, 0.1},
		{Pest, 0.1},
	}

	// Test individual probabilities
	for _, tt := range tests {
		t.Run(tt.berry.String(), func(t *testing.T) {
			t.Parallel()
			if got := FaceProb[tt.berry]; math.Abs(got-tt.want) > 1e-10 {
				t.Errorf("FaceProb[%v] = %v, want %v", tt.berry, got, tt.want)
			}
		})
	}

	// Test that all face probabilities sum to 1.0
	t.Run("sum to one", func(t *testing.T) {
		t.Parallel()
		sum := 0.0
		for b := Berry(0); b < NumBerryTypes; b++ {
			sum += FaceProb[b]
		}
		if math.Abs(sum-1.0) > 1e-10 {
			t.Errorf("FaceProb sum = %v, want 1.0", sum)
		}
	})
}
