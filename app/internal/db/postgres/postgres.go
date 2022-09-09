package postgres

import (
	"fmt"

	"github.com/todd-sudo/todo_system/internal/config"
	"github.com/todd-sudo/todo_system/pkg/logging"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(cfg *config.Config, log *logging.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.PostgreSQL.Username,
		cfg.PostgreSQL.Password,
		cfg.PostgreSQL.Host,
		cfg.PostgreSQL.Port,
		cfg.PostgreSQL.Database,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connection database: %v", err)
		return nil, err
	}

	return db, nil
}
