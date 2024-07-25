package grpc

import (
	"context"
	"errors"
	"github.com/Smile8MrBread/Chat/chat_service/internal/models"
	"github.com/Smile8MrBread/Chat/chat_service/internal/service"
	"github.com/Smile8MrBread/Chat/chat_service/internal/storage"
	"github.com/Smile8MrBread/Chat/chat_service/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"reflect"
	"strconv"
)

type Chat interface {
	IdentUser(ctx context.Context, id int64) (models.User, error)
	AddContact(ctx context.Context, id, contactId int64) error
	AllContacts(ctx context.Context, id int64) ([]int64, error)
	IsMessaged(ctx context.Context, id, contactId int64) error
	AllMessaged(ctx context.Context, id int64) ([]int64, error)
	CreateMessage(ctx context.Context, text, date string, userId, contactId int64) (int64, error)
	IdentMessage(ctx context.Context, messageId int64) (models.Message, error)
	AllMessages(ctx context.Context, userFrom, userTo int64) ([]int64, error)
}

type ServerAPI struct {
	gen.UnimplementedChatServer
	chat Chat
}

func Register(gRPC *grpc.Server, chatStr Chat) {
	gen.RegisterChatServer(gRPC, &ServerAPI{
		chat: chatStr,
	})
}

func (s *ServerAPI) IdentUser(ctx context.Context, req *gen.IdentUserRequest) (*gen.IdentUserResponse, error) {
	if err := validateIdentUser(req); err != nil {
		return nil, err
	}

	user, err := s.chat.IdentUser(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return nil, status.Error(codes.InvalidArgument, "User not found")
		}

		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &gen.IdentUserResponse{
		UserId:    user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
	}, nil
}

func (s *ServerAPI) AddContact(ctx context.Context, req *gen.AddContactRequest) (*gen.Nothing, error) {
	if err := validateAddContact(req); err != nil {
		return nil, err
	}

	err := s.chat.AddContact(ctx, req.GetUserId(), req.GetContactId())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.InvalidArgument, "User already added")
		}
		if errors.Is(err, service.ErrUserNotFound) {
			return nil, status.Error(codes.InvalidArgument, "User not found")
		}

		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &gen.Nothing{}, nil
}

func (s *ServerAPI) AllContacts(ctx context.Context, req *gen.AllContactsRequest) (*gen.AllContactsResponse, error) {
	if err := validateAllContacts(req); err != nil {
		return nil, err
	}

	contacts, err := s.chat.AllContacts(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "No users in contacts")
		}
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &gen.AllContactsResponse{ContactIds: contacts}, nil
}

func (s *ServerAPI) IsMessaged(ctx context.Context, req *gen.IsMessagedRequest) (*gen.Nothing, error) {
	if err := validateIsMessaged(req); err != nil {
		return nil, err
	}

	err := s.chat.IsMessaged(ctx, req.GetUserId(), req.GetContactId())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &gen.Nothing{}, nil
}

func (s *ServerAPI) AllMessaged(ctx context.Context, req *gen.AllMessagedRequest) (*gen.AllMessagedResponse, error) {
	if err := validateAllMessaged(req); err != nil {
		return nil, err
	}

	users, err := s.chat.AllMessaged(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &gen.AllMessagedResponse{UserIds: users}, nil
}

func (s *ServerAPI) CreateMessage(ctx context.Context, req *gen.CreateMessageRequest) (*gen.CreateMessageResponse, error) {
	if err := validateCreateMessage(req); err != nil {
		return nil, err
	}

	id, err := s.chat.CreateMessage(ctx, req.GetText(), req.GetDate(), req.GetUserFrom(), req.GetUserTo())
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &gen.CreateMessageResponse{MessageId: id}, nil
}

func (s *ServerAPI) IdentMessage(ctx context.Context, req *gen.IdentMessageRequest) (*gen.IdentMessageResponse, error) {
	if err := validateIdentMessage(req); err != nil {
		return nil, err
	}

	msg, err := s.chat.IdentMessage(ctx, req.GetMessageId())
	if err != nil {
		if errors.Is(err, storage.ErrMessageNotFound) {
			return nil, status.Error(codes.InvalidArgument, "Message not found")
		}

		return nil, status.Error(codes.Internal, "Internal error")
	}

	userFrom, err := strconv.Atoi(msg.UserFrom)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	userTo, err := strconv.Atoi(msg.UserTo)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &gen.IdentMessageResponse{
		MessageId: msg.Id,
		Text:      msg.Text,
		Date:      msg.Date,
		UserFrom:  int64(userFrom),
		UserTo:    int64(userTo),
	}, nil
}

func (s *ServerAPI) AllMessages(ctx context.Context, req *gen.AllMessagesRequest) (*gen.AllMessagesResponse, error) {
	if err := validateAllMessages(req); err != nil {
		return nil, err
	}

	ids, err := s.chat.AllMessages(ctx, req.GetUserFrom(), req.GetUserTo())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "User not found")
		}

		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &gen.AllMessagesResponse{MessageIds: ids}, nil
}

func validateIdentUser(req *gen.IdentUserRequest) error {
	if reflect.TypeOf(req.GetUserId()).Kind() != reflect.Int64 {
		return status.Error(codes.InvalidArgument, "Invalid id")
	}

	return nil
}

func validateAddContact(req *gen.AddContactRequest) error {
	if reflect.TypeOf(req.GetUserId()).Kind() != reflect.Int64 ||
		reflect.TypeOf(req.GetContactId()).Kind() != reflect.Int64 {
		return status.Error(codes.InvalidArgument, "Invalid id")
	}

	return nil
}

func validateAllContacts(req *gen.AllContactsRequest) error {
	if reflect.TypeOf(req.GetUserId()).Kind() != reflect.Int64 {
		return status.Error(codes.InvalidArgument, "Invalid id")
	}

	return nil
}

func validateIsMessaged(req *gen.IsMessagedRequest) error {
	if reflect.TypeOf(req.GetUserId()).Kind() != reflect.Int64 ||
		reflect.TypeOf(req.GetContactId()).Kind() != reflect.Int64 {
		return status.Error(codes.InvalidArgument, "Invalid id")
	}

	return nil
}

func validateAllMessaged(req *gen.AllMessagedRequest) error {
	if reflect.TypeOf(req.GetUserId()).Kind() != reflect.Int64 {
		return status.Error(codes.InvalidArgument, "Invalid id")
	}

	return nil
}

func validateCreateMessage(req *gen.CreateMessageRequest) error {
	if reflect.TypeOf(req.GetText()).Kind() != reflect.String &&
		reflect.TypeOf(req.GetDate()).Kind() != reflect.String &&
		reflect.TypeOf(req.GetUserTo()).Kind() != reflect.Int64 &&
		reflect.TypeOf(req.GetUserFrom()).Kind() != reflect.Int64 {
		return status.Error(codes.InvalidArgument, "Invalid message")
	}

	return nil
}

func validateIdentMessage(req *gen.IdentMessageRequest) error {
	if reflect.TypeOf(req.GetMessageId()).Kind() != reflect.Int64 {
		return status.Error(codes.InvalidArgument, "Invalid message")
	}

	return nil
}

func validateAllMessages(req *gen.AllMessagesRequest) error {
	if reflect.TypeOf(req.GetUserFrom()).Kind() != reflect.Int64 &&
		reflect.TypeOf(req.GetUserTo()).Kind() != reflect.Int64 {
		return status.Error(codes.InvalidArgument, "Invalid id")
	}

	return nil
}
