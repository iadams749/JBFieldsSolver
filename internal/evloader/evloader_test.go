package evloader

import (
	"os"
	"testing"
)

func TestLoadExisting(t *testing.T) {
	// This test requires a valid ev_table.json file to exist
	// Skip if not running in workspace root
	t.Parallel()

	const testPath = "../../ev_table.json"
	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		t.Skip("ev_table.json not found, skipping test")
	}

	table, err := Load(testPath)
	if err != nil {
		t.Fatalf("Load(%q) error = %v", testPath, err)
	}

	if table == nil {
		t.Fatal("Load() returned nil table")
	}
}

func TestLoadComputesIfMissing(t *testing.T) {
	t.Parallel()

	// Use a path that doesn't exist
	tmpFile := t.TempDir() + "/nonexistent_ev_table.json"

	table, err := Load(tmpFile)
	if err != nil {
		t.Fatalf("Load(%q) error = %v", tmpFile, err)
	}

	if table == nil {
		t.Fatal("Load() returned nil table")
	}

	// Verify the file was created
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Errorf("Load() did not create file at %q", tmpFile)
	}

	// Verify table has reasonable values
	allCatsEV := table.EV(511) // all 9 categories
	if allCatsEV < 100.0 || allCatsEV > 150.0 {
		t.Errorf("Computed table.EV(all categories) = %v, expected ~121.8", allCatsEV)
	}
}

func TestLoadInvalidFile(t *testing.T) {
	t.Parallel()

	// Create a file with invalid JSON
	tmpFile := t.TempDir() + "/invalid_ev_table.json"
	if err := os.WriteFile(tmpFile, []byte("invalid json content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Load should compute a new table since the file is invalid
	table, err := Load(tmpFile)
	if err != nil {
		t.Fatalf("Load(%q) error = %v", tmpFile, err)
	}

	if table == nil {
		t.Fatal("Load() returned nil table")
	}
}
