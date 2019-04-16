package main

import (
	"encoding/json"
	"io/ioutil"
	"net/url"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/crypto/acme/autocert"
)

type ProxyMapping struct {
	ReverseProxyGroup string   `json:"reverse-proxy-group"`
	TargetUrls        []string `json:"target-urls"`
}

func loadConfig(logger echo.Logger) []*ProxyMapping {
	file, err := ioutil.ReadFile("tlb.config")
	if err != nil {
		logger.Fatal("cannot read configuration file")
	}

	proxyMappings := []*ProxyMapping{}
	json.Unmarshal(file, &proxyMappings)

	return proxyMappings
}

func main() {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")

	targets := []*middleware.ProxyTarget{}
	for _, proxyMapping := range loadConfig(e.Logger) {
		for _, targetURL := range proxyMapping.TargetUrls {
			url, err := url.Parse(targetURL)
			if err != nil {
				e.Logger.Fatal(err)
			}
			targets = append(targets, &middleware.ProxyTarget{URL: url})
		}
		e.Logger.Info("Registering %s targets against group %s", targets, proxyMapping.ReverseProxyGroup)
		e.Group(proxyMapping.ReverseProxyGroup, middleware.Proxy(middleware.NewRandomBalancer(targets)))
	}

	e.Logger.Fatal(e.StartAutoTLS(":443"))
}
