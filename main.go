package main

import (
	"net/url"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	deployDaemonURL, err := url.Parse("http://localhost:9001")
	if err != nil {
		e.Logger.Fatal(err)
	}

	targets := []*middleware.ProxyTarget{
		{URL: deployDaemonURL},
	}

	e.Group("/deploy", middleware.Proxy(middleware.NewRandomBalancer(targets)))

	e.Logger.Fatal(e.Start(":80"))
}
