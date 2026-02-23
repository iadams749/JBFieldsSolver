package game

import "fmt"

// Berry represents a face on one of the 5 dice.
type Berry uint8

const (
	Jumbleberry Berry = iota
	Sugarberry
	Pickleberry
	Moonberry
	Pest
	NumBerryTypes // sentinel: always last
)

// BerryPoints maps each berry type to its point value.
var BerryPoints = [NumBerryTypes]int{
	Jumbleberry: 2,
	Sugarberry:  2,
	Pickleberry: 4,
	Moonberry:   7,
	Pest:        0,
}

var berryNames = [NumBerryTypes]string{
	Jumbleberry: "Jumbleberry",
	Sugarberry:  "Sugarberry",
	Pickleberry: "Pickleberry",
	Moonberry:   "Moonberry",
	Pest:        "Pest",
}

func (b Berry) String() string {
	if b < NumBerryTypes {
		return berryNames[b]
	}
	return fmt.Sprintf("Berry(%d)", b)
}

// NumDice is the number of dice rolled each turn.
const NumDice = 5

// Dice represents the counts of each berry type across all 5 dice.
// Dice[Jumbleberry] = number of dice showing Jumbleberry, etc.
// This is a fixed-size value type — copies are cheap (no heap allocation).
type Dice [NumBerryTypes]uint8

// Total returns the number of dice represented (should always be NumDice for a full roll).
func (d Dice) Total() int {
	sum := 0
	for _, c := range d {
		sum += int(c)
	}
	return sum
}

// Points returns the total point value of the dice.
func (d Dice) Points() int {
	sum := 0
	for b := Berry(0); b < NumBerryTypes; b++ {
		sum += int(d[b]) * BerryPoints[b]
	}
	return sum
}

func (d Dice) String() string {
	return fmt.Sprintf("{J:%d S:%d P:%d M:%d X:%d}",
		d[Jumbleberry], d[Sugarberry], d[Pickleberry], d[Moonberry], d[Pest])
}

// EnumerateAllDice returns all distinct dice outcomes (combinations with
// replacement of NumDice dice across NumBerryTypes faces).
// There are C(NumDice + NumBerryTypes - 1, NumBerryTypes - 1) = 126 outcomes.
func EnumerateAllDice() []Dice {
	var result []Dice

	// generate recursively builds dice combinations by deciding how many
	// dice show each berry type, from Jumbleberry (0) through Pest (4).
	//
	// berry:     the berry type we're currently assigning a count to
	// remaining: how many dice are still unassigned
	// current:   the dice combination being built up
	var generate func(berry Berry, remaining uint8, current Dice)
	generate = func(berry Berry, remaining uint8, current Dice) {
		// Base case: last berry type gets all remaining dice.
		// This avoids an unnecessary extra level of recursion.
		if berry == NumBerryTypes-1 {
			current[berry] = remaining
			result = append(result, current)
			return
		}

		// Try every valid count for this berry type (0 up to however many
		// dice are left), then recurse to assign the next berry type.
		for count := uint8(0); count <= remaining; count++ {
			current[berry] = count
			generate(berry+1, remaining-count, current)
		}
	}

	// Start with berry type 0, all 5 dice unassigned, and an empty Dice.
	generate(0, NumDice, Dice{})
	return result
}
