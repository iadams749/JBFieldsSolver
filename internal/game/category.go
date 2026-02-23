package game

import "fmt"

// Category represents one of the 9 scoring categories.
type Category uint8

const (
	CatJumbleberry Category = iota
	CatSugarberry
	CatPickleberry
	CatMoonberry
	CatBasketOfThree
	CatBasketOfFour
	CatBasketOfFive
	CatMixedBasket
	CatFreeRoll
	NumCategories // sentinel: always last
)

var categoryNames = [NumCategories]string{
	CatJumbleberry:   "Jumbleberry",
	CatSugarberry:    "Sugarberry",
	CatPickleberry:   "Pickleberry",
	CatMoonberry:     "Moonberry",
	CatBasketOfThree: "Basket of Three",
	CatBasketOfFour:  "Basket of Four",
	CatBasketOfFive:  "Basket of Five",
	CatMixedBasket:   "Mixed Basket",
	CatFreeRoll:      "Free Roll",
}

func (c Category) String() string {
	if c < NumCategories {
		return categoryNames[c]
	}
	return fmt.Sprintf("Category(%d)", c)
}

// CategorySet is a bitmask representing a set of categories.
// Bit i is set if category i is in the set.
type CategorySet uint16

const AllCategories CategorySet = (1 << NumCategories) - 1

// Has returns true if the category is in the set.
func (cs CategorySet) Has(c Category) bool {
	return cs&(1<<c) != 0
}

// Remove returns a new set with the given category removed.
func (cs CategorySet) Remove(c Category) CategorySet {
	return cs &^ (1 << c)
}

// Add returns a new set with the given category added.
func (cs CategorySet) Add(c Category) CategorySet {
	return cs | (1 << c)
}

// Count returns the number of categories in the set.
func (cs CategorySet) Count() int {
	// popcount for uint16
	n := 0
	bits := cs
	for bits != 0 {
		n++
		bits &= bits - 1
	}
	return n
}

// ForEach calls fn for each category in the set.
func (cs CategorySet) ForEach(fn func(Category)) {
	for c := Category(0); c < NumCategories; c++ {
		if cs.Has(c) {
			fn(c)
		}
	}
}
