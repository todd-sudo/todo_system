package service_rd

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
	jwtStorage "github.com/todd-sudo/todo_system/internal/storage/redis"
)

type RedisService interface {
	SetRefreshToken(ctx context.Context, username, tokenID, token string, expiresIn time.Duration) error
	GetRefreshToken(ctx context.Context, tokenID string) (string, error)
	DelRefreshToken(ctx context.Context, tokenID string) (int64, error)
}

type redisService struct {
	ctx        context.Context
	rc         *redis.Client
	jwtStorage jwtStorage.JWTStorage
}

func NewRedisService(ctx context.Context, rc *redis.Client, jwtStorage jwtStorage.JWTStorage) RedisService {
	return &redisService{
		ctx:        ctx,
		rc:         rc,
		jwtStorage: jwtStorage,
	}
}

func (s *redisService) SetRefreshToken(ctx context.Context, username, tokenID, token string, expiresIn time.Duration) error {

	if err := s.jwtStorage.SetRefreshToken(ctx, username, tokenID, token, expiresIn); err != nil {
		return err
	}
	return nil
}

func (s *redisService) GetRefreshToken(ctx context.Context, tokenID string) (string, error) {
	token, err := s.jwtStorage.GetRefreshToken(ctx, tokenID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *redisService) DelRefreshToken(ctx context.Context, tokenID string) (int64, error) {
	deleted, err := s.jwtStorage.DelRefreshToken(ctx, tokenID)
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
