package main

import (
	"os"

	"github.com/joho/godotenv"
)

func setup() {
	err := godotenv.Load()
	if err != nil {
		gracefulExit(err)
	}

	path, err := os.Getwd()
	if err != nil {
		gracefulExit(err)
	}

	gor.RootPath = path
	gor.DB.DatabaseType = os.Getenv("DATABASE_TYPE")
}
