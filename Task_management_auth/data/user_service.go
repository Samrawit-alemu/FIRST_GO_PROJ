package data

import (
	"context"
	"errors"
	"log"
	"taskmanager/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// ensureUsernameUnique creates a unique index on the username field.
func EnsureUsernameUnique() {
	collection := GetUsersCollection()
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"username": 1}, // 1 for ascending order
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Printf("Could not create unique index for username: %v\n", err)
	}
}

// In data/user_service.go

// CreateUser handles the logic of creating a new user.
func CreateUser(user models.User) (models.User, error) {

	collection := GetUsersCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}
	user.Password = string(hashedPassword)
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return models.User{}, err
	}
	if count == 0 {
		user.Role = "admin"
	} else {
		user.Role = "user"
	}

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return models.User{}, errors.New("username already exists")
		}
		return models.User{}, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

// LoginUser checks user credentials and returns the user if they are valid.
func LoginUser(username, password string) (models.User, error) {
	collection := GetUsersCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	// 1. Find the user by their username.
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return models.User{}, errors.New("invalid username or password")
	}

	// 2. Compare the provided password with the stored hash.
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return models.User{}, errors.New("invalid username or password")
	}

	// 3. Login successful. Return the user object.
	return user, nil
}

// PromoteUser changes a user's role to 'admin'.
func PromoteUser(userID string) (models.User, error) {
	collection := GetUsersCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return models.User{}, errors.New("invalid user ID format")
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"role": "admin"}}

	// We use FindOneAndUpdate to get the updated document back in one go.
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedUser models.User
	err = collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedUser)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}

	return updatedUser, nil
}
