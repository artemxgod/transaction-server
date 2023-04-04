package sqlstore

import (
	"database/sql"

	"github.com/artemxgod/transaction-server/internal/app/store"
	_ "github.com/lib/pq"
)

type Store struct {
	db             *sql.DB
	userRepository *UserRepository
}

func New(p_db *sql.DB) *Store {
	return &Store{
		db: p_db,
	}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
