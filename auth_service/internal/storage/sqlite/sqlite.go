package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Smile8MrBread/Chat/auth_service/internal/models"
	"github.com/Smile8MrBread/Chat/auth_service/internal/storage"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

func (s *Storage) SaveUser(ctx context.Context, firstName, lastName, login, avatar string,
	passHash []byte) (int64, error) {
	const op = "storage.sqlite.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO Users(first_name, last_name, login, avatar, password_hash) " +
		"VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return -1, fmt.Errorf("%s:%w", op, err)
	}

	res, err := stmt.ExecContext(ctx, firstName, lastName, login, avatar, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return -1, fmt.Errorf("%s:%w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%s:%w", op, err)
	}

	return id, nil
}

func (s *Storage) ProvideUser(ctx context.Context, login string) (models.User, error) {
	const op = "storage.sqlite.ProvideUser"

	stmt, err := s.db.Prepare("SELECT * FROM Users WHERE login = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s:%w", op, err)
	}

	row := stmt.QueryRowContext(ctx, login)

	var user models.User
	if err = row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Login, &user.Avatar, &user.PassHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
	}

	return user, nil
}
