package sqlimporter

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func createDSN() string {
	return ""
}

func Connect(driver, dsn string) (*sqlx.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	db, err := sqlx.ConnectContext(ctx, driver, dsn)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}

const DBNameDefault = "SQL_IMPORTER_DB_"

// CreateDB used to create database
// and import all queries located in a directories
func CreateDB(driver, dsn string) (*sqlx.DB, func() error, error) {
	defaultDrop := func() error {
		return nil
	}
	db, err := Connect(driver, dsn)
	if err != nil {
		return nil, defaultDrop, err
	}

	// create a new database
	// database name is always a random name
	unix := time.Now().Unix()
	randSource := rand.NewSource(unix)
	r := rand.New(randSource)
	dbName := DBNameDefault + strconv.Itoa(r.Int())
	// TODO: separate this, this is a dialect and might be not the same with other db
	createDBQuery := fmt.Sprintf(getDialect(driver, "create"), dbName)
	// exec create new b
	_, err = db.Exec(createDBQuery)
	if err != nil {
		return nil, defaultDrop, err
	}

	// use new db
	useDatabaseQuery := fmt.Sprintf(getDialect(driver, "use"), dbName)
	_, err = db.Exec(useDatabaseQuery)
	if err != nil {
		return nil, defaultDrop, err
	}
	return db, func() error {
		deleteDatabaseQuery := fmt.Sprintf(getDialect(driver, "drop"), dbName)
		_, err := db.Exec(deleteDatabaseQuery)
		if err != nil {
			return err
		}
		return db.Close()
	}, nil
}

// ImportSchemaFromFiles
func ImportSchemaFromFiles(db *sqlx.DB, filepath string) error {
	files, err := getFileList(filepath)
	if err != nil {
		return err
	}

	// an sql file will be executed as one batch of transaction
	for _, file := range files {
		sqlContents, err := parseFiles(file)
		if err != nil {
			return err
		}
		// end if empty
		if len(sqlContents) == 0 {
			return nil
		}

		tx, err := db.BeginTx(context.TODO(), nil)
		if err != nil {
			return err
		}

		var query string
		for key := range sqlContents {
			query = sqlContents[key]
			_, err = tx.ExecContext(context.TODO(), query)
			if err != nil {
				break
			}
		}

		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				return fmt.Errorf("Failed to rollback from file %s with error %s", file, errRollback.Error())
			}
			return fmt.Errorf("Failed to execute file %s with error %s and query: \n %s", file, err.Error(), query)
		}
		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("Failed to commit from file %s with error %s", file, err.Error())
		}

	}
	return nil
}
