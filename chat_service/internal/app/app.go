package app

import (
	"github.com/Smile8MrBread/Chat/chat_service/internal/app/grpcapp"
	"github.com/Smile8MrBread/Chat/chat_service/internal/service"
	"github.com/Smile8MrBread/Chat/chat_service/internal/storage/sqlite"
	"github.com/Smile8MrBread/Chat/chat_service/internal/transport/kafka/consumer"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log/slog"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, storagePath string, chatPort int, c *kafka.Consumer) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	chatService := service.New(log, storage, storage)

	oc := consumer.NewOrderConsumer(c, "createMess", chatService)
	gRPC := grpcapp.New(log, chatService, chatPort, oc)

	return &App{GRPCSrv: gRPC}
}
