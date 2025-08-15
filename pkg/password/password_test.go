package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword"
	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	// Hashing the same password should produce a different hash each time
	hashedPassword2, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEqual(t, hashedPassword, hashedPassword2)
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword"
	hashedPassword, _ := HashPassword(password)

	// Correct password
	assert.True(t, CheckPasswordHash(password, hashedPassword))

	// Incorrect password
	assert.False(t, CheckPasswordHash("wrongpassword", hashedPassword))

	// Empty password
	assert.False(t, CheckPasswordHash("", hashedPassword))
}
