# Jumbleberry Fields Solver

An optimal strategy solver for the dice game **Jumbleberry Fields**, using dynamic programming to compute expected values and recommend the best play at any game state.

## Live Solver

**[Try it in your browser →](https://iadams749.github.io/JBFieldsSolver/)**

The solver runs entirely client-side via WebAssembly — no server required. Enter your current dice, rolls remaining, and which categories you have left, and it instantly recommends the optimal play along with the expected game score.

## Game Overview

Jumbleberry Fields is a Yahtzee-style dice game with unique scoring rules:

### Dice
- **5 dice** with non-uniform probabilities:
  - **Jumbleberry** (J): 30% probability, 2 points each
  - **Sugarberry** (S): 30% probability, 2 points each
  - **Pickleberry** (P): 20% probability, 4 points each
  - **Moonberry** (M): 10% probability, 7 points each
  - **Pest** (X): 10% probability, 0 points

### Scoring Categories
Nine categories must be filled once each over nine rounds:
1. **Jumbleberry** - Score all Jumbleberries
2. **Sugarberry** - Score all Sugarberries
3. **Pickleberry** - Score all Pickleberries
4. **Moonberry** - Score all Moonberries
5. **Basket of Three** - 3+ of a kind scores all dice
6. **Basket of Four** - 4+ of a kind scores all dice
7. **Basket of Five** - 5 of a kind scores all dice
8. **Mixed Basket** - At least one of each berry (not Pest) scores all dice
9. **Free Roll** - Any combination scores all dice

### Gameplay
Each turn:
1. Roll all 5 dice
2. Up to 2 rerolls: choose which dice to keep/reroll
3. Score in one remaining category

## Features

- **Optimal Play Algorithm**: Uses precomputed expected values (EV) via dynamic programming
- **Interactive CLI**: Real-time solver recommendations during gameplay
- **HTTP API**: REST endpoint for programmatic access
- **Fast Lookups**: Precomputes all 512 category subset states and 126 dice outcomes
- **Comprehensive Output**: Shows best action, alternative options, and theoretical maximum

## Installation

Requires Go 1.26 or later:

```bash
# Build both executables
go build -o jbf-cli ./cmd/cli
go build -o jbf-api ./cmd/api
```

## Usage

### CLI (Interactive)

```bash
./jbf-cli
```

On first run, the solver computes the EV table (takes ~1-2 seconds), then saves it to `ev_table.json` for instant future loads.

**Example session:**
```
Dice: JJSPM
Rolls left (0-2): 2
Categories remaining: all

=== Solver Recommendation ===
Dice: {J:2 S:1 P:1 M:1 X:0}  |  Rolls left: 2  |  Categories: 9 remaining

Best action: REROLL
  Keep 1M  (reroll 4)
  Expected value: 94.23

Top reroll options:
  #1  Keep 1M                       reroll 4  EV:   94.23
  #2  Keep 1M 1P                    reroll 3  EV:   93.89
  #3  Keep nothing                  reroll 5  EV:   93.56
  ...
```

**Input formats:**

- **Dice:**
  - Sequence: `JJSPM` (exactly 5 letters)
  - Counts: `2J 1S 1P 1M` (space-separated)

- **Categories:**
  - `all` - All 9 categories
  - `all-j-s` - All except Jumbleberry and Sugarberry
  - `j,s,p,m,3k,4k,5k,mix,fr` - Comma-separated list
  - Shortcuts: `j/s/p/m` (berries), `3k/4k/5k` (baskets), `mix` (mixed), `fr` (free roll)

### API Server

```bash
# Start server on default port 8080
./jbf-api

# Or specify custom port and EV table path
./jbf-api -addr :3000 -ev ./my_ev_table.json
```

**API Endpoint:**

```http
POST /solve
Content-Type: application/json

{
  "dice": "JJSPM",
  "rolls_left": 2,
  "categories": "all-j"
}
```

**Response:**
```json
{
  "best_action": {
    "type": "reroll",
    "keep": "1M",
    "category": "",
    "ev": 94.23
  },
  "theoretical_max": 220.0,
  "category_options": [],
  "top_reroll_options": [
    {
      "keep": "1M",
      "num_rerolled": 4,
      "ev": 94.23
    }
  ]
}
```

## Architecture

### Package Structure

```
cmd/
  api/          HTTP API server
  cli/          Interactive command-line REPL
internal/
  ev/           Expected value table computation (core DP algorithm)
  evloader/     EV table loading/computation coordination
  game/         Game rules and types (dice, categories, scoring)
  solver/       Optimal decision algorithm and I/O formatting
```

### Algorithm

**Expected Value Computation** (computed once, cached):
1. Enumerate all 126 distinct 5-dice outcomes
2. For each of 512 category subsets (2^9 combinations):
   - Compute optimal value at 0 rolls left (choose category to score)
   - Compute optimal value at 1 roll left (choose dice to keep/reroll)
   - Compute optimal value at 2 rolls left (choose dice to keep/reroll)
   - Compute expected value over all initial roll outcomes
3. Store in lookup table: `EV[category_set]`

**Real-Time Solving** (instant lookup):
1. Build value layers for current category set (from precomputed table)
2. For player's specific dice, enumerate all keep decisions
3. Compute expected value using multinomial probabilities
4. Return best action (score or reroll) and alternatives

**Performance:**
- Precomputation: ~1-2 seconds (one-time)
- Query: <1ms per recommendation

## Implementation Details

### Memory Efficiency
- `GameState`: 12 bytes (value type, no heap allocation)
- `Dice`: 5 bytes (fixed-size array)
- `CategorySet`: 2 bytes (bitmask for 9 categories)
- Total EV table: 512 float64s = 4 KB

### Probability Model
Dice faces modeled with empirical probabilities:
- 3/10 each: Jumbleberry, Sugarberry
- 2/10: Pickleberry
- 1/10 each: Moonberry, Pest

Uses multinomial distribution for reroll outcome probabilities.

### Testing
Run tests:
```bash
go test ./...
```

Tests verify:
- Probability distributions sum to 1.0
- Dice indexing round-trips correctly
- All game rules implemented properly

## Expected Values

With optimal play from the start:
- **All 9 categories**: EV ≈ 121.8 points
- **Theoretical maximum** (perfect dice): 237 points

Individual category EVs vary by remaining set composition—strategic category ordering matters!

## License

This project is provided as-is for educational and entertainment purposes.
