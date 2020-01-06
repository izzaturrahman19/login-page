
package main

import (
	"github.com/Izzaturrahman19/login-page/controller"
	
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.POST("/register", controller.RegisterHandler)
	e.POST("/login", controller.LoginHandler)
	e.GET("/profile", controller.ProfileHandler)

	// Start server
	e.Logger.Fatal(e.Start(":7070"))
}