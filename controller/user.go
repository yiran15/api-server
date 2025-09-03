package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yiran15/api-server/base/constant"
	v1 "github.com/yiran15/api-server/service/v1"
)

type UserController interface {
	UserLoginController(c *gin.Context)
	UserLogoutController(c *gin.Context)
	UserCreateController(c *gin.Context)
	UserUpdateByAdminController(c *gin.Context)
	UserUpdateBySelfController(c *gin.Context)
	UserDeleteController(c *gin.Context)
	UserQueryController(c *gin.Context)
	UserListController(c *gin.Context)
	UserInfoController(c *gin.Context)
	OAuth2LoginController(c *gin.Context)
	OAuth2CallbackController(c *gin.Context)
	OAuth2ProviderController(c *gin.Context)
	OAuth2ActivateController(c *gin.Context)
}

type UserControllerImpl struct {
	userServicer v1.UserServicer
}

func NewUserController(userServicer v1.UserServicer) UserController {
	return &UserControllerImpl{
		userServicer: userServicer,
	}
}

// UserLoginController 用户登录
// @Summary 用户登录
// @Description 使用邮箱和密码登录，返回用户信息和 Token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body apitypes.UserLoginRequest true "登录请求参数"
// @Success 200 {object} apitypes.Response{data=apitypes.UserLoginResponse} "登录成功"
// @Router /api/v1/users/login [post]
func (u *UserControllerImpl) UserLoginController(c *gin.Context) {
	ResponseWithData(c, u.userServicer.Login, bindTypeJson)
}

// UserLogoutController 用户注销
// @Summary 用户注销
// @Description 用户注销，清空 Token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} apitypes.Response "注销成功"
// @Router /api/v1/user/logout [post]
func (u *UserControllerImpl) UserLogoutController(c *gin.Context) {
	ResponseNoBind(c, u.userServicer.Logout)
}

// UserCreateController 用户创建
// @Summary 用户创建
// @Description 创建用户同时可以设置角色
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body apitypes.UserCreateRequest true "创建请求参数"
// @Success 200 {object} apitypes.Response "创建成功"
// @Router /api/v1/user/register [post]
func (u *UserControllerImpl) UserCreateController(c *gin.Context) {
	ResponseOnlySuccess(c, u.userServicer.CreateUser, bindTypeJson)
}

// UserUpdateByAdminController 用户更新
// @Summary 用户更新
// @Description 更新用户信息，可以更新角色
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body apitypes.UserUpdateAdminRequest true "更新请求参数"
// @Success 200 {object} apitypes.Response "更新成功"
// @Router /api/v1/user/:id [put]
func (u *UserControllerImpl) UserUpdateByAdminController(c *gin.Context) {
	ResponseOnlySuccess(c, u.userServicer.UpdateUserByAdmin, bindTypeUri, bindTypeJson)
}

// UserUpdateBySelfController 用户更新自己的信息
// @Summary 用户更新自己的信息
// @Description 更新用户信息，不能更新角色
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body apitypes.UserUpdateSelfRequest true "更新请求参数"
// @Success 200 {object} apitypes.Response "更新成功"
// @Router /api/v1/user/self [put]
func (u *UserControllerImpl) UserUpdateBySelfController(c *gin.Context) {
	ResponseOnlySuccess(c, u.userServicer.UpdateUserBySelf, bindTypeJson)
}

// UserDeleteController 用户删除
// @Summary 用户删除
// @Description 删除用户，只能管理员删除
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body apitypes.IDRequest true "删除请求参数"
// @Success 200 {object} apitypes.Response "删除成功"
// @Router /api/v1/user/:id [delete]
func (u *UserControllerImpl) UserDeleteController(c *gin.Context) {
	ResponseOnlySuccess(c, u.userServicer.DeleteUser, bindTypeUri)
}

// UserQueryController 用户查询
// @Summary 用户查询
// @Description 使用 id 查询用户的信息和用户的角色
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body apitypes.IDRequest true "查询请求参数"
// @Success 200 {object} apitypes.Response{data=model.User} "查询成功"
// @Router /api/v1/user/:id [get]
func (u *UserControllerImpl) UserQueryController(c *gin.Context) {
	ResponseWithData(c, u.userServicer.QueryUser, bindTypeUri)
}

