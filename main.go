package main

import (
	"net/url"

	"github.com/tacusci/logging"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/crypto/acme/autocert"
)

type ProxyMapping struct {
	RequestURI string
	TargetURL  string
}

func main() {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	proxyMappings := []*ProxyMapping{
		&ProxyMapping{
			RequestURI: "/webhooks/push",
			TargetURL:  "http://localhost:9001",
		},
	}

	if len(proxyMappings) > 0 {
		targets := []*middleware.ProxyTarget{}
		e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
		for _, proxyMapping := range proxyMappings {
			url, err := url.Parse(proxyMapping.TargetURL)
			if err != nil {
				logging.ErrorAndExit(err.Error())
			}
			targets = append(targets, &middleware.ProxyTarget{URL: url})
			e.Group(proxyMapping.RequestURI, middleware.Proxy(middleware.NewRandomBalancer(targets)))
		}

		e.Logger.Fatal(e.StartAutoTLS(":443"))
		//e.Start(":80")
	}
}
