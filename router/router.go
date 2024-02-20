package router

import (
	"net/http"

	"github.com/labstack/echo-contrib/echoprometheus"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prokobit/auth-service/handler"
)

type Router struct {
	listenAddr string
}

func New(listenAddr string) *Router {
	return &Router{
		listenAddr: listenAddr,
	}
}

func (r *Router) Start() {
	e := echo.New()
	e.HideBanner = true
	e.Use(echojwt.WithConfig(echojwt.Config{
		Skipper:    skipHandler,
		SigningKey: []byte("secret"),
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(echoprometheus.NewMiddleware("app"))

	e.GET("/metrics", echoprometheus.NewHandler())
	e.GET("/", homeHandler)

	api := e.Group("/api")

	ah := handler.NewAccountHandler()
	accounts := api.Group("/accounts")
	accounts.POST("", ah.SignUp)
	accounts.POST("/login", ah.Login)

	e.Logger.Fatal(e.Start(r.listenAddr))
}

func skipHandler(c echo.Context) bool {
	if c.Request().Method == "POST" && (c.Path() == "/api/accounts" || c.Path() == "/api/accounts/login") {
		return true
	}
	return false
}

func homeHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
