package game

import (
	"testing"
)

func TestCategoryString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		category Category
		want     string
	}{
		{"Jumbleberry", CatJumbleberry, "Jumbleberry"},
		{"Sugarberry", CatSugarberry, "Sugarberry"},
		{"Pickleberry", CatPickleberry, "Pickleberry"},
		{"Moonberry", CatMoonberry, "Moonberry"},
		{"BasketOfThree", CatBasketOfThree, "Basket of Three"},
		{"BasketOfFour", CatBasketOfFour, "Basket of Four"},
		{"BasketOfFive", CatBasketOfFive, "Basket of Five"},
		{"MixedBasket", CatMixedBasket, "Mixed Basket"},
		{"FreeRoll", CatFreeRoll, "Free Roll"},
		{"Invalid", Category(99), "Category(99)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.category.String(); got != tt.want {
				t.Errorf("Category.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategorySetHas(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		set      CategorySet
		category Category
		want     bool
	}{
		{
			name:     "empty set",
			set:      0,
			category: CatJumbleberry,
			want:     false,
		},
		{
			name:     "single category present",
			set:      1 << CatJumbleberry,
			category: CatJumbleberry,
			want:     true,
		},
		{
			name:     "single category absent",
			set:      1 << CatJumbleberry,
			category: CatSugarberry,
			want:     false,
		},
		{
			name:     "all categories",
			set:      AllCategories,
			category: CatMoonberry,
			want:     true,
		},
		{
			name:     "multiple categories",
			set:      (1 << CatJumbleberry) | (1 << CatMoonberry),
			category: CatMoonberry,
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.set.Has(tt.category); got != tt.want {
				t.Errorf("CategorySet.Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategorySetRemove(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		set      CategorySet
		category Category
		wantHas  bool
	}{
		{
			name:     "remove from all",
			set:      AllCategories,
			category: CatJumbleberry,
			wantHas:  false,
		},
		{
			name:     "remove from single",
			set:      1 << CatJumbleberry,
			category: CatJumbleberry,
			wantHas:  false,
		},
		{
			name:     "remove non-existent",
			set:      1 << CatJumbleberry,
			category: CatSugarberry,
			wantHas:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.set.Remove(tt.category)
			if got := result.Has(tt.category); got != tt.wantHas {
				t.Errorf("CategorySet.Remove().Has() = %v, want %v", got, tt.wantHas)
			}
		})
	}
}

func TestCategorySetAdd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		set      CategorySet
		category Category
	}{
		{
			name:     "add to empty",
			set:      0,
			category: CatJumbleberry,
		},
		{
			name:     "add to existing",
			set:      1 << CatSugarberry,
			category: CatJumbleberry,
		},
		{
			name:     "add already present",
			set:      1 << CatJumbleberry,
			category: CatJumbleberry,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.set.Add(tt.category)
			if !result.Has(tt.category) {
				t.Errorf("CategorySet.Add() did not add category %v", tt.category)
			}
		})
	}
}

func TestCategorySetCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		set  CategorySet
		want int
	}{
		{
			name: "empty",
			set:  0,
			want: 0,
		},
		{
			name: "single",
			set:  1 << CatJumbleberry,
			want: 1,
		},
		{
			name: "all categories",
			set:  AllCategories,
			want: int(NumCategories),
		},
		{
			name: "three categories",
			set:  (1 << CatJumbleberry) | (1 << CatSugarberry) | (1 << CatMoonberry),
			want: 3,
		},
		{
			name: "all but one",
			set:  AllCategories &^ (1 << CatFreeRoll),
			want: int(NumCategories) - 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.set.Count(); got != tt.want {
				t.Errorf("CategorySet.Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategorySetForEach(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		set  CategorySet
		want []Category
	}{
		{
			name: "empty",
			set:  0,
			want: []Category{},
		},
		{
			name: "single",
			set:  1 << CatJumbleberry,
			want: []Category{CatJumbleberry},
		},
		{
			name: "three categories",
			set:  (1 << CatJumbleberry) | (1 << CatPickleberry) | (1 << CatFreeRoll),
			want: []Category{CatJumbleberry, CatPickleberry, CatFreeRoll},
		},
		{
			name: "all categories",
			set:  AllCategories,
			want: []Category{
				CatJumbleberry, CatSugarberry, CatPickleberry, CatMoonberry,
				CatBasketOfThree, CatBasketOfFour, CatBasketOfFive,
				CatMixedBasket, CatFreeRoll,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var got []Category
			tt.set.ForEach(func(c Category) {
				got = append(got, c)
			})
			if len(got) != len(tt.want) {
				t.Errorf("CategorySet.ForEach() visited %d categories, want %d", len(got), len(tt.want))
			}
			for i, c := range tt.want {
				if i >= len(got) || got[i] != c {
					t.Errorf("CategorySet.ForEach() category %d = %v, want %v", i, got[i], c)
				}
			}
		})
	}
}

func TestAllCategories(t *testing.T) {
	t.Parallel()

	// AllCategories should have all 9 categories
	if count := AllCategories.Count(); count != int(NumCategories) {
		t.Errorf("AllCategories.Count() = %d, want %d", count, NumCategories)
	}

	// Each category should be present
	for cat := Category(0); cat < NumCategories; cat++ {
		if !AllCategories.Has(cat) {
			t.Errorf("AllCategories.Has(%v) = false, want true", cat)
		}
	}
}
