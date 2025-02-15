package api

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yuugure-aikouka/kyoto-common/config"
	"github.com/yuugure-aikouka/kyoto-common/handler"
)

type Server struct {
	config  config.Config
	handler *handler.Handler
	router  *echo.Echo
}

func NewServer(config config.Config, handler *handler.Handler) *Server {
	server := &Server{
		config:  config,
		handler: handler,
	}

	server.setupRouter()

	return server
}

func (s *Server) setupRouter() {
	e := echo.New()

	// Server settings
	e.Server.ReadTimeout = s.config.ReadTimeout * time.Second
	e.Server.WriteTimeout = s.config.WriteTimeout * time.Second
	e.Server.IdleTimeout = s.config.IdleTimeout * time.Second

	// Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	rg := e.Group("/v1")
	rg.GET("/health", s.handler.HealthCheck)
	rg.GET("/users/:id/partners", s.handler.GetPartners)
	rg.GET("/users/:id/potential-partners", s.handler.GetPotentialPartners)

	s.router = e
}

func (s *Server) Start() error {
	return s.router.Start(s.config.Addr)
}

func (s *Server) Route() *echo.Echo {
	return s.router
}
