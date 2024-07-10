package database

import (
	"context"
	"database/sql"
)

type User struct {
	Id       int64
	Login    string
	Password []byte // Соответствует типу BLOB в SQL
}

func InsertUser(ctx context.Context, db *sql.DB, user *User) (int64, error) {
	const q = `INSERT INTO users (login, password) values ($1, $2)`
	result, err := db.ExecContext(ctx, q, user.Login, user.Password)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func SelectUserByID(ctx context.Context, db *sql.DB, id int64) (User, error) {
	user := User{}
	const q = "SELECT id, login, password FROM users WHERE id = $1"
	err := db.QueryRowContext(ctx, q, id).Scan(&user.Id, &user.Login, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}

func SelectUserByLogin(ctx context.Context, db *sql.DB, login string) (User, error) {
	user := User{}
	const q = "SELECT id, login, password FROM users WHERE login = $1"
	err := db.QueryRowContext(ctx, q, login).Scan(&user.Id, &user.Login, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}
