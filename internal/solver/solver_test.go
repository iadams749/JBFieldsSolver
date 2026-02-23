package solver

import (
	"testing"

	"github.com/iadams749/JBFieldsSolver/internal/game"
)

func TestFormatKeep(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		keep game.Dice
		want string
	}{
		{
			name: "nothing",
			keep: game.Dice{0, 0, 0, 0, 0},
			want: "nothing",
		},
		{
			name: "single moonberry",
			keep: game.Dice{0, 0, 0, 1, 0},
			want: "1M",
		},
		{
			name: "multiple types",
			keep: game.Dice{2, 1, 1, 1, 0},
			want: "2J 1S 1P 1M",
		},
		{
			name: "with pest",
			keep: game.Dice{1, 1, 1, 1, 1},
			want: "1J 1S 1P 1M 1X",
		},
		{
			name: "all moonberries",
			keep: game.Dice{0, 0, 0, 5, 0},
			want: "5M",
		},
		{
			name: "moonberry and pickleberry",
			keep: game.Dice{0, 0, 2, 1, 0},
			want: "2P 1M",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := FormatKeep(tt.keep); got != tt.want {
				t.Errorf("FormatKeep(%v) = %q, want %q", tt.keep, got, tt.want)
			}
		})
	}
}

func TestTheoreticalMax(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		dice      game.Dice
		rollsLeft int
		cs        game.CategorySet
		want      float64
	}{
		{
			name:      "all categories with rerolls",
			dice:      game.Dice{2, 1, 1, 1, 0},
			rollsLeft: 2,
			cs:        game.AllCategories,
			want:      237.0,
		},
		{
			name:      "all categories no rerolls best dice",
			dice:      game.Dice{0, 0, 0, 5, 0},
			rollsLeft: 0,
			cs:        game.AllCategories,
			want:      237.0, // 35 (moonberries) + 202 (remaining max scores without moonberry cat)
		},
		{
			name:      "single category remaining with rerolls",
			dice:      game.Dice{2, 1, 1, 1, 0},
			rollsLeft: 1,
			cs:        1 << game.CatFreeRoll,
			want:      35.0,
		},
		{
			name:      "single category no rerolls",
			dice:      game.Dice{2, 1, 1, 1, 0},
			rollsLeft: 0,
			cs:        1 << game.CatFreeRoll,
			want:      17.0,
		},
		{
			name:      "no categories remaining",
			dice:      game.Dice{2, 1, 1, 1, 0},
			rollsLeft: 2,
			cs:        0,
			want:      0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := theoreticalMax(tt.dice, tt.rollsLeft, tt.cs)
			if got != tt.want {
				t.Errorf("theoreticalMax(%v, %d, %v) = %v, want %v", tt.dice, tt.rollsLeft, tt.cs, got, tt.want)
			}
		})
	}
}

func TestMaxCategoryScore(t *testing.T) {
	t.Parallel()

	tests := []struct {
		category game.Category
		want     int
	}{
		{game.CatJumbleberry, 10},
		{game.CatSugarberry, 10},
		{game.CatPickleberry, 20},
		{game.CatMoonberry, 35},
		{game.CatBasketOfThree, 35},
		{game.CatBasketOfFour, 35},
		{game.CatBasketOfFive, 35},
		{game.CatMixedBasket, 22},
		{game.CatFreeRoll, 35},
	}

	for _, tt := range tests {
		t.Run(tt.category.String(), func(t *testing.T) {
			t.Parallel()
			if got := maxCategoryScore[tt.category]; got != tt.want {
				t.Errorf("maxCategoryScore[%v] = %v, want %v", tt.category, got, tt.want)
			}
		})
	}
}

func TestSumMaxScores(t *testing.T) {
	t.Parallel()

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
			cs:   1 << game.CatMoonberry,
			want: 35.0,
		},
		{
			name: "all categories",
			cs:   game.AllCategories,
			want: 237.0,
		},
		{
			name: "berry categories only",
			cs:   (1 << game.CatJumbleberry) | (1 << game.CatSugarberry) | (1 << game.CatPickleberry) | (1 << game.CatMoonberry),
			want: 75.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := sumMaxScores(tt.cs)
			if got != tt.want {
				t.Errorf("sumMaxScores(%v) = %v, want %v", tt.cs, got, tt.want)
			}
		})
	}
}
