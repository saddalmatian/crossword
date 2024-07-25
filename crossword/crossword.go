package crossword

import (
	"encoding/json"
	"fmt"
	"os"
)

type CharacterPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Direction represents the orientation of a word in the crossword puzzle.
type Direction int

// Prepare the response data
type GridResponse struct {
	GridXDim           int                  `json:"grid_x_dim"` // Horizontal size
	GridYDim           int                  `json:"grid_y_dim"` // Vertical size
	CharacterPositions []*CharacterPosition `json:"characterPositions"`
	Words              []*Word              `json:"words"`
	// Add clues or other relevant data here
}

const (
	// Horizontal direction.
	Horizontal Direction = iota
	// Vertical direction.
	Vertical
)

// String returns the string representation of the Direction.
func (d Direction) String() string {
	return [...]string{"Horizontal", "Vertical"}[d]
}

// Position represents the starting position of a word in the crossword puzzle.
type Position struct {
	x int
	y int
}

type Word struct {
	Word       string    `json:"-"`         // The answer word
	Position   Position  `json:"-"`         // Starting position (to be filled later)
	WordLength int       `json:"-"`         // Length of the word
	Clue       string    `json:"clue"`      // The question for the crossword puzzle
	Direction  Direction `json:"direction"` // Direction (to be filled later)
	Number     int       `json:"number"`    // Question number (add this field)
	StartX     int       `json:"startX"`    // Starting X-coordinate
	StartY     int       `json:"startY"`    // Starting Y-coordinate
}

// ByLen implements sort.Interface for []Word based on the word length.
type ByLen []*Word

func (a ByLen) Len() int           { return len(a) }
func (a ByLen) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByLen) Less(i, j int) bool { return len(a[i].Word) > len(a[j].Word) }

// NewWord creates a new Word instance with calculated length.
func NewWord(text, clue string) *Word {
	return &Word{
		Word:       text,
		Clue:       clue,
		WordLength: len(text),
	}
}

// SetPosition sets the position and direction of a word.
func (w *Word) SetPosition(x, y int, direction Direction) {
	w.Position.x = x
	w.Position.y = y
	w.Direction = direction
}

func (w *Word) ChangeDirection() {
	if w.Direction == Horizontal {
		w.Direction = Vertical
	} else {
		w.Direction = Horizontal
	}
}

// GenerateMapPosition generates a map of letter positions within the word.
func (w *Word) GenerateMapPosition() map[string][]int {
	positionMap := make(map[string][]int)
	for i, char := range w.Word {
		letter := string(char)
		positionMap[letter] = append(positionMap[letter], i)
	}
	return positionMap
}

// Grid represents the crossword puzzle grid.
type Grid struct {
	grid       [][]rune
	gridLimitX int
	gridLimitY int
}

// NewGrid creates a new crossword grid with the specified dimensions.
func NewGrid(limitX, limitY int) *Grid {
	grid := make([][]rune, limitX)
	for i := range grid {
		grid[i] = make([]rune, limitY)
	}
	return &Grid{
		grid:       grid,
		gridLimitX: limitX,
		gridLimitY: limitY,
	}
}

// ReduceGridSize reduces the size of the crossword grid to fit the actual content.
func (g *Grid) ReduceGridSize() *Grid {
	minX, minY, maxX, maxY := g.findGridBounds()
	newGridHeight := maxX - minX + 1
	newGridWidth := maxY - minY + 1
	newGrid := NewGrid(newGridHeight, newGridWidth)

	for x := 0; x < newGridHeight; x++ {
		for y := 0; y < newGridWidth; y++ {
			newGrid.grid[x][y] = g.grid[x+minX][y+minY]
		}
	}

	return newGrid
}

func (g *Grid) Grid() [][]rune {
	return g.grid
}

