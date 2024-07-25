package controller

import (
	"crossword/crossword"
	"crossword/repository"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CrosswordController struct {
	repo repository.CrosswordRepository
}

func NewCrosswordController(repo repository.CrosswordRepository) *CrosswordController {
	return &CrosswordController{repo: repo}
}

func (ctrl *CrosswordController) GenerateCrossword(c echo.Context) error {
	// 1. Get the chosen category from the request parameters
	category := c.QueryParam("category")

	// 2. Fetch words for the chosen category from the repository
	words, err := ctrl.repo.GetWordsByCategory(category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to get words for category '%s'", category)})
	}
	// 2. Create the crossword service
	crosswordService := crossword.NewCrosswordService()

	// 3. Delegate crossword generation to the service
	maskedGrid, grid := crosswordService.GenerateCrossword(words)

	// 4. Store grid in repo
	ctrl.repo.SetGrid(grid)

	return c.JSON(http.StatusOK, maskedGrid)
}
