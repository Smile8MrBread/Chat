package grpcapp

import (
	"fmt"
	chatgrpc "github.com/Smile8MrBread/Chat/chat_service/internal/transport/grpc"
	"github.com/Smile8MrBread/Chat/chat_service/internal/transport/kafka/consumer"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"strconv"
)

type App struct {
	log      *slog.Logger
	chat     *grpc.Server
	chatPort int
	coc      *consumer.OrderConsumer
}

func New(log *slog.Logger, chatServ chatgrpc.Chat, chatPort int, coc *consumer.OrderConsumer) *App {
	grpcServ := grpc.NewServer()

	chatgrpc.Register(grpcServ, chatServ)

	return &App{
		log:      log,
		chat:     grpcServ,
		chatPort: chatPort,
		coc:      coc,
	}
}

func (a *App) MustRunChat() {
	if err := a.runChat(); err != nil {
		panic(err)
	}
}

func (a *App) runChat() error {
	const op = "app.grpcapp.RunChat"

	a.log.With(
		slog.String("op", op),
		slog.Int("port", a.chatPort),
	)

	l, err := net.Listen("tcp", ":"+strconv.Itoa(a.chatPort))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.coc.Init("createMess", a.log)

	a.log.Info("Chat server is running", slog.String("address", l.Addr().String()))

	if err := a.chat.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) StopChat() {
	const op = "app.grpcapp.StopChat"

	a.log.With(
		slog.String("op", op),
	).Info("Chat server stopping", slog.Int("port", a.chatPort))

	a.chat.GracefulStop()
}
