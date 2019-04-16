package main

import (
	"net/url"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")

	deployDaemonURL, err := url.Parse("http://localhost:9001")
	if err != nil {
		e.Logger.Fatal(err)
	}

	targets := []*middleware.ProxyTarget{
		{URL: deployDaemonURL},
	}

	e.Group("/deploy", middleware.Proxy(middleware.NewRandomBalancer(targets)))

	e.Logger.Fatal(e.StartAutoTLS(":443"))
}
