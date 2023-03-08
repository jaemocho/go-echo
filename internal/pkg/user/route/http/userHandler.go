package http

import (
	"backend/config"
	model "backend/internal/pkg/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	db model.DBHandler
}

func NewUserHandler(echo *echo.Echo, cfg config.Config) *UserHandler {

	handler := &UserHandler{
		db: model.NewDBHandler(cfg),
	}

	user := echo.Group("/api/v1/user")
	{
		user.GET("", handler.getUsers)
		user.GET("/:id", handler.getUsersById)
		user.POST("", handler.createUser)
		user.DELETE("/:id", handler.deleteUsersById)
		user.PUT("/:id", handler.updateUserById)
	}

	return handler
}

// @Summary		Get users
// @Description	Get all user's info
// @name		getUsers
// @Accept		json
// @Produce		json
// @Success		200	{array}	model.User
// @Router		/api/v1/user [get]
// @Security    ApiKeyAuth
func (u *UserHandler) getUsers(c echo.Context) error {

	users := u.db.GetUsers()

	return c.JSON(http.StatusOK, users)
}

// @Summary		Get user by id
// @Description	Get user's info
// @name		getUsersById
// @Accept		json
// @Produce		json
// @Param		id	path		string	true	"id of the user"
// @Success		200	{object}	model.User
// @Router		/api/v1/user/{id} [get]
// @Security    ApiKeyAuth
func (u *UserHandler) getUsersById(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Error(err)
		return c.JSON(http.StatusBadRequest, c.Param("id"))
	}

	user := u.db.GetUserById(id)

	return c.JSON(http.StatusOK, user)
}

// @Summary		Create user
// @Description	Create new user
// @name		createUser
// @Accept		json
// @Produce		json
// @Param		userBody	body	model.User	true	"User Info Body"
// @Success		201
// @Router		/api/v1/user [post]
// @Security    ApiKeyAuth
func (u *UserHandler) createUser(c echo.Context) error {

	user := new(model.User)

	if err := c.Bind(user); err != nil {
		c.Error(err)
		return c.JSON(http.StatusBadRequest, nil)
	}

	if cnt := u.db.AddUser(user); cnt == 1 {
		return c.JSON(http.StatusCreated, nil)
	}

	return c.JSON(http.StatusBadRequest, nil)

}

// @Summary		delete user by id
// @Description	delete user's info
// @name		deleteUsersById
// @Accept		json
// @Produce		json
// @Param		id	path	string	true	"id of the user"
// @Success		200
// @Router		/api/v1/user/{id} [delete]
// @Security    ApiKeyAuth
func (u *UserHandler) deleteUsersById(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Error(err)
		return c.JSON(http.StatusBadRequest, c.Param("id"))
	}

	if cnt := u.db.DeleteUserById(id); cnt == 1 {
		return c.JSON(http.StatusOK, nil)
	}
	return c.JSON(http.StatusBadRequest, nil)
}

// @Summary		update user by id
// @Description	update user's info
// @name		updateUserById
// @Accept		json
// @Produce		json
// @Param		id			path	string		true	"id of the user"
// @Param		userBody	body	model.User	true	"User Info Body"
// @Success		200
// @Router		/api/v1/user/{id} [put]
// @Security    ApiKeyAuth
func (u *UserHandler) updateUserById(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Error(err)
		return c.JSON(http.StatusBadRequest, c.Param("id"))
	}

	user := new(model.User)

	if err := c.Bind(user); err != nil {
		c.Error(err)
		return c.JSON(http.StatusBadRequest, nil)
	}

	if cnt := u.db.UpdateUserById(id, user); cnt == 1 {
		return c.JSON(http.StatusOK, nil)
	}
	return c.JSON(http.StatusInternalServerError, nil)
}
