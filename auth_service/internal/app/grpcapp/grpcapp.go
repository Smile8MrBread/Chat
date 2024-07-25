package grpcapp

import (
	"fmt"
	authgrpc "github.com/Smile8MrBread/Chat/auth_service/internal/transport/grpc"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"strconv"
)

type App struct {
	log      *slog.Logger
	auth     *grpc.Server
	authPort int
}

func New(log *slog.Logger, authServ authgrpc.Auth, authPort int) *App {
	grpcServ := grpc.NewServer()

	authgrpc.Register(grpcServ, authServ)

	return &App{
		log:      log,
		auth:     grpcServ,
		authPort: authPort,
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
