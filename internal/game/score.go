package game

// Score returns the score for the given dice in the given category.
// Returns 0 if the dice don't meet the category's requirement.
func Score(d Dice, cat Category) int {
	switch cat {
	case CatJumbleberry:
		return int(d[Jumbleberry]) * BerryPoints[Jumbleberry]
	case CatSugarberry:
		return int(d[Sugarberry]) * BerryPoints[Sugarberry]
	case CatPickleberry:
		return int(d[Pickleberry]) * BerryPoints[Pickleberry]
	case CatMoonberry:
		return int(d[Moonberry]) * BerryPoints[Moonberry]
	case CatBasketOfThree:
		if hasNOfAKind(d, 3) {
			return d.Points()
		}
		return 0
	case CatBasketOfFour:
		if hasNOfAKind(d, 4) {
			return d.Points()
		}
		return 0
	case CatBasketOfFive:
		if hasNOfAKind(d, 5) {
			return d.Points()
		}
		return 0
	case CatMixedBasket:
		if d[Jumbleberry] >= 1 && d[Sugarberry] >= 1 &&
			d[Pickleberry] >= 1 && d[Moonberry] >= 1 {
			return d.Points()
		}
		return 0
	case CatFreeRoll:
		return d.Points()
	default:
		panic("invalid category")
	}
}

// hasNOfAKind returns true if any type (including Pest) has count >= n.
func hasNOfAKind(d Dice, n uint8) bool {
	for _, count := range d {
		if count >= n {
			return true
		}
	}
	return false
}
