# Obvious Service Examples

Examples showcasing the [server library](https://github.com/go-obvious/server) features and best practices.

## Examples

### 🚀 Basic HTTP Server (`cmd/hello/`)
A simple example showcasing basic API implementation with the server library.

```sh
cd cmd/hello
make build
make run
```

**Features Demonstrated:**
- Basic API registration
- HTTP routing with chi
- JSON responses
- Health checks and version endpoints

### 🔄 Lifecycle Management (`cmd/lifecycle/`)
Comprehensive example demonstrating advanced server features including configuration registry, lifecycle management, and graceful shutdown.

```sh
cd cmd/lifecycle  
make build
make run
```

**Features Demonstrated:**
- **Configuration Registry** - Self-registering configuration with validation
- **API Lifecycle Management** - Graceful startup and shutdown hooks
- **Enhanced Error Context** - Correlation ID tracking and structured errors
- **Database Service** - Mock database with connection management
- **Background Worker** - Job processing with lifecycle hooks
- **Graceful Shutdown** - SIGTERM/SIGINT handling with resource cleanup

**Environment Configuration:**
```sh
# Database settings
DATABASE_URL=mock://localhost:5432/testdb
DATABASE_MAX_CONNECTIONS=10
DATABASE_CONNECT_TIMEOUT=5s

# Worker settings  
WORKER_INTERVAL=30s
WORKER_MAX_JOBS=100
WORKER_ENABLE_PROCESSING=true
```

**API Endpoints:**
- `GET /api/database/users` - List users
- `POST /api/database/users` - Create user
- `GET /api/database/health` - Database health check
- `GET /api/worker/status` - Worker status and metrics
- `POST /api/worker/jobs` - Create background job
- `GET /api/worker/jobs` - List jobs

**Testing the APIs:**
```sh
make test-apis
```

**Custom Configuration:**
```sh
make run-with-config
```

## Development

### Running Examples

Each example includes its own Makefile with common targets:

```sh
make build      # Build the example
make run        # Run the server  
make clean      # Remove build artifacts
make help       # Show available targets
```

### Project Structure

```
examples/
├── cmd/
│   ├── hello/          # Basic HTTP server example
│   └── lifecycle/      # Advanced lifecycle example
└── internal/
    ├── build/          # Build information
    └── service/        # Service implementations
        ├── hello/      # Basic service
        ├── database/   # Database service with lifecycle
        └── worker/     # Background worker service
```

### Quick Start

1. **Basic Example:**
   ```sh
   cd cmd/hello && make run
   ```

2. **Advanced Example:**
   ```sh
   cd cmd/lifecycle && make run
   ```

3. **Test APIs:**
   ```sh
   # In another terminal
   curl http://localhost:8080/version
   curl http://localhost:8080/healthz
   ```

### Graceful Shutdown Demo

Run the lifecycle example and press `Ctrl+C` to see graceful shutdown in action:

```sh
cd cmd/lifecycle && make run
# Press Ctrl+C to trigger graceful shutdown
```

You'll see:
1. Configuration loading and validation
2. Database connection establishment  
3. Background worker startup
4. Server accepting requests
5. Graceful shutdown on SIGTERM
6. Background worker stopping
7. Database connections closing
8. Existing requests completing