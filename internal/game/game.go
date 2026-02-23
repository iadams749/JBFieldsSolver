package game

import "fmt"

// GameState represents the full state of a Jumbleberry Fields game.
//
// The struct is 12 bytes total and contains only value types —
// copying is a trivial memcpy with zero heap allocation or GC pressure.
//
// Layout (12 bytes):
//   - CurrentDice:      [5]uint8 = 5 bytes (counts per berry type)
//   - RollsLeft:        uint8    = 1 byte  (0-3)
//   - CategoriesLeft:   uint16   = 2 bytes (bitmask of 9 categories)
//   - Score:            uint16   = 2 bytes  (cumulative score)
//   - padding:                     2 bytes (alignment)
type GameState struct {
	CurrentDice    Dice        // counts of each berry type in current roll
	RollsLeft      uint8       // rerolls remaining this round (0-3)
	CategoriesLeft CategorySet // bitmask of unused scoring categories
	Score          uint16      // cumulative point total
}

// NewGame returns the initial game state: all categories available,
// 3 rolls to use, no dice rolled yet, score 0.
func NewGame() GameState {
	return GameState{
		RollsLeft:      3,
		CategoriesLeft: AllCategories,
	}
}

// Round returns the current round number (1-9).
func (gs GameState) Round() int {
	return int(NumCategories) - gs.CategoriesLeft.Count() + 1
}

// GameOver returns true if all 9 categories have been used.
func (gs GameState) GameOver() bool {
	return gs.CategoriesLeft == 0
}

func (gs GameState) String() string {
	return fmt.Sprintf("Round %d | Dice: %s | Rolls left: %d | Score: %d | Categories left: %d",
		gs.Round(), gs.CurrentDice, gs.RollsLeft, gs.Score, gs.CategoriesLeft.Count())
}
