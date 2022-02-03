package sqlstore_test

import (
	"github.com/ivanmalyi/RestApi/internal/app/model"
	"github.com/ivanmalyi/RestApi/internal/app/store"
	"github.com/ivanmalyi/RestApi/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseUrl)
	defer teardown("users")
	s := sqlstore.New(db)

	user := model.TestUser(t)
	err := s.User().Create(user)

	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseUrl)
	defer teardown("users")
	s := sqlstore.New(db)

	email := "ivan.malyi93@gmail.com"
	_, err := s.User().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	_ = s.User().Create(model.TestUser(t))
	user, err := s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, user)


}