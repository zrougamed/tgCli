# TigerGraph CLI (tgCli)

A comprehensive command-line interface tool for managing TigerGraph Cloud instances and server operations. TigerGraph CLI provides seamless integration with TigerGraph Cloud services and local server management capabilities.

## Features

### Cloud Management
- **Authentication**: Secure login to TigerGraph Cloud with credential management
- **Instance Control**: Start, stop, terminate, and archive cloud instances
- **Monitoring**: List and monitor all cloud instances with status tracking
- **Configuration**: Save and manage cloud credentials

### Server Operations
- **GSQL Terminal**: Interactive GSQL command execution with session management
- **Service Management**: Start and stop TigerGraph services (GPE/GSE/RESTPP)
- **Backup Operations**: Create comprehensive backups of TigerGraph databases
- **Multi-Version Support**: Compatible with TigerGraph versions 3.0.0 through 3.6.2

### Configuration Management
- **Server Aliases**: Create and manage server connection profiles
- **Credential Storage**: Secure storage of authentication credentials
- **Default Settings**: Set default configurations for streamlined operations
- **Cross-Platform**: Support for Windows, macOS, and Linux

## Installation

### Prerequisites
- Go 1.24 or higher
- TigerGraph Cloud account (for cloud operations)
- TigerGraph server access (for server operations)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/zrougamed/tgCli.git
cd tgCli

# Build the binary
make build

# Install globally (optional)
make install
```

### Download Binary
Pre-built binaries are available for multiple platforms:
- Linux (AMD64/ARM64)
- macOS (Intel/Apple Silicon)
- Windows (AMD64/ARM64)

## Usage

### Quick Start

```bash
# Display help
tg --help

# Check version
tg version

# Login to TigerGraph Cloud
tg cloud login -e your@email.com -p yourpassword

# List cloud instances
tg cloud list

# Add server configuration
tg conf add -a myserver -u tigergraph -p mypassword --host http://localhost

# Start GSQL session
tg server gsql -a myserver
```

### Cloud Operations

```bash
# Login with interactive prompts
tg cloud login

# Login with credentials and save
tg cloud login -e user@domain.com -p password -s y

# List active instances only
tg cloud list -a y

# Start a cloud instance
tg cloud start -i INSTANCE_ID

# Stop a cloud instance
tg cloud stop -i INSTANCE_ID

# Terminate a cloud instance
tg cloud terminate -i INSTANCE_ID

# Archive a cloud instance
tg cloud archive -i INSTANCE_ID
```

### Server Management

```bash
# Connect to GSQL using server alias
tg server gsql -a myserver

# Connect to GSQL with direct credentials
tg server gsql -u username -p password --host http://server:14240

# Create database backup
tg server backup -a myserver -t ALL

# Backup schema only
tg server backup -a myserver -t SCHEMA

# Backup data only
tg server backup -a myserver -t DATA

# Start TigerGraph services
tg server services --ops start

# Stop TigerGraph services
tg server services --ops stop
```

### Configuration Management

```bash
# Add new server configuration
tg conf add -a production \
    -u tigergraph \
    -p mypassword \
    --host https://mycluster.i.tgcloud.io \
    --gsPort 14240 \
    --restPort 9000 \
    -d y

# List all configurations
tg conf list

# Delete server configuration
tg conf delete -a myserver

# Configure TigerGraph Cloud credentials
tg conf tgcloud -e user@domain.com -p password
```

## Configuration

TigerGraph CLI stores configuration in `~/.tgcli/config.yml` and credentials in `~/.tgcli/creds.bank`.

### Configuration Structure

```yaml
tgcloud:
  user: "your@email.com"
  password: "encrypted_password"

machines:
  production:
    host: "https://cluster.i.tgcloud.io"
    user: "tigergraph"
    password: "encrypted_password"
    gsPort: "14240"
    restPort: "9000"
  
  development:
    host: "http://localhost"
    user: "tigergraph"
    password: "tigergraph"
    gsPort: "14240"
    restPort: "9000"

