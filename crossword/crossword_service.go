package crossword

import (
	"sort"
	"sync"
)

// CrosswordService (You can remove this if it's not needed)
type CrosswordService struct {
}

// NewCrosswordService (You can remove this if it's not needed)
func NewCrosswordService() *CrosswordService {
	return &CrosswordService{}
}

// GenerateCrossword generates a masked crossword grid from a list of words.
func (cs *CrosswordService) GenerateCrossword(words []*Word) (*GridResponse, *Grid) {
	// 1. Sort words by length
	sort.Sort(ByLen(words))

	// 2. Create grid
	grid := NewGrid(len(words[0].Word)*3, len(words[0].Word)*3)

	// 3. Place words (your existing logic)
	// Place the longest word at the center
	words[0].SetPosition(grid.Limit()/2-len(words[0].Word)/2, grid.Limit()/2, Horizontal)
	// Setting the values for the first word
	words[0].Number = 1
	words[0].StartX = words[0].Position.x
	words[0].StartY = words[0].Position.y

	grid.SetWord(*words[0])

	// Keep track of successfully placed words
	placedWords := []*Word{words[0]}

	var wg sync.WaitGroup
	wg.Add(len(words) - 1)

	for i := 1; i < len(words); i++ {
		go func(i int) {
			defer wg.Done()
			grid.SmartPosition(words[i], &placedWords, i+1)
		}(i)
	}

	wg.Wait()

	// 4. Reduce grid size
	newGrid := grid.ReduceGridSize()

	// 5. Mask the grid
	maskedGrid := newGrid.MaskGrid()
	// 5. Generate Character Positions from maskedGrid
	characterPositions := []*CharacterPosition{}
	for x, row := range maskedGrid {
		for y, cell := range row {
			if cell != 0 { // Assuming 0 represents an empty cell
				characterPositions = append(characterPositions, &CharacterPosition{
					X: x,
					Y: y,
				})
			}
		}
	}

	// Prepare the response data
	return &GridResponse{
		GridXDim:           len(maskedGrid),
		GridYDim:           len(maskedGrid[0]), // Assuming all rows have the same length
		CharacterPositions: characterPositions,
		Words:              placedWords,
		// Add clues or other relevant data here
	}, newGrid
}
