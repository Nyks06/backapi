package http

import (
	"context"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	webcore "github.com/nyks06/backapi"
)

const (
	// ContextKeyCurrentUser defines the name of the key in our echo.Context to get CurrentUser if any
	ContextKeyCurrentUser = "cUser_"
	// CookieSession defines the name of the cookie that is created for each user's session
	CookieSession = "sess_"
	// CookieFlash defines the name of the cookie that is created for each flash we want to display
	CookieFlash = "flash_"
)

type Server struct {
	Router *echo.Echo

	Middlewares *middlewares
	APIService  *webcore.APIService
}

func NewServer(APISe *webcore.APIService) *Server {
	return &Server{
		Router: echo.New(),

		Middlewares: &middlewares{
			APIService: APISe,
		},
		APIService: APISe,
	}
}

func (s *Server) SetupMiddlewares() error {
	// SetContext middleware is a mandatory one that SHOULD be called in first
	// It stores the stdlib context in the echo one with values potentially required later
	s.Router.Use(s.Middlewares.SetContext)

	// Handling of the logger configuration
	s.Router.Use(middleware.Logger())
	s.Router.Use(middleware.Recover())
	s.Router.Use(middleware.CORS())
	s.Router.Logger.SetLevel(log.DEBUG)

	// Handling of the RequestID generation
	// s.Router.Use(middleware.RequestID())

	// Add custom middlewares if any here

	// Auth middleware will store the CurrentUser in the stdlib context if possible
	// s.Router.Use(s.Middlewares.Auth)

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

	if err := s.SetupMiddlewares(); err != nil {
		return err
	}

	// Defines Handlers
	userHandler := UserHandler{
		APIService: s.APIService,
	}

	sessionHandler := SessionHandler{
		APIService: s.APIService,
	}

	// ticketHandler := TicketHandler{
	// 	APIService: s.APIService,
	// }

	// sportHandler := SportHandler{
	// 	APIService: s.APIService,
	// }

	// competitionHandler := CompetitionHandler{
	// 	APIService: s.APIService,
	// }

	// pronosticHandler := PronosticHandler{
	// 	APIService: s.APIService,
	// }

	// subscriptionHandler := SubscriptionHandler{
	// 	APIService: s.APIService,
	// }

	// contactHandler := ContactHandler{
	// 	APIService: s.APIService,
	// }

	//API Related routes
	apiV1Router := s.Router.Group("/api/v1")

	// Users related routes
	apiV1Router.POST("/users", userHandler.CreateUser)
	apiV1Router.GET("/users", userHandler.ListUsers)
	// apiV1Router.POST("/users/details", userHandler.UserUpdateDetails, s.Middlewares.MustBeAuth)
	// apiV1Router.POST("/users/contact", userHandler.UserUpdateContactSettings, s.Middlewares.MustBeAuth)
	// apiV1Router.POST("/users/change_password", userHandler.UserChangePassword, s.Middlewares.MustBeAuth)
	apiV1Router.GET("/users/me", userHandler.GetCUser)
	apiV1Router.GET("/users/:id", userHandler.GetUser)

	// Session related routes
	apiV1Router.POST("/login", sessionHandler.CreateSession)
	apiV1Router.POST("/logout", sessionHandler.RemoveSession, s.Middlewares.MustBeAuth)

	// apiV1Router.POST("/subscription", subscriptionHandler.Create, s.Middlewares.MustBeAuth)
	// apiV1Router.GET("/users/subscription", subscriptionHandler.Get, s.Middlewares.MustBeAuth)
	// apiV1Router.DELETE("/users/subscription/:id", subscriptionHandler.Cancel, s.Middlewares.MustBeAuth)

	// apiV1Router.POST("/ticket", ticketHandler.CreateTicket)
	// apiV1Router.GET("/ticket/:id", ticketHandler.GetTicket)
	// apiV1Router.DELETE("/ticket/:id", ticketHandler.DeleteTicket)
	// apiV1Router.POST("/ticket/:id", ticketHandler.UpdateTicket)
	// apiV1Router.GET("/ticket", ticketHandler.ListTicket)

	// apiV1Router.POST("/sport", sportHandler.CreateSport)
	// apiV1Router.GET("/sport/:id", sportHandler.GetSport)
	// apiV1Router.DELETE("/sport/:id", sportHandler.DeleteSport)
	// apiV1Router.POST("/sport/:id", ticketHandler.UpdateSport)
	// apiV1Router.GET("/sport", sportHandler.ListSport)

	// apiV1Router.POST("/competition", competitionHandler.CreateCompetition)
	// apiV1Router.GET("/competition/:id", competitionHandler.GetCompetition)
	// apiV1Router.DELETE("/competition/:id", competitionHandler.DeleteCompetition)
	// apiV1Router.POST("/competition/:id", ticketHandler.UpdateSport)
	// apiV1Router.GET("/competition", competitionHandler.ListCompetition)

	// apiV1Router.POST("/ticket", ticketHandler.CreateTicket)
	// apiV1Router.GET("/ticket/:id", ticketHandler.GetTicket)
	// apiV1Router.DELETE("/ticket/:id", ticketHandler.DeleteTicket)
	// apiV1Router.POST("/ticket/:id", ticketHandler.UpdateTicket)
	// apiV1Router.GET("/ticket", ticketHandler.ListTicket)

	// apiV1Router.POST("/pronostic", pronosticHandler.CreatePronostic)
	// apiV1Router.DELETE("/pronostic/:id", pronosticHandler.DeletePronostic)
	// apiV1Router.GET("/pronostic/:id", pronosticHandler.GetPronostic)
	// apiV1Router.GET("/pronostic", pronosticHandler.ListPronostic)
	// apiV1Router.POST("/pronostic/:id", pronosticHandler.UpdatePronostic)

	// apiV1Router.POST("/contact", contactHandler.SendMessage)

	return nil
}

func (s *Server) Start() error {

	// s.Router.AutoTLSManager.HostPolicy = autocert.HostWhitelist("www.yolirish.com", "yolirish.com")
	// s.Router.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
	// if err := s.Router.StartAutoTLS(":443"); err != nil {
	// 	return err
	// }

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	// Only for test purpose on localhost
	if err := s.Router.Start(":" + port); err != nil {
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
