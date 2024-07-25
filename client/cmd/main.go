package main

import (
	authGrpc "github.com/Smile8MrBread/Chat/auth_service/proto/gen"
	chatGrpc "github.com/Smile8MrBread/Chat/chat_service/proto/gen"
	"github.com/Smile8MrBread/Chat/client/internal/config"
	"github.com/Smile8MrBread/Chat/client/internal/logger"
	"github.com/Smile8MrBread/Chat/client/internal/transport/rest"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.MustLoad()
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	log := logger.SetupLogger(cfg.Env)

	connAuth, err := grpc.NewClient("auth_service:32100", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer connAuth.Close()
	if err != nil {
		panic(err)
	}

	connChat, err := grpc.NewClient("chat_service:32200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer connChat.Close()
	if err != nil {
		panic(err)
	}

	authClient := authGrpc.NewAuthClient(connAuth)
	chatClient := chatGrpc.NewChatClient(connChat)

	rest.StartServer(log, r, authClient, chatClient)
}
