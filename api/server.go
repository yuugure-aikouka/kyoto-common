package api

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	Addr string
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
	e.Server.ReadTimeout = 10 * time.Second
	e.Server.WriteTimeout = 30 * time.Second
	e.Server.IdleTimeout = time.Minute

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
