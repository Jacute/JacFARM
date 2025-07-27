package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const dbFilepath = "./database.db"

type Storage struct {
	db *sql.DB
}

func New() (*Storage, error) {
	db, err := sql.Open("sqlite3", dbFilepath)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Storage{db}, nil
}

func (s *Storage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
