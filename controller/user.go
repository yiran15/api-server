package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/yiran15/api-server/service"
)

type UserController interface {
	UserCreate(c *gin.Context) error
	UserUpdate(c *gin.Context) error
	UserDelete(c *gin.Context) error
	UserQuery(c *gin.Context) error
	UserList(c *gin.Context) error
}

type UserControllerImpl struct {
	userServicer service.UserServicer
}

func NewUserController(userServicer service.UserServicer) UserController {
	return &UserControllerImpl{
		userServicer: userServicer,
	}
}

func (u *UserControllerImpl) UserCreate(c *gin.Context) error {
	return nil
}

func (u *UserControllerImpl) UserUpdate(c *gin.Context) error {
	return nil
}

func (u *UserControllerImpl) UserDelete(c *gin.Context) error {
	return nil
}

func (u *UserControllerImpl) UserQuery(c *gin.Context) error {
	return nil
}

func (u *UserControllerImpl) UserList(c *gin.Context) error {
	return nil
}
