package model

import (
	"backend/config"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string    `json:"name"`
	Age      int       `json:"age"`
	Birthday time.Time `json:"birthday"`
}

type DBHandler interface {
	GetUsers() []*User
	AddUser(user *User) int64
	GetUserById(id int) *User
	DeleteUserById(id int) int64
	UpdateUserById(id int, user *User) int64
}

func NewDBHandler(cfg config.Config) DBHandler {
	return NewSqliteHandler(cfg)
}
