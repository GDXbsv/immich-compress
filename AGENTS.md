# Agent Guidelines for immich-compress

## Build & Test Commands

- `go build` - Build the application
- `go run .` - Run directly without building
- `go test ./...` - Run all tests
- `go test -v ./...` - Run tests with verbose output
- `go test -run TestName ./...` - Run specific test
- `go mod tidy` - Clean up dependencies
- `go mod download` - Download dependencies

## System Dependencies

The application requires the following system dependencies:

### Ubuntu/Debian

```bash
sudo apt-get install -y pkg-config libvips-dev
```

### macOS

```bash
brew install vips pkg-config
```

### Arch Linux

```bash
pacman -S pkgconfig libvips
```

### Fedora

```bash
dnf install pkgconfig vips-devel
```

## Advanced Setup

For development with full libvips features (recommended for CI):

### Ubuntu/Debian (Full Feature Set)

```bash
sudo apt-get update
sudo apt-get install -y \
  meson ninja-build \
  libglib2.0-dev libexpat-dev librsvg2-dev libpng-dev \
  libjpeg-turbo8-dev libimagequant-dev libfftw3-dev \
  libpoppler-glib-dev libxml2-dev \
  libopenslide-dev libcfitsio-dev liborc-0.4-dev libpango1.0-dev \
  libtiff5-dev libgsf-1-dev giflib-tools libwebp-dev libheif-dev \
  libopenjp2-7-dev libcgif-dev \
  gobject-introspection libgirepository1.0-dev \
  libmagickwand-dev libmatio-dev libnifti2-dev \
  libjxl-dev libzip-dev libarchive-dev \
  pkg-config

# Create missing NIfTI pkg-config file
sudo mkdir -p /usr/local/lib/pkgconfig
sudo tee /usr/local/lib/pkgconfig/niftiio.pc > /dev/null <<EOF
prefix=/usr
exec_prefix=\${prefix}
libdir=\${prefix}/lib/x86_64-linux-gnu
includedir=\${prefix}/include/nifti

Name: libniftiio
Description: nifti library
Version: 3.0.1
Requires:
Cflags: -I\${includedir}
Libs: -L\${libdir} -lniftiio -lznz
EOF

# Build latest libvips from source
export VIPS_VERSION=8.17.3
wget https://github.com/libvips/libvips/releases/download/v$VIPS_VERSION/vips-$VIPS_VERSION.tar.xz
tar xf vips-$VIPS_VERSION.tar.xz
cd vips-$VIPS_VERSION

PKG_CONFIG_PATH=/usr/local/lib/pkgconfig:/usr/lib/pkgconfig meson setup _build \
  --buildtype=release --strip --prefix=/usr/local --libdir=lib \
  -Dmagick=enabled -Dopenslide=enabled -Dintrospection=enabled -Djpeg-xl=enabled

ninja -C _build
sudo ninja -C _build install
sudo ldconfig
```

### Environment Variables

```bash
export CGO_CFLAGS_ALLOW="-Xpreprocessor"
```

## Dependency Notes

- **vipsgen**: Currently using v1.2.1 (latest version)
- **libvips**: Built from source (v8.17.3) with full feature support including ImageMagick, OpenSlide, JPEG-XL
- **CI Environment**: Uses ubuntu-24.04 with comprehensive libvips dependencies
- **CGO**: Requires CGO_CFLAGS_ALLOW="-Xpreprocessor" environment variable

## CI/CD Pipeline

- **GitHub Actions**: `.github/workflows/ci.yml`
- **Triggered on**: Push and Pull Request to main/master/develop branches
- **Jobs**:
  - **test**: Runs tests with coverage, race detection, and builds
  - **lint**: Code quality checks with golangci-lint
- **Go Version**: Uses Go ^1.24
- **Setup Action**: actions/setup-go@v4
- **Ubuntu Version**: Uses `ubuntu-latest` runner
- **System Dependencies**:
  - Comprehensive libvips dependencies built from source
  - Meson, Ninja build system
  - All major image format libraries (JPEG, PNG, TIFF, WebP, HEIF, JPEG-XL, etc.)
- **Caching**: libvips build caching and Go module caching for faster builds

## Code Style Guidelines

- Use **camelCase** for variable/function names, **PascalCase** for exported functions/types
- Import groups: stdlib, then third-party (alphabetical), then local imports
- Use struct tags for CLI flags (Cobra framework)
- Return errors with `fmt.Errorf("context: %w", err)` pattern
- Always close HTTP response bodies with `defer resp.Body.Close()`
- Use context.Context for all API operations
- Name struct fields with **camelCase** for JSON marshaling
- Use meaningful error messages and handle all errors explicitly

## Project Structure

- `cmd/` - CLI command definitions
- `compress/` - Core compression logic
- `immich/` - Auto-generated Immich API client
- Entry point: `main.go` → `cmd.Execute()`

## Naming Conventions

- Package names: lowercase, no underscores
- CLI commands: lowercase with hyphens
- Flag variables: `flag<Command><Name>` struct pattern
- Error variables: `<Action>Err` pattern

## Error Handling

- Always handle errors explicitly with `if err != nil`
- Use context cancellation where appropriate
- Log important operations but avoid sensitive data
- Exit gracefully with `os.Exit(1)` on fatal errors
