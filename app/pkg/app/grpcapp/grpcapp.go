package grpcapp

import (
	"fmt"
	authgrpc "github.com/Smile8MrBread/Chat/app/internal/transport/grpc/auth"
	chatgrpc "github.com/Smile8MrBread/Chat/app/internal/transport/grpc/chat"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"strconv"
)

type App struct {
	log      *slog.Logger
	auth     *grpc.Server
	chat     *grpc.Server
	authPort int
	chatPort int
}

func New(log *slog.Logger, authServ authgrpc.Auth, chatServ chatgrpc.Chat, authPort, chatPort int) *App {
	grpcServ := grpc.NewServer()

	authgrpc.Register(grpcServ, authServ)
	chatgrpc.Register(grpcServ, chatServ)

	return &App{
		log:      log,
		auth:     grpcServ,
		chat:     grpcServ,
		authPort: authPort,
		chatPort: chatPort,
	}
}

func (a *App) MustRunAuth() {
	if err := a.runAuth(); err != nil {
		panic(err)
	}
}

func (a *App) runAuth() error {
	const op = "app.grpcapp.RunAuth"

	a.log.With(
		slog.String("op", op),
		slog.Int("port", a.authPort),
	)

	l, err := net.Listen("tcp", ":"+strconv.Itoa(a.authPort))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("Auth server is running", slog.String("address", l.Addr().String()))

	if err := a.auth.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) StopAuth() {
	const op = "app.grpcapp.StopAuth"

	a.log.With(
		slog.String("op", op),
	).Info("Auth server stopping", slog.Int("port", a.authPort))

	a.auth.GracefulStop()
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
