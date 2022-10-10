package entity

type Item struct {
	ID          uint64 `gorm:"primary_key:auto_increment" json:"id"`
	Index       uint64 `gorm:"auto_increment" json:"index"`
	Title       string `gorm:"type:varchar(255)" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	FolderID    uint64 `gorm:"not null" json:"folder_id"`
	Folder      Folder `gorm:"foreignkey:FolderID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"folder"`
}
