package main

import (
	"bufio"
	"net/url"
	"os"
	"regexp"

	"github.com/tacusci/logging"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/crypto/acme/autocert"
)

var proxyLineParseRegex = `\bproxy[[:blank:]](?P<URI>\/\S+)[[:blank:]](?P<scheme>http:\/\/|https:\/\/)(?P<hostname>\S+)(?P<portnumber>\:[0-9]+)$`

type ProxyMapping struct {
	RequestURI string
	TargetURL  string
}

func loadProxyMapsFromConfig(configFileName string) ([]*ProxyMapping, error) {
	configFile, err := os.Open("tlb.config")
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	regexEngine := regexp.MustCompile(proxyLineParseRegex)

	proxyMappings := []*ProxyMapping{}

	scanner := bufio.NewScanner(configFile)
	for scanner.Scan() {
		proxyMapping := &ProxyMapping{}
		namedMatches := regexEngine.FindStringSubmatch(scanner.Text())
		if namedMatches != nil {
			names := regexEngine.SubexpNames()
			for i, name := range names {
				if i != 0 {
					namedMatch := namedMatches[i]
					switch name {
					case "URI":
						proxyMapping.RequestURI = namedMatch
					case "scheme":
						proxyMapping.TargetURL += namedMatch
					case "hostname":
						proxyMapping.TargetURL += namedMatch
					case "portnumber":
						proxyMapping.TargetURL += namedMatch
					}
				}
			}
		}
	}
	return proxyMappings, nil
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	proxyMappings, err := loadProxyMapsFromConfig("tlb.config")
	if err != nil {
		logging.ErrorAndExit(err.Error())
	}

	if len(proxyMappings) > 0 {
		targets := []*middleware.ProxyTarget{}
		e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
		for _, proxyMapping := range proxyMappings {
			if url, err := url.Parse(proxyMapping.TargetURL); err == nil {
				targets = append(targets, &middleware.ProxyTarget{URL: url})
			}
			e.Group(proxyMapping.RequestURI, middleware.Proxy(middleware.NewRandomBalancer(targets)))
		}

		//e.Logger.Fatal(e.StartAutoTLS(":443"))
		e.Start(":80")
	}

	//targets := []*middleware.ProxyTarget{}

	/*
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
	*/
}
