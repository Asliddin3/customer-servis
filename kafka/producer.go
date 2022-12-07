package kafka

import (
	"context"
	"fmt"
	"time"

	kafka "github.com/segmentio/kafka-go"

	"github.com/Asliddin3/customer-servis/config"

	// "github.com/casbin/casbin/v2/config"
	"github.com/Asliddin3/customer-servis/pkg/logger"
	"github.com/Asliddin3/customer-servis/pkg/messagebroker"
)

type KafkaProduce struct {
	kafkaWriter *kafka.Writer
	log         logger.Logger
}

func NewKafkaProducer(conf config.Config, log logger.Logger, topic string) messagebroker.Producer {
	connString := fmt.Sprintf("%s:%d", conf.KafkaHost, conf.KafkaPort)
	fmt.Println(connString)
	return &KafkaProduce{
		kafkaWriter: &kafka.Writer{
			Addr:         kafka.TCP(connString),
			Topic:        topic,
			BatchTimeout: time.Millisecond * 10,
			// RequiredAcks: kafka.RequireAll,
			// Async:        true,
		},
		log: log,
	}

}

func (p *KafkaProduce) Start() error {
	return nil
}

func (p *KafkaProduce) Stop() error {
	err := p.kafkaWriter.Close()
	if err != nil {
		return err
	}
	return nil
}

func (p *KafkaProduce) Produce(key, body []byte, logBody string) error {
	message := kafka.Message{
		Key:   key,
		Value: body,
	}
	fmt.Println(message)
	if err := p.kafkaWriter.WriteMessages(context.Background(), message); err != nil {
		return err
	}

	p.log.Info("Message produce(key/body): " + string(key) + "/" + logBody)
	return nil
}
