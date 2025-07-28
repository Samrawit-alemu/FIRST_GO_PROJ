package repositories

import (
	"context"
	"log"
	"os"
	"taskmanager/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoUserTestSuite struct {
	suite.Suite
	client   *mongo.Client
	userRepo IUserRepository
	dbName   string
}

func (s *MongoUserTestSuite) SetupSuite() {
	mongoURI := os.Getenv("MONGO_TEST_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to Mongo for testing: %v", err)
	}

	s.client = client
	s.dbName = "taskmanager_testdb"
	db := s.client.Database(s.dbName)
	s.userRepo = NewUserRepository(db)
}

func (s *MongoUserTestSuite) TearDownSuite() {
	err := s.client.Database(s.dbName).Drop(context.TODO())
	assert.NoError(s.T(), err, "Failed to drop test database")

	err = s.client.Disconnect(context.TODO())
	assert.NoError(s.T(), err, "Failed to disconnect from Mongo")
}

func (s *MongoUserTestSuite) BeforeTest(suiteName, testName string) {
	err := s.client.Database(s.dbName).Collection("users").Drop(context.TODO())
	assert.NoError(s.T(), err, "Failed to drop users collection before test")
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(MongoUserTestSuite))
}

func (s *MongoUserTestSuite) TestCreateAndFindByUsername() {
	assert := assert.New(s.T())
	ctx := context.Background()

	testUser := &domain.User{
		Username: "testuser",
		Password: "hashedpassword",
		Role:     "user",
	}

	err := s.userRepo.Create(ctx, testUser)
	assert.NoError(err, "Creating user should not produce an error")
	assert.NotNil(testUser.ID, "The user ID should be set by the Create method")

	foundUser, err := s.userRepo.FindByUsername(ctx, "testuser")
	assert.NoError(err, "Finding an existing user should not produce an error")
	assert.NotNil(foundUser, "Found user should not be nil")

	assert.Equal(testUser.ID, foundUser.ID)
	assert.Equal("testuser", foundUser.Username)
	assert.Equal("hashedpassword", foundUser.Password)
}

func (s *MongoUserTestSuite) TestFindByUsername_NotFound() {
	assert := assert.New(s.T())
	ctx := context.Background()

	foundUser, err := s.userRepo.FindByUsername(ctx, "nonexistentuser")

	assert.Error(err, "Should return an error for a non-existent user")
	assert.Nil(foundUser, "Should not return a user object")

	assert.Equal(mongo.ErrNoDocuments, err, "Error should be mongo.ErrNoDocuments")
}

func (s *MongoUserTestSuite) TestCreate_DuplicateUsername() {
	assert := assert.New(s.T())
	ctx := context.Background()

	user1 := &domain.User{Username: "duplicate", Password: "pw1", Role: "user"}
	user2 := &domain.User{Username: "duplicate", Password: "pw2", Role: "user"}

	err := s.userRepo.Create(ctx, user1)
	assert.NoError(err)

	err = s.userRepo.Create(ctx, user2)
	assert.Error(err, "Should return an error for a duplicate username")

	assert.True(mongo.IsDuplicateKeyError(err), "Error should be a duplicate key error")
}
