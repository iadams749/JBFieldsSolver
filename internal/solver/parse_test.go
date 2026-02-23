package solver

import (
	"testing"

	"github.com/iadams749/JBFieldsSolver/internal/game"
)

func TestParseDice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    game.Dice
		wantErr bool
	}{
		// Sequence format tests
		{
			name:    "sequence all same",
			input:   "JJJJJ",
			want:    game.Dice{5, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name:    "sequence mixed",
			input:   "JJSPM",
			want:    game.Dice{2, 1, 1, 1, 0},
			wantErr: false,
		},
		{
			name:    "sequence lowercase",
			input:   "jjspm",
			want:    game.Dice{2, 1, 1, 1, 0},
			wantErr: false,
		},
		{
			name:    "sequence with pests",
			input:   "JSPXX",
			want:    game.Dice{1, 1, 1, 0, 2},
			wantErr: false,
		},
		{
			name:    "sequence all pests",
			input:   "XXXXX",
			want:    game.Dice{0, 0, 0, 0, 5},
			wantErr: false,
		},
		{
			name:    "sequence invalid letter",
			input:   "JJSPA",
			want:    game.Dice{},
			wantErr: true,
		},
		{
			name:    "sequence too short",
			input:   "JJSP",
			want:    game.Dice{},
			wantErr: true,
		},
		{
			name:    "sequence too long",
			input:   "JJSSPP",
			want:    game.Dice{},
			wantErr: true,
		},

		// Count format tests
		{
			name:    "count standard",
			input:   "2J 1S 1P 1M",
			want:    game.Dice{2, 1, 1, 1, 0},
			wantErr: false,
		},
		{
			name:    "count with pest",
			input:   "1J 1S 1P 1M 1X",
			want:    game.Dice{1, 1, 1, 1, 1},
			wantErr: false,
		},
		{
			name:    "count all same",
			input:   "5M",
			want:    game.Dice{0, 0, 0, 5, 0},
			wantErr: false,
		},
		{
			name:    "count with zeros explicit",
			input:   "2J 2S 1P 0M 0X",
			want:    game.Dice{2, 2, 1, 0, 0},
			wantErr: false,
		},
		{
			name:    "count lowercase",
			input:   "2j 1s 1p 1m",
			want:    game.Dice{2, 1, 1, 1, 0},
			wantErr: false,
		},
		{
			name:    "count wrong total",
			input:   "3J 2S 2P",
			want:    game.Dice{},
			wantErr: true,
		},
		{
			name:    "count invalid berry",
			input:   "2J 1S 1P 1Z",
			want:    game.Dice{},
			wantErr: true,
		},
		{
			name:    "count invalid format",
			input:   "2J J2",
			want:    game.Dice{},
			wantErr: true,
		},

		// Edge cases
		{
			name:    "empty input",
			input:   "",
			want:    game.Dice{},
			wantErr: true,
		},
		{
			name:    "whitespace only",
			input:   "   ",
			want:    game.Dice{},
			wantErr: true,
		},
		{
			name:    "count with extra whitespace",
			input:   "  2J  1S  1P  1M  ",
			want:    game.Dice{2, 1, 1, 1, 0},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseDice(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDice(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseDice(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseCategories(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    game.CategorySet
		wantErr bool
	}{
		// "all" syntax
		{
			name:    "all categories",
			input:   "all",
			want:    game.AllCategories,
			wantErr: false,
		},
		{
			name:    "all uppercase",
			input:   "ALL",
			want:    game.AllCategories,
			wantErr: false,
		},
		{
			name:    "all minus one",
			input:   "all-j",
			want:    game.AllCategories.Remove(game.CatJumbleberry),
			wantErr: false,
		},
		{
			name:    "all minus multiple",
			input:   "all-j-s",
			want:    game.AllCategories.Remove(game.CatJumbleberry).Remove(game.CatSugarberry),
			wantErr: false,
		},
		{
			name:    "all minus with long names",
			input:   "all-jumbleberry-sugarberry",
			want:    game.AllCategories.Remove(game.CatJumbleberry).Remove(game.CatSugarberry),
			wantErr: false,
		},
		{
			name:    "all minus invalid category",
			input:   "all-j-invalid",
			want:    0,
			wantErr: true,
		},

		// Comma-separated list
		{
			name:    "single category short",
			input:   "j",
			want:    1 << game.CatJumbleberry,
			wantErr: false,
		},
		{
			name:    "single category long",
			input:   "jumbleberry",
			want:    1 << game.CatJumbleberry,
			wantErr: false,
		},
		{
			name:    "multiple categories",
			input:   "j,s,p",
			want:    (1 << game.CatJumbleberry) | (1 << game.CatSugarberry) | (1 << game.CatPickleberry),
			wantErr: false,
		},
		{
			name:    "basket categories",
			input:   "3k,4k,5k",
			want:    (1 << game.CatBasketOfThree) | (1 << game.CatBasketOfFour) | (1 << game.CatBasketOfFive),
			wantErr: false,
		},
		{
			name:    "mixed and free roll",
			input:   "mix,fr",
			want:    (1 << game.CatMixedBasket) | (1 << game.CatFreeRoll),
			wantErr: false,
		},
		{
			name:    "all four berries",
			input:   "j,s,p,m",
			want:    (1 << game.CatJumbleberry) | (1 << game.CatSugarberry) | (1 << game.CatPickleberry) | (1 << game.CatMoonberry),
			wantErr: false,
		},
		{
			name:    "mixed long and short names",
			input:   "jumbleberry,s,pickleberry,4k",
			want:    (1 << game.CatJumbleberry) | (1 << game.CatSugarberry) | (1 << game.CatPickleberry) | (1 << game.CatBasketOfFour),
			wantErr: false,
		},
		{
			name:    "uppercase",
			input:   "J,S,P",
			want:    (1 << game.CatJumbleberry) | (1 << game.CatSugarberry) | (1 << game.CatPickleberry),
			wantErr: false,
		},
		{
			name:    "with spaces",
			input:   "j, s, p",
			want:    (1 << game.CatJumbleberry) | (1 << game.CatSugarberry) | (1 << game.CatPickleberry),
			wantErr: false,
		},
		{
			name:    "invalid category in list",
			input:   "j,s,invalid",
			want:    0,
			wantErr: true,
		},

		// Edge cases
		{
			name:    "empty input",
			input:   "",
			want:    0,
			wantErr: true,
		},
		{
			name:    "whitespace only",
			input:   "   ",
			want:    0,
			wantErr: true,
		},
		{
			name:    "trailing comma",
			input:   "j,s,",
			want:    (1 << game.CatJumbleberry) | (1 << game.CatSugarberry),
			wantErr: false,
		},
		{
			name:    "empty elements",
			input:   "j,,s",
			want:    (1 << game.CatJumbleberry) | (1 << game.CatSugarberry),
			wantErr: false,
		},
		{
			name:    "only commas",
			input:   ",,",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseCategories(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCategories(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseCategories(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestCategoryAliases(t *testing.T) {
	t.Parallel()

	// Test that common aliases map correctly
	tests := []struct {
		alias string
		want  game.Category
	}{
		// Berry categories
		{"j", game.CatJumbleberry},
		{"jb", game.CatJumbleberry},
		{"jumbleberry", game.CatJumbleberry},
		{"s", game.CatSugarberry},
		{"sb", game.CatSugarberry},
		{"sugarberry", game.CatSugarberry},
		{"p", game.CatPickleberry},
		{"pb", game.CatPickleberry},
		{"pickleberry", game.CatPickleberry},
		{"m", game.CatMoonberry},
		{"mb", game.CatMoonberry},
		{"moonberry", game.CatMoonberry},

		// Basket categories
		{"3k", game.CatBasketOfThree},
		{"b3", game.CatBasketOfThree},
		{"three", game.CatBasketOfThree},
		{"basketofthree", game.CatBasketOfThree},
		{"4k", game.CatBasketOfFour},
		{"b4", game.CatBasketOfFour},
		{"four", game.CatBasketOfFour},
		{"basketoffour", game.CatBasketOfFour},
		{"5k", game.CatBasketOfFive},
		{"b5", game.CatBasketOfFive},
		{"five", game.CatBasketOfFive},
		{"basketoffive", game.CatBasketOfFive},

		// Special categories
		{"mix", game.CatMixedBasket},
		{"mixed", game.CatMixedBasket},
		{"mixedbasket", game.CatMixedBasket},
		{"fr", game.CatFreeRoll},
		{"free", game.CatFreeRoll},
		{"freeroll", game.CatFreeRoll},
	}

	for _, tt := range tests {
		t.Run(tt.alias, func(t *testing.T) {
			t.Parallel()
			cs, err := ParseCategories(tt.alias)
			if err != nil {
				t.Errorf("ParseCategories(%q) error = %v", tt.alias, err)
				return
			}
			if !cs.Has(tt.want) {
				t.Errorf("ParseCategories(%q) should contain %v", tt.alias, tt.want)
			}
			if cs.Count() != 1 {
				t.Errorf("ParseCategories(%q) should contain exactly 1 category, got %d", tt.alias, cs.Count())
			}
		})
	}
}
