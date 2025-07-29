# Golang RabbitMQ v2 - Clean Architecture

A production-ready Go application with RabbitMQ message processing, built using Clean Architecture principles.

## 🚀 Features

- **Clean Architecture** - Separation of concerns with domain, usecase, repository layers
- **RabbitMQ Integration** - Publisher confirms, consumer with retry mechanism, metrics monitoring
- **PostgreSQL** - Database with GORM ORM
- **Gin Framework** - High-performance HTTP router
- **Structured Logging** - Logrus with JSON formatting
- **Validation** - Request validation with go-playground/validator
- **File Logging** - Daily log rotation with Lumberjack
- **Docker Support** - Containerized deployment
- **Graceful Shutdown** - Proper resource cleanup
- **Health Checks** - Application and dependency health monitoring
- **Metrics** - Internal RabbitMQ metrics tracking

## 📁 Project Structure

```
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go               # Configuration management
│   ├── database/
│   │   └── postgres.go             # Database connection
│   ├── domain/
│   │   ├── user.go                 # User entity and interfaces
│   │   └── message.go              # Message entity and interfaces
│   ├── repository/
│   │   ├── user_repository.go      # User data access
│   │   └── message_repository.go   # Message data access
│   ├── usecase/
│   │   ├── user_usecase.go         # User business logic
│   │   └── message_usecase.go      # Message business logic
│   ├── handler/
│   │   ├── user_handler.go         # User HTTP handlers
│   │   └── message_handler.go      # Message HTTP handlers
│   └── middleware/
│       ├── cors.go                 # CORS middleware
│       ├── logging.go              # Request logging
│       └── recovery.go             # Panic recovery
├── pkg/
│   ├── logger/
│   │   └── logger.go               # Structured logging with file rotation
│   ├── rabbitmq/
│   │   ├── connection.go           # RabbitMQ connection management
│   │   ├── publisher.go            # Message publisher with confirms
│   │   ├── consumer.go             # Message consumer with worker pool
│   │   └── metrics.go              # Internal metrics tracking
│   └── validator/
│       └── validator.go            # Request validation
├── logs/                           # Log files directory
├── docker-compose.yml              # Docker services
├── Dockerfile                      # Application container
├── Makefile                        # Build automation and development commands
└── README.md                       # This file
```

## 🛠️ Technology Stack

- **Language**: Go 1.23
- **Web Framework**: Gin
- **Database**: PostgreSQL with GORM
- **Message Broker**: RabbitMQ
- **Logging**: Logrus with Lumberjack rotation
- **Validation**: go-playground/validator
- **Containerization**: Docker & Docker Compose

## 🚦 Getting Started

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- Make (optional)

### Quick Start with Docker

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd golang-rabbitmq-v2
   ```

2. **Start all services**
   ```bash
   make docker-up
   # or
   docker-compose up -d
   ```

3. **Check service status**
   ```bash
   make docker-ps
   # or
   docker-compose ps
   ```

4. **Test application health**
   ```bash
   make health
   # or
   curl http://localhost:8080/api/v1/monitoring/health
   ```

5. **Access RabbitMQ Management UI**
   ```bash
   make rabbitmq-ui
   # Opens: http://localhost:15672 (guest/guest)
   ```

### Local Development

1. **Quick setup for new developers**
   ```bash
   make quick-start
   ```

2. **Manual setup**
   ```bash
   # Setup environment
   make dev-setup
   
   # Install dependencies  
   make deps
   
   # Start PostgreSQL and RabbitMQ
   docker-compose up -d postgres rabbitmq
   
   # Run the application
   make run
   ```

3. **Development with live reload**
   ```bash
   # Install air for live reload
   make install-tools
   
   # Run with live reload
   make run-watch
   ```

## 📋 API Endpoints

### Users
- `POST /api/v1/users` - Create user
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user
- `GET /api/v1/users` - List users (with pagination)

### Messages
- `POST /api/v1/messages` - Create message
- `GET /api/v1/messages/:id` - Get message by ID
- `GET /api/v1/messages` - List messages (with filters)
- `POST /api/v1/messages/retry` - Retry failed messages
- `GET /api/v1/messages/stats` - Get message statistics

### Monitoring
- `GET /api/v1/monitoring/health` - Health check
- `GET /api/v1/monitoring/metrics` - RabbitMQ metrics

## 🔧 Configuration

Configuration is managed through environment variables. Copy `.env.example` to `.env` and adjust values:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=rabbitmq_app

# RabbitMQ Configuration
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_EXCHANGE=app_exchange
RABBITMQ_QUEUE=app_queue
RABBITMQ_ROUTING_KEY=app.message

# Server Configuration
SERVER_PORT=8080
GIN_MODE=debug

# Logging Configuration
LOG_LEVEL=info
LOG_FORMAT=json
LOG_FILE_ENABLED=true
LOG_DIR=logs
LOG_MAX_SIZE=100
LOG_MAX_BACKUPS=3
LOG_MAX_AGE=30
LOG_COMPRESS=true
```

## 📊 RabbitMQ Features

### Publisher Confirms
- Ensures message delivery to broker
- Timeout handling for confirmations
- Automatic retry on failures

### Consumer Features
- Worker pool for concurrent processing
- Graceful shutdown support
- Panic recovery
- Message acknowledgment/rejection
- Retry mechanism with backoff

