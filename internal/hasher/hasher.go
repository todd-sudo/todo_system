package hasher

import (
	"github.com/todd-sudo/todo_system/pkg/logging"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type SHA1Hasher struct {
	log logging.Logger
}

func NewSHA1Hasher(log logging.Logger) *SHA1Hasher {
	return &SHA1Hasher{log: log}
}

func (h *SHA1Hasher) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
