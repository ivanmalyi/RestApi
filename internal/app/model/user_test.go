package model_test

import (
	"github.com/ivanmalyi/RestApi/internal/app/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_BeforeCreate(t *testing.T) {
	user := model.TestUser(t)
	assert.NoError(t, user.BeforeCreate())
	assert.NotEmpty(t, user.EncryptedPassword)
}

func TestUser_Validate(t *testing.T) {
	testCases := []struct{
		name string
		genUser func() *model.User
		isValid bool
	}{
		{
			name: "valid",
			genUser: func() *model.User {
				return model.TestUser(t)
			},
			isValid: true,
		},
		{
			name: "with encrypted password",
			genUser: func() *model.User {
				user := model.TestUser(t)
				user.Password = ""
				user.EncryptedPassword = "encrypted_password"
				return user
			},
			isValid: true,
		},
		{
			name: "empty email",
			genUser: func() *model.User {
				user := model.TestUser(t)
				user.Email = ""
				return user
			},
			isValid: false,
		},
		{
			name: "invalid email",
			genUser: func() *model.User {
				user := model.TestUser(t)
				user.Email = "invalid"
				return user
			},
			isValid: false,
		},
		{
			name: "empty password",
			genUser: func() *model.User {
				user := model.TestUser(t)
				user.Password = ""
				return user
			},
			isValid: false,
		},
		{
			name: "short password",
			genUser: func() *model.User {
				user := model.TestUser(t)
				user.Password = "123"
				return user
			},
			isValid: false,
		},
	}

	for _, tc:=range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.genUser().Validate())
			} else {
				assert.Error(t, tc.genUser().Validate())
			}
		})
	}

	user := model.TestUser(t)
	assert.NoError(t, user.Validate())
}