package game

import (
	"strings"
	"testing"
)

func TestNewGame(t *testing.T) {
	t.Parallel()

	gs := NewGame()

	if gs.RollsLeft != 3 {
		t.Errorf("NewGame().RollsLeft = %d, want 3", gs.RollsLeft)
	}
	if gs.CategoriesLeft != AllCategories {
		t.Errorf("NewGame().CategoriesLeft = %v, want AllCategories", gs.CategoriesLeft)
	}
	if gs.Score != 0 {
		t.Errorf("NewGame().Score = %d, want 0", gs.Score)
	}
	if gs.CurrentDice.Total() != 0 {
		t.Errorf("NewGame().CurrentDice.Total() = %d, want 0", gs.CurrentDice.Total())
	}
}

func TestGameStateRound(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		categoriesLeft CategorySet
		want           int
	}{
		{
			name:           "round 1 - all categories",
			categoriesLeft: AllCategories,
			want:           1,
		},
		{
			name:           "round 5 - 5 categories left",
			categoriesLeft: (1 << CatJumbleberry) | (1 << CatSugarberry) | (1 << CatPickleberry) | (1 << CatMoonberry) | (1 << CatBasketOfThree),
			want:           5,
		},
		{
			name:           "round 9 - 1 category left",
			categoriesLeft: 1 << CatFreeRoll,
			want:           9,
		},
		{
			name:           "game over",
			categoriesLeft: 0,
			want:           10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gs := GameState{CategoriesLeft: tt.categoriesLeft}
			if got := gs.Round(); got != tt.want {
				t.Errorf("GameState.Round() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameStateGameOver(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		categoriesLeft CategorySet
		want           bool
	}{
		{
			name:           "not over - all categories",
			categoriesLeft: AllCategories,
			want:           false,
		},
		{
			name:           "not over - one category",
			categoriesLeft: 1 << CatFreeRoll,
			want:           false,
		},
		{
			name:           "game over",
			categoriesLeft: 0,
			want:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gs := GameState{CategoriesLeft: tt.categoriesLeft}
			if got := gs.GameOver(); got != tt.want {
				t.Errorf("GameState.GameOver() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameStateString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		state GameState
		want  []string // substrings that should appear
	}{
		{
			name: "new game",
			state: GameState{
				CurrentDice:    Dice{2, 1, 1, 1, 0},
				RollsLeft:      2,
				CategoriesLeft: AllCategories,
				Score:          15,
			},
			want: []string{"Round 1", "Rolls left: 2", "Score: 15", "Categories left: 9"},
		},
		{
			name: "mid game",
			state: GameState{
				CurrentDice:    Dice{0, 0, 0, 5, 0},
				RollsLeft:      0,
				CategoriesLeft: (1 << CatFreeRoll) | (1 << CatMoonberry),
				Score:          100,
			},
			want: []string{"Round 8", "Rolls left: 0", "Score: 100", "Categories left: 2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.state.String()
			for _, substr := range tt.want {
				if !strings.Contains(got, substr) {
					t.Errorf("GameState.String() = %q, want to contain %q", got, substr)
				}
			}
		})
	}
}
