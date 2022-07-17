package main

import (
	"goravel/handlers"
	"goravel/models"

	"github.com/denizumutdereli/gortisan"
)

type application struct {
	App      *gortisan.Gortisan
	Handlers *handlers.Handlers
	Models   models.Models
}

func main() {
	g := initApplication()
	g.App.ListenAndServe()
}
