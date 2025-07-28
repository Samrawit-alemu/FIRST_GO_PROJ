package infrastructure

import (
	"os"
	"taskmanager/domain"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestJWTService(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret_key_for_jwt")
	defer os.Unsetenv("JWT_SECRET")

	jwtService := NewJWTService()
	userID := primitive.NewObjectID()

	user := domain.User{
		ID:       userID,
		Username: "testuser",
		Role:     "admin",
	}

	// Test Token Generation
	tokenString, err := jwtService.GenerateToken(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Test Token Validation (Success)
	validatedToken, err := jwtService.ValidateToken(tokenString)
	assert.NoError(t, err)
	assert.NotNil(t, validatedToken)
	assert.True(t, validatedToken.Valid)

	// Check the claims
	claims, ok := validatedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, userID.Hex(), claims["user_id"])
	assert.Equal(t, "admin", claims["role"])

	// Test Token Validation (Failure - Malformed Token)
	_, err = jwtService.ValidateToken("this.is.a.bad.token")
	assert.Error(t, err)
}
