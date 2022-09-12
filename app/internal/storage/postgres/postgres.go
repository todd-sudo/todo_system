package postgres

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/todd-sudo/todo_system/pkg/logging"
	"gorm.io/gorm"
)

type Storage struct {
	ctx context.Context
	log logging.Logger
	rc  *redis.Client
	db  *gorm.DB
}

func NewStorage(ctx context.Context, db *gorm.DB, log logging.Logger, rc *redis.Client) *Storage {
	return &Storage{
		ctx: ctx,
		log: log,
		rc:  rc,
		db:  db,
	}
}
