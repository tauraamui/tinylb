package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/tacusci/logging"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/deploy", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		serveReverseProxy("http://localhost:9000/deploy", w, r)
	})
	router.GET("/health", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		logging.Info("Recieved health request")
		serveReverseProxy("http://localhost:9000/health", w, r)
	})

	http.ListenAndServe(":80", router)
}

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, w http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(target)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(w, req)
}
