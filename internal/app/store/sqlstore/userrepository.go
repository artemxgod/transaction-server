package sqlstore

import (
	"database/sql"

	"github.com/artemxgod/transaction-server/internal/app/model"
	"github.com/artemxgod/transaction-server/internal/app/store"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) CreateRecord(u *model.User) error {
	return r.store.db.QueryRow(
		"INSERT INTO users (name, balance) VALUES ($1, $2) RETURNING id",
		u.Name, u.Balance,
	).Scan(&u.ID)
}

func (r *UserRepository) FindByName(name string) (*model.User, error) {
	u := &model.User{}

	if err := r.store.db.QueryRow("SELECT id, name, balance FROM users WHERE name = $1",
		name).Scan(&u.ID, &u.Name, &u.Balance); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
	}
	return u, nil
}

func (r *UserRepository) Find(id int) (*model.User, error) {
	u := &model.User{}

	if err := r.store.db.QueryRow("SELECT id, name, balance FROM users WHERE id = $1",
		id).Scan(&u.ID, &u.Name, &u.Balance); err != nil {
		return nil, store.ErrRecordNotFound
	}
	return u, nil
}

func (r *UserRepository) ChangeBalance(id int, funds float64) (*model.User, error) {
	u, err := r.Find(id)
	if err != nil {
		return nil, err
	}

	new_balance := u.Balance + funds

	if err := r.store.db.QueryRow("UPDATE users SET balance = $1 WHERE id = $2", 
		new_balance, id).Err(); err != nil {
			return nil, err
		}
	u.Balance = new_balance
	return u, nil
} 
