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

// CreateRole 创建角色
// @Summary 创建角色
// @Description 创建角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param data body apitypes.RoleCreateRequest true "创建请求参数"
// @Success 200 {object} apitypes.Response "创建成功"
// @Router /api/v1/role [post]
func (s *roleController) CreateRole(c *gin.Context) {
	ResponseOnlySuccess(c, s.roleService.CreateRole, bindTypeJson)
}

// UpdateRole 更新角色
// @Summary 更新角色
// @Description 更新角色, 并且可以更新角色的权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param data body apitypes.RoleUpdateRequest true "更新请求参数"
// @Success 200 {object} apitypes.Response "更新成功"
// @Router /api/v1/role/:id [put]
func (s *roleController) UpdateRole(c *gin.Context) {
	ResponseOnlySuccess(c, s.roleService.UpdateRole, bindTypeUri, bindTypeJson)
}

// DeleteRole 删除角色
// @Summary 删除角色
// @Description 删除角色, 不能删除有用户的角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param data body apitypes.IDRequest true "删除请求参数"
// @Success 200 {object} apitypes.Response "删除成功"
// @Router /api/v1/role/:id [delete]
func (s *roleController) DeleteRole(c *gin.Context) {
	ResponseOnlySuccess(c, s.roleService.DeleteRole, bindTypeUri)
}

// QueryRole 查询角色
// @Summary 查询角色
// @Description 查询角色, 包括角色的权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param data body apitypes.IDRequest true "查询请求参数"
// @Success 200 {object} apitypes.Response{data=model.Role} "查询成功"
// @Router /api/v1/role/:id [get]
func (s *roleController) QueryRole(c *gin.Context) {
	ResponseWithData(c, s.roleService.QueryRole, bindTypeUri)
}

// ListRole 角色列表
// @Summary 角色列表
// @Description 使用分页查询角色的信息, 支持根据 name 查询
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param data query apitypes.RoleListRequest true "查询请求参数"
// @Success 200 {object} apitypes.Response{data=apitypes.RoleListResponse} "查询成功"
// @Router /api/v1/role/ [get]
func (s *roleController) ListRole(c *gin.Context) {
	ResponseWithData(c, s.roleService.ListRole, bindTypeUri, bindTypeQuery)
}
