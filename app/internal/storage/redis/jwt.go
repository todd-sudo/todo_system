package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
)

type JWTStorage interface {
	SetRefreshToken(ctx context.Context, userID string, token string, expiresIn time.Duration) error
	GetRefreshToken(ctx context.Context, userID string) (string, error)
}

type jwtStorage struct {
	ctx context.Context
	rc  *redis.Client
}

func NewJWTStorage(ctx context.Context, rc *redis.Client) JWTStorage {
	return &jwtStorage{
		ctx: ctx,
		rc:  rc,
	}
}

// SetRefreshToken - set refresh token in redis db
func (j *jwtStorage) SetRefreshToken(ctx context.Context, userID string, token string, expiresIn time.Duration) error {
	err := j.rc.Set(ctx, "token_"+userID, token, expiresIn).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetRefreshToken - get refresh token in redis db
func (j *jwtStorage) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	token, err := j.rc.Get(ctx, "token_"+userID).Result()
	if err != nil {
		return "", err
	}
	return token, nil
}
