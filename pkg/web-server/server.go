package web_server

import (
	"context"
	"errors"
	"fmt"
	loggerService "github.com/Miroshinsv/disko_go/pkg/logger-service"
	"github.com/gorilla/mux"
	"net/http"
)

var self IWebServer = nil

type Server struct {
	config    *Config
	server    *http.Server
	log       loggerService.ILogger
	isRunning bool
}

func (s *Server) RegisterRoutes(router *mux.Router) {
	s.server = &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
	}

	s.log.Debug("Registered new routes to web server", nil)
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	if s.server == nil {
		return errors.New("no routes defined for web-server")
	}

	if s.isRunning {
		return errors.New("web-server is already running")
	}

	var errCh = make(chan error, 1)

	go func(ec chan error) {
		ec <- s.server.ListenAndServe()
	}(errCh)

	s.isRunning = true

	s.log.Info("Web server has started", map[string]interface{}{
		"host": s.config.Host,
		"port": s.config.Port,
	})

	select {
	case err := <-errCh:
		if err == http.ErrServerClosed {
			return nil
		}

		s.log.Error("error on serving web application", err, nil)

		return err
	case <-ctx.Done():
		s.log.Warning("terminating web-server by context closing", nil)

		return nil

	}
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	if !s.isRunning {
		return errors.New("web-server was already stopped")
	}

	err := s.server.Shutdown(ctx)
	s.isRunning = false
	s.server = nil

	s.log.Info("Web server has stopped", nil)

	return err
}

func (s Server) IsRunning() bool {
	return s.isRunning
}

func MustNewWebServer(conf *Config, log loggerService.ILogger) IWebServer {
	if self != nil {
		panic("web-server already defined")
	}

	self = &Server{
		config: conf,
		log:    log,
	}

	return self
}

func GetWebServer() (IWebServer, error) {
	if self == nil {
		return &Server{}, errors.New("web-server never defined")
	}

	return self, nil
}
