hop by hop 
https://www.ory.sh/hop-by-hop-header-vulnerability-go-standard-library-reverse-proxy/

rp 
https://blog.joshsoftware.com/2021/05/25/simple-and-powerful-reverseproxy-in-go/

credited servers 
https://prabeshthapa.medium.com/learn-reverse-proxy-by-creating-one-yourself-using-go-87be2a29d1e


define proxy foward/reverse
benefits of each

differentiate from forward proxy 

## explanation of rp
X-Forwarded-For
The X-Forwarded-For (XFF) header is used to identify the originating IP address of a client connecting to a web server through an HTTP proxy or a load balancer.
This header helps maintain a record of the original client’s IP address, which is useful for security audits, logging, and geo-localization, among other purposes.

### X-Forwarded-Host
The X-Forwarded-Host (XFH) header is used to preserve the original Host header sent by the client, before any changes made by the proxy.
This is particularly useful in situations where a reverse proxy serves multiple backends and needs to forward the original Host information to servers that might generate different responses based on the perceived URL.

#### uses
1. Preserving the Original Host Information
Context: When a request is proxied, the Host header often gets changed to the host of the proxy or the next hop, not the original host that the client intended to contact.
Use: The X-Forwarded-Host header can preserve the client's intended host, enabling the server handling the request at the end of the proxy chain to know the original host.
2. Generating Accurate Redirects
Context: Web applications often generate redirects based on the host name from the incoming request's Host header.
Use: When behind a proxy, using the X-Forwarded-Host allows the application to generate accurate redirects that reflect the client's perspective, rather than redirecting to URLs based on the proxy's host.
3. Constructing Links for Responses
Context: Dynamically generated web pages often include links that should reflect the URL as seen by the user.
Use: By utilizing the X-Forwarded-Host, web applications can construct links that are valid and meaningful to the end user, even when the application is behind one or more proxies.
4. Security and Host Header Attacks
Context: The host header is a common vector for attacks such as cache poisoning and web application routing issues.
Use: While X-Forwarded-Host can be useful, it must be handled securely. Servers need to be configured to trust this header only when it's known to come from a trusted proxy, to prevent spoofing and related security issues.
5. Compliance with Virtual Hosting
Context: In environments where multiple applications or services are hosted (virtual hosting), the destination service is often determined by the Host header.
Use: The X-Forwarded-Host helps maintain fidelity of this information across proxy layers, ensuring requests are routed internally to the correct service based on the original client request.





implementation overview 
the rp creates a server at rhe port in the config file a new rp to handle requests for each resource 

`NewProxyV0` located in the proxyhandler.go file of the server package creates a proxy givena target url, and endpoint

```go
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
```

we configure transport of the proxy  because 

```go
// Transport is an implementation of [RoundTripper] that supports HTTP,
// HTTPS, and HTTP proxies (for either HTTP or HTTPS with CONNECT).

// By default, Transport caches connections for future re-use.
// This may leave many open connections when accessing many hosts.
// This behavior can be managed using [Transport.CloseIdleConnections] method
// and the [Transport.MaxIdleConnsPerHost] and [Transport.DisableKeepAlives] fields.
```
we then configure the Rewrite function that allows the request to be written as showcased by setting the target host andX forwarded host

the ``ProxyRequestHandlerV0`` returns a request handler in whcih the proxy ``ServeHTTP`` is called

```go
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
```
this is done so this handler can be  registered on the http request multiplexer 

in ther server.go
```go
mux := http.NewServeMux()

	// Register the health check endpoint.
	mux.HandleFunc("/ping", ping)

	for _, resource := range cfg.Resources {
		destURL, err := url.Parse(resource.Destination_URL)
		if err != nil {
			log.Printf("Error parsing URL '%s': %v", resource.Destination_URL, err)
			continue // Skip this resource if the URL is invalid
		}

		proxy := NewProxyV0(destURL, resource.Endpoint)

		// Register the handler using ProxyRequestHandlerV0, passing the created proxy, destination URL, and endpoint
		mux.HandleFunc(resource.Endpoint, ProxyRequestHandlerV0(proxy, destURL, resource.Endpoint))
	}

	// Initialize the HTTP server.
	server = &http.Server{
		Addr:    cfg.Server.Host + ":" + cfg.Server.Listen_port,
		Handler: mux,
	}

	fmt.Printf("Server configured to listen on %s\n", server.Addr)
```
 file we create a server  as well as a proxy and proxy handler for each condigured resource we then start the server 





dependencies to run and test 

explain how requets headers and resp headers are set 

bereak down hop by hop 

end with a test



