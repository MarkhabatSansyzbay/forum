package repository

import (
	"database/sql"
	"fmt"
)

const (
	usersTable   = "users"
	postsTable   = "posts"
	sessionTable = "sessions"
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
	_, err := db.Exec(fmt.Sprintf("CREATE TABLE if not exists %s (id INTEGER PRIMARY KEY, email varchar(319), username varchar(64), password text)", usersTable))
	if err != nil {
		return err
	}

	_, err = db.Exec(fmt.Sprintf("CREATE TABLE if not exists %s (id INTEGER PRIMARY KEY, user_id INTEGER, token varchar(32), expiration_date DATETIME, FOREIGN KEY(user_id) REFERENCES users(id))", sessionTable))
	if err != nil {
		return err
	}
	// FOR POSTS AND ...
	return nil
}
