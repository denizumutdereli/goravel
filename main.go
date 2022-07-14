package main

import (
	"goravel/handlers"

	"github.com/denizumutdereli/gortisan"
)

type application struct {
	App      *gortisan.Gortisan
	Handlers *handlers.Handlers
}

func main() {
	g := initApplication()
	g.App.ListenAndServe()
}
