package main

import (
	"crossword/controller"
	"crossword/repository"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {

	// Choose your repository implementation
	var repo, err = repository.NewInMemoryCrosswordRepository("database/")
	if err != nil {
		log.Fatal(err)
	}
	// Inject dependency
	ctrl := controller.NewCrosswordController(repo)

	e := echo.New()

	e.POST("/generate", ctrl.GenerateCrossword)

	e.Start(":8000")

	// ... (Start server)
}
