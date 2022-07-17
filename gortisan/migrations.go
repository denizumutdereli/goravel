package gortisan

import (
	"log"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (g *Gortisan) MigrateUp(dsn string) error {

	m, err := migrate.New("file://"+g.RootPath, dsn)

	if err != nil {
		return err
	}

	defer m.Close()

	if err := m.Up(); err != nil {
		log.Println("Error running migration: ", err)
		return err
	}

	return nil
}

func (g *Gortisan) MigrateDown(dsn string) error {
	m, err := migrate.New("file://"+g.RootPath, dsn)

	if err != nil {
		return err
	}

	defer m.Close()

	if err := m.Down(); err != nil {
		return err
	}

	return nil
}

func (g *Gortisan) Steps(n int, dsn string) error {

	m, err := migrate.New("file://"+g.RootPath, dsn)

	if err != nil {
		return err
	}

	defer m.Close()

	if err := m.Steps(n); err != nil {
		return err
	}

	return nil

}

func (g *Gortisan) MigrateForce(dsn string) error {
	m, err := migrate.New("file://"+g.RootPath, dsn)

	if err != nil {
		return err
	}

	defer m.Close()

	if err := m.Force(-1); err != nil {
		return err
	}

	return nil
}
