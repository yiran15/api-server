package controller

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/yiran15/api-server/service/v1"
)

type UserController interface {
	UserLogin(c *gin.Context)
	UserLogout(c *gin.Context)
	UserCreate(c *gin.Context)
	UserUpdate(c *gin.Context)
	UserDelete(c *gin.Context)
	UserQuery(c *gin.Context)
	UserList(c *gin.Context)
	UserUpdateRole(c *gin.Context)
	UserInfo(c *gin.Context)
}

type UserControllerImpl struct {
	userServicer v1.UserServicer
}

func NewUserController(userServicer v1.UserServicer) UserController {
	return &UserControllerImpl{
		userServicer: userServicer,
	}
}

func (u *UserControllerImpl) UserLogin(c *gin.Context) {
	ResponseWithData(c, u.userServicer.Login, bindTypeJson)
}

func (u *UserControllerImpl) UserLogout(c *gin.Context) {
	ResponseNoBind(c, u.userServicer.Logout)
}

func (u *UserControllerImpl) UserCreate(c *gin.Context) {
	ResponseOnlySuccess(c, u.userServicer.CreateUser, bindTypeJson)
}

func (u *UserControllerImpl) UserUpdate(c *gin.Context) {
	ResponseOnlySuccess(c, u.userServicer.UpdateUser, bindTypeUri, bindTypeJson)
}

func (u *UserControllerImpl) UserUpdateRole(c *gin.Context) {
	ResponseOnlySuccess(c, u.userServicer.UpdateRole, bindTypeUri, bindTypeShouldBind)
}

func (u *UserControllerImpl) UserDelete(c *gin.Context) {
	ResponseOnlySuccess(c, u.userServicer.DeleteUser, bindTypeUri)
}

func (u *UserControllerImpl) UserQuery(c *gin.Context) {
	ResponseWithData(c, u.userServicer.QueryUser, bindTypeUri)
}

func (u *UserControllerImpl) UserInfo(c *gin.Context) {
	ResponseWithDataNoBind(c, u.userServicer.Info)
}

func (u *UserControllerImpl) UserList(c *gin.Context) {
	ResponseWithData(c, u.userServicer.ListUser, bindTypeQuery)
}
