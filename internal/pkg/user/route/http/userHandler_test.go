package http

import (
	"backend/config"
	"backend/internal/pkg/model"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserSqlite(t *testing.T) {

	// test 완료 후 기존 file db 삭제
	os.Remove("./gorm.db")

	// test를 위한 echo/cfg/handelr 생성 및 설정
	e := echo.New()
	// sqlite db path config에서 불러올 수 없어서 강제 지정
	cfg := config.Config{
		SqliteDBPath: "./gorm.db",
	}
	// user handler 생성
	h := NewUserHandler(e, cfg)

	// 1. test createUser
	// user 생성, domain package 의 json 형태로 변경 가능한 User struct
	user := &model.User{
		Name:     "a",
		Age:      18,
		Birthday: time.Now(),
	}
	// json marshaling
	body, _ := json.Marshal(user)

	// json body 추가
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	// content-type
	req.Header.Set("Content-Type", "application/json")

	// echo context 생성
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, h.createUser(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}

	user = &model.User{
		Name:     "b",
		Age:      19,
		Birthday: time.Now(),
	}
	// json marshaling
	body, _ = json.Marshal(user)

	// json body 추가
	req = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	rec = httptest.NewRecorder()

	// content-type
	req.Header.Set("Content-Type", "application/json")

	// echo context 생성
	c = e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, h.createUser(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}

	// 2. 조회 테스트(all)
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()

	c = e.NewContext(req, rec)

	if assert.NoError(t, h.getUsers(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// rec 에서 읽어올 struct 생성
	users := []*model.User{}

	// json decoder를 이용하여 decoding 반환 값이 domain user의 slice형태
	err := json.NewDecoder(rec.Body).Decode(&users)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(users))

	for _, v := range users {
		t.Log(v.ID, v.Name, v.Age, v.Birthday)
	}

	// 3. 조회 테스트 by id
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()

	c = e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, h.getUsersById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// rec 에서 읽어올 struct 생성
	user = &model.User{}

	// json decoder를 이용하여 decoding 반환 값이 domain user
	err = json.NewDecoder(rec.Body).Decode(user)
	assert.NoError(t, err)
	assert.Equal(t, "a", user.Name)
	assert.Equal(t, 18, int(user.Age))

	// for _, v := range users {
	// 	t.Log(v.ID, v.Name, v.Age, v.Birthday)
	// }

	// 4. 삭제 테스트 by id
	req = httptest.NewRequest(http.MethodDelete, "/", nil)
	rec = httptest.NewRecorder()

	c = e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, h.deleteUsersById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// 5. 업데이트 테스트 by id

	body, _ = json.Marshal(&model.User{Name: "bbb", Age: 45, Birthday: time.Now()})

	req = httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(body))
	rec = httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	c = e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("2")

	if assert.NoError(t, h.updateUserById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// update 후 변경 내역 확인

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()

	c = e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("2")

	if assert.NoError(t, h.getUsersById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// rec 에서 읽어올 struct 생성
	user = &model.User{}

	// json decoder를 이용하여 decoding 반환 값이 domain user
	err = json.NewDecoder(rec.Body).Decode(user)
	assert.NoError(t, err)
	assert.Equal(t, "bbb", user.Name)
	assert.Equal(t, 45, int(user.Age))
}

func TestUserPostgres(t *testing.T) {

	// test 완료 후 기존 file db 삭제

	e := echo.New()
	// sqlite db path config에서 불러올 수 없어서 강제 지정
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
		DB: "postgre",
	}

	// user handler 생성
	h := NewUserHandler(e, cfg)

	// 1. test createUser
	// user 생성, domain package 의 json 형태로 변경 가능한 User struct
	user := &model.User{
		Name:     "a",
		Age:      18,
		Birthday: time.Now(),
	}
	// json marshaling
	body, _ := json.Marshal(user)

	// json body 추가
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	// content-type
	req.Header.Set("Content-Type", "application/json")

	// echo context 생성
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, h.createUser(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}

	user = &model.User{
		Name:     "b",
		Age:      19,
		Birthday: time.Now(),
	}
	// json marshaling
	body, _ = json.Marshal(user)

	// json body 추가
	req = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	rec = httptest.NewRecorder()

	// content-type
	req.Header.Set("Content-Type", "application/json")

	// echo context 생성
	c = e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, h.createUser(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}

	// 2. 조회 테스트(all)
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()

	c = e.NewContext(req, rec)

	if assert.NoError(t, h.getUsers(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// rec 에서 읽어올 struct 생성
	users := []*model.User{}

	// json decoder를 이용하여 decoding 반환 값이 domain user의 slice형태
	err := json.NewDecoder(rec.Body).Decode(&users)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(users))

	for _, v := range users {
		t.Log(v.ID, v.Name, v.Age, v.Birthday)
	}

	// 3. 조회 테스트 by id
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()

	c = e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, h.getUsersById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// rec 에서 읽어올 struct 생성
	user = &model.User{}

	// json decoder를 이용하여 decoding 반환 값이 domain user
	err = json.NewDecoder(rec.Body).Decode(user)
	assert.NoError(t, err)
	assert.Equal(t, "a", user.Name)
	assert.Equal(t, 18, int(user.Age))

	// for _, v := range users {
	// 	t.Log(v.ID, v.Name, v.Age, v.Birthday)
	// }

	// 4. 삭제 테스트 by id
	req = httptest.NewRequest(http.MethodDelete, "/", nil)
	rec = httptest.NewRecorder()

	c = e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, h.deleteUsersById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// 5. 업데이트 테스트 by id

	body, _ = json.Marshal(&model.User{Name: "bbb", Age: 45, Birthday: time.Now()})

	req = httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(body))
	rec = httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	c = e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("2")

	if assert.NoError(t, h.updateUserById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// update 후 변경 내역 확인

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()

	c = e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("2")

	if assert.NoError(t, h.getUsersById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// rec 에서 읽어올 struct 생성
	user = &model.User{}

	// json decoder를 이용하여 decoding 반환 값이 domain user
	err = json.NewDecoder(rec.Body).Decode(user)
	assert.NoError(t, err)
	assert.Equal(t, "bbb", user.Name)
	assert.Equal(t, 45, int(user.Age))
}
