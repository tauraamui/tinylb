package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/tacusci/logging"

	"github.com/julienschmidt/httprouter"
)

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func main() {
	router := httprouter.New()
	origin, _ := url.Parse("http://localhost:9001")
	path := "/*catchall"

	reverseProxy := httputil.NewSingleHostReverseProxy(origin)

	reverseProxy.Director = func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = origin.Scheme
		req.URL.Host = origin.Host

		wildcardIndex := strings.IndexAny(path, "*")
		proxyPath := singleJoiningSlash(origin.Path, req.URL.Path[wildcardIndex:])
		if strings.HasSuffix(proxyPath, "/") && len(proxyPath) > 1 {
			proxyPath = proxyPath[:len(proxyPath)-1]
		}
		req.URL.Path = proxyPath
	}

	router.Handle("GET", path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		reverseProxy.ServeHTTP(w, r)
	})

	err := http.ListenAndServe(":80", router)
	if err != nil {
		logging.ErrorAndExit(err.Error())
	}
}
