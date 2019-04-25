package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/tacusci/logging"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/acme/autocert"
)

var domainNameRegex = regexp.MustCompile(`([a-zA-Z]+(\.[a-zA-Z]{2,}){1,}|localhost(\:[0-9]+){0,})`)

type ProxyMapping struct {
	DomainContext string
	RequestURI    string
	TargetURL     string
}

func loadProxyMappings(reader io.Reader) ([]*ProxyMapping, error) {

	if reader == nil {
		return nil, errors.New("io.Reader instance is a nil pointer")
	}

	proxyMappings := []*ProxyMapping{}

	scanner := bufio.NewScanner(reader)
	configLineCount := 0
	for scanner.Scan() {
		configLineCount++
		configLine := scanner.Text()
		values := strings.Split(configLine, " ")
		if len(values) > 0 {
			proxyMapping := &ProxyMapping{}

			proxyDirectiveLookupIndex := 0
			if domainNameRegex.MatchString(strings.ToLower(values[0])) {
				proxyMapping.DomainContext = values[proxyDirectiveLookupIndex]
				proxyDirectiveLookupIndex++
			}
			if strings.ToLower(values[proxyDirectiveLookupIndex]) == "proxy" {
				if len(values) == proxyDirectiveLookupIndex+1 {
					return nil, fmt.Errorf("config line %d, missing request uri field for proxy mapping", configLineCount)
				}
				if len(values) > proxyDirectiveLookupIndex+2 {
					proxyMapping.RequestURI = values[proxyDirectiveLookupIndex+1]
					proxyMapping.TargetURL = values[proxyDirectiveLookupIndex+2]
					proxyMappings = append(proxyMappings, proxyMapping)
				} else {
					return nil, fmt.Errorf("config line %d, missing target url field for proxy mapping", configLineCount)
				}
			} else {
				return nil, fmt.Errorf("config line %d, unknown directive %s", configLineCount, values[0])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return proxyMappings, nil
}

type options struct {
	debug   bool
	autoSSL bool
}

func main() {

	opts := &options{}
	flag.BoolVar(&opts.debug, "dbg", false, "Set/turn on debug mode.")
	flag.BoolVar(&opts.autoSSL, "autossl", false, "Set/turn on auto ssl")

	flag.Parse()

	e := echo.New()

	if opts.debug {
		e.Logger.SetLevel(log.DEBUG)
	}
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	configFile, err := os.Open("tbl.config")
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer configFile.Close()

	proxyMappings, err := loadProxyMappings(configFile)
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
			e.Group(proxyMapping.RequestURI, middleware.ProxyWithConfig(middleware.ProxyConfig{
				Balancer: middleware.NewRandomBalancer(targets),
				Skipper: func(c echo.Context) bool {
					if proxyMapping.DomainContext == "" {
						return false
					}
					if strings.Contains(c.Request().Host, proxyMapping.DomainContext) {
						return false
					}
					return true
				},
			}))
			//e.Group(proxyMapping.RequestURI, middleware.Proxy(middleware.NewRandomBalancer(targets)))
		}

		e.Logger.Info("Started tinylb...")

		if opts.autoSSL {
			e.Use(middleware.HTTPSRedirect())
			e.Logger.Fatal(e.StartAutoTLS(":443"))
		} else {
			e.Start(":80")
		}
	}
}
