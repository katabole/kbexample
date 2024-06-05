package models

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
	var newUser User
	err := db.QueryRowx("INSERT INTO users (name) VALUES ($1) RETURNING *", u.Name).StructScan(&newUser)
	return &newUser, err
}

func (db *DB) GetUserByID(id int) (*User, error) {
	user := &User{}
	err := db.QueryRowx("SELECT * FROM users WHERE id=$1", id).StructScan(user)
	return user, err
}

func (db *DB) UpdateUser(u *User) error {
	_, err := db.Exec("UPDATE users SET name=$1 WHERE id=$2", u.Name, u.ID)
	return err
}

func (db *DB) DeleteUser(id int) error {
	_, err := db.Exec("DELETE FROM users WHERE id=$1", id)
	return err
}
