package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index:,unique,composite:email_deleted_at;index:,unique,composite:username_deleted_at"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Email     string         `json:"email" gorm:"index:,unique,composite:email_deleted_at"`
	UserName  string         `json:"username" gorm:"index:,unique,composite:username_deleted_at"`
	Password  string         `json:"password"`
}

type UserFriend struct {
	ID          uint `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index:,unique,composite:request_accept"`
	RequestId   int            `json:"request_id" gorm:"index:,unique,composite:request_accept"`
	RequestUser *User          `json:"request_user" gorm:"foreignKey:RequestId"`
	AcceptId    int            `json:"accept_id" gorm:"index:,unique,composite:request_accept"`
	AcceptUser  *User          `json:"accept_user" gorm:"foreignKey:AcceptId"`
	IsAccepted  bool           `json:"is_accepted"`
}
