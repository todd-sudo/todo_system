package service_pg

import (
	"context"

	"github.com/mashingan/smapping"
	"github.com/todd-sudo/todo_system/internal/dto"
	"github.com/todd-sudo/todo_system/internal/entity"
	"github.com/todd-sudo/todo_system/internal/storage/postgres"
	"github.com/todd-sudo/todo_system/pkg/logging"
)

type ItemService interface {
	AllItemsByFolder(ctx context.Context, m *dto.AllItemsByFolderDTO) ([]*entity.Item, error)
	AllItems(ctx context.Context, m *dto.AllItemDTO) ([]*entity.Item, error)
	InsertUpdateItem(ctx context.Context, m *dto.CreateItemDTO) (*entity.Item, error)
	DeleteItem(ctx context.Context, m *dto.DeleteItemDTO) error
}

type itemService struct {
	ctx     context.Context
	storage *postgres.Storage
	log     logging.Logger
}

func NewItemService(ctx context.Context, log logging.Logger, storage *postgres.Storage) ItemService {
	return &itemService{
		ctx:     ctx,
		log:     log,
		storage: storage,
	}
}

func (s *itemService) AllItemsByFolder(ctx context.Context, m *dto.AllItemsByFolderDTO) ([]*entity.Item, error) {
	items, err := s.storage.ItemStorage.AllItemsByFolder(
		ctx,
		m.FolderID,
		m.Limit,
		m.ExternalID,
		m.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *itemService) AllItems(ctx context.Context, m *dto.AllItemDTO) ([]*entity.Item, error) {
	items, err := s.storage.ItemStorage.AllItems(ctx, m.Username, m.Limit, m.ExternalID, m.CreatedAt)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *itemService) InsertUpdateItem(ctx context.Context, m *dto.CreateItemDTO) (*entity.Item, error) {
	itemDB := entity.Item{}
	if err := smapping.FillStruct(&itemDB, smapping.MapFields(m)); err != nil {
		return nil, err
	}
	item, err := s.storage.ItemStorage.InsertUpdateItem(ctx, &itemDB)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s *itemService) DeleteItem(ctx context.Context, m *dto.DeleteItemDTO) error {
	if err := s.storage.ItemStorage.DeleteItem(ctx, m.ItemID); err != nil {
		return err
	}
	return nil
}
