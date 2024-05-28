package models

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	DBName   string `envconfig:"NAME"`
	User     string `envconfig:"USER"`
	Password string `envconfig:"PASSWORD"`
	Host     string `envconfig:"HOST"`
	Port     int    `envconfig:"PORT"`
	SSLMode  string `envconfig:"SSL_MODE"`
}

func (c Config) ConnectionString() string {
	return fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d",
		c.DBName, c.User, c.Password, c.Host, c.Port)
}

func (c Config) URL() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}

type DB struct {
	conf Config

	*sqlx.DB
}

func NewDB(conf Config) (*DB, error) {
	db, err := sqlx.Connect("pgx", conf.ConnectionString())
	if err != nil {
		return nil, err
	}
	return &DB{conf: conf, DB: db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}
