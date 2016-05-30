package main

import (
	"./vibe/wrappers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	// "errors"
	"net/http"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "It's Vibe!")
	})
	handler := new(wrappers.Handlers)

	e.GET("/auth/:id", handler.Check)

	e.POST("/login", handler.Login)
	e.GET("/", handler.Accessible)
	e.POST("/register", handler.Register)

	r := e.Group("/account")
	r.Use(middleware.JWTAuthWithConfig(handler.JWTCheck()))
	r.GET("", handler.TokenResolve)
	r.POST("/update", handler.Update)
	r.GET("/info", handler.Get)
	r.GET("/info/:id", handler.Get)
	r.PUT("/:id", handler.Update)
	r.DELETE("", handler.Delete)

	// mux.HandleFunc("/auth/{provider}/callback", func(w http.ResponseWriter, r *http.Request) {
	// 	auth.Social(w, r)
	// })

	// mux.HandleFunc("/auth/{provider}", gothic.BeginAuthHandler)

	e.Run(standard.New(":1323"))
}
