package postgres

import (
	"context"

	"github.com/todd-sudo/todo_system/internal/entity"
	"github.com/todd-sudo/todo_system/pkg/logging"
	"gorm.io/gorm"
)

type UserStorage interface {
	InsertUser(ctx context.Context, user *entity.User) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	DeleteUser(ctx context.Context, username string) error
	FindUserByUsername(ctx context.Context, username string) (*entity.User, error)
	ProfileUser(ctx context.Context, username string) (*entity.User, error)
}

type userStorage struct {
	ctx        context.Context
	connection *gorm.DB
	log        logging.Logger
}

func NewUserStorage(ctx context.Context, db *gorm.DB, log logging.Logger) UserStorage {
	return &userStorage{
		ctx:        ctx,
		connection: db,
		log:        log,
	}
}

// InsertUser - insert user in db
func (db *userStorage) InsertUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	tx := db.connection.WithContext(ctx)
	res := tx.Save(&user)
	if res.Error != nil {
		db.log.Errorf("insert user error %v", res.Error)
		return nil, res.Error
	}
	return user, nil
}

// UpdateUser - update user in db
func (db *userStorage) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	tx := db.connection.WithContext(ctx)
	res := tx.Save(&user)
	if res.Error != nil {
		db.log.Errorf("update user error %v", res.Error)
		return nil, res.Error
	}
	return user, nil
}

// DeleteUser - delete user from db
func (db *userStorage) DeleteUser(ctx context.Context, username string) error {
	tx := db.connection.WithContext(ctx)
	var user *entity.User
	res := tx.Where(`username = ?`, username).Delete(&user)
	if res.Error != nil {
		db.log.Errorf("delete user error %v", res.Error)
		return res.Error
	}
	return nil
}

// FindByUsername - find user by 'username' from db
func (db *userStorage) FindUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	tx := db.connection.WithContext(ctx)
	var user *entity.User
	res := tx.Where("username = ?", username).Take(&user)
	if res.Error != nil {
		db.log.Errorf("find by username user error %v", res.Error)
		return nil, res.Error
	}
	return user, nil
}

// Вывод профиля пользователя
func (db *userStorage) ProfileUser(ctx context.Context, username string) (*entity.User, error) {
	tx := db.connection.WithContext(ctx).Debug()
	var user *entity.User
	//.Preload("Folders").Preload("Folders.User")
	res := tx.Where(`username = ?`, username).Find(&user)
	if res.Error != nil {
		db.log.Errorf("get profile user error %v", res.Error)
		return nil, res.Error
	}
	return user, nil
}
