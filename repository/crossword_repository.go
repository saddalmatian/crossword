package repository

import (
	"crossword/crossword" // Assuming this is your package path
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Interface for Crossword data access
type CrosswordRepository interface {
	GetWordsByCategory(category string) ([]*crossword.Word, error)
	SetGrid(grid *crossword.Grid)
	GetGrid() *crossword.Grid
	// Add more methods as needed...
}

// In-Memory Implementation
type InMemoryCrosswordRepository struct {
	grid       *crossword.Grid
	categories map[string][]*crossword.Word
}

// NewInMemoryCrosswordRepository creates and initializes a new
// in-memory repository.
func NewInMemoryCrosswordRepository(databaseDir string) (*InMemoryCrosswordRepository, error) {
	repo := &InMemoryCrosswordRepository{
		grid:       nil, // Or initialize with a default grid if needed
		categories: make(map[string][]*crossword.Word),
	}

	// Load data from the database directory
	err := repo.LoadAllDatabases(databaseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load databases: %w", err)
	}

	return repo, nil
}

// GetWordsByCategory retrieves words for a specific category.
func (repo *InMemoryCrosswordRepository) GetWordsByCategory(category string) ([]*crossword.Word, error) {
	words, ok := repo.categories[category]
	if !ok {
		return nil, fmt.Errorf("category not found: %s", category)
	}
	return words, nil
}

func (repo *InMemoryCrosswordRepository) SetGrid(g *crossword.Grid) {
	repo.grid = g
}
func (repo *InMemoryCrosswordRepository) GetGrid() *crossword.Grid {
	return repo.grid
}

// LoadAllDatabases reads all JSON files and organizes words by category.
func (repo *InMemoryCrosswordRepository) LoadAllDatabases(databaseDir string) error {
	err := filepath.Walk(databaseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(path) != ".json" {
			return nil // Skip directories and non-JSON files
		}

		file, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		data := make(map[string]string) // For word-clue pairs from JSON
		err = json.Unmarshal(file, &data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON from %s: %w", path, err)
		}

		// Get category (filename without extension)
		category := filepath.Base(path)
		category = category[:len(category)-len(filepath.Ext(category))]

		// Add words to the correct category
		for wordText, clue := range data {
			word := crossword.NewWord(wordText, clue)
			repo.categories[category] = append(repo.categories[category], word)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to read database files: %w", err)
	}

	return nil
}
