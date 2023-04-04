package sqlstore_test

import (
	"testing"

	"github.com/artemxgod/transaction-server/internal/app/model"
	"github.com/artemxgod/transaction-server/internal/app/store"
	"github.com/artemxgod/transaction-server/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")
	s := sqlstore.New(db)
	u := model.TestUser(t) 

	err := s.User().CreateRecord(u)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_FindByName(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")
	s := sqlstore.New(db)

	s.User().CreateRecord(model.TestUser(t))


	name := "Alex111"
	u, err := s.User().FindByName(name)

	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_Find(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)
	u1 := model.TestUser(t)
	s.User().CreateRecord(u1)


	u, err := s.User().Find(u1.ID)

	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail_non(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")
	s := sqlstore.New(db)

	name := "Nemo"
	_, err := s.User().FindByName(name)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

}