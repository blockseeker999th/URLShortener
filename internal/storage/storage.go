package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/blockseeker999th/URLShortener/models"
)

var (
	ErrURLNotFound        = errors.New("url not found")
	ErrURLExists          = errors.New("url exists")
	ErrSavingURL          = "error saving url"
	ErrMethodNotAllowed   = "method now allowed"
	ErrFailedToDecode     = "failed to decode request"
	ErrInvalidRequest     = "invalid request"
	ErrValidation         = "validation error"
	ErrSignUp             = "error registering a user"
	ErrCreatingSession    = "error creating session"
	ErrDeletingURL        = "error deleting URL"
	ErrInvalidCredentials = "invalid credentials"
	ErrFailedToGetURL     = "failed to get URL"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) SaveURL(urlToSave string, alias string, userId string) (*int64, error) {
	const op = "storage.SaveURL"

	qRes, err := s.db.Prepare(`INSERT INTO url (fullurl, alias, assignedToId) VALUES ($1, $2, $3) RETURNING id`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer qRes.Close()

	var id *int64
	err = qRes.QueryRow(urlToSave, alias, userId).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (*models.URL, error) {
	const op = "storage.GetURL"

	var u models.URL
	qRes, err := s.db.Prepare(`SELECT fullurl FROM url WHERE alias=$1`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer qRes.Close()

	err = qRes.QueryRow(alias).Scan(&u.Url)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrURLNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &u, nil
}

func (s *Storage) DeleteURL(alias string, userId string) error {
	const op = "storage.DeleteURL"

	qRes, err := s.db.Exec(`DELETE FROM url WHERE alias=$1 AND assignedToId=$2`, alias, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := qRes.RowsAffected()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return ErrURLNotFound
	}

	return nil
}

func (s *Storage) GetDuplicateAliasCheck(alias string) error {
	const op = "storage.GetAliasCheck"

	var duplicatedAlias string
	qRes := s.db.QueryRow(`SELECT * FROM url WHERE alias=$1`, alias)
	err := qRes.Scan(&duplicatedAlias)
	if err != sql.ErrNoRows {
		return fmt.Errorf("%s: %w ", op, err)
	}

	return nil
}

func (s *Storage) SignUpUser(user *models.User) (*models.User, error) {
	const op = "storage.SignUpUser"
	err := s.db.QueryRow("INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id", user.Username, user.Email, user.Password).Scan(&user.Id)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) SignInUser(loginData *models.LoginData) (*models.User, error) {
	const op = "storage.SignInUser"

	var authUser models.User
	err := s.db.QueryRow("SELECT id, email, password FROM users WHERE email=$1", loginData.Email).Scan(&authUser.Id, &authUser.Email, &authUser.Password)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &authUser, nil
}
