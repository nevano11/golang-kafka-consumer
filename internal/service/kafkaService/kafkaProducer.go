package kafkaService

import (
	"bufio"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
	"golang-kafka/internal/entity"
	"os"
	"strings"
)

type ProduceService interface {
	Produce(fio entity.Fio, message string)
	Shutdown()
}

type KafkaProducer struct {
	topicName  string
	configFile string
	producer   *kafka.Producer
}

func NewKafkaProduceService(topicName, configFile string) (*KafkaProducer, error) {
	kafkaService := &KafkaProducer{
		topicName:  topicName,
		configFile: configFile,
		producer:   nil,
	}
	configMap, err := kafkaService.readConfig()
	if err != nil {
		return nil, err
	}
	kafkaService.producer, err = kafka.NewProducer(&configMap)
	if err != nil {
		return nil, err
	}

	go func() {
		for e := range kafkaService.producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					logrus.Errorf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					logrus.Infof("Produced event to topic %s: value = %s\n",
						*ev.TopicPartition.Topic, string(ev.Value))
				}
			}
		}
	}()

	return kafkaService, nil
}

func (s *KafkaProducer) readConfig() (kafka.ConfigMap, error) {
	m := make(map[string]kafka.ConfigValue)

	file, err := os.Open(s.configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "#") && len(line) != 0 {
			before, after, found := strings.Cut(line, "=")
			if found {
				parameter := strings.TrimSpace(before)
				value := strings.TrimSpace(after)
				m[parameter] = value
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *KafkaProducer) Produce(fio entity.Fio, message string) {
	logrus.Debugf("Produce fio on topic %s with message %s", s.topicName, message)
	fioMarshalled, err := json.Marshal(fio)
	if err != nil {
		logrus.Errorf("Failed to marshal fio to send on topic %s: %s", s.topicName, err.Error())
		return
	}
	err = s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &s.topicName, Partition: kafka.PartitionAny},
		Key:            []byte(message),
		Value:          fioMarshalled,
	}, nil)

	if err != nil {
		logrus.Errorf("Failed to produce fio to send on topic %s: %s", s.topicName, err.Error())
		return
	}
}

func (s *KafkaProducer) Shutdown() {
	logrus.Debug("Shutdown kafkaProducer")
	s.producer.Flush(15 * 1000)
	s.producer.Close()
}
