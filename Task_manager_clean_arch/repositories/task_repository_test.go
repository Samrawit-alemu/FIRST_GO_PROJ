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

type MongoTaskTestSuite struct {
	suite.Suite
	taskRepo ITaskRepository
	userRepo IUserRepository
	dbName   string
	client   *mongo.Client
}

func (s *MongoTaskTestSuite) SetupSuite() {
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
	s.taskRepo = NewTaskRepository(db)
}

func (s *MongoTaskTestSuite) TearDownSuite() {
	err := s.client.Database(s.dbName).Drop(context.TODO())
	assert.NoError(s.T(), err, "Failed to drop test database")

	err = s.client.Disconnect(context.TODO())
	assert.NoError(s.T(), err, "Failed to disconnect from Mongo")
}

// This function runs the test suite.
func TestTaskRepository(t *testing.T) {
	suite.Run(t, new(MongoTaskTestSuite))
}

func (s *MongoTaskTestSuite) TestCreateAndGetTasks() {
	assert := assert.New(s.T())

	owner := &domain.User{Username: "taskowner", Password: "pw", Role: "user"}
	err := s.userRepo.Create(context.Background(), owner)
	assert.NoError(err)

	task1 := &domain.Task{
		Title:  "Task 1",
		Status: "Pending",
		UserID: owner.ID,
	}
	task2 := &domain.Task{
		Title:  "Task 2",
		Status: "In Progress",
		UserID: owner.ID,
	}

	// Test Create
	err = s.taskRepo.Create(context.Background(), task1)
	assert.NoError(err)
	assert.NotNil(task1.ID)
	err = s.taskRepo.Create(context.Background(), task2)
	assert.NoError(err)
	assert.NotNil(task2.ID)

	// Test GetAllByUserID
	tasks, err := s.taskRepo.GetAllByUserID(context.Background(), owner.ID)
	assert.NoError(err)
	assert.Len(tasks, 2)

	// Test GetByID
	foundTask, err := s.taskRepo.GetByID(context.Background(), task1.ID)
	assert.NoError(err)
	assert.Equal(task1.Title, foundTask.Title)
}
