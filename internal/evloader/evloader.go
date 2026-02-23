// Package evloader handles loading and computing the expected value table.
// It attempts to load from disk first, falling back to computing from scratch
// if the file doesn't exist.
package evloader

import (
	"fmt"
	"time"

	"github.com/iadams749/JBFieldsSolver/internal/ev"
)

// Load attempts to load the EV table from path. If the file doesn't exist
// or is invalid, it computes the table from scratch, saves it, and returns it.
func Load(path string) (*ev.Table, error) {
	table, err := ev.LoadJSON(path)
	if err == nil {
		fmt.Println("EV table loaded from", path)
		return table, nil
	}

	fmt.Println("EV table not found, computing...")
	start := time.Now()
	table = ev.Compute(func(size, total int) {
		elapsed := time.Since(start)
		fmt.Printf("  Completed size %d/%d  (%v elapsed)\n", size, total, elapsed.Round(time.Millisecond))
	})

	if err := table.SaveJSON(path); err != nil {
		return nil, fmt.Errorf("error saving EV table: %w", err)
	}
	fmt.Printf("EV table computed and saved to %s\n", path)
	return table, nil
}
