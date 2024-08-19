package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Smile8MrBread/Chat/chat_service/internal/models"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"reflect"
	"strconv"
)

type MessageCreator interface {
	CreateMessage(ctx context.Context, text, date string, userId, contactId int64) (int64, error)
}

type OrderConsumer struct {
	consumer *kafka.Consumer
	topic    string
	creator  MessageCreator
}

func NewOrderConsumer(c *kafka.Consumer, topic string, creator MessageCreator) *OrderConsumer {
	return &OrderConsumer{
		consumer: c,
		topic:    topic,
		creator:  creator,
	}
}

func (oc *OrderConsumer) Init(topic string, log *slog.Logger) {
	go func() {
		for {
			if err := oc.messListener(topic); err != nil {
				log.Error("Kafka consumer error", slog.String("error", err.Error()))
				panic(err)
			}
		}
	}()
}

func (oc *OrderConsumer) messListener(topic string) error {
	err := oc.consumer.Subscribe(topic, nil)
	if err != nil {
		return err
	}

	for {
		ev := oc.consumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			fmt.Println("Revived", string(e.Value))
			buf := e.Value
			data := models.Message{}

			err = json.Unmarshal(buf, &data)
			if err != nil {
				return err
			}

			userId, err := strconv.Atoi(data.UserFrom)
			if err != nil {
				return err
			}
			contactId, err := strconv.Atoi(data.UserTo)
			if err != nil {
				return err
			}

			_, err = oc.creator.CreateMessage(context.Background(), data.Text, data.Date, int64(userId), int64(contactId))
			if err != nil {
				return err
			}
		case kafka.Error:
			return err
		}
	}
}

func (oc *OrderConsumer) CreateMessage(ctx context.Context, text, date string, userId, contactId int64) (int64, error) {
	if reflect.TypeOf(text).Kind() != reflect.String &&
		reflect.TypeOf(date).Kind() != reflect.String &&
		reflect.TypeOf(userId).Kind() != reflect.Int64 &&
		reflect.TypeOf(contactId).Kind() != reflect.Int64 {
		return -1, status.Error(codes.InvalidArgument, "Invalid message")
	}

	id, err := oc.creator.CreateMessage(ctx, text, date, userId, contactId)
	if err != nil {
		return -1, status.Error(codes.Internal, "Internal error")
	}

	return id, nil
}
