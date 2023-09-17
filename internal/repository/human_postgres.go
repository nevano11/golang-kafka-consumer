package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"golang-kafka/internal/entity"
)

const (
	create_human_function = "create_human"
)

type HumanPostgres struct {
	db *sqlx.DB
}

func NewHumanPostgres(db *sqlx.DB) *HumanPostgres {
	return &HumanPostgres{
		db: db,
	}
}

func (a *HumanPostgres) CreateHuman(human entity.DbFio) (int, error) {
	newUserId := 0
	logrus.Debugf("Create human %s on HumanPostgres", human.String())

	row := a.db.QueryRow("select * from "+create_human_function+" ($1, $2, $3, $4, $5, $6)",
		human.Surname,
		human.FirstName,
		human.LastName,
		human.Age,
		human.Nationality,
		human.Gender)

	err := row.Scan(&newUserId)
	if err != nil {
		return newUserId, err
	}
	logrus.Debugf("Created human id = %d", newUserId)
	return newUserId, nil
}
