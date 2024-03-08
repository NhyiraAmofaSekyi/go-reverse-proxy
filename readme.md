# Simple Reverse Proxy Application

This application serves as a simple reverse proxy, forwarding requests to different backend services based on the provided configuration.

## Features

- Forward HTTP requests to multiple backend services.
- Easy configuration via a YAML file.
- Docker integration for running demo HTTP services.
- Makefile for convenient project management.

## Prerequisites

- Docker
- Go 

## Getting Started

To get started with this application, clone the repository to your local machine:

### Running Demo HTTP Services

To start demo HTTP services on your local machine, run:

```bash
make run
```

### Stopping Demo HTTP Services

```bash
make stop
```

### Running the Reverse Proxy Server

To run the reverse proxy server with a specific configuration, ensure you have a config.yaml file in the data directory, then run:

```bash
make run-proxy-server
```

### Configuration
Your config.yaml should specify the backend services and any other configurations required by the reverse proxy. Here's an example structure:

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
### Help

```bash
make help
```
