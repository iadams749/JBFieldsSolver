// This file provides human-readable text formatting
// for solver recommendations and game state.
package solver

import (
	"fmt"
	"io"
	"strings"

	"github.com/iadams749/JBFieldsSolver/internal/game"
)

// FormatRecommendation writes the solver's recommendation to w.
func FormatRecommendation(w io.Writer, rec Recommendation, dice game.Dice, rollsLeft int, cs game.CategorySet) {
	fmt.Fprintln(w)
	fmt.Fprintln(w, "=== Solver Recommendation ===")
	fmt.Fprintf(w, "Dice: %s  |  Rolls left: %d  |  Categories: %d remaining\n",
		dice, rollsLeft, cs.Count())
	fmt.Fprintln(w)

	switch rec.BestAction.Type {
	case ScoreAction:
		formatScoreRecommendation(w, rec)
	case RerollAction:
		formatRerollRecommendation(w, rec)
	}

	fmt.Fprintf(w, "Theoretical max: %.0f\n", rec.TheoreticalMax)
	fmt.Fprintln(w, "=============================")
}

func formatScoreRecommendation(w io.Writer, rec Recommendation) {
	best := rec.CategoryOptions[0]
	fmt.Fprintf(w, "Best action: SCORE in %s\n", best.Category)
	fmt.Fprintf(w, "  Score: %d  +  Future EV: %.2f  =  Total: %.2f\n",
		best.ImmediateScore, best.FutureEV, best.TotalValue)
	fmt.Fprintln(w)

	fmt.Fprintln(w, "All category options:")
	for i, opt := range rec.CategoryOptions {
		fmt.Fprintf(w, "  #%d  %-18s  score: %3d   future EV: %7.2f   total: %7.2f\n",
			i+1, opt.Category, opt.ImmediateScore, opt.FutureEV, opt.TotalValue)
	}
	fmt.Fprintln(w)
}

func formatRerollRecommendation(w io.Writer, rec Recommendation) {
	best := rec.TopRerollOptions[0]
	fmt.Fprintf(w, "Best action: REROLL\n")
	fmt.Fprintf(w, "  Keep %s  (reroll %d)\n", FormatKeep(best.Keep), best.NumRerolled)
	fmt.Fprintf(w, "  Expected value: %.2f\n", best.EV)
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Top reroll options:")
	for i, opt := range rec.TopRerollOptions {
		fmt.Fprintf(w, "  #%d  Keep %-24s  reroll %d  EV: %7.2f\n",
			i+1, FormatKeep(opt.Keep), opt.NumRerolled, opt.EV)
	}
	fmt.Fprintln(w)
}

// FormatKeep returns a human-readable string for a keep decision.
// e.g., "2M 1P" or "nothing" if keeping 0 dice.
func FormatKeep(keep game.Dice) string {
	if keep.Total() == 0 {
		return "nothing"
	}

	type entry struct {
		letter string
		count  uint8
	}
	letters := []entry{
		{"J", keep[game.Jumbleberry]},
		{"S", keep[game.Sugarberry]},
		{"P", keep[game.Pickleberry]},
		{"M", keep[game.Moonberry]},
		{"X", keep[game.Pest]},
	}

	var parts []string
	for _, e := range letters {
		if e.count > 0 {
			parts = append(parts, fmt.Sprintf("%d%s", e.count, e.letter))
		}
	}
	return strings.Join(parts, " ")
}
