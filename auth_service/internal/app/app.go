package app

import (
	"github.com/Smile8MrBread/Chat/auth_service/internal/app/grpcapp"
	"github.com/Smile8MrBread/Chat/auth_service/internal/service"
	"github.com/Smile8MrBread/Chat/auth_service/internal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, storagePath string, authPort int, tokenTTL time.Duration) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := service.New(log, storage, storage, tokenTTL)

	gRPC := grpcapp.New(log, authService, authPort)

	return &App{GRPCSrv: gRPC}
}
