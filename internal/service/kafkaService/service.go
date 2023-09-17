package kafkaService

import (
	"golang-kafka/internal/entity"
)

type Service struct {
	ConsumeService
	Processor
	ProduceService
}

type HumanSaver interface {
	CreateHuman(human entity.Human) (int, error)
}

func NewKafkaService(fioTopic, fioFailTopic, kafkaConfigPath string, saver HumanSaver) (*Service, error) {
	producer, err := NewKafkaProduceService(fioFailTopic, kafkaConfigPath)
	if err != nil {
		return nil, err
	}

	processor := NewFioProcessor(saver, producer.Produce)

	consumer, err := NewKafkaConsumeService(fioTopic, kafkaConfigPath, processor.ProcessFio)
	if err != nil {
		return nil, err
	}

	return &Service{
		ConsumeService: consumer,
		Processor:      processor,
		ProduceService: producer,
	}, nil
}

func (s *Service) Shutdown() {
	s.ConsumeService.Shutdown()
	s.ProduceService.Shutdown()
}
