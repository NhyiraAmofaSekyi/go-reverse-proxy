
---

# `rp` Reverse Proxy Tool Documentation

`rp` is a powerful CLI tool designed for setting up and managing reverse proxy configurations effortlessly. This documentation provides all the necessary information to get started with `rp`, including installation, usage examples, and troubleshooting.

## Getting Started

### Installation
Before you can use `rp`, ensure it is installed on your system. The installation method may vary based on your operating system and environment.

### Basic Commands

`rp` includes several commands and flags that help manage its operation. Here are the primary commands:

- `help`: Displays help about any command.
- `run`: Executes the reverse proxy server using a specified configuration.

### Configuration

`rp` relies on a configuration file to define various parameters for running the reverse proxy server. This configuration should be specified in a YAML file.

## Commands and Usage

### The `run` Command

The `run` command is used to start the reverse proxy server with configurations specified in a YAML file.

**Syntax:**

```bash
rp run --config /path/to/config.yaml
```

**Flags:**

- `--config, -c`: Specifies the path to the configuration file that `rp` should use to run the server.

**Example:**

```bash
rp run --config ./config/proxy_config.yaml
```

This command will start the reverse proxy server using the settings defined in `proxy_config.yaml`.

## Configuration File Example

Below is an example of what the configuration file might look like:

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

This configuration sets up a reverse proxy listening on localhost at port 8080, directing traffic to two backend services.

## Troubleshooting

If you encounter issues while using `rp`, ensure that:

- The configuration file path is correct.
- The configuration file syntax is correct for YAML.
- All required fields in the configuration are properly set.

For detailed errors, run the command with increased logging verbosity, if available, or check the system logs if the server fails to start.
