package service

import (
	"crypto/sha1"
	"fmt"

	"github.com/AlexPop69/todo-app"
	"github.com/AlexPop69/todo-app/pkg/repository"
)

const salt = "f7g6h5j4kl"

type AuthService struct {
	rep repository.Autorization
}

func NewAuthService(rep repository.Autorization) *AuthService {
	return &AuthService{rep: rep}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) {
	user.Password = s.generatePasswordHash(user.Password)

	return s.rep.CreateUser(user)
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
