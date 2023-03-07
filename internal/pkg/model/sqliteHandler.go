package model

import (
	"backend/config"

	"gorm.io/driver/sqlite" // Sqlite driver based on GGO
	// "github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"gorm.io/gorm"
)

type sqliteHandler struct {
	db *gorm.DB
}

func NewSqliteHandler(cfg config.Config) DBHandler {
	database, err := gorm.Open(sqlite.Open(cfg.SqliteDBPath), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	database.AutoMigrate(&User{})

	return &sqliteHandler{db: database}
}

func (s *sqliteHandler) AddUser(user *User) int64 {
	result := s.db.Create(&user)
	return result.RowsAffected
}

func (s *sqliteHandler) GetUsers() []*User {
	var users []*User
	s.db.Find(&users)
	return users
}

func (s *sqliteHandler) GetUserById(id int) *User {
	var user *User
	s.db.Find(&user, id)
	return user
}

func (s *sqliteHandler) DeleteUserById(id int) int64 {
	var user *User
	result := s.db.Delete(&user, id)
	return result.RowsAffected
}

func (s *sqliteHandler) UpdateUserById(id int, user *User) int64 {
	originUser := s.GetUserById(id)
	if originUser == nil {
		return 0
	}
	result := s.db.Model(&originUser).Updates(&user)
	return result.RowsAffected
}
