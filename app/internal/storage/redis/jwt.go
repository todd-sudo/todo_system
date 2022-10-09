package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
)

type JWTStorage interface {
	SetRefreshToken(ctx context.Context, username, tokenID, token string, expiresIn time.Duration) error
	GetRefreshToken(ctx context.Context, tokenID string) (string, error)
	DelRefreshToken(ctx context.Context, tokenID string) (int64, error)
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
func (j *jwtStorage) SetRefreshToken(
	ctx context.Context,
	username, tokenID, token string,
	expiresIn time.Duration,
) error {
	var cursor uint64
	var keys []string
	var err error

	keys, _, err = j.rc.Scan(ctx, cursor, fmt.Sprintf("%s_*", username), 0).Result()
	if err != nil {
		return err
	}

	if len(keys) == 5 {
		j.rc.Del(ctx, keys[len(keys)-1])
	}
	// if len(keys) > 5 {
	// 	for _, k := range keys {
	// 		j.rc.Del(ctx, k)
	// 	}
	// }

	err = j.rc.Set(ctx, username+"_"+tokenID, token, expiresIn).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetRefreshToken - get refresh token in redis db
func (j *jwtStorage) GetRefreshToken(ctx context.Context, tokenID string) (string, error) {
	var cursor uint64
	var keys []string
	var err error

	keys, _, err = j.rc.Scan(ctx, cursor, fmt.Sprintf("*_%s", tokenID), 0).Result()
	if err != nil {
		return "", err
	}

	if len(keys) == 0 {
		return "", err
	}
	key := keys[len(keys)-1]

	token, err := j.rc.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return token, nil
}

// DelRefreshToken - delete refresh token in redis db
func (j *jwtStorage) DelRefreshToken(ctx context.Context, tokenID string) (int64, error) {
	var cursor uint64
	var keys []string
	var err error
	keys, _, err = j.rc.Scan(ctx, cursor, fmt.Sprintf("*_%s", tokenID), 0).Result()
	if err != nil {
		return 0, err
	}

	if len(keys) == 0 {
		return 0, err
	}
	key := keys[len(keys)-1]

	deleted, err := j.rc.Del(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
