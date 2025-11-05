# immich-compress

A CLI tool for compressing existing photos and videos in Immich by downloading, compressing, and re-uploading them with the same metadata and tags.

![CI Status](https://github.com/your-org/immich-compress/workflows/CI/badge.svg)
![Go Version](https://img.shields.io/badge/Go-1.23%2B-blue)
![License](https://img.shields.io/badge/license-MIT-green)

## 🚧 Work in Progress

This project is currently under active development. While the core infrastructure is in place and tests are passing, some features may still be incomplete or subject to change.

## 📋 Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Development](#development)
- [System Dependencies](#system-dependencies)
- [CI/CD](#cicd)
- [Contributing](#contributing)
- [License](#license)

## ✨ Features

- **Smart Compression**: Automatically compresses photos and videos while preserving quality
- **Metadata Preservation**: Maintains original metadata and tags during compression
- **Parallel Processing**: Configurable parallel processing for better performance
- **Time-based Filtering**: Option to compress only assets after a specific timestamp
- **Immich Integration**: Seamless integration with existing Immich instances

## 🔧 Prerequisites

### System Dependencies

Before installing immich-compress, ensure you have the following system dependencies installed:

#### Ubuntu/Debian
```bash
sudo apt-get install -y pkg-config libvips-dev libvips-tools
```

#### macOS
```bash
brew install vips pkg-config
```

#### Arch Linux
```bash
pacman -S pkgconfig libvips
```

#### Fedora
```bash
dnf install pkgconfig vips-devel
```

### Go Requirements
- Go 1.23 or later
- Git

## 📦 Installation

### From Source

1. Clone the repository:
```bash
git clone https://github.com/your-org/immich-compress.git
cd immich-compress
```

2. Install system dependencies (see [System Dependencies](#system-dependencies))

3. Build the application:
```bash
go build
```

4. Install globally (optional):
```bash
go install
```

### Binary Installation

Download the latest binary from the [Releases](https://github.com/your-org/immich-compress/releases) page and add it to your PATH.

## 🚀 Usage

### Basic Usage

```bash
immich-compress compress --server https://your-immich-server.com --api-key YOUR_API_KEY
```

### Command Options

```bash
# Compress with custom parallel processing
immich-compress compress --server https://your-immich-server.com --api-key YOUR_API_KEY --parallel 8

# Compress only assets after a specific date
immich-compress compress --server https://your-immich-server.com --api-key YOUR_API_KEY --after 2024-01-01

# Show help
immich-compress --help
immich-compress compress --help
```

### Command Line Options

#### Global Options
- `--parallel, -p int`: Number of parallel processes (default: number of CPU cores)
- `--after, -t time`: Only compress assets after this timestamp

#### Compress Command
- `--server, -s string`: **Required** - Immich server address
- `--api-key, -a string`: **Required** - Immich server API key

### Environment Variables

You can also use environment variables instead of command line flags:

```bash
export IMMICH_SERVER="https://your-immich-server.com"
export IMMICH_API_KEY="YOUR_API_KEY"
immich-compress compress
```

## 🛠️ Development

### Project Structure

```
immich-compress/
├── cmd/                    # CLI command definitions
│   ├── root.go            # Root command setup
│   └── compress.go        # Compress command implementation
├── compress/              # Core compression logic
├── immich/                # Auto-generated Immich API client
├── .github/workflows/     # GitHub Actions CI/CD
├── main.go                # Application entry point
├── go.mod                 # Go module definition
└── README.md              # This file
```

### Building and Testing

```bash
# Build the application
go build

# Run the application
go run .

# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestName ./...

# Clean dependencies
go mod tidy

# Download dependencies
go mod download
```

### Getting API Key from Immich

1. Open your Immich web interface
2. Navigate to Account Settings
3. Create a new API key
4. Copy the generated key

### Testing with Immich

For testing purposes, you can:
1. Use a test Immich instance
2. Create a separate API key with limited permissions
3. Test with a small number of assets first

## 🔧 System Dependencies

### libvips

The application uses [libvips](https://libvips.github.io/libvips/) for image processing through the [vipsgen](https://github.com/cshum/vipsgen) Go bindings.

**Current Dependencies:**
- `pkg-config`: Build tool for managing library compile flags
- `libvips-dev`: Development headers for libvips
- `libvips-tools`: Additional libvips utilities

**Compatibility Notes:**
- Currently uses `vipsgen v1.1.3` for compatibility with Ubuntu CI
- May require newer libvips versions for advanced features

## 🚦 CI/CD

### GitHub Actions Pipeline

The project uses GitHub Actions for continuous integration:

- **Trigger**: Push and Pull Requests to main/master/develop branches
- **Jobs**:
  - **test**: Runs tests with coverage, race detection, and builds
  - **lint**: Code quality checks with golangci-lint
- **Platform**: Ubuntu Latest
- **Go Version**: Go stable (latest)
- **System Dependencies**: Automatically installed in CI

### Status Badges

Add these badges to your repository README:

```markdown
![CI](https://github.com/your-org/immich-compress/workflows/CI/badge.svg)
```

## 🤝 Contributing

We welcome contributions! Please see our contributing guidelines:

### Development Setup

1. Fork the repository
2. Clone your fork
3. Install system dependencies
4. Create a feature branch
5. Make your changes
6. Run tests and linting
7. Submit a pull request

### Code Style

- Follow Go best practices
- Use `golangci-lint` for code quality
- Write tests for new features
- Update documentation as needed

### Testing

- Ensure all tests pass
- Add tests for new functionality
- Test on multiple platforms if possible

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⚠️ Important Notes

- **Work in Progress**: This project is actively developed and may have breaking changes
- **Backup Recommended**: Always backup your Immich data before running compression
- **API Key Security**: Keep your API keys secure and never commit them to version control
- **Resource Usage**: Compression can be resource-intensive; monitor your system resources

## 🐛 Known Issues

- Some advanced vipsgen features may not work with current libvips version
- Large video files may require significant processing time
- Network timeouts may occur with slow connections

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/your-org/immich-compress/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-org/immich-compress/discussions)
- **Wiki**: [Project Wiki](https://github.com/your-org/immich-compress/wiki)

---

**Note**: This project is not affiliated with the official Immich project. It's a community tool designed to work with Immich's API.
