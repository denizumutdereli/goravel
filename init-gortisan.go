package main

import (
	"goravel/handlers"
	"log"
	"os"

	"github.com/denizumutdereli/gortisan"
)

func initApplication() *application {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	//init gortisan
	gor := &gortisan.Gortisan{}

	err = gor.New(path)

	if err != nil {
		log.Fatal(err)
	}

	gor.Appname = "goravel"

	myHandlers := &handlers.Handlers{
		App: gor,
	}

	//gor.InfoLog.Println("Debug is set to", gor.Debug)

	app := &application{
		App:      gor,
		Handlers: myHandlers,
	}

	app.App.Routes = app.routes()

	return app
}
