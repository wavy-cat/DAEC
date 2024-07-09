package database

import (
	"context"
	"database/sql"
)

type (
	Expression struct {
		Id      int64
		Status  string
		Result  float64
		Content string
	}
)

func InsertExpression(ctx context.Context, db *sql.DB, expression *Expression) (int64, error) {
	const q = `INSERT INTO expressions (status, result, content) values ($1, $2, $3)`
	result, err := db.ExecContext(ctx, q, expression.Status, expression.Result, expression.Content)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func SelectExpressions(ctx context.Context, db *sql.DB) ([]Expression, error) {
	var expressions []Expression
	const q = "SELECT id, status, result, content FROM expressions"

	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		e := Expression{}
		err := rows.Scan(&e.Id, &e.Status, &e.Result, &e.Content)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, e)
	}

	return expressions, nil
}

func SelectExpressionByID(ctx context.Context, db *sql.DB, id int64) (Expression, error) {
	u := Expression{}
	const q = "SELECT id, status, result, content FROM expressions WHERE id = $1"
	err := db.QueryRowContext(ctx, q, id).Scan(&u.Id, &u.Status, &u.Result, &u.Content)
	if err != nil {
		return u, err
	}

	return u, nil
}

func UpdateExpression(ctx context.Context, db *sql.DB, id int64, status string, result float64) error {
	const q = "UPDATE expressions SET status = $1, result = $2 WHERE id = $3"
	_, err := db.ExecContext(ctx, q, status, result, id)
	if err != nil {
		return err
	}

	return nil
}
