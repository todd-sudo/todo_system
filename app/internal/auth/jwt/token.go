package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/todd-sudo/todo_system/internal/config"
	"github.com/todd-sudo/todo_system/pkg/logging"
)

type JWTToken interface {
	CreateToken(ttlMinutes time.Duration, payload string, jwtKey string) (string, error)
	ValidateToken(token string, jwtKey string) (string, error)
}

type jwtToken struct {
	log logging.Logger
	cfg config.Config
}

// type tokenClaims

func NewJWTToken(log logging.Logger, cfg config.Config) JWTToken {
	return &jwtToken{
		log: log,
		cfg: cfg,
	}
}

type tokenClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

// CreateToken - create JWT token
func (j *jwtToken) CreateToken(ttlMinutes time.Duration, payload string, jwtKey string) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttlMinutes * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Username: payload,
	}).SignedString([]byte(jwtKey))
	if err != nil {
		j.log.Error(err)
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

// ValidateToken - validate JWT token
func (j *jwtToken) ValidateToken(token string, jwtKey string) (string, error) {
	accessToken, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(jwtKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := accessToken.Claims.(*tokenClaims)
	if !ok {
		return "", errors.New("token claims are not of type *tokenClaims")
	}

	return fmt.Sprintf("%v", claims.Username), nil
}
