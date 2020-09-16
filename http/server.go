package http

import (
	"context"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/nyks06/backapi"
	"github.com/spf13/viper"
)

const (
	// ContextKeyCurrentUser defines the name of the key in our echo.Context to get CurrentUser if any
	ContextKeyCurrentUser = "CUser_"
	// ContextKeySession defines the name of the key in our echo.Contexdt to get Session of the CurrentUser if any
	ContextKeySession = "Session_"
	// HeaderKeyAuthorization defines the name of the header that should be set in the request to authentify the user
	HeaderKeyAuthorization = "X-Authorization"
)

type Server struct {
	Router *echo.Echo
	Port   string

	Middlewares    *middlewares
	UserService    *backapi.UserService
	SessionService *backapi.SessionService
}

// NewServer returns a new instantiated Server with all the fields filled
// The logic and configuration itself is done only on the Setup step
func NewServer(
	UserService *backapi.UserService,
	SessionService *backapi.SessionService) *Server {
	return &Server{
		Router: echo.New(),
		Port:   viper.GetString("http.port"),
		Middlewares: &middlewares{
			UserService:    UserService,
			SessionService: SessionService,
		},
		UserService:    UserService,
		SessionService: SessionService,
	}
}

func (s *Server) setupMiddlewares() error {
	// SetContext middleware is a mandatory one that SHOULD be called in first
	// It stores the stdlib context in the echo one with values potentially required later
	// s.Router.Use(s.Middlewares.SetContext)

	// Handling of the logger configuration
	s.Router.Use(middleware.Logger())
	s.Router.Use(middleware.Recover())
	s.Router.Use(middleware.CORS())
	s.Router.Logger.SetLevel(log.DEBUG)

	// Handling of the RequestID generation
	s.Router.Use(middleware.RequestID())

	// Add custom middlewares if any here

	// Auth middleware will store the CurrentUser in the context if possible
	s.Router.Use(s.Middlewares.Auth)

	return nil
}

func (s *Server) Setup() error {
	// redirect all the queries to HTTPS
	// s.Router.Pre(middleware.HTTPSRedirect())
	s.Router.Pre(middleware.AddTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		Skipper: func(c echo.Context) bool {
			if c.Request().RequestURI == "/" {
				return false
			}
			return true
		},
	}))

	if err := s.setupMiddlewares(); err != nil {
		return err
	}

	// Defines Handlers
	userHandler := UserHandler{
		UserService: s.UserService,
	}

	sessionHandler := SessionHandler{
		SessionService: s.SessionService,
	}

	// Create a group unders /api/v1 using the Auth middleware
	// All the requests will then be potentially authenticated if the correct info are supplied
	apiV1Router := s.Router.Group("/api/v1", s.Middlewares.Auth)

	// Users related routes
	apiV1Router.POST("/users", userHandler.CreateUser)
	apiV1Router.GET("/users/me", userHandler.GetCUser, s.Middlewares.MustBeAuth)
	// apiV1Router.GET("/users", userHandler.ListUsers)
	// apiV1Router.POST("/users/details", userHandler.UserUpdateDetails, s.Middlewares.MustBeAuth)
	apiV1Router.POST("/users/me/change_password", userHandler.UserChangePassword, s.Middlewares.MustBeAuth)
	apiV1Router.POST("/users/change_email", userHandler.UserChangeEmail, s.Middlewares.MustBeAuth)
	// apiV1Router.GET("/users/:id", userHandler.GetUser)

	// Session related routes
	apiV1Router.POST("/session", sessionHandler.CreateSession)
	apiV1Router.DELETE("/session", sessionHandler.RemoveSession, s.Middlewares.MustBeAuth)

	return nil
}

func (s *Server) Start() error {
	if err := s.Router.Start(":" + s.Port); err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Router.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
