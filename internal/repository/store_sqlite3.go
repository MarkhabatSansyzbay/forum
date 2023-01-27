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
			Password TEXT
		);
		CREATE TABLE IF NOT EXISTS SESSIONS(
			ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			UserID INTEGER NOT NULL UNIQUE,
			Token VARCHAR(32) NOT NULL,
			ExpDate DATATIME NOT NULL
		);
		CREATE TABLE IF NOT EXISTS POSTS(
			ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			AuthorID INTEGER NOT NULL,
			Title TEXT NOT NULL,
			Content TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS COMMENTS(
			ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			AuthorID INTEGER NOT NULL,
			PostID INTEGER NOT NULL,
			Content TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS REACTIONS(
			ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			UserID INTEGER NOT NULL,
			PostID INTEGER,
			CommentID INTEGER,
			VOTE BLOB NOT NULL
		);
		CREATE TABLE IF NOT EXISTS CATEGORIES(
			ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			PostID INTEGER NOT NULL,
			Category TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS IMAGES(
			ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			PostID INTEGER,
			Image TEXT,
			FOREIGN KEY(PostID) REFERENCES POSTS(ID)
		)
	`
	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}
