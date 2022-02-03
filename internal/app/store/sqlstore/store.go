package sqlstore

import (
	"database/sql"
	"github.com/ivanmalyi/RestApi/internal/app/store"
	_"github.com/lib/pq"
)

type Store struct {
	db *sql.DB
	userRepository *UserRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (store *Store) User() store.UserRepository {
	if store.userRepository != nil {
		return store.userRepository
	}

	store.userRepository = &UserRepository {
		store: store,
	}

	return store.userRepository
}