// findGridBounds finds the minimum and maximum coordinates where letters exist in the grid.
func (g *Grid) findGridBounds() (int, int, int, int) {
	minX, minY, maxX, maxY := g.Limit(), g.Limit(), 0, 0

	for x := 0; x < g.Limit(); x++ {
		for y := 0; y < g.Limit(); y++ {
			if g.grid[x][y] != 0 {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	return minX, minY, maxX, maxY
}

// Limit return the limit of the grid
func (g *Grid) Limit() int {
	return g.gridLimitX
}

// Print the crossword grid to the console.
func (g *Grid) Print() {
	for _, row := range g.grid {
		for _, cell := range row {
			if cell == 0 {
				fmt.Print(". ") // Print '.' for empty cells
			} else {
				fmt.Printf("%c ", cell)
			}
		}
		fmt.Println()
	}
}

// SetWord places a word on the grid at its specified position and direction.
func (g *Grid) SetWord(w Word) {
	x, y := w.Position.x, w.Position.y
	for i := 0; i < w.WordLength; i++ {
		g.grid[x][y] = rune(w.Word[i])
		switch w.Direction {
		case Horizontal:
			y++
		case Vertical:
			x++
		}
	}
}

// SmartPosition intelligently positions a word in the crossword grid, avoiding collisions.
func (g *Grid) SmartPosition(w *Word, placedWords *[]*Word, number int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Err: ", r)
			fmt.Println("Continue to create Grid")
		}
	}()
	wordMapPosition := w.GenerateMapPosition()

	// Loop infinitely over the placedWords to check the possibility of the current word with the placedWords
	for {
		// Loop over the placedWords, the placeWords will be updating continuously
		for _, placedWord := range *placedWords {
			// If current placedWord direction is Vertical, change the current word direction to other, or the opposite
			if w.Direction == placedWord.Direction {
				w.ChangeDirection()
			}

			// Range over the placeWord word
			for placedWordLetterIdx, placedWordLetter := range placedWord.Word {
				// if the letter or current placedWord exist in the wordMapPosition above
				if wordPositions, ok := wordMapPosition[string(placedWordLetter)]; ok {
					// Check position of the current letter in current word map
					for _, positionIdx := range wordPositions {
						// if the word is Vertically direction, the current word position x will be the placed word position x - current letter in current word map
						// the current word position y will be the placed word position y + current position of checking placeword
						if w.Direction == Vertical {
							w.Position.x = placedWord.Position.x - positionIdx
							w.Position.y = placedWord.Position.y + placedWordLetterIdx
						} else { // Horizontal
							// in the opposite of the above
							w.Position.x = placedWord.Position.x + placedWordLetterIdx
							w.Position.y = placedWord.Position.y - positionIdx
						}

						// Check new position boundaries, make sure it will not exceed the grid limit
						if w.Position.x < 0 || w.Position.x >= g.Limit()-w.WordLength || w.Position.y < 0 || w.Position.y >= g.Limit()-w.WordLength {
							continue
						}

						// Check the possibility of the new word
						if g.CheckPossible(w) {
							g.SetWord(*w)
							*placedWords = append(*placedWords, w)
							// Set the placement details directly in the Word struct
							w.Number = number
							w.StartX = w.Position.x
							w.StartY = w.Position.y
							return
						}
					}
				}
			}
		}
	}
}

// CheckPossible checks if it's possible to place a word in the crossword grid without collisions.
func (g *Grid) CheckPossible(word *Word) bool {
	x, y := word.Position.x, word.Position.y

	// Check top and bottom boundaries for vertical words
	if word.Direction == Vertical && (x > 0 && (g.grid[x-1][y] != 0 || g.grid[x+word.WordLength][y] != 0)) {
		return false
	}

	// Check left and right boundaries for horizontal words
	if word.Direction == Horizontal && (y > 0 && (g.grid[x][y-1] != 0 || g.grid[x][y+word.WordLength] != 0)) {
		return false
	}

	for idx := 0; idx < word.WordLength; idx++ {
		// Check for collisions with filled cells and adjacent cells
		if g.grid[x][y] == 0 && ((word.Direction == Vertical && (y > 0 && g.grid[x][y-1] != 0 || y < g.gridLimitX-1 && g.grid[x][y+1] != 0)) ||
			(word.Direction == Horizontal && (x > 0 && g.grid[x-1][y] != 0 || x < g.gridLimitX-1 && g.grid[x+1][y] != 0))) {
			return false
		}

		// Check for collisions with existing letters in the word
		if g.grid[x][y] != 0 && g.grid[x][y] != rune(word.Word[idx]) {
			return false
		}

		// Move to the next position based on the word's direction
		switch word.Direction {
		case Horizontal:
			y++
		case Vertical:
			x++
		}
	}

	return true
}

// ReadWordsFromJSON reads crossword data from a JSON file.
func ReadWordsFromJSON(filename string) ([]*Word, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	data := make(map[string]string)
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}

	var words []*Word
	for answer, clue := range data {
		words = append(words, NewWord(
			answer,
			clue,
		))
	}

	return words, nil
}

// Helper function to mask the grid
func (g *Grid) MaskGrid() [][]rune {
	maskedGrid := make([][]rune, len(g.grid))
	for i := range g.grid {
		maskedGrid[i] = make([]rune, len(g.grid[i]))
		for j := range g.grid[i] {
			if g.grid[i][j] != 0 { // Check if the cell is not empty
				maskedGrid[i][j] = '#' // Replace with your desired placeholder
			}
		}
	}
	return maskedGrid
}
