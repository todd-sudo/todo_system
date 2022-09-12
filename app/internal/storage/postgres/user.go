package postgres

import (
	"context"

	"github.com/todd-sudo/todo_system/internal/entity"
	"github.com/todd-sudo/todo_system/pkg/logging"
	"gorm.io/gorm"
)

type UserStorage interface {
	InsertUser(ctx context.Context, user *entity.User) (*entity.User, error)
	// UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	// FindByEmail(ctx context.Context, email string) (*entity.User, error)
	// ProfileUser(ctx context.Context, userID string) (*entity.User, error)
}

type userStorage struct {
	ctx        context.Context
	connection *gorm.DB
	log        logging.Logger
}

//NewUserRepository is creates a new instance of UserRepository
func NewUserStorage(ctx context.Context, db *gorm.DB, log logging.Logger) UserStorage {
	return &userStorage{
		ctx:        ctx,
		connection: db,
		log:        log,
	}
}

func (db *userStorage) InsertUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	tx := db.connection.WithContext(ctx)
	res := tx.Save(&user)
	if res.Error != nil {
		db.log.Errorf("insert user error %v", res.Error)
		return nil, res.Error
	}
	return user, nil
}
