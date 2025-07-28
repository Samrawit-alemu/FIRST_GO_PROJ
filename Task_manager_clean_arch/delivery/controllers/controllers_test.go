package controllers_test

import (
	"context"
	"log"
	"os"
	"taskmanager/domain"
	"taskmanager/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoTestSuite is a test suite that sets up a connection to a real MongoDB
// and cleans up the database after the tests are done.
type MongoTestSuite struct {
	suite.Suite
	client   *mongo.Client
	userRepo repositories.IUserRepository
	dbName   string
}

// SetupSuite runs once before all tests in the suite.
func (s *MongoTestSuite) SetupSuite() {
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
	s.userRepo = repositories.NewUserRepository(db)
}

// TearDownSuite runs once after all tests in the suite are finished.
func (s *MongoTestSuite) TearDownSuite() {
	err := s.client.Database(s.dbName).Drop(context.TODO())
	assert.NoError(s.T(), err, "Failed to drop test database")

	err = s.client.Disconnect(context.TODO())
	assert.NoError(s.T(), err, "Failed to disconnect from Mongo")
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(MongoTestSuite))
}

func (s *MongoTestSuite) TestCreateAndFindByUsername() {
	assert := assert.New(s.T())
	testUser := &domain.User{
		Username: "testuser",
		Password: "hashedpassword",
		Role:     "user",
	}

	// --- Test Create ---
	err := s.userRepo.Create(context.Background(), testUser)
	assert.NoError(err)
	assert.NotNil(testUser.ID)

	// --- Test FindByUsername ---
	foundUser, err := s.userRepo.FindByUsername(context.Background(), "testuser")
	assert.NoError(err)
	assert.NotNil(foundUser)
	assert.Equal(testUser.Username, foundUser.Username)
	assert.Equal(testUser.Password, foundUser.Password)
}

func (s *MongoTestSuite) TestFindByUsername_NotFound() {
	assert := assert.New(s.T())

	foundUser, err := s.userRepo.FindByUsername(context.Background(), "nonexistentuser")
	assert.Error(err)
	assert.Nil(foundUser)
	assert.Equal(mongo.ErrNoDocuments, err)
}
