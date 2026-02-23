package game

import (
	"testing"
)

func TestScore(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		dice     Dice
		category Category
		want     int
	}{
		// Jumbleberry tests
		{
			name:     "jumbleberry with 5",
			dice:     Dice{5, 0, 0, 0, 0},
			category: CatJumbleberry,
			want:     10,
		},
		{
			name:     "jumbleberry with 0",
			dice:     Dice{0, 5, 0, 0, 0},
			category: CatJumbleberry,
			want:     0,
		},
		{
			name:     "jumbleberry with 2",
			dice:     Dice{2, 1, 1, 1, 0},
			category: CatJumbleberry,
			want:     4,
		},

		// Sugarberry tests
		{
			name:     "sugarberry with 5",
			dice:     Dice{0, 5, 0, 0, 0},
			category: CatSugarberry,
			want:     10,
		},
		{
			name:     "sugarberry with 0",
			dice:     Dice{5, 0, 0, 0, 0},
			category: CatSugarberry,
			want:     0,
		},
		{
			name:     "sugarberry with 3",
			dice:     Dice{0, 3, 2, 0, 0},
			category: CatSugarberry,
			want:     6,
		},

		// Pickleberry tests
		{
			name:     "pickleberry with 5",
			dice:     Dice{0, 0, 5, 0, 0},
			category: CatPickleberry,
			want:     20,
		},
		{
			name:     "pickleberry with 0",
			dice:     Dice{5, 0, 0, 0, 0},
			category: CatPickleberry,
			want:     0,
		},
		{
			name:     "pickleberry with 1",
			dice:     Dice{2, 2, 1, 0, 0},
			category: CatPickleberry,
			want:     4,
		},

		// Moonberry tests
		{
			name:     "moonberry with 5",
			dice:     Dice{0, 0, 0, 5, 0},
			category: CatMoonberry,
			want:     35,
		},
		{
			name:     "moonberry with 0",
			dice:     Dice{5, 0, 0, 0, 0},
			category: CatMoonberry,
			want:     0,
		},
		{
			name:     "moonberry with 2",
			dice:     Dice{2, 1, 0, 2, 0},
			category: CatMoonberry,
			want:     14,
		},

		// Basket of Three tests
		{
			name:     "3k with exactly 3",
			dice:     Dice{3, 2, 0, 0, 0},
			category: CatBasketOfThree,
			want:     10, // 3*2 + 2*2
		},
		{
			name:     "3k with 4",
			dice:     Dice{4, 1, 0, 0, 0},
			category: CatBasketOfThree,
			want:     10, // 4*2 + 1*2
		},
		{
			name:     "3k with 5",
			dice:     Dice{5, 0, 0, 0, 0},
			category: CatBasketOfThree,
			want:     10, // 5*2
		},
		{
			name:     "3k with pests",
			dice:     Dice{0, 0, 0, 2, 3},
			category: CatBasketOfThree,
			want:     14, // 2*7 + 3*0
		},
		{
			name:     "3k fails with 2",
			dice:     Dice{2, 2, 1, 0, 0},
			category: CatBasketOfThree,
			want:     0,
		},
		{
			name:     "3k fails with all different",
			dice:     Dice{1, 1, 1, 1, 1},
			category: CatBasketOfThree,
			want:     0,
		},

		// Basket of Four tests
		{
			name:     "4k with exactly 4",
			dice:     Dice{4, 1, 0, 0, 0},
			category: CatBasketOfFour,
			want:     10,
		},
		{
			name:     "4k with 5",
			dice:     Dice{0, 0, 0, 0, 5},
			category: CatBasketOfFour,
			want:     0, // all pests
		},
		{
			name:     "4k fails with 3",
			dice:     Dice{3, 2, 0, 0, 0},
			category: CatBasketOfFour,
			want:     0,
		},
		{
			name:     "4k with 4 moonberries",
			dice:     Dice{0, 0, 1, 4, 0},
			category: CatBasketOfFour,
			want:     32, // 1*4 + 4*7
		},

		// Basket of Five tests
		{
			name:     "5k with all same",
			dice:     Dice{0, 0, 5, 0, 0},
			category: CatBasketOfFive,
			want:     20,
		},
		{
			name:     "5k with all pests",
			dice:     Dice{0, 0, 0, 0, 5},
			category: CatBasketOfFive,
			want:     0,
		},
		{
			name:     "5k fails with 4",
			dice:     Dice{4, 1, 0, 0, 0},
			category: CatBasketOfFive,
			want:     0,
		},

		// Mixed Basket tests
		{
			name:     "mixed with all four berries",
			dice:     Dice{1, 1, 1, 2, 0},
			category: CatMixedBasket,
			want:     22, // 2+2+4+14
		},
		{
			name:     "mixed with one of each plus pest",
			dice:     Dice{1, 1, 1, 1, 1},
			category: CatMixedBasket,
			want:     15, // 2+2+4+7
		},
		{
			name:     "mixed fails without jumbleberry",
			dice:     Dice{0, 2, 2, 1, 0},
			category: CatMixedBasket,
			want:     0,
		},
		{
			name:     "mixed fails without sugarberry",
			dice:     Dice{2, 0, 2, 1, 0},
			category: CatMixedBasket,
			want:     0,
		},
		{
			name:     "mixed fails without pickleberry",
			dice:     Dice{2, 2, 0, 1, 0},
			category: CatMixedBasket,
			want:     0,
		},
		{
			name:     "mixed fails without moonberry",
			dice:     Dice{2, 2, 1, 0, 0},
			category: CatMixedBasket,
			want:     0,
		},

		// Free Roll tests
		{
			name:     "free roll with all jumbleberries",
			dice:     Dice{5, 0, 0, 0, 0},
			category: CatFreeRoll,
			want:     10,
		},
		{
			name:     "free roll with mixed",
			dice:     Dice{2, 1, 1, 1, 0},
			category: CatFreeRoll,
			want:     17, // 4+2+4+7
		},
		{
			name:     "free roll with all pests",
			dice:     Dice{0, 0, 0, 0, 5},
			category: CatFreeRoll,
			want:     0,
		},
		{
			name:     "free roll with all moonberries",
			dice:     Dice{0, 0, 0, 5, 0},
			category: CatFreeRoll,
			want:     35,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := Score(tt.dice, tt.category); got != tt.want {
				t.Errorf("Score(%v, %v) = %v, want %v", tt.dice, tt.category, got, tt.want)
			}
		})
	}
}

