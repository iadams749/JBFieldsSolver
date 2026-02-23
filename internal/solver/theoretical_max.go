// This file provides theoretical maximum score calculations,
// representing the upper bound if all future dice were perfect.
package solver

import "github.com/iadams749/JBFieldsSolver/internal/game"

// maxCategoryScore is the maximum achievable score for each category
// assuming perfect dice (you can roll any outcome you want).
var maxCategoryScore = [game.NumCategories]int{
	game.CatJumbleberry:   10, // 5 × 2
	game.CatSugarberry:    10, // 5 × 2
	game.CatPickleberry:   20, // 5 × 4
	game.CatMoonberry:     35, // 5 × 7
	game.CatBasketOfThree: 35, // 5 Moonberries
	game.CatBasketOfFour:  35, // 5 Moonberries
	game.CatBasketOfFive:  35, // 5 Moonberries
	game.CatMixedBasket:   22, // 1J+1S+1P+2M = 2+2+4+14
	game.CatFreeRoll:      35, // 5 Moonberries
}

// theoreticalMax computes the maximum possible score from the current state,
// assuming perfect dice on all remaining rerolls and future rounds.
func theoreticalMax(dice game.Dice, rollsLeft int, cs game.CategorySet) float64 {
	if rollsLeft > 0 {
		// With rerolls remaining, we could achieve any dice outcome.
		return sumMaxScores(cs)
	}
	// rollsLeft == 0: stuck with current dice for this round.
	best := 0.0
	cs.ForEach(func(cat game.Category) {
		total := float64(game.Score(dice, cat)) + sumMaxScores(cs.Remove(cat))
		if total > best {
			best = total
		}
	})
	return best
}

// sumMaxScores returns the sum of max scores for all categories in the set.
func sumMaxScores(cs game.CategorySet) float64 {
	sum := 0.0
	cs.ForEach(func(cat game.Category) {
		sum += float64(maxCategoryScore[cat])
	})
	return sum
}
