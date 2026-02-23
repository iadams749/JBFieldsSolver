package game

import (
	"testing"
)

func TestBerryString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		berry Berry
		want  string
	}{
		{"Jumbleberry", Jumbleberry, "Jumbleberry"},
		{"Sugarberry", Sugarberry, "Sugarberry"},
		{"Pickleberry", Pickleberry, "Pickleberry"},
		{"Moonberry", Moonberry, "Moonberry"},
		{"Pest", Pest, "Pest"},
		{"Invalid", Berry(99), "Berry(99)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.berry.String(); got != tt.want {
				t.Errorf("Berry.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiceTotal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		dice Dice
		want int
	}{
		{
			name: "empty dice",
			dice: Dice{},
			want: 0,
		},
		{
			name: "all jumbleberry",
			dice: Dice{5, 0, 0, 0, 0},
			want: 5,
		},
		{
			name: "mixed dice",
			dice: Dice{2, 1, 1, 1, 0},
			want: 5,
		},
		{
			name: "with pest",
			dice: Dice{1, 1, 1, 1, 1},
			want: 5,
		},
		{
			name: "three dice",
			dice: Dice{0, 2, 0, 1, 0},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.dice.Total(); got != tt.want {
				t.Errorf("Dice.Total() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDicePoints(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		dice Dice
		want int
	}{
		{
			name: "no dice",
			dice: Dice{},
			want: 0,
		},
		{
			name: "all jumbleberries",
			dice: Dice{5, 0, 0, 0, 0},
			want: 10, // 5 * 2
		},
		{
			name: "all sugarberries",
			dice: Dice{0, 5, 0, 0, 0},
			want: 10, // 5 * 2
		},
		{
			name: "all pickleberries",
			dice: Dice{0, 0, 5, 0, 0},
			want: 20, // 5 * 4
		},
		{
			name: "all moonberries",
			dice: Dice{0, 0, 0, 5, 0},
			want: 35, // 5 * 7
		},
		{
			name: "all pests",
			dice: Dice{0, 0, 0, 0, 5},
			want: 0, // 5 * 0
		},
		{
			name: "mixed standard",
			dice: Dice{2, 1, 1, 1, 0}, // 2J + 1S + 1P + 1M = 4 + 2 + 4 + 7
			want: 17,
		},
		{
			name: "with pest",
			dice: Dice{1, 1, 1, 1, 1}, // 1J + 1S + 1P + 1M = 2 + 2 + 4 + 7
			want: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.dice.Points(); got != tt.want {
				t.Errorf("Dice.Points() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiceString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		dice Dice
		want string
	}{
		{
			name: "empty",
			dice: Dice{},
			want: "{J:0 S:0 P:0 M:0 X:0}",
		},
		{
			name: "mixed",
			dice: Dice{2, 1, 1, 1, 0},
			want: "{J:2 S:1 P:1 M:1 X:0}",
		},
		{
			name: "with pest",
			dice: Dice{0, 0, 0, 0, 5},
			want: "{J:0 S:0 P:0 M:0 X:5}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.dice.String(); got != tt.want {
				t.Errorf("Dice.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnumerateAllDice(t *testing.T) {
	t.Parallel()

	allDice := EnumerateAllDice()

	// Should have exactly 126 outcomes (C(5+5-1, 5-1) = C(9,4))
	const expectedCount = 126
	if len(allDice) != expectedCount {
		t.Errorf("EnumerateAllDice() returned %d outcomes, want %d", len(allDice), expectedCount)
	}

	// Each outcome should have exactly 5 dice
	for i, dice := range allDice {
		if total := dice.Total(); total != NumDice {
			t.Errorf("EnumerateAllDice()[%d] has %d dice, want %d", i, total, NumDice)
		}
	}

	// All outcomes should be unique
	seen := make(map[Dice]bool)
	for i, dice := range allDice {
		if seen[dice] {
			t.Errorf("EnumerateAllDice()[%d] is duplicate: %v", i, dice)
		}
		seen[dice] = true
	}
}

func TestBerryPoints(t *testing.T) {
	t.Parallel()

	tests := []struct {
		berry Berry
		want  int
	}{
		{Jumbleberry, 2},
		{Sugarberry, 2},
		{Pickleberry, 4},
		{Moonberry, 7},
		{Pest, 0},
	}

	for _, tt := range tests {
		t.Run(tt.berry.String(), func(t *testing.T) {
			t.Parallel()
			if got := BerryPoints[tt.berry]; got != tt.want {
				t.Errorf("BerryPoints[%v] = %v, want %v", tt.berry, got, tt.want)
			}
		})
	}
}
