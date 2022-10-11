package service_pg

import (
	"context"

	"github.com/todd-sudo/todo_system/internal/hasher"
	"github.com/todd-sudo/todo_system/internal/storage/postgres"

	"github.com/todd-sudo/todo_system/pkg/logging"
)

type Service struct {
	UserService
	FolderService
	ItemService
}

func NewService(ctx context.Context, storage postgres.Storage, log logging.Logger, hasher hasher.PasswordHasher) *Service {
	return &Service{
		UserService:   NewUserService(ctx, log, &storage, hasher),
		FolderService: NewFolderService(ctx, log, &storage),
		ItemService:   NewItemService(ctx, log, &storage),
	}
}
