package database

import (
	"context"
	"database/sql"
)

func CreateTables(ctx context.Context, db *sql.DB) error {
	const (
		expressionsTable = `
			CREATE TABLE IF NOT EXISTS expressions(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				status TEXT NOT NULL,
				result REAL NOT NULL,
				content TEXT NOT NULL
			
-- 				FOREIGN KEY (user_id)  REFERENCES expressions (id)
			);`
	)

	if _, err := db.ExecContext(ctx, expressionsTable); err != nil {
		return err
	}

	return nil
}
