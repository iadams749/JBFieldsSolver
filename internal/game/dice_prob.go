package game

import "math"

// FaceProb holds the probability of each face appearing on a single die.
// Each die has 10 faces: 3 Jumbleberry, 3 Sugarberry, 2 Pickleberry, 1 Moonberry, 1 Pest.
var FaceProb = [NumBerryTypes]float64{
	Jumbleberry: 0.3,
	Sugarberry:  0.3,
	Pickleberry: 0.2,
	Moonberry:   0.1,
	Pest:        0.1,
}

// factorials[n] = n! for n in 0..5.
var factorials = [NumDice + 1]float64{1, 1, 2, 6, 24, 120}

// DiceProb returns the multinomial probability of rolling dice outcome d
// when rolling n dice.
//
// P(d | n) = n! / (d[0]!*d[1]!*...*d[4]!) * prod(FaceProb[i]^d[i])
func DiceProb(d Dice, n int) float64 {
	prob := factorials[n]
	for b := Berry(0); b < NumBerryTypes; b++ {
		prob /= factorials[d[b]]
		if d[b] > 0 {
			prob *= math.Pow(FaceProb[b], float64(d[b]))
		}
	}
	return prob
}

// RerollOutcome pairs a dice count tuple with its probability.
type RerollOutcome struct {
	Dice Dice
	Prob float64
}

// allDiceCache holds the 126 canonical dice outcomes, computed once at init.
var allDiceCache []Dice

// diceToIndex maps packDice(d) → index in allDiceCache.
// Pack space is 6^5 = 7776 (each count 0..5, base-6 encoding).
var diceToIndex [7776]uint16

// firstRollProb[i] = P(allDiceCache[i]) when rolling all 5 dice.
var firstRollProb [126]float64

// rerollCache[n] holds precomputed outcomes for rolling n dice (0 <= n <= 5).
var rerollCache [NumDice + 1][]RerollOutcome

func init() {
	allDiceCache = EnumerateAllDice()

	for i, d := range allDiceCache {
		diceToIndex[packDice(d)] = uint16(i)
	}

	for i, d := range allDiceCache {
		firstRollProb[i] = DiceProb(d, NumDice)
	}

	for n := 0; n <= NumDice; n++ {
		rerollCache[n] = enumerateRerolls(n)
	}
}

// packDice encodes a Dice into a unique integer via base-6 encoding.
func packDice(d Dice) int {
	return int(d[0]) + int(d[1])*6 + int(d[2])*36 + int(d[3])*216 + int(d[4])*1296
}

// DiceIndex returns the canonical index of d in AllDice().
func DiceIndex(d Dice) int {
	return int(diceToIndex[packDice(d)])
}

// AllDice returns the cached slice of all 126 distinct 5-dice outcomes.
func AllDice() []Dice {
	return allDiceCache
}

// NumAllDice is the number of distinct 5-dice outcomes (126).
func NumAllDice() int {
	return len(allDiceCache)
}

// Rerolls returns precomputed outcomes for rolling n dice.
func Rerolls(n int) []RerollOutcome {
	return rerollCache[n]
}

// FirstRollProb returns P(AllDice()[idx]) for the initial roll of all 5 dice.
func FirstRollProb(idx int) float64 {
	return firstRollProb[idx]
}

// AddDice returns the component-wise sum of two Dice.
func AddDice(a, b Dice) Dice {
	var result Dice
	for i := Berry(0); i < NumBerryTypes; i++ {
		result[i] = a[i] + b[i]
	}
	return result
}

// enumerateRerolls generates all possible outcomes when rolling n dice,
// each paired with its multinomial probability.
func enumerateRerolls(n int) []RerollOutcome {
	if n == 0 {
		return []RerollOutcome{{Dice: Dice{}, Prob: 1.0}}
	}
	var result []RerollOutcome
	var generate func(berry Berry, remaining uint8, current Dice)
	generate = func(berry Berry, remaining uint8, current Dice) {
		if berry == NumBerryTypes-1 {
			current[berry] = remaining
			result = append(result, RerollOutcome{
				Dice: current,
				Prob: DiceProb(current, n),
			})
			return
		}
		for count := uint8(0); count <= remaining; count++ {
			current[berry] = count
			generate(berry+1, remaining-count, current)
		}
	}
	generate(0, uint8(n), Dice{})
	return result
}
