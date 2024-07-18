package app

import (
	"github.com/Smile8MrBread/Chat/app/internal/service/auth"
	"github.com/Smile8MrBread/Chat/app/internal/service/chat"
	"github.com/Smile8MrBread/Chat/app/internal/storage/sqlite"
	"github.com/Smile8MrBread/Chat/app/pkg/app/grpcapp"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, storagePath string, authPort, chatPort int, tokenTTL time.Duration) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, tokenTTL)
	chatService := chat.New(log, storage, storage, storage)

	gRPC := grpcapp.New(log, authService, chatService, authPort, chatPort)

	return &App{GRPCSrv: gRPC}
}
