package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Smile8MrBread/Chat/app/internal/models"
	"github.com/Smile8MrBread/Chat/app/internal/storage"
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

func (s *Storage) IdentUser(ctx context.Context, id int64) (models.User, error) {
	const op = "storage.sqlite.IdentUser"

	stmt, err := s.db.Prepare("SELECT id, first_name, last_name, avatar FROM Users WHERE id = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s:%w", op, err)
	}

	row := stmt.QueryRowContext(ctx, id)

	var user models.User
	if err = row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Avatar); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
	}

	return user, nil
}

func (s *Storage) AddContact(ctx context.Context, id, contactId int64) error {
	const op = "storage.sqlite.AddContact"

	rows, err := s.db.QueryContext(ctx, "SELECT contact_id FROM Contacts WHERE user_id = ?", id)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	for rows.Next() {
		var n int64
		if err := rows.Scan(&n); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				break
			}
		}
		if n == contactId {
			return fmt.Errorf("%s:%w", op, storage.ErrUserExists)
		}
	}

	stmt, err := s.db.Prepare("INSERT INTO Contacts(user_id, contact_id) VALUES(?, ?)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, id, contactId)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (s *Storage) AllContacts(ctx context.Context, id int64) ([]int64, error) {
	const op = "storage.sqlite.AllContacts"

	rows, err := s.db.QueryContext(ctx, "SELECT contact_id FROM Contacts WHERE user_id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	var contactsId []int64
	for rows.Next() {
		var contactId int64
		if err = rows.Scan(&contactId); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
			}

			return nil, fmt.Errorf("%s: %w", op, err)
		}
		contactsId = append(contactsId, contactId)
	}

	if len(contactsId) == 0 {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
	}

	return contactsId, nil
}

func (s *Storage) IsMessaged(ctx context.Context, id, contactId int64) error {
	const op = "storage.sqlite.IsMessaged"

	stmt, err := s.db.PrepareContext(ctx, "UPDATE Contacts SET 'is_messaged' = 1"+
		" WHERE user_id = ? AND contact_id = ?")
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	if _, err = stmt.ExecContext(ctx, id, contactId); err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrNotFound) {
			return fmt.Errorf("%s:%w", op, storage.ErrUserNotFound)
		}

		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (s *Storage) AllMessaged(ctx context.Context, id int64) ([]int64, error) {
	const op = "storage.sqlite.AllMessaged"

	rows, err := s.db.QueryContext(ctx, "SELECT user_id FROM Contacts WHERE contact_id = ? AND is_messaged = 1", id)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	occured := map[int64]bool{}
	var users []int64

	for rows.Next() {
		var n int64
		if err = rows.Scan(&n); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
			}

			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if !occured[n] {
			users = append(users, n)
			occured[n] = true
		}
	}

	rows, err = s.db.QueryContext(ctx, "SELECT contact_id FROM Contacts WHERE user_id = ? AND is_messaged = 1", id)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	for rows.Next() {
		var n int64
		if err = rows.Scan(&n); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
			}

			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if !occured[n] {
			users = append(users, n)
			occured[n] = true
		}
	}

	return users, nil
}

func (s *Storage) CreateMessage(ctx context.Context, text, date string, userId, contactId int64) (int64, error) {
	const op = "storage.sqlite.CreateMessage"

	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO Messages(text, date, user_from, user_to) VALUES(?, ?, ?, ?)")
	if err != nil {
		return -1, fmt.Errorf("%s:%w", op, err)
	}

	res, err := stmt.ExecContext(ctx, text, date, userId, contactId)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%s:%w", op, err)
	}

	return id, nil
}

func (s *Storage) IdentMessage(ctx context.Context, messageId int64) (models.Message, error) {
	const op = "storage.sqlite.IdentMessage"

	stmt, err := s.db.PrepareContext(ctx, "SELECT * FROM Messages WHERE id = ?")
	if err != nil {
		return models.Message{}, fmt.Errorf("%s:%w", op, err)
	}

	row := stmt.QueryRowContext(ctx, messageId)

	var msg models.Message
	if err = row.Scan(&msg.Id, &msg.Text, &msg.Date, &msg.UserFrom, &msg.UserTo); err != nil {
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) && errors.Is(err, sql.ErrNoRows) {
			return models.Message{}, storage.ErrMessageNotFound
		}

		return models.Message{}, fmt.Errorf("%s:%w", op, err)
	}

	return msg, nil
}

func (s *Storage) AllMessages(ctx context.Context, userFrom, userTo int64) ([]int64, error) {
	const op = "storage.sqlite.AllMessages"

	rows, err := s.db.QueryContext(ctx, "SELECT id FROM Messages WHERE (user_from = ? AND user_to = ?) OR (user_to = ? AND user_from = ?)",
		userFrom, userTo, userFrom, userTo)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	var messages []int64
	for rows.Next() {
		var n int64
		if err = rows.Scan(&n); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
			}

			return nil, fmt.Errorf("%s: %w", op, err)
		}

		messages = append(messages, n)
	}

	return messages, nil
}
