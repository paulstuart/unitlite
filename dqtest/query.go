package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
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

type Queryor interface {
	QueryDB(queries ...string) ([][]string, error)
	Close() error
}

type dbx struct {
	db *sql.DB
}

// NewConnection creates a db connection
func NewConnection(dbname string, cluster []string) (*dbx, error) {
	store := getStore(cluster)
	driver, err := driver.New(store, driver.WithLogFunc(logFunc))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create dqlite driver")
	}
	sql.Register("dqlite", driver)

	db, err := sql.Open("dqlite", dbname)
	if err != nil {
		return nil, errors.Wrap(err, "can't open database")
	}

	return &dbx{db}, nil
}

func (d *dbx) Close() error {
	return d.db.Close()
}

func (d *dbx) QueryDB(queries ...string) ([][]string, error) {
	reply := make([][]string, 0, 32)
	for _, statement := range queries {
		action := strings.ToUpper(strings.Fields(statement)[0])
		if action != "SELECT" {
			return nil, fmt.Errorf("Invalid action: %q -- must use SELECT", action)
		}
		rows, err := d.db.Query(statement)
		if err != nil {
			return nil, errors.Wrap(err, "query failed")
		}
		defer rows.Close()
		columns, _ := rows.Columns()
		for rows.Next() {
			if len(reply) == 0 {
				reply = append(reply, columns)
			}
			buffer := make([]string, len(columns))
			scanTo := make([]interface{}, len(columns))
			for i := range buffer {
				scanTo[i] = &buffer[i]
			}
			if err := rows.Scan(scanTo...); err != nil {
				return nil, errors.Wrap(err, "failed to scan row")
			}
			reply = append(reply, buffer)
		}

	}
	return reply, nil
}

/*
	if _, err := db.Exec(statement); err != nil {
		return errors.Wrapf(err, "dbExec fail %d/%d", i+1, len(statements))
	}
*/

type Result struct {
	LastInsertID int64   `json:"last_insert_id,omitempty"`
	RowsAffected int64   `json:"rows_affected,omitempty"`
	Error        string  `json:"error,omitempty"`
	Time         float64 `json:"time,omitempty"`
}

// Rows represents the outcome of an operation that returns query data.
type Rows struct {
	Columns []string        `json:"columns,omitempty"`
	Types   []string        `json:"types,omitempty"`
	Values  [][]interface{} `json:"values,omitempty"`
	Error   string          `json:"error,omitempty"`
	Time    float64         `json:"time,omitempty"`
}

type ExecuteResponse struct {
	Results []Result
	Time    float64
	Raft    RaftResponse
}

type Executor interface {
	Execute(statements ...string) (*ExecuteResponse, error)
}

func (d *dbx) Execute(statements ...string) (*ExecuteResponse, error) {
	// TODO: add back atomic/timing bits

	results := make([]Result, 0, len(statements))

	for i, statement := range statements {
		resp, err := d.db.Exec(statement)
		if err != nil {
			return nil, errors.Wrapf(err, "dbExec fail %d/%d", i+1, len(statements))
		}
		lastID, _ := resp.LastInsertId()
		affected, _ := resp.RowsAffected()
		result := Result{LastInsertID: lastID, RowsAffected: affected}
		results = append(results, result)
	}

	return &ExecuteResponse{Results: results}, nil
}

/*
func handleExecute(connID uint64, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	queries := []string{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&queries); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var resp Response
	results, err := execer.Execute(queries...)
	if err != nil {
		resp.Error = err.Error()
	} else {
		resp.Results = results.Results
		resp.Raft = &results.Raft
	}
	writeResponse(w, r, &resp)
}
*/

func writeResponse(w http.ResponseWriter, r *http.Request, j *ExecuteResponse) {
	pretty, _ := isPretty(r)

	enc := json.NewEncoder(w)
	if pretty {
		enc.SetIndent("", "    ")
	}

	err := enc.Encode(j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/*

// isPretty returns whether the HTTP response body should be pretty-printed.
func isPretty(req *http.Request) (bool, error) {
	return queryParam(req, "pretty")
}

// queryParam returns whether the given query param is set to true.
func queryParam(req *http.Request, param string) (bool, error) {
	err := req.ParseForm()
	if err != nil {
		return false, err
	}
	if _, ok := req.Form[param]; ok {
		return true, nil
	}
	return false, nil
}

// durationParam returns the duration of the given query param, if set.
func durationParam(req *http.Request, param string) (time.Duration, bool, error) {
	q := req.URL.Query()
	t := strings.TrimSpace(q.Get(param))
	if t == "" {
		return 0, false, nil
	}
	dur, err := time.ParseDuration(t)
	if err != nil {
		return 0, false, err
	}
	return dur, true, nil
}

// stmtParam returns the value for URL param 'q', if present.
func stmtParam(req *http.Request) (string, error) {
	q := req.URL.Query()
	return strings.TrimSpace(q.Get("q")), nil
}
*/
