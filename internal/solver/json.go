// This file provides JSON serialization for solver recommendations,
// suitable for API responses and structured output.
package solver

// RecommendationJSON is the JSON-friendly representation of a Recommendation.
type RecommendationJSON struct {
	BestAction       ActionJSON        `json:"best_action"`
	TheoreticalMax   float64           `json:"theoretical_max"`
	CategoryOptions  []CategoryOptJSON `json:"category_options"`
	TopRerollOptions []RerollOptJSON   `json:"top_reroll_options"`
}

// ActionJSON is the JSON-friendly representation of an Action.
type ActionJSON struct {
	Type     string  `json:"type"`     // "score" or "reroll"
	Keep     string  `json:"keep"`     // e.g. "1M" or "" for score
	Category string  `json:"category"` // e.g. "Jumbleberry" or "" for reroll
	EV       float64 `json:"ev"`
}

// CategoryOptJSON is the JSON-friendly representation of a CategoryOption.
type CategoryOptJSON struct {
	Category       string  `json:"category"`
	ImmediateScore int     `json:"immediate_score"`
	FutureEV       float64 `json:"future_ev"`
	TotalValue     float64 `json:"total_value"`
}

// RerollOptJSON is the JSON-friendly representation of a RerollOption.
type RerollOptJSON struct {
	Keep        string  `json:"keep"`
	NumRerolled int     `json:"num_rerolled"`
	EV          float64 `json:"ev"`
}

// RecommendationToJSON converts a Recommendation to its JSON-friendly form.
func RecommendationToJSON(rec Recommendation) RecommendationJSON {
	var actionType string
	var keep, category string

	switch rec.BestAction.Type {
	case ScoreAction:
		actionType = "score"
		category = rec.BestAction.Category.String()
	case RerollAction:
		actionType = "reroll"
		keep = FormatKeep(rec.BestAction.Keep)
	}

	var catOpts []CategoryOptJSON
	for _, opt := range rec.CategoryOptions {
		catOpts = append(catOpts, CategoryOptJSON{
			Category:       opt.Category.String(),
			ImmediateScore: opt.ImmediateScore,
			FutureEV:       opt.FutureEV,
			TotalValue:     opt.TotalValue,
		})
	}

	var rerollOpts []RerollOptJSON
	for _, opt := range rec.TopRerollOptions {
		rerollOpts = append(rerollOpts, RerollOptJSON{
			Keep:        FormatKeep(opt.Keep),
			NumRerolled: opt.NumRerolled,
			EV:          opt.EV,
		})
	}

	// Ensure empty slices serialize as [] not null
	if catOpts == nil {
		catOpts = []CategoryOptJSON{}
	}
	if rerollOpts == nil {
		rerollOpts = []RerollOptJSON{}
	}

	return RecommendationJSON{
		BestAction: ActionJSON{
			Type:     actionType,
			Keep:     keep,
			Category: category,
			EV:       rec.BestAction.EV,
		},
		TheoreticalMax:   rec.TheoreticalMax,
		CategoryOptions:  catOpts,
		TopRerollOptions: rerollOpts,
	}
}
