package models

import "database/sql"

type User struct {
	ID   int    `db:"id" json:"id" formam:"id"`
	Name string `db:"name" json:"name" formam:"name"`
}

func (db *DB) GetUsers() ([]*User, error) {
	var users []*User
	err := db.Select(&users, "SELECT * FROM users ORDER BY id ASC")
	return users, err
}

func (db *DB) CreateUser(u *User) (*User, error) {
	var user User
	err := db.Get(&user, "INSERT INTO users (name) VALUES ($1) RETURNING *", u.Name)
	return &user, err
}

func (db *DB) GetUserByID(id int) (*User, error) {
	var user User
	err := db.Get(&user, "SELECT * FROM users WHERE id=$1", id)
	return &user, err
}

func (db *DB) UpdateUser(u *User) error {
	result, err := db.Exec("UPDATE users SET name=$1 WHERE id=$2", u.Name, u.ID)
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

func (db *DB) DeleteUser(id int) error {
	_, err := db.Exec("DELETE FROM users WHERE id=$1", id)
	return err
}
