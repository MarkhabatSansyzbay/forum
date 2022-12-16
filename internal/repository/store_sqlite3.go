package repository

import (
	"database/sql"
	"fmt"
)

func OpenSqliteDB(dbName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("./%s", dbName))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	if err = createTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS USERS(
			ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			Username TEXT NOT NULL UNIQUE,
			Email TEXT NOT NULL UNIQUE,
			Password TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS SESSIONS(
			ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			UserID INTEGER NOT NULL UNIQUE,
			Token VARCHAR(32) NOT NULL,
			ExpDate DATATIME NOT NULL
		);
	`
	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}
