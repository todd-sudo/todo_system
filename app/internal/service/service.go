package service

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/todd-sudo/todo_system/internal/storage/postgres"

	"github.com/todd-sudo/todo_system/pkg/logging"
)

type Service struct {
	ctx     context.Context
	storage postgres.Storage
	log     logging.Logger
	rc      *redis.Client
}

func NewService(ctx context.Context, storage postgres.Storage, log logging.Logger, rc *redis.Client) *Service {
	return &Service{
		ctx:     ctx,
		storage: storage,
		log:     log,
		rc:      rc,
	}
}
