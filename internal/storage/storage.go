package storage

import (
	"URLShortener/models"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url exists")
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) SaveURL(urlToSave string, alias string) (*int64, error) {
	const op = "storage.saveURL"

	qRes, err := s.db.Prepare(`INSERT INTO url (fullurl, alias) VALUES ($1, $2) RETURNING id`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer qRes.Close()

	var id *int64
	err = qRes.QueryRow(urlToSave, alias).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (*models.URL, error) {
	const op = "storage.getURL"

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

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.deleteURL"

	qRes, err := s.db.Exec(`DELETE FROM url WHERE alias=$1`, alias)
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
	const op = "storage.getAliasCheck"

	var duplicatedAlias string
	qRes := s.db.QueryRow(`SELECT * FROM url WHERE alias=$1`, alias)
	err := qRes.Scan(&duplicatedAlias)
	if err != sql.ErrNoRows {
		return fmt.Errorf("%s: %w ", op, err)
	}

	return nil
}
