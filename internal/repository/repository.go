package repository

import (
	"github.com/jmoiron/sqlx"
	"golang-kafka/internal/entity"
)

type Human interface {
	CreateHuman(human entity.DbFio) (int, error)
}

type Repository struct {
	Human
}

func NewRepository(database *sqlx.DB) *Repository {
	return &Repository{
		Human: NewHumanPostgres(database),
	}
}
