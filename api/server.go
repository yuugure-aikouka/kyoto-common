package api

import (
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
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
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	// Routes
	r := e.Group("/v1")
	r.GET("/health", s.handler.HealthCheck)
	r.POST("/auth/magic", s.handler.SendMagicLinkHandler)
	r.GET("/auth/magic", s.handler.VerifyMagicLinkHandler)
	r.POST("/auth/refresh", s.handler.RefreshAccessTokenHandler)
	r.DELETE("/auth/logout", s.handler.LogoutHandler)
	r.GET("/users/:id/partners", s.handler.GetPartners)
	r.GET("/users/:id/potential-partners", s.handler.GetPotentialPartners)

	// TODO: use authMiddleware like this
	// authMiddleware := middleware.AuthMiddleware(s.config.AccessTokenSecretKey)
	// protected := r.Group("", authMiddleware)
	// protected.GET("/users/:id/partners", s.handler.GetPartners)

	s.router = e
}

func (s *Server) Start() error {
	return s.router.Start(s.config.Addr)
}

func (s *Server) Route() *echo.Echo {
	return s.router
}
