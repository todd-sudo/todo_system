package entity

import "time"

type Item struct {
	ID          uint64    `gorm:"primary_key:auto_increment" json:"id"`
	ExternalID  string    `gorm:"type:varchar(255)" json:"external_id"`
	Title       string    `gorm:"type:varchar(255)" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	FolderID    uint64    `gorm:"not null" json:"folder_id"`
	Folder      Folder    `gorm:"foreignkey:FolderID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"folder"`
}
