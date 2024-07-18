package chat

import (
	"context"
	"errors"
	"fmt"
	"github.com/Smile8MrBread/Chat/app/internal/models"
	"github.com/Smile8MrBread/Chat/app/internal/storage"
	"log/slog"
	"strconv"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrMessageNotFound = errors.New("message not found")
)

type UserIdent interface {
	IdentUser(ctx context.Context, id int64) (models.User, error)
}

type ContactsWorker interface {
	AddContact(ctx context.Context, id, contactId int64) error
	AllContacts(ctx context.Context, id int64) ([]int64, error)
	IsMessaged(ctx context.Context, id, contactId int64) error
	AllMessaged(ctx context.Context, id int64) ([]int64, error)
}

type MessageWorker interface {
	CreateMessage(ctx context.Context, text, date string, userId, contactId int64) (int64, error)
	IdentMessage(ctx context.Context, messageId int64) (models.Message, error)
	AllMessages(ctx context.Context, userFrom, userTo int64) ([]int64, error)
}

type App struct {
	log            *slog.Logger
	userIdent      UserIdent
	workerContacts ContactsWorker
	workerMessage  MessageWorker
}

func New(log *slog.Logger, ident UserIdent, contacts ContactsWorker, messages MessageWorker) *App {
	return &App{
		log:            log,
		userIdent:      ident,
		workerContacts: contacts,
		workerMessage:  messages,
	}
}

func (a *App) IdentUser(ctx context.Context, id int64) (models.User, error) {
	const op = "service.auth.IdentUser"

	log := a.log.With(slog.String("op", op))
	log.Info("Identifying user")

	user, err := a.userIdent.IdentUser(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("User not found", err)
			return models.User{}, ErrUserNotFound
		}

		log.Error("Failed to get user", slog.String("error", err.Error()))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (a *App) AddContact(ctx context.Context, id, contactId int64) error {
	const op = "service.chat.AddContact"
	log := a.log.With(slog.String("op", op), slog.Int64("userId", id))
	log.Info("Adding contact")

	err := a.workerContacts.AddContact(ctx, id, contactId)
	if err != nil {
		log.Error("Failed to add contact", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) AllContacts(ctx context.Context, id int64) ([]int64, error) {
	const op = "service.chat.AllContacts"
	log := a.log.With(slog.String("op", op), slog.Int64("id", id))
	log.Info("Getting all contacts")

	contacts, err := a.workerContacts.AllContacts(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("No users", err)
			return nil, ErrUserNotFound
		}

		log.Error("Failed to get all contacts", err)
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return contacts, nil
}

func (a *App) IsMessaged(ctx context.Context, id, contactId int64) error {
	const op = "service.chat.IsMessaged"
	log := a.log.With(slog.String("op", op))
	log.Info("Trying message " + strconv.Itoa(int(id)) + " to " + strconv.Itoa(int(contactId)))

	err := a.workerContacts.IsMessaged(ctx, id, contactId)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("Failed to messaged", slog.String("error", err.Error()))
			return ErrUserNotFound
		}

		log.Error("Failed to messaged", slog.String("error", err.Error()))
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (a *App) AllMessaged(ctx context.Context, id int64) ([]int64, error) {
	const op = "service.chat.AllMessaged"
	log := a.log.With(slog.String("op", op))
	log.Info("Trying to get all messaged users of " + strconv.Itoa(int(id)))

	users, err := a.workerContacts.AllMessaged(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("Failed to get all messaged", slog.String("error", err.Error()))
			return nil, ErrUserNotFound
		}

		log.Error("Failed to get all messaged", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return users, nil
}

func (a *App) CreateMessage(ctx context.Context, text, date string, userId, contactId int64) (int64, error) {
	const op = "service.chat.CreateMessage"
	log := a.log.With(slog.String("op", op))
	log.Info("Trying to create message")

	messageId, err := a.workerMessage.CreateMessage(ctx, text, date, userId, contactId)
	if err != nil {
		log.Error("Failed to save message", slog.String("error", err.Error()))
		return -1, err
	}

	return messageId, nil
}

func (a *App) IdentMessage(ctx context.Context, messageId int64) (models.Message, error) {
	const op = "service.chat.IdentMessage"
	log := a.log.With(slog.String("op", op))
	log.Info("Ident message")

	msg, err := a.workerMessage.IdentMessage(ctx, messageId)
	if err != nil {
		if errors.Is(err, storage.ErrMessageNotFound) {
			return models.Message{}, ErrMessageNotFound
		}

		return models.Message{}, fmt.Errorf("%s:%w", op, err)
	}

	return msg, nil
}

func (a *App) AllMessages(ctx context.Context, userFrom, userTo int64) ([]int64, error) {
	const op = "service.chat.AllMessages"
	log := a.log.With(slog.String("op", op))
	log.Info("Getting all messages")

	ids, err := a.workerMessage.AllMessages(ctx, userFrom, userTo)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("User not found", slog.String("error", err.Error()))
			return nil, ErrUserNotFound
		}

		log.Error("Failed to get all messages", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return ids, nil
}
