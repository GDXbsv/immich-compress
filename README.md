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
- **Batch Limiting**: Option to limit the number of assets to process for testing or batch operations
- **UUID-based Selection**: Compress specific assets by their UUIDs for targeted operations
- **Image Quality Control**: Configurable image quality (1-100) with smart default of 80
- **Multiple Format Support**: Support for jpg, jpeg, jxl, webp, and heif image formats
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
# Basic compression of all assets
immich-compress compress --server https://your-immich-server.com --api-key YOUR_API_KEY

# Test compression with a limited number of assets (recommended for first run)
immich-compress compress --server https://your-immich-server.com --api-key YOUR_API_KEY --limit 50
```

### Command Options

```bash
# Compress with custom parallel processing
immich-compress compress --server https://your-immich-server.com --api-key YOUR_API_KEY --parallel 8

# Compress only assets after a specific date
immich-compress compress --server https://your-immich-server.com --api-key YOUR_API_KEY --after 2024-01-01

# Compress with a limited number of assets (useful for testing)
immich-compress compress --server https://your-immich-server.com --api-key YOUR_API_KEY --limit 100

# Compress specific assets by UUID
immich-compress compress --server https://your-immich-server.com --api-key YOUR_API_KEY --uuids "uuid1" --uuids "uuid2" --uuids "uuid3"

# Combine options - limited batch with parallel processing
immich-compress compress --server https://your-immich-server.com --api-key YOUR_API_KEY --parallel 4 --limit 50

# Compress with custom image quality
immich-compress compress --server https://your-immich-server.com --api-key YOUR_API_KEY --image-quality 90

# Compress to specific image format (WebP)
immich-compress compress --server https://your-immich-server.com --api-key YOUR_API_KEY --image-format webp

# Compress with high quality JPEG format
immich-compress compress --server https://your-immich-server.com --api-key YOUR_API_KEY --image-quality 95 --image-format jpeg

# Compress using modern JXL format with balanced quality
immich-compress compress --server https://your-immich-server.com --api-key YOUR_API_KEY --image-quality 80 --image-format jxl

# Batch process with format and quality settings
immich-compress compress --server https://your-immich-server.com --api-key YOUR_API_KEY --limit 100 --parallel 4 --image-quality 85 --image-format webp

# Show help
immich-compress --help
immich-compress compress --help
```

### Command Line Options

#### Global Options

- `--parallel, -p int`: Number of parallel processes (default: number of CPU cores)
- `--after, -t time`: Only compress assets after this timestamp
- `--limit, -l int`: Maximum number of assets to compress (default: 0 = no limit)

#### Compress Command

- `--server, -s string`: **Required** - Immich server address
- `--api-key, -a string`: **Required** - Immich server API key
- `--type, -i string`: Asset type to compress (IMAGE, VIDEO, ALL) (default: ALL)
- `--uuids, -u string`: Assets UUIDs (array)
- `--image-quality, -q int`: Image quality for compression (1-100) (default: 80)
- `--image-format, -f string`: Image format for compression (jpg, jpeg, jxl, webp, heif) (default: jpg)
- `--video-quality, -Q int`: Video quality for compression (1-100) Lower is higher quality (default: 25)
- `--video-format, -F string`: Video format for compression (av1, hevc) (default: av1)
- `--video-container, -c string`: Video container format (mkv, mp4) (default: mkv)
- `--diff-percents, -D int`: If size diff is lower than this percent files will not be replaced with new (default: 8)

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
3. Test with a small number of assets first using the `--limit` flag:
   ```bash
   # Test with only 10 assets
   immich-compress compress --server https://your-immich-server.com --api-key YOUR_API_KEY --limit 10
   ```
4. Combine with parallel processing for different test scenarios:
   ```bash
   # Test with limited assets and reduced parallelism
   immich-compress compress --server https://your-immich-server.com --api-key YOUR_API_KEY --limit 5 --parallel 2
   ```

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

## 💡 Best Practices

### Testing and Validation

- **Start Small**: Always test with a limited number of assets first using the `--limit` flag:

  ```bash
  # Test with only 10 assets
  immich-compress compress --server https://your-immich-server.com --api-key YOUR_API_KEY --limit 10
  ```

- **Gradual Scaling**: Once comfortable, gradually increase the limit and parallel settings:
  ```bash
  # Process 100 assets with moderate parallelism
  immich-compress compress --server https://your-immich-server.com --api-key YOUR_API_KEY --limit 100 --parallel 4
  ```

### Performance Optimization

- **Resource Monitoring**: Monitor your system resources (CPU, memory, network) during compression
- **Parallel Processing**: Adjust `--parallel` based on your system's capabilities (start conservative)
- **Time-based Filtering**: Use `--after` to process only recent assets for initial runs

### Image Compression Settings

- **Quality vs Size**: Balance image quality with file size using `--image-quality`:
  - Quality 70-80: Good balance for most use cases (default: 80)
  - Quality 90+: Minimal size reduction, better quality preservation
  - Quality 50-70: Smaller files, noticeable quality loss

- **Format Selection**: Choose optimal image formats with `--image-format`:
  - **jpg/jpeg**: Best compatibility, moderate compression
  - **webp**: Modern format, excellent compression, good quality
  - **jxl**: Cutting-edge format, superior compression, emerging support
  - **heif**: Apple ecosystem, good compression, limited browser support

- **Recommended Settings**:

  ```bash
  # For general use - WebP with balanced quality
  immich-compress compress --server ... --image-format webp --image-quality 80

  # For maximum compatibility - JPEG with high quality
  immich-compress compress --server ... --image-format jpeg --image-quality 90

  # For archival/quality preservation
  immich-compress compress --server ... --image-format jxl --image-quality 95
  ```

### Batch Operations

- **Large Batches**: For processing large numbers of assets:

  ```bash
  # Process assets in batches of 500
  immich-compress compress --server https://your-immich-server.com --api-key YOUR_API_KEY --limit 500 --parallel 8
  ```

- **Incremental Processing**: Run multiple times with different limits to process in chunks

### Production Deployment

- **Backup First**: Always backup your Immich data before running compression
- **API Keys**: Use dedicated API keys with minimal permissions for compression tasks
- **Resource Planning**: Ensure adequate disk space for temporary files during compression

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
