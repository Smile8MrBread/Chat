package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Smile8MrBread/Chat/chat_service/internal/models"
	"github.com/Smile8MrBread/Chat/chat_service/internal/storage"
	"log/slog"
	"net/http"
	"strconv"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrMessageNotFound = errors.New("message not found")
)

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
	workerContacts ContactsWorker
	workerMessage  MessageWorker
}

func New(log *slog.Logger, contacts ContactsWorker, messages MessageWorker) *App {
	return &App{
		log:            log,
		workerContacts: contacts,
		workerMessage:  messages,
	}
}

func (a *App) AddContact(ctx context.Context, id, contactId int64) error {
	const op = "service.chat.AddContact"
	log := a.log.With(slog.String("op", op), slog.Int64("userId", id))
	log.Info("Adding contact")

	req, err := http.NewRequest(http.MethodGet, "http://client:8080/identity/"+strconv.Itoa(int(contactId)), nil)
	if err != nil {
		log.Error("Failed to add contact", slog.String("error", err.Error()))
		return fmt.Errorf("%s:%w", op, err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Failed to add contact", slog.String("error", err.Error()))
		return fmt.Errorf("%s:%w", op, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 409 {
		log.Error("User not found")
		return fmt.Errorf("%s:%w", op, ErrUserNotFound)
	}
	err = a.workerContacts.AddContact(ctx, id, contactId)
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
			return nil, fmt.Errorf("%s:%w", op, ErrUserNotFound)
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
			return fmt.Errorf("%s:%w", op, ErrUserNotFound)
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
			return nil, fmt.Errorf("%s:%w", op, ErrUserNotFound)
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
			return nil, fmt.Errorf("%s:%w", op, ErrUserNotFound)
		}

		log.Error("Failed to get all messages", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return ids, nil
}
