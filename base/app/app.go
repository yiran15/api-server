package app

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/yiran15/api-server/base/server"
	"go.uber.org/zap"
)

type Options func(*Application)

func WithServer(server ...server.ServerInterface) Options {
	return func(app *Application) {
		app.servers = append(app.servers, server...)
	}
}

type Application struct {
	servers []server.ServerInterface
	wg      *sync.WaitGroup
}

func newApp(options ...Options) *Application {
	app := &Application{
		wg: &sync.WaitGroup{},
	}
	for _, option := range options {
		option(app)
	}
	return app
}

func NewApplication(e *gin.Engine) *Application {
	return newApp(
		WithServer(
			server.NewServer(e),
		),
	)
}

func (app *Application) Run(ctx context.Context) error {
	if len(app.servers) == 0 {
		return nil
	}
	errCh := make(chan error, 1)
	for _, s := range app.servers {
		go func(s server.ServerInterface) {
			errCh <- s.Start()
		}(s)
	}

	select {
	case err := <-errCh:
		app.Stop()
		return err
	case <-ctx.Done():
		app.Stop()
		return nil
	}
}

func (app *Application) Stop() {
	if len(app.servers) == 0 {
		return
	}
	for _, s := range app.servers {
		app.wg.Add(1)
		go func(s server.ServerInterface) {
			defer app.wg.Done()
			if err := s.Stop(); err != nil {
				zap.S().Errorf("stop server error %v", err)
			}
		}(s)
	}
	app.wg.Wait()
}
