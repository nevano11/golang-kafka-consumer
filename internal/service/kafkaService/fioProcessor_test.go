package kafkaService

import (
	"github.com/magiconair/properties/assert"
	"github.com/sirupsen/logrus"
	"golang-kafka/internal/entity"
	"testing"
)

type fakeHumanSaver struct{}

func (s *fakeHumanSaver) CreateHuman(human entity.Human) (int, error) {
	return 0, nil
}

func TestValidateFio(t *testing.T) {
	pr := NewFioProcessor(new(fakeHumanSaver), func(fio entity.Fio, message string) {
		logrus.Debug("Send on FioFailed")
	})

	testFio := entity.Fio{
		Surname:   "a",
		FirstName: "a",
		LastName:  "",
	}
	assert.Equal(t, pr.validateFio(testFio), true)

	testFio.Surname = ""
	assert.Equal(t, pr.validateFio(testFio), false)

	testFio.FirstName, testFio.Surname = "", "a"
	assert.Equal(t, pr.validateFio(testFio), false)
}
