package main

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/canonical/go-dqlite/driver"
	"github.com/pkg/errors"
)

func dbExec(dbname string, cluster []string, statements ...string) error {
	store := getStore(cluster)
	driver, err := driver.New(store, driver.WithLogFunc(logFunc))
	if err != nil {
		return errors.Wrapf(err, "failed to create dqlite driver")
	}
	sql.Register("dqlite", driver)

	db, err := sql.Open("dqlite", dbname)
	if err != nil {
		return errors.Wrap(err, "can't open database")
	}
	defer db.Close()

	for i, statement := range statements {
		action := strings.ToUpper(strings.Fields(statement)[0])
		if action == "SELECT" {
			rows, err := db.Query(statement)
			if err != nil {
				return errors.Wrap(err, "query failed")
			}
			defer rows.Close()
			columns, _ := rows.Columns()
			fmt.Println(columns)
			//buffer := make([]interface{}, len(columns))
			buffer := make([]string, len(columns))
			scanTo := make([]interface{}, len(columns))
			for i := range buffer {
				scanTo[i] = &buffer[i]
			}
			for rows.Next() {
				if err := rows.Scan(scanTo...); err != nil {
					return errors.Wrap(err, "failed to scan row")
				}
				fmt.Println(buffer)
			}
			continue
		}
		if _, err := db.Exec(statement); err != nil {
			return errors.Wrapf(err, "dbExec fail %d/%d", i+1, len(statements))
		}

	}
	return nil
}

func dbQuery(key string, cluster []string) error {
	store := getStore(cluster)
	driver, err := driver.New(store, driver.WithLogFunc(logFunc))
	if err != nil {
		return errors.Wrapf(err, "failed to create dqlite driver")
	}
	sql.Register("dqlite", driver)

	db, err := sql.Open("dqlite", "demo.db")
	if err != nil {
		return errors.Wrap(err, "can't open demo database")
	}
	defer db.Close()

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS model (key TEXT, value TEXT)"); err != nil {
		return errors.Wrap(err, "can't create demo table")
	}

	row := db.QueryRow("SELECT value FROM model WHERE key = ?", key)
	value := ""
	if err := row.Scan(&value); err != nil {
		return errors.Wrap(err, "failed to get key")
	}
	fmt.Println(value)

	return nil
}

func dbUpdate(key, value string, cluster []string) error {
	store := getStore(cluster)
	driver, err := driver.New(store, driver.WithLogFunc(logFunc))
	if err != nil {
		return errors.Wrapf(err, "failed to create dqlite driver")
	}
	sql.Register("dqlite", driver)

	db, err := sql.Open("dqlite", "demo.db")
	if err != nil {
		return errors.Wrap(err, "can't open demo database")
	}
	defer db.Close()

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS model (key TEXT, value TEXT)"); err != nil {
		return errors.Wrap(err, "can't create demo table")
	}

	if _, err := db.Exec("INSERT OR REPLACE INTO model(key, value) VALUES(?, ?)", key, value); err != nil {
		return errors.Wrap(err, "can't update key")
	}

	fmt.Println("done")

	return nil
}
