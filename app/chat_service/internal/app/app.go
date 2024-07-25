package app

import (
	"github.com/Smile8MrBread/Chat/chat_service/internal/app/grpcapp"
	"github.com/Smile8MrBread/Chat/chat_service/internal/service"
	"github.com/Smile8MrBread/Chat/chat_service/internal/storage/sqlite"
	"log/slog"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, storagePath string, chatPort int) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	chatService := service.New(log, storage, storage, storage)

	gRPC := grpcapp.New(log, chatService, chatPort)

	return &App{GRPCSrv: gRPC}
}
