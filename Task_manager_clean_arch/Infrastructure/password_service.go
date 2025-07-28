package infrastructure

import "golang.org/x/crypto/bcrypt"

type IPasswordService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type bcryptPasswordService struct{}

func NewPasswordService() IPasswordService {
	return &bcryptPasswordService{}
}

func (s *bcryptPasswordService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *bcryptPasswordService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
