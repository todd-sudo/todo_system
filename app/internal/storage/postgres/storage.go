package postgres

import (
	"context"

	"github.com/todd-sudo/todo_system/pkg/logging"
	"gorm.io/gorm"
)

type Storage struct {
	UserStorage
	FolderStorage
	ItemStorage
}

func NewStorage(ctx context.Context, db *gorm.DB, log logging.Logger) *Storage {
	return &Storage{
		UserStorage:   NewUserStorage(ctx, db, log),
		FolderStorage: NewFolderStorage(ctx, db, log),
		ItemStorage:   NewItemStorage(ctx, db, log),
	}
}
