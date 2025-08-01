package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yiran15/api-server/base/conf"
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
	if conf.GetServerLogLevel() == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	controller.NewValidator()

	r.RegisterRouter(engine)
	return engine, nil
}
