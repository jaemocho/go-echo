package model

import (
	"backend/config"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

type postgreHandler struct {
	db *gorm.DB
}

func NewPostgreHandler(cfg config.Config) DBHandler {

	database, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "user=" + cfg.Postgre.User + " password=" + cfg.Postgre.Password + " dbname=" + cfg.Postgre.DBName + " port=" + cfg.Postgre.Port + " sslmode=" + cfg.Postgre.SSLMode + " TimeZone=" + cfg.Postgre.TimeZone,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	database.AutoMigrate(&User{})

	return &postgreHandler{db: database}
}

func (p *postgreHandler) AddUser(user *User) int64 {
	result := p.db.Create(&user)
	return result.RowsAffected
}

func (p *postgreHandler) GetUsers() []*User {
	var users []*User
	p.db.Find(&users)
	return users
}

func (p *postgreHandler) GetUserById(id int) *User {
	var user *User
	p.db.Find(&user, id)
	return user
}

func (p *postgreHandler) DeleteUserById(id int) int64 {
	var user *User
	result := p.db.Delete(&user, id)
	return result.RowsAffected
}

func (p *postgreHandler) UpdateUserById(id int, user *User) int64 {
	originUser := p.GetUserById(id)
	if originUser == nil {
		return 0
	}
	result := p.db.Model(&originUser).Updates(&user)
	return result.RowsAffected
}
