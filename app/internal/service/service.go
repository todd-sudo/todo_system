package service

import (
	"context"

	"github.com/todd-sudo/todo_system/internal/repository"
	"github.com/todd-sudo/todo_system/pkg/logging"
)

type Service struct {
}

func NewService(ctx context.Context, r repository.Repository, log logging.Logger) *Service {
	return &Service{}
}
