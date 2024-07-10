package database

import (
	"context"
	"database/sql"
)

func CreateTables(ctx context.Context, db *sql.DB) error {
	const (
		usersTable = `
			CREATE TABLE IF NOT EXISTS users(
			    id INTEGER PRIMARY KEY AUTOINCREMENT,
			    login TEXT NOT NULL UNIQUE,
			    password BLOB NOT NULL
			);`

		expressionsTable = `
			CREATE TABLE IF NOT EXISTS expressions(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				status TEXT NOT NULL,
				result REAL NOT NULL,
				content TEXT NOT NULL,
				user_id INTEGER NOT NULL,
			
				FOREIGN KEY (user_id) REFERENCES users (id)
			);`
	)

	if _, err := db.ExecContext(ctx, usersTable); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, expressionsTable); err != nil {
		return err
	}

	return nil
}
