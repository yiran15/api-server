package controller

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/yiran15/api-server/service/v1"
)

type RoleController interface {
	CreateRole(c *gin.Context)
	UpdateRole(c *gin.Context)
	DeleteRole(c *gin.Context)
	QueryRole(c *gin.Context)
	ListRole(c *gin.Context)
}

type roleController struct {
	roleService v1.RoleServicer
}

func NewRoleController(roleService v1.RoleServicer) RoleController {
	return &roleController{
		roleService: roleService,
	}
}

func (s *roleController) CreateRole(c *gin.Context) {
	ResponseOnlySuccess(c, s.roleService.CreateRole, bindTypeJson)
}

func (s *roleController) UpdateRole(c *gin.Context) {
	ResponseOnlySuccess(c, s.roleService.UpdateRole, bindTypeUri, bindTypeJson)
}

func (s *roleController) DeleteRole(c *gin.Context) {
	ResponseOnlySuccess(c, s.roleService.DeleteRole, bindTypeUri)
}

func (s *roleController) QueryRole(c *gin.Context) {
	ResponseWithData(c, s.roleService.QueryRole, bindTypeUri)
}

func (s *roleController) ListRole(c *gin.Context) {
	ResponseWithData(c, s.roleService.ListRole, bindTypeUri, bindTypeQuery)
}
