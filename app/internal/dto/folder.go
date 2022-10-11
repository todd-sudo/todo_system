package dto

type CreateUpdateFolderDTO struct {
	Name   string `json:"name"`
	UserID uint64 `json:"user_id"`
}

type DeleteFolderDTO struct {
	FolderID uint64 `json:"folder_id"`
}
