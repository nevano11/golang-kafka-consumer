package repository

import (
	"github.com/jmoiron/sqlx"
	"golang-kafka/internal/entity"
	"golang-kafka/internal/service/filter"
)

type Human interface {
	CreateHuman(human entity.Human) (int, error)
	EditHuman(id int, human entity.Human) (int, error)
	DeleteHuman(id int) error
	GetHumanList(filter filter.Filter) ([]entity.Human, error)
}

type Repository struct {
	Human
}

func NewRepository(database *sqlx.DB) *Repository {
	return &Repository{
		Human: NewHumanPostgres(database),
	}
}
