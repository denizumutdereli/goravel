package models

import (
	"database/sql"
	"fmt"
	"os"

	db2 "github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
	"github.com/upper/db/v4/adapter/postgresql"
)

var db sql.DB
var upper db2.Session

type Models struct {
}

func New(databasePool *sql.DB) Models {

	db = *databasePool

	if os.Getenv("DATABASE_TYPE") == "mysql" || os.Getenv("DATABASE_TYPE") == "mariadb" {
		upper, _ = mysql.New(databasePool)
	} else {
		upper, _ = postgresql.New(databasePool)
	}

	return Models{}

}

func getInsertId(i db2.ID) int {
	idType := fmt.Sprintf("%T", i)

	if idType == "int64" {
		return int(i.(int64))
	}

	return i.(int)
}