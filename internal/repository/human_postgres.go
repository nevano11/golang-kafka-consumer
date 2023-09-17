package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"golang-kafka/internal/entity"
	filter "golang-kafka/internal/service/filter"
	"strings"
)

const (
	create_human_function = "create_human"
	edit_human_function   = "edit_human"
	delete_human_function = "delete_human"
)

type HumanPostgres struct {
	db *sqlx.DB
}

func NewHumanPostgres(db *sqlx.DB) *HumanPostgres {
	return &HumanPostgres{
		db: db,
	}
}

func (a *HumanPostgres) CreateHuman(human entity.Human) (int, error) {
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

func (a *HumanPostgres) EditHuman(id int, human entity.Human) (int, error) {
	logrus.Debugf("Edit human %s with id=%d on EditHuman", human.String(), id)

	humanId := 0
	row := a.db.QueryRow("select * from "+edit_human_function+" ($1, $2, $3, $4, $5, $6, $7)",
		id,
		human.Surname,
		human.FirstName,
		human.LastName,
		human.Age,
		human.Nationality,
		human.Gender)

	err := row.Scan(&humanId)
	if err != nil {
		return humanId, err
	}
	logrus.Debugf("Edited human with id = %d", humanId)
	return humanId, nil
}

func (a *HumanPostgres) DeleteHuman(id int) error {
	logrus.Debugf("Delete human with id=%d on DeleteHuman", id)

	humanId := 0
	row := a.db.QueryRow("select * from "+delete_human_function+" ($1)",
		id)

	err := row.Scan(&humanId)
	if err != nil {
		return err
	}
	logrus.Debugf("Deleted human with id = %d", humanId)
	return nil
}

func (a *HumanPostgres) GetHumanList(filter filter.Filter) ([]entity.Human, error) {
	logrus.Debugf("Get humans with filter=%s on GetHumanList", filter)

	query := strings.Builder{}
	query.WriteString("SELECT * FROM humans WHERE 1 = 1")
	query.WriteString(filter.OptionsToSql())

	var humans []entity.Human
	err := a.db.Select(&humans, query.String())
	if err != nil {
		return nil, err
	}
	logrus.Debugf("Humans: %s", humans)
	return humans, nil
}
