package service_rd

import (
	"context"

	"github.com/todd-sudo/todo_system/pkg/logging"
)

type Service struct {
}

func NewService(ctx context.Context, log logging.Logger) *Service {
	return &Service{}
}
