package teststore

import (
	"github.com/ivanmalyi/RestApi/internal/app/model"
	"github.com/ivanmalyi/RestApi/internal/app/store"
)

type UserRepository struct {
	store *Store
	users map[string]*model.User
}

func (userRepository *UserRepository) Create(user *model.User) error  {
	var err error
	err = user.Validate()
	if err != nil {
		return err
	}

	err = user.BeforeCreate()
	if err != nil {
		return err
	}

	userRepository.users[user.Email] = user
	user.ID = len(userRepository.users)

	return err
}

func (userRepository *UserRepository) FindByEmail(email string) (*model.User, error)  {
	user, ok := userRepository.users[email]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return user, nil
}