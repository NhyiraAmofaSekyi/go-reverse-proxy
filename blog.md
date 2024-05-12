

# A Simple Reverse Proxy in Golang
A proxy is a device or software which acts as an agent or intimediary  between a host and a network, there are two main types of proxies forward and reverse.

### Forward proxy
this proxy is used route traffic between a client and a server, one of the main benefits of this is privacy, as the proxy sends requests in th behalf of the user,this can protect the users privacy.

### Reverse proxy
The rest  of the article will focus on reverse proxies which as the name suggests perform the reverse function of the forward proxy, it's function is to serve requests to clients on the behalf of the server and this comes with a number of benefits 

1. load balancing
2. caching 
3. rate limiting
4. security
5. SSL encryption
6. Live activity Monitoring and logging

### Differences
A common point of confusion between these two types of proxies is their position. In client server communication they both sit between the client and the server  the diffference being, forward proxies route behalf of the client while reverse proxies route messages to the client on behalf of the server.


## implementation 
below I implement a simple reverse proxy in golang 

Dependencies
- Docker
- Go 
- Make

implementation overview 
the rp creates a server at the port in the config file a new rp to handle requests for each resource 

## configuration 
below is an example configuration file 

```yaml
server:
  host: "localhost"
  listen_port: "8080"
resources:
  - name: Server1
    endpoint: /server1
    destination_url: "http://localhost:9001"
  - name: Server2
    endpoint: /server2
    destination_url: "http://localhost:9002"
  - name: Server3
    endpoint: /server3
    destination_url: "http://localhost:9003"
```
we will create a server at port 8080
we will then route traffic to each of the above endpoints using the servers endpoint

## Creating the proxy
`NewProxyV0` located in the proxyhandler.go file of the server package creates a proxy given a target url, and endpoint

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

		logIncomingRequest(proxyReq.In)
		proxyReq.Out.URL.Host = target.Host
		proxyReq.Out.URL.Scheme = target.Scheme
		proxyReq.Out.Header.Set("X-Forwarded-Host", proxyReq.In.Host)
		log.Print("orginal host: ", proxyReq.In.Host)
		proxyReq.Out.Host = target.Host

		originalPath := proxyReq.Out.URL.Path
		logOutgoingRequest(proxyReq.Out)
		log.Print("orginal path: ", originalPath)
		log.Print("new path: ", proxyReq.Out.URL.Path)
		proxyReq.Out.URL.Path = strings.TrimPrefix(originalPath, endpoint)
	}
	proxy.ModifyResponse = func(r *http.Response) error {
		r.Header.Set("Server", "rp")
		return nil
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
we then configure the `Rewrite` function that allows the request to be written as showcased by setting the target `host` and `X-forwarded-host`
`
### X-Forwarded-Host
The X-Forwarded-Host (XFH) header is used to preserve the original Host header sent by the client, before any changes made by the proxy.
This is particularly useful in situations where a reverse proxy serves multiple backends and needs to forward the original Host information to servers that might generate different responses based on the perceived URL.

Finally we configure the ModifyResponse function and overwrite the Server header sent by the server behind the proxy, additionally er could cache the response recieved by the proxy 



the ``ProxyRequestHandlerV0`` returns a request handler in which the proxies ``ServeHTTP`` method is called 

```go
func ProxyRequestHandlerV0(proxy *httputil.ReverseProxy, url *url.URL, endpoint string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now() // Capture the start time

		fmt.Printf("[ Reverse Proxy ] Request received at %s at %s\n", r.URL, time.Now().UTC())

		proxy.ServeHTTP(w, r)

		duration := time.Since(startTime) // Calculate the duration
		fmt.Printf("[ Reverse Proxy ] Redirected request to %s in %v at %s\n", r.URL, duration, time.Now().UTC())
	}
}
```
this is done so this handler can can measure the time taken to route the request before the handler is registered on the http request multiplexer 

in ther server.go file
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
 we create a server  as well as a proxy and proxy handler for each configured resource we then start the server 

## Testing

to test the server we will run 3 docker containers each containing a webserver with the make run command 

```bash
$ make run
```
``` makefile
## run: starts demo http services
.PHONY: run
run: run-containers	


run-containers:
	docker run --rm -d -p 9001:80 --name server1 kennethreitz/httpbin
	docker run --rm -d -p 9002:80 --name server2 kennethreitz/httpbin
	docker run --rm -d -p 9003:80 --name server3 kennethreitz/httpbin
```
we then buld and run the reverse proxy with the path to our config file
```bash
$ make run-proxy-server
```
```makefile
## run: starts demo http services
.PHONY: run-proxy-server
run-proxy-server:
	cd cmd && go build -o ../bin/rp && ../bin/rp run -c "../data/config.yaml"
```

now we can test the reverse proxy by sending a request to the endpoint /server1 which is the first server on our config list

request example 
```bash
$ curl -I  http://localhost:8080/server1
```
response
```bash
HTTP/1.1 200 OK
Access-Control-Allow-Credentials: true
Access-Control-Allow-Origin: *
Content-Length: 9593
Content-Type: text/html; charset=utf-8
Date: Sat, 11 May 2024 17:33:06 GMT
Server: rp
```


### creds 
---
hop by hop 
https://www.ory.sh/hop-by-hop-header-vulnerability-go-standard-library-reverse-proxy/

rp 
https://blog.joshsoftware.com/2021/05/25/simple-and-powerful-reverseproxy-in-go/

credited servers 
https://prabeshthapa.medium.com/learn-reverse-proxy-by-creating-one-yourself-using-go-87be2a29d1e


