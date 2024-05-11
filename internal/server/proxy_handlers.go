package server

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

func NewProxy(target *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Transport = &http.Transport{
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		MaxIdleConnsPerHost: 100,
	}
	return proxy
}
func ProxyRequestHandler(proxy *httputil.ReverseProxy, url *url.URL, endpoint string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		startTime := time.Now() // Capture the start time

		fmt.Printf("[ Reverse Proxy ] Request received at %s at %s\n", r.URL, time.Now().UTC())
		// Update the headers to allow for SSL redirection
		r.URL.Host = url.Host
		r.URL.Scheme = url.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = url.Host
		// Trim reverseProxyRoutePrefix from the path
		path := r.URL.Path
		r.URL.Path = strings.TrimLeft(path, endpoint)
		// Note that ServeHttp is non-blocking and uses a go routine under the hood
		proxy.ServeHTTP(w, r)

		duration := time.Since(startTime) // Calculate the duration
		fmt.Printf("[ Reverse Proxy ] Redirected request to %s in %v at %s\n", r.URL, duration, time.Now().UTC())
	}
}

func NewProxyV0(target *url.URL, endpoint string) *httputil.ReverseProxy {
	proxy := &httputil.ReverseProxy{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			MaxIdleConnsPerHost: 100,
		},
	}

	// Use the Rewrite function to safely modify the outgoing request
	proxy.Rewrite = func(proxyReq *httputil.ProxyRequest) {
		// Example: Setting the X-Example-Header on the outgoing request
		logIncomingRequest(proxyReq.In)
		proxyReq.Out.URL.Host = target.Host
		proxyReq.Out.URL.Scheme = target.Scheme
		proxyReq.Out.Header.Set("X-Forwarded-Host", proxyReq.In.Host)
		log.Print("orginal host: ", proxyReq.In.Host)
		proxyReq.Out.Host = target.Host

		// Trim the endpoint prefix if needed
		originalPath := proxyReq.Out.URL.Path
		logOutgoingRequest(proxyReq.Out)
		log.Print("orginal path: ", originalPath)
		log.Print("new path: ", proxyReq.Out.URL.Path)
		proxyReq.Out.URL.Path = strings.TrimPrefix(originalPath, endpoint)
	}
	return proxy
}

func ProxyRequestHandlerV0(proxy *httputil.ReverseProxy, url *url.URL, endpoint string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now() // Capture the start time

		fmt.Printf("[ Reverse Proxy ] Request received at %s at %s\n", r.URL, time.Now().UTC())

		// Note: Modifications moved to the Rewrite function of the ReverseProxy setup
		proxy.ServeHTTP(w, r)

		duration := time.Since(startTime) // Calculate the duration
		fmt.Printf("[ Reverse Proxy ] Redirected request to %s in %v at %s\n", r.URL, duration, time.Now().UTC())
	}
}

func logOutgoingRequest(req *http.Request) {
	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Printf("Failed to dump outgoing request: %v", err)
	} else {
		log.Printf("Modified outgoing request: %v", string(dump))
	}
}
func logIncomingRequest(req *http.Request) {
	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Printf("Failed to dump outgoing request: %v", err)
	} else {
		log.Printf("incoming request: %v", string(dump))
	}
}
