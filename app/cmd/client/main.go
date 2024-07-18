package main

import (
	"github.com/Smile8MrBread/Chat/app/internal/config"
	"github.com/Smile8MrBread/Chat/app/internal/logger"
	"github.com/Smile8MrBread/Chat/app/internal/transport/rest"
	authGrpc "github.com/Smile8MrBread/Chat/protos/gen/auth"
	chatGrpc "github.com/Smile8MrBread/Chat/protos/gen/chat"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.MustLoad()
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	log := logger.SetupLogger(cfg.Env)

	connAuth, err := grpc.NewClient(":32100", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer connAuth.Close()
	if err != nil {
		panic(err)
	}

	connChat, err := grpc.NewClient(":32200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer connChat.Close()
	if err != nil {
		panic(err)
	}

	authClient := authGrpc.NewAuthClient(connAuth)
	chatClient := chatGrpc.NewChatClient(connChat)

	rest.StartServer(log, r, authClient, chatClient)
}
