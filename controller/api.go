package controller

import (
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/yiran15/api-server/base/apitypes"
	"github.com/yiran15/api-server/base/constant"
	v1 "github.com/yiran15/api-server/service/v1"
)

type ApiController interface {
	CreateApi(c *gin.Context)
	UpdateApi(c *gin.Context)
	DeleteApi(c *gin.Context)
	QueryApi(c *gin.Context)
	ListApi(c *gin.Context)
	GetServerApi(c *gin.Context)
}

type apiController struct {
	apiService v1.ApiServicer
}

func NewApiController(apiService v1.ApiServicer) ApiController {
	return &apiController{
		apiService: apiService,
	}
}

// CreateApi 创建 API
// @Summary 创建 API
// @Description 创建 API
// @Tags API管理
// @Accept json
// @Produce json
// @Param data body apitypes.ApiCreateRequest true "创建请求参数"
// @Success 200 {object} apitypes.Response "创建成功"
// @Router /api/v1/api [post]
func (receiver *apiController) CreateApi(c *gin.Context) {
	ResponseOnlySuccess(c, receiver.apiService.CreateApi, bindTypeJson)
}

// UpdateApi 更新 API
// @Summary 更新 API
// @Description 更新 API
// @Tags API管理
// @Accept json
// @Produce json
// @Param data body apitypes.ApiUpdateRequest true "更新请求参数"
// @Success 200 {object} apitypes.Response "更新成功"
// @Router /api/v1/api/:id [put]
func (receiver *apiController) UpdateApi(c *gin.Context) {
	ResponseOnlySuccess(c, receiver.apiService.UpdateApi, bindTypeUri, bindTypeJson)
}

// DeleteApi 删除 API
// @Summary 删除 API
// @Description 删除 API
// @Tags API管理
// @Accept json
// @Produce json
// @Param data body apitypes.IDRequest true "删除请求参数"
// @Success 200 {object} apitypes.Response "删除成功"
// @Router /api/v1/api/:id [delete]
func (receiver *apiController) DeleteApi(c *gin.Context) {
	ResponseOnlySuccess(c, receiver.apiService.DeleteApi, bindTypeUri)
}

// QueryApi 查询 API
// @Summary 查询 API
// @Description 查询 API
// @Tags API管理
// @Accept json
// @Produce json
// @Param data body apitypes.IDRequest true "查询请求参数"
// @Success 200 {object} apitypes.Response{data=model.Api} "查询成功"
// @Router /api/v1/api/:id [get]
func (receiver *apiController) QueryApi(c *gin.Context) {
	ResponseWithData(c, receiver.apiService.QueryApi, bindTypeUri)
}

// ListApi API列表
// @Summary API列表
// @Description 使用分页查询 API 的信息, 支持根据 name 查询
// @Tags API管理
// @Accept json
// @Produce json
// @Param data query apitypes.ApiListRequest true "查询请求参数"
// @Success 200 {object} apitypes.Response{data=apitypes.ApiListResponse} "查询成功"
// @Router /api/v1/api/ [get]
func (receiver *apiController) ListApi(c *gin.Context) {
	ResponseWithData(c, receiver.apiService.ListApi, bindTypeUri, bindTypeQuery)
}

// @Summary 获取所有api
// @Description 获取所有api
// @Tags API管理
// @Accept json
// @Produce json
// @Success 200 {object} apitypes.Response{data=apitypes.ServerApiData} "查询成功"
// @Router /api/v1/api/serverApi [get]
func (receiver *apiController) GetServerApi(c *gin.Context) {
	c.JSON(http.StatusOK, apitypes.NewResponse(http.StatusOK, "success", requestid.Get(c), constant.ApiData, nil))
}
