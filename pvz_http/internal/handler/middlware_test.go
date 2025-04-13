package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaksimovDenis/avito_pvz/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMiddlewareFunc(t *testing.T) {
	secretKey := "supersecretkey"
	tokenMaker := token.NewJWTMaker(secretKey)

	tests := []struct {
		name            string
		authHeader      string
		expectedStatus  int
		expectedError   string
		expectedUser    *token.UserClaims
		expectedUserKey string
	}{
		{
			name:           "No Authorization Header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Неавторизован",
		},
		{
			name:           "Invalid Authorization Header Format",
			authHeader:     "InvalidToken",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "invalid autorization header",
		},
		{
			name:           "Invalid Bearer Prefix",
			authHeader:     "Token 123",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "invalid autorization header",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			responseRecord := httptest.NewRecorder()

			testCtx, _ := gin.CreateTestContext(responseRecord)

			testCtx.Request = httptest.NewRequest("GET", "/", nil)

			testCtx.Request.Header.Set("Authorization", tt.authHeader)

			middleware := GetMiddlewareFunc(tokenMaker)
			middleware(testCtx)
			assert.Equal(t, tt.expectedStatus, responseRecord.Code)

			if tt.expectedError != "" {
				assert.JSONEq(t, fmt.Sprintf(`{"error": "%s"}`, tt.expectedError), responseRecord.Body.String())
			} else {
				user, exists := testCtx.Get(tt.expectedUserKey)
				if !assert.True(t, exists) {
					t.Fatalf("Expected user in context, but got none")
				}

				assert.Equal(t, tt.expectedUser.Email, user.(*token.UserClaims).Email)
			}
		})
	}
}

func TestVerifyClaimsFromAuthHeader(t *testing.T) {
	secretKey := "supersecretkey"
	tokenMaker := token.NewJWTMaker(secretKey)

	tests := []struct {
		name          string
		authHeader    string
		mockError     error
		expectedError string
	}{
		{
			name:          "No Authorization Header",
			authHeader:    "",
			expectedError: "Неавторизован",
		},
		{
			name:          "Invalid Authorization Header Format",
			authHeader:    "InvalidToken",
			expectedError: "invalid autorization header",
		},
		{
			name:          "Invalid Bearer Prefix",
			authHeader:    "Token 123",
			expectedError: "invalid autorization header",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			respRecord := httptest.NewRecorder()

			testCtx, _ := gin.CreateTestContext(respRecord)

			testCtx.Request = httptest.NewRequest("GET", "/", nil)
			testCtx.Request.Header.Set("Authorization", tt.authHeader)

			claims, err := verifyClaimsFromAuthHeader(testCtx, *tokenMaker)

			if tt.expectedError != "" {
				assert.Nil(t, claims)
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NotNil(t, claims)
				assert.Nil(t, err)
			}
		})
	}
}
