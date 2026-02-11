# Taskfile Build System

## Overview

The Taskfile build system provides a modern, cross-platform alternative to the traditional Makefile for building, installing, testing, and developing the Perplexity CLI application. It uses Task (taskfile.dev) to define and manage build tasks with improved capabilities including task dependencies, incremental builds, and automatic install location detection.

## Architecture

### Taskfile Configuration

The build system is defined in `Taskfile.yaml` with the following components:

**Global Variables:**
- `BUILD_DIR`: Output directory for compiled binaries (default: `build/`)
- `BINARY_NAME`: Name of the compiled binary (default: `pplx`)
- `MAIN_FILE`: Entry point for the application (default: `main.go`)

**Available Tasks:**

- `default` - Lists all available tasks
- `build` - Compiles the Go application to `build/pplx`
- `install` - Installs the binary to GOPATH/bin or /usr/local/bin
- `clean` - Removes the build directory
- `test` - Runs all tests with verbose output
- `test-coverage` - Runs tests with coverage reporting
- `fmt` - Formats Go code using gofmt
- `lint` - Runs golangci-lint for static analysis
- `deps` - Downloads and organizes Go module dependencies
- `run` - Builds and executes the application
- `verify` - Runs build and test tasks in sequence
- `all` - Runs build, test, and install tasks

### Key Features

**Incremental Builds:**
Taskfile tracks source files (`**/*.go`) for intelligent incremental builds, only recompiling when source files change.

**Auto-Detection Install Path:**
The install task automatically detects and prioritizes:
1. GOPATH/bin (Go workspace binaries directory)
2. /usr/local/bin (system-wide location)

**Task Dependencies:**
Tasks can depend on other tasks, ensuring proper execution order. For example, `verify` depends on `build` and `test`.

**Cross-Platform Support:**
Taskfile works consistently across macOS, Linux, and Windows, eliminating platform-specific Makefile syntax issues.

## Usage

### Prerequisites

Install Task (taskfile.dev) from the official website: https://taskfile.dev/installation/

### Common Commands

```bash
# List all available tasks
task --list

# Build the application
task build

# Install the binary to system
task install

# Clean build artifacts
task clean

# Run tests
task test

# Run tests with coverage
task test-coverage

# Format code
task fmt

# Run linter
task lint

# Build and run the application
task run

# Verify build and tests
task verify

# Build, test, and install
task all
```

### Build Output

The build process creates a `build/` directory containing the compiled binary:
```
build/
└── pplx
```

## Coexistence with Makefile

The Taskfile is additive and co-exists with the existing Makefile. Both build systems provide similar functionality:
- Makefile: Traditional, Unix-centric approach
- Taskfile: Modern, cross-platform with enhanced features

Users can choose either build system based on their preference and development environment.

## Design Decisions

**Build Directory Preservation:** The `build/` directory convention is maintained to match existing project structure.

**Root Build Command:** Uses `go build -o {{.BUILD_DIR}}/{{.BINARY_NAME}} .` to build from the project root.

**Dependency Management:** Leverages Task's built-in dependency system for proper task execution ordering.

**Source Tracking:** Includes all `.go` files for comprehensive incremental build detection.

## Implementation Details

**Build Task:**
```yaml
build:
  sources:
    - "**/*.go"
  generates:
    - "{{.BUILD_DIR}}/{{.BINARY_NAME}}"
  cmds:
    - go build -o {{.BUILD_DIR}}/{{.BINARY_NAME}} .
```

**Install Task:**
```yaml
install: build
  cmds:
    - |
      if [ -n "$GOPATH" ]; then
        install -m 0755 {{.BUILD_DIR}}/{{.BINARY_NAME}} $GOPATH/bin/
      else
        install -m 0755 {{.BUILD_DIR}}/{{.BINARY_NAME}} /usr/local/bin/
      fi
```

## Testing and Verification

The Taskfile should be validated after creation:

1. Run `task --list` to verify all tasks are defined
2. Run `task build` to confirm binary creation
3. Run `task clean` to verify cleanup
4. Run `task install` to test installation detection
5. Run `task test` to ensure tests execute
6. Run `task verify` to confirm task dependencies
7. Run incremental builds to verify source tracking

## Maintenance

The Taskfile should be updated when:
- New build requirements emerge
- Additional testing or linting tools are added
- Project structure changes affect build paths
- New development workflows need automation