### Metrics Monitoring
- Message publish/consume counters
- Processing time tracking
- Connection health monitoring
- Rate tracking (messages/second)
- Publisher confirm statistics

## 📝 Logging Features

### Daily Log Rotation
- **File Format**: `logs/app-YYYY-MM-DD.log` (e.g., `logs/app-2025-07-29.log`)
- **Auto Date Change**: New file created at midnight
- **Size-based Rotation**: Backup when file exceeds `LOG_MAX_SIZE`
- **Compression**: Old log files compressed with gzip
- **Cleanup**: Automatic deletion after `LOG_MAX_AGE` days

### Log Configuration
- **LOG_FILE_ENABLED**: Enable/disable file logging (default: true)
- **LOG_DIR**: Directory for log files (default: logs)
- **LOG_MAX_SIZE**: Max size per file in MB (default: 100)
- **LOG_MAX_BACKUPS**: Max backup files per day (default: 3)
- **LOG_MAX_AGE**: Max age in days (default: 30)
- **LOG_COMPRESS**: Compress old files (default: true)

### Log Structure Example
```
logs/
├── app-2025-07-29.log           # Today's log
├── app-2025-07-29.1.log         # Rotated backup (if size exceeded)
├── app-2025-07-28.log           # Yesterday's log
├── app-2025-07-27.log.gz        # Compressed old log
└── app-2025-07-26.log.gz        # Older compressed log
```

### Dual Output
- **Console**: Real-time logging to stdout
- **File**: Persistent logging with rotation
- **JSON Format**: Structured logging for production
- **Component Context**: Easy filtering by component

## 🏗️ Architecture Patterns

### Clean Architecture
- **Domain Layer**: Business entities and interfaces
- **Use Case Layer**: Business logic and orchestration
- **Repository Layer**: Data access abstraction
- **Handler Layer**: HTTP request/response handling

### Design Patterns Used
- Repository Pattern
- Dependency Injection
- Factory Pattern
- Observer Pattern (for metrics)

## 🐳 Docker Services

The `docker-compose.yml` includes:

- **PostgreSQL**: Database with persistent storage
- **RabbitMQ**: Message broker with management UI
- **App**: Go application with health checks and log volume mounting

Access RabbitMQ Management UI at: `http://localhost:15672` (guest/guest)

### Docker Log Management
- **Volume Mounting**: `app_logs:/app/logs` for persistent log storage
- **Log Access**: `docker-compose exec app ls -la logs/`
- **Log Viewing**: `docker logs rabbitmq_app --tail 20`
- **Follow Logs**: `docker logs rabbitmq_app -f`

## 🧪 Testing

### API Testing Examples

```bash
# Create a user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","phone":"1234567890"}'

# Create a message
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"content":"Hello RabbitMQ","type":"general"}'

# Get message statistics
curl http://localhost:8080/api/v1/messages/stats

# Get RabbitMQ metrics
curl http://localhost:8080/api/v1/monitoring/metrics

# Health check
curl http://localhost:8080/api/v1/monitoring/health
```

### Development Commands

```bash
# Show all available commands
make help

# Development setup
make dev-setup          # Setup development environment
make deps               # Install dependencies
make install-tools      # Install development tools

# Build and run
make build              # Build application
make run                # Run application
make run-watch          # Run with live reload

# Testing
make test               # Run tests
make test-coverage      # Run tests with coverage
make test-api           # Test API endpoints

# Code quality
make fmt                # Format code
make lint               # Lint code (requires golangci-lint)
make vet                # Run go vet

# Docker operations
make docker-up          # Start all services
make docker-down        # Stop all services
make docker-build       # Build Docker images
make docker-rebuild     # Rebuild and restart
make docker-logs        # View application logs
make docker-shell       # Access container shell

# Monitoring
make health             # Check application health
make metrics            # Show RabbitMQ metrics
make rabbitmq-ui        # Open RabbitMQ Management UI

# Cleanup
make clean              # Clean build artifacts
make clean-logs         # Clean log files
make docker-clean       # Clean Docker resources
```

## 📈 Monitoring & Observability

### Health Checks
- Database connectivity
- RabbitMQ connection status
- Application metrics

### Logging
- **Structured JSON logging** with Logrus
- **Daily file rotation** with automatic cleanup
- **Dual output**: Console + persistent file logging
- **Request/response logging** with detailed metrics
- **Error tracking** with stack traces
- **Component-based log contexts** for easy filtering

### Metrics
- Message processing rates
- Queue depths
- Error rates
- Connection statistics

## 🚀 Production Deployment

1. **Build and deploy**
   ```bash
   make docker-rebuild
   ```

2. **Scale application (if needed)**
   ```bash
   docker-compose up -d --scale app=3
   ```

3. **Monitor deployment**
   ```bash
   # Check container status
   make docker-ps
   
   # View logs
   make docker-logs
   
   # Check health
   make health
   
   # View metrics
   make metrics
   ```

4. **Useful production commands**
   ```bash
   make info               # Show project information
   make docker-shell       # Access container for debugging
   make db-reset           # Reset database if needed
   make rabbitmq-reset     # Reset RabbitMQ if needed
   ```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📝 License

This project is licensed under the MIT License.

## 🙋‍♂️ Support

For questions and support, please open an issue on GitHub.