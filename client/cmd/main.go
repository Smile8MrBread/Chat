package main

import (
	"client/internal/config"
	"client/internal/logger"
	"client/internal/transport/kafka/producer"
	"client/internal/transport/rest"
	authGrpc "github.com/Smile8MrBread/Chat/auth_service/proto/gen"
	chatGrpc "github.com/Smile8MrBread/Chat/chat_service/proto/gen"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	broker = "kafka:9092"
	topic  = "createMess"
)

func main() {
	cfg := config.MustLoad()
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	log := logger.SetupLogger(cfg.Env)

	connAuth, err := grpc.NewClient("auth:32100", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer connAuth.Close()
	if err != nil {
		panic(err)
	}

	connChat, err := grpc.NewClient("chat:32200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer connChat.Close()
	if err != nil {
		panic(err)
	}

	p, err := producer.Init(broker, "1", topic)

	if err != nil {
		panic(err)
	}

	authClient := authGrpc.NewAuthClient(connAuth)
	chatClient := chatGrpc.NewChatClient(connChat)

	rest.StartServer(log, r, authClient, chatClient, p)
}
