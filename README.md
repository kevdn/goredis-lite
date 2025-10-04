# GoRedis-Lite

A lightweight Redis-compatible in-memory key-value store written in Go, featuring I/O multiplexing for high-performance concurrent connections.

## Features

- **Redis Protocol Compatibility**: Supports Redis RESP protocol for seamless integration
- **I/O Multiplexing**: Uses epoll (Linux) and kqueue (macOS) for efficient event handling
- **Key Expiration**: Built-in TTL support with automatic key expiration
- **Concurrent Connections**: Handles up to 5000 concurrent connections
- **Cross-Platform**: Works on Linux and macOS

## Supported Commands

- `PING` - Test server connectivity
- `SET key value [EX seconds]` - Set a key-value pair with optional expiration
- `GET key` - Retrieve a value by key
- `TTL key` - Get time-to-live for a key
- `EXPIRE key seconds` - Set expiration time for an existing key
- `DEL key [key ...]` - Delete one or more keys
- `EXISTS key [key ...]` - Check if keys exist

## Quick Start

### Prerequisites

- Go 1.21 or later
- Linux or macOS

### Installation

1. Clone the repository:
```bash
git clone <your-repo-url>
cd goredis-lite
```

2. Build and run the server:
```bash
go run cmd/main.go
```

The server will start on port 3000.

### Testing with Redis CLI

To test the server, you can use the official Redis CLI:

1. Clone the Redis repository:
```bash
git clone https://github.com/redis/redis.git
cd redis
make
```

2. Start your GoRedis-Lite server:
```bash
go run cmd/main.go
```

3. In another terminal, connect using Redis CLI:
```bash
./src/redis-cli -p 3000
```

### Example Usage

```redis
127.0.0.1:3000> SET mykey "Hello World"
OK
127.0.0.1:3000> GET mykey
"Hello World"
127.0.0.1:3000> SET anotherkey "value" EX 60
OK
127.0.0.1:3000> TTL anotherkey
(integer) 60
127.0.0.1:3000> EXISTS mykey anotherkey
(integer) 2
127.0.0.1:3000> DEL mykey
(integer) 1
127.0.0.1:3000> GET mykey
(nil)
```

## Architecture

### I/O Multiplexing

The server uses platform-specific I/O multiplexing:
- **Linux**: epoll for efficient event handling
- **macOS**: kqueue for BSD-style event notification

### Key Components

- **Server**: TCP server with I/O multiplexing (`internal/server/`)
- **Core**: Command execution and RESP protocol handling (`internal/core/`)
- **Storage**: In-memory dictionary with expiration support (`internal/data_structure/`)
- **Config**: Server configuration (`internal/config/`)

### Expiration System

- Automatic key expiration using background cleanup
- TTL support with millisecond precision
- Configurable expiration frequency (100ms default)

## Configuration

Default configuration in `internal/config/config.go`:
- **Port**: 3000
- **Protocol**: TCP
- **Max Connections**: 20,000

## Development

### Project Structure

```
goredis-lite/
├── cmd/main.go                 # Application entry point
├── internal/
│   ├── config/                 # Configuration
│   ├── constant/               # Constants and timeouts
│   ├── core/                   # Core functionality
│   │   ├── executor.go         # Command execution
│   │   ├── resp.go             # RESP protocol
│   │   ├── expire.go           # Expiration logic
│   │   └── io_multiplexing/    # Platform-specific I/O
│   ├── data_structure/         # Storage implementation
│   └── server/                 # Server implementation
└── README.md
```

### Building

```bash
# Build the binary
go build -o goredis-lite cmd/main.go

# Run the binary
./goredis-lite
```

## Performance

- Supports up to 20,000 concurrent connections
- Efficient I/O multiplexing with 50ms timeout
- Automatic cleanup of expired keys every 100ms
- Memory-efficient key-value storage

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is open source. Please check the license file for details.

## Acknowledgments

Inspired by Redis and built with Go's excellent networking capabilities.