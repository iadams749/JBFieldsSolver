// Package main provides an interactive command-line REPL for the Jumbleberry Fields solver.
// Users enter their current dice, rolls remaining, and available categories,
// and receive optimal play recommendations.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/iadams749/JBFieldsSolver/internal/evloader"
	"github.com/iadams749/JBFieldsSolver/internal/game"
	"github.com/iadams749/JBFieldsSolver/internal/solver"
)

const evTablePath = "ev_table.json"

func main() {
	table, err := evloader.Load(evTablePath)
	if err != nil {
		fmt.Printf("Fatal: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("EV with all categories: %.2f\n\n", table.EV(game.AllCategories))

	// REPL loop
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("=== Jumbleberry Fields Solver ===")
	fmt.Println("Enter your game state to get optimal play advice.")
	fmt.Println("Type 'quit' or 'exit' at any prompt to quit.")
	fmt.Println()
	fmt.Println("--- Dice ---")
	fmt.Println("  Letters: J=Jumbleberry  S=Sugarberry  P=Pickleberry  M=Moonberry  X=Pest")
	fmt.Println("  Sequence format:  JJSPM       (one letter per die, 5 total)")
	fmt.Println("  Count format:     2J 1S 1P 1M (space-separated, omitted types = 0)")
	fmt.Println()
	fmt.Println("--- Rolls Left ---")
	fmt.Println("  0 = no rerolls (must score)    1 = one reroll left    2 = two rerolls left")
	fmt.Println()
	fmt.Println("--- Categories ---")
	fmt.Println("  Shorthand: j  s  p  m  3k  4k  5k  mix  fr")
	fmt.Println("  Examples:  'all'          all 9 categories")
	fmt.Println("             'all-j-s'      all except Jumbleberry and Sugarberry")
	fmt.Println("             'j,m,3k,fr'    only those 4 categories")
	fmt.Println()

	for {
		// Prompt for dice
		fmt.Print("Dice: ")
		if !scanner.Scan() {
			break
		}
		diceInput := strings.TrimSpace(scanner.Text())
		if diceInput == "quit" || diceInput == "exit" {
			break
		}

		dice, err := solver.ParseDice(diceInput)
		if err != nil {
			fmt.Printf("  Error: %v\n\n", err)
			continue
		}

		// Prompt for rolls left
		fmt.Print("Rolls left (0-2): ")
		if !scanner.Scan() {
			break
		}
		rollsInput := strings.TrimSpace(scanner.Text())
		if rollsInput == "quit" || rollsInput == "exit" {
			break
		}

		rollsLeft, err := strconv.Atoi(rollsInput)
		if err != nil || rollsLeft < 0 || rollsLeft > 2 {
			fmt.Println("  Error: rolls left must be 0, 1, or 2")
			fmt.Println()
			continue
		}

		// Prompt for categories
		fmt.Print("Categories remaining: ")
		if !scanner.Scan() {
			break
		}
		catInput := strings.TrimSpace(scanner.Text())
		if catInput == "quit" || catInput == "exit" {
			break
		}

		cs, err := solver.ParseCategories(catInput)
		if err != nil {
			fmt.Printf("  Error: %v\n\n", err)
			continue
		}

		// Solve and display
		rec := solver.Solve(dice, rollsLeft, cs, table)
		solver.FormatRecommendation(os.Stdout, rec, dice, rollsLeft, cs)
		fmt.Println()
	}
}
