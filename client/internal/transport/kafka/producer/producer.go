package producer

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type OrderProducer struct {
	producer *kafka.Producer
	topic    string
}

func NewOrderProducer(p *kafka.Producer, topic string) *OrderProducer {
	return &OrderProducer{
		producer: p,
		topic:    topic,
	}
}

func (op *OrderProducer) MessOrder(data []byte) error {
	deliverCh := make(chan kafka.Event, 10000)
	defer close(deliverCh)

	err := op.producer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &op.topic,
				Partition: kafka.PartitionAny,
			},
			Value: data,
		},
		deliverCh,
	)
	if err != nil {
		return err
	}

	<-deliverCh
	fmt.Printf("Produced to order: %s", data)

	return nil
}

func Init(broker, groupId, topic string) (*OrderProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"client.id":         groupId,
		"acks":              "all"})
	if err != nil {
		return nil, err
	}

	op := NewOrderProducer(p, topic)

	return op, nil
}
