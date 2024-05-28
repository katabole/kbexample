package models

type User struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

func (db *DB) GetUserByID(id int) (*User, error) {
	user := &User{}
	err := db.Get(user, "SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (db *DB) SaveUser(u *User) error {
	_, err := db.Exec("INSERT INTO users (id, name) VALUES ($1, $2)", u.ID, u.Name)
	return err
}
