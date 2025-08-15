package jwt

import (
	"testing"
	"time"

	
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateTokens(t *testing.T) {
	userID := uuid.New()
	userRole := "user"
	secret := "supersecretkey"
	accessTokenTTL := time.Minute * 15
	refreshTokenTTL := time.Hour * 24 * 7

	accessToken, refreshToken, err := GenerateTokens(userID, userRole, secret, accessTokenTTL, refreshTokenTTL)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	// Validate access token
	claims, err := ValidateToken(accessToken, secret)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, userRole, claims.Role)
	assert.True(t, claims.ExpiresAt.Time.After(time.Now()))

	// Validate refresh token
	claims, err = ValidateToken(refreshToken, secret)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, userRole, claims.Role)
	assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
}

func TestValidateToken(t *testing.T) {
	userID := uuid.New()
	userRole := "admin"
	secret := "anothersecretkey"
	accessTokenTTL := time.Minute * 1

	accessToken, _, _ := GenerateTokens(userID, userRole, secret, accessTokenTTL, time.Minute*5) // Refresh token not used here

	// Valid token
	claims, err := ValidateToken(accessToken, secret)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, userRole, claims.Role)

	// Invalid secret
	_, err = ValidateToken(accessToken, "wrongsecret")
	assert.Error(t, err)

	// Expired token (simulate by waiting)
	time.Sleep(accessTokenTTL + time.Second) // Wait for token to expire
	_, err = ValidateToken(accessToken, secret)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is expired")
}
