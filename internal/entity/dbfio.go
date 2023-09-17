package entity

import "fmt"

type DbFio struct {
	Id          int    `json:"-"`
	Surname     string `json:"surname"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Age         int    `json:"age"`
	Nationality string `json:"nationality"`
	Gender      string `json:"gender"`
}

func (f DbFio) String() string {
	return fmt.Sprintf(
		"DbFio:{Id:%d, Surname:%s, FirstName:%s, LastName:%s, Age:%d, Nationality:%s, Gender:%s}",
		f.Id, f.Surname, f.FirstName, f.LastName, f.Age, f.Nationality, f.Gender)
}
