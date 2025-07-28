package usecases

import (
	"context"
	"taskmanager/domain"
	"taskmanager/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

// TestRegister_Success_FirstUserIsAdmin tests  for registration
// where the user is the first one and should become an admin.
func TestRegister_Success_FirstUserIsAdmin(t *testing.T) {

	// 1. Create instances of our mocks. These are our "stunt doubles".
	mockUserRepo := new(mocks.IUserRepository)
	mockPasswordSvc := new(mocks.IPasswordService)
	mockJwtSvc := new(mocks.IJWTService)

	// 2. Define the input we will pass to the function we are testing.
	username := "adminuser"
	password := "password123"
	hashedPassword := "hashed_password" // A fake hash for our test.

	// 3. Set up the "expectations" for our mocks. This is the most important part.

	mockUserRepo.On("FindByUsername", mock.Anything, username).Return(nil, mongo.ErrNoDocuments)

	mockUserRepo.On("Count", mock.Anything).Return(int64(0), nil)

	mockPasswordSvc.On("HashPassword", password).Return(hashedPassword, nil)

	mockUserRepo.On("Create", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
		return user.Role == "admin" && user.Username == username && user.Password == hashedPassword
	})).Return(nil)

	usecase := NewUserUsecase(mockUserRepo, mockPasswordSvc, mockJwtSvc)
	createdUser, err := usecase.Register(context.Background(), username, password)

	// Use testify's assertion library to make our checks clean and readable.
	assert.NoError(t, err)                     // We assert that no error was returned.
	assert.NotNil(t, createdUser)              // We assert that we got a user object back.
	assert.Equal(t, "admin", createdUser.Role) // We assert that the user's role is "admin".
	assert.Equal(t, username, createdUser.Username)

	mockUserRepo.AssertExpectations(t)
	mockPasswordSvc.AssertExpectations(t)
}

// TestRegister_Failure_UsernameExists tests the edge case where the username is already taken.
func TestRegister_Failure_UsernameExists(t *testing.T) {
	mockUserRepo := new(mocks.IUserRepository)
	mockPasswordSvc := new(mocks.IPasswordService)
	mockJwtSvc := new(mocks.IJWTService)

	username := "existinguser"
	password := "password123"

	mockUserRepo.On("FindByUsername", mock.Anything, username).Return(&domain.User{}, nil)

	usecase := NewUserUsecase(mockUserRepo, mockPasswordSvc, mockJwtSvc)
	createdUser, err := usecase.Register(context.Background(), username, password)

	// --- ASSERT ---
	assert.Error(t, err)                                    // We assert that an error WAS returned.
	assert.Nil(t, createdUser)                              // We assert that no user object was returned.
	assert.Equal(t, "username already exists", err.Error()) // We check for the specific error message.

	mockPasswordSvc.AssertNotCalled(t, "HashPassword", mock.Anything)
	mockUserRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}
