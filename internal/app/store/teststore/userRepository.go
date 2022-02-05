package teststore

import (
	"github.com/ivanmalyi/RestApi/internal/app/model"
	"github.com/ivanmalyi/RestApi/internal/app/store"
)

type UserRepository struct {
	store *Store
	users map[int]*model.User
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

	user.ID = len(userRepository.users) + 1
	userRepository.users[user.ID] = user

	return err
}

func (userRepository *UserRepository) FindByEmail(email string) (*model.User, error)  {
	for _,user := range userRepository.users {
		if email == user.Email {
			return user, nil
		}
	}

	return nil, store.ErrRecordNotFound
}

func (userRepository *UserRepository) Find(id int) (*model.User, error)  {
	user, ok := userRepository.users[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return user, nil
}