package service

import (
	"context"

	"github.com/todd-sudo/todo_system/internal/storage/postgres"
	"github.com/todd-sudo/todo_system/pkg/logging"
)

type Service struct {
}

func NewService(ctx context.Context, r postgres.Storage, log logging.Logger) *Service {
	return &Service{}
}
