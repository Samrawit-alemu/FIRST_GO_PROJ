package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordService(t *testing.T) {
	passwordService := NewPasswordService()

	password := "my_secret_password"

	// Test Hashing
	hashedPassword, err := passwordService.HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, hashedPassword)

	// Test successful check
	match := passwordService.CheckPasswordHash(password, hashedPassword)
	assert.True(t, match, "Password should match the hash")

	// Test failed check
	wrongPassword := "wrong_password"
	match = passwordService.CheckPasswordHash(wrongPassword, hashedPassword)
	assert.False(t, match, "Wrong password should not match the hash")
}
