// This file provides parsing utilities for converting user input
// into game types (Dice and CategorySet).
package solver

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/iadams749/JBFieldsSolver/internal/game"
)

// berryLetters maps single characters to Berry types.
var berryLetters = map[byte]game.Berry{
	'J': game.Jumbleberry,
	'j': game.Jumbleberry,
	'S': game.Sugarberry,
	's': game.Sugarberry,
	'P': game.Pickleberry,
	'p': game.Pickleberry,
	'M': game.Moonberry,
	'm': game.Moonberry,
	'X': game.Pest,
	'x': game.Pest,
}

// ParseDice parses a dice string into a Dice value.
//
// Supported formats:
//   - Count format: "2J 1S 1P 1M 0X" (space-separated, missing types default to 0)
//   - Sequence format: "JJSPM" (exactly 5 characters, one per die)
func ParseDice(input string) (game.Dice, error) {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return game.Dice{}, fmt.Errorf("empty dice input")
	}

	// Check if it's sequence format (exactly 5 letters, no digits or spaces)
	if len(input) == game.NumDice && isAllLetters(input) {
		return parseDiceSequence(input)
	}

	return parseDiceCounts(input)
}

func isAllLetters(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// parseDiceSequence parses "JJSPM" format.
func parseDiceSequence(input string) (game.Dice, error) {
	var d game.Dice
	for i := 0; i < len(input); i++ {
		b, ok := berryLetters[input[i]]
		if !ok {
			return game.Dice{}, fmt.Errorf("unknown die face %q at position %d (use J, S, P, M, or X)", string(input[i]), i+1)
		}
		d[b]++
	}
	return d, nil
}

// parseDiceCounts parses "2J 1S 1P 1M 0X" format.
func parseDiceCounts(input string) (game.Dice, error) {
	var d game.Dice
	tokens := strings.Fields(input)
	for _, token := range tokens {
		if len(token) < 2 {
			return game.Dice{}, fmt.Errorf("invalid token %q (expected format like '2J' or '1M')", token)
		}

		// Split into number and letter parts
		letterIdx := 0
		for letterIdx < len(token) && (token[letterIdx] >= '0' && token[letterIdx] <= '9') {
			letterIdx++
		}
		if letterIdx == 0 || letterIdx >= len(token) {
			return game.Dice{}, fmt.Errorf("invalid token %q (expected format like '2J' or '1M')", token)
		}

		count, err := strconv.Atoi(token[:letterIdx])
		if err != nil {
			return game.Dice{}, fmt.Errorf("invalid count in %q: %v", token, err)
		}

		letter := token[letterIdx:]
		if len(letter) != 1 {
			return game.Dice{}, fmt.Errorf("invalid berry letter in %q (use J, S, P, M, or X)", token)
		}

		b, ok := berryLetters[letter[0]]
		if !ok {
			return game.Dice{}, fmt.Errorf("unknown berry %q in %q (use J, S, P, M, or X)", letter, token)
		}

		d[b] += uint8(count)
	}

	total := d.Total()
	if total != game.NumDice {
		return game.Dice{}, fmt.Errorf("dice must sum to %d, got %d", game.NumDice, total)
	}

	return d, nil
}

// categoryAliases maps short names and full names to Category values.
var categoryAliases = map[string]game.Category{
	"jumbleberry":   game.CatJumbleberry,
	"j":             game.CatJumbleberry,
	"jb":            game.CatJumbleberry,
	"sugarberry":    game.CatSugarberry,
	"s":             game.CatSugarberry,
	"sb":            game.CatSugarberry,
	"pickleberry":   game.CatPickleberry,
	"p":             game.CatPickleberry,
	"pb":            game.CatPickleberry,
	"moonberry":     game.CatMoonberry,
	"m":             game.CatMoonberry,
	"mb":            game.CatMoonberry,
	"basketofthree": game.CatBasketOfThree,
	"3k":            game.CatBasketOfThree,
	"b3":            game.CatBasketOfThree,
	"three":         game.CatBasketOfThree,
	"basketoffour":  game.CatBasketOfFour,
	"4k":            game.CatBasketOfFour,
	"b4":            game.CatBasketOfFour,
	"four":          game.CatBasketOfFour,
	"basketoffive":  game.CatBasketOfFive,
	"5k":            game.CatBasketOfFive,
	"b5":            game.CatBasketOfFive,
	"five":          game.CatBasketOfFive,
	"mixedbasket":   game.CatMixedBasket,
	"mix":           game.CatMixedBasket,
	"mixed":         game.CatMixedBasket,
	"freeroll":      game.CatFreeRoll,
	"fr":            game.CatFreeRoll,
	"free":          game.CatFreeRoll,
}

// ParseCategories parses a category set from user input.
//
// Supported formats:
//   - "all" → AllCategories
//   - "all-j-s" → all except Jumbleberry and Sugarberry
//   - "j,s,p,m,3k,4k,5k,mix,fr" → explicit comma-separated list
func ParseCategories(input string) (game.CategorySet, error) {
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return 0, fmt.Errorf("empty categories input")
	}

	// Check for "all" or "all-..." syntax
	if strings.HasPrefix(input, "all") {
		cs := game.AllCategories
		rest := input[3:]
		if rest == "" {
			return cs, nil
		}
		// Parse minus-separated removals: "all-j-s-mix"
		parts := strings.Split(rest, "-")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			cat, ok := categoryAliases[part]
			if !ok {
				return 0, fmt.Errorf("unknown category %q", part)
			}
			cs = cs.Remove(cat)
		}
		return cs, nil
	}

	// Comma-separated list
	var cs game.CategorySet
	parts := strings.Split(input, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		cat, ok := categoryAliases[part]
		if !ok {
			return 0, fmt.Errorf("unknown category %q", part)
		}
		cs = cs.Add(cat)
	}

	if cs == 0 {
		return 0, fmt.Errorf("no valid categories specified")
	}

	return cs, nil
}
