# Feature Overview

This document provides an overview of the Perplexity CLI system and describes all auxiliary feature files in the `docs/features/` directory.

## System Description

The Perplexity CLI is a Go-based command-line interface tool that provides seamless interaction with the Perplexity Sonar API. The application uses the Cobra framework for CLI commands and Viper for configuration management.

Key capabilities include:
- One-shot queries via `pplx run` command (without session persistence)
- Interactive conversation mode with REPL interface (with colored separators)
- Session persistence in `~/.pplx/sessions/` with JSON-based storage (interactive mode only)
- Session listing, search, and show functionality with short Base62 IDs
- Citation parsing and formatted reference display
- Optional configuration file support (works with environment variables only)
- Cross-platform build system with Taskfile and Makefile support

## Feature Files

The following documentation files describe specific features of the system:

- `taskfile-build-system.md` - Documents the Taskfile.yaml implementation, a modern alternative to the traditional Makefile that provides improved cross-platform support, incremental builds, and task dependencies for building, installing, testing, and developing the CLI application.

- `perplexity-cli.md` - Comprehensive documentation of the core Perplexity CLI implementation, including the API client architecture, session management system, interactive REPL interface with colored separators, citation parsing, optional configuration file handling, and all available commands for querying, searching, and managing conversations.