func TestHasNOfAKind(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		dice Dice
		n    uint8
		want bool
	}{
		{
			name: "has 3 of a kind",
			dice: Dice{3, 2, 0, 0, 0},
			n:    3,
			want: true,
		},
		{
			name: "has 4 of a kind",
			dice: Dice{4, 1, 0, 0, 0},
			n:    4,
			want: true,
		},
		{
			name: "has 5 of a kind",
			dice: Dice{5, 0, 0, 0, 0},
			n:    5,
			want: true,
		},
		{
			name: "does not have 3 of a kind",
			dice: Dice{2, 2, 1, 0, 0},
			n:    3,
			want: false,
		},
		{
			name: "does not have 4 of a kind",
			dice: Dice{3, 2, 0, 0, 0},
			n:    4,
			want: false,
		},
		{
			name: "pests count for n-of-a-kind",
			dice: Dice{0, 0, 0, 0, 5},
			n:    5,
			want: true,
		},
		{
			name: "has exactly n",
			dice: Dice{0, 3, 2, 0, 0},
			n:    3,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := hasNOfAKind(tt.dice, tt.n); got != tt.want {
				t.Errorf("hasNOfAKind(%v, %v) = %v, want %v", tt.dice, tt.n, got, tt.want)
			}
		})
	}
}
