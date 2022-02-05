package teststore_test

import (
	"github.com/ivanmalyi/RestApi/internal/app/model"
	"github.com/ivanmalyi/RestApi/internal/app/store"
	"github.com/ivanmalyi/RestApi/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	s := teststore.New()

	user := model.TestUser(t)
	err := s.User().Create(user)

	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	s := teststore.New()

	email := "ivan.malyi93@gmail.com"
	_, err := s.User().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	_ = s.User().Create(model.TestUser(t))
	user, err := s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestUserRepository_Find(t *testing.T) {
	s := teststore.New()
	id := 1
	_, err := s.User().Find(id)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	testUser := model.TestUser(t)
	_ = s.User().Create(testUser)
	user, err := s.User().Find(testUser.ID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}