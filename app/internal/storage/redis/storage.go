package redis

import (
	"context"

	"github.com/go-redis/redis/v9"
)

type Storage struct {
	JWTStorage
}

func NewStorage(ctx context.Context, rc *redis.Client) *Storage {
	return &Storage{
		JWTStorage: NewJWTStorage(ctx, rc),
	}
}
