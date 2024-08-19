package main

import (
	"github.com/Smile8MrBread/Chat/chat_service/internal/app"
	"github.com/Smile8MrBread/Chat/chat_service/internal/config"
	"github.com/Smile8MrBread/Chat/chat_service/internal/logger"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

var (
	broker = "kafka:9092"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)
	log.Info("Starting chat_service service")

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          "1",
		"auto.offset.reset": "smallest",
	})
	if err != nil {
		panic(err)
	}

	application := app.New(log, cfg.StoragePath, cfg.ChatGRPC.Port, consumer)
	go application.GRPCSrv.MustRunChat()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	stopSign := <-stop

	log.Info("Stopping chat_service service", slog.Any("signal", stopSign))

	application.GRPCSrv.StopChat()
}
