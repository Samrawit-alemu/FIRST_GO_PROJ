package usecases

import (
	"context"
	"errors"
	"taskmanager/domain"
	"taskmanager/infrastructure"
	"taskmanager/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IUserUsecase interface {
	Register(ctx context.Context, username, password string) (*domain.User, error)
	Login(ctx context.Context, username, password string) (string, error)
	Promote(ctx context.Context, userID string) (*domain.User, error)
}

type userUsecase struct {
	userRepo        repositories.IUserRepository
	passwordService infrastructure.IPasswordService
	jwtService      infrastructure.IJWTService
}

func NewUserUsecase(repo repositories.IUserRepository, ps infrastructure.IPasswordService, js infrastructure.IJWTService) IUserUsecase {
	return &userUsecase{
		userRepo:        repo,
		passwordService: ps,
		jwtService:      js,
	}
}

func (uc *userUsecase) Register(ctx context.Context, username, password string) (*domain.User, error) {
	_, err := uc.userRepo.FindByUsername(ctx, username)
	if err != mongo.ErrNoDocuments {
		if err == nil {
			return nil, errors.New("username already exists")
		}
		return nil, err
	}

	hashedPassword, err := uc.passwordService.HashPassword(password)
	if err != nil {
		return nil, err
	}

	userCount, err := uc.userRepo.Count(ctx)
	if err != nil {
		return nil, err
	}
	role := "user"
	if userCount == 0 {
		role = "admin"
	}

	user := &domain.User{
		Username: username,
		Password: hashedPassword,
		Role:     role,
	}

	err = uc.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *userUsecase) Login(ctx context.Context, username, password string) (string, error) {
	user, err := uc.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	if !uc.passwordService.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid username or password")
	}

	return uc.jwtService.GenerateToken(*user)
}

func (uc *userUsecase) Promote(ctx context.Context, userID string) (*domain.User, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user, err := uc.userRepo.FindByID(ctx, objectID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.Role = "admin"

	err = uc.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
