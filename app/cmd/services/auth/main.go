package main

import (
	"github.com/Smile8MrBread/Chat/app/internal/app"
	"github.com/Smile8MrBread/Chat/app/internal/config"
	"github.com/Smile8MrBread/Chat/app/internal/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)
	log.Info("Starting auth service")

	application := app.New(log, cfg.StoragePath, cfg.AuthGRPC.Port, cfg.ChatGRPC.Port, cfg.TokenTTL)
	go application.GRPCSrv.MustRunAuth()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	stopSign := <-stop

	log.Info("Stopping auth service", slog.Any("signal", stopSign))

	application.GRPCSrv.StopAuth()
}
