package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateCity(t *testing.T) {
	tests := []struct {
		name      string
		city      string
		wantErr   bool
		errString string
	}{
		{
			name:    "Valid city Moscow",
			city:    "Москва",
			wantErr: false,
		},
		{
			name:    "Valid city Kazan",
			city:    "Казань",
			wantErr: false,
		},
		{
			name:    "Valid city Saint Petersburg",
			city:    "Санкт-Петербург",
			wantErr: false,
		},
		{
			name:      "Invalid city Novosibirsk",
			city:      "Новосибирск",
			wantErr:   true,
			errString: "в данном городе пока нет доступных ПВЗ",
		},
		{
			name:      "Invalid empty string",
			city:      "",
			wantErr:   true,
			errString: "в данном городе пока нет доступных ПВЗ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCity(tt.city)
			if tt.wantErr {
				require.Error(t, err)
				require.EqualError(t, err, tt.errString)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
