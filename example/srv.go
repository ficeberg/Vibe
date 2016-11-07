package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/Festum/Vibe/wrappers"
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
	e.POST("/register", handler.Register)

	r := e.Group("/account")
	r.Use(middleware.JWTWithConfig(handler.JWTCheck()))
	r.GET("", handler.TokenResolve)
	r.POST("/update", handler.Update)
	r.GET("/info", handler.Get)
	r.GET("/info/:id", handler.Get)
	r.PUT("/:id", handler.Update)
	r.DELETE("", handler.Delete)

	e.Run(standard.New(":1323"))
}
