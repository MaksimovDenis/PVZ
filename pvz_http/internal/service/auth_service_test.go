package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateData(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		password    string
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "empty email",
			email:       "",
			password:    "pass123",
			wantErr:     true,
			expectedErr: errors.New("укажите почту"),
		},
		{
			name:        "empty password",
			email:       "user@mail.com",
			password:    "",
			wantErr:     true,
			expectedErr: errors.New("укажите пароль"),
		},
		{
			name:        "email equals password",
			email:       "samevalue",
			password:    "samevalue",
			wantErr:     true,
			expectedErr: errors.New("почта и пароль совпадают"),
		},
		{
			name:     "valid input",
			email:    "admin@mail.ru",
			password: "securepassword",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateData(tt.email, tt.password)
			if tt.wantErr {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateRole(t *testing.T) {
	tests := []struct {
		name        string
		role        string
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "invalid role - admin",
			role:        "admin",
			wantErr:     true,
			expectedErr: errors.New("Неверный формат роли пользователя"),
		},
		{
			name:        "empty role",
			role:        "",
			wantErr:     true,
			expectedErr: errors.New("Неверный формат роли пользователя"),
		},
		{
			name:    "valid employee",
			role:    "employee",
			wantErr: false,
		},
		{
			name:    "valid moderator",
			role:    "moderator",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRole(tt.role)
			if tt.wantErr {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
