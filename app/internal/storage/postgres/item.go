package postgres

import (
	"context"

	"github.com/todd-sudo/todo_system/internal/entity"
	"github.com/todd-sudo/todo_system/pkg/logging"
	"gorm.io/gorm"
)

type ItemStorage interface {
	AllItemsByFolder(ctx context.Context, folderID int) ([]*entity.Item, error)
	AllItems(ctx context.Context, username string) ([]*entity.Item, error)
	InsertItem(ctx context.Context, item *entity.Item) (*entity.Item, error)
	UpdateItem(ctx context.Context, item *entity.Item) (*entity.Item, error)
	DeleteItem(ctx context.Context, itemID int) error
}

type itemStorage struct {
	ctx        context.Context
	connection *gorm.DB
	log        logging.Logger
}

func NewItemStorage(ctx context.Context, db *gorm.DB, log logging.Logger) ItemStorage {
	return &itemStorage{
		ctx:        ctx,
		connection: db,
		log:        log,
	}
}

// AllItemsByFolder - get all items by folder
func (db *itemStorage) AllItemsByFolder(ctx context.Context, folderID int) ([]*entity.Item, error) {
	tx := db.connection.WithContext(ctx)
	var items []*entity.Item
	if err := tx.Preload("Folder").Joins("Folder").Where(
		`"id" = ?`,
		folderID,
	).Find(&items).Error; err != nil {
		db.log.Errorf("get all items by folder_id error %v", err.Error())
		return nil, err
	}
	return items, nil
}

// AllItems - get all items by username
func (db *itemStorage) AllItems(ctx context.Context, username string) ([]*entity.Item, error) {
	tx := db.connection.WithContext(ctx)
	var items []*entity.Item
	if err := tx.Preload("User").Joins("User").Where(
		`"username" = ?`,
		username,
	).Find(&items).Error; err != nil {
		db.log.Errorf("get all items by username error %v", err.Error())
		return nil, err
	}
	return items, nil
}

// InsertItem - insert item in db
func (db *itemStorage) InsertItem(ctx context.Context, item *entity.Item) (*entity.Item, error) {
	tx := db.connection.WithContext(ctx)
	if err := tx.Save(&item).Error; err != nil {
		db.log.Errorf("insert item error %v", err.Error())
		return nil, err
	}
	return item, nil
}

// UpdateItem - update item in db
func (db *itemStorage) UpdateItem(ctx context.Context, item *entity.Item) (*entity.Item, error) {
	tx := db.connection.WithContext(ctx)
	if err := tx.Save(&item).Error; err != nil {
		db.log.Errorf("update item error %v", err.Error())
		return nil, err
	}
	return item, nil
}

// DeleteItem - delete item from db
func (db *itemStorage) DeleteItem(ctx context.Context, itemID int) error {
	tx := db.connection.WithContext(ctx)
	var item *entity.Item
	if err := tx.Where(`id = ?`, itemID).Delete(&item).Error; err != nil {
		db.log.Errorf("delete item error %v", err.Error())
		return err
	}
	return nil
}
