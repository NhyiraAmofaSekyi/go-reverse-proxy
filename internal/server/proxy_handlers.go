package server

import (
	"fmt"
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
