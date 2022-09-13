package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	ID          uint64    `gorm:"primary_key:auto_increment" json:"id"`
	Username    string    `gorm:"column:username;unique" json:"username"`
	Password    string    `gorm:"column:password" json:"password"`
	FirstName   string    `gorm:"column:first_name" json:"first_name"`
	LastName    string    `gorm:"column:last_name" json:"last_name"`
	CreatedAt   time.Time `gorm:"column:create_at" json:"create_at"`
	Avatar      string    `gorm:"column:avatar" json:"avatar"`
	IsSuperuser bool      `gorm:"column:is_superuser;default:false" json:"is_superuser"`
	IsStaff     bool      `gorm:"column:is_staff;default:false" json:"is_staff"`
	Folders     []*Folder `json:"folders"`
}
