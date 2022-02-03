package store

import "github.com/ivanmalyi/RestApi/internal/app/model"

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
}
