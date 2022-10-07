package postgres

import (
	"context"

	"github.com/todd-sudo/todo_system/pkg/logging"
	"gorm.io/gorm"
)

type Storage struct {
	UserStorage
}

func NewStorage(ctx context.Context, db *gorm.DB, log logging.Logger) *Storage {
	return &Storage{
		UserStorage: NewUserStorage(ctx, db, log),
	}
}