default: "production"
```

## Command Reference

### Global Flags
- `-d, --debug`: Enable debug mode for verbose output

### Cloud Commands
- `tg cloud login`: Authenticate with TigerGraph Cloud
- `tg cloud list`: List all cloud instances
- `tg cloud start`: Start a cloud instance
- `tg cloud stop`: Stop a cloud instance
- `tg cloud terminate`: Terminate a cloud instance
- `tg cloud archive`: Archive a cloud instance

### Server Commands
- `tg server gsql`: Launch interactive GSQL terminal
- `tg server backup`: Create database backups
- `tg server services`: Manage TigerGraph services

### Configuration Commands
- `tg conf add`: Add server configuration
- `tg conf delete`: Remove server configuration
- `tg conf list`: Display all configurations
- `tg conf tgcloud`: Configure cloud credentials

## Development

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Create release packages
make package

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Run linter
make lint

# Show all available targets
make help
```

### Project Structure

```
tgCli/
├── cmd/
│   ├── main.go              # Application entry point
│   └── main_test.go         # Main application tests
├── internal/
│   ├── cloud/
│   │   ├── cloud.go         # Cloud operations
│   │   └── cloud_test.go    # Cloud operations tests
│   ├── config/
│   │   ├── config.go        # Configuration management
│   │   └── config_test.go   # Configuration tests
│   ├── helpers/
│   │   ├── helpers.go       # Utility functions
│   │   └── helpers_test.go  # Helper function tests
│   ├── models/
│   │   ├── models.go        # Data structures
│   │   └── models_test.go   # Model tests
│   └── server/
│       ├── server.go        # Server operations
│       └── server_test.go   # Server operation tests
├── pkg/
│   └── constants/
│       ├── constants.go     # Application constants
│       └── constants_test.go # Constants tests
├── LICENSE                  # Software license
├── Makefile                 # Build automation with comprehensive targets
├── README.md                # Project documentation
├── go.mod                   # Go module definition
└── go.sum                   # Dependency checksums
```

## Compatibility

### TigerGraph Versions
- TigerGraph 3.0.0 through 3.6.2
- TigerGraph Cloud (all current versions)

### Platforms
- Linux (AMD64, ARM64)
- macOS (Intel, Apple Silicon)
- Windows (AMD64, ARM64)

## Security

- Credentials are stored securely with appropriate file permissions
- Authentication tokens are managed automatically
- Password input uses secure terminal input methods
- Configuration files use restricted access permissions (0600)

## Troubleshooting

### Common Issues

**Authentication Failed**
```bash
# Re-login to refresh credentials
tg cloud login
```

**Connection Timeout**
```bash
# Check server connectivity and ports
tg server gsql --host http://your-server:14240
```

**Configuration Not Found**
```bash
# List available configurations
tg conf list

# Add new configuration
tg conf add -a myserver
```

## Contributing

This project uses a restrictive license. Please contact the author for contribution guidelines and licensing terms.

## License

This software is proprietary and confidential. Unauthorized copying, distribution, or use is strictly prohibited. See LICENSE file for complete terms.

## Support

- **Documentation**: [TigerGraph Documentation](https://docs.tigergraph.com)
- **Community**: [TigerGraph Community](https://community.tigergraph.com)
- **Discord**: [TigerGraph Discord](https://discord.gg/GkEmvDqB)

## Author

**Mohamed Zrouga**
- Website: [zrouga.email](https://zrouga.email) 
- Linkedin: [linkedin.com/in/zrouga-mohamed/](https://www.linkedin.com/in/zrouga-mohamed/)
- GitHub: [@zrougamed](https://github.com/zrougamed)
- Repository: [tgCli](https://github.com/zrougamed/tgCli)

## Version

Current Version: 0.1.1

---

Copyright (c) 2025 Mohamed Zrouga. All rights reserved.