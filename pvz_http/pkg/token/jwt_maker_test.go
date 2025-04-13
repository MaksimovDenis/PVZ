package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const secretKey = "testsecret"

func TestJWTMaker(t *testing.T) {
	maker := NewJWTMaker(secretKey)

	uesrId, err := uuid.NewRandom()
	require.NoError(t, err)

	tests := []struct {
		name         string
		id           uuid.UUID
		email        string
		role         string
		duration     time.Duration
		shouldExpire bool
		shouldFail   bool
		signingNone  bool
	}{
		{
			name:     "Valid Token",
			id:       uesrId,
			email:    "test@mail.ru",
			duration: time.Minute,
		},
		{
			name:         "Expired Token",
			id:           uesrId,
			email:        "expiredMail@gmail.com",
			role:         "moderator",
			duration:     -time.Minute,
			shouldExpire: true,
			shouldFail:   true,
		},
		{
			name:        "Invalid Signing Method",
			id:          uesrId,
			email:       "",
			role:        "employee",
			duration:    time.Minute,
			signingNone: true,
			shouldFail:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tokenStr string

			var claims *UserClaims

			var err error

			if tt.signingNone {
				claims, _ = NewUserClaims(tt.id, tt.email, tt.role, tt.duration)
				token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
				tokenStr, _ = token.SignedString(jwt.UnsafeAllowNoneSignatureType)
			} else {
				tokenStr, claims, err = maker.CreateToken(tt.id, tt.email, tt.role, tt.duration)
				assert.NoError(t, err, "CreateToken should not return an error")
				assert.NotEmpty(t, tokenStr, "Token should not be empty")
				assert.NotNil(t, claims, "Claims should not be nil")
			}

			parsedClaims, err := maker.VerifyToken(tokenStr)

			if tt.shouldFail {
				assert.Error(t, err, "VerifyToken should return an error")
				assert.Nil(t, parsedClaims, "Claims should be nil for invalid tokens")
			} else {
				assert.NoError(t, err, "VerifyToken should not return an error")
				assert.NotNil(t, parsedClaims, "Claims should not be nil")
				assert.Equal(t, tt.id, parsedClaims.ID, "UserID should match")
				assert.Equal(t, tt.email, parsedClaims.Email, "Email should match")

				if tt.shouldExpire {
					assert.True(t, parsedClaims.ExpiresAt.Before(time.Now()), "Token should be expired")
				} else {
					assert.WithinDuration(t, time.Now().Add(tt.duration),
						parsedClaims.ExpiresAt.Time, time.Second, "Expiration time should be correct")
				}
			}
		})
	}
}
