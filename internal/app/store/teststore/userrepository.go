package teststore

import (
	"github.com/artemxgod/transaction-server/internal/app/model"
	"github.com/artemxgod/transaction-server/internal/app/store"
)

type UserRepository struct {
	store *Store
	users map[int]*model.User
}

/* Creates new record in database, returns modified User struct that contains ID
Requires user struct with email and password*/
func (r *UserRepository) CreateRecord(u *model.User) error {
	u.ID = len(r.users) + 1
	r.users[u.ID] = u

	return nil
}

func (r *UserRepository) Find(ID int) (*model.User, error) {
	u, ok := r.users[ID]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return u, nil
}

func (r *UserRepository) FindByName(name string) (*model.User, error) {

	for _, u := range r.users {
		if u.Name == name {
			return u, nil
		}
	}
		
	return nil, store.ErrRecordNotFound
}

// TODO?
func (r* UserRepository) ChangeBalance(id int, funds float64) (*model.User, error) {
	return nil, nil
}