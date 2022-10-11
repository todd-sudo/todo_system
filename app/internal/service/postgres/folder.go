package service_pg

import (
	"context"

	"github.com/mashingan/smapping"
	"github.com/todd-sudo/todo_system/internal/dto"
	"github.com/todd-sudo/todo_system/internal/entity"
	"github.com/todd-sudo/todo_system/internal/storage/postgres"
	"github.com/todd-sudo/todo_system/pkg/logging"
)

type FolderService interface {
	AllFolders(ctx context.Context, username string) ([]*entity.Folder, error)
	InsertUpdateFolder(ctx context.Context, folderDTO *dto.CreateUpdateFolderDTO) (*entity.Folder, error)
	DeleteFolder(ctx context.Context, folderDTO *dto.DeleteFolderDTO) error
}

type folderService struct {
	ctx     context.Context
	storage *postgres.Storage
	log     logging.Logger
}

func NewFolderService(ctx context.Context, log logging.Logger, storage *postgres.Storage) FolderService {
	return &folderService{
		ctx:     ctx,
		log:     log,
		storage: storage,
	}
}

func (s *folderService) AllFolders(ctx context.Context, username string) ([]*entity.Folder, error) {
	folders, err := s.storage.FolderStorage.AllFolders(ctx, username)
	if err != nil {
		return nil, err
	}
	return folders, nil
}

func (s *folderService) InsertUpdateFolder(ctx context.Context, folderDTO *dto.CreateUpdateFolderDTO) (*entity.Folder, error) {
	folderDB := entity.Folder{}
	if err := smapping.FillStruct(&folderDB, smapping.MapFields(folderDTO)); err != nil {
		return nil, err
	}
	folder, err := s.storage.FolderStorage.InsertUpdateFolder(ctx, &folderDB)
	if err != nil {
		return nil, err
	}
	return folder, nil
}

func (s *folderService) DeleteFolder(ctx context.Context, folderDTO *dto.DeleteFolderDTO) error {
	if err := s.storage.FolderStorage.DeleteFolder(ctx, int(folderDTO.FolderID)); err != nil {
		return err
	}
	return nil
}
