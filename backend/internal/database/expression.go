package database

import (
	"context"
	"database/sql"
)

type Expression struct {
	Id      int64
	Status  string
	Result  float64
	Content string
	UserId  int64
}

func InsertExpression(ctx context.Context, db *sql.DB, expression *Expression) (int64, error) {
	const q = `INSERT INTO expressions (status, result, content, user_id) values ($1, $2, $3, $4)`
	result, err := db.ExecContext(ctx, q, expression.Status, expression.Result, expression.Content, expression.UserId)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func SelectUserExpressions(ctx context.Context, db *sql.DB, userId int64) ([]Expression, error) {
	var expressions []Expression
	const q = "SELECT id, status, result, content, user_id FROM expressions WHERE user_id = $1"

	rows, err := db.QueryContext(ctx, q, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		e := Expression{}
		err := rows.Scan(&e.Id, &e.Status, &e.Result, &e.Content, &e.UserId)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, e)
	}

	return expressions, nil
}

func SelectUserExpressionByID(ctx context.Context, db *sql.DB, userId int64, id int64) (Expression, error) {
	expression := Expression{}
	const q = "SELECT id, status, result, content, user_id FROM expressions WHERE id = $1 AND user_id = $2"
	err := db.QueryRowContext(ctx, q, id, userId).
		Scan(&expression.Id, &expression.Status, &expression.Result, &expression.Content, &expression.UserId)
	if err != nil {
		return expression, err
	}

	return expression, nil
}

func UpdateExpression(ctx context.Context, db *sql.DB, expression *Expression) error {
	const q = "UPDATE expressions SET status = $1, result = $2 WHERE id = $3"
	_, err := db.ExecContext(ctx, q, expression.Status, expression.Result, expression.Id)
	if err != nil {
		return err
	}

	return nil
}

func SelectAllExpressions(ctx context.Context, db *sql.DB) ([]Expression, error) {
	var expressions []Expression
	const q = "SELECT id, status, result, content, user_id FROM expressions"

	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		e := Expression{}
		err := rows.Scan(&e.Id, &e.Status, &e.Result, &e.Content, &e.UserId)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, e)
	}

	return expressions, nil
}
