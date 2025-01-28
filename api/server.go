package api

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type Server struct {
	config Config
	router *echo.Echo
}

func NewServer(config Config) *Server {
	server := &Server{
		config: config,
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
	rg.GET("/health", s.healthCheckHandler)

	s.router = e
}

func (s *Server) Start() error {
	return s.router.Start(s.config.Addr)
}
