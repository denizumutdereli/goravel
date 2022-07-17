package main

import (
	"errors"
	"fmt"
	"time"
)

func doMake(arg2, arg3 string) error {

	switch arg2 {
	case "migration":
		dbType := gor.DB.DatabaseType
		if arg3 == "" {
			gracefulExit(errors.New("you must give the migration a name"))
		}

		fileName := fmt.Sprintf("%d_%s", time.Now().UnixMicro(), arg3)

		upFile := gor.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
		downFile := gor.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

		err := copyFilefromTemplate("templates/migrations/migration."+dbType+".up.sql", upFile)
		if err != nil {
			gracefulExit(err)
		}

		err = copyFilefromTemplate("templates/migrations/migration."+dbType+".down.sql", downFile)
		if err != nil {
			gracefulExit(err)
		}
	}

	return nil
}
