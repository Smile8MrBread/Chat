package grpc

import (
	"context"
	"errors"
	"github.com/Smile8MrBread/Chat/chat_service/internal/models"
	"github.com/Smile8MrBread/Chat/chat_service/internal/service"
	"github.com/Smile8MrBread/Chat/chat_service/internal/storage"
	chatGrpc "github.com/Smile8MrBread/Chat/chat_service/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"reflect"
	"strconv"
)

type Chat interface {
	AddContact(ctx context.Context, id, contactId int64) error
	AllContacts(ctx context.Context, id int64) ([]int64, error)
	IsMessaged(ctx context.Context, id, contactId int64) error
	AllMessaged(ctx context.Context, id int64) ([]int64, error)
	CreateMessage(ctx context.Context, text, date string, userId, contactId int64) (int64, error)
	IdentMessage(ctx context.Context, messageId int64) (models.Message, error)
	AllMessages(ctx context.Context, userFrom, userTo int64) ([]int64, error)
}

type ServerAPI struct {
	chatGrpc.UnimplementedChatServer
	chat Chat
}

func Register(gRPC *grpc.Server, chatStr Chat) {
	chatGrpc.RegisterChatServer(gRPC, &ServerAPI{
		chat: chatStr,
	})
}

func (s *ServerAPI) AddContact(ctx context.Context, req *chatGrpc.AddContactRequest) (*chatGrpc.Nothing, error) {
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

	return &chatGrpc.Nothing{}, nil
}

func (s *ServerAPI) AllContacts(ctx context.Context, req *chatGrpc.AllContactsRequest) (*chatGrpc.AllContactsResponse, error) {
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

	return &chatGrpc.AllContactsResponse{ContactIds: contacts}, nil
}

func (s *ServerAPI) IsMessaged(ctx context.Context, req *chatGrpc.IsMessagedRequest) (*chatGrpc.Nothing, error) {
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

	return &chatGrpc.Nothing{}, nil
}

func (s *ServerAPI) AllMessaged(ctx context.Context, req *chatGrpc.AllMessagedRequest) (*chatGrpc.AllMessagedResponse, error) {
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

	return &chatGrpc.AllMessagedResponse{UserIds: users}, nil
}

func (s *ServerAPI) IdentMessage(ctx context.Context, req *chatGrpc.IdentMessageRequest) (*chatGrpc.IdentMessageResponse, error) {
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

	return &chatGrpc.IdentMessageResponse{
		MessageId: msg.Id,
		Text:      msg.Text,
		Date:      msg.Date,
		UserFrom:  int64(userFrom),
		UserTo:    int64(userTo),
	}, nil
}

func (s *ServerAPI) AllMessages(ctx context.Context, req *chatGrpc.AllMessagesRequest) (*chatGrpc.AllMessagesResponse, error) {
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

	return &chatGrpc.AllMessagesResponse{MessageIds: ids}, nil
}

func validateAddContact(req *chatGrpc.AddContactRequest) error {
	if reflect.TypeOf(req.GetUserId()).Kind() != reflect.Int64 ||
		reflect.TypeOf(req.GetContactId()).Kind() != reflect.Int64 {
		return status.Error(codes.InvalidArgument, "Invalid id")
	}

	return nil
}

func validateAllContacts(req *chatGrpc.AllContactsRequest) error {
	if reflect.TypeOf(req.GetUserId()).Kind() != reflect.Int64 {
		return status.Error(codes.InvalidArgument, "Invalid id")
	}

	return nil
}

func validateIsMessaged(req *chatGrpc.IsMessagedRequest) error {
	if reflect.TypeOf(req.GetUserId()).Kind() != reflect.Int64 ||
		reflect.TypeOf(req.GetContactId()).Kind() != reflect.Int64 {
		return status.Error(codes.InvalidArgument, "Invalid id")
	}

	return nil
}

func validateAllMessaged(req *chatGrpc.AllMessagedRequest) error {
	if reflect.TypeOf(req.GetUserId()).Kind() != reflect.Int64 {
		return status.Error(codes.InvalidArgument, "Invalid id")
	}

	return nil
}

func validateIdentMessage(req *chatGrpc.IdentMessageRequest) error {
	if reflect.TypeOf(req.GetMessageId()).Kind() != reflect.Int64 {
		return status.Error(codes.InvalidArgument, "Invalid message")
	}

	return nil
}

func validateAllMessages(req *chatGrpc.AllMessagesRequest) error {
	if reflect.TypeOf(req.GetUserFrom()).Kind() != reflect.Int64 &&
		reflect.TypeOf(req.GetUserTo()).Kind() != reflect.Int64 {
		return status.Error(codes.InvalidArgument, "Invalid id")
	}

	return nil
}
