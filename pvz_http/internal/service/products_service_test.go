package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateProductType(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		expectedErr error
	}{
		{
			name:    "valid electronics",
			input:   "электроника",
			wantErr: false,
		},
		{
			name:    "valid clothing",
			input:   "одежда",
			wantErr: false,
		},
		{
			name:    "valid shoes",
			input:   "обувь",
			wantErr: false,
		},
		{
			name:        "invalid product",
			input:       "еда",
			wantErr:     true,
			expectedErr: errors.New("данный тип товара не поддерживается"),
		},
		{
			name:        "empty string",
			input:       "",
			wantErr:     true,
			expectedErr: errors.New("данный тип товара не поддерживается"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateProductType(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
