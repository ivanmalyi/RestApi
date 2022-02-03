package store

import "github.com/ivanmalyi/RestApi/internal/app/model"

type UserRepository struct {
	store *Store
}

func (userRepository *UserRepository) Create(user *model.User) (*model.User, error)  {
	var err error
	err = user.Validate()
	if err != nil {
		return nil, err
	}

	err = user.BeforeCreate()
	if err != nil {
		return nil, err
	}

	err = userRepository.store.db.QueryRow(
		`insert into users (email, encrypted_password) values ($1, $2) RETURNING id`,
    	user.Email, user.EncryptedPassword,
		).Scan(&user.ID)

	if err != nil {
		return nil, err
	}
	return user, nil
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
		return nil, err
	}
	return user, nil
}