package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID int `json:"id"`
	Email string `json:"email"`
	Password string `json:"password,omitempty"`
	EncryptedPassword string `json:"-"`
}

func (user *User) Validate() error {
	return validation.ValidateStruct(
		user,
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Password, validation.By(requiredIf(user.EncryptedPassword=="")), validation.Length(6, 100)))
}

func (user *User) BeforeCreate() error {
	if len(user.Password) > 0 {
		enc, err := encryptedString(user.Password)
		if err != nil {
			return err
		}

		user.EncryptedPassword = enc
	}
	return nil
}

func encryptedString(s string) (string, error)  {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (user *User) Sanitize() {
	user.Password = ""
}

func (user *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(password)) == nil
}