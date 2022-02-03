package sqlstore

import (
	"database/sql"
	"github.com/ivanmalyi/RestApi/internal/app/model"
	"github.com/ivanmalyi/RestApi/internal/app/store"
)

type UserRepository struct {
	store *Store
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

	return userRepository.store.db.QueryRow(
		`insert into users (email, encrypted_password) values ($1, $2) RETURNING id`,
    	user.Email, user.EncryptedPassword,
		).Scan(&user.ID)
}

func (userRepository *UserRepository) FindByEmail(email string) (*model.User, error)  {
	user := &model.User{}

	err := userRepository.store.db.QueryRow(
		`select id, email, encrypted_password 
		from users
		where email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.EncryptedPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return user, nil
}