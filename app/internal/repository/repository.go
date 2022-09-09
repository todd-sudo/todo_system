package repository

import (
	"context"

	"github.com/todd-sudo/todo_system/pkg/logging"
	"gorm.io/gorm"
)

type Repository struct {
}

func NewRepository(ctx context.Context, db *gorm.DB, log logging.Logger) *Repository {
	return &Repository{
		// City: NewCityRepository(ctx, db, log),
		// Shop: NewShopRepository(ctx, db, log),
	}
}
