package api

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	db "github.com/yuugure-aikouka/kyoto-common/db/store"
	"github.com/yuugure-aikouka/kyoto-common/utils"
)

type Server struct {
	config utils.Config
	store  db.Store
	router *echo.Echo
}

func NewServer(config utils.Config, store db.Store) *Server {
	server := &Server{
		config: config,
		store:  store,
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

	// Validator
	e.Validator = &Validator{validator: validator.New()}

	// Routes
	rg := e.Group("/v1")
	rg.GET("/health", s.healthCheckHandler)
	rg.GET("/users/:id/partners", s.getPartnersHandler)
	rg.GET("/users/:id/potential-partners", s.getPotentialPartnersHandler)

	s.router = e
}

func (s *Server) Start() error {
	return s.router.Start(s.config.Addr)
}

func (s *Server) Route() *echo.Echo {
	return s.router
}
