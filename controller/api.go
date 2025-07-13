package controller

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/yiran15/api-server/service/v1"
)

type ApiController interface {
	CreateApi(c *gin.Context)
	UpdateApi(c *gin.Context)
	DeleteApi(c *gin.Context)
	QueryApi(c *gin.Context)
	ListApi(c *gin.Context)
}

type apiController struct {
	apiService v1.ApiServicer
}

func NewApiController(apiService v1.ApiServicer) ApiController {
	return &apiController{
		apiService: apiService,
	}
}

func (s *apiController) CreateApi(c *gin.Context) {
	ResponseOnlySuccess(c, s.apiService.CreateApi, bindTypeJson)
}

func (s *apiController) UpdateApi(c *gin.Context) {
	ResponseOnlySuccess(c, s.apiService.UpdateApi, bindTypeUri, bindTypeJson)
}

func (s *apiController) DeleteApi(c *gin.Context) {
	ResponseOnlySuccess(c, s.apiService.DeleteApi, bindTypeUri)
}

func (s *apiController) QueryApi(c *gin.Context) {
	ResponseWithData(c, s.apiService.QueryApi, bindTypeUri)
}

func (s *apiController) ListApi(c *gin.Context) {
	ResponseWithData(c, s.apiService.ListApi, bindTypeUri, bindTypeQuery)
}
