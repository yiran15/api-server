package server

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yiran15/api-server/base/apitypes"
	"github.com/yiran15/api-server/base/conf"
	"github.com/yiran15/api-server/base/constant"
	"github.com/yiran15/api-server/base/router"
	"github.com/yiran15/api-server/controller"
	"go.uber.org/zap"
)

const (
	defaultShutdownTimeout = 30 * time.Second
)

type ServerInterface interface {
	Start() error
	Stop() error
}

type Server struct {
	shutdown time.Duration
	server   *http.Server
}

func NewServer(server *gin.Engine) *Server {
	return &Server{
		shutdown: defaultShutdownTimeout,
		server: &http.Server{
			Addr:    conf.GetServerBind(),
			Handler: server,
		},
	}
}

func (s *Server) Start() (err error) {
	zap.S().Infof("start server, addr: %s", s.server.Addr)
	if err = s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdown)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func NewHttpServer(r router.RouterInterface) (*gin.Engine, error) {
	if conf.GetLogLevel() == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	controller.NewValidator()

	r.RegisterRouter(engine)
	var apiData apitypes.ServerApiData
	apiData.ApiInfo = make(map[string][]apitypes.ApiInfo)
	for _, v := range engine.Routes() {
		if v.Path == "/swagger/*any" || v.Path == "/oauth2/login" || v.Path == "/oauth2/callback" || v.Path == "/oauth2/provider" {
			continue
		}
		api := strings.TrimPrefix(v.Path, "/")
		apiType := strings.Split(api, "/")[2]
		_, ok := apiData.ApiInfo[apiType]
		if !ok {
			apiData.ApiInfo[apiType] = make([]apitypes.ApiInfo, 0)
			apiData.ApiType = append(apiData.ApiType, apiType)
		}
		apiData.ApiInfo[apiType] = append(apiData.ApiInfo[apiType], apitypes.ApiInfo{
			Method:  v.Method,
			Path:    v.Path,
			Handler: v.Handler,
		})
	}
	constant.ApiData = apiData
	return engine, nil
}
