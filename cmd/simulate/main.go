// Package main simulates many games of Jumbleberry Fields using optimal play
// to validate the theoretical expected value and compute standard deviation.
package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"time"

	"github.com/iadams749/JBFieldsSolver/internal/ev"
	"github.com/iadams749/JBFieldsSolver/internal/game"
	"github.com/iadams749/JBFieldsSolver/internal/solver"
)

func main() {
	numGames := flag.Int("n", 100000, "number of games to simulate")
	evPath := flag.String("ev", "ev_table.json", "path to EV table JSON")
	seed := flag.Uint64("seed", 0, "random seed (0 = use current time)")
	flag.Parse()

	table, err := ev.LoadJSON(*evPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading EV table: %v\n", err)
		fmt.Fprintln(os.Stderr, "Run the CLI or API first to generate ev_table.json")
		os.Exit(1)
	}

	var rng *rand.Rand
	if *seed != 0 {
		rng = rand.New(rand.NewPCG(*seed, 0))
	} else {
		rng = rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0))
	}

	theoreticalEV := table.EV(game.AllCategories)
	fmt.Printf("Theoretical EV (all categories): %.4f\n", theoreticalEV)
	fmt.Printf("Simulating %d games...\n\n", *numGames)

	scores := make([]int, *numGames)
	start := time.Now()

	// Track best game
	bestScore := 0
	var bestBreakdown []categoryResult

	for i := range *numGames {
		score, breakdown := simulateGame(rng, table)
		scores[i] = score
		if score > bestScore {
			bestScore = score
			bestBreakdown = breakdown
		}
		if (i+1)%(*numGames/10) == 0 {
			elapsed := time.Since(start)
			fmt.Printf("  %d/%d games  (%v elapsed)\n", i+1, *numGames, elapsed.Round(time.Millisecond))
		}
	}

	elapsed := time.Since(start)

	// Compute statistics
	sum := 0
	for _, s := range scores {
		sum += s
	}
	mean := float64(sum) / float64(*numGames)

	varSum := 0.0
	for _, s := range scores {
		diff := float64(s) - mean
		varSum += diff * diff
	}
	stddev := math.Sqrt(varSum / float64(*numGames))
	stderr := stddev / math.Sqrt(float64(*numGames))

	minScore, maxScore := scores[0], scores[0]
	for _, s := range scores[1:] {
		if s < minScore {
			minScore = s
		}
		if s > maxScore {
			maxScore = s
		}
	}

	// Histogram buckets
	bucketSize := 10
	buckets := make(map[int]int)
	for _, s := range scores {
		bucket := (s / bucketSize) * bucketSize
		buckets[bucket]++
	}

	fmt.Println()
	fmt.Println("=== Results ===")
	fmt.Printf("Games simulated:  %d\n", *numGames)
	fmt.Printf("Time elapsed:     %v\n", elapsed.Round(time.Millisecond))
	fmt.Println()
	fmt.Printf("Theoretical EV:   %.4f\n", theoreticalEV)
	fmt.Printf("Simulated mean:   %.4f\n", mean)
	fmt.Printf("Difference:       %+.4f\n", mean-theoreticalEV)
	fmt.Printf("Std error:        %.4f\n", stderr)
	fmt.Println()
	fmt.Printf("Standard dev:     %.4f\n", stddev)
	fmt.Printf("Min score:        %d\n", minScore)
	fmt.Printf("Max score:        %d\n", maxScore)
	fmt.Println()

	// Print best game breakdown
	fmt.Println("Best game breakdown:")
	for i, cr := range bestBreakdown {
		fmt.Printf("  Round %d:  %-18s  scored %3d  with %s\n",
			i+1, cr.category, cr.score, formatDice(cr.dice))
	}
	fmt.Printf("  %-26s  = %d\n", "TOTAL", bestScore)
	fmt.Println()

	// Print histogram
	fmt.Println("Score distribution:")
	minBucket := (minScore / bucketSize) * bucketSize
	maxBucket := (maxScore / bucketSize) * bucketSize
	maxCount := 0
	for _, c := range buckets {
		if c > maxCount {
			maxCount = c
		}
	}
	barWidth := 50
	for b := minBucket; b <= maxBucket; b += bucketSize {
		count := buckets[b]
		barLen := count * barWidth / maxCount
		bar := ""
		for range barLen {
			bar += "█"
		}
		pct := float64(count) / float64(*numGames) * 100
		fmt.Printf("  %3d-%3d: %s %d (%.1f%%)\n", b, b+bucketSize-1, bar, count, pct)
	}
}

// categoryResult records what dice were scored in a category for one game.
type categoryResult struct {
	category game.Category
	dice     game.Dice
	score    int
}

// simulateGame plays one full game using optimal strategy,
// returning the total score and the per-category breakdown.
func simulateGame(rng *rand.Rand, table *ev.Table) (int, []categoryResult) {
	cs := game.AllCategories
	totalScore := 0
	var breakdown []categoryResult

	for round := 0; round < int(game.NumCategories); round++ {
		dice := rollAllDice(rng)

		// Up to 2 rerolls
		for rollsLeft := 2; rollsLeft > 0; rollsLeft-- {
			rec := solver.Solve(dice, rollsLeft, cs, table)
			if rec.BestAction.Type == solver.ScoreAction {
				break // solver says score now
			}
			// Reroll: keep the recommended dice, reroll the rest
			keep := rec.BestAction.Keep
			numReroll := game.NumDice - keep.Total()
			rerolled := rollNDice(rng, numReroll)
			dice = game.AddDice(keep, rerolled)
		}

		// Must score now (rollsLeft == 0)
		rec := solver.Solve(dice, 0, cs, table)
		cat := rec.BestAction.Category
		score := game.Score(dice, cat)
		totalScore += score
		breakdown = append(breakdown, categoryResult{
			category: cat,
			dice:     dice,
			score:    score,
		})
		cs = cs.Remove(cat)
	}

	return totalScore, breakdown
}

// rollAllDice rolls 5 dice and returns the result.
func rollAllDice(rng *rand.Rand) game.Dice {
	return rollNDice(rng, game.NumDice)
}

// rollNDice rolls n dice and returns the result as a Dice count array.
func rollNDice(rng *rand.Rand, n int) game.Dice {
	var d game.Dice
	for range n {
		face := rollOneDie(rng)
		d[face]++
	}
	return d
}

// rollOneDie returns a random berry face based on FaceProb distribution.
func rollOneDie(rng *rand.Rand) game.Berry {
	r := rng.Float64()
	cumulative := 0.0
	for b := game.Berry(0); b < game.NumBerryTypes; b++ {
		cumulative += game.FaceProb[b]
		if r < cumulative {
			return b
		}
	}
	return game.Pest // rounding safety
}

// formatDice returns a human-readable string for a dice outcome, e.g. "3M 2P".
func formatDice(d game.Dice) string {
	letters := [game.NumBerryTypes]string{"J", "S", "P", "M", "X"}
	result := ""
	for b := game.Berry(0); b < game.NumBerryTypes; b++ {
		if d[b] > 0 {
			if result != "" {
				result += " "
			}
			result += fmt.Sprintf("%d%s", d[b], letters[b])
		}
	}
	if result == "" {
		return "(empty)"
	}
	return result
}
