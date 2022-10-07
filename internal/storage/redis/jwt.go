package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
)

type JWTStorage interface {
	SetRefreshToken(ctx context.Context, userID string, token string, expiresIn time.Duration) error
	GetRefreshToken(ctx context.Context, userID string) (string, error)
	DelRefreshToken(ctx context.Context, username string) (int64, error)
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
func (j *jwtStorage) SetRefreshToken(ctx context.Context, username string, token string, expiresIn time.Duration) error {
	err := j.rc.Set(ctx, "token_"+username, token, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetRefreshToken - get refresh token in redis db
func (j *jwtStorage) GetRefreshToken(ctx context.Context, username string) (string, error) {
	token, err := j.rc.Get(ctx, "token_"+username).Result()
	if err != nil {
		return "", err
	}
	return token, nil
}

// DelRefreshToken - delete refresh token in redis db
func (j *jwtStorage) DelRefreshToken(ctx context.Context, username string) (int64, error) {
	deleted, err := j.rc.Del(ctx, "token_"+username).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
