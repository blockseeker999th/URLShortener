package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/blockseeker999th/URLShortener/internal/config"

	_ "github.com/lib/pq"
)

type PostgreSQLStorage struct {
	db *sql.DB
}

func ConnectDB(config *config.Config) *PostgreSQLStorage {
	const op = "storage.db.ConnectDB"

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	fmt.Println("connectionString: ", connectionString)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("error path: %s", op)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("error pinging PostgreSQL: %v", err)
	}

	return &PostgreSQLStorage{db: db}
}

func (s *PostgreSQLStorage) InitNewPostgreSQLStorage() (*sql.DB, error) {
	if err := s.createURLTable(); err != nil {
		return nil, err
	}

	return s.db, nil
}

func (s *PostgreSQLStorage) createURLTable() error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS url (
    id SERIAL PRIMARY KEY,
    alias VARCHAR(255) NOT NULL UNIQUE,
    url VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);
`)

	if err != nil {
		return errors.New(err.Error())
	}

	_, err = s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);`)

	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}
