# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go library for interacting with the Steam platform, providing a third-party API client for Steam authentication, points system, reactions, and Steam Guard functionality. The library is structured with a clean separation between the high-level client interface and lower-level data access objects.

## Key Development Commands

### Go Module Operations
```bash
# Build the project
go build .

# Run the main demo
go run main.go

# Test the library
go test ./...

# Run specific example
go run examples/basic/main.go
go run examples/advanced/main.go

# Build and run the test session (separate module)
cd test_session && go run session_demo.go
```

### Module Management
```bash
# Initialize/update dependencies
go mod tidy

# Verify dependencies
go mod verify

# Download dependencies
go mod download
```

## Architecture Overview

### Core Components

1. **Steam/client.go** - High-level client interface
   - Provides user-friendly API wrapper
   - Handles configuration and error management
   - Main entry point: `Steam.NewClient(config)`

2. **Steam/Dao/** - Data Access Layer
   - `dao.go` - HTTP client and request handling
   - `login.go` - Authentication and session management
   - `point.go` - Points system operations
   - `user.go` - User information management
   - `time.go` - Steam time synchronization

3. **Steam/Model/** - Data Models
   - Response structures for Steam API calls
   - Login and authentication models

4. **Steam/Protoc/** - Protocol Buffers
   - Steam API communication protocols
   - Generated from .proto files

5. **Steam/Utils/** - Utility Functions
   - Steam Guard token generation
   - Cryptographic helpers

### Key Patterns

- **Layered Architecture**: Clear separation between client interface, data access, and protocol handling
- **Error Handling**: Centralized error management through `Steam/Errors/`
- **Configuration**: Flexible client configuration with proxy support
- **Session Management**: Automatic cookie and token handling
- **Retry Logic**: Built-in retry mechanisms for network requests

### Authentication Flow

1. Get RSA public key from Steam (`getRSA`)
2. Encrypt password using RSA (`encryptPassword`)
3. Begin authentication session (`beginAuthSessionViaCredentials`)
4. Handle 2FA if required (Steam Guard codes)
5. Finalize login across multiple Steam domains (`finalizeLogin`)

### Important Files to Understand

- `Steam/client.go:85-119` - Main login implementation
- `Steam/Dao/login.go:634-654` - Core login logic
- `Steam/Dao/dao.go:184-239` - HTTP client setup with proxy support
- `examples/basic/main.go` - Simple usage example
- `examples/advanced/main.go` - Advanced interactive usage

### Testing

The project includes:
- Basic usage examples in `examples/`
- Interactive demo in `session_demo.go`
- Test module in `test_session/` (separate go.mod)

### Dependencies

Key external dependencies:
- `google.golang.org/protobuf` - Protocol buffer support
- `github.com/antchfx/htmlquery` - HTML parsing
- `golang.org/x/net` - Extended networking

### Security Considerations

- Passwords are RSA-encrypted before transmission
- Steam Guard integration for 2FA
- Cookie-based session management
- TLS configuration with `InsecureSkipVerify: true` (development only)

## Common Development Tasks

When modifying this codebase:

1. **Adding new Steam API endpoints**: Add protobuf definitions in `Steam/Protoc/`, implement in appropriate Dao files, expose through client interface
2. **Error handling**: Use the centralized error system in `Steam/Errors/`
3. **Authentication changes**: Modify login flow in `Steam/Dao/login.go`
4. **Testing**: Use the examples and test_session module for verification