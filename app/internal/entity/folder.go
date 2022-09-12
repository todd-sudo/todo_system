package entity

import "gorm.io/gorm"

type Folder struct {
	gorm.Model

	ID     uint64  `gorm:"primary_key:auto_increment" json:"id"`
	Name   string  `gorm:"column:name" json:"name"`
	UserId uint64  `gorm:"not null" json:"user_id"`
	User   *User   `gorm:"foreignkey:UserID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"user"`
	Items  []*Item `json:"items"`
}
