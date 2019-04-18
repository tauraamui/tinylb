package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/tacusci/logging"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/acme/autocert"
)

type ProxyMapping struct {
	RequestURI string
	TargetURL  string
}

func loadProxyMappings() ([]*ProxyMapping, error) {
	configFile, err := os.Open("tbl.config")
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	proxyMappings := []*ProxyMapping{}

	scanner := bufio.NewScanner(configFile)
	configLineCount := 0
	for scanner.Scan() {
		configLineCount++
		configLine := scanner.Text()
		values := strings.Split(configLine, " ")
		if len(values) > 0 {
			if strings.ToLower(values[0]) == "proxy" {
				if len(values) == 1 {
					return nil, fmt.Errorf("config line %d, missing request uri field for proxy mapping", configLineCount)
				}
				if len(values) > 2 {
					proxyMappings = append(proxyMappings, &ProxyMapping{RequestURI: values[1], TargetURL: values[2]})
				} else {
					return nil, fmt.Errorf("config line %d, missing target url field for proxy mapping", configLineCount)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return proxyMappings, nil
}

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	proxyMappings, err := loadProxyMappings()
	if err != nil {
		e.Logger.Error(err)
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
			e.Logger.Debug(fmt.Sprintf("Mapping URI group: %s to target endpoint: %s", proxyMapping.RequestURI, proxyMapping.TargetURL))
			e.Group(proxyMapping.RequestURI, middleware.Proxy(middleware.NewRandomBalancer(targets)))
		}

		e.Logger.Info("Started tinylb...")

		e.Logger.Fatal(e.StartAutoTLS(":443"))
		//e.Start(":80")
	}
}
