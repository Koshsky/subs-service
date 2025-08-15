package models

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestUserValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name    string
		user    User
		isValid bool
	}{
		{
			name: "Valid user",
			user: User{
				Email:    "test@example.com",
				Password: "password123",
			},
			isValid: true,
		},
		{
			name: "Invalid email format - missing @",
			user: User{
				Email:    "invalid-email",
				Password: "password123",
			},
			isValid: false,
		},
		{
			name: "Invalid email format - missing domain",
			user: User{
				Email:    "test@",
				Password: "password123",
			},
			isValid: false,
		},
		{
			name: "Invalid email format - missing local part",
			user: User{
				Email:    "@example.com",
				Password: "password123",
			},
			isValid: false,
		},
		{
			name: "Invalid email format - spaces",
			user: User{
				Email:    "test @example.com",
				Password: "password123",
			},
			isValid: false,
		},
		{
			name: "Empty email",
			user: User{
				Email:    "",
				Password: "password123",
			},
			isValid: false,
		},
		{
			name: "Empty password",
			user: User{
				Email:    "test@example.com",
				Password: "",
			},
			isValid: false,
		},
		{
			name: "Password too short",
			user: User{
				Email:    "test@example.com",
				Password: "123",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.user)
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
