package dto

type AllItemsByFolderDTO struct {
	FolderID   int    `json:"folder_id"`
	Limit      int    `json:"limit"`
	ExternalID string `json:"external_id"`
	CreatedAt  string `json:"created_at"`
}

type AllItemDTO struct {
	Username   string `json:"username"`
	Limit      int    `json:"limit"`
	ExternalID string `json:"external_id"`
	CreatedAt  string `json:"created_at"`
}

type CreateItemDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	FolderID    int    `json:"folder_id"`
}

type DeleteItemDTO struct {
	ItemID int `json:"item_id"`
}
