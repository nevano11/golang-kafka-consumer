package entity

import "fmt"

type Human struct {
	Id          int    `json:"id" db:"id"`
	Surname     string `json:"surname" db:"surname"`
	FirstName   string `json:"first_name" db:"name"`
	LastName    string `json:"last_name" db:"patronymic"`
	Age         int    `json:"age" db:"age"`
	Nationality string `json:"nationality" db:"nationality"`
	Gender      string `json:"gender" db:"gender"`
}

func (h Human) String() string {
	return fmt.Sprintf(
		"Human:{Id:%d, Surname:%s, FirstName:%s, LastName:%s, Age:%d, Nationality:%s, Gender:%s}",
		h.Id, h.Surname, h.FirstName, h.LastName, h.Age, h.Nationality, h.Gender)
}
