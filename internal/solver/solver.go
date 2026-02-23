// Package solver implements the core optimal play algorithm for Jumbleberry Fields.
// It uses dynamic programming with precomputed expected values to determine
// whether to score or reroll, and which category or keep decision maximizes EV.
package solver

import (
	"math"
	"sort"

	"github.com/iadams749/JBFieldsSolver/internal/ev"
	"github.com/iadams749/JBFieldsSolver/internal/game"
)

// ActionType indicates what kind of action the solver recommends.
type ActionType int

const (
	ScoreAction ActionType = iota
	RerollAction
)

// Action represents the solver's recommended move.
type Action struct {
	Type     ActionType
	Keep     game.Dice     // which dice to keep (only for RerollAction)
	Category game.Category // which category to score (only for ScoreAction)
	EV       float64       // expected value of taking this action
}

// CategoryOption is one possible scoring choice when rollsLeft == 0.
type CategoryOption struct {
	Category       game.Category
	ImmediateScore int
	FutureEV       float64
	TotalValue     float64
}

// RerollOption is one possible keep/reroll choice when rollsLeft > 0.
type RerollOption struct {
	Keep        game.Dice
	NumRerolled int
	EV          float64
}

// Recommendation is the full solver output for a given game state.
type Recommendation struct {
	BestAction       Action
	TheoreticalMax   float64
	CategoryOptions  []CategoryOption // populated when rollsLeft == 0
	TopRerollOptions []RerollOption   // populated when rollsLeft > 0 (top 10)
}

// Solve computes the optimal action for the given game state.
// rollsLeft: 0 = must score, 1 = one reroll left, 2 = two rerolls left.
func Solve(dice game.Dice, rollsLeft int, cs game.CategorySet, table *ev.Table) Recommendation {
	if rollsLeft == 0 {
		return solveScoring(dice, cs, table)
	}
	return solveReroll(dice, rollsLeft, cs, table)
}

// solveScoring handles the case where the player must score (rollsLeft == 0).
func solveScoring(dice game.Dice, cs game.CategorySet, table *ev.Table) Recommendation {
	var options []CategoryOption
	bestVal := math.Inf(-1)
	var bestCat game.Category

	cs.ForEach(func(cat game.Category) {
		imm := game.Score(dice, cat)
		fut := table.EV(cs.Remove(cat))
		total := float64(imm) + fut
		options = append(options, CategoryOption{
			Category:       cat,
			ImmediateScore: imm,
			FutureEV:       fut,
			TotalValue:     total,
		})
		if total > bestVal {
			bestVal = total
			bestCat = cat
		}
	})

	sort.Slice(options, func(i, j int) bool {
		return options[i].TotalValue > options[j].TotalValue
	})

	return Recommendation{
		BestAction: Action{
			Type:     ScoreAction,
			Category: bestCat,
			EV:       bestVal,
		},
		TheoreticalMax:  theoreticalMax(dice, 0, cs),
		CategoryOptions: options,
	}
}

// solveReroll handles the case where the player can reroll (rollsLeft > 0).
func solveReroll(dice game.Dice, rollsLeft int, cs game.CategorySet, table *ev.Table) Recommendation {
	allDice := game.AllDice()
	numDice := len(allDice)

	// Build v0: the scoring layer (rollsLeft == 0) for all 126 dice outcomes.
	v0 := make([]float64, numDice)
	for i, d := range allDice {
		bestVal := math.Inf(-1)
		cs.ForEach(func(cat game.Category) {
			val := float64(game.Score(d, cat)) + table.EV(cs.Remove(cat))
			if val > bestVal {
				bestVal = val
			}
		})
		v0[i] = bestVal
	}

	// Build successive reroll layers v1..v_{rollsLeft-1}.
	layers := make([][]float64, rollsLeft)
	layers[0] = v0
	for r := 1; r < rollsLeft; r++ {
		layers[r] = make([]float64, numDice)
		ev.ComputeRerollLayer(allDice, layers[r-1], layers[r])
	}

	// For the user's specific dice, enumerate all keep decisions
	// and evaluate each against the appropriate layer.
	prevLayer := layers[rollsLeft-1]
	var allOptions []RerollOption
	bestEV := math.Inf(-1)
	var bestKeep game.Dice

	ev.EnumerateKeeps(dice, func(keep game.Dice, numKept int) {
		numRerolled := game.NumDice - numKept
		if numRerolled == 0 {
			return // skip "keep all" — that's a score decision, not a reroll
		}

		var keepEV float64
		rerolls := game.Rerolls(numRerolled)
		for _, ro := range rerolls {
			resultDice := game.AddDice(keep, ro.Dice)
			resultIdx := game.DiceIndex(resultDice)
			keepEV += ro.Prob * prevLayer[resultIdx]
		}

		allOptions = append(allOptions, RerollOption{
			Keep:        keep,
			NumRerolled: numRerolled,
			EV:          keepEV,
		})

		if keepEV > bestEV {
			bestEV = keepEV
			bestKeep = keep
		}
	})

	// Check if scoring now (keeping all dice) beats every reroll option.
	// If so, return a score recommendation instead.
	scoreRec := solveScoring(dice, cs, table)
	if scoreRec.BestAction.EV >= bestEV {
		return scoreRec
	}

	sort.Slice(allOptions, func(i, j int) bool {
		return allOptions[i].EV > allOptions[j].EV
	})

	topN := 10
	if len(allOptions) < topN {
		topN = len(allOptions)
	}

	return Recommendation{
		BestAction: Action{
			Type: RerollAction,
			Keep: bestKeep,
			EV:   bestEV,
		},
		TheoreticalMax:   theoreticalMax(dice, rollsLeft, cs),
		TopRerollOptions: allOptions[:topN],
	}
}
