package crudService

import (
	"golang-kafka/internal/entity"
	"golang-kafka/internal/repository"
	"golang-kafka/internal/service/filter"
)

type Service struct {
	repos *repository.Repository
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		repos: repos,
	}
}

func (s *Service) CreateHuman(human entity.Human) (int, error) {
	humanId, err := s.repos.CreateHuman(human)
	if err != nil {
		return 0, err
	}
	return humanId, nil
}

func (s *Service) EditHuman(id int, human entity.Human) error {
	_, err := s.repos.EditHuman(id, human)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteHuman(id int) error {
	return s.repos.DeleteHuman(id)
}

func (s *Service) GetHumanList(filter filter.Filter) ([]entity.Human, error) {
	return s.repos.GetHumanList(filter)
}
