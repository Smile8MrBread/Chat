package auth

import (
	"context"
	"errors"
	"github.com/Smile8MrBread/Chat/app/internal/service/auth"
	"github.com/Smile8MrBread/Chat/protos/gen/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, login, password string) (token string, err error)
	Registration(ctx context.Context, firstName, lastName, login, password, avatar string) (id int64, err error)
}

type ServerAPI struct {
	authGrpc.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, authStr Auth) {
	authGrpc.RegisterAuthServer(gRPC, &ServerAPI{auth: authStr})
}

func (s *ServerAPI) Login(ctx context.Context, req *authGrpc.LoginRequest) (*authGrpc.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) || errors.Is(err, auth.ErrUserNotFound) {
			return nil, status.Error(codes.InvalidArgument, "Invalid login or password")
		}

		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &authGrpc.LoginResponse{
		Token: token,
	}, nil
}

func (s *ServerAPI) Registration(ctx context.Context, req *authGrpc.RegisterRequest) (*authGrpc.RegisterResponse, error) {
	if err := validateRegistration(req); err != nil {
		return nil, err
	}

	id, err := s.auth.Registration(ctx, req.GetFirstName(), req.GetLastName(), req.GetLogin(), req.GetPassword(), req.GetAvatar())
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user exists")
		}

		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &authGrpc.RegisterResponse{
		UserId: id,
	}, nil
}

func validateLogin(req *authGrpc.LoginRequest) error {
	if req.GetLogin() == "" || len(req.GetLogin()) > 255 {
		return status.Error(codes.InvalidArgument, "Invalid login or password")
	}

	if req.GetPassword() == "" || len(req.GetPassword()) > 255 {
		return status.Error(codes.InvalidArgument, "Invalid login or password")
	}

	return nil
}

func validateRegistration(req *authGrpc.RegisterRequest) error {
	if req.GetLogin() == "" || len(req.GetLogin()) > 255 {
		return status.Error(codes.InvalidArgument, "Invalid login")
	}

	if req.GetPassword() == "" || len(req.GetPassword()) > 255 {
		return status.Error(codes.InvalidArgument, "Invalid password")
	}

	if req.GetFirstName() == "" || len(req.GetFirstName()) > 255 {
		return status.Error(codes.InvalidArgument, "Invalid first name")
	}

	if req.GetLastName() == "" || len(req.GetLastName()) > 255 {
		return status.Error(codes.InvalidArgument, "Invalid last name")
	}

	return nil
}
