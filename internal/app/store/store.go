package store

import (
	"database/sql"
	_"github.com/lib/pq"
)

type Store struct {
	config *Config
	db *sql.DB
}

func New(config *Config) *Store  {
	return &Store{
		config: config,
	}
}

func (store *Store) Open() error {
	db, err := sql.Open("postgres", store.config.DatabaseUrl)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	store.db = db

	return nil
}

func (store *Store)Close() {
	_ = store.db.Close()
}