// UserInfoController 用户获取自己的信息
// @Summary 用户获取自己的信息
// @Description 使用 id 查询用户的信息和用户的角色
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} apitypes.Response{data=model.User} "查询成功"
// @Router /api/v1/user/info [get]
func (u *UserControllerImpl) UserInfoController(c *gin.Context) {
	ResponseWithDataNoBind(c, u.userServicer.Info)
}

// UserListController 用户列表
// @Summary 用户列表
// @Description 使用分页查询用户的信息, 支持根据 name, email, mobile, department 查询
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data query apitypes.UserListRequest true "查询请求参数"
// @Success 200 {object} apitypes.Response{data=apitypes.UserListResponse} "登录成功"
// @Router /api/v1/user/ [get]
func (u *UserControllerImpl) UserListController(c *gin.Context) {
	ResponseWithData(c, u.userServicer.ListUser, bindTypeQuery)
}

// OAuth2LoginController OAuth 登录
// @Summary OAuth 登录
// @Description 使用 OAuth 登录，返回用户信息和 Token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 302 {string} string "重定向到 OAuth 登录页面"
// @Router /api/v1/oauth2/login [get]
func (u *UserControllerImpl) OAuth2LoginController(c *gin.Context) {
	session := sessions.Default(c)
	state := uuid.New().String()
	session.Set("state", state)
	provider := c.Query("provider")
	if provider != "" {
		session.Set("provider", provider)
	}

	if err := session.Save(); err != nil {
		responseError(c, fmt.Errorf("save session failed: %w", err))
		return
	}
	url, err := u.userServicer.OAuth2Login(provider, state)
	if err != nil {
		responseError(c, err)
		return
	}
	c.Redirect(http.StatusFound, url)
}

// OAuth2CallbackController OAuth 回调
// @Summary OAuth 回调
// @Description 使用 OAuth 回调，返回用户信息和 Token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data query apitypes.OAuthLoginRequest true "回调请求参数"
// @Success 200 {object} apitypes.Response{data=apitypes.UserLoginResponse} "登录成功"
// @Router /api/v1/oauth2/callback [get]
func (u *UserControllerImpl) OAuth2CallbackController(c *gin.Context) {
	session := sessions.Default(c)
	stateSession := session.Get("state")
	providerSession := session.Get("provider")
	state := c.Query("state")
	if state == "" {
		responseError(c, errors.New("state is empty"))
		return
	}
	if state != stateSession {
		responseError(c, errors.New("state invalid"))
		return
	}
	ctx := context.WithValue(c.Request.Context(), constant.ProviderContextKey, providerSession)
	c.Request = c.Request.WithContext(ctx)
	ResponseWithData(c, u.userServicer.OAuth2Callback, bindTypeQuery)
}

// OAuth2ProviderController OAuth2 提供商列表
// @Summary OAuth2 提供商列表
// @Description 获取 OAuth2 提供商列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} apitypes.Response{data=[]string} "获取成功"
// @Router /api/v1/oauth2/provider [get]
func (u *UserControllerImpl) OAuth2ProviderController(c *gin.Context) {
	ResponseWithDataNoBind(c, u.userServicer.OAuth2Provider)
}

// OAuth2ActivateController OAuth2 激活
// @Summary OAuth2 激活
// @Description 使用 OAuth2 激活，返回用户信息和 Token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body apitypes.OAuthActivateRequest true "激活请求参数"
// @Success 200 {object} apitypes.Response{data=apitypes.UserLoginResponse} "激活成功"
// @Router /api/v1/oauth2/:id [post]
func (u *UserControllerImpl) OAuth2ActivateController(c *gin.Context) {
	session := sessions.Default(c)
	stateSession := session.Get("state")
	state := c.Query("state")
	if state == "" {
		responseError(c, errors.New("state is empty"))
		return
	}
	if state != stateSession {
		responseError(c, errors.New("state invalid"))
		return
	}
	ResponseWithData(c, u.userServicer.OAuth2Activate, bindTypeUri, bindTypeJson)
}
