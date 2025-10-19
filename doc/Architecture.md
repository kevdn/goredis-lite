# GoRedis-Lite Architecture Documentation

## Overview

GoRedis-Lite supports two distinct architectures for handling Redis-compatible operations:

1. **I/O Multiplexing Architecture** - Single-threaded, event-driven
2. **Share-Nothing Architecture** - Multi-worker, partitioned data storage

## Architecture Comparison

| Feature | I/O Multiplexing | Share-Nothing |
|---------|------------------|---------------|
| **Threading Model** | Single-threaded | Multi-threaded |
| **Data Storage** | Global shared storage | Partitioned per worker |
| **Scalability** | Vertical (CPU-bound) | Horizontal (worker-bound) |
| **Locking** | No locks needed | No locks needed |
| **Key Distribution** | All keys in one store | Consistent hashing |
| **Memory Usage** | Lower overhead | Higher overhead |
| **Latency** | Lower (no context switching) | Higher (worker dispatch) |
| **Throughput** | Limited by single CPU | Scales with CPU cores |

## 1. I/O Multiplexing Architecture

### Overview
Single-threaded server using platform-specific I/O multiplexing (epoll on Linux, kqueue on macOS) for efficient event handling.

### Key Components
- **Single Server Thread**: Handles all connections and commands
- **Global Storage**: One shared `dictStore` for all data
- **Event Loop**: Processes I/O events asynchronously
- **Active Expiration**: Background cleanup of expired keys

### Flow Diagram
See `IOMultiplexing_flow.puml` for detailed sequence diagram.

### Usage
```go
// In cmd/main.go
go server.RunIoMultiplexingServer(&wg)
```

### Benefits
- **Low Latency**: No context switching overhead
- **Simple**: Single-threaded, no concurrency issues
- **Memory Efficient**: Minimal overhead
- **Predictable**: Deterministic performance

### Limitations
- **CPU Bound**: Limited to single CPU core
- **Scalability**: Cannot scale beyond single thread performance

## 2. Share-Nothing Architecture

### Overview
Multi-worker architecture where each worker maintains isolated storage, providing horizontal scalability and consistent key partitioning.

### Key Components
- **Multiple Workers**: Each with isolated `dictStore`
- **I/O Handlers**: Round-robin connection assignment
- **Key Partitioning**: Consistent hashing for key distribution
- **Worker Isolation**: No shared state between workers

### Flow Diagram
See `SharedNothing_flow.puml` for detailed sequence diagram.

### Usage
```go
// In cmd/main.go
s := server.NewServer()
go s.StartSingleListener(&wg)  // or StartMultiListeners(&wg)
```

### Benefits
- **Horizontal Scalability**: Scales with CPU cores
- **No Locking**: Each worker operates independently
- **Consistent Partitioning**: Keys always go to same worker
- **Fault Isolation**: Worker failure doesn't affect others

### Limitations
- **Higher Latency**: Worker dispatch overhead
- **Memory Overhead**: Multiple storage instances
- **Complexity**: More complex architecture

## Configuration

### Worker Configuration
```go
// In internal/config/config.go
numCores := runtime.NumCPU()        // e.g., 8 cores
numIOHandlers := numCores / 2       // e.g., 4 I/O handlers
numWorkers := numCores / 2         // e.g., 4 workers
```

### Key Partitioning
```go
func (s *Server) getPartitionID(key string) int {
    hasher := fnv.New32a()
    hasher.Write([]byte(key))
    return int(hasher.Sum32()) % s.numWorkers
}
```

## Performance Considerations

### I/O Multiplexing
- **Best for**: Low-latency applications, single-core systems
- **CPU Usage**: Single-threaded, predictable
- **Memory**: Minimal overhead
- **Throughput**: Limited by single CPU performance

### Share-Nothing
- **Best for**: High-throughput applications, multi-core systems
- **CPU Usage**: Utilizes all available cores
- **Memory**: Higher overhead due to multiple workers
- **Throughput**: Scales linearly with CPU cores

## Switching Between Architectures

### Enable I/O Multiplexing
```go
// In cmd/main.go
go server.RunIoMultiplexingServer(&wg)
```

### Enable Share-Nothing
```go
// In cmd/main.go
s := server.NewServer()
go s.StartSingleListener(&wg)
```

## Monitoring and Profiling

### Profiling Setup
1. Install Graphviz: `brew install graphviz` (macOS) or `sudo apt-get install graphviz` (Linux)
2. Start server: `go run cmd/main.go`
3. Start profiling: `go tool pprof http://localhost:6060/debug/pprof/profile?seconds=20`
4. In pprof CLI: type `web` to see visualization

### Key Metrics to Monitor
- **Connection Count**: Active connections per I/O handler
- **Worker Load**: Tasks per worker per second
- **Memory Usage**: Per-worker memory consumption
- **Latency**: Command execution time
- **Throughput**: Commands per second

## Best Practices

### I/O Multiplexing
- Use for applications requiring low latency
- Monitor single-threaded CPU usage
- Consider connection pooling for high concurrency

### Share-Nothing
- Use for applications requiring high throughput
- Monitor worker load distribution
- Ensure even key distribution for optimal performance
- Consider worker-specific metrics and monitoring
