package models

import (
	"context"
	"database/sql"
)

type User struct {
	ID   int    `db:"id" json:"id" formam:"id"`
	Name string `db:"name" json:"name" formam:"name"`
}

func (db *DB) GetUsers(ctx context.Context) ([]*User, error) {
	var users []*User
	err := db.SelectContext(ctx, &users, "SELECT * FROM users ORDER BY id ASC")
	return users, err
}

func (db *DB) CreateUser(ctx context.Context, u *User) (*User, error) {
	var newUser User
	err := db.QueryRowxContext(ctx, "INSERT INTO users (name) VALUES ($1) RETURNING *", u.Name).StructScan(&newUser)
	return &newUser, err
}

func (db *DB) GetUserByID(ctx context.Context, id int) (*User, error) {
	user := &User{}
	err := db.QueryRowxContext(ctx, "SELECT * FROM users WHERE id=$1", id).StructScan(user)
	return user, err
}

func (db *DB) UpdateUser(ctx context.Context, u *User) error {
	result, err := db.ExecContext(ctx, "UPDATE users SET name=$1 WHERE id=$2", u.Name, u.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (db *DB) DeleteUser(ctx context.Context, id int) error {
	_, err := db.ExecContext(ctx, "DELETE FROM users WHERE id=$1", id)
	return err
}
