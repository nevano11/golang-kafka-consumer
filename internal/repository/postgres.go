package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type DbConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DbName   string
	SslMode  string
}

func NewPostgresDb(config DbConfig) (*sqlx.DB, error) {
	c := config
	db, err := sqlx.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			c.Host, c.Port, c.Username, c.DbName, c.Password, c.SslMode))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
