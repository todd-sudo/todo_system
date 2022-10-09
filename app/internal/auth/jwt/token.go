package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/todd-sudo/todo_system/internal/config"
	"github.com/todd-sudo/todo_system/pkg/logging"
)

type JWTToken interface {
	CreateToken(ttlMinutes time.Duration, payload string, jwtKey string) (string, string, error)
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
	// Uuid     string `json:"uuid"`
}

// CreateToken - create JWT token
func (j *jwtToken) CreateToken(ttlMinutes time.Duration, payload string, jwtKey string) (string, string, error) {
	tokenID := uuid.New().String()
	ttlDur := time.Duration(ttlMinutes * time.Minute)
	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttlDur).Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        tokenID,
		},
		Username: payload,
		// Uuid:     tokenID,
	})

	token, err := tokenJWT.SignedString([]byte(jwtKey))
	if err != nil {
		j.log.Error(err)
		return "", "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, tokenID, nil
}

// ValidateToken - validate JWT token
func (j *jwtToken) ValidateToken(token string, jwtKey string) (string, error) {
	jwtToken, _ := jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(jwtKey), nil
	})
	claims, ok := jwtToken.Claims.(*tokenClaims)
	if !ok {
		return "", errors.New("token claims are not of type *tokenClaims")
	}

	return fmt.Sprintf("%v", claims.Username), nil
}
