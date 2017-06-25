package main

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/zengfu/web/whyengineer.com/api"
)

func main() {
	api.LaunchMqtt("tcp://0.0.0.0:1883")
	api.LaunchMqtt("ws://0.0.0.0:1884")
	defer api.CloseMqtt()
	e := echo.New()
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	e.Static("/", "static")
	e.GET("/", func(c echo.Context) error {

		return c.File("./template/index.html")
	})
	e.GET("/register", func(c echo.Context) error {

		return c.File("./template/register.html")
	})
	e.GET("/mqtt", func(c echo.Context) error {

		return c.File("./template/mqtt.html")
	})
	e.GET("/api/allmqtt", api.GetAllClients)
	e.GET("/api/checksession", api.CheckSession)
	e.GET("/api/cleansession", api.CLeanSession)
	e.GET("/api/checkusername", api.CheckUsername)
	e.POST("/api/signin", api.Signin)
	e.POST("/api/login", api.Login)
	e.Logger.Fatal(e.Start("0.0.0.0:8000"))
}
