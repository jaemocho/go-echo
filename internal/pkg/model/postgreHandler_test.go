package model

import (
	"backend/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPostgreHandler(t *testing.T) {
	assert := assert.New(t)

	var cnt int64

	cfg := config.Config{
		Postgre: config.Postgre{
			IP:       "127.0.0.1",
			Port:     "5432",
			DBName:   "postgres",
			User:     "postgres",
			Password: "echogorm",
			SSLMode:  "disable",
			TimeZone: "Asia/Seoul",
		},
	}
	h := NewPostgreHandler(cfg)

	cnt = h.AddUser(&User{Name: "a", Age: 38, Birthday: time.Now()})
	assert.Equal(1, int(cnt))

	users := h.GetUsers()
	assert.Equal("a", users[0].Name)
	assert.Equal(38, users[0].Age)

	cnt = h.AddUser(&User{Name: "b", Age: 38, Birthday: time.Now()})
	assert.Equal(1, int(cnt))
	cnt = h.AddUser(&User{Name: "c", Age: 38, Birthday: time.Now()})
	assert.Equal(1, int(cnt))

	users = h.GetUsers()
	assert.Equal(3, len(users))

	user := h.GetUserById(1)
	assert.Equal("a", user.Name)

	cnt = h.DeleteUserById(1)
	assert.Equal(1, int(cnt))

	users = h.GetUsers()
	assert.Equal(2, len(users))

	user = h.GetUserById(1)
	assert.Equal(uint(0), user.ID)
	assert.Equal("", user.Name)

	cnt = h.UpdateUserById(2, &User{Name: "bbb", Age: 40, Birthday: time.Now()})
	assert.Equal(1, int(cnt))

	user = h.GetUserById(2)
	assert.Equal("bbb", user.Name)
	assert.Equal(40, user.Age)

}
