package postgres

import (
	"context"

	"github.com/todd-sudo/todo_system/internal/entity"
	"github.com/todd-sudo/todo_system/pkg/logging"
	"gorm.io/gorm"
)

type FolderStorage interface {
	AllFolders(ctx context.Context, username string) ([]*entity.Folder, error)
	InsertUpdateFolder(ctx context.Context, folder *entity.Folder) (*entity.Folder, error)
	DeleteFolder(ctx context.Context, folderID int) error
}

type folderStorage struct {
	ctx        context.Context
	connection *gorm.DB
	log        logging.Logger
}

func NewFolderStorage(ctx context.Context, db *gorm.DB, log logging.Logger) FolderStorage {
	return &folderStorage{
		ctx:        ctx,
		connection: db,
		log:        log,
	}
}

// AllFolders - get all folders by username
func (db *folderStorage) AllFolders(ctx context.Context, username string) ([]*entity.Folder, error) {
	tx := db.connection.WithContext(ctx)
	var folders []*entity.Folder
	if err := tx.Preload("User").Joins("User").Where(
		`"username" = ?`,
		username,
	).Find(&folders).Error; err != nil {
		return nil, err
	}
	return folders, nil
}

// InsertFolder - insert folder in db
func (db *folderStorage) InsertUpdateFolder(ctx context.Context, folder *entity.Folder) (*entity.Folder, error) {
	tx := db.connection.WithContext(ctx)
	if err := tx.Save(&folder).Error; err != nil {
		db.log.Errorf("insert folder error %v", err.Error())
		return nil, err
	}
	return folder, nil
}

// DeleteFolder - delete folder from db
func (db *folderStorage) DeleteFolder(ctx context.Context, folderID int) error {
	tx := db.connection.WithContext(ctx)
	var folder *entity.Folder
	if err := tx.Where(`id = ?`, folderID).Delete(&folder).Error; err != nil {
		db.log.Errorf("delete folder error %v", err.Error())
		return err
	}
	return nil
}
