# KPF - Kubernetes Port Forwarder

A Go-based CLI tool that automatically sets up and manages port forwarding from local ports to pods running in a Kubernetes cluster.

## Features

> **Key Advantage**: All port-forward processes are automatically terminated when the script stops — no lingering `kubectl port-forward` processes to clean up manually.

- **Automatic Process Cleanup**: Automatically terminates all port-forward processes on exit (SIGINT/SIGTERM)
- **Automatic Pod Discovery**: Finds pods in a Kubernetes namespace based on name prefix matching
- **Multiple Port Forward Rules**: Define multiple port forwarding rules in a config file
- **Graceful Shutdown**: Handles SIGINT/SIGTERM signals for clean shutdown
- **Built-in Config Editor**: Edit configuration using your preferred editor
- **Dual Output Logging**: Logs to both console and file

## Installation

```bash
git clone https://github.com/alex-campulungeanu/kpf
cd kpf
make build
```

The binary will be built at `./dist/kpf`.

## Configuration

KPF reads configuration from `~/.config/kpf/config.json`. On first run, it will create a default config file if one doesn't exist.

### Config File Format

```json
{
  "namespace": "your-namespace",
  "port_forward_rules": [
    {"prefix": "postgres-", "port": "5432"},
    {"prefix": "redis-", "port": "6379"},
    {"prefix": "api-", "port": "8080"}
  ]
}
```

| Field | Description |
|-------|-------------|
| `namespace` | The Kubernetes namespace to search for pods |
| `port_forward_rules` | Array of port forwarding rules |
| `prefix` | Pod name prefix to match (e.g., `postgres-` matches `postgres-deployment-abc123`) |
| `port` | Local and target port for forwarding |

### Editing Config

```bash
# Set your editor
export EDITOR=vim  # or nano, code, etc.

# Open config in editor
./dist/kpf -edit
```

## Usage

```bash
# Run port forwarding
./dist/kpf

# Edit configuration
./dist/kpf -edit
```

### Prerequisites

- Go 1.24.0 or later
- `kubectl` installed and configured
- `$KUBECONFIG` or `~/.kube/config` pointing to your cluster

## Project Structure

```
kpf/
├── main.go                  # Entry point
├── config/                  # Configuration management
├── helpers/                 # kubectl and port-forward helpers
├── dlogger/                 # Dual logger (console + file)
├── util/                    # Utility functions
├── bash_implementation/     # Bash script alternative
└── data/                    # Runtime data and logs
```

## Development

```bash
# Run tests
go test ./...

# Build
make build

```

## Bash Implementation

For quick scripts or systems without Go, use the bash version in `bash_implementation/`:

```bash
cp bash_implementation/.env.example bash_implementation/.env
# Edit .env with your settings
chmod +x bash_implementation/port_forwarding.sh
./bash_implementation/port_forwarding.sh
```
