package service

import (
	"golang-kafka/internal/entity"
)

type KafkaService struct {
	ConsumeService
	Processor
	ProduceService
}

type HumanSaver interface {
	CreateHuman(human entity.DbFio) (int, error)
}

func NewKafkaService(fioTopic, fioFailTopic, kafkaConfigPath string, saver HumanSaver) (*KafkaService, error) {
	producer, err := NewKafkaProduceService(fioFailTopic, kafkaConfigPath)
	if err != nil {
		return nil, err
	}

	processor := NewFioProcessor(saver, producer.Produce)

	consumer, err := NewKafkaConsumeService(fioTopic, kafkaConfigPath, processor.ProcessFio)
	if err != nil {
		return nil, err
	}

	return &KafkaService{
		ConsumeService: consumer,
		Processor:      processor,
		ProduceService: producer,
	}, nil
}

func (s *KafkaService) Shutdown() {
	s.ConsumeService.Shutdown()
	s.ProduceService.Shutdown()
}
