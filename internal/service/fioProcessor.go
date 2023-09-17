package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang-kafka/internal/entity"
	"golang-kafka/internal/entity/dao"
	"io"
	"net/http"
	"sync"
)

type Processor interface {
	ProcessFio(fio entity.Fio)
}

type FioProcessor struct {
	onNonValid func(fio entity.Fio, message string)
	saver      HumanSaver
}

func NewFioProcessor(saver HumanSaver, onNonValid func(fio entity.Fio, message string)) *FioProcessor {
	return &FioProcessor{
		saver:      saver,
		onNonValid: onNonValid,
	}
}

func (s *FioProcessor) ProcessFio(fio entity.Fio) {
	logrus.Infof("Fio to consume: %s", fio.String())

	// validate
	logrus.Debugf("Validating: %s", fio.String())
	isValid := s.validateFio(fio)
	if !isValid {
		logrus.Warningf("Fio is invalid: %s", fio.String())
		s.onNonValid(fio, "Invalid fio")
		return
	}

	// age, sex, nationality
	logrus.Debugf("Add fields on: %s", fio.String())
	dbfio, err := s.addNewFields(fio)
	if err != nil {
		logrus.Warningf("Failed to add new fields: %s", fio.String())
		s.onNonValid(fio, err.Error())
		return
	}

	// add on repository
	logrus.Debugf("Saving on repository: %s", dbfio.String())
	_, err = s.saver.CreateHuman(dbfio)
	if err != nil {
		logrus.Warningf("Failed to save on db: %s", dbfio.String())
		s.onNonValid(fio, err.Error())
		return
	}
}

func (s *FioProcessor) validateFio(fio entity.Fio) bool {
	isGood := func(s string) bool {
		return len(s) > 0
	}
	// last name is not required
	return isGood(fio.Surname) && isGood(fio.FirstName)
}

func (s *FioProcessor) addNewFields(fio entity.Fio) (entity.DbFio, error) {
	// result entity
	dbfio := entity.DbFio{
		Id:          0,
		Surname:     fio.Surname,
		FirstName:   fio.FirstName,
		LastName:    fio.LastName,
		Age:         0,
		Nationality: "",
		Gender:      "",
	}
	// waitgroup
	var wg sync.WaitGroup
	hasFailOnApi := false

	// apis
	apis := make(map[string]func(b []byte))
	apis[fmt.Sprintf("https://api.agify.io/?name=%s", fio.FirstName)] =
		func(b []byte) {
			var target dao.Agify
			if err := json.Unmarshal(b, &target); err != nil {
				hasFailOnApi = true
				logrus.Error("failed to unmarshal Agify" + string(b))
			} else {
				dbfio.Age = target.Age
			}
		}
	apis[fmt.Sprintf("https://api.genderize.io/?name=%s", fio.FirstName)] =
		func(b []byte) {
			var target dao.Genderize
			if err := json.Unmarshal(b, &target); err != nil {
				hasFailOnApi = true
				logrus.Error("failed to unmarshal Genderize " + string(b))
			} else {
				dbfio.Gender = target.Gender
			}
		}
	apis[fmt.Sprintf("https://api.nationalize.io/?name=%s", fio.FirstName)] =
		func(b []byte) {
			var target dao.Nationalize
			if err := json.Unmarshal(b, &target); err != nil {
				hasFailOnApi = true
				logrus.Error("failed to unmarshal Nationalize" + string(b))
			} else {
				if len(target.Country) > 0 {
					dbfio.Nationality = target.Country[0].CountryId
				} else {
					hasFailOnApi = true
					logrus.Error("failed on Nationalize api. target.Country.len == 0")
				}
			}
		}
	// parallel send request goroutine
	sendRequest := func(url string, entityFiller func(b []byte)) {
		defer wg.Done()
		logrus.Debugf("sending get-request on url %s", url)
		res, err := http.Get(url)
		if err != nil {
			logrus.Errorf("error making http request to url%s: %s\n", url, err)
			hasFailOnApi = true
			return
		}
		defer res.Body.Close()

		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			logrus.Errorf("failed to read body%s: %s\n", url, err)
			hasFailOnApi = true
			return
		}
		entityFiller(bytes)
	}

	// send get requests
	wg.Add(len(apis))
	for k, v := range apis {
		go sendRequest(k, v)
	}
	wg.Wait()

	if hasFailOnApi {
		return dbfio, errors.New("failed on requests")
	}
	return dbfio, nil
}
