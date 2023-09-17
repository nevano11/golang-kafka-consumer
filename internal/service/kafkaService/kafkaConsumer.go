package kafkaService

import (
	"bufio"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
	"golang-kafka/internal/entity"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type ConsumeService interface {
	Shutdown()
}

type kafkaConsumer struct {
	topicName  string
	configFile string
	consumer   *kafka.Consumer
	onConsume  func(fio entity.Fio)
}

func NewKafkaConsumeService(topicName, configFile string, onConsume func(fio entity.Fio)) (*kafkaConsumer, error) {
	kafkaService := &kafkaConsumer{
		topicName:  topicName,
		configFile: configFile,
		onConsume:  onConsume,
		consumer:   nil,
	}
	configMap, err := kafkaService.readConfig()
	if err != nil {
		return nil, err
	}
	kafkaService.consumer, err = kafka.NewConsumer(&configMap)
	if err != nil {
		return nil, err
	}

	topic := kafkaService.topicName
	err = kafkaService.consumer.SubscribeTopics([]string{topic}, nil)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	go func() {
		run := true
		for run {
			select {
			case sig := <-sigchan:
				logrus.Infof("Caught signal %v: terminating\n", sig)
				run = false
			default:
				ev, err := kafkaService.consumer.ReadMessage(100 * time.Millisecond)
				if err != nil {
					// Waiting message
					continue
				}
				logrus.Infof("Consumed event from topic %s: value = %s\n",
					*ev.TopicPartition.Topic, string(ev.Value))
				var fio entity.Fio
				if err := json.Unmarshal(ev.Value, &fio); err != nil {
					logrus.Errorf("Failed to marshal info. In obj = %s", ev.Value)
				} else {
					kafkaService.onConsume(fio)
				}
			}
		}
	}()

	return kafkaService, nil
}

func (s *kafkaConsumer) readConfig() (kafka.ConfigMap, error) {
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

func (s *kafkaConsumer) Shutdown() {
	logrus.Debug("Shutdown kafkaConsumer")
	s.consumer.Close()
}
