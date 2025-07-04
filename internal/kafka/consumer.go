package kafka

import (
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

const (
	sessionTimeout = 45000
	notTimeout     = -1
)

type Consumer struct {
	consumer *kafka.Consumer
	stop     bool
}

func NewConsumer(address []string, topic, consumerGroup string) (*Consumer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(address, ","),
		"group.id":                 consumerGroup,
		"session.timeout.ms":       sessionTimeout,
		"enable.auto.offset.store": false,
		"enable.auto.commit":       true,
		"auto.commit.interval.ms":  5000,
		"auto.offset.reset":        "earliest",
	}
	c, err := kafka.NewConsumer(conf)
	if err != nil {
		return nil, err
	}
	err = c.Subscribe(topic, nil)
	if err != nil {
		return nil, err
	}
	return &Consumer{consumer: c}, nil
}

func (c *Consumer) Start() {
	for {
		if c.stop {
			break
		}
		kafkaMsg, err := c.consumer.ReadMessage(notTimeout)
		if err != nil {
			logrus.Println(err.Error())
		}
		if kafkaMsg == nil {
			continue
		}
		logrus.Infof("Message: %d %s", kafkaMsg.TopicPartition.Offset, string(kafkaMsg.Value))
		_, err = c.consumer.StoreMessage(kafkaMsg)
		if err != nil {
			logrus.Println(err.Error())
			continue
		}
	}
}

func (c *Consumer) Stop() error {
	c.stop = true
	_, err := c.consumer.Commit()
	if err != nil {
		return err
	}
	logrus.Infof("Commited offet")
	return c.consumer.Close()
}
