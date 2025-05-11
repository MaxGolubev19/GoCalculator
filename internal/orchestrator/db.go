package orchestrator

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "data/rpn.db")
	if err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	const (
		usersTable = `
	CREATE TABLE IF NOT EXISTS users (
		login TEXT PRIMARY KEY NOT NULL,
		password_hash TEXT NOT NULL
	);`

		expressionsTable = `
	CREATE TABLE IF NOT EXISTS expressions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		expression TEXT NOT NULL,
		user_login TEXT NOT NULL,
		status TEXT NOT NULL,
		result float64 NOT NULL,
		FOREIGN KEY (user_login) REFERENCES users(users_login)
	);`
	)

	if _, err := db.Exec(usersTable); err != nil {
		return err
	}

	if _, err := db.Exec(expressionsTable); err != nil {
		return err
	}

	return nil
}
