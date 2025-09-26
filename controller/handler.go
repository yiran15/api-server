package controller

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/yiran15/api-server/base/apitypes"
	"github.com/yiran15/api-server/base/constant"
	"gorm.io/gorm"
)

type bindType int

const (
	bindTypeUri bindType = iota
	bindTypeJson
	bindTypeQuery
	bindTypeShouldBind
)

func bindWithSources(c *gin.Context, req any, sources ...bindType) (success bool) {
	for _, src := range sources {
		var err error
		switch src {
		case bindTypeUri:
			err = c.ShouldBindUri(req)
		case bindTypeJson:
			err = c.ShouldBindJSON(req)
		case bindTypeQuery:
			err = c.ShouldBindQuery(req)
		case bindTypeShouldBind:
			err = c.ShouldBind(req)
		default:
			continue
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				responseParamError(c, err, "request body is empty")
				return false
			}
			responseParamError(c, err, translateErrors(err))
			return false
		}
	}

	requestID := requestid.Get(c)
	if requestID != "" {
		ctx := context.WithValue(c.Request.Context(), constant.RequestIDContextKey, requestID)
		c.Request = c.Request.WithContext(ctx)
	}

	return true
}

type HandlerData[T any, R any] func(ctx context.Context, req *T) (R, error)

func ResponseWithData[T any, R any](c *gin.Context, handler HandlerData[T, R], bindType ...bindType) {
	var (
		data R
		err  error
	)
	req := new(T)
	if !bindWithSources(c, req, bindType...) {
		return
	}

	if data, err = handler(c.Request.Context(), req); err != nil {
		responseError(c, err)
		return
	}

	responseSuccess(c, data)
}

type HandlerErr[T any] func(ctx context.Context, req *T) error

func ResponseOnlySuccess[T any](c *gin.Context, handler HandlerErr[T], bindTypes ...bindType) {
	req := new(T)
	if !bindWithSources(c, req, bindTypes...) {
		return
	}

	if err := handler(c.Request.Context(), req); err != nil {
		responseError(c, err)
		return
	}

	responseSuccess(c, nil)
}

type Handler[R any] func(ctx context.Context) (R, error)

func ResponseWithDataNoBind[R any](c *gin.Context, handler Handler[R]) {
	var (
		data R
		err  error
	)
	if data, err = handler(c.Request.Context()); err != nil {
		responseError(c, err)
		return
	}

	responseSuccess(c, data)
}

type HandlerErrNoBind func(ctx context.Context) error

func ResponseNoBind(c *gin.Context, handler HandlerErrNoBind) {
	if err := handler(c.Request.Context()); err != nil {
		responseError(c, err)
		return
	}
	responseSuccess(c, nil)
}

func responseError(c *gin.Context, err error) {
	code, err := getErr(err)
	c.JSON(code, apitypes.NewResponse(code, "", requestid.Get(c), nil, err.Error()))
	c.Error(err)
}

func responseSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, apitypes.NewResponse(0, "success", requestid.Get(c), data, nil))
}

func responseParamError(c *gin.Context, err error, errors any) {
	c.JSON(http.StatusBadRequest, apitypes.NewResponse(http.StatusBadRequest, "parameter error", requestid.Get(c), nil, errors))
	c.Error(err)
}

func getErr(err error) (int, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound, errors.New("object not found")
	}

	if code, ok, err := mysqlErr(err); ok {
		return code, err
	}

	return defaultErr(err)
}

func mysqlErr(err error) (int, bool, error) {
	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		return 0, false, err
	}

	switch mysqlErr.Number {
	case 1062:
		return http.StatusBadRequest, true, errors.New("object already exists")
	default:
		return http.StatusInternalServerError, true, err
	}
}

func defaultErr(err error) (int, error) {
	return http.StatusInternalServerError, err
}
