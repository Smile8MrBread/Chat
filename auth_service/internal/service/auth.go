package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Smile8MrBread/Chat/auth_service/internal/lib/jwt"
	"github.com/Smile8MrBread/Chat/auth_service/internal/models"
	"github.com/Smile8MrBread/Chat/auth_service/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user exists")
	ErrUserNotFound       = errors.New("user not found")
)

type UserSaver interface {
	SaveUser(ctx context.Context, firstName, lastName, login, avatar string, passHash []byte) (int64, error)
}

type UserProvider interface {
	ProvideUser(ctx context.Context, login string) (models.User, error)
}

type App struct {
	log          *slog.Logger
	userProvider UserProvider
	userSaver    UserSaver
	tokenTTL     time.Duration
}

func New(log *slog.Logger, userProvider UserProvider, userSaver UserSaver, tokenTTL time.Duration) *App {
	return &App{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		tokenTTL:     tokenTTL,
	}
}

func (a *App) Login(ctx context.Context, login, password string) (token string, err error) {
	const op = "service.auth_service.Login"

	log := a.log.With(slog.String("op", op), slog.String("login", login))
	log.Info("Attempting to login")

	user, err := a.userProvider.ProvideUser(ctx, login)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("User not found", slog.String("error", err.Error()))
			return "", ErrUserNotFound
		}

		log.Error("Failed to get user", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Warn("Invalid credentials", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s:%w", op, ErrInvalidCredentials)
	}

	log.Info("Login successful")

	token, err = jwt.NewToken(user, a.tokenTTL)
	if err != nil {
		log.Error("Failed to create token", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s, %w", op, err)
	}

	return token, nil
}

func (a *App) Registration(ctx context.Context, firstName, lastName, login, password, avatar string) (id int64, err error) {
	const op = "service.auth_service.Registration"

	log := a.log.With(slog.String("op", op), slog.String("login", login))
	log.Info("Registration new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to generate passHash", slog.String("error", err.Error()))

		return -1, fmt.Errorf("%s:%w", op, err)
	}

	id, err = a.userSaver.SaveUser(ctx, firstName, lastName, login, avatar, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Error("User exists", slog.String("error", err.Error()))
			return -1, ErrUserExists
		}

		log.Error("Failed to save user", slog.String("error", err.Error()))
		return -1, fmt.Errorf("%s:%w", op, err)
	}

	return id, nil
}
