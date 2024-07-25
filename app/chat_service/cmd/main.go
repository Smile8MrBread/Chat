package main

import (
	"github.com/Smile8MrBread/Chat/chat_service/internal/app"
	"github.com/Smile8MrBread/Chat/chat_service/internal/config"
	"github.com/Smile8MrBread/Chat/chat_service/internal/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)
	log.Info("Starting chat_service service")

	application := app.New(log, cfg.StoragePath, cfg.ChatGRPC.Port)
	go application.GRPCSrv.MustRunChat()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	stopSign := <-stop

	log.Info("Stopping chat_service service", slog.Any("signal", stopSign))

	application.GRPCSrv.StopChat()
}
