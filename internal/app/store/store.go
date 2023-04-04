package store

import "github.com/artemxgod/transaction-server/internal/app/model"

type Store interface {
	User() UserRepository
}

type UserRepository interface {
	CreateRecord(*model.User) error
	FindByName(name string) (*model.User, error)
	Find(id int) (*model.User, error)
	ChangeBalance(id int, funds float64) (*model.User, error)
}
