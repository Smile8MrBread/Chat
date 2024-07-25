package main

import (
	"github.com/Smile8MrBread/Chat/auth_service/internal/app"
	"github.com/Smile8MrBread/Chat/auth_service/internal/config"
	"github.com/Smile8MrBread/Chat/auth_service/internal/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)
	log.Info("Starting auth_service service")

	application := app.New(log, cfg.StoragePath, cfg.AuthGRPC.Port, cfg.TokenTTL)
	go application.GRPCSrv.MustRunAuth()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	stopSign := <-stop

	log.Info("Stopping auth_service service", slog.Any("signal", stopSign))

	application.GRPCSrv.StopAuth()
}
