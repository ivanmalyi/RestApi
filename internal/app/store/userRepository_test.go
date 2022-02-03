package store_test

import (
	"github.com/ivanmalyi/RestApi/internal/app/model"
	"github.com/ivanmalyi/RestApi/internal/app/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	s, teardown := store.TestStore(t, databaseUrl)
	defer teardown("users")

	user, err := s.User().Create(model.TestUser(t))

	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	s, teardown := store.TestStore(t, databaseUrl)
	defer teardown("users")
	email := "ivan.malyi93@gmail.com"

	_, err := s.User().FindByEmail(email)
	assert.Error(t, err)

	_,_ = s.User().Create(model.TestUser(t))
	user, err := s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, user)


}